// Package cache 基础设施层 - Redis缓存
// 核心目的：提供Redis缓存与分布式锁能力
// 模块功能：
//   - 号源锁定（分布式锁，防止超卖）
//   - 会话缓存（JWT Token黑名单）
//   - 热点数据缓存（大屏快照、指标数据）
//   - WebSocket频道消息中转（Redis Pub/Sub）
package cache
