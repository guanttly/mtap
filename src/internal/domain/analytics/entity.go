// Package analytics 统计分析领域 - 实体定义
package analytics

import "time"

// DashboardSnapshot 大屏快照（聚合根）
type DashboardSnapshot struct {
	ID           string
	CampusID     string
	Timestamp    time.Time
	SlotUsage    SlotUsageData
	DeviceStatus []DeviceStatusData
	WaitTrend    []WaitTrendPoint
	Alerts       []AlertItem
}

// SlotUsageData 号源占用数据
type SlotUsageData struct {
	TotalSlots     int
	UsedSlots      int
	ExpiredSlots   int
	AvailableSlots int
	UsageRate      float64
}

// DeviceStatusData 设备状态数据
type DeviceStatusData struct {
	DeviceID   string
	DeviceName string
	Status     string // idle / in_use / maintenance
	QueueCount int
}

// WaitTrendPoint 等待时长趋势点
type WaitTrendPoint struct {
	Time       time.Time
	AvgWaitMin float64
}

// AlertItem 告警项
type AlertItem struct {
	Type     string // device_queue_overflow / slot_exhausted
	Message  string
	DeviceID string
	Value    int
}

// DeviceDetail 设备当日详情
type DeviceDetail struct {
	DeviceID     string
	DeviceName   string
	Date         time.Time
	TimeSlots    []SlotSummary
	QueueSummary QueueSummary
}

// SlotSummary 时段号源摘要
type SlotSummary struct {
	Hour      int
	Total     int
	Used      int
	UsageRate float64
}

// QueueSummary 队列摘要
type QueueSummary struct {
	TotalCheckedIn int
	AvgWaitMin     float64
	MaxWaitMin     float64
	NoShowCount    int
}

// DateRange 日期范围
type DateRange struct {
	Start time.Time
	End   time.Time
}

// Report 报表（聚合根）
type Report struct {
	ID          string
	ReportType  string // daily / weekly / monthly
	Dimensions  []string
	DateRange   DateRange
	Status      string // generating / ready / failed
	Format      string // xlsx / pdf
	FilePath    string
	FileSize    int64
	GeneratedAt *time.Time
	CreatedBy   string
	CreatedAt   time.Time
}
