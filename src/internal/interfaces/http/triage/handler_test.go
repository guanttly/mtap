// Package triage HTTP handler 集成测试
package triage

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appTriage "github.com/euler/mtap/internal/application/triage"
	domain "github.com/euler/mtap/internal/domain/triage"
)

// ─── Mock 仓库 ────────────────────────────────────────────────────────────────

type htMockCheckInRepo struct{ records map[string]*domain.CheckIn }

func newHtCheckInRepo() *htMockCheckInRepo {
	return &htMockCheckInRepo{records: make(map[string]*domain.CheckIn)}
}
func (m *htMockCheckInRepo) Save(_ context.Context, c *domain.CheckIn) error {
	m.records[c.ID] = c
	return nil
}
func (m *htMockCheckInRepo) FindByAppointmentID(_ context.Context, appointmentID string) (*domain.CheckIn, error) {
	for _, c := range m.records {
		if c.AppointmentID == appointmentID {
			return c, nil
		}
	}
	return nil, nil
}
func (m *htMockCheckInRepo) FindByPatientAndDate(_ context.Context, _, _ string) ([]*domain.CheckIn, error) {
	return nil, nil
}

type htMockWaitingQueueRepo struct {
	queues map[string]*domain.WaitingQueue
}

func newHtWaitingQueueRepo() *htMockWaitingQueueRepo {
	return &htMockWaitingQueueRepo{queues: make(map[string]*domain.WaitingQueue)}
}
func (m *htMockWaitingQueueRepo) Save(_ context.Context, q *domain.WaitingQueue) error {
	m.queues[q.ID] = q
	return nil
}
func (m *htMockWaitingQueueRepo) Update(_ context.Context, q *domain.WaitingQueue) error {
	m.queues[q.ID] = q
	return nil
}
func (m *htMockWaitingQueueRepo) FindByRoomID(_ context.Context, roomID string) (*domain.WaitingQueue, error) {
	for _, q := range m.queues {
		if q.RoomID == roomID {
			return q, nil
		}
	}
	return nil, nil
}
func (m *htMockWaitingQueueRepo) FindOrCreateByRoom(_ context.Context, roomID, deviceID, departmentID string) (*domain.WaitingQueue, error) {
	for _, q := range m.queues {
		if q.RoomID == roomID {
			return q, nil
		}
	}
	q := domain.NewWaitingQueue(roomID, deviceID, departmentID)
	m.queues[q.ID] = q
	return q, nil
}

type htMockQueueEntryRepo struct{ entries map[string]*domain.QueueEntry }

func newHtQueueEntryRepo() *htMockQueueEntryRepo {
	return &htMockQueueEntryRepo{entries: make(map[string]*domain.QueueEntry)}
}
func (m *htMockQueueEntryRepo) Save(_ context.Context, e *domain.QueueEntry) error {
	m.entries[e.ID] = e
	return nil
}
func (m *htMockQueueEntryRepo) Update(_ context.Context, e *domain.QueueEntry) error {
	m.entries[e.ID] = e
	return nil
}
func (m *htMockQueueEntryRepo) FindByQueueID(_ context.Context, queueID string, status domain.EntryStatus) ([]*domain.QueueEntry, error) {
	var out []*domain.QueueEntry
	for _, e := range m.entries {
		if e.QueueID == queueID && e.Status == status {
			out = append(out, e)
		}
	}
	return out, nil
}

type htMockExecRepo struct {
	execs map[string]*domain.ExamExecution
}

func newHtExecRepo() *htMockExecRepo {
	return &htMockExecRepo{execs: make(map[string]*domain.ExamExecution)}
}
func (m *htMockExecRepo) Save(_ context.Context, e *domain.ExamExecution) error {
	m.execs[e.ID] = e
	return nil
}
func (m *htMockExecRepo) Update(_ context.Context, e *domain.ExamExecution) error {
	m.execs[e.ID] = e
	return nil
}
func (m *htMockExecRepo) FindByID(_ context.Context, id string) (*domain.ExamExecution, error) {
	return m.execs[id], nil
}
func (m *htMockExecRepo) FindByAppointmentItemID(_ context.Context, itemID string) (*domain.ExamExecution, error) {
	for _, e := range m.execs {
		if e.AppointmentItemID == itemID {
			return e, nil
		}
	}
	return nil, nil
}
func (m *htMockExecRepo) FindByDevice(_ context.Context, deviceID string) ([]*domain.ExamExecution, error) {
	var out []*domain.ExamExecution
	for _, e := range m.execs {
		if e.DeviceID == deviceID {
			out = append(out, e)
		}
	}
	return out, nil
}

