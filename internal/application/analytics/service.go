// Package analytics 应用层 - analytics应用服务
// 核心目的：编排analytics领域服务，管理事务边界，实现用例
// 模块功能：
//   - 接收接口层DTO，转换为领域对象
//   - 调用领域服务执行业务逻辑
//   - 管理数据库事务
//   - 发布领域事件
//   - 权限校验与审计日志
package analytics
