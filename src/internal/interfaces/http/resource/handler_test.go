package resource

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appRes "github.com/euler/mtap/internal/application/resource"
	"github.com/euler/mtap/pkg/errors"
	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type htMockCampusRepo struct{ campuses []appRes.CampusResp }

func (r *htMockCampusRepo) List(_ context.Context) ([]appRes.CampusResp, error) {
	return r.campuses, nil
}

type htMockDeptRepo struct{ depts []appRes.DepartmentResp }

func (r *htMockDeptRepo) List(_ context.Context, campusID string) ([]appRes.DepartmentResp, error) {
	return r.depts, nil
}

type htMockDeviceRepo struct {
	devices map[string]appRes.DeviceResp
	err     error
}

func newHtMockDeviceRepo() *htMockDeviceRepo {
	return &htMockDeviceRepo{devices: make(map[string]appRes.DeviceResp)}
}
func (r *htMockDeviceRepo) Create(_ context.Context, d appRes.DeviceResp) error {
	if r.err != nil {
		return r.err
	}
	r.devices[d.ID] = d
	return nil
}
func (r *htMockDeviceRepo) Get(_ context.Context, id string) (*appRes.DeviceResp, error) {
	if r.err != nil {
		return nil, r.err
	}
	d, ok := r.devices[id]
	if !ok {
		return nil, errors.New(errors.ErrNotFound)
	}
	return &d, nil
}
func (r *htMockDeviceRepo) List(_ context.Context) ([]appRes.DeviceResp, error) {
	if r.err != nil {
		return nil, r.err
	}
	var out []appRes.DeviceResp
	for _, d := range r.devices {
		out = append(out, d)
	}
	return out, nil
}
func (r *htMockDeviceRepo) Update(_ context.Context, id string, d appRes.DeviceResp) error {
	if r.err != nil {
		return r.err
	}
	r.devices[id] = d
	return nil
}
func (r *htMockDeviceRepo) Delete(_ context.Context, id string) error {
	if r.err != nil {
		return r.err
	}
	delete(r.devices, id)
	return nil
}

type htMockExamItemRepo struct {
	items map[string]appRes.ExamItemResp
	err   error
}

func newHtMockExamItemRepo() *htMockExamItemRepo {
	return &htMockExamItemRepo{items: make(map[string]appRes.ExamItemResp)}
}
func (r *htMockExamItemRepo) Create(_ context.Context, e appRes.ExamItemResp) error {
	if r.err != nil {
		return r.err
	}
	r.items[e.ID] = e
	return nil
}
func (r *htMockExamItemRepo) Get(_ context.Context, id string) (*appRes.ExamItemResp, error) {
	if r.err != nil {
		return nil, r.err
	}
	e, ok := r.items[id]
	if !ok {
		return nil, errors.New(errors.ErrNotFound)
	}
	return &e, nil
}
func (r *htMockExamItemRepo) List(_ context.Context) ([]appRes.ExamItemResp, error) {
	if r.err != nil {
		return nil, r.err
	}
	var out []appRes.ExamItemResp
	for _, e := range r.items {
		out = append(out, e)
	}
	return out, nil
}
func (r *htMockExamItemRepo) Update(_ context.Context, id string, e appRes.ExamItemResp) error {
	if r.err != nil {
		return r.err
	}
	r.items[id] = e
	return nil
}
func (r *htMockExamItemRepo) Delete(_ context.Context, id string) error {
	if r.err != nil {
		return r.err
	}
	delete(r.items, id)
	return nil
}
func (r *htMockExamItemRepo) ListFastingIDs(_ context.Context, ids []string) ([]string, error) {
	return nil, r.err
}
func (r *htMockExamItemRepo) GetDurationMin(_ context.Context, id string) (int, error) {
	e, ok := r.items[id]
	if !ok {
		return 0, errors.New(errors.ErrNotFound)
	}
	return e.DurationMin, nil
}

type htMockAliasRepo struct {
	aliases []appRes.AliasResp
	err     error
}

func (r *htMockAliasRepo) Create(_ context.Context, a appRes.AliasResp) error {
	if r.err != nil {
		return r.err
	}
	r.aliases = append(r.aliases, a)
	return nil
}
func (r *htMockAliasRepo) List(_ context.Context, examItemID string) ([]appRes.AliasResp, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.aliases, nil
}
func (r *htMockAliasRepo) Delete(_ context.Context, aliasID string) error {
	if r.err != nil {
		return r.err
	}
	return nil
}

type htMockSlotPoolRepo struct {
	pools []appRes.SlotPoolResp
	err   error
}

