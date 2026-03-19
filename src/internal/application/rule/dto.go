// Package rule 应用层 - 规则引擎DTO
package rule

import "time"

// === 冲突规则 DTO ===

type CreateConflictRuleReq struct {
	ItemAID     string `json:"item_a_id" binding:"required"`
	ItemBID     string `json:"item_b_id" binding:"required"`
	MinInterval int    `json:"min_interval" binding:"gte=0,lte=720"`
	Level       string `json:"level" binding:"required,oneof=forbid warning"`
}

type UpdateConflictRuleReq struct {
	MinInterval *int    `json:"min_interval" binding:"omitempty,gte=0,lte=720"`
	Level       *string `json:"level" binding:"omitempty,oneof=forbid warning"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type ConflictRuleResp struct {
	ID           string    `json:"id"`
	ItemAID      string    `json:"item_a_id"`
	ItemBID      string    `json:"item_b_id"`
	MinInterval  int       `json:"min_interval"`
	IntervalUnit string    `json:"interval_unit"`
	Level        string    `json:"level"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// === 冲突包 DTO ===

type CreateConflictPackageReq struct {
	Name        string   `json:"name" binding:"required,max=30"`
	ItemIDs     []string `json:"item_ids" binding:"required,min=2"`
	MinInterval int      `json:"min_interval" binding:"gte=0,lte=720"`
	Level       string   `json:"level" binding:"required,oneof=forbid warning"`
}

type UpdateConflictPackageReq struct {
	Name        *string  `json:"name" binding:"omitempty,max=30"`
	ItemIDs     []string `json:"item_ids" binding:"omitempty,min=2"`
	MinInterval *int     `json:"min_interval" binding:"omitempty,gte=0,lte=720"`
	Level       *string  `json:"level" binding:"omitempty,oneof=forbid warning"`
}

type ConflictPackageResp struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	MinInterval  int                       `json:"min_interval"`
	IntervalUnit string                    `json:"interval_unit"`
	Level        string                    `json:"level"`
	Status       string                    `json:"status"`
	Items        []ConflictPackageItemResp `json:"items"`
	CreatedAt    time.Time                 `json:"created_at"`
}

type ConflictPackageItemResp struct {
	ExamItemID string `json:"exam_item_id"`
}

// === 依赖规则 DTO ===

type CreateDependencyRuleReq struct {
	PreItemID     string `json:"pre_item_id" binding:"required"`
	PostItemID    string `json:"post_item_id" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=mandatory recommended"`
	ValidityHours int    `json:"validity_hours" binding:"required,gt=0"`
}

