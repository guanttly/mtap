import http from './request';
export const triageApi = {
    // 签到
    kioskCheckIn: (qrCodeData) => http.post('/triage/checkin', { method: 'kiosk', qr_code_data: qrCodeData }),
    nurseCheckIn: (data) => http.post('/triage/checkin', { method: 'nurse', ...data }),
    // 队列管理
    getQueueStatus: (roomId) => http.get(`/triage/queue/${roomId}`),
    callNext: (roomId) => http.post(`/triage/call/${roomId}/next`),
    recall: (roomId) => http.post(`/triage/call/${roomId}/recall`),
    missAndRequeue: (roomId) => http.post(`/triage/call/${roomId}/miss`),
    // 检查执行
    startExam: (id) => http.post(`/triage/exam/${id}/start`),
    completeExam: (id) => http.post(`/triage/exam/${id}/complete`),
    undoExam: (id, reason) => http.post(`/triage/exam/${id}/undo`, { reason }),
    getExamExecution: (id) => http.get(`/triage/exam/${id}`),
};
