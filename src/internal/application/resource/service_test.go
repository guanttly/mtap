package resource

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// ============================================================
// Mock Repositories
// ============================================================

type mockCampusRepo struct {
	items []CampusResp
}

func (m *mockCampusRepo) List(_ context.Context) ([]CampusResp, error) { return m.items, nil }

type mockDeptRepo struct {
	items []DepartmentResp
}

func (m *mockDeptRepo) List(_ context.Context, campusID string) ([]DepartmentResp, error) {
	if campusID == "" {
		return m.items, nil
	}
	var result []DepartmentResp
	for _, d := range m.items {
		if d.CampusID == campusID {
			result = append(result, d)
		}
	}
	return result, nil
}

type mockDeviceRepo struct {
	items map[string]DeviceResp
}

func newMockDeviceRepo() *mockDeviceRepo { return &mockDeviceRepo{items: make(map[string]DeviceResp)} }

func (m *mockDeviceRepo) Create(_ context.Context, d DeviceResp) error {
	m.items[d.ID] = d
	return nil
}
func (m *mockDeviceRepo) Get(_ context.Context, id string) (*DeviceResp, error) {
	if d, ok := m.items[id]; ok {
		return &d, nil
	}
	return nil, nil
}
func (m *mockDeviceRepo) List(_ context.Context) ([]DeviceResp, error) {
	var result []DeviceResp
	for _, d := range m.items {
		result = append(result, d)
	}
	return result, nil
}
func (m *mockDeviceRepo) Update(_ context.Context, id string, d DeviceResp) error {
	m.items[id] = d
	return nil
}
func (m *mockDeviceRepo) Delete(_ context.Context, id string) error {
	delete(m.items, id)
	return nil
}

type mockExamRepo struct {
	items map[string]ExamItemResp
}

func newMockExamRepo() *mockExamRepo { return &mockExamRepo{items: make(map[string]ExamItemResp)} }

func (m *mockExamRepo) Create(_ context.Context, e ExamItemResp) error {
	m.items[e.ID] = e
	return nil
}
func (m *mockExamRepo) Get(_ context.Context, id string) (*ExamItemResp, error) {
	if e, ok := m.items[id]; ok {
		return &e, nil
	}
	return nil, nil
}
func (m *mockExamRepo) List(_ context.Context) ([]ExamItemResp, error) {
	var result []ExamItemResp
	for _, e := range m.items {
		result = append(result, e)
	}
	return result, nil
}
func (m *mockExamRepo) Update(_ context.Context, id string, e ExamItemResp) error {
	m.items[id] = e
	return nil
}
func (m *mockExamRepo) Delete(_ context.Context, id string) error {
	delete(m.items, id)
	return nil
}
func (m *mockExamRepo) ListFastingIDs(_ context.Context, ids []string) ([]string, error) {
	var result []string
	for _, id := range ids {
		if e, ok := m.items[id]; ok && e.IsFasting {
			result = append(result, id)
		}
	}
	return result, nil
}
func (m *mockExamRepo) GetDurationMin(_ context.Context, id string) (int, error) {
	if e, ok := m.items[id]; ok {
		return e.DurationMin, nil
	}
	return 0, nil
}

type mockAliasRepo struct {
	items map[string]AliasResp
}

func newMockAliasRepo() *mockAliasRepo { return &mockAliasRepo{items: make(map[string]AliasResp)} }

func (m *mockAliasRepo) Create(_ context.Context, a AliasResp) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockAliasRepo) List(_ context.Context, examItemID string) ([]AliasResp, error) {
	var result []AliasResp
	for _, a := range m.items {
		if a.ExamItemID == examItemID {
			result = append(result, a)
		}
	}
	return result, nil
}
func (m *mockAliasRepo) Delete(_ context.Context, id string) error {
	delete(m.items, id)
	return nil
}

type mockSlotPoolRepo struct {
	items []SlotPoolResp
}

