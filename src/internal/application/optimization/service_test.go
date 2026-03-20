package optimization

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/euler/mtap/internal/domain/optimization"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// ============================================================
// Mock Repositories
// ============================================================

type mockMetricRepo struct {
	byID   map[string]*domain.EfficiencyMetric
	byCode map[string]*domain.EfficiencyMetric
}

func newMockMetricRepo() *mockMetricRepo {
	return &mockMetricRepo{
		byID:   make(map[string]*domain.EfficiencyMetric),
		byCode: make(map[string]*domain.EfficiencyMetric),
	}
}
func (m *mockMetricRepo) Save(_ context.Context, me *domain.EfficiencyMetric) error {
	m.byID[me.ID] = me
	m.byCode[me.Code] = me
	return nil
}
func (m *mockMetricRepo) FindByID(_ context.Context, id string) (*domain.EfficiencyMetric, error) {
	return m.byID[id], nil
}
func (m *mockMetricRepo) FindByCode(_ context.Context, code string) (*domain.EfficiencyMetric, error) {
	return m.byCode[code], nil
}
func (m *mockMetricRepo) List(_ context.Context) ([]*domain.EfficiencyMetric, error) {
	all := make([]*domain.EfficiencyMetric, 0, len(m.byID))
	for _, v := range m.byID {
		all = append(all, v)
	}
	return all, nil
}
func (m *mockMetricRepo) Update(_ context.Context, me *domain.EfficiencyMetric) error {
	m.byID[me.ID] = me
	m.byCode[me.Code] = me
	return nil
}

type mockSnapshotRepo struct {
	items map[string][]*domain.MetricSnapshot
}

func newMockSnapshotRepo() *mockSnapshotRepo {
	return &mockSnapshotRepo{items: make(map[string][]*domain.MetricSnapshot)}
}
func (m *mockSnapshotRepo) Save(_ context.Context, s *domain.MetricSnapshot) error {
	m.items[s.MetricID] = append(m.items[s.MetricID], s)
	return nil
}
func (m *mockSnapshotRepo) FindByMetricID(_ context.Context, metricID string, limit int) ([]*domain.MetricSnapshot, error) {
	snaps := m.items[metricID]
	if limit > 0 && len(snaps) > limit {
		snaps = snaps[:limit]
	}
	return snaps, nil
}
func (m *mockSnapshotRepo) FindRecent90Days(_ context.Context, metricID string) ([]*domain.MetricSnapshot, error) {
	return m.items[metricID], nil
}

type mockAlertRepo struct {
	items map[string]*domain.BottleneckAlert
}

func newMockAlertRepo() *mockAlertRepo {
	return &mockAlertRepo{items: make(map[string]*domain.BottleneckAlert)}
}
func (m *mockAlertRepo) Save(_ context.Context, a *domain.BottleneckAlert) error {
	m.items[a.ID] = a
	return nil
}
func (m *mockAlertRepo) FindByID(_ context.Context, id string) (*domain.BottleneckAlert, error) {
	return m.items[id], nil
}
func (m *mockAlertRepo) List(_ context.Context, status string, page, size int) ([]*domain.BottleneckAlert, int64, error) {
	var all []*domain.BottleneckAlert
	for _, a := range m.items {
		if status == "" || string(a.Status) == status {
			all = append(all, a)
		}
	}
	total := int64(len(all))
	start := (page - 1) * size
	if start >= len(all) {
		return []*domain.BottleneckAlert{}, total, nil
	}
	end := start + size
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}
func (m *mockAlertRepo) Update(_ context.Context, a *domain.BottleneckAlert) error {
	m.items[a.ID] = a
	return nil
}

type mockStrategyRepo struct {
	items map[string]*domain.OptimizationStrategy
}

