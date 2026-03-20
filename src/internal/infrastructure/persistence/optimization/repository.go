// Package optimization 基础设施层 - optimization 仓储实现（GORM）
package optimization

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	domain "github.com/euler/mtap/internal/domain/optimization"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// Repositories optimization 模块仓储集合
type Repositories struct {
	db *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{db: db}
}

func (r *Repositories) MetricRepo() domain.EfficiencyMetricRepository { return &metricRepo{r.db} }
func (r *Repositories) SnapshotRepo() domain.MetricSnapshotRepository { return &snapshotRepo{r.db} }
func (r *Repositories) AlertRepo() domain.BottleneckAlertRepository   { return &alertRepo{r.db} }
func (r *Repositories) StrategyRepo() domain.OptimizationStrategyRepository {
	return &strategyRepo{r.db}
}
func (r *Repositories) TrialRepo() domain.TrialRunRepository           { return &trialRunRepo{r.db} }
func (r *Repositories) EvalRepo() domain.EvaluationReportRepository    { return &evalReportRepo{r.db} }
func (r *Repositories) ROIRepo() domain.ROIReportRepository            { return &roiReportRepo{r.db} }
func (r *Repositories) ScanRepo() domain.PerformanceScanRepository     { return &scanRepo{r.db} }
func (r *Repositories) DecayRepo() domain.StrategyDecayAlertRepository { return &decayRepo{r.db} }

// ─── metricRepo ───────────────────────────────────────────────────────────────

type metricRepo struct{ db *gorm.DB }

func (r *metricRepo) Save(ctx context.Context, m *domain.EfficiencyMetric) error {
	return r.db.WithContext(ctx).Create(&po.EfficiencyMetricPO{
		ID: m.ID, Name: m.Name, Code: m.Code, CalcFormula: m.CalcFormula,
		Unit: m.Unit, NormalMean: m.NormalMean, NormalStdDev: m.NormalStdDev,
		NormalMin: m.NormalMin, NormalMax: m.NormalMax, IsCustom: m.IsCustom,
	}).Error
}

func (r *metricRepo) FindByID(ctx context.Context, id string) (*domain.EfficiencyMetric, error) {
	var p po.EfficiencyMetricPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return metricFromPO(&p), nil
}

func (r *metricRepo) FindByCode(ctx context.Context, code string) (*domain.EfficiencyMetric, error) {
	var p po.EfficiencyMetricPO
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return metricFromPO(&p), nil
}

func (r *metricRepo) List(ctx context.Context) ([]*domain.EfficiencyMetric, error) {
	var pos []po.EfficiencyMetricPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.EfficiencyMetric, 0, len(pos))
	for i := range pos {
		result = append(result, metricFromPO(&pos[i]))
	}
	return result, nil
}

func (r *metricRepo) Update(ctx context.Context, m *domain.EfficiencyMetric) error {
	return r.db.WithContext(ctx).Model(&po.EfficiencyMetricPO{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"normal_mean": m.NormalMean, "normal_stddev": m.NormalStdDev,
		"normal_min": m.NormalMin, "normal_max": m.NormalMax,
	}).Error
}

