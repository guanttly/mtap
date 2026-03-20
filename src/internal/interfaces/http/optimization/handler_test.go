package optimization

import (
"bytes"
"context"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
appOpt "github.com/euler/mtap/internal/application/optimization"
domain "github.com/euler/mtap/internal/domain/optimization"
"github.com/euler/mtap/pkg/errors"
)

func init() { gin.SetMode(gin.TestMode) }

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type htMockMetricRepo struct {
store map[string]*domain.EfficiencyMetric
}

func newHtMockMetricRepo() *htMockMetricRepo {
return &htMockMetricRepo{store: make(map[string]*domain.EfficiencyMetric)}
}
func (r *htMockMetricRepo) Save(_ context.Context, m *domain.EfficiencyMetric) error {
r.store[m.ID] = m
return nil
}
func (r *htMockMetricRepo) FindByID(_ context.Context, id string) (*domain.EfficiencyMetric, error) {
m, ok := r.store[id]
if !ok {
return nil, errors.New(errors.ErrNotFound)
}
return m, nil
}
func (r *htMockMetricRepo) FindByCode(_ context.Context, code string) (*domain.EfficiencyMetric, error) {
for _, m := range r.store {
if m.Code == code {
return m, nil
}
}
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockMetricRepo) List(_ context.Context) ([]*domain.EfficiencyMetric, error) {
var out []*domain.EfficiencyMetric
for _, m := range r.store {
out = append(out, m)
}
return out, nil
}
func (r *htMockMetricRepo) Update(_ context.Context, m *domain.EfficiencyMetric) error {
r.store[m.ID] = m
return nil
}

type htMockSnapshotRepo struct{}

func (r *htMockSnapshotRepo) Save(_ context.Context, s *domain.MetricSnapshot) error { return nil }
func (r *htMockSnapshotRepo) FindByMetricID(_ context.Context, metricID string, limit int) ([]*domain.MetricSnapshot, error) {
return nil, nil
}
func (r *htMockSnapshotRepo) FindRecent90Days(_ context.Context, metricID string) ([]*domain.MetricSnapshot, error) {
return nil, nil
}

type htMockAlertRepo struct {
store map[string]*domain.BottleneckAlert
}

func newHtMockAlertRepo() *htMockAlertRepo {
return &htMockAlertRepo{store: make(map[string]*domain.BottleneckAlert)}
}
func (r *htMockAlertRepo) Save(_ context.Context, a *domain.BottleneckAlert) error {
r.store[a.ID] = a
return nil
}
func (r *htMockAlertRepo) FindByID(_ context.Context, id string) (*domain.BottleneckAlert, error) {
a, ok := r.store[id]
if !ok {
return nil, nil // service 检查 nil 返回 ErrNotFound
}
return a, nil
}
func (r *htMockAlertRepo) List(_ context.Context, status string, page, size int) ([]*domain.BottleneckAlert, int64, error) {
var out []*domain.BottleneckAlert
for _, a := range r.store {
if status == "" || string(a.Status) == status {
out = append(out, a)
}
}
return out, int64(len(out)), nil
}
func (r *htMockAlertRepo) Update(_ context.Context, a *domain.BottleneckAlert) error {
r.store[a.ID] = a
return nil
}

type htMockStrategyRepo struct {
store map[string]*domain.OptimizationStrategy
}

func newHtMockStrategyRepo() *htMockStrategyRepo {
return &htMockStrategyRepo{store: make(map[string]*domain.OptimizationStrategy)}
}
func (r *htMockStrategyRepo) Save(_ context.Context, s *domain.OptimizationStrategy) error {
r.store[s.ID] = s
return nil
}
func (r *htMockStrategyRepo) FindByID(_ context.Context, id string) (*domain.OptimizationStrategy, error) {
s, ok := r.store[id]
if !ok {
return nil, nil // service 检查 nil 返回 ErrNotFound
}
return s, nil
}
func (r *htMockStrategyRepo) List(_ context.Context, category, status string, page, size int) ([]*domain.OptimizationStrategy, int64, error) {
var out []*domain.OptimizationStrategy
for _, s := range r.store {
out = append(out, s)
}
return out, int64(len(out)), nil
}
func (r *htMockStrategyRepo) Update(_ context.Context, s *domain.OptimizationStrategy) error {
r.store[s.ID] = s
return nil
}
func (r *htMockStrategyRepo) ListActiveTrials(_ context.Context) ([]*domain.OptimizationStrategy, error) {
return nil, nil
}
func (r *htMockStrategyRepo) ListPromoted(_ context.Context) ([]*domain.OptimizationStrategy, error) {
return nil, nil
}
func (r *htMockStrategyRepo) CountPendingReview(_ context.Context) (int64, error) {
return 0, nil
}

