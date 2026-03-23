// Package po 基础设施层 - resource持久化对象（GORM Model）
package po

import "time"

type CampusPO struct {
	ID        string    `gorm:"column:id;primaryKey;size:36"`
	Name      string    `gorm:"column:name;not null;size:50;uniqueIndex"`
	Code      string    `gorm:"column:code;not null;size:20;uniqueIndex"`
	Address   string    `gorm:"column:address;size:200"`
	Status    string    `gorm:"column:status;not null;size:10;default:active;index"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CampusPO) TableName() string { return "campuses" }

type DepartmentPO struct {
	ID        string     `gorm:"column:id;primaryKey;size:36"`
	CampusID  string     `gorm:"column:campus_id;not null;size:36;index"`
	Name      string     `gorm:"column:name;not null;size:50"`
	Code      string     `gorm:"column:code;not null;size:20;uniqueIndex"`
	Floor     string     `gorm:"column:floor;size:20"`
	Status    string     `gorm:"column:status;not null;size:10;default:active;index"`
	SyncedAt  *time.Time `gorm:"column:synced_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (DepartmentPO) TableName() string { return "departments" }

type DevicePO struct {
	ID                 string    `gorm:"column:id;primaryKey;size:36"`
	Name               string    `gorm:"column:name;not null;size:100;index"`
	Model              string    `gorm:"column:model;size:50"`
	Manufacturer       string    `gorm:"column:manufacturer;size:50"`
	SupportedExamTypes string    `gorm:"column:supported_exam_types;type:json;default:'[]'"` // JSON array
	MaxDailySlots      int       `gorm:"column:max_daily_slots;not null;default:50"`
	CampusID           string    `gorm:"column:campus_id;size:36;index"`
	DepartmentID       string    `gorm:"column:department_id;size:36;index"`
	Status             string    `gorm:"column:status;not null;size:10;default:active;index"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DevicePO) TableName() string { return "devices" }

type ExamItemPO struct {
	ID          string    `gorm:"column:id;primaryKey;size:36"`
	Name        string    `gorm:"column:name;not null;size:100;uniqueIndex"`
	DurationMin int       `gorm:"column:duration_min;not null"`
	IsFasting   bool      `gorm:"column:is_fasting;not null;default:false"`
	FastingDesc string    `gorm:"column:fasting_desc;size:200"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ExamItemPO) TableName() string { return "exam_items" }

type ItemAliasPO struct {
	ID         string    `gorm:"column:id;primaryKey;size:36"`
	ExamItemID string    `gorm:"column:exam_item_id;not null;size:36;index"`
	Alias      string    `gorm:"column:alias;not null;size:50;uniqueIndex"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ItemAliasPO) TableName() string { return "item_aliases" }

type SlotPoolPO struct {
	ID                 string    `gorm:"column:id;primaryKey;size:36"`
	Name               string    `gorm:"column:name;not null;size:60;uniqueIndex"`
	Type               string    `gorm:"column:type;not null;size:20"` // public/department/doctor
	Status             string    `gorm:"column:status;not null;size:10;default:active;index"`
	AllocationRatio    float64   `gorm:"column:allocation_ratio;not null;default:0"`
	OverflowEnabled    bool      `gorm:"column:overflow_enabled;not null;default:false"`
	OverflowTargetPool string    `gorm:"column:overflow_target_pool;size:36"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SlotPoolPO) TableName() string { return "slot_pools" }

type DoctorPO struct {
	ID           string     `gorm:"column:id;primaryKey;size:36"`
	DepartmentID string     `gorm:"column:department_id;not null;size:36;index"`
	HISCode      string     `gorm:"column:his_code;size:30"`
	Name         string     `gorm:"column:name;not null;size:30"`
	Title        string     `gorm:"column:title;size:20"`
	Gender       string     `gorm:"column:gender;not null;size:10;default:unknown"`
	Status       string     `gorm:"column:status;not null;size:10;default:active;index"`
	SyncedAt     *time.Time `gorm:"column:synced_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (DoctorPO) TableName() string { return "doctors" }

type ScheduleTemplatePO struct {
	ID           string    `gorm:"column:id;primaryKey;size:36"`
	Name         string    `gorm:"column:name;not null;size:50;uniqueIndex"`
	RepeatType   string    `gorm:"column:repeat_type;not null;size:20"`    // once/daily/weekly
	SlotPattern  string    `gorm:"column:slot_pattern;not null;type:json"` // JSON
	SkipWeekends bool      `gorm:"column:skip_weekends;not null;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ScheduleTemplatePO) TableName() string { return "schedule_templates" }

type SchedulePO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	DeviceID      string    `gorm:"column:device_id;not null;size:36;index"`
	Date          time.Time `gorm:"column:date;not null;index"`        // 只使用日期部分
	StartTime     string    `gorm:"column:start_time;not null;size:5"` // HH:mm
	EndTime       string    `gorm:"column:end_time;not null;size:5"`
	Status        string    `gorm:"column:status;not null;size:15;default:normal;index"`
	SuspendReason string    `gorm:"column:suspend_reason;size:200"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SchedulePO) TableName() string { return "schedules" }

type TimeSlotPO struct {
	ID               string     `gorm:"column:id;primaryKey;size:36"`
	DeviceID         string     `gorm:"column:device_id;not null;size:36;index"`
	Date             time.Time  `gorm:"column:date;not null;index"`
	ExamItemID       string     `gorm:"column:exam_item_id;size:36;index"`
	PoolType         string     `gorm:"column:pool_type;not null;size:15;default:public;index"`
	StartAt          time.Time  `gorm:"column:start_at;not null;index"`
	EndAt            time.Time  `gorm:"column:end_at;not null"`
	StandardDuration int        `gorm:"column:standard_duration;not null;default:0"`
	AdjustedDuration int        `gorm:"column:adjusted_duration;not null;default:0"`
	Status           string     `gorm:"column:status;not null;size:15;default:available;index"`
	LockedBy         string     `gorm:"column:locked_by;size:36;index"`
	LockUntil        *time.Time `gorm:"column:lock_until;index"`
	Remaining        int        `gorm:"column:remaining;not null;default:1"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (TimeSlotPO) TableName() string { return "time_slots" }
