export type DeviceStatus = 'online' | 'offline' | 'maintenance'
export type ScheduleStatus = 'normal' | 'suspended'
export type SlotStatus = 'available' | 'locked' | 'booked' | 'expired' | 'suspended'

export interface Device {
  id: string
  department_id: string
  name: string
  model: string
  manufacturer: string
  supported_exam_types: string[]
  max_daily_slots: number
  status: DeviceStatus
}

export interface ExamItem {
  id: string
  name: string
  duration_min: number
  is_fasting: boolean
  fasting_desc: string
  aliases: ItemAlias[]
}

export interface ItemAlias {
  alias: string
  exam_item_id: string
}

export interface SlotPool {
  id: string
  name: string
  pool_type: 'public' | 'department' | 'doctor'
  department_id: string
  allocation_ratio: number
  overflow_enabled: boolean
  overflow_target_id: string
}

export interface Schedule {
  id: string
  device_id: string
  device_name?: string
  work_date: string
  start_time: string
  end_time: string
  slot_minutes?: number
  status: ScheduleStatus
  slots?: TimeSlot[]
}

export interface TimeSlot {
  id: string
  schedule_id: string
  exam_item_id: string
  start_time: string
  end_time: string
  standard_duration: number
  adjusted_duration: number
  status: SlotStatus
  pool_type: string
  lock_until?: string
}
