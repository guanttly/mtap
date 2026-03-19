package resource

import "time"

type CreateDeviceReq struct {
	Name         string `json:"name" binding:"required,max=100"`
	CampusID     string `json:"campus_id"`
	DepartmentID string `json:"department_id"`
}

type DeviceResp struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	CampusID     string    `json:"campus_id"`
	DepartmentID string    `json:"department_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateExamItemReq struct {
	Name        string `json:"name" binding:"required,max=100"`
	DurationMin int    `json:"duration_min" binding:"required,gt=0"`
	IsFasting   bool   `json:"is_fasting"`
	FastingDesc string `json:"fasting_desc" binding:"omitempty,max=200"`
}

type ExamItemResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DurationMin int    `json:"duration_min"`
	IsFasting   bool   `json:"is_fasting"`
	FastingDesc string `json:"fasting_desc,omitempty"`
}

type CreateAliasReq struct {
	ExamItemID string `json:"exam_item_id" binding:"required"`
	Alias      string `json:"alias" binding:"required,max=50"`
}

type AliasResp struct {
	ID        string `json:"id"`
	ExamItemID string `json:"exam_item_id"`
	Alias     string `json:"alias"`
}

type CreateSlotPoolReq struct {
	Name string `json:"name" binding:"required,max=60"`
	Type string `json:"type" binding:"required,oneof=public department doctor"`
}

type SlotPoolResp struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type GenerateScheduleReq struct {
	DeviceID    string `json:"device_id" binding:"required"`
	Date        string `json:"date" binding:"omitempty"`       // YYYY-MM-DD（单日生成）
	StartDate   string `json:"start_date" binding:"omitempty"` // YYYY-MM-DD（批量生成起）
	EndDate     string `json:"end_date" binding:"omitempty"`   // YYYY-MM-DD（批量生成止）
	StartTime   string `json:"start_time" binding:"required"` // HH:mm
	EndTime     string `json:"end_time" binding:"required"`   // HH:mm
	SlotMinutes int    `json:"slot_minutes" binding:"required,gt=0"`
	ExamItemID  string `json:"exam_item_id"` // 可选：绑定检查项目
	PoolType    string `json:"pool_type" binding:"omitempty,oneof=public department doctor outpatient inpatient referral"`
	SkipWeekends bool  `json:"skip_weekends"`
}

type TimeSlotResp struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	ExamItemID string   `json:"exam_item_id,omitempty"`
	PoolType  string    `json:"pool_type"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	Status    string    `json:"status"`
	StandardDuration int `json:"standard_duration"`
	AdjustedDuration int `json:"adjusted_duration"`
	LockedBy  string    `json:"locked_by,omitempty"`
	LockUntil *time.Time `json:"lock_until,omitempty"`
	Remaining int       `json:"remaining"`
}

type LockSlotReq struct {
	PatientID string `json:"patient_id" binding:"required"`
}

type QueryAvailableSlotsReq struct {
	DeviceID   string `form:"device_id" binding:"required"`
	Date       string `form:"date" binding:"required"` // YYYY-MM-DD
	ExamItemID string `form:"exam_item_id"`
	PoolType   string `form:"pool_type"`
	PatientAge int    `form:"patient_age"`
}

type SuspendScheduleReq struct {
	DeviceID   string `json:"device_id" binding:"required"`
	Date       string `json:"date" binding:"required"` // YYYY-MM-DD
	StartTime  string `json:"start_time" binding:"required"` // HH:mm
	EndTime    string `json:"end_time" binding:"required"`   // HH:mm
	Reason     string `json:"reason" binding:"required,max=200"`
}

type SubstituteScheduleReq struct {
	SourceDeviceID string `json:"source_device_id" binding:"required"`
	TargetDeviceID string `json:"target_device_id" binding:"required"`
	Date           string `json:"date" binding:"required"` // YYYY-MM-DD
}

type AddExtraSlotsReq struct {
	DeviceID    string `json:"device_id" binding:"required"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	StartTime   string `json:"start_time" binding:"required"` // HH:mm
	EndTime     string `json:"end_time" binding:"required"`   // HH:mm
	SlotMinutes int    `json:"slot_minutes" binding:"required,gt=0"`
	Reason      string `json:"reason" binding:"required,max=200"`
	ExamItemID  string `json:"exam_item_id"`
	PoolType    string `json:"pool_type" binding:"omitempty,oneof=public department doctor outpatient inpatient referral"`
}
