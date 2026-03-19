// Package resource 资源管理领域 - 实体定义
// 核心目的：定义基础信息与资源排班管理的核心实体
// 模块功能：
//   - Device: 设备聚合根（归属科室/院区、支持检查类型列表）
//   - Schedule: 排班计划聚合根（设备 + 日期 + 工作时段 + 排班模板）
//   - TimeSlot: 号源时段实体（属于Schedule，起止时间 + 关联检查类型 + 剩余量）
//   - SlotPool: 号源池聚合根（公共池/科室池/医生专池，含配额与溢出规则）
//   - ExamItem: 检查项目聚合根（标准耗时、空腹标记）
//   - Campus: 院区实体
//   - Department: 科室实体
//   - Doctor: 医生实体
package resource
