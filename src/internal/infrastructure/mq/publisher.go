// Package mq 基础设施层 - 消息发布者
// 核心目的：将领域事件发布到消息队列
// 模块功能：
//   - PublishEvent: 发布领域事件到RabbitMQ
//   - 消息序列化（JSON）
//   - 发布确认与重试
package mq
