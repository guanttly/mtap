import http from './request';
export const analyticsApi = {
    getDashboard: (campusId) => http.get('/analytics/dashboard', { params: { campus_id: campusId } }),
    getDeviceDetail: (deviceId, date) => http.get(`/analytics/dashboard/device/${deviceId}`, { params: { date } }),
    generateReport: (data) => http.post('/analytics/reports', data),
    getReport: (id) => http.get(`/analytics/reports/${id}`),
    listReports: (params) => http.get('/analytics/reports', { params }),
    exportReport: (id, format) => http.get(`/analytics/reports/${id}/export`, { params: { format }, responseType: 'blob' }),
};
