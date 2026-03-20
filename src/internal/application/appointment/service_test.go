package appointment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/euler/mtap/internal/domain/appointment"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// ============================================================
// Mock Repositories
// ============================================================

type mockApptRepo struct {
	items map[string]*domain.Appointment
}

func newMockApptRepo() *mockApptRepo {
	return &mockApptRepo{items: make(map[string]*domain.Appointment)}
}

func (m *mockApptRepo) Save(_ context.Context, a *domain.Appointment) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockApptRepo) Update(_ context.Context, a *domain.Appointment) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockApptRepo) FindByID(_ context.Context, id string) (*domain.Appointment, error) {
	return m.items[id], nil
}
func (m *mockApptRepo) FindByPatientID(_ context.Context, patientID string, _, _ int) ([]*domain.Appointment, int64, error) {
	var result []*domain.Appointment
	for _, a := range m.items {
		if a.PatientID == patientID {
			result = append(result, a)
		}
	}
	return result, int64(len(result)), nil
}
func (m *mockApptRepo) FindByStatus(_ context.Context, status domain.AppointmentStatus, _, _ int) ([]*domain.Appointment, int64, error) {
	var result []*domain.Appointment
	for _, a := range m.items {
		if status == "" || a.Status == status {
			result = append(result, a)
		}
	}
	return result, int64(len(result)), nil
}
func (m *mockApptRepo) FindConfirmedPendingTimeout(_ context.Context) ([]*domain.Appointment, error) {
	return nil, nil
}

type mockApptItemRepo struct{}

func (m *mockApptItemRepo) FindByAppointmentID(_ context.Context, _ string) ([]*domain.AppointmentItem, error) {
	return nil, nil
}
func (m *mockApptItemRepo) UpdateStatus(_ context.Context, _ string, _ domain.ItemStatus) error {
	return nil
}

type mockCredRepo struct {
	creds map[string]*domain.Credential
}

func newMockCredRepo() *mockCredRepo {
	return &mockCredRepo{creds: make(map[string]*domain.Credential)}
}
func (m *mockCredRepo) Save(_ context.Context, c *domain.Credential) error {
	m.creds[c.AppointmentID] = c
	return nil
}
func (m *mockCredRepo) FindByAppointmentID(_ context.Context, appointmentID string) (*domain.Credential, error) {
	return m.creds[appointmentID], nil
}

type mockBlacklistRepo struct {
	items map[string]*domain.Blacklist
}

func newMockBlacklistRepo() *mockBlacklistRepo {
	return &mockBlacklistRepo{items: make(map[string]*domain.Blacklist)}
}
func (m *mockBlacklistRepo) Save(_ context.Context, b *domain.Blacklist) error {
	m.items[b.PatientID] = b
	return nil
}
func (m *mockBlacklistRepo) Update(_ context.Context, b *domain.Blacklist) error {
	m.items[b.PatientID] = b
	return nil
}
func (m *mockBlacklistRepo) FindByPatientID(_ context.Context, patientID string) (*domain.Blacklist, error) {
	return m.items[patientID], nil
}
func (m *mockBlacklistRepo) FindAll(_ context.Context, _, _ int) ([]*domain.Blacklist, int64, error) {
	var result []*domain.Blacklist
	for _, b := range m.items {
		result = append(result, b)
	}
	return result, int64(len(result)), nil
}
func (m *mockBlacklistRepo) FindExpired(_ context.Context) ([]*domain.Blacklist, error) {
	return nil, nil
}

type mockNoShowRepo struct {
	counts map[string]int
}

func newMockNoShowRepo() *mockNoShowRepo {
	return &mockNoShowRepo{counts: make(map[string]int)}
}
func (m *mockNoShowRepo) Save(_ context.Context, r *domain.NoShowRecord) error {
	m.counts[r.PatientID]++
	return nil
}
func (m *mockNoShowRepo) CountByPatientIDInWindow(_ context.Context, patientID string, _ int) (int, error) {
	return m.counts[patientID], nil
}
func (m *mockNoShowRepo) FindByPatientID(_ context.Context, _ string) ([]*domain.NoShowRecord, error) {
	return nil, nil
}

