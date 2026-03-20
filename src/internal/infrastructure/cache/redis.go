// Package cache 基础设施层 - Redis缓存
// 核心目的：提供Redis缓存与分布式锁能力
// 模块功能：
//   - 号源锁定（分布式锁，防止超卖）
//   - 会话缓存（JWT Token黑名单）
//   - 热点数据缓存（大屏快照、指标数据）
//   - WebSocket频道消息中转（Redis Pub/Sub）
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端封装
type Client struct {
	rdb *redis.Client
}

// NewClient 创建 Redis 客户端
func NewClient(addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

// Close 关闭连接
func (c *Client) Close() error { return c.rdb.Close() }

// ── 号源分布式锁 ─────────────────────────────────────────────────

// LockSlot 加号源锁（NX + PX）
// 成功返回 true，已被锁定返回 false
func (c *Client) LockSlot(ctx context.Context, slotID, patientID string, ttl time.Duration) (bool, error) {
	key := slotKey(slotID)
	ok, err := c.rdb.SetNX(ctx, key, patientID, ttl).Result()
	return ok, err
}

// UnlockSlot 释放号源锁（仅允许锁定者释放）
func (c *Client) UnlockSlot(ctx context.Context, slotID, patientID string) error {
	key := slotKey(slotID)
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // 已不存在
	}
	if err != nil {
		return err
	}
	if val != patientID {
		return fmt.Errorf("slot %s locked by another patient", slotID)
	}
	return c.rdb.Del(ctx, key).Err()
}

// ── JWT Token 黑名单 ─────────────────────────────────────────────

// BlacklistToken 将 Token 加入黑名单（登出时调用）
func (c *Client) BlacklistToken(ctx context.Context, jti string, expireAt time.Time) error {
	ttl := time.Until(expireAt)
	if ttl <= 0 {
		return nil
	}
	return c.rdb.Set(ctx, jwtBlackKey(jti), "1", ttl).Err()
}

// IsTokenBlacklisted 检查 Token 是否在黑名单
func (c *Client) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, jwtBlackKey(jti)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// ── 通用缓存 Get/Set/Del ─────────────────────────────────────────

// Set 设置缓存
func (c *Client) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存
func (c *Client) Get(ctx context.Context, key string, dest any) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

// Del 删除缓存
func (c *Client) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// ── Pub/Sub（WebSocket 频道中转）────────────────────────────────

// Publish 发布消息到频道
func (c *Client) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.rdb.Publish(ctx, channel, data).Err()
}

// Subscribe 订阅频道，返回消息 channel
func (c *Client) Subscribe(ctx context.Context, channels ...string) <-chan *redis.Message {
	sub := c.rdb.Subscribe(ctx, channels...)
	return sub.Channel()
}

// ── Sentinel 大屏快照缓存 ────────────────────────────────────────

// SetDashboardSnapshot 缓存大屏快照（TTL 10s 刷新）
func (c *Client) SetDashboardSnapshot(ctx context.Context, campusID string, data any) error {
	return c.Set(ctx, dashboardKey(campusID), data, 10*time.Second)
}

// GetDashboardSnapshot 获取大屏快照
func (c *Client) GetDashboardSnapshot(ctx context.Context, campusID string, dest any) error {
	return c.Get(ctx, dashboardKey(campusID), dest)
}

// ── 错误定义 ─────────────────────────────────────────────────────

// ErrCacheMiss 缓存未命中哨兵错误
var ErrCacheMiss = fmt.Errorf("cache miss")

// ── key 辅助函数 ─────────────────────────────────────────────────

func slotKey(slotID string) string        { return "slot:lock:" + slotID }
func jwtBlackKey(jti string) string       { return "jwt:black:" + jti }
func dashboardKey(campusID string) string { return "dashboard:snap:" + campusID }
