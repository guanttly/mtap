// Package po 基础设施层 - appointment 持久化对象（GORM Model）
package po

import "time"

// AppointmentPO 预约单表
type AppointmentPO struct {
	ID              string     `gorm:"column:id;primaryKey;size:36"`
	PatientID       string     `gorm:"column:patient_id;not null;size:36;index:idx_appt_patient"`
	Mode            string     `gorm:"column:mode;not null;size:10"`
	Status          string     `gorm:"column:status;not null;size:15;index:idx_appt_status"`
	OverrideBy      string     `gorm:"column:override_by;size:36"`
	OverrideReason  string     `gorm:"column:override_reason;size:200"`
	PaymentVerified bool       `gorm:"column:payment_verified;not null;default:false"`
	ChangeCount     int        `gorm:"column:change_count;not null;default:0"`
	CancelReason    string     `gorm:"column:cancel_reason;size:200"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	CancelledAt     *time.Time `gorm:"column:cancelled_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime"`

	Items []AppointmentItemPO `gorm:"foreignKey:AppointmentID;references:ID"`
}

func (AppointmentPO) TableName() string { return "appointments" }

// AppointmentItemPO 预约单检查项目表
type AppointmentItemPO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	AppointmentID string    `gorm:"column:appointment_id;not null;size:36;index:idx_appt_items_appt"`
	ExamItemID    string    `gorm:"column:exam_item_id;not null;size:36"`
	SlotID        string    `gorm:"column:slot_id;not null;size:36"`
	DeviceID      string    `gorm:"column:device_id;not null;size:36"`
	StartTime     time.Time `gorm:"column:start_time;not null"`
	EndTime       time.Time `gorm:"column:end_time;not null"`
	Status        string    `gorm:"column:status;not null;size:15;default:pending"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AppointmentItemPO) TableName() string { return "appointment_items" }

// AppointmentCredentialPO 预约凭证表
type AppointmentCredentialPO struct {
	ID                string    `gorm:"column:id;primaryKey;size:36"`
	AppointmentID     string    `gorm:"column:appointment_id;not null;size:36;uniqueIndex"`
	QRCodeData        string    `gorm:"column:qr_code_data;not null;type:text"`
	PatientNameMasked string    `gorm:"column:patient_name_masked;size:30"`
	ExamSummary       string    `gorm:"column:exam_summary;type:text"`
	NoticeContent     string    `gorm:"column:notice_content;type:text"`
	GeneratedAt       time.Time `gorm:"column:generated_at;not null"`
}

func (AppointmentCredentialPO) TableName() string { return "appointment_credentials" }

// AppointmentChangeLogPO 预约变更日志表
type AppointmentChangeLogPO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	AppointmentID string    `gorm:"column:appointment_id;not null;size:36;index:idx_change_log_appt"`
	ChangeType    string    `gorm:"column:change_type;not null;size:15"`
	OldSlotID     string    `gorm:"column:old_slot_id;size:36"`
	NewSlotID     string    `gorm:"column:new_slot_id;size:36"`
	Reason        string    `gorm:"column:reason;size:200"`
	OperatorID    string    `gorm:"column:operator_id;not null;size:36"`
	ChangedAt     time.Time `gorm:"column:changed_at;not null;autoCreateTime"`
}

func (AppointmentChangeLogPO) TableName() string { return "appointment_change_logs" }

// BlacklistPO 黑名单表
type BlacklistPO struct {
	ID            string     `gorm:"column:id;primaryKey;size:36"`
	PatientID     string     `gorm:"column:patient_id;not null;size:36;index:idx_blacklist_patient"`
	TriggerTime   time.Time  `gorm:"column:trigger_time;not null"`
	ExpiresAt     time.Time  `gorm:"column:expires_at;not null;index:idx_blacklist_expires"`
	Status        string     `gorm:"column:status;not null;size:10;default:active"`
	ReleasedAt    *time.Time `gorm:"column:released_at"`
	ReleaseReason string     `gorm:"column:release_reason;size:200"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (BlacklistPO) TableName() string { return "blacklists" }

// NoShowRecordPO 爽约记录表
type NoShowRecordPO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	PatientID     string    `gorm:"column:patient_id;not null;size:36;index:idx_noshow_patient"`
	AppointmentID string    `gorm:"column:appointment_id;not null;size:36"`
	OccurredAt    time.Time `gorm:"column:occurred_at;not null;autoCreateTime"`
}

func (NoShowRecordPO) TableName() string { return "no_show_records" }

// AppealPO 申诉表
type AppealPO struct {
	ID          string     `gorm:"column:id;primaryKey;size:36"`
	BlacklistID string     `gorm:"column:blacklist_id;not null;size:36;index:idx_appeal_blacklist"`
	Reason      string     `gorm:"column:reason;not null;size:500"`
	Status      string     `gorm:"column:status;not null;size:10;default:pending"`
	ReviewedBy  string     `gorm:"column:reviewed_by;size:36"`
	ReviewedAt  *time.Time `gorm:"column:reviewed_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (AppealPO) TableName() string { return "appeals" }
