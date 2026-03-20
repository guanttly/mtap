// Package analytics 应用层 - 统计分析数据传输对象
package analytics

import "time"

// ─── 请求 DTO ────────────────────────────────────────────────────────────

// GetDashboardReq 大屏数据请求
type GetDashboardReq struct {
	CampusID string `form:"campus_id"`
}

// GetDeviceDetailReq 设备详情请求
type GetDeviceDetailReq struct {
	Date string `form:"date"` // YYYY-MM-DD
}

// CreateReportReq 生成报表请求
type CreateReportReq struct {
	ReportType string   `json:"report_type" binding:"required,oneof=daily weekly monthly"`
	Metrics    []string `json:"metrics"`
	CampusID   string   `json:"campus_id"`
	DateStart  string   `json:"date_start" binding:"required"`
	DateEnd    string   `json:"date_end" binding:"required"`
	Format     string   `json:"format" binding:"omitempty,oneof=xlsx pdf"`
}

// ListReportsReq 报表列表请求
type ListReportsReq struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────

// SlotUsageResp 号源占用响应
type SlotUsageResp struct {
	TotalSlots     int     `json:"total_slots"`
	UsedSlots      int     `json:"used_slots"`
	ExpiredSlots   int     `json:"expired_slots"`
	AvailableSlots int     `json:"available_slots"`
	UsageRate      float64 `json:"usage_rate"`
}

// DeviceStatusResp 设备状态响应
type DeviceStatusResp struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Status     string `json:"status"`
	QueueCount int    `json:"queue_count"`
}

// WaitTrendResp 等待趋势响应
type WaitTrendResp struct {
	Time       time.Time `json:"time"`
	AvgWaitMin float64   `json:"avg_wait_min"`
}

// AlertItemResp 告警响应
type AlertItemResp struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	DeviceID string `json:"device_id,omitempty"`
	Value    int    `json:"value,omitempty"`
}

// DashboardResp 大屏数据响应
type DashboardResp struct {
	Timestamp    time.Time          `json:"timestamp"`
	CampusID     string             `json:"campus_id"`
	SlotUsage    SlotUsageResp      `json:"slot_usage"`
	DeviceStatus []DeviceStatusResp `json:"device_status"`
	WaitTrend    []WaitTrendResp    `json:"wait_trend"`
	Alerts       []AlertItemResp    `json:"alerts"`
}

// SlotSummaryResp 时段号源摘要响应
type SlotSummaryResp struct {
	Hour      int     `json:"hour"`
	Total     int     `json:"total"`
	Used      int     `json:"used"`
	UsageRate float64 `json:"usage_rate"`
}

// QueueSummaryResp 队列摘要响应
type QueueSummaryResp struct {
	TotalCheckedIn int     `json:"total_checked_in"`
	AvgWaitMin     float64 `json:"avg_wait_min"`
	MaxWaitMin     float64 `json:"max_wait_min"`
	NoShowCount    int     `json:"no_show_count"`
}

// DeviceDetailResp 设备详情响应
type DeviceDetailResp struct {
	DeviceID     string            `json:"device_id"`
	DeviceName   string            `json:"device_name"`
	Date         string            `json:"date"`
	TimeSlots    []SlotSummaryResp `json:"time_slots"`
	QueueSummary QueueSummaryResp  `json:"queue_summary"`
}

// ReportResp 报表响应
type ReportResp struct {
	ID          string     `json:"id"`
	ReportType  string     `json:"report_type"`
	Status      string     `json:"status"`
	Format      string     `json:"format"`
	FilePath    string     `json:"file_path,omitempty"`
	FileSize    int64      `json:"file_size,omitempty"`
	GeneratedAt *time.Time `json:"generated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
