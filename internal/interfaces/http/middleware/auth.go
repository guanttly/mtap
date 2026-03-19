// Package middleware 接口层 - 认证中间件
// 核心目的：JWT Token验证与用户身份提取
// 模块功能：
//   - AuthMiddleware: 验证请求携带的JWT Token有效性
//   - 从Token中提取用户信息注入Context
//   - Token过期/无效时返回401
package middleware
