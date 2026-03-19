// Package http 接口层 - 路由注册总入口
// 核心目的：注册所有HTTP路由，挂载中间件
// 模块功能：
//   - RegisterRoutes: 注册 /api/v1/ 下所有模块路由
//   - 挂载认证、限流、日志、CORS等全局中间件
package http
