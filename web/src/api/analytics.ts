import http, { type PageResult } from './request'
import type { DashboardSnapshot, DeviceDetail, GenerateReportRequest, Report } from '@/types/analytics'

export const analyticsApi = {
  getDashboard: (campusId?: string) =>
    http.get<any, DashboardSnapshot>('/analytics/dashboard', { params: { campus_id: campusId } }),
  getDeviceDetail: (deviceId: string, date: string) =>
    http.get<any, DeviceDetail>(`/analytics/dashboard/device/${deviceId}`, { params: { date } }),
  generateReport: (data: GenerateReportRequest) =>
    http.post<any, Report>('/analytics/reports', data),
  getReport: (id: string) =>
    http.get<any, Report>(`/analytics/reports/${id}`),
  listReports: (params?: Record<string, unknown>) =>
    http.get<any, PageResult<Report>>('/analytics/reports', { params }),
  exportReport: (id: string, format: string) =>
    http.get<any, Blob>(`/analytics/reports/${id}/export`, { params: { format }, responseType: 'blob' }),
}
