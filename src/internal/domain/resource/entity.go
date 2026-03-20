// Package resource 资源管理领域 - 实体定义
// 核心目的：定义基础信息与资源排班管理的核心实体
// 模块功能：
//   - Device: 设备聚合根（归属科室/院区、支持检查类型列表）
//   - Schedule: 排班计划聚合根（设备 + 日期 + 工作时段 + 排班模板）
//   - TimeSlot: 号源时段实体（属于Schedule，起止时间 + 关联检查类型 + 剩余量）
//   - SlotPool: 号源池聚合根（公共池/科室池/医生专池，含配额与溢出规则）
//   - ExamItem: 检查项目聚合根（标准耗时、空腹标记）
//   - Campus: 院区实体
//   - Department: 科室实体
//   - Doctor: 医生实体
package resource

import (
	"time"

	"github.com/google/uuid"
)

// ── 院区 ─────────────────────────────────────────────────────────

// Campus 院区实体
type Campus struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	Status    string    `json:"status"` // active / inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCampus 创建院区
func NewCampus(name, code, address string) *Campus {
	now := time.Now()
	return &Campus{
		ID:        uuid.New().String(),
		Name:      name,
		Code:      code,
		Address:   address,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ── 科室 ─────────────────────────────────────────────────────────

// Department 科室实体
type Department struct {
	ID        string    `json:"id"`
	CampusID  string    `json:"campus_id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Floor     string    `json:"floor"`
	Status    string    `json:"status"` // active / inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewDepartment 创建科室
func NewDepartment(campusID, name, code, floor string) *Department {
	now := time.Now()
	return &Department{
		ID:        uuid.New().String(),
		CampusID:  campusID,
		Name:      name,
		Code:      code,
		Floor:     floor,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ── 医生 ─────────────────────────────────────────────────────────

// Doctor 医生实体
type Doctor struct {
	ID           string    `json:"id"`
	DepartmentID string    `json:"department_id"`
	HISCode      string    `json:"his_code"` // HIS 系统编码
	Name         string    `json:"name"`
	Title        string    `json:"title"`  // 职称
	Gender       string    `json:"gender"` // M / F
	Status       string    `json:"status"` // active / inactive
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewDoctor 创建医生
func NewDoctor(deptID, hisCode, name, title, gender string) *Doctor {
	now := time.Now()
	return &Doctor{
		ID:           uuid.New().String(),
		DepartmentID: deptID,
		HISCode:      hisCode,
		Name:         name,
		Title:        title,
		Gender:       gender,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// ── 设备 ─────────────────────────────────────────────────────────

// DeviceStatus 设备状态
type DeviceStatus string

const (
	DeviceStatusActive      DeviceStatus = "active"
	DeviceStatusMaintenance DeviceStatus = "maintenance"
	DeviceStatusRetired     DeviceStatus = "retired"
)

// Device 设备聚合根
type Device struct {
	ID                 string       `json:"id"`
	DepartmentID       string       `json:"department_id"`
	CampusID           string       `json:"campus_id"`
	Name               string       `json:"name"`
	Model              string       `json:"model"`
	Manufacturer       string       `json:"manufacturer"`
	SupportedExamTypes []string     `json:"supported_exam_types"` // e.g. ["CT", "US"]
	MaxDailySlots      int          `json:"max_daily_slots"`
	Status             DeviceStatus `json:"status"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

// NewDevice 创建设备
func NewDevice(deptID, campusID, name, model, manufacturer string, maxSlots int, examTypes []string) *Device {
	now := time.Now()
	return &Device{
		ID:                 uuid.New().String(),
		DepartmentID:       deptID,
		CampusID:           campusID,
		Name:               name,
		Model:              model,
		Manufacturer:       manufacturer,
		SupportedExamTypes: examTypes,
		MaxDailySlots:      maxSlots,
		Status:             DeviceStatusActive,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// ── 检查项目 ─────────────────────────────────────────────────────

// ItemAlias 检查项目别名
type ItemAlias struct {
	ID         string    `json:"id"`
	ExamItemID string    `json:"exam_item_id"`
	Alias      string    `json:"alias"`
	CreatedAt  time.Time `json:"created_at"`
}

// ExamItem 检查项目聚合根
type ExamItem struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DurationMin int         `json:"duration_min"` // 标准耗时（分钟）
	IsFasting   bool        `json:"is_fasting"`
	FastingDesc string      `json:"fasting_desc"` // 空腹说明，最长200字符
	Aliases     []ItemAlias `json:"aliases"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewExamItem 创建检查项目
func NewExamItem(name string, durationMin int, isFasting bool, fastingDesc string) *ExamItem {
	now := time.Now()
	return &ExamItem{
		ID:          uuid.New().String(),
		Name:        name,
		DurationMin: durationMin,
		IsFasting:   isFasting,
		FastingDesc: fastingDesc,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddAlias 添加别名
func (e *ExamItem) AddAlias(alias string) {
	e.Aliases = append(e.Aliases, ItemAlias{
		ID:         uuid.New().String(),
		ExamItemID: e.ID,
		Alias:      alias,
		CreatedAt:  time.Now(),
	})
}

// MatchName 检查名称或别名是否匹配
func (e *ExamItem) MatchName(input string) bool {
	if e.Name == input {
		return true
	}
	for _, a := range e.Aliases {
		if a.Alias == input {
			return true
		}
	}
	return false
}

// ── 号源池 ────────────────────────────────────────────────────────

// SlotPoolType 号源池类型
type SlotPoolType string

const (
	SlotPoolPublic     SlotPoolType = "public"     // 公共池
	SlotPoolDepartment SlotPoolType = "department" // 科室池
	SlotPoolDoctor     SlotPoolType = "doctor"     // 医生专池
)

// SlotPool 号源池聚合根
type SlotPool struct {
	ID                 string       `json:"id"`
	Name               string       `json:"name"`
	Type               SlotPoolType `json:"type"`
	AllocationRatio    float64      `json:"allocation_ratio"` // 分配比例 0~1
	OverflowEnabled    bool         `json:"overflow_enabled"`
	OverflowTargetPool string       `json:"overflow_target_pool"` // 溢出目标池ID
	Status             string       `json:"status"`               // active / inactive
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

// NewSlotPool 创建号源池
func NewSlotPool(name string, poolType SlotPoolType, ratio float64) *SlotPool {
	now := time.Now()
	return &SlotPool{
		ID:              uuid.New().String(),
		Name:            name,
		Type:            poolType,
		AllocationRatio: ratio,
		Status:          "active",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// ── 排班 ─────────────────────────────────────────────────────────

// ScheduleStatus 排班状态
type ScheduleStatus string

const (
	ScheduleStatusNormal     ScheduleStatus = "normal"
	ScheduleStatusSuspended  ScheduleStatus = "suspended"
	ScheduleStatusSubstitute ScheduleStatus = "substitute"
)

// Schedule 排班计划聚合根
type Schedule struct {
	ID            string         `json:"id"`
	DeviceID      string         `json:"device_id"`
	WorkDate      time.Time      `json:"work_date"`
	StartTime     string         `json:"start_time"` // HH:mm
	EndTime       string         `json:"end_time"`   // HH:mm
	Status        ScheduleStatus `json:"status"`
	SuspendReason string         `json:"suspend_reason,omitempty"`
	Slots         []TimeSlot     `json:"slots,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// NewSchedule 创建排班
func NewSchedule(deviceID string, date time.Time, startTime, endTime string) *Schedule {
	now := time.Now()
	return &Schedule{
		ID:        uuid.New().String(),
		DeviceID:  deviceID,
		WorkDate:  date,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    ScheduleStatusNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Suspend 停诊
func (s *Schedule) Suspend(reason string) {
	s.Status = ScheduleStatusSuspended
	s.SuspendReason = reason
	s.UpdatedAt = time.Now()
}

// SubstituteTo 替代排班（改到另一台设备）
func (s *Schedule) SubstituteTo(targetDeviceID string) {
	s.DeviceID = targetDeviceID
	s.Status = ScheduleStatusSubstitute
	s.UpdatedAt = time.Now()
}

// ── 号源时段 ─────────────────────────────────────────────────────

// TimeSlotStatus 号源状态
type TimeSlotStatus string

const (
	TimeSlotAvailable TimeSlotStatus = "available"
	TimeSlotLocked    TimeSlotStatus = "locked"
	TimeSlotBooked    TimeSlotStatus = "booked"
	TimeSlotCompleted TimeSlotStatus = "completed"
	TimeSlotSuspended TimeSlotStatus = "suspended"
	TimeSlotExpired   TimeSlotStatus = "expired"
)

// TimeSlot 号源时段实体
type TimeSlot struct {
	ID               string         `json:"id"`
	ScheduleID       string         `json:"schedule_id"`
	DeviceID         string         `json:"device_id"`
	Date             time.Time      `json:"date"`
	ExamItemID       string         `json:"exam_item_id"`
	PoolType         string         `json:"pool_type"` // public / department / doctor
	StartAt          time.Time      `json:"start_at"`
	EndAt            time.Time      `json:"end_at"`
	StandardDuration int            `json:"standard_duration"` // 分钟
	AdjustedDuration int            `json:"adjusted_duration"` // 调整后耗时（分钟）
	Status           TimeSlotStatus `json:"status"`
	LockedBy         string         `json:"locked_by,omitempty"` // 锁定的患者ID
	LockUntil        *time.Time     `json:"lock_until,omitempty"`
	Remaining        int            `json:"remaining"` // 剩余容量（通常为1）
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// NewTimeSlot 创建号源时段
func NewTimeSlot(scheduleID, deviceID, examItemID, poolType string, date, startAt, endAt time.Time, stdDuration int) *TimeSlot {
	now := time.Now()
	return &TimeSlot{
		ID:               uuid.New().String(),
		ScheduleID:       scheduleID,
		DeviceID:         deviceID,
		Date:             date,
		ExamItemID:       examItemID,
		PoolType:         poolType,
		StartAt:          startAt,
		EndAt:            endAt,
		StandardDuration: stdDuration,
		AdjustedDuration: stdDuration,
		Status:           TimeSlotAvailable,
		Remaining:        1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Lock 锁定号源
func (t *TimeSlot) Lock(patientID string, until time.Time) bool {
	if t.Status != TimeSlotAvailable || t.Remaining <= 0 {
		return false
	}
	t.Status = TimeSlotLocked
	t.LockedBy = patientID
	t.LockUntil = &until
	t.Remaining--
	t.UpdatedAt = time.Now()
	return true
}

// Release 释放号源
func (t *TimeSlot) Release() {
	t.Status = TimeSlotAvailable
	t.LockedBy = ""
	t.LockUntil = nil
	t.Remaining++
	t.UpdatedAt = time.Now()
}

// Book 预订（锁定 → 已预约）
func (t *TimeSlot) Book() {
	t.Status = TimeSlotBooked
	t.UpdatedAt = time.Now()
}

// IsExpiredLock 判断锁定是否已超时
func (t *TimeSlot) IsExpiredLock() bool {
	return t.Status == TimeSlotLocked && t.LockUntil != nil && time.Now().After(*t.LockUntil)
}
