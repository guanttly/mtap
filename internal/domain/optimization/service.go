// Package optimization 智能效能优化领域 - 领域服务
// 核心目的：实现效能优化闭环的核心业务逻辑
// 模块功能：
//   - AnomalyDetectionService: 异常检测引擎（均值±2σ检测、环比突变、趋势恶化）
//   - BottleneckAttributionService: 瓶颈归因分析（资源/流程/规则/行为四类归因）
//   - StrategyGenerationService: 优化策略生成（A/B/C分类，含预期收益量化）
//   - TrialExecutionService: 试运行执行（灰度下发、基线快照、即时回滚）
//   - TrialMonitorService: 试运行监控（实时对比、紧急回滚触发）
//   - EvaluationService: 效果评估（基线对比、达标判定、报告生成）
//   - DecayTrackingService: 策略衰减追踪（长期效果监测、衰减告警）
//   - PerformanceScanService: 周期性全量效能扫描
package optimization
