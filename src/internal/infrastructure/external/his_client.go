// Package external 基础设施层 - HIS系统客户端
// 核心目的：对接HIS系统的HL7/RESTful接口
// 模块功能：
//   - 患者信息查询、检查项目同步、医嘱数据拉取
//   - 增量同步与全量同步
//   - 超时重试（3次，间隔1分钟）
package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HISPatient HIS 患者信息
type HISPatient struct {
	PatientID string    `json:"patient_id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	Phone     string    `json:"phone"`
	CardNo    string    `json:"card_no"` // 就诊卡号
}

// HISExamItem HIS 检查项目
type HISExamItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	DurationMin int    `json:"duration_min"`
	IsFasting   bool   `json:"is_fasting"`
}

// HISOrder HIS 医嘱/申请单
type HISOrder struct {
	OrderID   string    `json:"order_id"`
	PatientID string    `json:"patient_id"`
	ExamCode  string    `json:"exam_code"`
	DoctorID  string    `json:"doctor_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // pending / scheduled / completed
}

// HISClient HIS 系统 HTTP 客户端
type HISClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	maxRetry   int
}

// NewHISClient 创建 HIS 客户端
func NewHISClient(baseURL, apiKey string, timeout time.Duration, maxRetry int) *HISClient {
	return &HISClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetry: maxRetry,
	}
}

// GetPatient 查询患者信息
func (c *HISClient) GetPatient(ctx context.Context, patientID string) (*HISPatient, error) {
	var patient HISPatient
	if err := c.get(ctx, fmt.Sprintf("/patients/%s", patientID), &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

// ListExamItems 同步检查项目列表（全量）
func (c *HISClient) ListExamItems(ctx context.Context) ([]HISExamItem, error) {
	var items []HISExamItem
	if err := c.get(ctx, "/exam-items", &items); err != nil {
		return nil, err
	}
	return items, nil
}

// ListOrdersSince 拉取指定时间后的医嘱增量数据
func (c *HISClient) ListOrdersSince(ctx context.Context, since time.Time) ([]HISOrder, error) {
	url := fmt.Sprintf("/orders?since=%s", since.Format(time.RFC3339))
	var orders []HISOrder
	if err := c.get(ctx, url, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// ── 内部辅助 ─────────────────────────────────────────────────────

func (c *HISClient) get(ctx context.Context, path string, dest any) error {
	url := c.baseURL + path
	var lastErr error
	for i := 0; i <= c.maxRetry; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if i < c.maxRetry {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Minute):
				}
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("HIS API error %d: %s", resp.StatusCode, string(body))
		}
		return json.NewDecoder(resp.Body).Decode(dest)
	}
	return fmt.Errorf("HIS request failed after %d retries: %w", c.maxRetry, lastErr)
}