func newMockStrategyRepo() *mockStrategyRepo {
	return &mockStrategyRepo{items: make(map[string]*domain.OptimizationStrategy)}
}
func (m *mockStrategyRepo) Save(_ context.Context, s *domain.OptimizationStrategy) error {
	m.items[s.ID] = s
	return nil
}
func (m *mockStrategyRepo) FindByID(_ context.Context, id string) (*domain.OptimizationStrategy, error) {
	return m.items[id], nil
}
func (m *mockStrategyRepo) List(_ context.Context, category, status string, page, size int) ([]*domain.OptimizationStrategy, int64, error) {
	var all []*domain.OptimizationStrategy
	for _, s := range m.items {
		if (category == "" || string(s.Category) == category) && (status == "" || string(s.Status) == status) {
			all = append(all, s)
		}
	}
	total := int64(len(all))
	start := (page - 1) * size
	if start >= len(all) {
		return []*domain.OptimizationStrategy{}, total, nil
	}
	end := start + size
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}
func (m *mockStrategyRepo) Update(_ context.Context, s *domain.OptimizationStrategy) error {
	m.items[s.ID] = s
	return nil
}
func (m *mockStrategyRepo) ListActiveTrials(_ context.Context) ([]*domain.OptimizationStrategy, error) {
	var result []*domain.OptimizationStrategy
	for _, s := range m.items {
		if s.Status == domain.StatusTrialRunning || s.Status == domain.StatusTrialRunningB {
			result = append(result, s)
		}
	}
	return result, nil
}
func (m *mockStrategyRepo) ListPromoted(_ context.Context) ([]*domain.OptimizationStrategy, error) {
	var result []*domain.OptimizationStrategy
	for _, s := range m.items {
		if s.Status == domain.StatusPromoted || s.Status == domain.StatusTracking {
			result = append(result, s)
		}
	}
	return result, nil
}
func (m *mockStrategyRepo) CountPendingReview(_ context.Context) (int64, error) {
	var count int64
	for _, s := range m.items {
		if s.Status == domain.StatusPendingReview {
			count++
		}
	}
	return count, nil
}

type mockTrialRepo struct {
	items   map[string]*domain.TrialRun
	byStrat map[string]*domain.TrialRun
}

func newMockTrialRepo() *mockTrialRepo {
	return &mockTrialRepo{
		items:   make(map[string]*domain.TrialRun),
		byStrat: make(map[string]*domain.TrialRun),
	}
}
func (m *mockTrialRepo) Save(_ context.Context, t *domain.TrialRun) error {
	m.items[t.ID] = t
	m.byStrat[t.StrategyID] = t
	return nil
}
func (m *mockTrialRepo) FindByID(_ context.Context, id string) (*domain.TrialRun, error) {
	return m.items[id], nil
}
func (m *mockTrialRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.TrialRun, error) {
	return m.byStrat[strategyID], nil
}
func (m *mockTrialRepo) Update(_ context.Context, t *domain.TrialRun) error {
	m.items[t.ID] = t
	m.byStrat[t.StrategyID] = t
	return nil
}

type mockEvalRepo struct {
	items   map[string]*domain.EvaluationReport
	byStrat map[string]*domain.EvaluationReport
}

func newMockEvalRepo() *mockEvalRepo {
	return &mockEvalRepo{
		items:   make(map[string]*domain.EvaluationReport),
		byStrat: make(map[string]*domain.EvaluationReport),
	}
}
func (m *mockEvalRepo) Save(_ context.Context, r *domain.EvaluationReport) error {
	m.items[r.ID] = r
	m.byStrat[r.StrategyID] = r
	return nil
}
func (m *mockEvalRepo) FindByID(_ context.Context, id string) (*domain.EvaluationReport, error) {
	return m.items[id], nil
}
func (m *mockEvalRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.EvaluationReport, error) {
	return m.byStrat[strategyID], nil
}

type mockROIRepo struct {
	items   map[string]*domain.ROIReport
	byStrat map[string]*domain.ROIReport
}

func newMockROIRepo() *mockROIRepo {
	return &mockROIRepo{
		items:   make(map[string]*domain.ROIReport),
		byStrat: make(map[string]*domain.ROIReport),
	}
}
func (m *mockROIRepo) Save(_ context.Context, r *domain.ROIReport) error {
	m.items[r.ID] = r
	m.byStrat[r.StrategyID] = r
	return nil
}
func (m *mockROIRepo) FindByID(_ context.Context, id string) (*domain.ROIReport, error) {
	return m.items[id], nil
}
func (m *mockROIRepo) FindByStrategyID(_ context.Context, strategyID string) (*domain.ROIReport, error) {
	return m.byStrat[strategyID], nil
}
func (m *mockROIRepo) Update(_ context.Context, r *domain.ROIReport) error {
	m.items[r.ID] = r
	m.byStrat[r.StrategyID] = r
	return nil
}

type mockScanRepo struct {
	items map[string]*domain.PerformanceScan
}

