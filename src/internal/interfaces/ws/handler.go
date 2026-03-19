// Package ws 接口层 - WebSocket消息处理
// 核心目的：处理WebSocket连接与消息推送
// 模块功能：
//   - ScreenHandler: 分诊大屏推送（候诊队列 + 呼叫信息，10秒刷新）
//   - DashboardHandler: 监控大屏推送（号源占用率 + 设备状态 + 等待时长）
//   - 连接断开重连机制（每5秒重连）
package ws
