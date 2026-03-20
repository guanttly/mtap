package appointment

import (
"bytes"
"context"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/gin-gonic/gin"
appAppt "github.com/euler/mtap/internal/application/appointment"
domain "github.com/euler/mtap/internal/domain/appointment"
"github.com/euler/mtap/pkg/errors"
)

func init() { gin.SetMode(gin.TestMode) }

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type htMockApptRepo struct {
store map[string]*domain.Appointment
err   error
}

func newHtMockApptRepo() *htMockApptRepo {
return &htMockApptRepo{store: make(map[string]*domain.Appointment)}
}
func (r *htMockApptRepo) Save(_ context.Context, a *domain.Appointment) error {
if r.err != nil {
return r.err
}
r.store[a.ID] = a
return nil
}
func (r *htMockApptRepo) Update(_ context.Context, a *domain.Appointment) error {
if r.err != nil {
return r.err
}
r.store[a.ID] = a
return nil
}
func (r *htMockApptRepo) FindByID(_ context.Context, id string) (*domain.Appointment, error) {
if r.err != nil {
return nil, r.err
}
a, ok := r.store[id]
if !ok {
return nil, nil
}
return a, nil
}
func (r *htMockApptRepo) FindByPatientID(_ context.Context, patientID string, page, size int) ([]*domain.Appointment, int64, error) {
if r.err != nil {
return nil, 0, r.err
}
var out []*domain.Appointment
for _, a := range r.store {
if a.PatientID == patientID {
out = append(out, a)
}
}
return out, int64(len(out)), nil
}
func (r *htMockApptRepo) FindByStatus(_ context.Context, status domain.AppointmentStatus, page, size int) ([]*domain.Appointment, int64, error) {
if r.err != nil {
return nil, 0, r.err
}
var out []*domain.Appointment
for _, a := range r.store {
if status == "" || a.Status == status {
out = append(out, a)
}
}
return out, int64(len(out)), nil
}
func (r *htMockApptRepo) FindConfirmedPendingTimeout(_ context.Context) ([]*domain.Appointment, error) {
return nil, r.err
}

type htMockApptItemRepo struct{}

func (r *htMockApptItemRepo) FindByAppointmentID(_ context.Context, appointmentID string) ([]*domain.AppointmentItem, error) {
return nil, nil
}
func (r *htMockApptItemRepo) UpdateStatus(_ context.Context, itemID string, status domain.ItemStatus) error {
return nil
}

type htMockCredRepo struct {
store map[string]*domain.Credential
}

func newHtMockCredRepo() *htMockCredRepo {
return &htMockCredRepo{store: make(map[string]*domain.Credential)}
}
func (r *htMockCredRepo) Save(_ context.Context, c *domain.Credential) error {
r.store[c.AppointmentID] = c
return nil
}
func (r *htMockCredRepo) FindByAppointmentID(_ context.Context, appointmentID string) (*domain.Credential, error) {
c, ok := r.store[appointmentID]
if !ok {
return nil, nil
}
return c, nil
}

type htMockBlacklistRepo struct{}

func (r *htMockBlacklistRepo) Save(_ context.Context, b *domain.Blacklist) error { return nil }
func (r *htMockBlacklistRepo) Update(_ context.Context, b *domain.Blacklist) error {
return nil
}
func (r *htMockBlacklistRepo) FindByPatientID(_ context.Context, patientID string) (*domain.Blacklist, error) {
return nil, nil // 默认不在黑名单
}
func (r *htMockBlacklistRepo) FindAll(_ context.Context, page, size int) ([]*domain.Blacklist, int64, error) {
return nil, 0, nil
}
func (r *htMockBlacklistRepo) FindExpired(_ context.Context) ([]*domain.Blacklist, error) {
return nil, nil
}

type htMockNoShowRepo struct{}

func (r *htMockNoShowRepo) Save(_ context.Context, ns *domain.NoShowRecord) error { return nil }
func (r *htMockNoShowRepo) CountByPatientIDInWindow(_ context.Context, patientID string, days int) (int, error) {
return 0, nil
}
func (r *htMockNoShowRepo) FindByPatientID(_ context.Context, patientID string) ([]*domain.NoShowRecord, error) {
return nil, nil
}

type htMockAppealRepo struct {
store map[string]*domain.Appeal
}

func newHtMockAppealRepo() *htMockAppealRepo {
return &htMockAppealRepo{store: make(map[string]*domain.Appeal)}
}
func (r *htMockAppealRepo) Save(_ context.Context, a *domain.Appeal) error {
r.store[a.ID] = a
return nil
}
func (r *htMockAppealRepo) Update(_ context.Context, a *domain.Appeal) error {
r.store[a.ID] = a
return nil
}
func (r *htMockAppealRepo) FindByID(_ context.Context, id string) (*domain.Appeal, error) {
a, ok := r.store[id]
if !ok {
return nil, errors.New(errors.ErrNotFound)
}
return a, nil
}
func (r *htMockAppealRepo) FindByBlacklistID(_ context.Context, blacklistID string) (*domain.Appeal, error) {
return nil, nil
}

