import http, { type PageResult } from './request'
import type { EfficiencyMetric, BottleneckAlert, OptimizationStrategy, TrialRun, EvaluationReport, ROIReport, PerformanceScan } from '@/types/optimization'

export const optimizationApi = {
  // 指标
  listMetrics: () => http.get<any, PageResult<EfficiencyMetric>>('/optimization/metrics'),
  getMetricTrend: (code: string, params?: Record<string, unknown>) => http.get<any, unknown[]>(`/optimization/metrics/${code}/trend`, { params }),

  // 告警
  listAlerts: (params?: Record<string, unknown>) => http.get<any, PageResult<BottleneckAlert>>('/optimization/alerts', { params }),
  dismissAlert: (id: string, reason: string) => http.put<any, void>(`/optimization/alerts/${id}/dismiss`, { reason }),

  // 策略
  listStrategies: (params?: Record<string, unknown>) => http.get<any, PageResult<OptimizationStrategy>>('/optimization/strategies', { params }),
  getStrategy: (id: string) => http.get<any, OptimizationStrategy>(`/optimization/strategies/${id}`),
  approveStrategy: (id: string, data: { trial_days?: number, gray_scope?: unknown }) =>
    http.post<any, void>(`/optimization/strategies/${id}/approve`, data),
  rejectStrategy: (id: string, reason: string) => http.post<any, void>(`/optimization/strategies/${id}/reject`, { reason }),
  rollbackStrategy: (id: string) => http.post<any, void>(`/optimization/strategies/${id}/rollback`),
  promoteStrategy: (id: string) => http.post<any, void>(`/optimization/strategies/${id}/promote`),

  // 试运行
  getTrialMonitor: (id: string) => http.get<any, TrialRun>(`/optimization/trials/${id}/monitor`),

  // 评估报告
  getEvaluation: (id: string) => http.get<any, EvaluationReport>(`/optimization/evaluations/${id}`),

  // ROI报告
  getROIReport: (id: string) => http.get<any, ROIReport>(`/optimization/roi-reports/${id}`),
  submitROIResult: (id: string, data: { approved: boolean, comment: string }) =>
    http.post<any, void>(`/optimization/roi-reports/${id}/result`, data),

  // 周期扫描
  listScans: (params?: Record<string, unknown>) => http.get<any, PageResult<PerformanceScan>>('/optimization/scans', { params }),
  getScan: (id: string) => http.get<any, PerformanceScan>(`/optimization/scans/${id}`),
}
