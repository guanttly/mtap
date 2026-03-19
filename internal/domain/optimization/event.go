// Package optimization 智能效能优化领域 - 领域事件
// 核心目的：定义效能优化的领域事件
// 模块功能：
//   - BottleneckDetected: 异常检测发现瓶颈
//   - StrategyApproved: 管理员批准策略
//   - TrialStarted: 策略试运行开始
//   - TrialEmergencyRollback: 关键指标恶化超阈值触发紧急回滚
//   - TrialCompleted: 试运行周期到期
//   - StrategyPromoted: 策略转正/常态化
//   - StrategyDecayed: 已转正策略效果衰减
//   - PerformanceScanCompleted: 周期效能扫描完成
package optimization
