// Package analytics HTTP handler 集成测试
package analytics

import (
"bytes"
"context"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

appAnalytics "github.com/euler/mtap/internal/application/analytics"
domain "github.com/euler/mtap/internal/domain/analytics"
)

// ─── Mock 仓库 ────────────────────────────────────────────────────────────────

type htMockDashboardRepo struct{}

func (m *htMockDashboardRepo) GetSlotUsage(_ context.Context, _ string, _ time.Time) (domain.SlotUsageData, error) {
return domain.SlotUsageData{TotalSlots: 100, UsedSlots: 60, AvailableSlots: 35}, nil
}
func (m *htMockDashboardRepo) GetDeviceStatus(_ context.Context, _ string) ([]domain.DeviceStatusData, error) {
return []domain.DeviceStatusData{
{DeviceID: "D001", DeviceName: "CT-1", Status: "idle", QueueCount: 0},
}, nil
}
func (m *htMockDashboardRepo) GetWaitTrend(_ context.Context, _ string, _ int) ([]domain.WaitTrendPoint, error) {
return []domain.WaitTrendPoint{
{Time: time.Now(), AvgWaitMin: 10},
}, nil
}
func (m *htMockDashboardRepo) SaveSnapshot(_ context.Context, _ *domain.DashboardSnapshot) error {
return nil
}
func (m *htMockDashboardRepo) GetDeviceDetail(_ context.Context, deviceID string, _ time.Time) (*domain.DeviceDetail, error) {
if deviceID == "not-exist" {
return nil, nil
}
return &domain.DeviceDetail{
DeviceID:   deviceID,
DeviceName: "CT-1",
TimeSlots:  []domain.SlotSummary{{Hour: 9, Total: 10, Used: 6}},
}, nil
}

type htMockReportRepo struct{ reports map[string]*domain.Report }

func newHtReportRepo() *htMockReportRepo {
return &htMockReportRepo{reports: make(map[string]*domain.Report)}
}
func (m *htMockReportRepo) Save(_ context.Context, r *domain.Report) error {
m.reports[r.ID] = r
return nil
}
func (m *htMockReportRepo) FindByID(_ context.Context, id string) (*domain.Report, error) {
return m.reports[id], nil
}
func (m *htMockReportRepo) List(_ context.Context, _, _ int) ([]*domain.Report, int64, error) {
var out []*domain.Report
for _, r := range m.reports {
out = append(out, r)
}
return out, int64(len(out)), nil
}
func (m *htMockReportRepo) Update(_ context.Context, r *domain.Report) error {
m.reports[r.ID] = r
return nil
}

// ─── 辅助函数 ─────────────────────────────────────────────────────────────────

func newTestAnalyticsHandler() *Handler {
dashRepo := &htMockDashboardRepo{}
reportRepo := newHtReportRepo()
svc := appAnalytics.NewAnalyticsAppService(dashRepo, reportRepo)
return NewHandler(svc)
}

func setupAnalyticsRouter(h *Handler) *gin.Engine {
gin.SetMode(gin.TestMode)
r := gin.New()
r.Use(func(c *gin.Context) {
c.Set("user_id", "test-admin")
c.Set("role", "admin")
c.Next()
})
api := r.Group("/api/v1")
h.RegisterRoutes(api)
return r
}

func doAnalyticsRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseAnalyticsBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
t.Helper()
var m map[string]interface{}
require.NoError(t, json.NewDecoder(w.Body).Decode(&m))
return m
}

// ─── 测试用例 ─────────────────────────────────────────────────────────────────

func TestAnalyticsHandler_GetDashboard_OK(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

w := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/dashboard", nil)
assert.Equal(t, http.StatusOK, w.Code)
body := parseAnalyticsBody(t, w)
assert.Equal(t, float64(0), body["code"])
data := body["data"].(map[string]interface{})
assert.NotNil(t, data["slot_usage"])
}

func TestAnalyticsHandler_GetDeviceDetail_OK(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

w := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/dashboard/device/D001?date=2024-01-15", nil)
assert.Equal(t, http.StatusOK, w.Code)
body := parseAnalyticsBody(t, w)
assert.Equal(t, float64(0), body["code"])
}

func TestAnalyticsHandler_GetDeviceDetail_BadDate(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

// 非法日期格式 → 应返回非 200
w := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/dashboard/device/D001?date=20240115", nil)
assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_CreateReport_OK(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

w := doAnalyticsRequest(r, http.MethodPost, "/api/v1/analytics/reports", map[string]interface{}{
"report_type": "daily",
"date_start":  "2024-01-01",
"date_end":    "2024-01-31",
})
assert.Equal(t, http.StatusCreated, w.Code)
body := parseAnalyticsBody(t, w)
assert.Equal(t, float64(0), body["code"])
data := body["data"].(map[string]interface{})
assert.NotEmpty(t, data["id"])
}

func TestAnalyticsHandler_CreateReport_BadRequest(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

// report_type 缺失 → binding 验证失败
w := doAnalyticsRequest(r, http.MethodPost, "/api/v1/analytics/reports", map[string]interface{}{
"date_start": "2024-01-01",
"date_end":   "2024-01-31",
})
assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsHandler_GetReport_OK(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

// 先创建
w1 := doAnalyticsRequest(r, http.MethodPost, "/api/v1/analytics/reports", map[string]interface{}{
"report_type": "weekly",
"date_start":  "2024-01-01",
"date_end":    "2024-01-07",
})
require.Equal(t, http.StatusCreated, w1.Code)
id := parseAnalyticsBody(t, w1)["data"].(map[string]interface{})["id"].(string)

w2 := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/reports/"+id, nil)
assert.Equal(t, http.StatusOK, w2.Code)
body := parseAnalyticsBody(t, w2)
assert.Equal(t, float64(0), body["code"])
}

func TestAnalyticsHandler_GetReport_NotFound(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

w := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/reports/non-existent", nil)
assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAnalyticsHandler_ListReports_OK(t *testing.T) {
h := newTestAnalyticsHandler()
r := setupAnalyticsRouter(h)

doAnalyticsRequest(r, http.MethodPost, "/api/v1/analytics/reports", map[string]interface{}{
"report_type": "monthly",
"date_start":  "2024-01-01",
"date_end":    "2024-01-31",
})

w := doAnalyticsRequest(r, http.MethodGet, "/api/v1/analytics/reports?page=1&page_size=10", nil)
assert.Equal(t, http.StatusOK, w.Code)
body := parseAnalyticsBody(t, w)
assert.Equal(t, float64(0), body["code"])
}
