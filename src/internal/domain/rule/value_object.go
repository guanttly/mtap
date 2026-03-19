package rule

// ConflictLevel 冲突级别
type ConflictLevel string

const (
	ConflictLevelForbid  ConflictLevel = "forbid"  // 禁止级 - 不可跳过
	ConflictLevelWarning ConflictLevel = "warning" // 警告级 - 确认后继续
)

// IsValid 校验冲突级别是否合法
func (l ConflictLevel) IsValid() bool {
	return l == ConflictLevelForbid || l == ConflictLevelWarning
}

// DependencyType 依赖类型
type DependencyType string

const (
	DependencyTypeMandatory   DependencyType = "mandatory"   // 强制依赖
	DependencyTypeRecommended DependencyType = "recommended" // 推荐依赖
)

// IsValid 校验依赖类型是否合法
func (t DependencyType) IsValid() bool {
	return t == DependencyTypeMandatory || t == DependencyTypeRecommended
}

// SortingType 排序策略类型
type SortingType string

const (
	SortingTypeShortestWait SortingType = "shortest_wait" // 等待时间最短
	SortingTypeNearest      SortingType = "nearest"       // 距离最近
	SortingTypePriority     SortingType = "priority"      // 指定优先级
)

// RuleStatus 规则状态
type RuleStatus string

const (
	RuleStatusActive   RuleStatus = "active"
	RuleStatusInactive RuleStatus = "inactive"
)

// EffectiveScope 生效范围值对象
type EffectiveScope struct {
	CampusIDs     []string `json:"campus_ids,omitempty"`
	DepartmentIDs []string `json:"department_ids,omitempty"`
	DeviceIDs     []string `json:"device_ids,omitempty"`
}

// IsEmpty 范围是否为空
func (s EffectiveScope) IsEmpty() bool {
	return len(s.CampusIDs) == 0 && len(s.DepartmentIDs) == 0 && len(s.DeviceIDs) == 0
}

// FastingItem 空腹项目标记值对象
type FastingItem struct {
	ExamItemID  string `json:"exam_item_id"`
	IsFasting   bool   `json:"is_fasting"`
	Description string `json:"description"`
}

// AdaptConditionType 适配条件类型
type AdaptConditionType string

const (
	AdaptConditionAge       AdaptConditionType = "age"
	AdaptConditionGender    AdaptConditionType = "gender"
	AdaptConditionPregnancy AdaptConditionType = "pregnancy"
)

// AdaptAction 适配动作
type AdaptAction string

const (
	AdaptActionFilterDevice AdaptAction = "filter_device"
	AdaptActionFilterSlot   AdaptAction = "filter_slot"
	AdaptActionFilterDoctor AdaptAction = "filter_doctor"
)

// OrderSource 开单来源
type OrderSource string

const (
	OrderSourceOutpatient OrderSource = "outpatient" // 门诊
	OrderSourceInpatient  OrderSource = "inpatient"  // 住院
	OrderSourceReferral   OrderSource = "referral"   // 转诊
)
