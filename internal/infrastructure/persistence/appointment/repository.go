// Package appointment 基础设施层 - appointment仓储实现
// 核心目的：实现appointment领域层定义的仓储接口（基于GORM）
// 模块功能：
//   - CRUD操作的具体数据库实现
//   - PO ↔ 领域实体转换
//   - 复杂查询与分页
package appointment