func metricFromPO(p *po.EfficiencyMetricPO) *domain.EfficiencyMetric {
	return &domain.EfficiencyMetric{
		ID: p.ID, Name: p.Name, Code: p.Code, CalcFormula: p.CalcFormula,
		Unit: p.Unit, NormalMean: p.NormalMean, NormalStdDev: p.NormalStdDev,
		NormalMin: p.NormalMin, NormalMax: p.NormalMax, IsCustom: p.IsCustom,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

// ─── snapshotRepo ───────────────────────────────────────────────────────────

type snapshotRepo struct{ db *gorm.DB }

func (r *snapshotRepo) Save(ctx context.Context, s *domain.MetricSnapshot) error {
	dims, _ := json.Marshal(s.Dimensions)
	return r.db.WithContext(ctx).Create(&po.MetricSnapshotPO{
		ID: s.ID, MetricID: s.MetricID, Value: s.Value,
		Dimensions: string(dims), SampledAt: s.SampledAt,
	}).Error
}

func (r *snapshotRepo) FindByMetricID(ctx context.Context, metricID string, limit int) ([]*domain.MetricSnapshot, error) {
	var pos []po.MetricSnapshotPO
	if err := r.db.WithContext(ctx).Where("metric_id = ?", metricID).
		Order("sampled_at DESC").Limit(limit).Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return snapshotsFromPOs(pos), nil
}

func (r *snapshotRepo) FindRecent90Days(ctx context.Context, metricID string) ([]*domain.MetricSnapshot, error) {
	var pos []po.MetricSnapshotPO
	cutoff := time.Now().AddDate(0, 0, -90)
	if err := r.db.WithContext(ctx).Where("metric_id = ? AND sampled_at >= ?", metricID, cutoff).
		Order("sampled_at DESC").Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return snapshotsFromPOs(pos), nil
}

func snapshotsFromPOs(pos []po.MetricSnapshotPO) []*domain.MetricSnapshot {
	result := make([]*domain.MetricSnapshot, 0, len(pos))
	for _, p := range pos {
		var dims map[string]string
		_ = json.Unmarshal([]byte(p.Dimensions), &dims)
		result = append(result, &domain.MetricSnapshot{
			ID: p.ID, MetricID: p.MetricID, Value: p.Value,
			Dimensions: dims, SampledAt: p.SampledAt,
		})
	}
	return result
}

// ─── alertRepo ─────────────────────────────────────────────────────────────────

type alertRepo struct{ db *gorm.DB }

func (r *alertRepo) Save(ctx context.Context, a *domain.BottleneckAlert) error {
	hypo, _ := json.Marshal(a.RootCauseHypotheses)
	scope, _ := json.Marshal(a.AffectedScope)
	return r.db.WithContext(ctx).Create(&po.BottleneckAlertPO{
		ID: a.ID, MetricID: a.MetricID, AlertType: string(a.AlertType),
		Severity: a.Severity, DeviationPct: a.DeviationPct,
		ConsecutiveCount: a.ConsecutiveCount, AffectedScope: string(scope),
		RootCauseHypotheses: string(hypo), SuggestedCategory: string(a.SuggestedCategory),
		Status: string(a.Status), DismissReason: a.DismissReason, CreatedAt: a.CreatedAt,
	}).Error
}

func (r *alertRepo) FindByID(ctx context.Context, id string) (*domain.BottleneckAlert, error) {
	var p po.BottleneckAlertPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return alertFromPO(&p), nil
}

func (r *alertRepo) List(ctx context.Context, status string, page, size int) ([]*domain.BottleneckAlert, int64, error) {
	var total int64
	var pos []po.BottleneckAlertPO
	db := r.db.WithContext(ctx).Model(&po.BottleneckAlertPO{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if err := db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.BottleneckAlert, 0, len(pos))
	for i := range pos {
		result = append(result, alertFromPO(&pos[i]))
	}
	return result, total, nil
}

func (r *alertRepo) Update(ctx context.Context, a *domain.BottleneckAlert) error {
	hypo, _ := json.Marshal(a.RootCauseHypotheses)
	return r.db.WithContext(ctx).Model(&po.BottleneckAlertPO{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"status": string(a.Status), "dismiss_reason": a.DismissReason,
		"root_cause_hypotheses": string(hypo), "suggested_category": string(a.SuggestedCategory),
	}).Error
}

func alertFromPO(p *po.BottleneckAlertPO) *domain.BottleneckAlert {
	var hypo []domain.Hypothesis
	_ = json.Unmarshal([]byte(p.RootCauseHypotheses), &hypo)
	var scope map[string]string
	_ = json.Unmarshal([]byte(p.AffectedScope), &scope)
	return &domain.BottleneckAlert{
		ID: p.ID, MetricID: p.MetricID, AlertType: domain.AlertType(p.AlertType),
		Severity: p.Severity, DeviationPct: p.DeviationPct,
		ConsecutiveCount: p.ConsecutiveCount, AffectedScope: scope,
		RootCauseHypotheses: hypo, SuggestedCategory: domain.StrategyCategory(p.SuggestedCategory),
		Status: domain.AlertStatus(p.Status), DismissReason: p.DismissReason, CreatedAt: p.CreatedAt,
	}
}

// ─── strategyRepo ────────────────────────────────────────────────────────────

type strategyRepo struct{ db *gorm.DB }

func (r *strategyRepo) Save(ctx context.Context, s *domain.OptimizationStrategy) error {
	ap, _ := json.Marshal(s.ApprovalFlow)
	ce, _ := json.Marshal(s.CostEstimate)
	return r.db.WithContext(ctx).Create(&po.OptimizationStrategyPO{
		ID: s.ID, Title: s.Title, Category: string(s.Category), Status: string(s.Status),
		AlertID: s.AlertID, CurrentValue: s.CurrentValue, TargetValue: s.TargetValue,
		ExpectedBenefit: s.ExpectedBenefit, RiskNote: s.RiskNote,
		CostEstimate: string(ce), ApprovalFlow: string(ap),
		RejectReason: s.RejectReason, CooldownUntil: s.CooldownUntil, PromotedAt: s.PromotedAt,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}).Error
}

func (r *strategyRepo) FindByID(ctx context.Context, id string) (*domain.OptimizationStrategy, error) {
	var p po.OptimizationStrategyPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return strategyFromPO(&p), nil
}

func (r *strategyRepo) List(ctx context.Context, category, status string, page, size int) ([]*domain.OptimizationStrategy, int64, error) {
	var total int64
	var pos []po.OptimizationStrategyPO
	db := r.db.WithContext(ctx).Model(&po.OptimizationStrategyPO{})
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if err := db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.OptimizationStrategy, 0, len(pos))
	for i := range pos {
		result = append(result, strategyFromPO(&pos[i]))
	}
	return result, total, nil
}

func (r *strategyRepo) Update(ctx context.Context, s *domain.OptimizationStrategy) error {
	ap, _ := json.Marshal(s.ApprovalFlow)
	return r.db.WithContext(ctx).Model(&po.OptimizationStrategyPO{}).Where("id = ?", s.ID).
		Updates(map[string]interface{}{
			"status": string(s.Status), "reject_reason": s.RejectReason,
			"approval_flow": string(ap), "cooldown_until": s.CooldownUntil,
			"promoted_at": s.PromotedAt, "updated_at": time.Now(),
		}).Error
}

func (r *strategyRepo) ListActiveTrials(ctx context.Context) ([]*domain.OptimizationStrategy, error) {
	var pos []po.OptimizationStrategyPO
	if err := r.db.WithContext(ctx).
		Where("status IN ?", []string{string(domain.StatusTrialRunning), string(domain.StatusTrialRunningB)}).
		Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.OptimizationStrategy, 0, len(pos))
	for i := range pos {
		result = append(result, strategyFromPO(&pos[i]))
	}
	return result, nil
}

func (r *strategyRepo) ListPromoted(ctx context.Context) ([]*domain.OptimizationStrategy, error) {
	var pos []po.OptimizationStrategyPO
	if err := r.db.WithContext(ctx).
		Where("status IN ?", []string{
			string(domain.StatusPromoted), string(domain.StatusNormalized), string(domain.StatusTracking),
		}).Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.OptimizationStrategy, 0, len(pos))
	for i := range pos {
		result = append(result, strategyFromPO(&pos[i]))
	}
	return result, nil
}

