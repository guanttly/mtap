package resource

import "time"

type CampusResp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address,omitempty"`
	Status  string `json:"status"`
}

type DepartmentResp struct {
	ID       string `json:"id"`
	CampusID string `json:"campus_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Floor    string `json:"floor,omitempty"`
	Status   string `json:"status"`
}

type ScheduleResp struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

type ListSchedulesReq struct {
	DeviceID  string `form:"device_id"`
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
}

type CreateDeviceReq struct {
	Name               string   `json:"name" binding:"required,max=100"`
	CampusID           string   `json:"campus_id"`
	DepartmentID       string   `json:"department_id"`
	Model              string   `json:"model" binding:"omitempty,max=50"`
	Manufacturer       string   `json:"manufacturer" binding:"omitempty,max=50"`
	SupportedExamTypes []string `json:"supported_exam_types"`
	MaxDailySlots      int      `json:"max_daily_slots" binding:"omitempty,gte=1"`
}

type DeviceResp struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	CampusID           string    `json:"campus_id"`
	DepartmentID       string    `json:"department_id"`
	Model              string    `json:"model,omitempty"`
	Manufacturer       string    `json:"manufacturer,omitempty"`
	SupportedExamTypes []string  `json:"supported_exam_types"`
	MaxDailySlots      int       `json:"max_daily_slots"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
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
	ID         string `json:"id"`
	ExamItemID string `json:"exam_item_id"`
	Alias      string `json:"alias"`
}

type CreateSlotPoolReq struct {
	Name               string  `json:"name" binding:"required,max=60"`
	Type               string  `json:"type" binding:"required,oneof=public department doctor"`
	AllocationRatio    float64 `json:"allocation_ratio" binding:"omitempty,gte=0,lte=1"`
	OverflowEnabled    bool    `json:"overflow_enabled"`
	OverflowTargetPool string  `json:"overflow_target_pool"`
}

type SlotPoolResp struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	Status             string  `json:"status"`
	AllocationRatio    float64 `json:"allocation_ratio"`
	OverflowEnabled    bool    `json:"overflow_enabled"`
	OverflowTargetPool string  `json:"overflow_target_pool,omitempty"`
}

type GenerateScheduleReq struct {
	DeviceID     string `json:"device_id" binding:"required"`
	TemplateID   string `json:"template_id"`                    // 可选：从模板加载排班参数
	Date         string `json:"date" binding:"omitempty"`       // YYYY-MM-DD（单日生成）
	StartDate    string `json:"start_date" binding:"omitempty"` // YYYY-MM-DD（批量生成起）
	EndDate      string `json:"end_date" binding:"omitempty"`   // YYYY-MM-DD（批量生成止）
	StartTime    string `json:"start_time" binding:"required"`  // HH:mm
	EndTime      string `json:"end_time" binding:"required"`    // HH:mm
	SlotMinutes  int    `json:"slot_minutes" binding:"required,gt=0"`
	ExamItemID   string `json:"exam_item_id"` // 可选：绑定检查项目
	PoolType     string `json:"pool_type" binding:"omitempty,oneof=public department doctor outpatient inpatient referral"`
	SkipWeekends bool   `json:"skip_weekends"`
}

type TimeSlotResp struct {
	ID               string     `json:"id"`
	DeviceID         string     `json:"device_id"`
	ExamItemID       string     `json:"exam_item_id,omitempty"`
	PoolType         string     `json:"pool_type"`
	StartAt          time.Time  `json:"start_at"`
	EndAt            time.Time  `json:"end_at"`
	Status           string     `json:"status"`
	StandardDuration int        `json:"standard_duration"`
	AdjustedDuration int        `json:"adjusted_duration"`
	LockedBy         string     `json:"locked_by,omitempty"`
	LockUntil        *time.Time `json:"lock_until,omitempty"`
	Remaining        int        `json:"remaining"`
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
	DeviceID  string `json:"device_id" binding:"required"`
	Date      string `json:"date" binding:"required"`       // YYYY-MM-DD
	StartTime string `json:"start_time" binding:"required"` // HH:mm
	EndTime   string `json:"end_time" binding:"required"`   // HH:mm
	Reason    string `json:"reason" binding:"required,max=200"`
}