// ─── 辅助函数 ─────────────────────────────────────────────────────────────────

type triageTestEnv struct {
	h        *Handler
	r        *gin.Engine
	execRepo *htMockExecRepo
}

func newTestTriageEnv() *triageTestEnv {
	checkInRepo := newHtCheckInRepo()
	queueRepo := newHtWaitingQueueRepo()
	entryRepo := newHtQueueEntryRepo()
	execRepo := newHtExecRepo()
	svc := appTriage.NewTriageAppService(checkInRepo, queueRepo, entryRepo, execRepo)
	h := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-nurse")
		c.Set("role", "nurse")
		c.Next()
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return &triageTestEnv{h: h, r: r, execRepo: execRepo}
}

func doTriageRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseTriageBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&m))
	return m
}

// ─── 测试用例 ─────────────────────────────────────────────────────────────────

func TestTriageHandler_NurseCheckIn_OK(t *testing.T) {
	env := newTestTriageEnv()

	w := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/checkin/nurse", map[string]interface{}{
		"appointment_id": "APPT-001",
	})
	assert.Equal(t, http.StatusOK, w.Code)
	body := parseTriageBody(t, w)
	assert.Equal(t, float64(0), body["code"])
	data := body["data"].(map[string]interface{})
	assert.NotEmpty(t, data["check_in_id"])
}

func TestTriageHandler_NurseCheckIn_BadRequest(t *testing.T) {
	env := newTestTriageEnv()

	// appointment_id 缺失
	w := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/checkin/nurse", map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTriageHandler_NurseCheckIn_Duplicate(t *testing.T) {
	env := newTestTriageEnv()

	// 第一次签到
	doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/checkin/nurse", map[string]interface{}{
		"appointment_id": "APPT-DUP",
	})
	// 第二次重复签到 → 应返回错误（冲突）
	w := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/checkin/nurse", map[string]interface{}{
		"appointment_id": "APPT-DUP",
	})
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestTriageHandler_GetQueueStatus_Empty(t *testing.T) {
	env := newTestTriageEnv()

	w := doTriageRequest(env.r, http.MethodGet, "/api/v1/triage/queue/ROOM-001", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	body := parseTriageBody(t, w)
	assert.Equal(t, float64(0), body["code"])
}

func TestTriageHandler_GetQueueStatus_WithEntries(t *testing.T) {
	env := newTestTriageEnv()

	// 先签到，创建队列
	doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/checkin/nurse", map[string]interface{}{
		"appointment_id": "APPT-Q01",
		"room_id":        "ROOM-002",
	})

	w := doTriageRequest(env.r, http.MethodGet, "/api/v1/triage/queue/ROOM-002", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTriageHandler_CallNext_EmptyQueue(t *testing.T) {
	env := newTestTriageEnv()

	// 队列为空时呼叫 → 返回业务错误（not found）
	w := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/call/ROOM-999/next", nil)
	// 空队列返回 404 或业务错误
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestTriageHandler_GetExamExecution_NotFound(t *testing.T) {
	env := newTestTriageEnv()

	w := doTriageRequest(env.r, http.MethodGet, "/api/v1/triage/exam/non-existent", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTriageHandler_StartExam_NotFound(t *testing.T) {
	env := newTestTriageEnv()

	w := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/exam/non-existent/start", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTriageHandler_StartAndCompleteExam_OK(t *testing.T) {
	env := newTestTriageEnv()

	// 预插入一条 checked_in 状态的 ExamExecution
	exec := &domain.ExamExecution{
		ID:                "EXEC-001",
		AppointmentItemID: "ITEM-001",
		PatientID:         "P001",
		DeviceID:          "D001",
		Status:            domain.ExamCheckedIn,
	}
	_ = env.execRepo.Save(context.Background(), exec)

	// 开始检查
	w1 := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/exam/ITEM-001/start", nil)
	assert.Equal(t, http.StatusOK, w1.Code)
	body1 := parseTriageBody(t, w1)
	assert.Equal(t, float64(0), body1["code"])

	// 完成检查
	w2 := doTriageRequest(env.r, http.MethodPost, "/api/v1/triage/exam/ITEM-001/complete", nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	body2 := parseTriageBody(t, w2)
	assert.Equal(t, float64(0), body2["code"])
}
