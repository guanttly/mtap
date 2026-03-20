// Package mq 基础设施层 - 消息发布者（Kafka）
// 核心目的：将领域事件发布到 Kafka Topic
// 模块功能：
//   - Publish: 根据路由键将事件发布到对应 Topic
//   - 消息序列化（JSON）
//   - 连接复用与优雅关闭
package mq

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// Topic Kafka Topic 常量
const (
	TopicNotification = "mtap.notification"
	TopicOptimization = "mtap.optimization"
	TopicAnalytics    = "mtap.analytics"
)

// 路由键常量（消息 Key，用于事件类型标识）
const (
	RoutingKeyAppointmentConfirmed = "appointment.confirmed"
	RoutingKeyAppointmentCancelled = "appointment.cancelled"
	RoutingKeyScheduleChanged      = "schedule.changed"
	RoutingKeyBottleneckDetected   = "bottleneck.detected"
	RoutingKeyTrialStarted         = "trial.started"
	RoutingKeyTrialCompleted       = "trial.completed"
	RoutingKeyStrategyApproved     = "strategy.approved"
	RoutingKeyStrategyPromoted     = "strategy.promoted"
)

// routingKeyToTopic 将路由键前缀映射到 Kafka Topic
func routingKeyToTopic(routingKey string) string {
	switch {
	case strings.HasPrefix(routingKey, "appointment."),
		strings.HasPrefix(routingKey, "schedule."):
		return TopicNotification
	case strings.HasPrefix(routingKey, "bottleneck."),
		strings.HasPrefix(routingKey, "trial."),
		strings.HasPrefix(routingKey, "strategy."):
		return TopicOptimization
	default:
		return TopicAnalytics
	}
}

// Event 领域事件
type Event struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	OccurredAt time.Time      `json:"occurred_at"`
	Payload    map[string]any `json:"payload"`
}

// Publisher Kafka 消息发布者
type Publisher struct {
	writers map[string]*kafka.Writer // topic -> writer
}

// NewPublisher 创建消息发布者
// brokers: Kafka broker 地址列表，如 []string{"localhost:9092"}
func NewPublisher(brokers []string) *Publisher {
	writers := map[string]*kafka.Writer{
		TopicNotification: newWriter(brokers, TopicNotification),
		TopicOptimization: newWriter(brokers, TopicOptimization),
		TopicAnalytics:    newWriter(brokers, TopicAnalytics),
	}
	return &Publisher{writers: writers}
}

func newWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireOne,
		WriteTimeout:           10 * time.Second,
	}
}

// Publish 发布事件到对应 Topic
// routingKey 用于确定目标 Topic 并作为消息 Key 便于分区路由
func (p *Publisher) Publish(ctx context.Context, routingKey string, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	topic := routingKeyToTopic(routingKey)
	w, ok := p.writers[topic]
	if !ok {
		w = p.writers[TopicAnalytics]
	}
	return w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(routingKey),
		Value: body,
	})
}

// Close 关闭所有 Writer
func (p *Publisher) Close() {
	for _, w := range p.writers {
		_ = w.Close()
	}
}
