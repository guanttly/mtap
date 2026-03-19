// Package rule 规则引擎领域 - 仓储接口
// 核心目的：定义规则引擎聚合根的持久化接口（依赖倒置）
// 模块功能：
//   - ConflictRuleRepository: 冲突规则仓储
//   - ConflictPackageRepository: 冲突包仓储
//   - DependencyRuleRepository: 依赖规则仓储
//   - PriorityTagRepository: 优先级标签仓储
//   - SortingStrategyRepository: 排序策略仓储
//   - PatientAdaptRuleRepository: 患者属性适配规则仓储
//   - SourceControlRepository: 开单来源控制仓储
package rule
