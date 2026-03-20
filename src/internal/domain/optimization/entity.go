// Package optimization 智能效能优化领域 - 实体定义
package optimization

import (
	"fmt"
	"time"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// EfficiencyMetric 效率指标（聚合根）
type EfficiencyMetric struct {
	ID           string
	Name         string
	Code         string
	CalcFormula  string
	Unit         string // % / min / count
	NormalMean   float64
	NormalStdDev float64
	NormalMin    float64
	NormalMax    float64
	IsCustom     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UpdateBaseline 更新指标基线（均値±2σ作为正常区间）
func (m *EfficiencyMetric) UpdateBaseline(mean, stddev float64) {
	m.NormalMean = mean
	m.NormalStdDev = stddev
	m.NormalMin = mean - 2*stddev
	m.NormalMax = mean + 2*stddev
	m.UpdatedAt = time.Now()
}

// MetricSnapshot 指标快照
type MetricSnapshot struct {
	ID         string
	MetricID   string
	Value      float64
	Dimensions map[string]string
	SampledAt  time.Time
}

// BottleneckAlert 瓶颈告警（聚合根）
type BottleneckAlert struct {
	ID                  string
	MetricID            string
	AlertType           AlertType
	Severity            string // warning / critical
	DeviationPct        float64
	ConsecutiveCount    int
	AffectedScope       map[string]string
	RootCauseHypotheses []Hypothesis
	SuggestedCategory   StrategyCategory
	Status              AlertStatus
	DismissReason       string
	CreatedAt           time.Time
}

// Dismiss 标记误报
func (a *BottleneckAlert) Dismiss(reason string) error {
	if a.Status != AlertActive {
		return bizErr.NewWithDetail(bizErr.ErrOptStatusInvalid, "告警状态不允许标记误报")
	}
	a.Status = AlertDismissed
	a.DismissReason = reason
	return nil
}

// OptimizationStrategy 优化策略（聚合根）
type OptimizationStrategy struct {
	ID              string
	Title           string
	Category        StrategyCategory
	Status          StrategyStatus
	AlertID         string
	CurrentValue    string
	TargetValue     string
	ExpectedBenefit string
	RiskNote        string
	CostEstimate    *CostEstimate
	ApprovalFlow    ApprovalFlow
	RejectReason    string
	CooldownUntil   *time.Time
	PromotedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Approve 审批策略
func (s *OptimizationStrategy) Approve(approverID string) error {
	if s.Status != StatusPendingReview {
		return bizErr.NewWithDetail(bizErr.ErrOptStatusInvalid, fmt.Sprintf("策略状态 %s 不允许审批", s.Status))
	}
	now := time.Now()
	for i, node := range s.ApprovalFlow.Approvers {
		if node.ApproverID == approverID && node.Status == "pending" {
			s.ApprovalFlow.Approvers[i].Status = "approved"
			s.ApprovalFlow.Approvers[i].Timestamp = &now
		}
	}
	// 检查是否全部审批
	allApproved := true
	for _, node := range s.ApprovalFlow.Approvers {
		if node.Status != "approved" {
			allApproved = false
			break
		}
	}
	if allApproved {
		s.ApprovalFlow.Status = "approved"
		switch s.Category {
		case CategoryA:
			s.Status = StatusTrialRunning
		case CategoryB:
			s.Status = StatusTrialRunningB
		default: // C类
			s.Status = StatusSubmittedApproval
		}
	}
	s.UpdatedAt = time.Now()
	return nil
}

// Reject 驳回策略
func (s *OptimizationStrategy) Reject(reason string) error {
	if s.Status != StatusPendingReview {
		return bizErr.NewWithDetail(bizErr.ErrOptStatusInvalid, "只有待审核的策略可以驳回")
	}
	if reason == "" {
		return bizErr.New(bizErr.ErrOptRejectReasonReq)
	}
	s.Status = StatusRejected
	s.RejectReason = reason
	cooldown := time.Now().Add(30 * 24 * time.Hour)
	s.CooldownUntil = &cooldown
	s.UpdatedAt = time.Now()
	return nil
}

// Rollback 手动回滚
func (s *OptimizationStrategy) Rollback() error {
	if s.Status != StatusTrialRunning && s.Status != StatusTrialRunningB {
		return bizErr.NewWithDetail(bizErr.ErrOptStatusInvalid, "只有试运行中的策略可以回滚")
	}
	s.Status = StatusRolledBack
	s.UpdatedAt = time.Now()
	return nil
}

// Promote 策略转正/常态化
func (s *OptimizationStrategy) Promote() error {
	if s.Status != StatusPendingEval {
		return bizErr.NewWithDetail(bizErr.ErrOptStatusInvalid, "只有待评估的策略可以转正")
	}
	now := time.Now()
	if s.Category == CategoryA {
		s.Status = StatusPromoted
	} else {
		s.Status = StatusNormalized
	}
	s.PromotedAt = &now
	s.Status = StatusTracking
	s.UpdatedAt = time.Now()
	return nil
}

// TrialRun 试运行记录
type TrialRun struct {
	ID                         string
	StrategyID                 string
	GrayScope                  GrayScope
	TrialDays                  int
	Baseline                   BaselineSnapshot
	ConfigBackup               map[string]interface{}
	EmergencyRollbackThreshold float64
	StartedAt                  time.Time
	EndsAt                     time.Time
	Status                     TrialStatus
	RolledBackAt               *time.Time
	RollbackReason             string
}

// IsRunning 试运行是否进行中
func (t *TrialRun) IsRunning() bool {
	return t.Status == TrialRunning && time.Now().Before(t.EndsAt)
}

// CheckEmergencyRollback 检测是否需要紧急回滚（指标恶化超阈値）
func (t *TrialRun) CheckEmergencyRollback(currentMetrics map[string]float64) bool {
	for code, baseline := range t.Baseline.Metrics {
		current, ok := currentMetrics[code]
		if !ok || baseline == 0 {
			continue
		}
		changePct := (current - baseline) / baseline
		// 负向变化超阈値则触发回滚
		if changePct < -t.EmergencyRollbackThreshold {
			return true
		}
	}
	return false
}

// EvaluationReport 评估报告
type EvaluationReport struct {
	ID               string
	StrategyID       string
	TrialRunID       string
	ReportType       string // trial / long_term_30 / long_term_90 / long_term_180 / post_deploy
	BaselineMetrics  map[string]float64
	TrialMetrics     map[string]float64
	ChangePct        map[string]float64
	IsQualified      bool
	QualifyThreshold float64
	ActualCost       *float64
	ActualROI        *float64
	Recommendation   string // promote / normalize / retry / abandon
	GeneratedAt      time.Time
}

// ROIReport ROI 论证报告（C类专用）
type ROIReport struct {
	ID                    string
	StrategyID            string
	CurrentBottleneck     string
	TriedOptimizations    []string
	InvestmentItems       []InvestItem
	TotalInvestment       float64
	ExpectedAnnualRevenue float64
	PaybackPeriodMonths   int
	RiskFactors           []string
	PDFPath               string
	ApprovalResult        string // approved / rejected / pending
	ApprovalResultBy      string
	CreatedAt             time.Time
}

// PerformanceScan 周期性效能扫描
type PerformanceScan struct {
	ID            string
	ScanWeek      string // 2025-W12
	ScannedAt     time.Time
	Metrics       []ScanMetricResult
	Opportunities []ScanOpportunity
}

// StrategyDecayAlert 策略衰减告警
type StrategyDecayAlert struct {
	ID                  string
	StrategyID          string
	MetricCode          string
	OriginalImprovement float64
	CurrentImprovement  float64
	DecayPct            float64
	DetectedAt          time.Time
}
