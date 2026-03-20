// Package appointment 预约服务领域 - 实体定义
package appointment

import (
	"time"

	"github.com/google/uuid"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// Appointment 预约单聚合根
type Appointment struct {
	ID              string            `json:"id"`
	PatientID       string            `json:"patient_id"`
	Mode            AppointmentMode   `json:"mode"`
	Status          AppointmentStatus `json:"status"`
	Items           []AppointmentItem `json:"items"`
	OverrideBy      string            `json:"override_by"`
	OverrideReason  string            `json:"override_reason"`
	PaymentVerified bool              `json:"payment_verified"`
	ChangeCount     int               `json:"change_count"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ConfirmedAt     *time.Time        `json:"confirmed_at"`
	CancelledAt     *time.Time        `json:"cancelled_at"`
	CancelReason    string            `json:"cancel_reason"`
}

// AppointmentItem 预约单中的单个检查项目
type AppointmentItem struct {
	ID            string     `json:"id"`
	AppointmentID string     `json:"appointment_id"`
	ExamItemID    string     `json:"exam_item_id"`
	SlotID        string     `json:"slot_id"`
	DeviceID      string     `json:"device_id"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	Status        ItemStatus `json:"status"`
}

// NewAppointment 创建预约单
func NewAppointment(patientID string, mode AppointmentMode) (*Appointment, error) {
	if patientID == "" {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "患者ID不能为空")
	}
	if !mode.IsValid() {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "无效的预约模式")
	}
	now := time.Now()
	return &Appointment{
		ID:        uuid.New().String(),
		PatientID: patientID,
		Mode:      mode,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// AddItem 添加预约项目
func (a *Appointment) AddItem(examItemID, slotID, deviceID string, start, end time.Time) {
	a.Items = append(a.Items, AppointmentItem{
		ID:            uuid.New().String(),
		AppointmentID: a.ID,
		ExamItemID:    examItemID,
		SlotID:        slotID,
		DeviceID:      deviceID,
		StartTime:     start,
		EndTime:       end,
		Status:        ItemStatusPending,
	})
}

// Confirm 确认预约
func (a *Appointment) Confirm() error {
	if a.Status != StatusPending {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "仅待确认状态的预约可以确认")
	}
	now := time.Now()
	a.Status = StatusConfirmed
	a.ConfirmedAt = &now
	a.UpdatedAt = now
	return nil
}

// MarkPaid 标记已缴费
func (a *Appointment) MarkPaid() error {
	if a.Status != StatusConfirmed && a.Status != StatusPayVerifying {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "状态不允许标记缴费")
	}
	a.PaymentVerified = true
	a.Status = StatusPaid
	a.UpdatedAt = time.Now()
	return nil
}

// Cancel 取消预约
func (a *Appointment) Cancel(reason string) error {
	switch a.Status {
	case StatusPaid, StatusConfirmed, StatusPayVerifying:
		// 检查是否距检查不足2小时
		for _, item := range a.Items {
			if time.Until(item.StartTime) < 2*time.Hour {
				return bizErr.New(bizErr.ErrApptTooCloseToExam)
			}
		}
	default:
		return bizErr.NewWithDetail(bizErr.ErrConflict, "当前状态不允许取消")
	}
	now := time.Now()
	a.Status = StatusCancelled
	a.CancelledAt = &now
	a.CancelReason = reason
	a.UpdatedAt = now
	return nil
}

// Reschedule 改约
func (a *Appointment) Reschedule() error {
	if a.Status != StatusPaid {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "仅已缴费状态的预约可以改约")
	}
	if a.ChangeCount >= 3 {
		return bizErr.New(bizErr.ErrApptChangeLimitReached)
	}
	// 检查是否距检查不足2小时
	for _, item := range a.Items {
		if time.Until(item.StartTime) < 2*time.Hour {
			return bizErr.New(bizErr.ErrApptTooCloseToExam)
		}
	}
	a.Status = StatusRescheduling
	a.UpdatedAt = time.Now()
	return nil
}

// CompleteReschedule 改约完成
func (a *Appointment) CompleteReschedule() {
	a.ChangeCount++
	a.Status = StatusPaid
	a.UpdatedAt = time.Now()
}

// MarkNoShow 标记爽约
func (a *Appointment) MarkNoShow() error {
	if a.Status != StatusCheckedIn && a.Status != StatusPaid {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "状态不允许标记爽约")
	}
	a.Status = StatusNoShow
	a.UpdatedAt = time.Now()
	return nil
}

