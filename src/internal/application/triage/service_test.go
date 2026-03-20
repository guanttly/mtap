package triage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/euler/mtap/internal/domain/triage"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// ============================================================
// Mock Repositories
// ============================================================

type mockCheckInRepo struct {
	items map[string]*domain.CheckIn // appointmentID → CheckIn
}

func newMockCheckInRepo() *mockCheckInRepo {
	return &mockCheckInRepo{items: make(map[string]*domain.CheckIn)}
}

func (m *mockCheckInRepo) Save(_ context.Context, c *domain.CheckIn) error {
	m.items[c.AppointmentID] = c
	return nil
}
func (m *mockCheckInRepo) FindByAppointmentID(_ context.Context, appointmentID string) (*domain.CheckIn, error) {
	return m.items[appointmentID], nil
}
func (m *mockCheckInRepo) FindByPatientAndDate(_ context.Context, _, _ string) ([]*domain.CheckIn, error) {
	return nil, nil
}

type mockQueueRepo struct {
	queues map[string]*domain.WaitingQueue // roomID → queue
}

func newMockQueueRepo() *mockQueueRepo {
	return &mockQueueRepo{queues: make(map[string]*domain.WaitingQueue)}
}

func (m *mockQueueRepo) Save(_ context.Context, q *domain.WaitingQueue) error {
	m.queues[q.RoomID] = q
	return nil
}
func (m *mockQueueRepo) Update(_ context.Context, q *domain.WaitingQueue) error {
	m.queues[q.RoomID] = q
	return nil
}
func (m *mockQueueRepo) FindByRoomID(_ context.Context, roomID string) (*domain.WaitingQueue, error) {
	return m.queues[roomID], nil
}
func (m *mockQueueRepo) FindOrCreateByRoom(_ context.Context, roomID, deviceID, departmentID string) (*domain.WaitingQueue, error) {
	if q, ok := m.queues[roomID]; ok {
		return q, nil
	}
	q := domain.NewWaitingQueue(roomID, deviceID, departmentID)
	m.queues[roomID] = q
	return q, nil
}

type mockQueueEntryRepo struct {
	items map[string]*domain.QueueEntry
}

func newMockQueueEntryRepo() *mockQueueEntryRepo {
	return &mockQueueEntryRepo{items: make(map[string]*domain.QueueEntry)}
}

func (m *mockQueueEntryRepo) Save(_ context.Context, e *domain.QueueEntry) error {
	m.items[e.ID] = e
	return nil
}
func (m *mockQueueEntryRepo) Update(_ context.Context, e *domain.QueueEntry) error {
	m.items[e.ID] = e
	return nil
}
func (m *mockQueueEntryRepo) FindByQueueID(_ context.Context, queueID string, status domain.EntryStatus) ([]*domain.QueueEntry, error) {
	var result []*domain.QueueEntry
	for _, e := range m.items {
		if e.QueueID == queueID && (status == "" || e.Status == status) {
			result = append(result, e)
		}
	}
	return result, nil
}

type mockExecRepo struct {
	items map[string]*domain.ExamExecution // appointmentItemID → exec
}

func newMockExecRepo() *mockExecRepo {
	return &mockExecRepo{items: make(map[string]*domain.ExamExecution)}
}

func (m *mockExecRepo) Save(_ context.Context, e *domain.ExamExecution) error {
	m.items[e.AppointmentItemID] = e
	return nil
}
func (m *mockExecRepo) Update(_ context.Context, e *domain.ExamExecution) error {
	m.items[e.AppointmentItemID] = e
	return nil
}
func (m *mockExecRepo) FindByID(_ context.Context, id string) (*domain.ExamExecution, error) {
	for _, e := range m.items {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, nil
}
func (m *mockExecRepo) FindByAppointmentItemID(_ context.Context, itemID string) (*domain.ExamExecution, error) {
	return m.items[itemID], nil
}
func (m *mockExecRepo) FindByDevice(_ context.Context, _ string) ([]*domain.ExamExecution, error) {
	return nil, nil
}

// ============================================================
// Helper
// ============================================================

func newTestTriageService() (*TriageAppService, *mockExecRepo) {
	execRepo := newMockExecRepo()
	svc := NewTriageAppService(
		newMockCheckInRepo(),
		newMockQueueRepo(),
		newMockQueueEntryRepo(),
		execRepo,
	)
	return svc, execRepo
}

// seedExamExecution 在仓储中预置一个 checked_in 状态的检查执行记录
func seedExamExecution(execRepo *mockExecRepo, itemID, patientID, deviceID string) {
	exec := domain.NewExamExecution(itemID, patientID, deviceID)
	execRepo.items[itemID] = exec
}

// ============================================================
// NurseCheckIn Tests
// ============================================================

func TestNurseCheckIn_OK(t *testing.T) {
	svc, _ := newTestTriageService()
	resp, err := svc.NurseCheckIn(context.Background(), NurseCheckInReq{
		AppointmentID: "APPT001",
		Remark:        "紧急签到",
	}, "nurse01")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.CheckInID)
	assert.Equal(t, 1, resp.QueueNumber)
}

