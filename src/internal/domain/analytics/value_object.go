// Package analytics 统计分析领域 - 值对象
// 核心目的：定义统计分析领域的枚举常量与不可变值对象
// 模块功能：
//   - ReportType:   报表类型（日报/周报/月报）
//   - ReportStatus: 报表生成状态（生成中/就绪/失败）
//   - ExportFormat: 导出格式（xlsx/pdf）
//   - AlertType:    告警类型（设备队列溢出/号源耗尽）
//   - MetricDimension: 统计指标维度
package analytics

// ReportType 报表类型
type ReportType string

const (
	ReportTypeDaily   ReportType = "daily"   // 日报
	ReportTypeWeekly  ReportType = "weekly"  // 周报
	ReportTypeMonthly ReportType = "monthly" // 月报
)

// ReportStatus 报表状态
type ReportStatus string

const (
	ReportStatusGenerating ReportStatus = "generating" // 生成中
	ReportStatusReady      ReportStatus = "ready"      // 就绪（可下载）
	ReportStatusFailed     ReportStatus = "failed"     // 生成失败
)

// ExportFormat 导出格式
type ExportFormat string

const (
	ExportFormatXLSX ExportFormat = "xlsx"
	ExportFormatPDF  ExportFormat = "pdf"
)

// AlertType 告警类型
type AlertType string

const (
	AlertTypeDeviceQueueOverflow AlertType = "device_queue_overflow" // 设备队列溢出（>20人）
	AlertTypeSlotExhausted       AlertType = "slot_exhausted"        // 号源耗尽（剩余<5%）
)

// MetricDimension 统计指标维度
type MetricDimension string

const (
	MetricSlotUsage    MetricDimension = "slot_usage"    // 号源利用率
	MetricDeviceUsage  MetricDimension = "device_usage"  // 设备使用率
	MetricAvgWait      MetricDimension = "avg_wait"      // 平均等待时长
	MetricNoShowRate   MetricDimension = "no_show_rate"  // 爽约率
	MetricOverrideRate MetricDimension = "override_rate" // 人工干预率
)
