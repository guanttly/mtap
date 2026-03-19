// Package rule 应用层 - DTO与领域对象的装配器
package rule

import (
	domain "github.com/euler/mtap/internal/domain/rule"
)

// === ConflictRule 装配 ===

func ConflictRuleToResp(r *domain.ConflictRule) *ConflictRuleResp {
	return &ConflictRuleResp{
		ID:           r.ID,
		ItemAID:      r.ItemAID,
		ItemBID:      r.ItemBID,
		MinInterval:  r.MinInterval,
		IntervalUnit: r.IntervalUnit,
		Level:        string(r.Level),
		Status:       string(r.Status),
		CreatedBy:    r.CreatedBy,
		CreatedAt:    r.CreatedAt,
	}
}

func ConflictRulesToResp(rules []*domain.ConflictRule) []ConflictRuleResp {
	result := make([]ConflictRuleResp, len(rules))
	for i, r := range rules {
		result[i] = *ConflictRuleToResp(r)
	}
	return result
}

// === ConflictPackage 装配 ===

func ConflictPackageToResp(pkg *domain.ConflictPackage) *ConflictPackageResp {
	items := make([]ConflictPackageItemResp, len(pkg.Items))
	for i, item := range pkg.Items {
		items[i] = ConflictPackageItemResp{ExamItemID: item.ExamItemID}
	}
	return &ConflictPackageResp{
		ID:           pkg.ID,
		Name:         pkg.Name,
		MinInterval:  pkg.MinInterval,
		IntervalUnit: pkg.IntervalUnit,
		Level:        string(pkg.Level),
		Status:       string(pkg.Status),
		Items:        items,
		CreatedAt:    pkg.CreatedAt,
	}
}

// === DependencyRule 装配 ===

func DependencyRuleToResp(r *domain.DependencyRule) *DependencyRuleResp {
	return &DependencyRuleResp{
		ID:            r.ID,
		PreItemID:     r.PreItemID,
		PostItemID:    r.PostItemID,
		Type:          string(r.Type),
		ValidityHours: r.ValidityHours,
		Status:        string(r.Status),
		CreatedAt:     r.CreatedAt,
	}
}

func DependencyRulesToResp(rules []*domain.DependencyRule) []DependencyRuleResp {
	result := make([]DependencyRuleResp, len(rules))
	for i, r := range rules {
		result[i] = *DependencyRuleToResp(r)
	}
	return result
}

// === PriorityTag 装配 ===

func PriorityTagToResp(t *domain.PriorityTag) *PriorityTagResp {
	return &PriorityTagResp{
		ID:       t.ID,
		Name:     t.Name,
		Weight:   t.Weight,
		Color:    t.Color,
		IsPreset: t.IsPreset,
	}
}

func PriorityTagsToResp(tags []*domain.PriorityTag) []PriorityTagResp {
	result := make([]PriorityTagResp, len(tags))
	for i, t := range tags {
		result[i] = *PriorityTagToResp(t)
	}
	return result
}

// === SortingStrategy 装配 ===

func ScopeToDTO(s domain.EffectiveScope) EffectiveScopeDTO {
	return EffectiveScopeDTO{
		CampusIDs:     s.CampusIDs,
		DepartmentIDs: s.DepartmentIDs,
		DeviceIDs:     s.DeviceIDs,
	}
}

func ScopeFromDTO(s EffectiveScopeDTO) domain.EffectiveScope {
	return domain.EffectiveScope{
		CampusIDs:     s.CampusIDs,
		DepartmentIDs: s.DepartmentIDs,
		DeviceIDs:     s.DeviceIDs,
	}
}

func SortingStrategyToResp(s *domain.SortingStrategy) *SortingStrategyResp {
	return &SortingStrategyResp{
		ID:        s.ID,
		Type:      string(s.Type),
		Scope:     ScopeToDTO(s.Scope),
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		Status:    string(s.Status),
	}
}

// === PatientAdaptRule 装配 ===

func PatientAdaptRuleToResp(r *domain.PatientAdaptRule) *PatientAdaptRuleResp {
	return &PatientAdaptRuleResp{
		ID:             r.ID,
		ConditionType:  string(r.ConditionType),
		ConditionValue: r.ConditionValue,
		Action:         string(r.Action),
		ActionParams:   r.ActionParams,
		Priority:       r.Priority,
		Status:         string(r.Status),
	}
}

func PatientAdaptRulesToResp(rules []*domain.PatientAdaptRule) []PatientAdaptRuleResp {
	out := make([]PatientAdaptRuleResp, 0, len(rules))
	for _, r := range rules {
		out = append(out, *PatientAdaptRuleToResp(r))
	}
	return out
}

// === SourceControl 装配 ===

func SourceControlToResp(c *domain.SourceControl) *SourceControlResp {
	return &SourceControlResp{
		ID:                   c.ID,
		SourceType:           string(c.SourceType),
		SlotPoolID:           c.SlotPoolID,
		AllocationRatio:      c.AllocationRatio,
		OverflowEnabled:      c.OverflowEnabled,
		OverflowTargetPoolID: c.OverflowTargetPoolID,
		Status:               string(c.Status),
	}
}

func SourceControlsToResp(controls []*domain.SourceControl) []SourceControlResp {
	out := make([]SourceControlResp, 0, len(controls))
	for _, c := range controls {
		out = append(out, *SourceControlToResp(c))
	}
	return out
}

// === ConflictResult 装配 ===

func ConflictResultToResp(r *domain.ConflictResult) ConflictResultResp {
	return ConflictResultResp{
		ItemAID:        r.ItemAID,
		ItemBID:        r.ItemBID,
		Level:          string(r.Level),
		MinInterval:    r.MinInterval,
		ActualInterval: r.ActualInterval,
		Reason:         r.Reason,
	}
}

func ConflictResultsToResp(results []domain.ConflictResult) []ConflictResultResp {
	resp := make([]ConflictResultResp, len(results))
	for i, r := range results {
		resp[i] = ConflictResultToResp(&r)
	}
	return resp
}

// === DependencyResult 装配 ===

func DependencyResultToResp(r *domain.DependencyResult) DependencyResultResp {
	return DependencyResultResp{
		PostItemID:    r.PostItemID,
		PreItemID:     r.PreItemID,
		Type:          string(r.Type),
		Status:        string(r.Status),
		ValidityHours: r.ValidityHours,
		CompletedAt:   r.CompletedAt,
		Reason:        r.Reason,
	}
}

func DependencyResultsToResp(results []domain.DependencyResult) []DependencyResultResp {
	resp := make([]DependencyResultResp, len(results))
	for i, r := range results {
		resp[i] = DependencyResultToResp(&r)
	}
	return resp
}