func (m *mockSlotPoolRepo) Create(_ context.Context, p SlotPoolResp) error {
	m.items = append(m.items, p)
	return nil
}
func (m *mockSlotPoolRepo) List(_ context.Context) ([]SlotPoolResp, error) { return m.items, nil }

type mockDoctorRepo struct {
	items map[string]DoctorResp
}

func newMockDoctorRepo() *mockDoctorRepo { return &mockDoctorRepo{items: map[string]DoctorResp{}} }
func (m *mockDoctorRepo) Create(_ context.Context, d DoctorResp) error {
	m.items[d.ID] = d
	return nil
}
func (m *mockDoctorRepo) Get(_ context.Context, id string) (*DoctorResp, error) {
	v, ok := m.items[id]
	if !ok {
		return nil, nil
	}
	return &v, nil
}
func (m *mockDoctorRepo) List(_ context.Context, _ string) ([]DoctorResp, error) {
	out := make([]DoctorResp, 0, len(m.items))
	for _, v := range m.items {
		out = append(out, v)
	}
	return out, nil
}
func (m *mockDoctorRepo) Update(_ context.Context, id string, d DoctorResp) error {
	m.items[id] = d
	return nil
}

type mockScheduleTemplateRepo struct {
	items map[string]ScheduleTemplateResp
}

func newMockScheduleTemplateRepo() *mockScheduleTemplateRepo {
	return &mockScheduleTemplateRepo{items: map[string]ScheduleTemplateResp{}}
}
func (m *mockScheduleTemplateRepo) Create(_ context.Context, t ScheduleTemplateResp) error {
	m.items[t.ID] = t
	return nil
}
func (m *mockScheduleTemplateRepo) Get(_ context.Context, id string) (*ScheduleTemplateResp, error) {
	v, ok := m.items[id]
	if !ok {
		return nil, nil
	}
	return &v, nil
}
func (m *mockScheduleTemplateRepo) List(_ context.Context) ([]ScheduleTemplateResp, error) {
	out := make([]ScheduleTemplateResp, 0, len(m.items))
	for _, v := range m.items {
		out = append(out, v)
	}
	return out, nil
}
func (m *mockScheduleTemplateRepo) Delete(_ context.Context, id string) error {
	delete(m.items, id)
	return nil
}

type mockScheduleRepo struct {
	schedules []string // IDs
}

func (m *mockScheduleRepo) Create(_ context.Context, _ string, _ time.Time, _, _ string) (string, error) {
	id := "sched-001"
	m.schedules = append(m.schedules, id)
	return id, nil
}
func (m *mockScheduleRepo) Suspend(_ context.Context, _ string, _ time.Time, _ string) error {
	return nil
}
func (m *mockScheduleRepo) Substitute(_ context.Context, _, _ string, _ time.Time) error { return nil }
func (m *mockScheduleRepo) List(_ context.Context, _ string, _, _ time.Time) ([]ScheduleResp, error) {
	return nil, nil
}

type mockTimeSlotRepo struct {
	slots map[string]TimeSlotResp
}

