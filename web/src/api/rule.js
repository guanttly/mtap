import http from './request'

export const ruleApi = {
  // 冲突规则
  createConflictRule: data => http.post('/rules/conflicts', data),
  listConflictRules: params => http.get('/rules/conflicts', { params }),
  getConflictRule: id => http.get(`/rules/conflicts/${id}`),
  updateConflictRule: (id, data) => http.put(`/rules/conflicts/${id}`, data),
  deleteConflictRule: id => http.delete(`/rules/conflicts/${id}`),
  // 冲突包
  createConflictPackage: data => http.post('/rules/conflict-packages', data),
  listConflictPackages: params => http.get('/rules/conflict-packages', { params }),
  updateConflictPackage: (id, data) => http.put(`/rules/conflict-packages/${id}`, data),
  deleteConflictPackage: id => http.delete(`/rules/conflict-packages/${id}`),
  // 依赖规则
  createDependencyRule: data => http.post('/rules/dependencies', data),
  listDependencyRules: params => http.get('/rules/dependencies', { params }),
  updateDependencyRule: (id, data) => http.put(`/rules/dependencies/${id}`, data),
  deleteDependencyRule: id => http.delete(`/rules/dependencies/${id}`),
  // 优先级标签
  createPriorityTag: data => http.post('/rules/priority-tags', data),
  listPriorityTags: () => http.get('/rules/priority-tags'),
  updatePriorityTag: (id, data) => http.put(`/rules/priority-tags/${id}`, data),
  deletePriorityTag: id => http.delete(`/rules/priority-tags/${id}`),
  // 排序策略
  saveSortingStrategy: data => http.post('/rules/sorting-strategies', data),
  listSortingStrategies: () => http.get('/rules/sorting-strategies'),
  deleteSortingStrategy: id => http.delete(`/rules/sorting-strategies/${id}`),
  // 患者属性适配
  listPatientAdaptRules: () => http.get('/rules/patient-adapt'),
  savePatientAdaptRules: data => http.post('/rules/patient-adapt', data),
  // 设置开单来源控制
  saveSourceControls: data => http.post('/rules/source-controls', data),
  listSourceControls: () => http.get('/rules/source-controls'),
  // 规则综合校验
  checkRules: data => http.post('/rules/check', data),
}