type SubstituteScheduleReq struct {
	SourceDeviceID string `json:"source_device_id" binding:"required"`
	TargetDeviceID string `json:"target_device_id" binding:"required"`
	Date           string `json:"date" binding:"required"` // YYYY-MM-DD
}

type AddExtraSlotsReq struct {
	DeviceID    string `json:"device_id" binding:"required"`
	Date        string `json:"date" binding:"required"`       // YYYY-MM-DD
	StartTime   string `json:"start_time" binding:"required"` // HH:mm
	EndTime     string `json:"end_time" binding:"required"`   // HH:mm
	SlotMinutes int    `json:"slot_minutes" binding:"required,gt=0"`
	Reason      string `json:"reason" binding:"required,max=200"`
	ExamItemID  string `json:"exam_item_id"`
	PoolType    string `json:"pool_type" binding:"omitempty,oneof=public department doctor outpatient inpatient referral"`
}

// UpdateDeviceReq 更新设备请求
type UpdateDeviceReq struct {
	Name               string   `json:"name" binding:"omitempty,max=100"`
	CampusID           string   `json:"campus_id"`
	DepartmentID       string   `json:"department_id"`
	Model              string   `json:"model" binding:"omitempty,max=50"`
	Manufacturer       string   `json:"manufacturer" binding:"omitempty,max=50"`
	SupportedExamTypes []string `json:"supported_exam_types"`
	MaxDailySlots      int      `json:"max_daily_slots" binding:"omitempty,gte=1"`
	Status             string   `json:"status" binding:"omitempty,oneof=active maintenance retired"`
}

// DoctorResp 医生响应
type DoctorResp struct {
	ID           string `json:"id"`
	DepartmentID string `json:"department_id"`
	HISCode      string `json:"his_code,omitempty"`
	Name         string `json:"name"`
	Title        string `json:"title,omitempty"`
	Gender       string `json:"gender"`
	Status       string `json:"status"`
}

// CreateDoctorReq 创建医生请求
type CreateDoctorReq struct {
	DepartmentID string `json:"department_id" binding:"required"`
	HISCode      string `json:"his_code" binding:"omitempty,max=30"`
	Name         string `json:"name" binding:"required,max=30"`
	Title        string `json:"title" binding:"omitempty,max=20"`
	Gender       string `json:"gender" binding:"omitempty,oneof=M F unknown"`
}

// UpdateDoctorReq 更新医生请求
type UpdateDoctorReq struct {
	HISCode string `json:"his_code" binding:"omitempty,max=30"`
	Name    string `json:"name" binding:"omitempty,max=30"`
	Title   string `json:"title" binding:"omitempty,max=20"`
	Gender  string `json:"gender" binding:"omitempty,oneof=M F unknown"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// SlotPatternReq 排班模板号源配置
type SlotPatternReq struct {
	StartTime   string `json:"start_time" binding:"required"` // HH:mm
	EndTime     string `json:"end_time" binding:"required"`   // HH:mm
	SlotMinutes int    `json:"slot_minutes" binding:"required,gt=0"`
	ExamItemID  string `json:"exam_item_id"`
	PoolType    string `json:"pool_type" binding:"omitempty,oneof=public department doctor outpatient inpatient referral"`
}

// ScheduleTemplateResp 排班模板响应
type ScheduleTemplateResp struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	RepeatType   string         `json:"repeat_type"`
	SlotPattern  SlotPatternReq `json:"slot_pattern"`
	SkipWeekends bool           `json:"skip_weekends"`
}

// CreateScheduleTemplateReq 创建排班模板请求
type CreateScheduleTemplateReq struct {
	Name         string         `json:"name" binding:"required,max=50"`
	RepeatType   string         `json:"repeat_type" binding:"required,oneof=once daily weekly"`
	SlotPattern  SlotPatternReq `json:"slot_pattern" binding:"required"`
	SkipWeekends bool           `json:"skip_weekends"`
}

// UpdateExamItemReq 更新检查项目请求
type UpdateExamItemReq struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	DurationMin int    `json:"duration_min" binding:"omitempty,gt=0"`
	IsFasting   *bool  `json:"is_fasting"`
	FastingDesc string `json:"fasting_desc" binding:"omitempty,max=200"`
}