// MarkCheckedIn 标记已签到
func (a *Appointment) MarkCheckedIn() error {
	if a.Status != StatusPaid {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "仅已缴费状态可以签到")
	}
	a.Status = StatusCheckedIn
	a.UpdatedAt = time.Now()
	return nil
}

// Release 释放号源后标记已释放
func (a *Appointment) Release() {
	a.Status = StatusReleased
	a.UpdatedAt = time.Now()
}

// Credential 预约凭证实体
type Credential struct {
	ID                string    `json:"id"`
	AppointmentID     string    `json:"appointment_id"`
	QRCodeData        string    `json:"qr_code_data"`
	PatientNameMasked string    `json:"patient_name_masked"`
	ExamSummary       string    `json:"exam_summary"`
	NoticeContent     string    `json:"notice_content"`
	GeneratedAt       time.Time `json:"generated_at"`
}

// NewCredential 创建预约凭证
func NewCredential(appointmentID, qrCodeData, patientNameMasked, examSummary, noticeContent string) *Credential {
	return &Credential{
		ID:                uuid.New().String(),
		AppointmentID:     appointmentID,
		QRCodeData:        qrCodeData,
		PatientNameMasked: patientNameMasked,
		ExamSummary:       examSummary,
		NoticeContent:     noticeContent,
		GeneratedAt:       time.Now(),
	}
}

// Blacklist 黑名单聚合根
type Blacklist struct {
	ID            string          `json:"id"`
	PatientID     string          `json:"patient_id"`
	TriggerTime   time.Time       `json:"trigger_time"`
	ExpiresAt     time.Time       `json:"expires_at"`
	Status        BlacklistStatus `json:"status"`
	ReleasedAt    *time.Time      `json:"released_at"`
	ReleaseReason string          `json:"release_reason"`
	NoShowRecords []NoShowRecord  `json:"no_show_records"`
	Appeal        *Appeal         `json:"appeal"`
}

// NewBlacklist 创建黑名单记录
func NewBlacklist(patientID string, validDays int) *Blacklist {
	now := time.Now()
	return &Blacklist{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		TriggerTime: now,
		ExpiresAt:   now.AddDate(0, 0, validDays),
		Status:      BlacklistActive,
	}
}

// IsExpired 是否已过期
func (b *Blacklist) IsExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// CanAppointOnline 是否可以在线预约
func (b *Blacklist) CanAppointOnline() bool {
	return b.Status != BlacklistActive || b.IsExpired()
}

// Release 解除黑名单
func (b *Blacklist) Release(reason string) error {
	if b.Status != BlacklistActive {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "黑名单已解除或过期")
	}
	now := time.Now()
	b.Status = BlacklistReleased
	b.ReleasedAt = &now
	b.ReleaseReason = reason
	return nil
}

// NoShowRecord 爽约记录实体
type NoShowRecord struct {
	ID            string    `json:"id"`
	PatientID     string    `json:"patient_id"`
	AppointmentID string    `json:"appointment_id"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// NewNoShowRecord 创建爽约记录
func NewNoShowRecord(patientID, appointmentID string) *NoShowRecord {
	return &NoShowRecord{
		ID:            uuid.New().String(),
		PatientID:     patientID,
		AppointmentID: appointmentID,
		OccurredAt:    time.Now(),
	}
}

// Appeal 申诉实体
type Appeal struct {
	ID          string       `json:"id"`
	BlacklistID string       `json:"blacklist_id"`
	Reason      string       `json:"reason"`
	Status      AppealStatus `json:"status"`
	ReviewedBy  string       `json:"reviewed_by"`
	ReviewedAt  *time.Time   `json:"reviewed_at"`
	CreatedAt   time.Time    `json:"created_at"`
}

// NewAppeal 创建申诉
func NewAppeal(blacklistID, reason string) (*Appeal, error) {
	if len(reason) == 0 || len(reason) > 500 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "申诉原因长度应在 1~500 字符之间")
	}
	return &Appeal{
		ID:          uuid.New().String(),
		BlacklistID: blacklistID,
		Reason:      reason,
		Status:      AppealPending,
		CreatedAt:   time.Now(),
	}, nil
}

// Review 审核申诉
func (a *Appeal) Review(reviewerID string, approved bool) error {
	if a.Status != AppealPending {
		return bizErr.NewWithDetail(bizErr.ErrConflict, "申诉已审核")
	}
	now := time.Now()
	a.ReviewedBy = reviewerID
	a.ReviewedAt = &now
	if approved {
		a.Status = AppealApproved
	} else {
		a.Status = AppealRejected
	}
	return nil
}