type htMockTrialRepo struct{}

func (r *htMockTrialRepo) Save(_ context.Context, t *domain.TrialRun) error { return nil }
func (r *htMockTrialRepo) FindByID(_ context.Context, id string) (*domain.TrialRun, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockTrialRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.TrialRun, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockTrialRepo) Update(_ context.Context, t *domain.TrialRun) error { return nil }

type htMockEvalRepo struct{}

func (r *htMockEvalRepo) Save(_ context.Context, e *domain.EvaluationReport) error { return nil }
func (r *htMockEvalRepo) FindByID(_ context.Context, id string) (*domain.EvaluationReport, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockEvalRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.EvaluationReport, error) {
return nil, errors.New(errors.ErrNotFound)
}

type htMockROIRepo struct{}

func (r *htMockROIRepo) Save(_ context.Context, roi *domain.ROIReport) error { return nil }
func (r *htMockROIRepo) FindByID(_ context.Context, id string) (*domain.ROIReport, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockROIRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.ROIReport, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockROIRepo) Update(_ context.Context, roi *domain.ROIReport) error { return nil }

type htMockScanRepo struct{}

func (r *htMockScanRepo) Save(_ context.Context, s *domain.PerformanceScan) error { return nil }
func (r *htMockScanRepo) FindByID(_ context.Context, id string) (*domain.PerformanceScan, error) {
return nil, errors.New(errors.ErrNotFound)
}
func (r *htMockScanRepo) List(_ context.Context, page, size int) ([]*domain.PerformanceScan, int64, error) {
return nil, 0, nil
}

type htMockDecayRepo struct{}

func (r *htMockDecayRepo) Save(_ context.Context, a *domain.StrategyDecayAlert) error { return nil }
func (r *htMockDecayRepo) FindByStrategyID(_ context.Context, strategyID string) ([]*domain.StrategyDecayAlert, error) {
return nil, nil
}

// ---------------------------------------------------------------------------
// Test env helpers
// ---------------------------------------------------------------------------

type optTestEnv struct {
router      *gin.Engine
metricRepo  *htMockMetricRepo
alertRepo   *htMockAlertRepo
strategyRepo *htMockStrategyRepo
}

func newOptTestEnv() *optTestEnv {
metricRepo := newHtMockMetricRepo()
alertRepo := newHtMockAlertRepo()
strategyRepo := newHtMockStrategyRepo()

svc := appOpt.NewOptimizationAppService(
metricRepo,
&htMockSnapshotRepo{},
alertRepo,
strategyRepo,
&htMockTrialRepo{},
&htMockEvalRepo{},
&htMockROIRepo{},
&htMockScanRepo{},
&htMockDecayRepo{},
)
h := NewHandler(svc)

r := gin.New()
r.Use(func(c *gin.Context) {
c.Set("user_id", "admin-001")
c.Set("role", "admin")
c.Next()
})
v1 := r.Group("/api/v1")
h.RegisterRoutes(v1)

return &optTestEnv{router: r, metricRepo: metricRepo, alertRepo: alertRepo, strategyRepo: strategyRepo}
}

func doOptRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
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

func parseOptBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
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

func TestOptHandler_ListMetrics_Empty(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/metrics", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestOptHandler_ListMetrics_WithData(t *testing.T) {
env := newOptTestEnv()
// 预先放入指标数据
env.metricRepo.store["m-001"] = &domain.EfficiencyMetric{
ID: "m-001", Name: "设备利用率", Code: "device_util",
Unit: "%", NormalMean: 80, NormalStdDev: 5,
CreatedAt: time.Now(), UpdatedAt: time.Now(),
}
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/metrics", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
body := parseOptBody(t, w)
data, ok := body["data"].([]interface{})
if !ok || len(data) == 0 {
t.Errorf("expected non-empty metric list, got %v", body["data"])
}
}

func TestOptHandler_ListAlerts_Empty(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/alerts", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestOptHandler_GetAlert_NotFound(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/alerts/nonexistent", nil)
if w.Code != http.StatusNotFound {
t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
}
}

func TestOptHandler_ListStrategies_Empty(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/strategies", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}

func TestOptHandler_GetStrategy_NotFound(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/strategies/nonexistent", nil)
if w.Code != http.StatusNotFound {
t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
}
}

func TestOptHandler_ListScans_Empty(t *testing.T) {
env := newOptTestEnv()
w := doOptRequest(env.router, http.MethodGet, "/api/v1/optimization/scans", nil)
if w.Code != http.StatusOK {
t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
}
}
