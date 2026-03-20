export type StrategyCategory = 'A' | 'B' | 'C'
export type StrategyStatus = 'pending_review' | 'rejected' | 'trial_running' | 'trial_running_b' | 'submitted_approval' | 'pending_eval' | 'promoted' | 'normalized' | 'rolled_back' | 'archived' | 'post_validation' | 'tracking' | 'decayed'
export type AlertType = 'consecutive_deviation' | 'sudden_change' | 'trend_degradation'

export interface EfficiencyMetric {
  id: string
  name: string
  code: string
  calc_formula: string
  unit: string
  normal_mean: number
  normal_std_dev: number
  normal_min: number
  normal_max: number
  is_custom: boolean
  latest_value?: number
  latest_sampled_at?: string
}

export interface MetricSnapshot {
  id: string
  metric_id: string
  value: number
  dimensions: Record<string, string>
  sampled_at: string
}

export interface BottleneckAlert {
  id: string
  metric_id: string
  metric_name?: string
  alert_type: AlertType
  severity: 'low' | 'medium' | 'high' | 'critical'
  deviation_pct: number
  consecutive_count: number
  affected_scope: string
  root_cause_hypotheses: string[]
  suggested_category: StrategyCategory
  status: 'open' | 'dismissed' | 'resolved'
  created_at: string
}

export interface OptimizationStrategy {
  id: string
  title: string
  category: StrategyCategory
  status: StrategyStatus
  alert_id: string
  current_value: string
  target_value: string
  expected_benefit: string
  risk_note: string
  trial_run?: TrialRun
  eval_report?: EvaluationReport
  roi_report?: ROIReport
  created_at: string
  updated_at: string
}

export interface TrialRun {
  id: string
  strategy_id: string
  trial_days: number
  started_at: string
  ends_at: string
  status: 'running' | 'completed' | 'rolled_back'
  emergency_rollback_threshold: number
  gray_scope: {
    department_ids: string[]
    device_ids: string[]
    time_periods: string[]
  }
}

export interface EvaluationReport {
  id: string
  strategy_id: string
  trial_run_id: string
  baseline_metrics: Record<string, number>
  trial_metrics: Record<string, number>
  change_pct: Record<string, number>
  is_qualified: boolean
  qualify_threshold: number
  recommendation: string
  generated_at: string
}

export interface ROIReport {
  id: string
  strategy_id: string
  current_bottleneck: string
  total_investment: number
  expected_annual_revenue: number
  payback_period_months: number
  risk_factors: string[]
  approval_result?: string
  pdf_path?: string
}

export interface PerformanceScan {
  id: string
  scan_week: string
  scanned_at: string
  opportunities: Array<{
    metric_code: string
    metric_name: string
    current_value: number
    normal_value: number
    deviation_pct: number
    suggested_category: StrategyCategory
  }>
}
