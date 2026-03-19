// Package appointment 预约服务领域 - 值对象
// 核心目的：定义预约领域的不可变值对象
// 模块功能：
//   - AppointmentPlan: 预约方案（时间 + 设备 + 地点组合）
//   - ChangeLog: 改约/取消变更记录（操作时间 + 操作人 + 原号源 + 新号源）
//   - AppointmentMode: 预约模式枚举（自动/组合/人工干预）
//   - AppointmentStatus: 预约状态枚举
package appointment