func newMockScanRepo() *mockScanRepo {
	return &mockScanRepo{items: make(map[string]*domain.PerformanceScan)}
}
func (m *mockScanRepo) Save(_ context.Context, s *domain.PerformanceScan) error {
	m.items[s.ID] = s
	return nil
}
func (m *mockScanRepo) FindByID(_ context.Context, id string) (*domain.PerformanceScan, error) {
	return m.items[id], nil
}
func (m *mockScanRepo) List(_ context.Context, page, size int) ([]*domain.PerformanceScan, int64, error) {
	all := make([]*domain.PerformanceScan, 0, len(m.items))
	for _, s := range m.items {
		all = append(all, s)
	}
	total := int64(len(all))
	start := (page - 1) * size
	if start >= len(all) {
		return []*domain.PerformanceScan{}, total, nil
	}
	end := start + size
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

type mockDecayRepo struct {
	byStrat map[string][]*domain.StrategyDecayAlert
}

func newMockDecayRepo() *mockDecayRepo {
	return &mockDecayRepo{byStrat: make(map[string][]*domain.StrategyDecayAlert)}
}
func (m *mockDecayRepo) Save(_ context.Context, a *domain.StrategyDecayAlert) error {
	m.byStrat[a.StrategyID] = append(m.byStrat[a.StrategyID], a)
	return nil
}
func (m *mockDecayRepo) FindByStrategyID(_ context.Context, strategyID string) ([]*domain.StrategyDecayAlert, error) {
	return m.byStrat[strategyID], nil
}

// ============================================================
// Helper
// ============================================================

func newTestOptimizationService() *OptimizationAppService {
	return NewOptimizationAppService(
		newMockMetricRepo(),
		newMockSnapshotRepo(),
		newMockAlertRepo(),
		newMockStrategyRepo(),
		newMockTrialRepo(),
		newMockEvalRepo(),
		newMockROIRepo(),
		newMockScanRepo(),
		newMockDecayRepo(),
	)
}

func seedMetric(svc *OptimizationAppService, id, code, name string) {
	m := &domain.EfficiencyMetric{ID: id, Code: code, Name: name, Unit: "%"}
	_ = svc.metricRepo.Save(context.Background(), m)
}

func seedAlert(svc *OptimizationAppService) *domain.BottleneckAlert {
	// 预置指标，避免归因分析时 nil pointer
	seedMetric(svc, "M001", "avg_wait_min", "平均等待时间")
	a := &domain.BottleneckAlert{
		ID:        uuid.New().String(),
		MetricID:  "M001",
		AlertType: domain.AlertConsecutiveDeviation,
		Severity:  "warning",
		Status:    domain.AlertActive,
		// 预置假设，避免 GetAlert 触发归因分析（归因分析需要数据库中有真实数据）
		RootCauseHypotheses: []domain.Hypothesis{
			{RootCause: "测试假设", Confidence: 1.0},
		},
	}
	_ = svc.alertRepo.Save(context.Background(), a)
	return a
}

func seedStrategy(svc *OptimizationAppService) *domain.OptimizationStrategy {
	st := &domain.OptimizationStrategy{
		ID:       uuid.New().String(),
		Title:    "优化CT排班间隔",
		Category: domain.CategoryA,
		Status:   domain.StatusPendingReview,
		ApprovalFlow: domain.ApprovalFlow{
			Type: "single",
			Approvers: []domain.ApprovalNode{
				{ApproverID: "approver001", ApproverRole: "scheduler_admin", Status: "pending"},
			},
			Status: "pending",
		},
	}
	_ = svc.strategyRepo.Save(context.Background(), st)
	return st
}

// ============================================================
// Metrics Tests
// ============================================================

func TestListMetrics_Empty(t *testing.T) {
	svc := newTestOptimizationService()
	list, err := svc.ListMetrics(context.Background())
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestListMetrics_WithData(t *testing.T) {
	svc := newTestOptimizationService()
	seedMetric(svc, "M001", "slot_usage_rate", "号源利用率")
	seedMetric(svc, "M002", "avg_wait_min", "平均等待时间")

	list, err := svc.ListMetrics(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetMetricTrend_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetMetricTrend(context.Background(), "nonexistent_code", 7)
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetMetricTrend_OK(t *testing.T) {
	svc := newTestOptimizationService()
	seedMetric(svc, "M001", "slot_usage_rate", "号源利用率")

	list, err := svc.GetMetricTrend(context.Background(), "slot_usage_rate", 7)
	require.NoError(t, err)
	assert.Empty(t, list) // 无快照数据
}

// ============================================================
// Alerts Tests
// ============================================================

func TestListAlerts_OK(t *testing.T) {
	svc := newTestOptimizationService()
	seedAlert(svc)
	seedAlert(svc)

	list, total, err := svc.ListAlerts(context.Background(), ListAlertsReq{Page: 1, Size: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

func TestGetAlert_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetAlert(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetAlert_OK(t *testing.T) {
	svc := newTestOptimizationService()
	alert := seedAlert(svc)

	resp, err := svc.GetAlert(context.Background(), alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alert.ID, resp.ID)
	assert.Equal(t, "active", resp.Status)
}

func TestDismissAlert_OK(t *testing.T) {
	svc := newTestOptimizationService()
	alert := seedAlert(svc)

	err := svc.DismissAlert(context.Background(), alert.ID, DismissAlertReq{Reason: "误报，已确认正常"})
	require.NoError(t, err)

	resp, _ := svc.GetAlert(context.Background(), alert.ID)
	assert.Equal(t, "dismissed", resp.Status)
}

func TestDismissAlert_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	err := svc.DismissAlert(context.Background(), "nonexistent", DismissAlertReq{Reason: "test"})
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// Strategy Tests
// ============================================================

func TestListStrategies_OK(t *testing.T) {
	svc := newTestOptimizationService()
	seedStrategy(svc)
	seedStrategy(svc)

	list, total, err := svc.ListStrategies(context.Background(), ListStrategiesReq{Page: 1, Size: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

func TestGetStrategy_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetStrategy(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetStrategy_OK(t *testing.T) {
	svc := newTestOptimizationService()
	st := seedStrategy(svc)

	resp, err := svc.GetStrategy(context.Background(), st.ID)
	require.NoError(t, err)
	assert.Equal(t, st.ID, resp.ID)
	assert.Equal(t, "pending_review", resp.Status)
}

func TestApproveStrategy_OK(t *testing.T) {
	svc := newTestOptimizationService()
	st := seedStrategy(svc)

	resp, err := svc.ApproveStrategy(context.Background(), st.ID, "approver001", ApproveStrategyReq{
		TrialDays: 7,
		GrayScope: GrayScopeDTO{DeviceIDs: []string{"DEV001"}},
	})
	require.NoError(t, err)
	// A类策略审批后进入试运行
	assert.Equal(t, "trial_running", resp.Status)
}

func TestRejectStrategy_OK(t *testing.T) {
	svc := newTestOptimizationService()
	st := seedStrategy(svc)

	err := svc.RejectStrategy(context.Background(), st.ID, RejectStrategyReq{Reason: "配置变更影响过大，暂缓实施"})
	require.NoError(t, err)

	resp, _ := svc.GetStrategy(context.Background(), st.ID)
	assert.Equal(t, "rejected", resp.Status)
}

func TestRollbackStrategy_OK(t *testing.T) {
	svc := newTestOptimizationService()
	st := seedStrategy(svc)

	// 先审批让策略进入试运行状态
	_, err := svc.ApproveStrategy(context.Background(), st.ID, "approver001", ApproveStrategyReq{TrialDays: 7})
	require.NoError(t, err)

	// 手动回滚
	err = svc.RollbackStrategy(context.Background(), st.ID)
	require.NoError(t, err)

	resp, _ := svc.GetStrategy(context.Background(), st.ID)
	assert.Equal(t, "rolled_back", resp.Status)
}

func TestRejectStrategy_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	err := svc.RejectStrategy(context.Background(), "nonexistent", RejectStrategyReq{Reason: "not found"})
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// TrialMonitor / Evaluation / ROI Tests
// ============================================================

func TestGetTrialMonitor_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetTrialMonitor(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetEvaluation_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetEvaluation(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestGetROIReport_NotFound(t *testing.T) {
	svc := newTestOptimizationService()
	_, err := svc.GetROIReport(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// ListScans Tests
// ============================================================

func TestListScans_Empty(t *testing.T) {
	svc := newTestOptimizationService()
	list, total, err := svc.ListScans(context.Background(), ListScansReq{Page: 1, Size: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, list)
}
