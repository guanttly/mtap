// Package main 通知工作器
// 订阅 MQ 通知队列，将预约/排班事件推送给患者（SMS/微信）
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	"github.com/euler/mtap/internal/infrastructure/external"
	"github.com/euler/mtap/internal/infrastructure/mq"
)

// notifyPayload 通知事件载荷（公共字段）
type notifyPayload struct {
	Phone           string `json:"phone"`
	PatientName     string `json:"patient_name"`
	ExamName        string `json:"exam_name"`
	AppointmentTime string `json:"appointment_time"`
	Reason          string `json:"reason"`
}

func main() {
	cfgPath := os.Getenv("MTAP_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}
	cfg, err := infraConfig.Load(cfgPath)
	if err != nil {
		log.Printf("load config failed (using defaults): %v", err)
		cfg, _ = infraConfig.Load("")
	}

	// 初始化消息客户端
	msgClient := external.NewMessageClient(
		cfg.Message.SMSKey, // smsGateway（使用 SMSKey 作为网关地址或空）
		"",                 // wechatURL（可后续配置）
		cfg.Message.SMSKey,
		cfg.Message.SMSSecret,
	)

	// 初始化 Kafka 消费者
	brokers := cfg.Kafka.Brokers
	if len(brokers) == 0 {
		brokers = []string{"localhost:9092"}
	}
	consumer, err := mq.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("[NotifyWorker] 初始化 Kafka 消费者失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 订阅通知 Topic
	groupID := cfg.Kafka.GroupID + ".notify"
	if err := consumer.Subscribe(ctx, mq.QueueNotification, groupID, makeHandler(msgClient)); err != nil {
		log.Fatalf("[NotifyWorker] 订阅通知 Topic 失败: %v", err)
	}

	log.Println("[NotifyWorker] 通知工作器已启动，等待消息...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[NotifyWorker] 收到退出信号，正在停止...")
	cancel()
	log.Println("[NotifyWorker] 已停止")
}

// makeHandler 构造 MQ 事件处理函数
func makeHandler(msgClient *external.MessageClient) mq.Handler {
	return func(ctx context.Context, event mq.Event) error {
		// 解析公共载荷
		raw, err := json.Marshal(event.Payload)
		if err != nil {
			log.Printf("[NotifyWorker] 序列化事件载荷失败: %v", err)
			return err
		}
		var p notifyPayload
		if err := json.Unmarshal(raw, &p); err != nil {
			log.Printf("[NotifyWorker] 解析通知载荷失败: %v", err)
			return err
		}

		switch event.Type {
		case mq.RoutingKeyAppointmentConfirmed:
			log.Printf("[NotifyWorker] 发送预约确认通知: %s %s", p.PatientName, p.Phone)
			return msgClient.SendAppointmentConfirmed(ctx, p.Phone, p.PatientName, p.ExamName, p.AppointmentTime)

		case mq.RoutingKeyAppointmentCancelled:
			log.Printf("[NotifyWorker] 发送预约取消通知: %s %s", p.PatientName, p.Phone)
			return msgClient.SendAppointmentCancelled(ctx, p.Phone, p.PatientName, p.ExamName)

		case mq.RoutingKeyScheduleChanged:
			log.Printf("[NotifyWorker] 发送排班变更通知: %s %s", p.PatientName, p.Phone)
			return msgClient.SendScheduleChanged(ctx, p.Phone, p.PatientName, p.Reason)

		default:
			log.Printf("[NotifyWorker] 未知事件类型 %s，跳过", event.Type)
			return nil
		}
	}
}