type mockAppealRepo struct {
	items map[string]*domain.Appeal
}

func newMockAppealRepo() *mockAppealRepo {
	return &mockAppealRepo{items: make(map[string]*domain.Appeal)}
}
func (m *mockAppealRepo) Save(_ context.Context, a *domain.Appeal) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockAppealRepo) Update(_ context.Context, a *domain.Appeal) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockAppealRepo) FindByID(_ context.Context, id string) (*domain.Appeal, error) {
	return m.items[id], nil
}
func (m *mockAppealRepo) FindByBlacklistID(_ context.Context, blacklistID string) (*domain.Appeal, error) {
	for _, a := range m.items {
		if a.BlacklistID == blacklistID {
			return a, nil
		}
	}
	return nil, nil
}

// ============================================================
// Helper
// ============================================================

func newTestApptService() *AppointmentAppService {
	return NewAppointmentAppService(
		newMockApptRepo(),
		&mockApptItemRepo{},
		newMockCredRepo(),
		newMockBlacklistRepo(),
		newMockNoShowRepo(),
		newMockAppealRepo(),
	)
}

// ============================================================
// GetAppointment Tests
// ============================================================

func TestGetAppointment_NotFound(t *testing.T) {
	svc := newTestApptService()
	_, err := svc.GetAppointment(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetAppointment_OK(t *testing.T) {
	svc := newTestApptService()
	resp, err := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID:   "P001",
		ExamItemIDs: []string{"E1", "E2"},
	}, "operator1")
	require.NoError(t, err)

	appt, err := svc.GetAppointment(context.Background(), resp.AppointmentID)
	require.NoError(t, err)
	assert.Equal(t, "P001", appt.PatientID)
}

// ============================================================
// ListAppointments Tests
// ============================================================

func TestListAppointments_All(t *testing.T) {
	svc := newTestApptService()
	_, _ = svc.AutoAppointment(context.Background(), AutoAppointmentReq{PatientID: "P001", ExamItemIDs: []string{"E1"}}, "op")
	_, _ = svc.AutoAppointment(context.Background(), AutoAppointmentReq{PatientID: "P002", ExamItemIDs: []string{"E2"}}, "op")

	list, total, err := svc.ListAppointments(context.Background(), ListAppointmentReq{Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

func TestListAppointments_ByPatient(t *testing.T) {
	svc := newTestApptService()
	_, _ = svc.AutoAppointment(context.Background(), AutoAppointmentReq{PatientID: "P001", ExamItemIDs: []string{"E1"}}, "op")
	_, _ = svc.AutoAppointment(context.Background(), AutoAppointmentReq{PatientID: "P002", ExamItemIDs: []string{"E2"}}, "op")

	list, total, err := svc.ListAppointments(context.Background(), ListAppointmentReq{PatientID: "P001", Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "P001", list[0].PatientID)
}

// ============================================================
// AutoAppointment Tests
// ============================================================

func TestAutoAppointment_OK(t *testing.T) {
	svc := newTestApptService()
	resp, err := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID:   "P001",
		ExamItemIDs: []string{"E1"},
	}, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AppointmentID)
}

func TestAutoAppointment_Blacklisted(t *testing.T) {
	svc := newTestApptService()
	// 连续爽约 3 次触发黑名单（状态流转: pending→confirmed→paid→no_show）
	for i := 0; i < 3; i++ {
		resp, _ := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
			PatientID: "P_BL", ExamItemIDs: []string{"E1"},
		}, "op")
		_, _ = svc.ConfirmAppointment(context.Background(), resp.AppointmentID)
		_, _ = svc.MarkPaid(context.Background(), resp.AppointmentID)
		_ = svc.RecordNoShow(context.Background(), resp.AppointmentID)
	}
	// 第 4 次应被拦截
	_, err := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P_BL", ExamItemIDs: []string{"E1"},
	}, "op")
	assert.True(t, bizErr.Is(err, bizErr.ErrApptBlacklisted))
}

