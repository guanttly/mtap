// Package optimization 应用层 - 效能优化应用服务
package optimization

import (
	"context"
	"time"

	"github.com/google/uuid"

	domain "github.com/euler/mtap/internal/domain/optimization"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// OptimizationAppService 效能优化应用服务
type OptimizationAppService struct {
	metricRepo   domain.EfficiencyMetricRepository
	snapshotRepo domain.MetricSnapshotRepository
	alertRepo    domain.BottleneckAlertRepository
	strategyRepo domain.OptimizationStrategyRepository
	trialRepo    domain.TrialRunRepository
	evalRepo     domain.EvaluationReportRepository
	roiRepo      domain.ROIReportRepository
	scanRepo     domain.PerformanceScanRepository
	decayRepo    domain.StrategyDecayAlertRepository

	anomalySvc  *domain.AnomalyDetectionService
	attrSvc     *domain.BottleneckAttributionService
	strategySvc *domain.StrategyGenerationService
	evalSvc     *domain.EvaluationService
	decaySvc    *domain.DecayTrackingService
}

// NewOptimizationAppService 创建应用服务
func NewOptimizationAppService(
	metricRepo domain.EfficiencyMetricRepository,
	snapshotRepo domain.MetricSnapshotRepository,
	alertRepo domain.BottleneckAlertRepository,
	strategyRepo domain.OptimizationStrategyRepository,
	trialRepo domain.TrialRunRepository,
	evalRepo domain.EvaluationReportRepository,
	roiRepo domain.ROIReportRepository,
	scanRepo domain.PerformanceScanRepository,
	decayRepo domain.StrategyDecayAlertRepository,
) *OptimizationAppService {
	attrSvc := domain.NewBottleneckAttributionService(metricRepo, alertRepo)
	svc := &OptimizationAppService{
		metricRepo: metricRepo, snapshotRepo: snapshotRepo,
		alertRepo: alertRepo, strategyRepo: strategyRepo,
		trialRepo: trialRepo, evalRepo: evalRepo,
		roiRepo: roiRepo, scanRepo: scanRepo, decayRepo: decayRepo,
		anomalySvc:  domain.NewAnomalyDetectionService(metricRepo, snapshotRepo, alertRepo),
		attrSvc:     attrSvc,
		strategySvc: domain.NewStrategyGenerationService(alertRepo, strategyRepo, attrSvc),
		evalSvc:     domain.NewEvaluationService(strategyRepo, trialRepo, evalRepo, snapshotRepo, metricRepo),
		decaySvc:    domain.NewDecayTrackingService(strategyRepo, evalRepo, decayRepo, snapshotRepo, metricRepo),
	}
	return svc
}

// ─── 指标 ───────────────────────────────────────────────────────────────────

// ListMetrics 获取效率指标列表
func (s *OptimizationAppService) ListMetrics(ctx context.Context) ([]*MetricResp, error) {
	metrics, err := s.metricRepo.List(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	resps := make([]*MetricResp, 0, len(metrics))
	for _, m := range metrics {
		resps = append(resps, ToMetricResp(m))
	}
	return resps, nil
}

// GetMetricTrend 获取指标趋势
func (s *OptimizationAppService) GetMetricTrend(ctx context.Context, code string, days int) ([]*MetricSnapshotResp, error) {
	if days <= 0 {
		days = 7
	}
	metric, err := s.metricRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if metric == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	snaps, err := s.snapshotRepo.FindByMetricID(ctx, metric.ID, days*24)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	resps := make([]*MetricSnapshotResp, 0, len(snaps))
	for _, snap := range snaps {
		resps = append(resps, &MetricSnapshotResp{Value: snap.Value, SampledAt: snap.SampledAt})
	}
	return resps, nil
}

// ─── 告警 ───────────────────────────────────────────────────────────────────

// ListAlerts 获取告警列表
func (s *OptimizationAppService) ListAlerts(ctx context.Context, req ListAlertsReq) ([]*AlertResp, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	alerts, total, err := s.alertRepo.List(ctx, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]*AlertResp, 0, len(alerts))
	for _, a := range alerts {
		resps = append(resps, ToAlertResp(a))
	}
	return resps, total, nil
}

// GetAlert 获取告警详情（含归因分析）
func (s *OptimizationAppService) GetAlert(ctx context.Context, id string) (*AlertResp, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if alert == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	// 尝试归因分析（如果还未归因）
	if len(alert.RootCauseHypotheses) == 0 {
		_, _ = s.attrSvc.Analyze(ctx, id)
		alert, _ = s.alertRepo.FindByID(ctx, id)
	}
	return ToAlertResp(alert), nil
}

// DismissAlert 标记告警误报
func (s *OptimizationAppService) DismissAlert(ctx context.Context, id string, req DismissAlertReq) error {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if alert == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := alert.Dismiss(req.Reason); err != nil {
		return err
	}
	return s.alertRepo.Update(ctx, alert)
}

// ─── 策略 ───────────────────────────────────────────────────────────────────

// ListStrategies 获取策略列表
func (s *OptimizationAppService) ListStrategies(ctx context.Context, req ListStrategiesReq) ([]*StrategyResp, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	strategies, total, err := s.strategyRepo.List(ctx, req.Category, req.Status, req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]*StrategyResp, 0, len(strategies))
	for _, st := range strategies {
		resps = append(resps, ToStrategyResp(st))
	}
	return resps, total, nil
}

