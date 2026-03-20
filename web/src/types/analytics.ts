export type ReportType = 'daily_summary' | 'device_efficiency' | 'item_distribution' | 'wait_trend' | 'no_show'
export type ReportStatus = 'generating' | 'ready' | 'failed'
export type ReportFormat = 'xlsx' | 'pdf' | 'csv'
export type AlertLevel = 'info' | 'warning' | 'critical'

export interface SlotUsage {
  total_slots: number
  used_slots: number
  available_slots: number
  usage_rate: number
}

export interface DeviceStatusItem {
  device_id: string
  device_name: string
  status: 'idle' | 'busy' | 'offline'
  queue_count: number
}

export interface WaitTrendPoint {
  time: string
  avg_wait_min: number
}

export interface DashboardAlert {
  type: AlertLevel
  message: string
  device_id?: string
  value?: number
}

export interface DashboardSnapshot {
  id: string
  timestamp: string
  slot_usage: SlotUsage
  device_status: DeviceStatusItem[]
  wait_trend: WaitTrendPoint[]
  alerts: DashboardAlert[]
}

export interface DeviceDetail {
  device_id: string
  device_name: string
  date: string
  total_exams: number
  avg_duration_min: number
  utilization_rate: number
  hourly_breakdown: Array<{ hour: number, count: number, avg_wait_min: number }>
}

export interface Report {
  id: string
  report_type: ReportType
  dimensions?: Record<string, unknown>
  date_start: string
  date_end: string
  status: ReportStatus
  file_path?: string
  file_size?: number
  format: ReportFormat
  generated_at?: string
  created_at: string
}

export interface GenerateReportRequest {
  report_type: ReportType
  date_start: string
  date_end: string
  dimensions?: Record<string, unknown>
  format?: ReportFormat
}
