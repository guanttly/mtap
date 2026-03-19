// Package triage 分诊执行领域 - 值对象
// 核心目的：定义分诊领域的不可变值对象
// 模块功能：
//   - CheckInMethod: 签到方式枚举（自助机扫码/护士站手动/NFC）
//   - ExamStatus: 检查状态枚举（已签到/候诊中/检查中/检查完成）
//   - CallResult: 呼叫结果枚举（到达/重叫/过号）
package triage
