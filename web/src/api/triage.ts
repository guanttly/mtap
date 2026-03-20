import http from './request'
import type { CheckInResult, QueueStatus, CallResult, ExamExecution } from '@/types/triage'

export const triageApi = {
  // 签到
  kioskCheckIn: (qrCodeData: string) => http.post<any, CheckInResult>('/triage/checkin', { method: 'kiosk', qr_code_data: qrCodeData }),
  nurseCheckIn: (data: { appointment_id: string, patient_id: string, remark?: string }) =>
    http.post<any, CheckInResult>('/triage/checkin', { method: 'nurse', ...data }),

  // 队列管理
  getQueueStatus: (roomId: string) => http.get<any, QueueStatus>(`/triage/queue/${roomId}`),
  callNext: (roomId: string) => http.post<any, CallResult>(`/triage/call/${roomId}/next`),
  recall: (roomId: string) => http.post<any, CallResult>(`/triage/call/${roomId}/recall`),
  missAndRequeue: (roomId: string) => http.post<any, { miss_count: number, is_no_show: boolean }>(`/triage/call/${roomId}/miss`),

  // 检查执行
  startExam: (id: string) => http.post<any, ExamExecution>(`/triage/exam/${id}/start`),
  completeExam: (id: string) => http.post<any, ExamExecution>(`/triage/exam/${id}/complete`),
  undoExam: (id: string, reason: string) => http.post<any, void>(`/triage/exam/${id}/undo`, { reason }),
  getExamExecution: (id: string) => http.get<any, ExamExecution>(`/triage/exam/${id}`),
}