func newMockTimeSlotRepo() *mockTimeSlotRepo {
	return &mockTimeSlotRepo{slots: make(map[string]TimeSlotResp)}
}
func (m *mockTimeSlotRepo) BulkCreate(_ context.Context, slots []TimeSlotResp) error {
	for _, s := range slots {
		m.slots[s.ID] = s
	}
	return nil
}
func (m *mockTimeSlotRepo) ListByDeviceAndDate(_ context.Context, deviceID string, date time.Time) ([]TimeSlotResp, error) {
	var result []TimeSlotResp
	for _, s := range m.slots {
		if s.DeviceID == deviceID && s.StartAt.Format("2006-01-02") == date.Format("2006-01-02") {
			result = append(result, s)
		}
	}
	return result, nil
}
func (m *mockTimeSlotRepo) QueryAvailable(_ context.Context, deviceID string, date time.Time, examItemID, poolType string) ([]TimeSlotResp, error) {
	var result []TimeSlotResp
	for _, s := range m.slots {
		if s.DeviceID == deviceID && s.StartAt.Format("2006-01-02") == date.Format("2006-01-02") && s.Status == "available" {
			result = append(result, s)
		}
	}
	return result, nil
}
func (m *mockTimeSlotRepo) Lock(_ context.Context, slotID, _ string, _ time.Time) error {
	if s, ok := m.slots[slotID]; ok {
		s.Status = "locked"
		m.slots[slotID] = s
	}
	return nil
}
func (m *mockTimeSlotRepo) Release(_ context.Context, slotID, _ string, _ bool) error {
	if s, ok := m.slots[slotID]; ok {
		s.Status = "available"
		m.slots[slotID] = s
	}
	return nil
}
func (m *mockTimeSlotRepo) SuspendRange(_ context.Context, _ string, _ time.Time, _, _ time.Time, _ string) (int64, error) {
	return 0, nil
}
func (m *mockTimeSlotRepo) UpdateDeviceByDate(_ context.Context, _, _ string, _ time.Time) (int64, error) {
	return 0, nil
}
func (m *mockTimeSlotRepo) HasOverlap(_ context.Context, _ string, _ time.Time, _, _ time.Time) (bool, error) {
	return false, nil
}

// ============================================================
// Helper
// ============================================================

func newTestResourceService() *Service {
	return NewService(
		&mockCampusRepo{items: []CampusResp{{ID: "C1", Name: "主院区", Code: "MAIN", Status: "active"}}},
		&mockDeptRepo{items: []DepartmentResp{{ID: "D1", CampusID: "C1", Name: "超声科", Code: "US", Status: "active"}}},
		newMockDeviceRepo(),
		newMockExamRepo(),
		newMockAliasRepo(),
		&mockSlotPoolRepo{},
		&mockScheduleRepo{},
		newMockTimeSlotRepo(),
		newMockDoctorRepo(),
		newMockScheduleTemplateRepo(),
	)
}

// ============================================================
// Device Tests
// ============================================================

