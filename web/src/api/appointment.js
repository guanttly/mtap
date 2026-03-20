import http from './request';
export const appointmentApi = {
    // 预约列表
    listAppointments: (params) => http.get('/appointments', { params }),
    getAppointment: (id) => http.get(`/appointments/${id}`),
    // 预约操作
    autoAppointment: (data) => http.post('/appointments/auto', data),
    comboAppointment: (data) => http.post('/appointments/combo', data),
    manualAppointment: (data) => http.post('/appointments/manual', data),
    reschedule: (id, data) => http.put(`/appointments/${id}/reschedule`, data),
    cancel: (id, reason) => http.put(`/appointments/${id}/cancel`, { reason }),
    confirm: (id) => http.put(`/appointments/${id}/confirm`),
    markPaid: (id) => http.put(`/appointments/${id}/paid`),
    // 凭证
    getCredential: (id) => http.get(`/appointments/${id}/credential`),
    // 黑名单
    listBlacklist: (params) => http.get('/appointments/blacklist', { params }),
    removeFromBlacklist: (id) => http.delete(`/appointments/blacklist/${id}`),
    // 申诉
    submitAppeal: (blacklistId, reason) => http.post(`/appointments/blacklist/${blacklistId}/appeals`, { reason }),
    reviewAppeal: (appealId, data) => http.put(`/appointments/appeals/${appealId}/review`, data),
};
