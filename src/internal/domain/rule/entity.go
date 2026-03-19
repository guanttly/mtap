// Package rule 规则引擎领域层 - 实体与聚合根定义
package rule

import (
	"time"

	"github.com/google/uuid"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// ConflictRule 冲突规则聚合根
type ConflictRule struct {
	ID           string        `json:"id"`
	ItemAID      string        `json:"item_a_id"`
	ItemBID      string        `json:"item_b_id"`
	MinInterval  int           `json:"min_interval"`
	IntervalUnit string        `json:"interval_unit"` // hour
	Level        ConflictLevel `json:"level"`
	Status       RuleStatus    `json:"status"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// NewConflictRule 创建冲突规则
func NewConflictRule(itemAID, itemBID string, minInterval int, level ConflictLevel, createdBy string) (*ConflictRule, error) {
	if itemAID == itemBID {
		return nil, bizErr.New(bizErr.ErrRuleSameItem)
	}
	if minInterval < 0 || minInterval > 720 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "最小间隔应在0~720小时之间")
	}
	if !level.IsValid() {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "冲突级别无效")
	}

	// 确保 itemAID < itemBID 以保证唯一性
	if itemAID > itemBID {
		itemAID, itemBID = itemBID, itemAID
	}

	now := time.Now()
	return &ConflictRule{
		ID:           uuid.New().String(),
		ItemAID:      itemAID,
		ItemBID:      itemBID,
		MinInterval:  minInterval,
		IntervalUnit: "hour",
		Level:        level,
		Status:       RuleStatusActive,
		CreatedBy:    createdBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate 校验冲突规则的业务不变量
func (r *ConflictRule) Validate() error {
	if r.ItemAID == r.ItemBID {
		return bizErr.New(bizErr.ErrRuleSameItem)
	}
	if r.MinInterval < 0 || r.MinInterval > 720 {
		return bizErr.NewWithDetail(bizErr.ErrInvalidParam, "最小间隔应在0~720小时之间")
	}
	if !r.Level.IsValid() {
		return bizErr.NewWithDetail(bizErr.ErrInvalidParam, "冲突级别无效")
	}
	return nil
}

// ConflictPackage 冲突包聚合根
type ConflictPackage struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	MinInterval  int                 `json:"min_interval"`
	IntervalUnit string              `json:"interval_unit"`
	Level        ConflictLevel       `json:"level"`
	Status       RuleStatus          `json:"status"`
	Items        []ConflictPackageItem `json:"items"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// ConflictPackageItem 冲突包项目
type ConflictPackageItem struct {
	PackageID  string `json:"package_id"`
	ExamItemID string `json:"exam_item_id"`
}

// NewConflictPackage 创建冲突包
func NewConflictPackage(name string, itemIDs []string, minInterval int, level ConflictLevel) (*ConflictPackage, error) {
	if len(name) == 0 || len(name) > 30 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "冲突包名称长度应在1~30之间")
	}
	if len(itemIDs) < 2 {
		return nil, bizErr.New(bizErr.ErrRulePkgTooFew)
	}

	pkgID := uuid.New().String()
	items := make([]ConflictPackageItem, len(itemIDs))
	for i, id := range itemIDs {
		items[i] = ConflictPackageItem{PackageID: pkgID, ExamItemID: id}
	}

	now := time.Now()
	return &ConflictPackage{
		ID:           pkgID,
		Name:         name,
		MinInterval:  minInterval,
		IntervalUnit: "hour",
		Level:        level,
		Status:       RuleStatusActive,
		Items:        items,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// AddItem 添加项目到冲突包
func (p *ConflictPackage) AddItem(examItemID string) error {
	for _, item := range p.Items {
		if item.ExamItemID == examItemID {
			return bizErr.NewWithDetail(bizErr.ErrDuplicate, "项目已在冲突包中")
		}
	}
	p.Items = append(p.Items, ConflictPackageItem{PackageID: p.ID, ExamItemID: examItemID})
	p.UpdatedAt = time.Now()
	return nil
}

// RemoveItem 从冲突包移除项目
func (p *ConflictPackage) RemoveItem(examItemID string) error {
	newItems := make([]ConflictPackageItem, 0, len(p.Items))
	found := false
	for _, item := range p.Items {
		if item.ExamItemID == examItemID {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}
	if !found {
		return bizErr.New(bizErr.ErrNotFound)
	}
	p.Items = newItems
	p.UpdatedAt = time.Now()
	return nil
}

// IsValid 冲突包是否有效（至少2个项目）
func (p *ConflictPackage) IsValid() bool {
	return len(p.Items) >= 2
}

// DependencyRule 依赖规则聚合根
type DependencyRule struct {
	ID            string         `json:"id"`
	PreItemID     string         `json:"pre_item_id"`
	PostItemID    string         `json:"post_item_id"`
	Type          DependencyType `json:"type"`
	ValidityHours int            `json:"validity_hours"`
	Status        RuleStatus     `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// NewDependencyRule 创建依赖规则
func NewDependencyRule(preItemID, postItemID string, depType DependencyType, validityHours int) (*DependencyRule, error) {
	if preItemID == postItemID {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "前后置项目不能相同")
	}
	if validityHours <= 0 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "时效必须大于0")
	}
	if !depType.IsValid() {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "依赖类型无效")
	}

	now := time.Now()
	return &DependencyRule{
		ID:            uuid.New().String(),
		PreItemID:     preItemID,
		PostItemID:    postItemID,
		Type:          depType,
		ValidityHours: validityHours,
		Status:        RuleStatusActive,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// PriorityTag 优先级标签聚合根
type PriorityTag struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Weight    int        `json:"weight"`
	Color     string     `json:"color"`
	IsPreset  bool       `json:"is_preset"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewPriorityTag 创建优先级标签
func NewPriorityTag(name string, weight int, color string) (*PriorityTag, error) {
	if len(name) == 0 || len(name) > 20 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "标签名称长度应在1~20之间")
	}
	if weight < 1 || weight > 100 {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "权重应在1~100之间")
	}

	now := time.Now()
	return &PriorityTag{
		ID:        uuid.New().String(),
		Name:      name,
		Weight:    weight,
		Color:     color,
		IsPreset:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CanDelete 判断标签是否可删除
func (t *PriorityTag) CanDelete() bool {
	return !t.IsPreset
}

// SortingStrategy 排序策略聚合根
type SortingStrategy struct {
	ID        string         `json:"id"`
	Type      SortingType    `json:"type"`
	Scope     EffectiveScope `json:"scope"`
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	Status    RuleStatus     `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// NewSortingStrategy 创建排序策略
func NewSortingStrategy(sortType SortingType, scope EffectiveScope, startDate, endDate time.Time) (*SortingStrategy, error) {
	if !endDate.After(startDate) {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "结束日期必须晚于开始日期")
	}

	now := time.Now()
	return &SortingStrategy{
		ID:        uuid.New().String(),
		Type:      sortType,
		Scope:     scope,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    RuleStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsEffective 策略在指定日期是否生效
func (s *SortingStrategy) IsEffective(date time.Time) bool {
	return s.Status == RuleStatusActive &&
		!date.Before(s.StartDate) &&
		!date.After(s.EndDate)
}

// PatientAdaptRule 患者属性适配规则
type PatientAdaptRule struct {
	ID             string             `json:"id"`
	ConditionType  AdaptConditionType `json:"condition_type"`
	ConditionValue string             `json:"condition_value"`
	Action         AdaptAction        `json:"action"`
	ActionParams   map[string]string  `json:"action_params"`
	Priority       int                `json:"priority"`
	Status         RuleStatus         `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// SourceControl 开单来源控制
type SourceControl struct {
	ID                   string      `json:"id"`
	SourceType           OrderSource `json:"source_type"`
	SlotPoolID           string      `json:"slot_pool_id"`
	AllocationRatio      float64     `json:"allocation_ratio"`
	OverflowEnabled      bool        `json:"overflow_enabled"`
	OverflowTargetPoolID string      `json:"overflow_target_pool_id,omitempty"`
	Status               RuleStatus  `json:"status"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}
