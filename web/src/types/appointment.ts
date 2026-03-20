export type AppointmentMode = 'auto' | 'combo' | 'manual'
export type AppointmentStatus = 'pending' | 'confirmed' | 'pay_verifying' | 'paid' | 'rescheduling' | 'cancelled' | 'checked_in' | 'examining' | 'completed' | 'no_show' | 'released'
export type BlacklistStatus = 'active' | 'released' | 'expired'
export type AppealStatus = 'pending' | 'approved' | 'rejected'

export interface Appointment {
  id: string
  patient_id: string
  mode: AppointmentMode
  status: AppointmentStatus
  items: AppointmentItem[]
  override_by?: string
  override_reason?: string
  payment_verified: boolean
  change_count: number
  created_at: string
  confirmed_at?: string
  cancel_reason?: string
}

export interface AppointmentItem {
  id: string
  appointment_id: string
  exam_item_id: string
  exam_item_name?: string
  slot_id: string
  device_id: string
  device_name?: string
  start_time: string
  end_time: string
  status: string
}

export interface Credential {
  id: string
  appointment_id: string
  qr_code_data: string
  patient_name_masked: string
  exam_summary: string
  notice_content: string
  generated_at: string
}

export interface BlacklistRecord {
  id: string
  patient_id: string
  trigger_time: string
  expires_at: string
  status: BlacklistStatus
  no_show_count: number
  appeal?: Appeal
}

export interface NoShowRecord {
  id: string
  patient_id: string
  appointment_id: string
  occurred_at: string
}

export interface Appeal {
  id: string
  blacklist_id: string
  reason: string
  status: AppealStatus
  reviewed_by?: string
  reviewed_at?: string
}

export interface AutoAppointmentReq {
  patient_id: string
  exam_item_ids: string[]
  preferences?: SchedulePreference
}

export interface SchedulePreference {
  preferred_date_range?: { start: string, end: string }
  preferred_time_period?: string
}
