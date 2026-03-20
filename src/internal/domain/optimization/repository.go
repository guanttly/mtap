// Package optimization 智能效能优化领域 - 仓储接口
package optimization

import "context"

// EfficiencyMetricRepository 效率指标仓储
type EfficiencyMetricRepository interface {
	Save(ctx context.Context, m *EfficiencyMetric) error
	FindByID(ctx context.Context, id string) (*EfficiencyMetric, error)
	FindByCode(ctx context.Context, code string) (*EfficiencyMetric, error)
	List(ctx context.Context) ([]*EfficiencyMetric, error)
	Update(ctx context.Context, m *EfficiencyMetric) error
}

// MetricSnapshotRepository 指标快照仓储
type MetricSnapshotRepository interface {
	Save(ctx context.Context, s *MetricSnapshot) error
	FindByMetricID(ctx context.Context, metricID string, limit int) ([]*MetricSnapshot, error)
	FindRecent90Days(ctx context.Context, metricID string) ([]*MetricSnapshot, error)
}

// BottleneckAlertRepository 瓶颈告警仓储
type BottleneckAlertRepository interface {
	Save(ctx context.Context, a *BottleneckAlert) error
	FindByID(ctx context.Context, id string) (*BottleneckAlert, error)
	List(ctx context.Context, status string, page, size int) ([]*BottleneckAlert, int64, error)
	Update(ctx context.Context, a *BottleneckAlert) error
}

// OptimizationStrategyRepository 优化策略仓储
type OptimizationStrategyRepository interface {
	Save(ctx context.Context, s *OptimizationStrategy) error
	FindByID(ctx context.Context, id string) (*OptimizationStrategy, error)
	List(ctx context.Context, category, status string, page, size int) ([]*OptimizationStrategy, int64, error)
	Update(ctx context.Context, s *OptimizationStrategy) error
	ListActiveTrials(ctx context.Context) ([]*OptimizationStrategy, error)
	ListPromoted(ctx context.Context) ([]*OptimizationStrategy, error)
	CountPendingReview(ctx context.Context) (int64, error)
}

// TrialRunRepository 试运行仓储
type TrialRunRepository interface {
	Save(ctx context.Context, t *TrialRun) error
	FindByID(ctx context.Context, id string) (*TrialRun, error)
	FindByStrategyID(ctx context.Context, strategyID string) (*TrialRun, error)
	Update(ctx context.Context, t *TrialRun) error
}

// EvaluationReportRepository 评估报告仓储
type EvaluationReportRepository interface {
	Save(ctx context.Context, r *EvaluationReport) error
	FindByID(ctx context.Context, id string) (*EvaluationReport, error)
	FindByStrategyID(ctx context.Context, strategyID string) (*EvaluationReport, error)
}

// ROIReportRepository ROI 报告仓储
type ROIReportRepository interface {
	Save(ctx context.Context, r *ROIReport) error
	FindByID(ctx context.Context, id string) (*ROIReport, error)
	FindByStrategyID(ctx context.Context, strategyID string) (*ROIReport, error)
	Update(ctx context.Context, r *ROIReport) error
}

// PerformanceScanRepository 效能扫描仓储
type PerformanceScanRepository interface {
	Save(ctx context.Context, s *PerformanceScan) error
	FindByID(ctx context.Context, id string) (*PerformanceScan, error)
	List(ctx context.Context, page, size int) ([]*PerformanceScan, int64, error)
}

// StrategyDecayAlertRepository 策略衰减告警仓储
type StrategyDecayAlertRepository interface {
	Save(ctx context.Context, a *StrategyDecayAlert) error
	FindByStrategyID(ctx context.Context, strategyID string) ([]*StrategyDecayAlert, error)
}
