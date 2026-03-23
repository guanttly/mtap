// Package po 基础设施层 - rule持久化对象（GORM Model）
package po

import "time"

// ConflictRulePO 冲突规则表
type ConflictRulePO struct {
	ID           string    `gorm:"column:id;primaryKey;size:36"`
	ItemAID      string    `gorm:"column:item_a_id;not null;size:36;index:idx_conflict_rules_item_a"`
	ItemBID      string    `gorm:"column:item_b_id;not null;size:36;index:idx_conflict_rules_item_b"`
	MinInterval  int       `gorm:"column:min_interval;not null;default:0"`
	IntervalUnit string    `gorm:"column:interval_unit;not null;size:10;default:hour"`
	Level        string    `gorm:"column:level;not null;size:10"`
	Status       string    `gorm:"column:status;not null;size:10;default:active;index:idx_conflict_rules_status"`
	CreatedBy    string    `gorm:"column:created_by;size:36"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ConflictRulePO) TableName() string { return "conflict_rules" }

// ConflictPackagePO 冲突包表
type ConflictPackagePO struct {
	ID           string                  `gorm:"column:id;primaryKey;size:36"`
	Name         string                  `gorm:"column:name;not null;size:30;uniqueIndex"`
	MinInterval  int                     `gorm:"column:min_interval;not null;default:0"`
	IntervalUnit string                  `gorm:"column:interval_unit;not null;size:10;default:hour"`
	Level        string                  `gorm:"column:level;not null;size:10"`
	Status       string                  `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt    time.Time               `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time               `gorm:"column:updated_at;autoUpdateTime"`
	Items        []ConflictPackageItemPO `gorm:"foreignKey:PackageID;references:ID"`
}

func (ConflictPackagePO) TableName() string { return "conflict_packages" }

// ConflictPackageItemPO 冲突包项目关联表
type ConflictPackageItemPO struct {
	ID         string    `gorm:"column:id;primaryKey;size:36"`
	PackageID  string    `gorm:"column:package_id;not null;size:36;index:idx_cpi_package"`
	ExamItemID string    `gorm:"column:exam_item_id;not null;size:36"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ConflictPackageItemPO) TableName() string { return "conflict_package_items" }

// DependencyRulePO 依赖规则表
type DependencyRulePO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	PreItemID     string    `gorm:"column:pre_item_id;not null;size:36"`
	PostItemID    string    `gorm:"column:post_item_id;not null;size:36"`
	Type          string    `gorm:"column:type;not null;size:15"`
	ValidityHours int       `gorm:"column:validity_hours;not null;default:72"`
	Status        string    `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DependencyRulePO) TableName() string { return "dependency_rules" }

// PriorityTagPO 优先级标签表
type PriorityTagPO struct {
	ID        string    `gorm:"column:id;primaryKey;size:36"`
	Name      string    `gorm:"column:name;not null;size:20;uniqueIndex"`
	Weight    int       `gorm:"column:weight;not null"`
	Color     string    `gorm:"column:color;not null;size:7"`
	IsPreset  bool      `gorm:"column:is_preset;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PriorityTagPO) TableName() string { return "priority_tags" }

// SortingStrategyPO 排序策略表
type SortingStrategyPO struct {
	ID            string    `gorm:"column:id;primaryKey;size:36"`
	Type          string    `gorm:"column:type;not null;size:20"`
	ScopeCampuses string    `gorm:"column:scope_campuses;type:text"` // JSON
	ScopeDepts    string    `gorm:"column:scope_depts;type:text"`    // JSON
	ScopeDevices  string    `gorm:"column:scope_devices;type:text"`  // JSON
	StartDate     time.Time `gorm:"column:start_date;not null"`
	EndDate       time.Time `gorm:"column:end_date;not null"`
	Status        string    `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SortingStrategyPO) TableName() string { return "sorting_strategies" }

// PatientAdaptRulePO 患者属性适配规则表
type PatientAdaptRulePO struct {
	ID             string    `gorm:"column:id;primaryKey;size:36"`
	ConditionType  string    `gorm:"column:condition_type;not null;size:20"`
	ConditionValue string    `gorm:"column:condition_value;not null;size:100"`
	Action         string    `gorm:"column:action;not null;size:20"`
	ActionParams   string    `gorm:"column:action_params;not null;type:text"` // JSON
	Priority       int       `gorm:"column:priority;not null;default:0"`
	Status         string    `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PatientAdaptRulePO) TableName() string { return "patient_adapt_rules" }

// SourceControlPO 开单来源控制表
type SourceControlPO struct {
	ID                   string    `gorm:"column:id;primaryKey;size:36"`
	SourceType           string    `gorm:"column:source_type;not null;size:15"`
	SlotPoolID           string    `gorm:"column:slot_pool_id;not null;size:36"`
	AllocationRatio      float64   `gorm:"column:allocation_ratio;not null"`
	OverflowEnabled      bool      `gorm:"column:overflow_enabled;not null;default:false"`
	OverflowTargetPoolID string    `gorm:"column:overflow_target_pool;size:36"`
	Status               string    `gorm:"column:status;not null;size:10;default:active"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SourceControlPO) TableName() string { return "source_controls" }
