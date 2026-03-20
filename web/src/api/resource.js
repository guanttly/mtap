import http from './request';
export const resourceApi = {
    // 设备
    listDevices: (params) => http.get('/resources/devices', { params }),
    getDevice: (id) => http.get(`/resources/devices/${id}`),
    createDevice: (data) => http.post('/resources/devices', data),
    updateDevice: (id, data) => http.put(`/resources/devices/${id}`, data),
    deleteDevice: (id) => http.delete(`/resources/devices/${id}`),
    // 检查项目
    listExamItems: (params) => http.get('/resources/exam-items', { params }),
    createExamItem: (data) => http.post('/resources/exam-items', data),
    updateExamItem: (id, data) => http.put(`/resources/exam-items/${id}`, data),
    deleteExamItem: (id) => http.delete(`/resources/exam-items/${id}`),
    // 项目别名
    addAlias: (examItemId, alias) => http.post(`/resources/exam-items/${examItemId}/aliases`, { alias }),
    deleteAlias: (examItemId, alias) => http.delete(`/resources/exam-items/${examItemId}/aliases/${alias}`),
    // 号源池
    listSlotPools: () => http.get('/resources/slot-pools'),
    createSlotPool: (data) => http.post('/resources/slot-pools', data),
    updateSlotPool: (id, data) => http.put(`/resources/slot-pools/${id}`, data),
    // 院区 & 科室
    listCampuses: () => http.get('/resources/campuses'),
    listDepartments: (campusId) => http.get('/resources/departments', { params: { campus_id: campusId } }),
    // 排班
    listSchedules: (params) => http.get('/resources/schedules', { params }),
    createSchedule: (data) => http.post('/resources/schedules', data),
    generateSchedule: (data) => http.post('/resources/schedules/generate', data),
    suspendSchedule: (deviceId, date, startTime, endTime, reason) => http.post('/resources/schedules/suspend', { device_id: deviceId, date, start_time: startTime, end_time: endTime, reason }),
    substituteSchedule: (sourceDeviceId, targetDeviceId, date) => http.post('/resources/schedules/substitute', { source_device_id: sourceDeviceId, target_device_id: targetDeviceId, date }),
    addExtraSlots: (scheduleId, data) => http.post('/resources/schedules/add-slots', data),
    // 号源
    listTimeSlots: (params) => http.get('/resources/slots', { params }),
    listAvailableSlots: (params) => http.get('/resources/slots/available', { params }),
    lockSlot: (id, patientId) => http.post(`/resources/slots/${id}/lock`, { patient_id: patientId }),
    releaseSlot: (id) => http.post(`/resources/slots/${id}/release`),
};