// GetStrategy 获取策略详情
func (s *OptimizationAppService) GetStrategy(ctx context.Context, id string) (*StrategyResp, error) {
	st, err := s.strategyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if st == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToStrategyResp(st), nil
}

// ApproveStrategy 审批策略（A类启动试运行，C类提交院级）
func (s *OptimizationAppService) ApproveStrategy(ctx context.Context, id, approverID string, req ApproveStrategyReq) (*StrategyResp, error) {
	st, err := s.strategyRepo.FindByID(ctx, id)
	if err != nil || st == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	// 设置审批者 ID
	for i := range st.ApprovalFlow.Approvers {
		if st.ApprovalFlow.Approvers[i].Status == "pending" && st.ApprovalFlow.Approvers[i].ApproverID == "" {
			st.ApprovalFlow.Approvers[i].ApproverID = approverID
			break
		}
	}
	if err := st.Approve(approverID); err != nil {
		return nil, err
	}
	if err := s.strategyRepo.Update(ctx, st); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	// A/B类策略审批通过后创建试运行记录
	if st.Status == domain.StatusTrialRunning || st.Status == domain.StatusTrialRunningB {
		trial := &domain.TrialRun{
			ID: uuid.New().String(), StrategyID: st.ID,
			GrayScope: domain.GrayScope{
				DepartmentIDs: req.GrayScope.DepartmentIDs,
				DeviceIDs:     req.GrayScope.DeviceIDs,
				TimePeriods:   req.GrayScope.TimePeriods,
			},
			TrialDays:                  req.TrialDays,
			Baseline:                   domain.BaselineSnapshot{Metrics: map[string]float64{}, SnapshotAt: time.Now()},
			EmergencyRollbackThreshold: 0.15,
			StartedAt:                  time.Now(),
			EndsAt:                     time.Now().AddDate(0, 0, req.TrialDays),
			Status:                     domain.TrialRunning,
		}
		_ = s.trialRepo.Save(ctx, trial)
	}
	return ToStrategyResp(st), nil
}

// RejectStrategy 驳回策略
func (s *OptimizationAppService) RejectStrategy(ctx context.Context, id string, req RejectStrategyReq) error {
	st, err := s.strategyRepo.FindByID(ctx, id)
	if err != nil || st == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := st.Reject(req.Reason); err != nil {
		return err
	}
	return s.strategyRepo.Update(ctx, st)
}

// RollbackStrategy 手动回滚策略
func (s *OptimizationAppService) RollbackStrategy(ctx context.Context, id string) error {
	st, err := s.strategyRepo.FindByID(ctx, id)
	if err != nil || st == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := st.Rollback(); err != nil {
		return err
	}
	if err := s.strategyRepo.Update(ctx, st); err != nil {
		return err
	}
	// 更新试运行状态
	trial, _ := s.trialRepo.FindByStrategyID(ctx, id)
	if trial != nil {
		now := time.Now()
		trial.Status = domain.TrialRolledBack
		trial.RolledBackAt = &now
		trial.RollbackReason = "手动回滚"
		_ = s.trialRepo.Update(ctx, trial)
	}
	return nil
}

// PromoteStrategy 策略转正
func (s *OptimizationAppService) PromoteStrategy(ctx context.Context, id string) error {
	return s.evalSvc.PromoteStrategy(ctx, id)
}

// ─── 试运行与评估 ───────────────────────────────────────────────────────────

// GetTrialMonitor 获取试运行监控数据
func (s *OptimizationAppService) GetTrialMonitor(ctx context.Context, id string) (*TrialRunResp, error) {
	trial, err := s.trialRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if trial == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToTrialRunResp(trial), nil
}

// GetEvaluation 获取评估报告
func (s *OptimizationAppService) GetEvaluation(ctx context.Context, id string) (*EvaluationResp, error) {
	eval, err := s.evalRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if eval == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToEvaluationResp(eval), nil
}

// GetROIReport 获取 ROI 报告
func (s *OptimizationAppService) GetROIReport(ctx context.Context, id string) (*ROIReportResp, error) {
	roi, err := s.roiRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if roi == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToROIReportResp(roi), nil
}

// ExportROIReport 导出ROI报告PDF（返回报告基础信息，实际PDF由前端/存储服务生成）
func (s *OptimizationAppService) ExportROIReport(ctx context.Context, id string) (*ROIReportResp, error) {
	roi, err := s.roiRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if roi == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToROIReportResp(roi), nil
}

// SubmitROIApprovalResult 回填院级审批结果
func (s *OptimizationAppService) SubmitROIApprovalResult(ctx context.Context, id, approverID string, req ROIApprovalResultReq) error {
	roi, err := s.roiRepo.FindByID(ctx, id)
	if err != nil || roi == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	roi.ApprovalResult = req.Result
	roi.ApprovalResultBy = approverID
	return s.roiRepo.Update(ctx, roi)
}

// ListScans 获取周期扫描列表
func (s *OptimizationAppService) ListScans(ctx context.Context, req ListScansReq) ([]*ScanResp, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	scans, total, err := s.scanRepo.List(ctx, req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]*ScanResp, 0, len(scans))
	for _, sc := range scans {
		resps = append(resps, ToScanResp(sc))
	}
	return resps, total, nil
}

// GetScan 获取扫描详情
func (s *OptimizationAppService) GetScan(ctx context.Context, id string) (*ScanResp, error) {
	sc, err := s.scanRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if sc == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToScanResp(sc), nil
}
