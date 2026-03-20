// Package mq 基础设施层 - 消息消费者（Kafka）
// 核心目的：从 Kafka Topic 消费领域事件
// 模块功能：
//   - 通知类事件消费（短信/微信推送）
//   - 统计类事件消费（数据采集入库）
//   - 消费确认与错误处理（提交 offset）
package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// 队列/消费者组名称常量（对应原 RabbitMQ queue 概念）
const (
	QueueNotification = TopicNotification
	QueueAnalytics    = TopicAnalytics
	QueueOptimization = TopicOptimization
)

// Handler 事件处理函数类型
type Handler func(ctx context.Context, event Event) error

// Consumer Kafka 消费者
type Consumer struct {
	brokers []string
	readers []*kafka.Reader
}

// NewConsumer 创建消费者
// brokers: Kafka broker 地址列表，如 []string{"localhost:9092"}
func NewConsumer(brokers []string) (*Consumer, error) {
	return &Consumer{brokers: brokers}, nil
}

// Subscribe 订阅 Topic
// topic: Kafka Topic 名称（使用 TopicXxx 常量）
// groupID: 消费者组 ID，相同组内负载均衡消费
// handler: 事件处理函数
func (c *Consumer) Subscribe(ctx context.Context, topic, groupID string, handler Handler) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       1e6, // 1MB
		MaxWait:        500 * time.Millisecond,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})
	c.readers = append(c.readers, r)

	go func() {
		defer r.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			m, err := r.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("[MQ] FetchMessage error (topic=%s): %v", topic, err)
				time.Sleep(time.Second)
				continue
			}

			var event Event
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("[MQ] 解析消息失败 (topic=%s offset=%d): %v", topic, m.Offset, err)
				// 解析失败直接提交，避免无限卡住
				_ = r.CommitMessages(ctx, m)
				continue
			}

			if err := handler(ctx, event); err != nil {
				log.Printf("[MQ] 处理消息失败 type=%s: %v", event.Type, err)
				// Kafka 无 Nack，重试由业务层决定；此处仍提交 offset 避免 poison pill
				// 生产环境可写入死信 Topic
			}
			_ = r.CommitMessages(ctx, m)
		}
	}()
	return nil
}

// Close 关闭所有 Reader
func (c *Consumer) Close() {
	for _, r := range c.readers {
		_ = r.Close()
	}
}
