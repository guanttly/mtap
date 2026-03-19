// Package ws 接口层 - WebSocket连接管理
// 核心目的：管理WebSocket连接的生命周期
// 模块功能：
//   - Hub: 连接池管理（注册/注销/广播）
//   - 按频道分组（分诊大屏频道、监控大屏频道）
//   - 心跳检测与自动重连
package ws
