// Package triage 分诊执行领域 - 实体定义
// 核心目的：定义签到、候诊、呼叫、检查执行全流程的核心实体
// 模块功能：
//   - CheckIn: 签到记录聚合根（签到方式 + 时间 + 备注）
//   - WaitingQueue: 候诊队列聚合根（按诊室维度，管理队列排序）
//   - QueueEntry: 队列条目实体（属于WaitingQueue，含签到时间 + 呼叫次数 + 过号次数）
//   - ExamExecution: 检查执行聚合根（状态流转：已签到→候诊中→检查中→检查完成）
package triage