func TestCreateDevice_OK(t *testing.T) {
	svc := newTestResourceService()
	resp, err := svc.CreateDevice(context.Background(), CreateDeviceReq{
		Name: "超声仪 A", CampusID: "C1", DepartmentID: "D1",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "超声仪 A", resp.Name)
	assert.Equal(t, "active", resp.Status)
}

func TestListDevices_Empty(t *testing.T) {
	svc := newTestResourceService()
	list, err := svc.ListDevices(context.Background())
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestListDevices_AfterCreate(t *testing.T) {
	svc := newTestResourceService()
	_, _ = svc.CreateDevice(context.Background(), CreateDeviceReq{Name: "设备1", CampusID: "C1"})
	_, _ = svc.CreateDevice(context.Background(), CreateDeviceReq{Name: "设备2", CampusID: "C1"})
	list, err := svc.ListDevices(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestUpdateDevice_OK(t *testing.T) {
	svc := newTestResourceService()
	d, _ := svc.CreateDevice(context.Background(), CreateDeviceReq{Name: "旧名称", CampusID: "C1"})
	updated, err := svc.UpdateDevice(context.Background(), d.ID, UpdateDeviceReq{Name: "新名称"})
	require.NoError(t, err)
	assert.Equal(t, "新名称", updated.Name)
}

func TestUpdateDevice_NotFound(t *testing.T) {
	svc := newTestResourceService()
	_, err := svc.UpdateDevice(context.Background(), "nonexistent", UpdateDeviceReq{Name: "X"})
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestDeleteDevice_OK(t *testing.T) {
	svc := newTestResourceService()
	d, _ := svc.CreateDevice(context.Background(), CreateDeviceReq{Name: "设备", CampusID: "C1"})
	err := svc.DeleteDevice(context.Background(), d.ID)
	assert.NoError(t, err)

	list, _ := svc.ListDevices(context.Background())
	assert.Empty(t, list)
}

func TestDeleteDevice_NotFound(t *testing.T) {
	svc := newTestResourceService()
	err := svc.DeleteDevice(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// ExamItem Tests
// ============================================================

func TestCreateExamItem_OK(t *testing.T) {
	svc := newTestResourceService()
	resp, err := svc.CreateExamItem(context.Background(), CreateExamItemReq{
		Name: "腹部超声", DurationMin: 20, IsFasting: true, FastingDesc: "检查前8小时禁食",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "腹部超声", resp.Name)
	assert.True(t, resp.IsFasting)
}

func TestListExamItems_AfterCreate(t *testing.T) {
	svc := newTestResourceService()
	_, _ = svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "CT平扫", DurationMin: 15})
	_, _ = svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "MRI", DurationMin: 45})
	list, err := svc.ListExamItems(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestUpdateExamItem_OK(t *testing.T) {
	svc := newTestResourceService()
	item, _ := svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "旧项目", DurationMin: 30})
	isFasting := true
	updated, err := svc.UpdateExamItem(context.Background(), item.ID, UpdateExamItemReq{
		Name: "新项目", DurationMin: 45, IsFasting: &isFasting,
	})
	require.NoError(t, err)
	assert.Equal(t, "新项目", updated.Name)
	assert.Equal(t, 45, updated.DurationMin)
	assert.True(t, updated.IsFasting)
}

func TestUpdateExamItem_NotFound(t *testing.T) {
	svc := newTestResourceService()
	_, err := svc.UpdateExamItem(context.Background(), "nonexistent", UpdateExamItemReq{Name: "X"})
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestDeleteExamItem_OK(t *testing.T) {
	svc := newTestResourceService()
	item, _ := svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "项目", DurationMin: 20})
	err := svc.DeleteExamItem(context.Background(), item.ID)
	assert.NoError(t, err)
}

// ============================================================
// Alias Tests
// ============================================================

func TestCreateAlias_OK(t *testing.T) {
	svc := newTestResourceService()
	item, _ := svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "腹部超声", DurationMin: 20})
	alias, err := svc.CreateAlias(context.Background(), CreateAliasReq{
		ExamItemID: item.ID, Alias: "B超",
	})
	require.NoError(t, err)
	assert.Equal(t, "B超", alias.Alias)
	assert.Equal(t, item.ID, alias.ExamItemID)
}

func TestListAliases_AfterCreate(t *testing.T) {
	svc := newTestResourceService()
	item, _ := svc.CreateExamItem(context.Background(), CreateExamItemReq{Name: "CT", DurationMin: 15})
	_, _ = svc.CreateAlias(context.Background(), CreateAliasReq{ExamItemID: item.ID, Alias: "CT扫描"})
	_, _ = svc.CreateAlias(context.Background(), CreateAliasReq{ExamItemID: item.ID, Alias: "断层扫描"})

	aliases, err := svc.ListAliases(context.Background(), item.ID)
	require.NoError(t, err)
	assert.Len(t, aliases, 2)
}

// ============================================================
// SlotPool Tests
// ============================================================

func TestCreateSlotPool_OK(t *testing.T) {
	svc := newTestResourceService()
	resp, err := svc.CreateSlotPool(context.Background(), CreateSlotPoolReq{
		Name: "公共号源池", Type: "public",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "public", resp.Type)
	assert.Equal(t, "active", resp.Status)
}

func TestListSlotPools_OK(t *testing.T) {
	svc := newTestResourceService()
	_, _ = svc.CreateSlotPool(context.Background(), CreateSlotPoolReq{Name: "公共池", Type: "public"})
	_, _ = svc.CreateSlotPool(context.Background(), CreateSlotPoolReq{Name: "科室池", Type: "department"})
	list, err := svc.ListSlotPools(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

// ============================================================
// Schedule & Slot Tests
// ============================================================

func TestGenerateSchedule_SingleDay(t *testing.T) {
	svc := newTestResourceService()
	slots, err := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID:    "DEV001",
		Date:        "2026-04-01",
		StartTime:   "08:00",
		EndTime:     "12:00",
		SlotMinutes: 30,
		PoolType:    "public",
	})
	require.NoError(t, err)
	// 08:00~12:00，每30分钟一个号，共8个
	assert.Len(t, slots, 8)
	assert.Equal(t, "available", slots[0].Status)
	assert.Equal(t, 30, slots[0].StandardDuration)
}

func TestGenerateSchedule_BatchDays(t *testing.T) {
	svc := newTestResourceService()
	slots, err := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID:     "DEV001",
		StartDate:    "2026-04-01",
		EndDate:      "2026-04-03",
		StartTime:    "08:00",
		EndTime:      "09:00",
		SlotMinutes:  30,
		SkipWeekends: false,
	})
	require.NoError(t, err)
	// 3天 × 2个号/天 = 6个
	assert.Len(t, slots, 6)
}

