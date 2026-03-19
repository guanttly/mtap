// Package resource 资源管理领域 - 值对象
// 核心目的：定义资源管理领域的不可变值对象
// 模块功能：
//   - ItemAlias: 项目别名映射（标准名称 + 别名列表）
//   - WorkPeriod: 工作时段（起止时间）
//   - ScheduleTemplate: 排班模板（周重复/自定义周期）
//   - SlotAllocation: 号源配额分配（门诊/住院/公共比例）
//   - AgeFactor: 年龄折算系数（儿童+10%/老年+15%）
package resource