func (r *strategyRepo) CountPendingReview(ctx context.Context) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&po.OptimizationStrategyPO{}).
		Where("status = ?", string(domain.StatusPendingReview)).Count(&count).Error
}

func strategyFromPO(p *po.OptimizationStrategyPO) *domain.OptimizationStrategy {
	var af domain.ApprovalFlow
	_ = json.Unmarshal([]byte(p.ApprovalFlow), &af)
	var ce *domain.CostEstimate
	_ = json.Unmarshal([]byte(p.CostEstimate), &ce)
	return &domain.OptimizationStrategy{
		ID: p.ID, Title: p.Title, Category: domain.StrategyCategory(p.Category),
		Status: domain.StrategyStatus(p.Status), AlertID: p.AlertID,
		CurrentValue: p.CurrentValue, TargetValue: p.TargetValue,
		ExpectedBenefit: p.ExpectedBenefit, RiskNote: p.RiskNote,
		CostEstimate: ce, ApprovalFlow: af,
		RejectReason: p.RejectReason, CooldownUntil: p.CooldownUntil, PromotedAt: p.PromotedAt,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

// ─── trialRunRepo ─────────────────────────────────────────────────────────────

type trialRunRepo struct{ db *gorm.DB }

func (r *trialRunRepo) Save(ctx context.Context, t *domain.TrialRun) error {
	gs, _ := json.Marshal(t.GrayScope)
	bs, _ := json.Marshal(t.Baseline)
	cb, _ := json.Marshal(t.ConfigBackup)
	return r.db.WithContext(ctx).Create(&po.TrialRunPO{
		ID: t.ID, StrategyID: t.StrategyID, GrayScope: string(gs),
		TrialDays: t.TrialDays, BaselineSnapshot: string(bs), ConfigBackup: string(cb),
		EmergencyRollbackThreshold: t.EmergencyRollbackThreshold,
		StartedAt:                  t.StartedAt, EndsAt: t.EndsAt, Status: string(t.Status),
		RolledBackAt: t.RolledBackAt, RollbackReason: t.RollbackReason,
	}).Error
}

func (r *trialRunRepo) FindByID(ctx context.Context, id string) (*domain.TrialRun, error) {
	var p po.TrialRunPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return trialFromPO(&p), nil
}

func (r *trialRunRepo) FindByStrategyID(ctx context.Context, strategyID string) (*domain.TrialRun, error) {
	var p po.TrialRunPO
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).
		Order("started_at DESC").First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return trialFromPO(&p), nil
}

func (r *trialRunRepo) Update(ctx context.Context, t *domain.TrialRun) error {
	return r.db.WithContext(ctx).Model(&po.TrialRunPO{}).Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"status": string(t.Status), "rolled_back_at": t.RolledBackAt,
			"rollback_reason": t.RollbackReason,
		}).Error
}