func TestGenerateSchedule_InvalidDateFormat(t *testing.T) {
	svc := newTestResourceService()
	_, err := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID: "DEV001", Date: "invalid", StartTime: "08:00", EndTime: "12:00", SlotMinutes: 30,
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrInvalidParam))
}

func TestGenerateSchedule_EndBeforeStart(t *testing.T) {
	svc := newTestResourceService()
	_, err := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID: "DEV001", Date: "2026-04-01", StartTime: "12:00", EndTime: "08:00", SlotMinutes: 30,
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrInvalidParam))
}

func TestQueryAvailableSlots_AgeAdjustment(t *testing.T) {
	svc := newTestResourceService()
	// 先生成号源
	_, _ = svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID: "DEV001", Date: "2026-04-01", StartTime: "08:00", EndTime: "09:00", SlotMinutes: 30,
	})

	// 普通成人
	adult, err := svc.QueryAvailableSlots(context.Background(), "DEV001", "2026-04-01", "", "public", 30)
	require.NoError(t, err)
	if len(adult) > 0 {
		assert.Equal(t, 30, adult[0].AdjustedDuration)
	}

	// 儿童（<14）：+10%
	child, err := svc.QueryAvailableSlots(context.Background(), "DEV001", "2026-04-01", "", "public", 10)
	require.NoError(t, err)
	if len(child) > 0 {
		assert.Equal(t, 33, child[0].AdjustedDuration) // 30 * 1.10 = 33
	}

	// 老年（>70）：+15%
	elder, err := svc.QueryAvailableSlots(context.Background(), "DEV001", "2026-04-01", "", "public", 75)
	require.NoError(t, err)
	if len(elder) > 0 {
		assert.Equal(t, 35, elder[0].AdjustedDuration) // 30 * 1.15 ≈ 35
	}
}

func TestLockSlot_OK(t *testing.T) {
	svc := newTestResourceService()
	slots, _ := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID: "DEV001", Date: "2026-04-01", StartTime: "08:00", EndTime: "08:30", SlotMinutes: 30,
	})
	require.Len(t, slots, 1)

	err := svc.LockSlot(context.Background(), slots[0].ID, LockSlotReq{PatientID: "P001"}, false)
	assert.NoError(t, err)
}

func TestReleaseSlot_OK(t *testing.T) {
	svc := newTestResourceService()
	slots, _ := svc.GenerateSchedule(context.Background(), GenerateScheduleReq{
		DeviceID: "DEV001", Date: "2026-04-01", StartTime: "08:00", EndTime: "08:30", SlotMinutes: 30,
	})
	require.Len(t, slots, 1)

	_ = svc.LockSlot(context.Background(), slots[0].ID, LockSlotReq{PatientID: "P001"}, false)
	err := svc.ReleaseSlot(context.Background(), slots[0].ID, "P001", false)
	assert.NoError(t, err)
}

func TestListCampuses_OK(t *testing.T) {
	svc := newTestResourceService()
	list, err := svc.ListCampuses(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "主院区", list[0].Name)
}

func TestListDepartments_OK(t *testing.T) {
	svc := newTestResourceService()
	list, err := svc.ListDepartments(context.Background(), "C1")
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "超声科", list[0].Name)
}
