// Package analytics 统计分析领域 - 领域服务
// 核心目的：实现实时监控数据聚合与报表生成的业务逻辑
// 模块功能：
//   - DashboardService: 大屏数据实时聚合（号源利用率、设备状态、等待时长）
//   - ReportService: 报表生成与导出（Excel/PDF，异步模式）
//   - AlertService: 告警阈值检测（设备等待队列超限等）
package analytics
