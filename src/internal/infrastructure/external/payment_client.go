// Package external 基础设施层 - 收费系统客户端
// 核心目的：对接HIS收费系统验证缴费状态
// 模块功能：
//   - 缴费状态查询（RESTful JSON）
//   - 超时处理（5秒阈值，降级为"待校验"）
package external
