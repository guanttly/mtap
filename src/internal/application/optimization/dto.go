// Package optimization 应用层 - 效能优化数据传输对象
package optimization

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────

// ListMetricsReq 指标列表请求
type ListMetricsReq struct{}

// GetMetricTrendReq 指标趋势请求
type GetMetricTrendReq struct {
	Days int `form:"days"` // 7 / 30 / 90
}

// ListAlertsReq 告警列表请求
type ListAlertsReq struct {
	Status string `form:"status"`
	Page   int    `form:"page"`
	Size   int    `form:"size"`
}

// DismissAlertReq 标记告警请求
type DismissAlertReq struct {
	Reason string `json:"reason" binding:"required"`
}

// ListStrategiesReq 策略列表请求
type ListStrategiesReq struct {
	Category string `form:"category"`
	Status   string `form:"status"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
}

// ApproveStrategyReq 审批策略请求
type ApproveStrategyReq struct {
	TrialDays int          `json:"trial_days" binding:"required,min=1,max=90"`
	GrayScope GrayScopeDTO `json:"gray_scope"`
	Comment   string       `json:"comment"`
}

// GrayScopeDTO 灰度范围 DTO
type GrayScopeDTO struct {
	DepartmentIDs []string `json:"department_ids"`
	DeviceIDs     []string `json:"device_ids"`
	TimePeriods   []string `json:"time_periods"`
}

// RejectStrategyReq 驳回策略请求
type RejectStrategyReq struct {
	Reason string `json:"reason" binding:"required"`
}

// ROIApprovalResultReq 回填院级审批结果
type ROIApprovalResultReq struct {
	Result string `json:"result" binding:"required,oneof=approved rejected pending"`
}

// ListScansReq 扫描列表请求
type ListScansReq struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────

// MetricResp 指标响应
type MetricResp struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Unit       string  `json:"unit"`
	NormalMin  float64 `json:"normal_min"`
	NormalMax  float64 `json:"normal_max"`
	NormalMean float64 `json:"normal_mean"`
}

// MetricSnapshotResp 指标快照响应
type MetricSnapshotResp struct {
	Value     float64   `json:"value"`
	SampledAt time.Time `json:"sampled_at"`
}

// AlertResp 告警响应
type AlertResp struct {
	ID                string    `json:"id"`
	MetricID          string    `json:"metric_id"`
	AlertType         string    `json:"alert_type"`
	Severity          string    `json:"severity"`
	DeviationPct      float64   `json:"deviation_pct"`
	ConsecutiveCount  int       `json:"consecutive_count"`
	SuggestedCategory string    `json:"suggested_category"`
	Status            string    `json:"status"`
	DismissReason     string    `json:"dismiss_reason,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

// ApprovalNodeResp 审批节点响应
type ApprovalNodeResp struct {
	ApproverRole string     `json:"approver_role"`
	Status       string     `json:"status"`
	Timestamp    *time.Time `json:"timestamp,omitempty"`
	Comment      string     `json:"comment,omitempty"`
}

// StrategyResp 策略响应
type StrategyResp struct {
	ID              string             `json:"id"`
	Title           string             `json:"title"`
	Category        string             `json:"category"`
	Status          string             `json:"status"`
	AlertID         string             `json:"alert_id,omitempty"`
	CurrentValue    string             `json:"current_value"`
	TargetValue     string             `json:"target_value"`
	ExpectedBenefit string             `json:"expected_benefit"`
	RiskNote        string             `json:"risk_note"`
	ApprovalType    string             `json:"approval_type"`
	Approvers       []ApprovalNodeResp `json:"approvers"`
	RejectReason    string             `json:"reject_reason,omitempty"`
	PromotedAt      *time.Time         `json:"promoted_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
}

// TrialRunResp 试运行响应
type TrialRunResp struct {
	ID         string    `json:"id"`
	StrategyID string    `json:"strategy_id"`
	TrialDays  int       `json:"trial_days"`
	StartedAt  time.Time `json:"started_at"`
	EndsAt     time.Time `json:"ends_at"`
	Status     string    `json:"status"`
}

// EvaluationResp 评估报告响应
type EvaluationResp struct {
	ID               string             `json:"id"`
	StrategyID       string             `json:"strategy_id"`
	ReportType       string             `json:"report_type"`
	BaselineMetrics  map[string]float64 `json:"baseline_metrics"`
	TrialMetrics     map[string]float64 `json:"trial_metrics"`
	ChangePct        map[string]float64 `json:"change_pct"`
	IsQualified      bool               `json:"is_qualified"`
	QualifyThreshold float64            `json:"qualify_threshold"`
	Recommendation   string             `json:"recommendation"`
	GeneratedAt      time.Time          `json:"generated_at"`
}

// ROIReportResp ROI报告响应
type ROIReportResp struct {
	ID                    string    `json:"id"`
	StrategyID            string    `json:"strategy_id"`
	CurrentBottleneck     string    `json:"current_bottleneck"`
	TotalInvestment       float64   `json:"total_investment"`
	ExpectedAnnualRevenue float64   `json:"expected_annual_revenue"`
	PaybackPeriodMonths   int       `json:"payback_period_months"`
	ApprovalResult        string    `json:"approval_result"`
	CreatedAt             time.Time `json:"created_at"`
}

// ScanResp 扫描报告响应
type ScanResp struct {
	ID        string    `json:"id"`
	ScanWeek  string    `json:"scan_week"`
	ScannedAt time.Time `json:"scanned_at"`
}
