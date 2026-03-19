// Package appointment 预约服务领域 - 领域事件
// 核心目的：定义预约服务的领域事件
// 模块功能：
//   - AppointmentConfirmed: 预约确认锁定号源
//   - AppointmentCancelled: 取消预约
//   - NoShowTriggered: 爽约发生
//   - BlacklistTriggered: 累计爽约达阈值
package appointment
