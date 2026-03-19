// Package external 基础设施层 - HIS系统客户端
// 核心目的：对接HIS系统的HL7/RESTful接口
// 模块功能：
//   - 患者信息查询、检查项目同步、医嘱数据拉取
//   - 增量同步与全量同步
//   - 超时重试（3次，间隔1分钟）
package external
