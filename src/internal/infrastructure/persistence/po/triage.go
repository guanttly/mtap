// Package po 基础设施层 - triage 持久化对象（GORM Model）
package po

import "time"

// CheckInPO 签到记录表
type CheckInPO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	AppointmentID string    `gorm:"column:appointment_id;not null;size:36;uniqueIndex"`
	PatientID     string    `gorm:"column:patient_id;not null;size:36;index:idx_checkin_patient"`
	Method        string    `gorm:"column:method;not null;size:10"`
	CheckInTime   time.Time `gorm:"column:check_in_time;not null"`
	IsLate        bool      `gorm:"column:is_late;not null;default:false"`
	Remark        string    `gorm:"column:remark;size:100"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CheckInPO) TableName() string { return "check_ins" }

// WaitingQueuePO 候诊队列表
type WaitingQueuePO struct {
	ID           string    `gorm:"column:id;primaryKey;size:36"`
	RoomID       string    `gorm:"column:room_id;not null;size:36;uniqueIndex"`
	DeviceID     string    `gorm:"column:device_id;not null;size:36"`
	DepartmentID string    `gorm:"column:department_id;not null;size:36"`
	Status       string    `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (WaitingQueuePO) TableName() string { return "waiting_queues" }

// QueueEntryPO 队列条目表
type QueueEntryPO struct {
	ID                string     `gorm:"column:id;primaryKey;size:36"`
	QueueID           string     `gorm:"column:queue_id;not null;size:36;index:idx_queue_entries_queue"`
	PatientID         string     `gorm:"column:patient_id;not null;size:36"`
	PatientNameMasked string     `gorm:"column:patient_name_masked;size:30"`
	AppointmentID     string     `gorm:"column:appointment_id;not null;size:36"`
	CheckInID         string     `gorm:"column:check_in_id;not null;size:36"`
	QueueNumber       int        `gorm:"column:queue_number;not null"`
	Status            string     `gorm:"column:status;not null;size:15;default:waiting"`
	CallCount         int        `gorm:"column:call_count;not null;default:0"`
	MissCount         int        `gorm:"column:miss_count;not null;default:0"`
	EnteredAt         time.Time  `gorm:"column:entered_at;not null"`
	CalledAt          *time.Time `gorm:"column:called_at"`
	CompletedAt       *time.Time `gorm:"column:completed_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (QueueEntryPO) TableName() string { return "queue_entries" }

// ExamExecutionPO 检查执行表
type ExamExecutionPO struct {
	ID                string     `gorm:"column:id;primaryKey;size:36"`
	AppointmentItemID string     `gorm:"column:appointment_item_id;not null;size:36;uniqueIndex"`
	PatientID         string     `gorm:"column:patient_id;not null;size:36"`
	DeviceID          string     `gorm:"column:device_id;not null;size:36;index:idx_exec_device"`
	Status            string     `gorm:"column:status;not null;size:15;default:checked_in"`
	StartedAt         *time.Time `gorm:"column:started_at"`
	CompletedAt       *time.Time `gorm:"column:completed_at"`
	Duration          int        `gorm:"column:duration;default:0"`
	OperatorID        string     `gorm:"column:operator_id;size:36"`
	UndoDeadline      *time.Time `gorm:"column:undo_deadline"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (ExamExecutionPO) TableName() string { return "exam_executions" }
