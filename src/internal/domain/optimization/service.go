// Package optimization 智能效能优化领域 - 领域服务
package optimization

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// ─── AnomalyDetectionService ───────────────────────────────────────────────

// AnomalyDetectionService 异常检测引擎（μ±2σ + 环比突变 + 趋势恶化）
type AnomalyDetectionService struct {
	metricRepo   EfficiencyMetricRepository
	snapshotRepo MetricSnapshotRepository
	alertRepo    BottleneckAlertRepository
}

func NewAnomalyDetectionService(
	metricRepo EfficiencyMetricRepository,
	snapshotRepo MetricSnapshotRepository,
	alertRepo BottleneckAlertRepository,
) *AnomalyDetectionService {
	return &AnomalyDetectionService{metricRepo: metricRepo, snapshotRepo: snapshotRepo, alertRepo: alertRepo}
}

// RunFullScan 全指标异常检测扫描
func (s *AnomalyDetectionService) RunFullScan(ctx context.Context) ([]BottleneckAlert, error) {
	metrics, err := s.metricRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	var alerts []BottleneckAlert
	for _, m := range metrics {
		recent, err := s.snapshotRepo.FindRecent90Days(ctx, m.ID)
		if err != nil || len(recent) < 3 {
			continue
		}
		// 连续偏离检测
		consecutiveOut := 0
		var lastDeviation float64
		for _, snap := range recent[:minInt(3, len(recent))] {
			if snap.Value < m.NormalMin || snap.Value > m.NormalMax {
				consecutiveOut++
				if m.NormalMean != 0 {
					lastDeviation = math.Abs(snap.Value-m.NormalMean) / math.Abs(m.NormalMean) * 100
				}
			}
		}
		if consecutiveOut >= 3 {
			alert := BottleneckAlert{
				ID: uuid.New().String(), MetricID: m.ID,
				AlertType: AlertConsecutiveDeviation, Severity: "warning",
				DeviationPct: lastDeviation, ConsecutiveCount: consecutiveOut,
				Status: AlertActive, CreatedAt: time.Now(),
			}
			_ = s.alertRepo.Save(ctx, &alert)
			alerts = append(alerts, alert)
		}
		// 环比突变检测（单日 >20%）
		if len(recent) >= 2 {
			prev, curr := recent[1].Value, recent[0].Value
			if prev != 0 {
				changePct := math.Abs(curr-prev) / math.Abs(prev) * 100
				if changePct > 20 {
					alert := BottleneckAlert{
						ID: uuid.New().String(), MetricID: m.ID,
						AlertType: AlertSuddenChange, Severity: "warning",
						DeviationPct: changePct, Status: AlertActive, CreatedAt: time.Now(),
					}
					_ = s.alertRepo.Save(ctx, &alert)
					alerts = append(alerts, alert)
				}
			}
		}
	}
	return alerts, nil
}

// ─── BottleneckAttributionService ──────────────────────────────────────────

// BottleneckAttributionService 瓶颈归因分析服务
type BottleneckAttributionService struct {
	metricRepo EfficiencyMetricRepository
	alertRepo  BottleneckAlertRepository
}

func NewBottleneckAttributionService(
	metricRepo EfficiencyMetricRepository,
	alertRepo BottleneckAlertRepository,
) *BottleneckAttributionService {
	return &BottleneckAttributionService{metricRepo: metricRepo, alertRepo: alertRepo}
}

// AttributionReport 归因报告
type AttributionReport struct {
	AlertID           string
	Hypotheses        []Hypothesis
	SuggestedCategory StrategyCategory
}

