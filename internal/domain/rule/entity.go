// Package rule 规则引擎领域 - 实体定义
// 核心目的：定义预约规则引擎的核心实体与聚合根
// 模块功能：
//   - ConflictRule: 冲突规则聚合根（项目对 + 最小间隔 + 冲突级别）
//   - ConflictPackage: 冲突包聚合根（一组互斥检查项目的逻辑分组）
//   - DependencyRule: 依赖规则聚合根（前置项目 + 后续项目 + 依赖类型）
//   - PriorityTag: 优先级标签聚合根（名称 + 权重）
//   - SortingStrategy: 排序策略聚合根（类型 + 范围 + 时段）
//   - PatientAdaptRule: 患者属性适配规则聚合根
//   - SourceControl: 开单来源控制规则聚合根
package rule
