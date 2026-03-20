import type { Appeal, Appointment, AutoAppointmentReq, BlacklistRecord, Credential } from '@/types/appointment'
import http, { type PageResult } from './request'

export const appointmentApi = {
  // 预约列表
  listAppointments: (params?: Record<string, unknown>) => http.get<any, PageResult<Appointment>>('/appointments', { params }),
  getAppointment: (id: string) => http.get<any, Appointment>(`/appointments/${id}`),

  // 预约操作
  autoAppointment: (data: AutoAppointmentReq) => http.post<any, { plans: unknown[], failed_items: unknown[] }>('/appointments/auto', data),
  comboAppointment: (data: unknown) => http.post<any, Appointment>('/appointments/combo', data),
  manualAppointment: (data: unknown) => http.post<any, Appointment>('/appointments/manual', data),
  reschedule: (id: string, data: unknown) => http.put<any, Appointment>(`/appointments/${id}/reschedule`, data),
  cancel: (id: string, reason: string) => http.put<any, void>(`/appointments/${id}/cancel`, { reason }),
  confirm: (id: string) => http.put<any, void>(`/appointments/${id}/confirm`),
  markPaid: (id: string) => http.put<any, void>(`/appointments/${id}/paid`),

  // 凭证
  getCredential: (id: string) => http.get<any, Credential>(`/appointments/${id}/credential`),

  // 黑名单
  listBlacklist: (params?: Record<string, unknown>) => http.get<any, PageResult<BlacklistRecord>>('/appointments/blacklist', { params }),
  removeFromBlacklist: (id: string) => http.delete<any, void>(`/appointments/blacklist/${id}`),

  // 申诉
  submitAppeal: (blacklistId: string, reason: string) => http.post<any, Appeal>(`/appointments/blacklist/${blacklistId}/appeals`, { reason }),
  reviewAppeal: (appealId: string, data: { approved: boolean, comment: string }) =>
    http.put<any, void>(`/appointments/appeals/${appealId}/review`, data),
}
