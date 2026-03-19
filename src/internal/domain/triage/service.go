// Package triage 分诊执行领域 - 领域服务
// 核心目的：实现签到、队列管理、状态流转的核心业务逻辑
// 模块功能：
//   - CheckInService: 多终端签到（二维码验证 + 时间窗口校验 + 队列入队）
//   - QueueManagementService: 队列管理（呼叫下一个 + 重叫 + 过号重排）
//   - ExamStatusService: 检查状态流转（状态机校验 + 误操作撤销）
//   - ScreenPushService: 分诊大屏推送数据组装
package triage
