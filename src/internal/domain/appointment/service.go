// Package appointment 预约服务领域 - 领域服务
// 核心目的：实现预约编排与闭环管控的核心业务逻辑
// 模块功能：
//   - AutoAppointmentService: 一键自动预约（调用规则引擎 + 号源匹配 + 方案生成）
//   - ComboAppointmentService: 组合预约（多项目批量导入 + 方案对比）
//   - ManualOverrideService: 人工干预预约（权限校验 + 日志记录）
//   - RescheduleService: 改约服务（释放原号源 + 锁定新号源 + 事务回滚）
//   - BlacklistService: 黑名单管理（爽约计数 + 触发/解除 + 申诉审核）
//   - PaymentVerificationService: 缴费校验服务
package appointment
