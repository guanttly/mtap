// Package optimization 智能效能优化领域 - 值对象
// 核心目的：定义效能优化领域的不可变值对象
// 模块功能：
//   - StrategyCategory: 策略分类枚举（A类软策略/B类弹性资源策略/C类硬资源策略）
//   - StrategyStatus: 策略状态枚举（待审核/试运行中/试行中/已提交审批/待评估/已转正/已常态化/已回滚/已驳回/已归档/长期追踪/策略衰减）
//   - BaselineSnapshot: 基线快照（试运行前各指标数值）
//   - GrayScope: 灰度范围（科室/设备/时段维度限定）
//   - ResourceActionList: 资源调配执行清单/归还清单（B类专用）
//   - CostEstimate: 成本评估值对象（人力/耗材/能耗）
//   - ROIMetrics: 投入产出比指标
package optimization
