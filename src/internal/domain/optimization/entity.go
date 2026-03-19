// Package optimization 智能效能优化领域 - 实体定义
// 核心目的：定义效能优化闭环的核心实体与聚合根
// 模块功能：
//   - EfficiencyMetric: 效率指标聚合根（指标名称 + 计算口径 + 阈值 + 正常区间）
//   - MetricSnapshot: 指标快照实体（属于EfficiencyMetric，按采样周期存储数值）
//   - BottleneckAlert: 瓶颈告警聚合根（异常指标 + 偏离程度 + 持续时长 + 归因报告）
//   - OptimizationStrategy: 优化策略聚合根（核心聚合，含A/B/C分类 + 完整状态机）
//   - TrialRun: 试运行记录实体（灰度范围 + 周期 + 基线快照，属于Strategy）
//   - EvaluationReport: 评估报告实体（基线vs试运行对比 + 达标判定）
//   - ROIReport: ROI论证报告实体（C类专用，投资测算 + 收益预估 + 回收期）
//   - PerformanceScan: 周期性效能扫描结果聚合根
//   - StrategyDecayAlert: 策略衰减告警实体
package optimization
