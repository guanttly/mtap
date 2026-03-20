// Package triage 应用层 - 分诊服务数据传输对象
package triage

import "time"

// === 请求 DTO ===

// KioskCheckInReq 自动机扫码签到请求
type KioskCheckInReq struct {
	QRCodeData string `json:"qr_code_data" binding:"required"` // 预约凭证中的二维码
}

// NurseCheckInReq 护士手动签到请求
type NurseCheckInReq struct {
	AppointmentID string `json:"appointment_id" binding:"required"`
	Remark        string `json:"remark"        binding:"max=100"`
}

// NFCCheckInReq NFC读卡签到请求
type NFCCheckInReq struct {
	CardID string `json:"card_id" binding:"required"`
}

// UndoExamReq 撤销误操作请求
type UndoExamReq struct {
	Reason string `json:"reason" binding:"required,max=200"`
}

// === 响应 DTO ===

// CheckInResp 签到响应
type CheckInResp struct {
	CheckInID     string `json:"check_in_id"`
	QueueNumber   int    `json:"queue_number"`
	EstimatedWait int    `json:"estimated_wait"` // 预计等候分钟
	RoomLocation  string `json:"room_location"`
	IsLate        bool   `json:"is_late"`
}

// QueueEntryResp 队列条目响应
type QueueEntryResp struct {
	ID                string     `json:"id"`
	QueueNumber       int        `json:"queue_number"`
	PatientNameMasked string     `json:"patient_name_masked"`
	Status            string     `json:"status"`
	CallCount         int        `json:"call_count"`
	MissCount         int        `json:"miss_count"`
	EnteredAt         time.Time  `json:"entered_at"`
	CalledAt          *time.Time `json:"called_at,omitempty"`
}

// QueueStatusResp 队列状态响应
type QueueStatusResp struct {
	RoomID         string           `json:"room_id"`
	WaitingCount   int              `json:"waiting_count"`
	AverageWait    int              `json:"average_wait"`
	CurrentCalling *QueueEntryResp  `json:"current_calling,omitempty"`
	Entries        []QueueEntryResp `json:"entries"`
}

// ExamExecutionResp 检查执行响应
type ExamExecutionResp struct {
	ID                string     `json:"id"`
	AppointmentItemID string     `json:"appointment_item_id"`
	PatientID         string     `json:"patient_id"`
	DeviceID          string     `json:"device_id"`
	Status            string     `json:"status"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	Duration          int        `json:"duration"`
}