// Analyze 对告警进行归因分析
func (s *BottleneckAttributionService) Analyze(ctx context.Context, alertID string) (*AttributionReport, error) {
	alert, err := s.alertRepo.FindByID(ctx, alertID)
	if err != nil {
		return nil, err
	}
	if alert == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	metric, err := s.metricRepo.FindByID(ctx, alert.MetricID)
	if err != nil {
		return nil, err
	}

	report := &AttributionReport{AlertID: alertID}
	switch metric.Code {
	case "device_usage_rate":
		report.Hypotheses = []Hypothesis{
			{RootCause: "设备排班不足", Confidence: 0.7, Evidence: []string{"设备利用率持续超过90%"}},
			{RootCause: "患者分流不均", Confidence: 0.3, Evidence: []string{"部分设备队列超长"}},
		}
		report.SuggestedCategory = CategoryB
	case "avg_wait_min":
		report.Hypotheses = []Hypothesis{
			{RootCause: "号源超配", Confidence: 0.6, Evidence: []string{"预约等待时长超过30分钟"}},
			{RootCause: "签到高峰错位", Confidence: 0.4, Evidence: []string{"早高峰集中签到"}},
		}
		report.SuggestedCategory = CategoryA
	case "conflict_trigger_rate":
		report.Hypotheses = []Hypothesis{
			{RootCause: "冲突规则过严", Confidence: 0.8, Evidence: []string{"冲突触发率高于历史均值"}},
		}
		report.SuggestedCategory = CategoryA
	case "no_show_rate":
		report.Hypotheses = []Hypothesis{
			{RootCause: "通知机制不完善", Confidence: 0.5, Evidence: []string{"爽约率高于行业均值"}},
			{RootCause: "预约时段不合理", Confidence: 0.5, Evidence: []string{"特定时段爽约率更高"}},
		}
		report.SuggestedCategory = CategoryA
	default:
		report.Hypotheses = []Hypothesis{
			{RootCause: "待进一步分析", Confidence: 1.0,
				Evidence: []string{fmt.Sprintf("指标 %s 超出正常区间（偏离%.1f%%）", metric.Code, alert.DeviationPct)}},
		}
		report.SuggestedCategory = CategoryA
	}

	// 更新告警归因
	alert.RootCauseHypotheses = report.Hypotheses
	alert.SuggestedCategory = report.SuggestedCategory
	_ = s.alertRepo.Update(ctx, alert)
	return report, nil
}

// ─── StrategyGenerationService ─────────────────────────────────────────────

// StrategyGenerationService 策略生成服务
type StrategyGenerationService struct {
	alertRepo    BottleneckAlertRepository
	strategyRepo OptimizationStrategyRepository
	attrSvc      *BottleneckAttributionService
}

func NewStrategyGenerationService(
	alertRepo BottleneckAlertRepository,
	strategyRepo OptimizationStrategyRepository,
	attrSvc *BottleneckAttributionService,
) *StrategyGenerationService {
	return &StrategyGenerationService{alertRepo: alertRepo, strategyRepo: strategyRepo, attrSvc: attrSvc}
}

