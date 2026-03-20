// Package analytics 统计分析领域 - 领域事件
package analytics

import "time"

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
}

// SnapshotSavedEvent 大屏快照已保存事件
// 触发时机：DashboardService.GetSnapshot 落库成功后
type SnapshotSavedEvent struct {
	SnapshotID string    `json:"snapshot_id"`
	CampusID   string    `json:"campus_id"`
	AlertCount int       `json:"alert_count"` // 本次快照中告警数
	OccurredAt time.Time `json:"occurred_at"`
}

func (e SnapshotSavedEvent) EventName() string { return "analytics.snapshot_saved" }

// ReportGeneratedEvent 报表生成完成事件
// 触发时机：ReportService.Generate 完成后
type ReportGeneratedEvent struct {
	ReportID   string    `json:"report_id"`
	ReportType string    `json:"report_type"`
	Format     string    `json:"format"`
	CreatedBy  string    `json:"created_by"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e ReportGeneratedEvent) EventName() string { return "analytics.report_generated" }

// AlertTriggeredEvent 监控告警触发事件
// 触发时机：大屏快照中检测到异常告警项时
type AlertTriggeredEvent struct {
	AlertType  string    `json:"alert_type"` // device_queue_overflow / slot_exhausted
	Message    string    `json:"message"`
	CampusID   string    `json:"campus_id"`
	DeviceID   string    `json:"device_id,omitempty"`
	Value      int       `json:"value,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e AlertTriggeredEvent) EventName() string { return "analytics.alert_triggered" }
