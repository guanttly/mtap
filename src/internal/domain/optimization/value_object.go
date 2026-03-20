// Package optimization 智能效能优化领域 - 值对象
package optimization

import "time"

// StrategyCategory 策略分类
type StrategyCategory string

const (
	CategoryA StrategyCategory = "A" // A类·软策略（纯配置调整）
	CategoryB StrategyCategory = "B" // B类·弹性资源策略（现有资源重新调配）
	CategoryC StrategyCategory = "C" // C类·硬资源策略（新增物理资源）
)

// StrategyStatus 策略状态
type StrategyStatus string

const (
	StatusPendingReview     StrategyStatus = "pending_review"     // 待审核
	StatusRejected          StrategyStatus = "rejected"           // 已驳回
	StatusTrialRunning      StrategyStatus = "trial_running"      // 试运行中(A类)
	StatusTrialRunningB     StrategyStatus = "trial_running_b"    // 试行中(B类)
	StatusSubmittedApproval StrategyStatus = "submitted_approval" // 已提交院级审批(C类)
	StatusPendingEval       StrategyStatus = "pending_eval"       // 待评估
	StatusPromoted          StrategyStatus = "promoted"           // 已转正(A类)
	StatusNormalized        StrategyStatus = "normalized"         // 已常态化(B类)
	StatusRolledBack        StrategyStatus = "rolled_back"        // 已回滚
	StatusArchived          StrategyStatus = "archived"           // 已归档
	StatusPostValidation    StrategyStatus = "post_validation"    // 待投产验证(C类)
	StatusTracking          StrategyStatus = "tracking"           // 长期追踪
	StatusDecayed           StrategyStatus = "decayed"            // 策略衰减
)

// AlertType 告警类型
type AlertType string

const (
	AlertConsecutiveDeviation AlertType = "consecutive_deviation" // 连续偏离
	AlertSuddenChange         AlertType = "sudden_change"         // 环比突变
	AlertTrendDegradation     AlertType = "trend_degradation"     // 趋势恶化
)

// AlertStatus 告警状态
type AlertStatus string

const (
	AlertActive    AlertStatus = "active"
	AlertDismissed AlertStatus = "dismissed"
	AlertResolved  AlertStatus = "resolved"
)

// TrialStatus 试运行状态
type TrialStatus string

const (
	TrialRunning    TrialStatus = "running"
	TrialCompleted  TrialStatus = "completed"
	TrialRolledBack TrialStatus = "rolled_back"
)

// GrayScope 灰度范围
type GrayScope struct {
	DepartmentIDs []string
	DeviceIDs     []string
	TimePeriods   []string // 如 "08:00-12:00"
}

// ApprovalNode 审批节点
type ApprovalNode struct {
	ApproverID   string
	ApproverRole string
	Status       string // pending / approved / rejected
	Timestamp    *time.Time
	Comment      string
}

// ApprovalFlow 审批流
type ApprovalFlow struct {
	Type      string // single(A类) / joint(B类) / external(C类)
	Approvers []ApprovalNode
	Status    string
}

// CostEstimate 成本评估（B类专用）
type CostEstimate struct {
	LaborCost    float64
	MaterialCost float64
	EnergyCost   float64
	TotalCost    float64
	ROI          float64
}

// BaselineSnapshot 基线快照
type BaselineSnapshot struct {
	Metrics    map[string]float64
	SnapshotAt time.Time
}

// Hypothesis 归因假设
type Hypothesis struct {
	RootCause  string
	Confidence float64
	Evidence   []string
}

// ScanOpportunity 效能优化机会
type ScanOpportunity struct {
	Code              string
	Description       string
	MetricCode        string
	CurrentValue      float64
	ExpectedValue     float64
	SuggestedCategory StrategyCategory
}

// ScanMetricResult 扫描指标结果
type ScanMetricResult struct {
	MetricCode   string
	CurrentValue float64
	NormalMin    float64
	NormalMax    float64
	Status       string // normal / warning / critical
}

// InvestItem ROI 投资明细
type InvestItem struct {
	Item   string
	Amount float64
	Remark string
}