// GenerateFromAlert 根据告警生成优化策略
func (s *StrategyGenerationService) GenerateFromAlert(ctx context.Context, alertID string) (*OptimizationStrategy, error) {
	count, err := s.strategyRepo.CountPendingReview(ctx)
	if err != nil {
		return nil, err
	}
	if count >= 10 {
		return nil, bizErr.New(bizErr.ErrOptStrategyLimit)
	}

	attribution, err := s.attrSvc.Analyze(ctx, alertID)
	if err != nil {
		return nil, err
	}
	alert, _ := s.alertRepo.FindByID(ctx, alertID)
	category := attribution.SuggestedCategory

	var title, currentVal, targetVal, benefit string
	var approvalType string
	var approvers []ApprovalNode

	switch category {
	case CategoryA:
		title = "调整规则配置降低异常指标"
		currentVal = fmt.Sprintf("偏离 %.1f%%", alert.DeviationPct)
		targetVal = "恢复至正常区间"
		benefit = "预期减少 20% 冲突触发，提升号源利用率"
		approvalType = "single"
		approvers = []ApprovalNode{{ApproverRole: "admin", Status: "pending"}}
	case CategoryB:
		title = "弹性调配设备资源"
		currentVal = fmt.Sprintf("偏离 %.1f%%", alert.DeviationPct)
		targetVal = "设备利用率 < 85%"
		benefit = "预期缩短平均等待时长 10 分钟"
		approvalType = "joint"
		approvers = []ApprovalNode{
			{ApproverRole: "admin", Status: "pending"},
			{ApproverRole: "scheduler_admin", Status: "pending"},
		}
	default: // C
		title = "新增物理设备资源"
		currentVal = fmt.Sprintf("偏离 %.1f%%", alert.DeviationPct)
		targetVal = "满足日检查量需求"
		benefit = "根本性解决设备瓶颈"
		approvalType = "external"
		approvers = []ApprovalNode{{ApproverRole: "admin", Status: "pending"}}
	}

	strategy := &OptimizationStrategy{
		ID: uuid.New().String(), Title: title, Category: category,
		Status: StatusPendingReview, AlertID: alertID,
		CurrentValue: currentVal, TargetValue: targetVal, ExpectedBenefit: benefit,
		RiskNote:     "需在低峰期执行，观察 3 天指标变化",
		ApprovalFlow: ApprovalFlow{Type: approvalType, Approvers: approvers, Status: "pending"},
		CreatedAt:    time.Now(), UpdatedAt: time.Now(),
	}
	if err := s.strategyRepo.Save(ctx, strategy); err != nil {
		return nil, err
	}
	return strategy, nil
}

// ─── EvaluationService ─────────────────────────────────────────────────────

// EvaluationService 效果评估服务
type EvaluationService struct {
	strategyRepo OptimizationStrategyRepository
	trialRepo    TrialRunRepository
	evalRepo     EvaluationReportRepository
	snapshotRepo MetricSnapshotRepository
	metricRepo   EfficiencyMetricRepository
}

func NewEvaluationService(
	strategyRepo OptimizationStrategyRepository,
	trialRepo TrialRunRepository,
	evalRepo EvaluationReportRepository,
	snapshotRepo MetricSnapshotRepository,
	metricRepo EfficiencyMetricRepository,
) *EvaluationService {
	return &EvaluationService{
		strategyRepo: strategyRepo, trialRepo: trialRepo, evalRepo: evalRepo,
		snapshotRepo: snapshotRepo, metricRepo: metricRepo,
	}
}

// GenerateReport 试运行到期后生成评估报告
func (s *EvaluationService) GenerateReport(ctx context.Context, strategyID string) (*EvaluationReport, error) {
	strategy, err := s.strategyRepo.FindByID(ctx, strategyID)
	if err != nil || strategy == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	trial, err := s.trialRepo.FindByStrategyID(ctx, strategyID)
	if err != nil || trial == nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrOptEvalFailed, "试运行记录不存在")
	}

	trialMetrics := make(map[string]float64)
	changePct := make(map[string]float64)
	qualified, total := 0, 0

	for code, baseVal := range trial.Baseline.Metrics {
		metric, err := s.metricRepo.FindByCode(ctx, code)
		if err != nil || metric == nil {
			continue
		}
		snaps, err := s.snapshotRepo.FindByMetricID(ctx, metric.ID, 10)
		if err != nil || len(snaps) == 0 {
			continue
		}
		var sum float64
		for _, snap := range snaps {
			sum += snap.Value
		}
		trialVal := sum / float64(len(snaps))
		trialMetrics[code] = trialVal
		if baseVal != 0 {
			pct := (trialVal - baseVal) / math.Abs(baseVal) * 100
			changePct[code] = pct
			total++
			if pct >= 5.0 {
				qualified++
			}
		}
	}

	isQualified := total > 0 && float64(qualified)/float64(total) >= 0.5
	recommendation := "promote"
	if !isQualified {
		recommendation = "retry"
	}
	if strategy.Category == CategoryB {
		recommendation = "normalize"
	}

	report := &EvaluationReport{
		ID: uuid.New().String(), StrategyID: strategyID, TrialRunID: trial.ID,
		ReportType: "trial", BaselineMetrics: trial.Baseline.Metrics,
		TrialMetrics: trialMetrics, ChangePct: changePct,
		IsQualified: isQualified, QualifyThreshold: 5.0,
		Recommendation: recommendation, GeneratedAt: time.Now(),
	}
	if err := s.evalRepo.Save(ctx, report); err != nil {
		return nil, err
	}
	strategy.Status = StatusPendingEval
	_ = s.strategyRepo.Update(ctx, strategy)
	return report, nil
}