func trialFromPO(p *po.TrialRunPO) *domain.TrialRun {
	var gs domain.GrayScope
	_ = json.Unmarshal([]byte(p.GrayScope), &gs)
	var bs domain.BaselineSnapshot
	_ = json.Unmarshal([]byte(p.BaselineSnapshot), &bs)
	return &domain.TrialRun{
		ID: p.ID, StrategyID: p.StrategyID, GrayScope: gs,
		TrialDays: p.TrialDays, Baseline: bs,
		EmergencyRollbackThreshold: p.EmergencyRollbackThreshold,
		StartedAt:                  p.StartedAt, EndsAt: p.EndsAt, Status: domain.TrialStatus(p.Status),
		RolledBackAt: p.RolledBackAt, RollbackReason: p.RollbackReason,
	}
}

// ─── evalReportRepo ───────────────────────────────────────────────────────────

type evalReportRepo struct{ db *gorm.DB }

func (r *evalReportRepo) Save(ctx context.Context, e *domain.EvaluationReport) error {
	bm, _ := json.Marshal(e.BaselineMetrics)
	tm, _ := json.Marshal(e.TrialMetrics)
	cp, _ := json.Marshal(e.ChangePct)
	return r.db.WithContext(ctx).Create(&po.EvaluationReportPO{
		ID: e.ID, StrategyID: e.StrategyID, TrialRunID: e.TrialRunID,
		ReportType: e.ReportType, BaselineMetrics: string(bm),
		TrialMetrics: string(tm), ChangePct: string(cp),
		IsQualified: e.IsQualified, QualifyThreshold: e.QualifyThreshold,
		ActualCost: e.ActualCost, ActualROI: e.ActualROI,
		Recommendation: e.Recommendation, GeneratedAt: e.GeneratedAt,
	}).Error
}

