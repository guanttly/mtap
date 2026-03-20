export type CheckInMethod = 'kiosk' | 'nurse' | 'nfc'
export type EntryStatus = 'waiting' | 'calling' | 'examining' | 'completed' | 'missed' | 'no_show'
export type ExamStatus = 'checked_in' | 'waiting' | 'ongoing' | 'done'

export interface CheckInResult {
  check_in_id: string
  queue_number: number
  estimated_wait: number
  room_location: string
  is_late: boolean
}

export interface QueueEntry {
  id: string
  queue_id: string
  patient_id: string
  patient_name_masked: string
  appointment_id: string
  queue_number: number
  status: EntryStatus
  call_count: number
  miss_count: number
  entered_at: string
  called_at?: string
}

export interface QueueStatus {
  room_id: string
  room_name: string
  current_calling?: {
    entry_id: string
    patient_name_masked: string
    queue_number: number
    call_count: number
  }
  waiting_count: number
  average_wait: number
  entries: QueueEntry[]
}

export interface CallResult {
  entry_id: string
  patient_name_masked: string
  queue_number: number
  room_name: string
  call_count: number
}

export interface ExamExecution {
  id: string
  appointment_item_id: string
  patient_id: string
  device_id: string
  status: ExamStatus
  started_at?: string
  completed_at?: string
  duration?: number
  operator_id: string
}
