import http, { type PageResult } from './request'
import type {
  ConflictRule,
  ConflictPackage,
  CreateConflictPackageReq,
  UpdateConflictPackageReq,
  DependencyRule,
  PriorityTag,
  SortingStrategy,
  CheckRulesReq,
  CheckRulesResp,
  PatientAdaptRule,
  SourceControl,
} from '@/types/rule'

export const ruleApi = {
  // 冲突规则
  createConflictRule: (data: Partial<ConflictRule>) => http.post<any, ConflictRule>('/rules/conflicts', data),
  listConflictRules: (params?: Record<string, unknown>) => http.get<any, PageResult<ConflictRule>>('/rules/conflicts', { params }),
  getConflictRule: (id: string) => http.get<any, ConflictRule>(`/rules/conflicts/${id}`),
  updateConflictRule: (id: string, data: { min_interval?: number, level?: string, status?: string }) => http.put<any, ConflictRule>(`/rules/conflicts/${id}`, data),
  deleteConflictRule: (id: string) => http.delete<any, void>(`/rules/conflicts/${id}`),

  // 冲突包
  createConflictPackage: (data: CreateConflictPackageReq) => http.post<any, ConflictPackage>('/rules/conflict-packages', data),
  listConflictPackages: (params?: Record<string, unknown>) => http.get<any, PageResult<ConflictPackage>>('/rules/conflict-packages', { params }),
  updateConflictPackage: (id: string, data: UpdateConflictPackageReq) => http.put<any, ConflictPackage>(`/rules/conflict-packages/${id}`, data),
  deleteConflictPackage: (id: string) => http.delete<any, void>(`/rules/conflict-packages/${id}`),

  // 依赖规则
  createDependencyRule: (data: Partial<DependencyRule>) => http.post<any, DependencyRule>('/rules/dependencies', data),
  listDependencyRules: (params?: Record<string, unknown>) => http.get<any, PageResult<DependencyRule>>('/rules/dependencies', { params }),
  updateDependencyRule: (id: string, data: { type?: string, validity_hours?: number, status?: string }) => http.put<any, DependencyRule>(`/rules/dependencies/${id}`, data),
  deleteDependencyRule: (id: string) => http.delete<any, void>(`/rules/dependencies/${id}`),

  // 优先级标签
  createPriorityTag: (data: Partial<PriorityTag>) => http.post<any, PriorityTag>('/rules/priority-tags', data),
  listPriorityTags: () => http.get<any, PageResult<PriorityTag>>('/rules/priority-tags'),
  updatePriorityTag: (id: string, data: { name?: string, weight?: number, color?: string }) => http.put<any, PriorityTag>(`/rules/priority-tags/${id}`, data),
  deletePriorityTag: (id: string) => http.delete<any, void>(`/rules/priority-tags/${id}`),

  // 排序策略
  saveSortingStrategy: (data: Partial<SortingStrategy>) => http.post<any, SortingStrategy>('/rules/sorting-strategies', data),
  getSortingStrategy: () => http.get<any, SortingStrategy>('/rules/sorting-strategies'),

  // 患者属性适配
  listPatientAdaptRules: () => http.get<any, { items: PatientAdaptRule[] }>('/rules/patient-adapt'),
  savePatientAdaptRules: (data: Partial<PatientAdaptRule>[]) => http.post<any, void>('/rules/patient-adapt', data),

  // 设置开单来源控制
  saveSourceControls: (data: Partial<SourceControl>[]) => http.post<any, void>('/rules/source-controls', data),
  listSourceControls: () => http.get<any, { items: SourceControl[] }>('/rules/source-controls'),

  // 规则综合校验
  checkRules: (data: CheckRulesReq) => http.post<any, CheckRulesResp>('/rules/check', data),
}