func (r *evalReportRepo) FindByID(ctx context.Context, id string) (*domain.EvaluationReport, error) {
	var p po.EvaluationReportPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return evalFromPO(&p), nil
}

func (r *evalReportRepo) FindByStrategyID(ctx context.Context, strategyID string) (*domain.EvaluationReport, error) {
	var p po.EvaluationReportPO
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).
		Order("generated_at DESC").First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return evalFromPO(&p), nil
}

func evalFromPO(p *po.EvaluationReportPO) *domain.EvaluationReport {
	var bm, tm, cp map[string]float64
	_ = json.Unmarshal([]byte(p.BaselineMetrics), &bm)
	_ = json.Unmarshal([]byte(p.TrialMetrics), &tm)
	_ = json.Unmarshal([]byte(p.ChangePct), &cp)
	return &domain.EvaluationReport{
		ID: p.ID, StrategyID: p.StrategyID, TrialRunID: p.TrialRunID,
		ReportType: p.ReportType, BaselineMetrics: bm, TrialMetrics: tm, ChangePct: cp,
		IsQualified: p.IsQualified, QualifyThreshold: p.QualifyThreshold,
		ActualCost: p.ActualCost, ActualROI: p.ActualROI,
		Recommendation: p.Recommendation, GeneratedAt: p.GeneratedAt,
	}
}

// ─── roiReportRepo ────────────────────────────────────────────────────────────

type roiReportRepo struct{ db *gorm.DB }

func (r *roiReportRepo) Save(ctx context.Context, roi *domain.ROIReport) error {
	to, _ := json.Marshal(roi.TriedOptimizations)
	ii, _ := json.Marshal(roi.InvestmentItems)
	rf, _ := json.Marshal(roi.RiskFactors)
	return r.db.WithContext(ctx).Create(&po.ROIReportPO{
		ID: roi.ID, StrategyID: roi.StrategyID,
		CurrentBottleneck: roi.CurrentBottleneck, TriedOptimizations: string(to),
		InvestmentItems: string(ii), TotalInvestment: roi.TotalInvestment,
		ExpectedAnnualRevenue: roi.ExpectedAnnualRevenue, PaybackPeriodMonths: roi.PaybackPeriodMonths,
		RiskFactors: string(rf), PDFPath: roi.PDFPath,
		ApprovalResult: roi.ApprovalResult, ApprovalResultBy: roi.ApprovalResultBy,
	}).Error
}

func (r *roiReportRepo) FindByID(ctx context.Context, id string) (*domain.ROIReport, error) {
	var p po.ROIReportPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return roiFromPO(&p), nil
}

func (r *roiReportRepo) FindByStrategyID(ctx context.Context, strategyID string) (*domain.ROIReport, error) {
	var p po.ROIReportPO
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return roiFromPO(&p), nil
}

func (r *roiReportRepo) Update(ctx context.Context, roi *domain.ROIReport) error {
	return r.db.WithContext(ctx).Model(&po.ROIReportPO{}).Where("id = ?", roi.ID).
		Updates(map[string]interface{}{
			"approval_result": roi.ApprovalResult, "approval_result_by": roi.ApprovalResultBy,
			"pdf_path": roi.PDFPath,
		}).Error
}

