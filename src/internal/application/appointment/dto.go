// Package appointment 应用层 - 预约服务数据传输对象
package appointment

import "time"

// === 请求 DTO ===

// AutoAppointmentReq 一键自动预约请求
type AutoAppointmentReq struct {
	PatientID   string   `json:"patient_id"    binding:"required"`
	ExamItemIDs []string `json:"exam_item_ids" binding:"required,min=1,max=10"`
}

// ComboAppointmentReq 组合预约计算请求
type ComboAppointmentReq struct {
	PatientID    string   `json:"patient_id"    binding:"required"`
	ExamItemIDs  []string `json:"exam_item_ids" binding:"required,min=1,max=10"`
	PreferPeriod string   `json:"prefer_period"` // morning / afternoon / any
	StartDate    string   `json:"start_date"`    // 格式 2006-01-02，默认今日
	EndDate      string   `json:"end_date"`      // 格式 2006-01-02，默认今日+90天
}

// ComboConfirmReq 确认组合预约方案请求
type ComboConfirmReq struct {
	AppointmentID string `json:"appointment_id" binding:"required"`
	PlanID        string `json:"plan_id"        binding:"required"`
}

// ManualAppointmentReq 人工干预预约请求
type ManualAppointmentReq struct {
	PatientID   string `json:"patient_id"   binding:"required"`
	ExamItemID  string `json:"exam_item_id" binding:"required"`
	SlotID      string `json:"slot_id"      binding:"required"`
	Reason      string `json:"reason"       binding:"required,max=200"`
	AckConflict bool   `json:"ack_conflict"`
}

// ConfirmAppointmentReq 确认预约请求（patient 端）
type ConfirmAppointmentReq struct{}

// RescheduleReq 改约请求
type RescheduleReq struct {
	NewSlotID string `json:"new_slot_id" binding:"required"`
	Reason    string `json:"reason"      binding:"max=200"`
}

// CancelReq 取消预约请求
type CancelReq struct {
	Reason string `json:"reason" binding:"required,max=200"`
}

// AppealReq 提交申诉请求
type AppealReq struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

// ReviewAppealReq 审核申诉请求
type ReviewAppealReq struct {
	Approved bool   `json:"approved"`
	Comment  string `json:"comment" binding:"max=200"`
}

// ListAppointmentReq 预约列表查询请求
type ListAppointmentReq struct {
	PatientID string `form:"patient_id"`
	Status    string `form:"status"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

// ListBlacklistReq 黑名单列表查询请求
type ListBlacklistReq struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// === 响应 DTO ===

// AppointmentResp 预约单响应
type AppointmentResp struct {
	ID              string                `json:"id"`
	PatientID       string                `json:"patient_id"`
	Mode            string                `json:"mode"`
	Status          string                `json:"status"`
	PaymentVerified bool                  `json:"payment_verified"`
	ChangeCount     int                   `json:"change_count"`
	Items           []AppointmentItemResp `json:"items"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	ConfirmedAt     *time.Time            `json:"confirmed_at,omitempty"`
}

// AppointmentItemResp 预约项目响应
type AppointmentItemResp struct {
	ID         string    `json:"id"`
	ExamItemID string    `json:"exam_item_id"`
	SlotID     string    `json:"slot_id"`
	DeviceID   string    `json:"device_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
}

// AutoAppointmentResp 一键自动预约响应
type AutoAppointmentResp struct {
	AppointmentID string     `json:"appointment_id"`
	Plans         []PlanResp `json:"plans"`
	Warnings      []string   `json:"warnings,omitempty"`
}

// PlanResp 预约方案响应
type PlanResp struct {
	PlanID       string         `json:"plan_id"`
	PlanType     string         `json:"plan_type"`
	Items        []PlanItemResp `json:"items"`
	TotalMinutes int            `json:"total_minutes"`
	WaitMinutes  int            `json:"wait_minutes"`
	TripCount    int            `json:"trip_count"`
}

// PlanItemResp 方案单项响应
type PlanItemResp struct {
	ExamItemID   string    `json:"exam_item_id"`
	ExamItemName string    `json:"exam_item_name"`
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	RoomName     string    `json:"room_name"`
	SlotID       string    `json:"slot_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsFasting    bool      `json:"is_fasting"`
}

// CredentialResp 预约凭证响应
type CredentialResp struct {
	ID                string    `json:"id"`
	AppointmentID     string    `json:"appointment_id"`
	QRCodeURL         string    `json:"qr_code_url"` // base64 image data URL
	PatientNameMasked string    `json:"patient_name_masked"`
	ExamSummary       string    `json:"exam_summary"`
	NoticeContent     string    `json:"notice_content"`
	GeneratedAt       time.Time `json:"generated_at"`
}

// BlacklistResp 黑名单响应
type BlacklistResp struct {
	ID               string    `json:"id"`
	PatientID        string    `json:"patient_id"`
	TriggerTime      time.Time `json:"trigger_time"`
	ExpiresAt        time.Time `json:"expires_at"`
	Status           string    `json:"status"`
	NoShowCount      int       `json:"no_show_count"`
	HasPendingAppeal bool      `json:"has_pending_appeal"`
}

// AppealResp 申诉响应
type AppealResp struct {
	ID          string     `json:"id"`
	BlacklistID string     `json:"blacklist_id"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	ReviewedBy  string     `json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
