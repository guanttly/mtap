import type { Device, ExamItem, ItemAlias, Schedule, SlotPool, TimeSlot } from '@/types/resource'
import http, { type PageResult } from './request'

export const resourceApi = {
  // 设备
  listDevices: (params?: Record<string, unknown>) => http.get<any, PageResult<Device>>('/resources/devices', { params }),
  getDevice: (id: string) => http.get<any, Device>(`/resources/devices/${id}`),
  createDevice: (data: Partial<Device>) => http.post<any, Device>('/resources/devices', data),
  updateDevice: (id: string, data: Partial<Device>) => http.put<any, Device>(`/resources/devices/${id}`, data),
  deleteDevice: (id: string) => http.delete<any, void>(`/resources/devices/${id}`),

  // 检查项目
  listExamItems: (params?: Record<string, unknown>) => http.get<any, PageResult<ExamItem>>('/resources/exam-items', { params }),
  createExamItem: (data: Partial<ExamItem>) => http.post<any, ExamItem>('/resources/exam-items', data),
  updateExamItem: (id: string, data: Partial<ExamItem>) => http.put<any, ExamItem>(`/resources/exam-items/${id}`, data),
  deleteExamItem: (id: string) => http.delete<any, void>(`/resources/exam-items/${id}`),

  // 项目别名
  addAlias: (examItemId: string, alias: string) => http.post<any, ItemAlias>(`/resources/exam-items/${examItemId}/aliases`, { alias }),
  deleteAlias: (examItemId: string, alias: string) => http.delete<any, void>(`/resources/exam-items/${examItemId}/aliases/${alias}`),

  // 号源池
  listSlotPools: () => http.get<any, PageResult<SlotPool>>('/resources/slot-pools'),
  createSlotPool: (data: Partial<SlotPool>) => http.post<any, SlotPool>('/resources/slot-pools', data),
  updateSlotPool: (id: string, data: Partial<SlotPool>) => http.put<any, SlotPool>(`/resources/slot-pools/${id}`, data),

  // 院区 & 科室
  listCampuses: () => http.get<any, { items: Array<{ id: string, name: string, code: string }> }>('/resources/campuses'),
  listDepartments: (campusId?: string) => http.get<any, { items: Array<{ id: string, name: string, campus_id: string }> }>('/resources/departments', { params: { campus_id: campusId } }),

  // 排班
  listSchedules: (params?: Record<string, unknown>) => http.get<any, PageResult<Schedule>>('/resources/schedules', { params }),
  createSchedule: (data: Partial<Schedule>) => http.post<any, Schedule>('/resources/schedules', data),
  generateSchedule: (data: {
    device_id: string
    start_date: string
    end_date?: string
    date?: string
    start_time: string
    end_time: string
    slot_minutes: number
    exam_item_id?: string
    pool_type?: string
    skip_weekends?: boolean
  }) => http.post<any, { generated_count: number, total_slots: number }>('/resources/schedules/generate', data),
  suspendSchedule: (deviceId: string, date: string, startTime: string, endTime: string, reason: string) =>
    http.post<any, { released_slots: number }>('/resources/schedules/suspend', { device_id: deviceId, date, start_time: startTime, end_time: endTime, reason }),
  substituteSchedule: (sourceDeviceId: string, targetDeviceId: string, date: string) =>
    http.post<any, { moved_slots: number }>('/resources/schedules/substitute', { source_device_id: sourceDeviceId, target_device_id: targetDeviceId, date }),
  addExtraSlots: (scheduleId: string, data: unknown) => http.post<any, TimeSlot[]>('/resources/schedules/add-slots', data),

  // 号源
  listTimeSlots: (params?: Record<string, unknown>) => http.get<any, PageResult<TimeSlot>>('/resources/slots', { params }),
  listAvailableSlots: (params?: Record<string, unknown>) => http.get<any, TimeSlot[]>('/resources/slots/available', { params }),
  lockSlot: (id: string, patientId: string) => http.post<any, void>(`/resources/slots/${id}/lock`, { patient_id: patientId }),
  releaseSlot: (id: string) => http.post<any, void>(`/resources/slots/${id}/release`),
}