func roiFromPO(p *po.ROIReportPO) *domain.ROIReport {
	var to []string
	_ = json.Unmarshal([]byte(p.TriedOptimizations), &to)
	var ii []domain.InvestItem
	_ = json.Unmarshal([]byte(p.InvestmentItems), &ii)
	var rf []string
	_ = json.Unmarshal([]byte(p.RiskFactors), &rf)
	return &domain.ROIReport{
		ID: p.ID, StrategyID: p.StrategyID,
		CurrentBottleneck: p.CurrentBottleneck, TriedOptimizations: to,
		InvestmentItems: ii, TotalInvestment: p.TotalInvestment,
		ExpectedAnnualRevenue: p.ExpectedAnnualRevenue, PaybackPeriodMonths: p.PaybackPeriodMonths,
		RiskFactors: rf, PDFPath: p.PDFPath,
		ApprovalResult: p.ApprovalResult, ApprovalResultBy: p.ApprovalResultBy, CreatedAt: p.CreatedAt,
	}
}

// ─── scanRepo ───────────────────────────────────────────────────────────────────

type scanRepo struct{ db *gorm.DB }

func (r *scanRepo) Save(ctx context.Context, s *domain.PerformanceScan) error {
	mets, _ := json.Marshal(s.Metrics)
	opps, _ := json.Marshal(s.Opportunities)
	return r.db.WithContext(ctx).Create(&po.PerformanceScanPO{
		ID: s.ID, ScanWeek: s.ScanWeek, Metrics: string(mets),
		Opportunities: string(opps), ScannedAt: s.ScannedAt,
	}).Error
}

func (r *scanRepo) FindByID(ctx context.Context, id string) (*domain.PerformanceScan, error) {
	var p po.PerformanceScanPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return scanFromPO(&p), nil
}

func (r *scanRepo) List(ctx context.Context, page, size int) ([]*domain.PerformanceScan, int64, error) {
	var total int64
	var pos []po.PerformanceScanPO
	db := r.db.WithContext(ctx).Model(&po.PerformanceScanPO{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if err := db.Order("scanned_at DESC").Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.PerformanceScan, 0, len(pos))
	for i := range pos {
		result = append(result, scanFromPO(&pos[i]))
	}
	return result, total, nil
}

func scanFromPO(p *po.PerformanceScanPO) *domain.PerformanceScan {
	var mets []domain.ScanMetricResult
	_ = json.Unmarshal([]byte(p.Metrics), &mets)
	var opps []domain.ScanOpportunity
	_ = json.Unmarshal([]byte(p.Opportunities), &opps)
	return &domain.PerformanceScan{
		ID: p.ID, ScanWeek: p.ScanWeek, ScannedAt: p.ScannedAt,
		Metrics: mets, Opportunities: opps,
	}
}

// ─── decayRepo ─────────────────────────────────────────────────────────────────

type decayRepo struct{ db *gorm.DB }

func (r *decayRepo) Save(ctx context.Context, a *domain.StrategyDecayAlert) error {
	return r.db.WithContext(ctx).Create(&po.StrategyDecayAlertPO{
		ID: a.ID, StrategyID: a.StrategyID, MetricCode: a.MetricCode,
		OriginalImprovement: a.OriginalImprovement, CurrentImprovement: a.CurrentImprovement,
		DecayPct: a.DecayPct, DetectedAt: a.DetectedAt,
	}).Error
}

func (r *decayRepo) FindByStrategyID(ctx context.Context, strategyID string) ([]*domain.StrategyDecayAlert, error) {
	var pos []po.StrategyDecayAlertPO
	if err := r.db.WithContext(ctx).Where("strategy_id = ?", strategyID).
		Order("detected_at DESC").Find(&pos).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]*domain.StrategyDecayAlert, 0, len(pos))
	for _, p := range pos {
		result = append(result, &domain.StrategyDecayAlert{
			ID: p.ID, StrategyID: p.StrategyID, MetricCode: p.MetricCode,
			OriginalImprovement: p.OriginalImprovement, CurrentImprovement: p.CurrentImprovement,
			DecayPct: p.DecayPct, DetectedAt: p.DetectedAt,
		})
	}
	return result, nil
}
