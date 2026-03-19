// Package rule 规则引擎领域 - 领域服务
// 核心目的：实现规则引擎的核心业务逻辑
// 模块功能：
//   - ConflictDetectionService: 冲突检测服务（遍历冲突包、检查互斥关系）
//   - DependencyValidationService: 依赖校验服务（前置条件时效验证）
//   - TimeWindowCalculator: 多项目时间窗口计算（最优组合方案生成）
//   - PatientAdaptService: 患者属性适配服务（设备/时段过滤）
//   - PriorityScorer: 优先级排序权重计算
package rule