func TestNurseCheckIn_DuplicateRejected(t *testing.T) {
	svc, _ := newTestTriageService()
	_, err := svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT002"}, "nurse01")
	require.NoError(t, err)

	// 第二次签到同一预约应被拒绝
	_, err = svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT002"}, "nurse01")
	assert.True(t, bizErr.Is(err, bizErr.ErrTriageAlreadyCheckedIn))
}

func TestNurseCheckIn_QueueNumberIncremental(t *testing.T) {
	svc, _ := newTestTriageService()

	resp1, _ := svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT010"}, "nurse01")
	resp2, _ := svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT011"}, "nurse01")

	assert.Less(t, resp1.QueueNumber, resp2.QueueNumber)
}

// ============================================================
// CallNext Tests
// ============================================================

func TestCallNext_OK(t *testing.T) {
	svc, _ := newTestTriageService()
	// 护士签到后候诊队列（空 roomID）有一条记录
	_, _ = svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT100"}, "n1")

	entry, err := svc.CallNext(context.Background(), "")
	require.NoError(t, err)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, "calling", string(entry.Status))
}

func TestCallNext_EmptyQueue(t *testing.T) {
	svc, _ := newTestTriageService()
	_, err := svc.CallNext(context.Background(), "nonexistent-room")
	assert.True(t, bizErr.Is(err, bizErr.ErrTriageQueueEmpty))
}

// ============================================================
// Recall Tests
// ============================================================

func TestRecall_OK(t *testing.T) {
	svc, _ := newTestTriageService()
	_, _ = svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT200"}, "n1")

	// 先呼叫
	_, err := svc.CallNext(context.Background(), "")
	require.NoError(t, err)

	// 再重叫
	entry, err := svc.Recall(context.Background(), "")
	require.NoError(t, err)
	assert.Equal(t, 2, entry.CallCount)
}

// ============================================================
// GetQueueStatus Tests
// ============================================================

func TestGetQueueStatus_Empty(t *testing.T) {
	svc, _ := newTestTriageService()
	status, err := svc.GetQueueStatus(context.Background(), "room-X")
	require.NoError(t, err)
	assert.Equal(t, 0, status.WaitingCount)
	assert.Nil(t, status.CurrentCalling)
}

func TestGetQueueStatus_WithEntries(t *testing.T) {
	svc, _ := newTestTriageService()
	_, _ = svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT300"}, "n1")
	_, _ = svc.NurseCheckIn(context.Background(), NurseCheckInReq{AppointmentID: "APPT301"}, "n1")

	status, err := svc.GetQueueStatus(context.Background(), "")
	require.NoError(t, err)
	assert.Equal(t, 2, status.WaitingCount)
}

// ============================================================
// ExamExecution Tests
// ============================================================

func TestStartAndCompleteExam_OK(t *testing.T) {
	svc, execRepo := newTestTriageService()
	itemID := "ITEM001"
	// 预置 checked_in 状态的检查执行记录（实际由签到流程创建）
	seedExamExecution(execRepo, itemID, "P001", "DEV001")

	err := svc.StartExam(context.Background(), itemID, "operator1")
	require.NoError(t, err)

	err = svc.CompleteExam(context.Background(), itemID, "operator1")
	require.NoError(t, err)

	exec, err := svc.GetExamExecution(context.Background(), itemID)
	require.NoError(t, err)
	assert.Equal(t, "done", exec.Status)
	assert.NotNil(t, exec.CompletedAt)
}

func TestUndoExam_AfterStart(t *testing.T) {
	svc, execRepo := newTestTriageService()
	itemID := "ITEM002"
	seedExamExecution(execRepo, itemID, "P002", "DEV001")

	err := svc.StartExam(context.Background(), itemID, "op1")
	require.NoError(t, err)

	// 撤销开始操作（在撤销窗口内）
	err = svc.UndoExam(context.Background(), itemID, "op1", "误操作")
	require.NoError(t, err)

	exec, err := svc.GetExamExecution(context.Background(), itemID)
	require.NoError(t, err)
	// 撤销后状态退回到 checked_in
	assert.Equal(t, "checked_in", exec.Status)
}

func TestGetExamExecution_NotFound(t *testing.T) {
	svc, _ := newTestTriageService()
	_, err := svc.GetExamExecution(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestStartExam_NotFound(t *testing.T) {
	svc, _ := newTestTriageService()
	err := svc.StartExam(context.Background(), "nonexistent", "op1")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}
