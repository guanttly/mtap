// Package resource 资源管理领域 - 领域服务
// 核心目的：实现号源生成与排班管理的核心业务逻辑
// 模块功能：
//   - SlotGenerationService: 动态号源生成（按设备切分时段、年龄折算）
//   - ScheduleService: 排班管理（批量生成、临时停诊、替班、追加号源）
//   - SlotPoolService: 号源池管理（分配、锁定、释放、溢出）
//   - ItemMappingService: 项目别名映射解析
package resource