func (r *htMockSlotPoolRepo) Create(_ context.Context, p appRes.SlotPoolResp) error {
	if r.err != nil {
		return r.err
	}
	r.pools = append(r.pools, p)
	return nil
}
func (r *htMockSlotPoolRepo) List(_ context.Context) ([]appRes.SlotPoolResp, error) {
	return r.pools, r.err
}

type htMockScheduleRepo struct {
	schedules []appRes.ScheduleResp
	err       error
}

func (r *htMockScheduleRepo) Create(_ context.Context, deviceID string, date time.Time, startTime, endTime string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return "sched-001", nil
}
func (r *htMockScheduleRepo) Suspend(_ context.Context, deviceID string, date time.Time, reason string) error {
	return r.err
}
func (r *htMockScheduleRepo) Substitute(_ context.Context, sourceDeviceID, targetDeviceID string, date time.Time) error {
	return r.err
}
func (r *htMockScheduleRepo) List(_ context.Context, deviceID string, startDate, endDate time.Time) ([]appRes.ScheduleResp, error) {
	return r.schedules, r.err
}

type htMockTimeSlotRepo struct {
	slots []appRes.TimeSlotResp
	err   error
}

func (r *htMockTimeSlotRepo) BulkCreate(_ context.Context, slots []appRes.TimeSlotResp) error {
	if r.err != nil {
		return r.err
	}
	r.slots = append(r.slots, slots...)
	return nil
}
func (r *htMockTimeSlotRepo) ListByDeviceAndDate(_ context.Context, deviceID string, date time.Time) ([]appRes.TimeSlotResp, error) {
	return r.slots, r.err
}
func (r *htMockTimeSlotRepo) QueryAvailable(_ context.Context, deviceID string, date time.Time, examItemID, poolType string) ([]appRes.TimeSlotResp, error) {
	return r.slots, r.err
}
func (r *htMockTimeSlotRepo) Lock(_ context.Context, slotID string, patientID string, lockUntil time.Time) error {
	return r.err
}
func (r *htMockTimeSlotRepo) Release(_ context.Context, slotID string, patientID string, allowForce bool) error {
	return r.err
}
func (r *htMockTimeSlotRepo) SuspendRange(_ context.Context, deviceID string, date time.Time, startAt, endAt time.Time, reason string) (int64, error) {
	return 0, r.err
}
func (r *htMockTimeSlotRepo) UpdateDeviceByDate(_ context.Context, sourceDeviceID, targetDeviceID string, date time.Time) (int64, error) {
	return 0, r.err
}
func (r *htMockTimeSlotRepo) HasOverlap(_ context.Context, deviceID string, date time.Time, startAt, endAt time.Time) (bool, error) {
	return false, r.err
}

type htMockDoctorRepo struct{ doctors []appRes.DoctorResp }

func (r *htMockDoctorRepo) Create(_ context.Context, d appRes.DoctorResp) error {
	r.doctors = append(r.doctors, d)
	return nil
}
func (r *htMockDoctorRepo) Get(_ context.Context, id string) (*appRes.DoctorResp, error) {
	for _, d := range r.doctors {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, nil
}
func (r *htMockDoctorRepo) List(_ context.Context, _ string) ([]appRes.DoctorResp, error) {
	return r.doctors, nil
}
func (r *htMockDoctorRepo) Update(_ context.Context, id string, d appRes.DoctorResp) error {
	return nil
}

type htMockScheduleTemplateRepo struct{ templates []appRes.ScheduleTemplateResp }

func (r *htMockScheduleTemplateRepo) Create(_ context.Context, t appRes.ScheduleTemplateResp) error {
	r.templates = append(r.templates, t)
	return nil
}
func (r *htMockScheduleTemplateRepo) Get(_ context.Context, id string) (*appRes.ScheduleTemplateResp, error) {
	return nil, nil
}
func (r *htMockScheduleTemplateRepo) List(_ context.Context) ([]appRes.ScheduleTemplateResp, error) {
	return r.templates, nil
}
func (r *htMockScheduleTemplateRepo) Delete(_ context.Context, id string) error { return nil }

// ---------------------------------------------------------------------------
// Test env helpers
// ---------------------------------------------------------------------------

type resourceTestEnv struct {
	router     *gin.Engine
	deviceRepo *htMockDeviceRepo
	examRepo   *htMockExamItemRepo
}

func newResourceTestEnv() *resourceTestEnv {
	devRepo := newHtMockDeviceRepo()
	examRepo := newHtMockExamItemRepo()

	svc := appRes.NewService(
		&htMockCampusRepo{},
		&htMockDeptRepo{},
		devRepo,
		examRepo,
		&htMockAliasRepo{},
		&htMockSlotPoolRepo{},
		&htMockScheduleRepo{},
		&htMockTimeSlotRepo{},
		&htMockDoctorRepo{},
		&htMockScheduleTemplateRepo{},
	)
	h := NewHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-admin")
		c.Set("role", "admin")
		c.Next()
	})
	v1 := r.Group("/api/v1")
	h.RegisterRoutes(v1)

	return &resourceTestEnv{router: r, deviceRepo: devRepo, examRepo: examRepo}
}

func doResourceRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseResourceBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&m); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return m
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestResourceHandler_CreateDevice_OK(t *testing.T) {
	env := newResourceTestEnv()
	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/devices", map[string]interface{}{
		"name":          "CT Scanner A",
		"campus_id":     "campus-001",
		"department_id": "dept-001",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	body := parseResourceBody(t, w)
	data := body["data"].(map[string]interface{})
	if data["name"] != "CT Scanner A" {
		t.Errorf("unexpected name: %v", data["name"])
	}
}

func TestResourceHandler_CreateDevice_BadRequest(t *testing.T) {
	env := newResourceTestEnv()
	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/devices", map[string]interface{}{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestResourceHandler_ListDevices_OK(t *testing.T) {
	env := newResourceTestEnv()
	// 先创建一个设备
	doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/devices", map[string]interface{}{
		"name": "MRI-001",
	})
	w := doResourceRequest(env.router, http.MethodGet, "/api/v1/resources/devices", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := parseResourceBody(t, w)
	data, ok := body["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Errorf("expected non-empty device list, got %v", body["data"])
	}
}

func TestResourceHandler_CreateExamItem_OK(t *testing.T) {
	env := newResourceTestEnv()
	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/exam-items", map[string]interface{}{
		"name":         "CT Scan",
		"duration_min": 30,
		"is_fasting":   false,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	body := parseResourceBody(t, w)
	data := body["data"].(map[string]interface{})
	if data["name"] != "CT Scan" {
		t.Errorf("unexpected name: %v", data["name"])
	}
}

func TestResourceHandler_CreateExamItem_BadRequest(t *testing.T) {
	env := newResourceTestEnv()
	// 缺少 name 和 duration_min
	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/exam-items", map[string]interface{}{
		"is_fasting": false,
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestResourceHandler_ListExamItems_OK(t *testing.T) {
	env := newResourceTestEnv()
	doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/exam-items", map[string]interface{}{
		"name":         "X-Ray",
		"duration_min": 15,
	})
	w := doResourceRequest(env.router, http.MethodGet, "/api/v1/resources/exam-items", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := parseResourceBody(t, w)
	data, ok := body["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Errorf("expected non-empty exam item list, got %v", body["data"])
	}
}

func TestResourceHandler_CreateAlias_OK(t *testing.T) {
	env := newResourceTestEnv()
	// 先创建检查项目
	wr := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/exam-items", map[string]interface{}{
		"name":         "CT Scan",
		"duration_min": 30,
	})
	if wr.Code != http.StatusCreated {
		t.Fatalf("create exam item failed: %s", wr.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(wr.Body).Decode(&resp)
	data := resp["data"].(map[string]interface{})
	examItemID := data["id"].(string)

	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/item-aliases", map[string]interface{}{
		"exam_item_id": examItemID,
		"alias":        "CT",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestResourceHandler_CreateSlotPool_OK(t *testing.T) {
	env := newResourceTestEnv()
	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/slot-pools", map[string]interface{}{
		"name": "公共号源池",
		"type": "public",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	body := parseResourceBody(t, w)
	data := body["data"].(map[string]interface{})
	if data["name"] != "公共号源池" {
		t.Errorf("unexpected name: %v", data["name"])
	}
}

func TestResourceHandler_ListSlotPools_OK(t *testing.T) {
	env := newResourceTestEnv()
	doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/slot-pools", map[string]interface{}{
		"name": "公共池",
		"type": "public",
	})
	w := doResourceRequest(env.router, http.MethodGet, "/api/v1/resources/slot-pools", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestResourceHandler_ListCampuses_OK(t *testing.T) {
	env := newResourceTestEnv()
	w := doResourceRequest(env.router, http.MethodGet, "/api/v1/resources/campuses", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestResourceHandler_GenerateSchedule_OK(t *testing.T) {
	env := newResourceTestEnv()
	// 先创建设备
	wr := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/devices", map[string]interface{}{
		"name": "CT-001",
	})
	if wr.Code != http.StatusCreated {
		t.Fatalf("create device: %s", wr.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(wr.Body).Decode(&resp)
	devID := resp["data"].(map[string]interface{})["id"].(string)

	w := doResourceRequest(env.router, http.MethodPost, "/api/v1/resources/schedules/generate", map[string]interface{}{
		"device_id":    devID,
		"date":         "2025-08-01",
		"start_time":   "08:00",
		"end_time":     "12:00",
		"slot_minutes": 30,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