type UpdateDependencyRuleReq struct {
	Type          *string `json:"type" binding:"omitempty,oneof=mandatory recommended"`
	ValidityHours *int    `json:"validity_hours" binding:"omitempty,gt=0"`
	Status        *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type DependencyRuleResp struct {
	ID            string    `json:"id"`
	PreItemID     string    `json:"pre_item_id"`
	PostItemID    string    `json:"post_item_id"`
	Type          string    `json:"type"`
	ValidityHours int       `json:"validity_hours"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// === 优先级标签 DTO ===

type CreatePriorityTagReq struct {
	Name   string `json:"name" binding:"required,max=20"`
	Weight int    `json:"weight" binding:"required,gte=1,lte=100"`
	Color  string `json:"color" binding:"required"`
}

type UpdatePriorityTagReq struct {
	Name   *string `json:"name" binding:"omitempty,max=20"`
	Weight *int    `json:"weight" binding:"omitempty,gte=1,lte=100"`
	Color  *string `json:"color"`
}

type PriorityTagResp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Weight   int    `json:"weight"`
	Color    string `json:"color"`
	IsPreset bool   `json:"is_preset"`
}

// === 排序策略 DTO ===

type SaveSortingStrategyReq struct {
	Type      string        `json:"type" binding:"required,oneof=shortest_wait nearest priority"`
	Scope     EffectiveScopeDTO `json:"scope" binding:"required"`
	StartDate time.Time     `json:"start_date" binding:"required"`
	EndDate   time.Time     `json:"end_date" binding:"required"`
}

type SortingStrategyResp struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Scope     EffectiveScopeDTO `json:"scope"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"end_date"`
	Status    string          `json:"status"`
}

type EffectiveScopeDTO struct {
	CampusIDs     []string `json:"campus_ids,omitempty" form:"campus_ids[]"`
	DepartmentIDs []string `json:"department_ids,omitempty" form:"department_ids[]"`
	DeviceIDs     []string `json:"device_ids,omitempty" form:"device_ids[]"`
}

// === 患者属性适配 DTO ===

type SavePatientAdaptRuleReq struct {
	ConditionType  string            `json:"condition_type" binding:"required,oneof=age gender pregnancy"`
	ConditionValue string            `json:"condition_value" binding:"required"`
	Action         string            `json:"action" binding:"required,oneof=filter_device filter_slot filter_doctor"`
	ActionParams   map[string]string `json:"action_params" binding:"required"`
	Priority       int               `json:"priority"`
}

type PatientAdaptRuleResp struct {
	ID             string            `json:"id"`
	ConditionType  string            `json:"condition_type"`
	ConditionValue string            `json:"condition_value"`
	Action         string            `json:"action"`
	ActionParams   map[string]string `json:"action_params"`
	Priority       int               `json:"priority"`
	Status         string            `json:"status"`
}

// === 开单来源控制 DTO ===

type SaveSourceControlReq struct {
	SourceType           string   `json:"source_type" binding:"required,oneof=outpatient inpatient referral"`
	SlotPoolID           string   `json:"slot_pool_id" binding:"required"`
	AllocationRatio      float64  `json:"allocation_ratio" binding:"required,gte=0,lte=1"`
	OverflowEnabled      bool     `json:"overflow_enabled"`
	OverflowTargetPoolID string   `json:"overflow_target_pool_id"`
}

type SourceControlResp struct {
	ID                   string  `json:"id"`
	SourceType           string  `json:"source_type"`
	SlotPoolID           string  `json:"slot_pool_id"`
	AllocationRatio      float64 `json:"allocation_ratio"`
	OverflowEnabled      bool    `json:"overflow_enabled"`
	OverflowTargetPoolID string  `json:"overflow_target_pool_id,omitempty"`
	Status               string  `json:"status"`
}

// === 规则校验 DTO ===

type RuleCheckReq struct {
	PatientID   string      `json:"patient_id" binding:"required"`
	PatientAttr PatientAttr `json:"patient_attr"`
	ExamItemIDs []string    `json:"exam_item_ids" binding:"required,min=1"`
	OrderSource string      `json:"order_source" binding:"required,oneof=outpatient inpatient referral"`
}

type PatientAttr struct {
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	IsPregnant bool   `json:"is_pregnant"`
}

type RuleCheckResp struct {
	Conflicts       []ConflictResultResp   `json:"conflicts"`
	Dependencies    []DependencyResultResp `json:"dependencies"`
	HasForbidden    bool                   `json:"has_forbidden"`
	HasBlocked      bool                   `json:"has_blocked"`
	FastingItems    []string               `json:"fasting_items"`
	PriorityScore   int                    `json:"priority_score"`
	FilteredPoolIDs []string               `json:"filtered_pool_ids"`
	Warnings        []string               `json:"warnings"`
}

type ConflictResultResp struct {
	ItemAID        string `json:"item_a_id"`
	ItemBID        string `json:"item_b_id"`
	Level          string `json:"level"`
	MinInterval    int    `json:"min_interval"`
	ActualInterval int    `json:"actual_interval"`
	Reason         string `json:"reason"`
}

type DependencyResultResp struct {
	PostItemID    string     `json:"post_item_id"`
	PreItemID     string     `json:"pre_item_id"`
	Type          string     `json:"type"`
	Status        string     `json:"status"`
	ValidityHours int        `json:"validity_hours"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	Reason        string     `json:"reason"`
}

// === 通用查询 ===

type ListReq struct {
	Page     int    `form:"page" binding:"gte=1"`
	PageSize int    `form:"page_size" binding:"gte=1,lte=100"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive"`
	Keyword  string `form:"keyword"`
}

func (r *ListReq) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 20
	}
}
