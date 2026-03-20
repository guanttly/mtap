// Package po 基础设施层 - optimization 持久化对象
package po

import "time"

// EfficiencyMetricPO 效率指标
type EfficiencyMetricPO struct {
	ID           string    `gorm:"primaryKey;column:id;size:36"`
	Name         string    `gorm:"column:name;size:50;not null"`
	Code         string    `gorm:"column:code;size:30;not null;uniqueIndex"`
	CalcFormula  string    `gorm:"column:calc_formula;type:text"`
	Unit         string    `gorm:"column:unit;size:10"`
	NormalMean   float64   `gorm:"column:normal_mean"`
	NormalStdDev float64   `gorm:"column:normal_stddev"`
	NormalMin    float64   `gorm:"column:normal_min"`
	NormalMax    float64   `gorm:"column:normal_max"`
	IsCustom     bool      `gorm:"column:is_custom;not null;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (EfficiencyMetricPO) TableName() string { return "efficiency_metrics" }

// MetricSnapshotPO 指标快照
type MetricSnapshotPO struct {
	ID         string    `gorm:"primaryKey;column:id;size:36"`
	MetricID   string    `gorm:"column:metric_id;size:36;not null;index:idx_metric_snap_metric_time"`
	Value      float64   `gorm:"column:value;not null"`
	Dimensions string    `gorm:"column:dimensions;type:text"` // JSON
	SampledAt  time.Time `gorm:"column:sampled_at;not null;index:idx_metric_snap_metric_time"`
}

func (MetricSnapshotPO) TableName() string { return "metric_snapshots" }

// BottleneckAlertPO 瓶颈告警
type BottleneckAlertPO struct {
	ID                  string    `gorm:"primaryKey;column:id;size:36"`
	MetricID            string    `gorm:"column:metric_id;size:36;not null"`
	AlertType           string    `gorm:"column:alert_type;size:30;not null"`
	Severity            string    `gorm:"column:severity;size:10;not null;default:warning"`
	DeviationPct        float64   `gorm:"column:deviation_pct;not null"`
	ConsecutiveCount    int       `gorm:"column:consecutive_count"`
	AffectedScope       string    `gorm:"column:affected_scope;type:text"`        // JSON
	RootCauseHypotheses string    `gorm:"column:root_cause_hypotheses;type:text"` // JSON
	SuggestedCategory   string    `gorm:"column:suggested_category;size:1"`
	Status              string    `gorm:"column:status;size:15;not null;default:active"`
	DismissReason       string    `gorm:"column:dismiss_reason;size:500"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (BottleneckAlertPO) TableName() string { return "bottleneck_alerts" }

// OptimizationStrategyPO 优化策略
type OptimizationStrategyPO struct {
	ID              string     `gorm:"primaryKey;column:id;size:36"`
	Title           string     `gorm:"column:title;size:100;not null"`
	Category        string     `gorm:"column:category;size:1;not null;index:idx_strategy_category"`
	Status          string     `gorm:"column:status;size:25;not null;default:pending_review;index:idx_strategy_status"`
	AlertID         string     `gorm:"column:alert_id;size:36"`
	CurrentValue    string     `gorm:"column:current_value;type:text"`
	TargetValue     string     `gorm:"column:target_value;type:text"`
	ExpectedBenefit string     `gorm:"column:expected_benefit;type:text"`
	RiskNote        string     `gorm:"column:risk_note;type:text"`
	CostEstimate    string     `gorm:"column:cost_estimate;type:text"` // JSON
	ApprovalFlow    string     `gorm:"column:approval_flow;type:text"` // JSON
	RejectReason    string     `gorm:"column:reject_reason;size:500"`
	CooldownUntil   *time.Time `gorm:"column:cooldown_until"`
	PromotedAt      *time.Time `gorm:"column:promoted_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (OptimizationStrategyPO) TableName() string { return "optimization_strategies" }

// TrialRunPO 试运行记录
type TrialRunPO struct {
	ID                         string     `gorm:"primaryKey;column:id;size:36"`
	StrategyID                 string     `gorm:"column:strategy_id;size:36;not null"`
	GrayScope                  string     `gorm:"column:gray_scope;type:text"` // JSON
	TrialDays                  int        `gorm:"column:trial_days;not null"`
	BaselineSnapshot           string     `gorm:"column:baseline_snapshot;type:text"` // JSON
	ConfigBackup               string     `gorm:"column:config_backup;type:text"`     // JSON
	EmergencyRollbackThreshold float64    `gorm:"column:emergency_threshold;default:0.15"`
	StartedAt                  time.Time  `gorm:"column:started_at;not null"`
	EndsAt                     time.Time  `gorm:"column:ends_at;not null"`
	Status                     string     `gorm:"column:status;size:15;not null;default:running"`
	RolledBackAt               *time.Time `gorm:"column:rolled_back_at"`
	RollbackReason             string     `gorm:"column:rollback_reason;size:500"`
}

func (TrialRunPO) TableName() string { return "trial_runs" }

// EvaluationReportPO 评估报告
type EvaluationReportPO struct {
	ID               string    `gorm:"primaryKey;column:id;size:36"`
	StrategyID       string    `gorm:"column:strategy_id;size:36;not null"`
	TrialRunID       string    `gorm:"column:trial_run_id;size:36"`
	ReportType       string    `gorm:"column:report_type;size:20;not null;default:trial"`
	BaselineMetrics  string    `gorm:"column:baseline_metrics;type:text"` // JSON
	TrialMetrics     string    `gorm:"column:trial_metrics;type:text"`    // JSON
	ChangePct        string    `gorm:"column:change_pct;type:text"`       // JSON
	IsQualified      bool      `gorm:"column:is_qualified"`
	QualifyThreshold float64   `gorm:"column:qualify_threshold"`
	ActualCost       *float64  `gorm:"column:actual_cost"`
	ActualROI        *float64  `gorm:"column:actual_roi"`
	Recommendation   string    `gorm:"column:recommendation;size:20"`
	GeneratedAt      time.Time `gorm:"column:generated_at;not null"`
}

func (EvaluationReportPO) TableName() string { return "evaluation_reports" }

// ROIReportPO ROI 论证报告
type ROIReportPO struct {
	ID                    string    `gorm:"primaryKey;column:id;size:36"`
	StrategyID            string    `gorm:"column:strategy_id;size:36;not null"`
	CurrentBottleneck     string    `gorm:"column:current_bottleneck;type:text;not null"`
	TriedOptimizations    string    `gorm:"column:tried_optimizations;type:text"` // JSON
	InvestmentItems       string    `gorm:"column:investment_items;type:text"`    // JSON
	TotalInvestment       float64   `gorm:"column:total_investment;not null"`
	ExpectedAnnualRevenue float64   `gorm:"column:expected_annual_rev"`
	PaybackPeriodMonths   int       `gorm:"column:payback_months"`
	RiskFactors           string    `gorm:"column:risk_factors;type:text"` // JSON
	PDFPath               string    `gorm:"column:pdf_path;size:500"`
	ApprovalResult        string    `gorm:"column:approval_result;size:15"`
	ApprovalResultBy      string    `gorm:"column:approval_result_by;size:36"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ROIReportPO) TableName() string { return "roi_reports" }

// PerformanceScanPO 周期效能扫描
type PerformanceScanPO struct {
	ID            string    `gorm:"primaryKey;column:id;size:36"`
	ScanWeek      string    `gorm:"column:scan_week;size:10;not null"`
	Metrics       string    `gorm:"column:metrics;type:text"`       // JSON
	Opportunities string    `gorm:"column:opportunities;type:text"` // JSON
	ScannedAt     time.Time `gorm:"column:scanned_at;not null;autoCreateTime"`
}

func (PerformanceScanPO) TableName() string { return "performance_scans" }

// StrategyDecayAlertPO 策略衰减告警
type StrategyDecayAlertPO struct {
	ID                  string    `gorm:"primaryKey;column:id;size:36"`
	StrategyID          string    `gorm:"column:strategy_id;size:36;not null"`
	MetricCode          string    `gorm:"column:metric_code;size:30;not null"`
	OriginalImprovement float64   `gorm:"column:original_improvement;not null"`
	CurrentImprovement  float64   `gorm:"column:current_improvement;not null"`
	DecayPct            float64   `gorm:"column:decay_pct;not null"`
	DetectedAt          time.Time `gorm:"column:detected_at;not null;autoCreateTime"`
}

func (StrategyDecayAlertPO) TableName() string { return "strategy_decay_alerts" }