// ============================================================
// ManualAppointment Tests
// ============================================================

func TestManualAppointment_OK(t *testing.T) {
	svc := newTestApptService()
	resp, err := svc.ManualAppointment(context.Background(), ManualAppointmentReq{
		PatientID:  "P001",
		ExamItemID: "E1",
		SlotID:     "SLOT001",
		Reason:     "紧急患者人工干预",
	}, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "manual", string(resp.Mode))
}

// ============================================================
// ConfirmAppointment Tests
// ============================================================

func TestConfirmAppointment_OK(t *testing.T) {
	svc := newTestApptService()
	created, err := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P001", ExamItemIDs: []string{"E1"},
	}, "op")
	require.NoError(t, err)

	confirmed, err := svc.ConfirmAppointment(context.Background(), created.AppointmentID)
	require.NoError(t, err)
	assert.Equal(t, "confirmed", string(confirmed.Status))
}

func TestConfirmAppointment_NotFound(t *testing.T) {
	svc := newTestApptService()
	_, err := svc.ConfirmAppointment(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// CancelAppointment Tests
// ============================================================

func TestCancelAppointment_OK(t *testing.T) {
	svc := newTestApptService()
	created, _ := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P001", ExamItemIDs: []string{"E1"},
	}, "op")
	_, _ = svc.ConfirmAppointment(context.Background(), created.AppointmentID)

	err := svc.CancelAppointment(context.Background(), created.AppointmentID, CancelReq{Reason: "患者主动取消"}, "op")
	assert.NoError(t, err)

	appt, _ := svc.GetAppointment(context.Background(), created.AppointmentID)
	assert.Equal(t, "cancelled", string(appt.Status))
}

func TestCancelAppointment_NotFound(t *testing.T) {
	svc := newTestApptService()
	err := svc.CancelAppointment(context.Background(), "nonexistent", CancelReq{Reason: "test"}, "op")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// Credential Tests
// ============================================================

func TestGetCredential_GeneratesNew(t *testing.T) {
	svc := newTestApptService()
	created, _ := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P001", ExamItemIDs: []string{"E1"},
	}, "op")

	cred, err := svc.GetCredential(context.Background(), created.AppointmentID)
	require.NoError(t, err)
	assert.NotEmpty(t, cred.QRCodeURL)
}

func TestGetCredential_ReturnsCached(t *testing.T) {
	svc := newTestApptService()
	created, _ := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P001", ExamItemIDs: []string{"E1"},
	}, "op")

	cred1, _ := svc.GetCredential(context.Background(), created.AppointmentID)
	cred2, _ := svc.GetCredential(context.Background(), created.AppointmentID)
	assert.Equal(t, cred1.QRCodeURL, cred2.QRCodeURL)
}

// ============================================================
// NoShow / Blacklist Tests
// ============================================================

func TestRecordNoShow_OK(t *testing.T) {
	svc := newTestApptService()
	created, _ := svc.AutoAppointment(context.Background(), AutoAppointmentReq{
		PatientID: "P001", ExamItemIDs: []string{"E1"},
	}, "op")
	_, _ = svc.ConfirmAppointment(context.Background(), created.AppointmentID)
	// 缴费后才可标记爽约（状态流转: pending→confirmed→paid→no_show）
	_, _ = svc.MarkPaid(context.Background(), created.AppointmentID)

	err := svc.RecordNoShow(context.Background(), created.AppointmentID)
	assert.NoError(t, err)

	appt, _ := svc.GetAppointment(context.Background(), created.AppointmentID)
	assert.Equal(t, "no_show", string(appt.Status))
}

func TestRecordNoShow_NotFound(t *testing.T) {
	svc := newTestApptService()
	err := svc.RecordNoShow(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}
