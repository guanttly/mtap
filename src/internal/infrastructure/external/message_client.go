// Package external 基础设施层 - 消息平台客户端
// 核心目的：对接短信/微信消息推送平台
// 模块功能：
//   - 短信发送（主备通道切换）
//   - 微信模板消息推送
//   - 发送失败重试队列
package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SMSPayload 短信发送参数
type SMSPayload struct {
	PhoneNumbers []string          `json:"phone_numbers"`
	TemplateCode string            `json:"template_code"`
	Params       map[string]string `json:"params"`
}

// WechatPayload 微信模板消息参数
type WechatPayload struct {
	OpenID     string            `json:"open_id"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	URL        string            `json:"url,omitempty"`
}

// MessageClient 消息平台客户端
type MessageClient struct {
	smsGateway string // 短信网关 URL
	wechatURL  string // 微信推送服务 URL
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

// NewMessageClient 创建消息客户端
func NewMessageClient(smsGateway, wechatURL, apiKey, apiSecret string) *MessageClient {
	return &MessageClient{
		smsGateway: smsGateway,
		wechatURL:  wechatURL,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SendSMS 发送短信
func (c *MessageClient) SendSMS(ctx context.Context, payload SMSPayload) error {
	if c.smsGateway == "" {
		return nil // 未配置网关时静默跳过（开发模式）
	}
	return c.post(ctx, c.smsGateway+"/sms/send", payload)
}

// SendWechat 发送微信模板消息
func (c *MessageClient) SendWechat(ctx context.Context, payload WechatPayload) error {
	if c.wechatURL == "" {
		return nil
	}
	return c.post(ctx, c.wechatURL+"/message/template", payload)
}

// SendAppointmentConfirmed 发送预约确认通知
func (c *MessageClient) SendAppointmentConfirmed(ctx context.Context, phone, patientName, examName, appointmentTime string) error {
	return c.SendSMS(ctx, SMSPayload{
		PhoneNumbers: []string{phone},
		TemplateCode: "SMS_APPT_CONFIRMED",
		Params: map[string]string{
			"name": patientName,
			"exam": examName,
			"time": appointmentTime,
		},
	})
}

// SendAppointmentCancelled 发送预约取消通知
func (c *MessageClient) SendAppointmentCancelled(ctx context.Context, phone, patientName, examName string) error {
	return c.SendSMS(ctx, SMSPayload{
		PhoneNumbers: []string{phone},
		TemplateCode: "SMS_APPT_CANCELLED",
		Params: map[string]string{
			"name": patientName,
			"exam": examName,
		},
	})
}

// SendScheduleChanged 发送排班变更通知
func (c *MessageClient) SendScheduleChanged(ctx context.Context, phone, patientName, reason string) error {
	return c.SendSMS(ctx, SMSPayload{
		PhoneNumbers: []string{phone},
		TemplateCode: "SMS_SCHEDULE_CHANGED",
		Params: map[string]string{
			"name":   patientName,
			"reason": reason,
		},
	})
}

func (c *MessageClient) post(ctx context.Context, url string, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("message client request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("message api error %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