// ---------------------------------------------------------------------------
// Test env helpers
// ---------------------------------------------------------------------------

type apptTestEnv struct {
router   *gin.Engine
apptRepo *htMockApptRepo
credRepo *htMockCredRepo
}

func newApptTestEnv() *apptTestEnv {
apptRepo := newHtMockApptRepo()
credRepo := newHtMockCredRepo()
svc := appAppt.NewAppointmentAppService(
apptRepo,
&htMockApptItemRepo{},
credRepo,
&htMockBlacklistRepo{},
&htMockNoShowRepo{},
newHtMockAppealRepo(),
)
h := NewHandler(svc)

r := gin.New()
r.Use(func(c *gin.Context) {
c.Set("user_id", "op-001")
c.Set("role", "admin")
c.Next()
})
v1 := r.Group("/api/v1")
h.RegisterRoutes(v1)

return &apptTestEnv{router: r, apptRepo: apptRepo, credRepo: credRepo}
}

func doApptRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseApptBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
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

func TestApptHandler_AutoAppointment_OK(t *testing.T) {
env := newApptTestEnv()
w := doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/auto", map[string]interface{}{
"patient_id":    "P-001",
"exam_item_ids": []string{"exam-001"},
})
if w.Code != http.StatusCreated {
t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
}
body := parseApptBody(t, w)
data := body["data"].(map[string]interface{})
// AutoAppointmentResp 返回 appointment_id 字段
if _, ok := data["appointment_id"]; !ok {
t.Errorf("expected appointment_id in response, got %v", data)
}
}

func TestApptHandler_AutoAppointment_BadRequest(t *testing.T) {
env := newApptTestEnv()
w := doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/auto", map[string]interface{}{})
if w.Code != http.StatusBadRequest {
t.Fatalf("expected 400, got %d", w.Code)
}
}

func TestApptHandler_GetAppointment_NotFound(t *testing.T) {
env := newApptTestEnv()
w := doApptRequest(env.router, http.MethodGet, "/api/v1/appointments/nonexistent", nil)
if w.Code != http.StatusNotFound {
t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
}
}

func TestApptHandler_GetAppointment_OK(t *testing.T) {
env := newApptTestEnv()
// 先创建预约单，AutoAppointmentResp 返回 appointment_id
wr := doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/auto", map[string]interface{}{
"patient_id":    "P-002",
"exam_item_ids": []string{"exam-002"},
})
if wr.Code != http.StatusCreated {
t.Fatalf("auto appointment failed: %s", wr.Body.String())
}
var resp map[string]interface{}
json.NewDecoder(wr.Body).Decode(&resp)
apptID := resp["data"].(map[string]interface{})["appointment_id"].(string)

w := doApptRequest(env.router, http.MethodGet, "/api/v1/appointments/"+apptID, nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
body := parseApptBody(t, w)
data := body["data"].(map[string]interface{})
if data["id"] != apptID {
t.Errorf("expected id=%s, got %v", apptID, data["id"])
}
}

func TestApptHandler_ListAppointments_OK(t *testing.T) {
env := newApptTestEnv()
doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/auto", map[string]interface{}{
"patient_id":    "P-003",
"exam_item_ids": []string{"exam-001"},
})
w := doApptRequest(env.router, http.MethodGet, "/api/v1/appointments", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestApptHandler_ManualAppointment_OK(t *testing.T) {
env := newApptTestEnv()
w := doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/manual", map[string]interface{}{
"patient_id":   "P-004",
"exam_item_id": "exam-001",
"slot_id":      "slot-001",
"reason":       "urgent case",
})
if w.Code != http.StatusCreated {
t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
}
}

func TestApptHandler_CancelAppointment_OK(t *testing.T) {
env := newApptTestEnv()
// 先创建预约
wr := doApptRequest(env.router, http.MethodPost, "/api/v1/appointments/auto", map[string]interface{}{
"patient_id":    "P-005",
"exam_item_ids": []string{"exam-001"},
})
if wr.Code != http.StatusCreated {
t.Fatalf("create failed: %s", wr.Body.String())
}
var resp map[string]interface{}
json.NewDecoder(wr.Body).Decode(&resp)
apptID := resp["data"].(map[string]interface{})["appointment_id"].(string)

// 确认预约（Pending -> Confirmed）
wc := doApptRequest(env.router, http.MethodPut, "/api/v1/appointments/"+apptID+"/confirm", map[string]interface{}{})
if wc.Code != http.StatusOK {
t.Fatalf("confirm failed: %d %s", wc.Code, wc.Body.String())
}

// 再取消（Confirmed -> Cancelled）
w := doApptRequest(env.router, http.MethodPut, "/api/v1/appointments/"+apptID+"/cancel", map[string]interface{}{
"reason": "patient request",
})
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestApptHandler_ListBlacklists_OK(t *testing.T) {
env := newApptTestEnv()
w := doApptRequest(env.router, http.MethodGet, "/api/v1/appointments/blacklist", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}
