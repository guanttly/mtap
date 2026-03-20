import http from './request';
export const optimizationApi = {
    // 指标
    listMetrics: () => http.get('/optimization/metrics'),
    getMetricTrend: (code, params) => http.get(`/optimization/metrics/${code}/trend`, { params }),
    // 告警
    listAlerts: (params) => http.get('/optimization/alerts', { params }),
    dismissAlert: (id, reason) => http.put(`/optimization/alerts/${id}/dismiss`, { reason }),
    // 策略
    listStrategies: (params) => http.get('/optimization/strategies', { params }),
    getStrategy: (id) => http.get(`/optimization/strategies/${id}`),
    approveStrategy: (id, data) => http.post(`/optimization/strategies/${id}/approve`, data),
    rejectStrategy: (id, reason) => http.post(`/optimization/strategies/${id}/reject`, { reason }),
    rollbackStrategy: (id) => http.post(`/optimization/strategies/${id}/rollback`),
    promoteStrategy: (id) => http.post(`/optimization/strategies/${id}/promote`),
    // 试运行
    getTrialMonitor: (id) => http.get(`/optimization/trials/${id}/monitor`),
    // 评估报告
    getEvaluation: (id) => http.get(`/optimization/evaluations/${id}`),
    // ROI报告
    getROIReport: (id) => http.get(`/optimization/roi-reports/${id}`),
    submitROIResult: (id, data) => http.post(`/optimization/roi-reports/${id}/result`, data),
    // 周期扫描
    listScans: (params) => http.get('/optimization/scans', { params }),
    getScan: (id) => http.get(`/optimization/scans/${id}`),
};
