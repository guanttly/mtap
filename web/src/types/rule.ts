export type ConflictLevel = 'forbid' | 'warning'
export type RuleStatus = 'active' | 'inactive'
export type DependencyType = 'mandatory' | 'recommended'

export interface ConflictRule {
  id: string
  item_a_id: string
  item_b_id: string
  item_a_name?: string
  item_b_name?: string
  min_interval: number
  interval_unit: string
  level: ConflictLevel
  status: RuleStatus
  created_by: string
  created_at: string
  updated_at: string
}

export interface ConflictPackage {
  id: string
  name: string
  min_interval: number
  interval_unit: string
  level: ConflictLevel
  status: RuleStatus
  items: ConflictPackageItem[]
  created_at: string
}

/** 创建冲突包请求（与后端 CreateConflictPackageReq 对齐） */
export interface CreateConflictPackageReq {
  name: string
  item_ids: string[]
  min_interval?: number
  interval_unit?: string
  level: ConflictLevel
}

/** 更新冲突包请求（与后端 UpdateConflictPackageReq 对齐） */
export interface UpdateConflictPackageReq {
  name?: string
  item_ids?: string[]
  min_interval?: number
  level?: ConflictLevel
}

export interface ConflictPackageItem {
  package_id: string
  exam_item_id: string
  exam_item_name?: string
}

export interface DependencyRule {
  id: string
  pre_item_id: string
  post_item_id: string
  pre_item_name?: string
  post_item_name?: string
  type: DependencyType
  validity_hours: number
  status: RuleStatus
  created_at: string
}

export interface PriorityTag {
  id: string
  name: string
  weight: number
  color: string
  is_preset: boolean
  created_at: string
}

export interface SortingStrategy {
  id: string
  type: 'shortest_wait' | 'nearest' | 'priority'
  scope_campuses: string[]
  scope_depts: string[]
  scope_devices: string[]
  start_date: string
  end_date: string
  status: RuleStatus
}

export interface PatientAdaptRule {
  id: string
  condition_type: 'age' | 'gender' | 'pregnancy'
  condition_value: string
  action: 'filter_device' | 'filter_slot' | 'filter_doctor'
  action_params: Record<string, string>
  priority: number
  status: RuleStatus
}

export interface SourceControl {
  id: string
  source_type: 'outpatient' | 'inpatient' | 'referral'
  slot_pool_id: string
  allocation_ratio: number
  overflow_enabled: boolean
  overflow_target_pool_id?: string
  status: RuleStatus
}

export interface CheckRulesReq {
  patient_id: string
  exam_item_ids: string[]
}

export interface CheckRulesResp {
  conflicts: ConflictResult[]
  dependencies: DependencyResult[]
  fasting_items: string[]
  warnings: string[]
}

export interface ConflictResult {
  item_a_id: string
  item_a_name: string
  item_b_id: string
  item_b_name: string
  level: ConflictLevel
  min_interval: number
  actual_interval: number
  reason: string
  rule_id: string
}

export interface DependencyResult {
  post_item_id: string
  post_item_name: string
  pre_item_id: string
  pre_item_name: string
  type: DependencyType
  status: 'passed' | 'blocked' | 'expired' | 'unknown'
  validity_hours: number
  completed_at?: string
  reason: string
}
