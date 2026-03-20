// Package external 基础设施层 - 收费系统客户端
// 核心目的：对接HIS收费系统验证缴费状态
// 模块功能：
//   - 缴费状态查询（RESTful JSON）
//   - 超时处理（5秒阈值，降级为"待校验"）
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PaymentStatus 缴费状态
type PaymentStatus string

const (
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPending PaymentStatus = "pending" // 降级状态：待校验
)

// PaymentInfo 缴费信息
type PaymentInfo struct {
	OrderID   string        `json:"order_id"`
	PatientID string        `json:"patient_id"`
	Amount    float64       `json:"amount"`
	Status    PaymentStatus `json:"status"`
	PaidAt    *time.Time    `json:"paid_at,omitempty"`
}

// PaymentClient 收费系统客户端
type PaymentClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewPaymentClient 创建收费系统客户端
func NewPaymentClient(baseURL, apiKey string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // 5秒超时，超时降级
		},
	}
}

// QueryPaymentStatus 查询缴费状态
// 超时情况下返回 PaymentStatusPending（降级处理）
func (c *PaymentClient) QueryPaymentStatus(ctx context.Context, appointmentID string) (PaymentStatus, *PaymentInfo, error) {
	if c.baseURL == "" {
		// 未配置时返回已缴费（开发/测试模式）
		return PaymentStatusPaid, nil, nil
	}

	url := fmt.Sprintf("%s/payments/appointment/%s", c.baseURL, appointmentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return PaymentStatusPending, nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// 超时降级为待校验
		return PaymentStatusPending, nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return PaymentStatusUnpaid, nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return PaymentStatusPending, nil, fmt.Errorf("payment api error %d: %s", resp.StatusCode, string(body))
	}

	var info PaymentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return PaymentStatusPending, nil, err
	}
	return info.Status, &info, nil
}

// VerifyPayment 校验预约缴费是否有效（简化接口）
func (c *PaymentClient) VerifyPayment(ctx context.Context, appointmentID string) (bool, error) {
	status, _, err := c.QueryPaymentStatus(ctx, appointmentID)
	if err != nil {
		return false, err
	}
	return status == PaymentStatusPaid, nil
}
