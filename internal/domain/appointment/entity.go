// Package appointment 预约服务领域 - 实体定义
// 核心目的：定义预约全生命周期的核心实体与聚合根
// 模块功能：
//   - Appointment: 预约单聚合根（核心聚合，含状态机：待确认→已确认→已签到→检查中→已完成/已取消/爽约）
//   - Credential: 预约凭证实体（二维码 + 注意事项模板）
//   - Blacklist: 黑名单聚合根（患者 + 触发时间 + 有效期 + 状态）
//   - NoShowRecord: 爽约记录实体（属于Blacklist统计）
package appointment