// PromoteStrategy 策略转正
func (s *EvaluationService) PromoteStrategy(ctx context.Context, strategyID string) error {
	strategy, err := s.strategyRepo.FindByID(ctx, strategyID)
	if err != nil || strategy == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := strategy.Promote(); err != nil {
		return err
	}
	return s.strategyRepo.Update(ctx, strategy)
}

// ─── DecayTrackingService ───────────────────────────────────────────────────

// DecayTrackingService 策略衰减追踪服务
type DecayTrackingService struct {
	strategyRepo OptimizationStrategyRepository
	evalRepo     EvaluationReportRepository
	decayRepo    StrategyDecayAlertRepository
	snapshotRepo MetricSnapshotRepository
	metricRepo   EfficiencyMetricRepository
}

func NewDecayTrackingService(
	strategyRepo OptimizationStrategyRepository,
	evalRepo EvaluationReportRepository,
	decayRepo StrategyDecayAlertRepository,
	snapshotRepo MetricSnapshotRepository,
	metricRepo EfficiencyMetricRepository,
) *DecayTrackingService {
	return &DecayTrackingService{
		strategyRepo: strategyRepo, evalRepo: evalRepo, decayRepo: decayRepo,
		snapshotRepo: snapshotRepo, metricRepo: metricRepo,
	}
}

// CheckDecay 检测已转正策略的效果是否衰减
func (s *DecayTrackingService) CheckDecay(ctx context.Context) ([]StrategyDecayAlert, error) {
	strategies, err := s.strategyRepo.ListPromoted(ctx)
	if err != nil {
		return nil, err
	}
	var decayAlerts []StrategyDecayAlert
	for _, strategy := range strategies {
		evalReport, err := s.evalRepo.FindByStrategyID(ctx, strategy.ID)
		if err != nil || evalReport == nil {
			continue
		}
		for code, originalChange := range evalReport.ChangePct {
			if originalChange <= 0 {
				continue
			}
			metric, err := s.metricRepo.FindByCode(ctx, code)
			if err != nil || metric == nil {
				continue
			}
			snaps, err := s.snapshotRepo.FindByMetricID(ctx, metric.ID, 5)
			if err != nil || len(snaps) == 0 {
				continue
			}
			var sum float64
			for _, snap := range snaps {
				sum += snap.Value
			}
			currentVal := sum / float64(len(snaps))
			baseVal := evalReport.BaselineMetrics[code]
			if baseVal == 0 {
				continue
			}
			currentChange := (currentVal - baseVal) / math.Abs(baseVal) * 100
			decayPct := originalChange - currentChange
			if decayPct > originalChange*0.5 { // 改善效果损失超 50%
				da := StrategyDecayAlert{
					ID: uuid.New().String(), StrategyID: strategy.ID,
					MetricCode: code, OriginalImprovement: originalChange,
					CurrentImprovement: currentChange, DecayPct: decayPct,
					DetectedAt: time.Now(),
				}
				_ = s.decayRepo.Save(ctx, &da)
				strategy.Status = StatusDecayed
				_ = s.strategyRepo.Update(ctx, strategy)
				decayAlerts = append(decayAlerts, da)
			}
		}
	}
	return decayAlerts, nil
}

// ─── 辅助函数 ───────────────────────────────────────────────────────────────

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
