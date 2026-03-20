// Package optimization 应用层 - 效能优化装配器
package optimization

import (
	domain "github.com/euler/mtap/internal/domain/optimization"
)

// ToMetricResp 指标实体 → DTO
func ToMetricResp(m *domain.EfficiencyMetric) *MetricResp {
	return &MetricResp{
		ID: m.ID, Name: m.Name, Code: m.Code, Unit: m.Unit,
		NormalMin: m.NormalMin, NormalMax: m.NormalMax, NormalMean: m.NormalMean,
	}
}

// ToAlertResp 告警实体 → DTO
func ToAlertResp(a *domain.BottleneckAlert) *AlertResp {
	return &AlertResp{
		ID: a.ID, MetricID: a.MetricID, AlertType: string(a.AlertType),
		Severity: a.Severity, DeviationPct: a.DeviationPct,
		ConsecutiveCount:  a.ConsecutiveCount,
		SuggestedCategory: string(a.SuggestedCategory),
		Status:            string(a.Status), DismissReason: a.DismissReason, CreatedAt: a.CreatedAt,
	}
}

// ToStrategyResp 策略实体 → DTO
func ToStrategyResp(s *domain.OptimizationStrategy) *StrategyResp {
	approvers := make([]ApprovalNodeResp, 0, len(s.ApprovalFlow.Approvers))
	for _, node := range s.ApprovalFlow.Approvers {
		approvers = append(approvers, ApprovalNodeResp{
			ApproverRole: node.ApproverRole, Status: node.Status,
			Timestamp: node.Timestamp, Comment: node.Comment,
		})
	}
	return &StrategyResp{
		ID: s.ID, Title: s.Title, Category: string(s.Category),
		Status: string(s.Status), AlertID: s.AlertID,
		CurrentValue: s.CurrentValue, TargetValue: s.TargetValue,
		ExpectedBenefit: s.ExpectedBenefit, RiskNote: s.RiskNote,
		ApprovalType: s.ApprovalFlow.Type, Approvers: approvers,
		RejectReason: s.RejectReason, PromotedAt: s.PromotedAt, CreatedAt: s.CreatedAt,
	}
}

// ToTrialRunResp 试运行实体 → DTO
func ToTrialRunResp(t *domain.TrialRun) *TrialRunResp {
	return &TrialRunResp{
		ID: t.ID, StrategyID: t.StrategyID, TrialDays: t.TrialDays,
		StartedAt: t.StartedAt, EndsAt: t.EndsAt, Status: string(t.Status),
	}
}

// ToEvaluationResp 评估报告实体 → DTO
func ToEvaluationResp(e *domain.EvaluationReport) *EvaluationResp {
	return &EvaluationResp{
		ID: e.ID, StrategyID: e.StrategyID, ReportType: e.ReportType,
		BaselineMetrics: e.BaselineMetrics, TrialMetrics: e.TrialMetrics, ChangePct: e.ChangePct,
		IsQualified: e.IsQualified, QualifyThreshold: e.QualifyThreshold,
		Recommendation: e.Recommendation, GeneratedAt: e.GeneratedAt,
	}
}

// ToROIReportResp ROI报告实体 → DTO
func ToROIReportResp(r *domain.ROIReport) *ROIReportResp {
	return &ROIReportResp{
		ID: r.ID, StrategyID: r.StrategyID,
		CurrentBottleneck:     r.CurrentBottleneck,
		TotalInvestment:       r.TotalInvestment,
		ExpectedAnnualRevenue: r.ExpectedAnnualRevenue,
		PaybackPeriodMonths:   r.PaybackPeriodMonths,
		ApprovalResult:        r.ApprovalResult, CreatedAt: r.CreatedAt,
	}
}

// ToScanResp 扫描实体 → DTO
func ToScanResp(s *domain.PerformanceScan) *ScanResp {
	return &ScanResp{ID: s.ID, ScanWeek: s.ScanWeek, ScannedAt: s.ScannedAt}
}
