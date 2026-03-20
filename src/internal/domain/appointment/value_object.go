// Package appointment 预约服务领域 - 值对象
package appointment

import "time"

// AppointmentMode 预约模式
type AppointmentMode string

const (
	ModeAuto   AppointmentMode = "auto"   // 一键自动预约
	ModeCombo  AppointmentMode = "combo"  // 组合预约
	ModeManual AppointmentMode = "manual" // 人工干预
)

func (m AppointmentMode) IsValid() bool {
	switch m {
	case ModeAuto, ModeCombo, ModeManual:
		return true
	}
	return false
}

// AppointmentStatus 预约状态
type AppointmentStatus string

const (
	StatusPending      AppointmentStatus = "pending"       // 待确认
	StatusConfirmed    AppointmentStatus = "confirmed"     // 已确认
	StatusPayVerifying AppointmentStatus = "pay_verifying" // 缴费待校验
	StatusPaid         AppointmentStatus = "paid"          // 已缴费
	StatusRescheduling AppointmentStatus = "rescheduling"  // 改约中
	StatusCancelled    AppointmentStatus = "cancelled"     // 已取消
	StatusCheckedIn    AppointmentStatus = "checked_in"    // 已签到
	StatusExamining    AppointmentStatus = "examining"     // 检查中
	StatusCompleted    AppointmentStatus = "completed"     // 已完成
	StatusNoShow       AppointmentStatus = "no_show"       // 爽约
	StatusReleased     AppointmentStatus = "released"      // 已释放
)

// ItemStatus 预约项目状态
type ItemStatus string

const (
	ItemStatusPending   ItemStatus = "pending"    // 待检查
	ItemStatusCheckedIn ItemStatus = "checked_in" // 已签到
	ItemStatusExamining ItemStatus = "examining"  // 检查中
	ItemStatusCompleted ItemStatus = "completed"  // 已完成
	ItemStatusCancelled ItemStatus = "cancelled"  // 已取消
)

// BlacklistStatus 黑名单状态
type BlacklistStatus string

const (
	BlacklistActive   BlacklistStatus = "active"   // 生效中
	BlacklistReleased BlacklistStatus = "released" // 已解除
	BlacklistExpired  BlacklistStatus = "expired"  // 已过期
)

// AppealStatus 申诉状态
type AppealStatus string

const (
	AppealPending  AppealStatus = "pending"  // 待审核
	AppealApproved AppealStatus = "approved" // 已通过
	AppealRejected AppealStatus = "rejected" // 已驳回
)

// AppointmentPlan 预约方案（值对象）
type AppointmentPlan struct {
	PlanID       string     `json:"plan_id"`
	PlanType     string     `json:"plan_type"` // shortest_time / least_trips
	Items        []PlanItem `json:"items"`
	TotalMinutes int        `json:"total_minutes"`
	WaitMinutes  int        `json:"wait_minutes"`
	TripCount    int        `json:"trip_count"`
}

// PlanItem 方案中的单个检查安排
type PlanItem struct {
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

// ChangeLog 改约/取消变更记录（值对象）
type ChangeLog struct {
	ID            string    `json:"id"`
	AppointmentID string    `json:"appointment_id"`
	ChangeType    string    `json:"change_type"` // reschedule / cancel
	OldSlotID     string    `json:"old_slot_id"`
	NewSlotID     string    `json:"new_slot_id"`
	Reason        string    `json:"reason"`
	OperatorID    string    `json:"operator_id"`
	ChangedAt     time.Time `json:"changed_at"`
}
