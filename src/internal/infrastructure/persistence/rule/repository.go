// Package rule 基础设施层 - rule仓储实现（GORM）
package rule

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "github.com/euler/mtap/internal/domain/rule"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

type Repositories struct {
	DB *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{DB: db}
}

// === ConflictRuleRepository ===

type ConflictRuleRepo struct{ db *gorm.DB }

func (r *Repositories) ConflictRuleRepo() *ConflictRuleRepo { return &ConflictRuleRepo{db: r.DB} }

func (r *ConflictRuleRepo) Save(ctx context.Context, rule *domain.ConflictRule) error {
	return r.db.WithContext(ctx).Create(&po.ConflictRulePO{
		ID:           rule.ID,
		ItemAID:      rule.ItemAID,
		ItemBID:      rule.ItemBID,
		MinInterval:  rule.MinInterval,
		IntervalUnit: rule.IntervalUnit,
		Level:        string(rule.Level),
		Status:       string(rule.Status),
		CreatedBy:    rule.CreatedBy,
		CreatedAt:    rule.CreatedAt,
		UpdatedAt:    rule.UpdatedAt,
	}).Error
}

func (r *ConflictRuleRepo) FindByID(ctx context.Context, id string) (*domain.ConflictRule, error) {
	var p po.ConflictRulePO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return conflictRuleFromPO(p), nil
}

func (r *ConflictRuleRepo) FindByItemPair(ctx context.Context, itemAID, itemBID string) (*domain.ConflictRule, error) {
	// 保持与领域层一致：pairKey 归一化
	if itemAID > itemBID {
		itemAID, itemBID = itemBID, itemAID
	}
	var p po.ConflictRulePO
	err := r.db.WithContext(ctx).
		Where("item_a_id = ? AND item_b_id = ?", itemAID, itemBID).
		First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return conflictRuleFromPO(p), nil
}

func (r *ConflictRuleRepo) FindAll(ctx context.Context, status domain.RuleStatus) ([]*domain.ConflictRule, error) {
	var ps []po.ConflictRulePO
	q := r.db.WithContext(ctx)
	if status != "" {
		q = q.Where("status = ?", string(status))
	}
	if err := q.Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.ConflictRule, 0, len(ps))
	for _, p := range ps {
		out = append(out, conflictRuleFromPO(p))
	}
	return out, nil
}

func (r *ConflictRuleRepo) Update(ctx context.Context, rule *domain.ConflictRule) error {
	return r.db.WithContext(ctx).Model(&po.ConflictRulePO{}).Where("id = ?", rule.ID).Updates(map[string]interface{}{
		"min_interval":  rule.MinInterval,
		"interval_unit": rule.IntervalUnit,
		"level":         string(rule.Level),
		"status":        string(rule.Status),
		"updated_at":    rule.UpdatedAt,
	}).Error
}

func (r *ConflictRuleRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&po.ConflictRulePO{}, "id = ?", id).Error
}

// === ConflictPackageRepository ===

type ConflictPackageRepo struct{ db *gorm.DB }

func (r *Repositories) ConflictPackageRepo() *ConflictPackageRepo {
	return &ConflictPackageRepo{db: r.DB}
}

func (r *ConflictPackageRepo) Save(ctx context.Context, pkg *domain.ConflictPackage) error {
	p := po.ConflictPackagePO{
		ID:           pkg.ID,
		Name:         pkg.Name,
		MinInterval:  pkg.MinInterval,
		IntervalUnit: pkg.IntervalUnit,
		Level:        string(pkg.Level),
		Status:       string(pkg.Status),
		CreatedAt:    pkg.CreatedAt,
		UpdatedAt:    pkg.UpdatedAt,
	}
	items := make([]po.ConflictPackageItemPO, 0, len(pkg.Items))
	for _, it := range pkg.Items {
		items = append(items, po.ConflictPackageItemPO{
			ID:         uuid.New().String(),
			PackageID:  pkg.ID,
			ExamItemID: it.ExamItemID,
		})
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ConflictPackageRepo) FindByID(ctx context.Context, id string) (*domain.ConflictPackage, error) {
	var p po.ConflictPackagePO
	err := r.db.WithContext(ctx).Preload("Items").First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return conflictPackageFromPO(p), nil
}

func (r *ConflictPackageRepo) FindByName(ctx context.Context, name string) (*domain.ConflictPackage, error) {
	var p po.ConflictPackagePO
	err := r.db.WithContext(ctx).Preload("Items").First(&p, "name = ?", name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return conflictPackageFromPO(p), nil
}

func (r *ConflictPackageRepo) FindAll(ctx context.Context) ([]*domain.ConflictPackage, error) {
	var ps []po.ConflictPackagePO
	if err := r.db.WithContext(ctx).Preload("Items").Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.ConflictPackage, 0, len(ps))
	for _, p := range ps {
		out = append(out, conflictPackageFromPO(p))
	}
	return out, nil
}

func (r *ConflictPackageRepo) Update(ctx context.Context, pkg *domain.ConflictPackage) error {
	// 简化实现：更新主表字段；items 由上层做“删除重建”式更新
	return r.db.WithContext(ctx).Model(&po.ConflictPackagePO{}).Where("id = ?", pkg.ID).Updates(map[string]interface{}{
		"name":          pkg.Name,
		"min_interval":  pkg.MinInterval,
		"interval_unit": pkg.IntervalUnit,
		"level":         string(pkg.Level),
		"status":        string(pkg.Status),
		"updated_at":    pkg.UpdatedAt,
	}).Error
}

func (r *ConflictPackageRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&po.ConflictPackageItemPO{}, "package_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&po.ConflictPackagePO{}, "id = ?", id).Error
	})
}

// === DependencyRuleRepository ===

type DependencyRuleRepo struct{ db *gorm.DB }

func (r *Repositories) DependencyRuleRepo() *DependencyRuleRepo { return &DependencyRuleRepo{db: r.DB} }

func (r *DependencyRuleRepo) Save(ctx context.Context, rule *domain.DependencyRule) error {
	return r.db.WithContext(ctx).Create(&po.DependencyRulePO{
		ID:            rule.ID,
		PreItemID:     rule.PreItemID,
		PostItemID:    rule.PostItemID,
		Type:          string(rule.Type),
		ValidityHours: rule.ValidityHours,
		Status:        string(rule.Status),
		CreatedAt:     rule.CreatedAt,
		UpdatedAt:     rule.UpdatedAt,
	}).Error
}

func (r *DependencyRuleRepo) FindByID(ctx context.Context, id string) (*domain.DependencyRule, error) {
	var p po.DependencyRulePO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return dependencyRuleFromPO(p), nil
}

func (r *DependencyRuleRepo) FindByPostItem(ctx context.Context, postItemID string) ([]*domain.DependencyRule, error) {
	var ps []po.DependencyRulePO
	if err := r.db.WithContext(ctx).Where("post_item_id = ?", postItemID).Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.DependencyRule, 0, len(ps))
	for _, p := range ps {
		out = append(out, dependencyRuleFromPO(p))
	}
	return out, nil
}

func (r *DependencyRuleRepo) FindAll(ctx context.Context, status domain.RuleStatus) ([]*domain.DependencyRule, error) {
	var ps []po.DependencyRulePO
	q := r.db.WithContext(ctx)
	if status != "" {
		q = q.Where("status = ?", string(status))
	}
	if err := q.Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.DependencyRule, 0, len(ps))
	for _, p := range ps {
		out = append(out, dependencyRuleFromPO(p))
	}
	return out, nil
}

func (r *DependencyRuleRepo) Update(ctx context.Context, rule *domain.DependencyRule) error {
	return r.db.WithContext(ctx).Model(&po.DependencyRulePO{}).Where("id = ?", rule.ID).Updates(map[string]interface{}{
		"type":           string(rule.Type),
		"validity_hours": rule.ValidityHours,
		"status":         string(rule.Status),
		"updated_at":     rule.UpdatedAt,
	}).Error
}

func (r *DependencyRuleRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&po.DependencyRulePO{}, "id = ?", id).Error
}

// === PriorityTagRepository ===

type PriorityTagRepo struct{ db *gorm.DB }

func (r *Repositories) PriorityTagRepo() *PriorityTagRepo { return &PriorityTagRepo{db: r.DB} }

func (r *PriorityTagRepo) Save(ctx context.Context, tag *domain.PriorityTag) error {
	return r.db.WithContext(ctx).Create(&po.PriorityTagPO{
		ID:        tag.ID,
		Name:      tag.Name,
		Weight:    tag.Weight,
		Color:     tag.Color,
		IsPreset:  tag.IsPreset,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}).Error
}

func (r *PriorityTagRepo) FindByID(ctx context.Context, id string) (*domain.PriorityTag, error) {
	var p po.PriorityTagPO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return priorityTagFromPO(p), nil
}

func (r *PriorityTagRepo) FindByName(ctx context.Context, name string) (*domain.PriorityTag, error) {
	var p po.PriorityTagPO
	err := r.db.WithContext(ctx).First(&p, "name = ?", name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return priorityTagFromPO(p), nil
}

func (r *PriorityTagRepo) FindAll(ctx context.Context) ([]*domain.PriorityTag, error) {
	var ps []po.PriorityTagPO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.PriorityTag, 0, len(ps))
	for _, p := range ps {
		out = append(out, priorityTagFromPO(p))
	}
	return out, nil
}

func (r *PriorityTagRepo) Update(ctx context.Context, tag *domain.PriorityTag) error {
	return r.db.WithContext(ctx).Model(&po.PriorityTagPO{}).Where("id = ?", tag.ID).Updates(map[string]interface{}{
		"name":       tag.Name,
		"weight":     tag.Weight,
		"color":      tag.Color,
		"is_preset":  tag.IsPreset,
		"updated_at": tag.UpdatedAt,
	}).Error
}

func (r *PriorityTagRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&po.PriorityTagPO{}, "id = ?", id).Error
}

// === SortingStrategyRepository ===

type SortingStrategyRepo struct{ db *gorm.DB }

func (r *Repositories) SortingStrategyRepo() *SortingStrategyRepo {
	return &SortingStrategyRepo{db: r.DB}
}

func (r *SortingStrategyRepo) Save(ctx context.Context, strategy *domain.SortingStrategy) error {
	campusJSON, _ := json.Marshal(strategy.Scope.CampusIDs)
	deptJSON, _ := json.Marshal(strategy.Scope.DepartmentIDs)
	deviceJSON, _ := json.Marshal(strategy.Scope.DeviceIDs)
	return r.db.WithContext(ctx).Create(&po.SortingStrategyPO{
		ID:            strategy.ID,
		Type:          string(strategy.Type),
		ScopeCampuses: string(campusJSON),
		ScopeDepts:    string(deptJSON),
		ScopeDevices:  string(deviceJSON),
		StartDate:     strategy.StartDate,
		EndDate:       strategy.EndDate,
		Status:        string(strategy.Status),
		CreatedAt:     strategy.CreatedAt,
		UpdatedAt:     strategy.UpdatedAt,
	}).Error
}

func (r *SortingStrategyRepo) FindByID(ctx context.Context, id string) (*domain.SortingStrategy, error) {
	var p po.SortingStrategyPO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return sortingStrategyFromPO(p), nil
}

func (r *SortingStrategyRepo) FindAll(ctx context.Context) ([]*domain.SortingStrategy, error) {
	var ps []po.SortingStrategyPO
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.SortingStrategy, 0, len(ps))
	for _, p := range ps {
		out = append(out, sortingStrategyFromPO(p))
	}
	return out, nil
}

func (r *SortingStrategyRepo) FindByScope(ctx context.Context, scope domain.EffectiveScope) ([]*domain.SortingStrategy, error) {
	// 简化：先全量查，再用 JSON scope 精确匹配（避免依赖特定 DB 的 JSONB 能力）
	var ps []po.SortingStrategyPO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.SortingStrategy, 0)
	for _, p := range ps {
		s := sortingStrategyFromPO(p)
		if equalScope(s.Scope, scope) {
			out = append(out, s)
		}
	}
	return out, nil
}

func (r *SortingStrategyRepo) Update(ctx context.Context, strategy *domain.SortingStrategy) error {
	campusJSON, _ := json.Marshal(strategy.Scope.CampusIDs)
	deptJSON, _ := json.Marshal(strategy.Scope.DepartmentIDs)
	deviceJSON, _ := json.Marshal(strategy.Scope.DeviceIDs)
	return r.db.WithContext(ctx).Model(&po.SortingStrategyPO{}).Where("id = ?", strategy.ID).Updates(map[string]interface{}{
		"type":           string(strategy.Type),
		"scope_campuses": string(campusJSON),
		"scope_depts":    string(deptJSON),
		"scope_devices":  string(deviceJSON),
		"start_date":     strategy.StartDate,
		"end_date":       strategy.EndDate,
		"status":         string(strategy.Status),
		"updated_at":     strategy.UpdatedAt,
	}).Error
}

func (r *SortingStrategyRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&po.SortingStrategyPO{}, "id = ?", id).Error
}

// === PatientAdaptRuleRepository ===

type PatientAdaptRuleRepo struct{ db *gorm.DB }

func (r *Repositories) PatientAdaptRuleRepo() *PatientAdaptRuleRepo {
	return &PatientAdaptRuleRepo{db: r.DB}
}

func (r *PatientAdaptRuleRepo) SaveAll(ctx context.Context, rules []*domain.PatientAdaptRule) error {
	ps := make([]po.PatientAdaptRulePO, 0, len(rules))
	for _, rule := range rules {
		paramsJSON, _ := json.Marshal(rule.ActionParams)
		ps = append(ps, po.PatientAdaptRulePO{
			ID:             rule.ID,
			ConditionType:  string(rule.ConditionType),
			ConditionValue: rule.ConditionValue,
			Action:         string(rule.Action),
			ActionParams:   string(paramsJSON),
			Priority:       rule.Priority,
			Status:         string(rule.Status),
			CreatedAt:      rule.CreatedAt,
			UpdatedAt:      rule.UpdatedAt,
		})
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&po.PatientAdaptRulePO{}, "1=1").Error; err != nil {
			return err
		}
		if len(ps) == 0 {
			return nil
		}
		return tx.Create(&ps).Error
	})
}

func (r *PatientAdaptRuleRepo) FindAll(ctx context.Context) ([]*domain.PatientAdaptRule, error) {
	var ps []po.PatientAdaptRulePO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.PatientAdaptRule, 0, len(ps))
	for _, p := range ps {
		out = append(out, patientAdaptRuleFromPO(p))
	}
	return out, nil
}

func (r *PatientAdaptRuleRepo) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&po.PatientAdaptRulePO{}, "1=1").Error
}

// === SourceControlRepository ===

type SourceControlRepo struct{ db *gorm.DB }

func (r *Repositories) SourceControlRepo() *SourceControlRepo { return &SourceControlRepo{db: r.DB} }

func (r *SourceControlRepo) SaveAll(ctx context.Context, controls []*domain.SourceControl) error {
	ps := make([]po.SourceControlPO, 0, len(controls))
	for _, c := range controls {
		ps = append(ps, po.SourceControlPO{
			ID:                   c.ID,
			SourceType:           string(c.SourceType),
			SlotPoolID:           c.SlotPoolID,
			AllocationRatio:      c.AllocationRatio,
			OverflowEnabled:      c.OverflowEnabled,
			OverflowTargetPoolID: c.OverflowTargetPoolID,
			Status:               string(c.Status),
			CreatedAt:            c.CreatedAt,
			UpdatedAt:            c.UpdatedAt,
		})
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&po.SourceControlPO{}, "1=1").Error; err != nil {
			return err
		}
		if len(ps) == 0 {
			return nil
		}
		return tx.Create(&ps).Error
	})
}

func (r *SourceControlRepo) FindAll(ctx context.Context) ([]*domain.SourceControl, error) {
	var ps []po.SourceControlPO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.SourceControl, 0, len(ps))
	for _, p := range ps {
		out = append(out, sourceControlFromPO(p))
	}
	return out, nil
}

func (r *SourceControlRepo) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&po.SourceControlPO{}, "1=1").Error
}

// === PO -> Domain converters ===

func conflictRuleFromPO(p po.ConflictRulePO) *domain.ConflictRule {
	return &domain.ConflictRule{
		ID:           p.ID,
		ItemAID:      p.ItemAID,
		ItemBID:      p.ItemBID,
		MinInterval:  p.MinInterval,
		IntervalUnit: p.IntervalUnit,
		Level:        domain.ConflictLevel(p.Level),
		Status:       domain.RuleStatus(p.Status),
		CreatedBy:    p.CreatedBy,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func conflictPackageFromPO(p po.ConflictPackagePO) *domain.ConflictPackage {
	items := make([]domain.ConflictPackageItem, 0, len(p.Items))
	for _, it := range p.Items {
		items = append(items, domain.ConflictPackageItem{
			PackageID:  it.PackageID,
			ExamItemID: it.ExamItemID,
		})
	}
	return &domain.ConflictPackage{
		ID:           p.ID,
		Name:         p.Name,
		MinInterval:  p.MinInterval,
		IntervalUnit: p.IntervalUnit,
		Level:        domain.ConflictLevel(p.Level),
		Status:       domain.RuleStatus(p.Status),
		Items:        items,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func dependencyRuleFromPO(p po.DependencyRulePO) *domain.DependencyRule {
	return &domain.DependencyRule{
		ID:            p.ID,
		PreItemID:     p.PreItemID,
		PostItemID:    p.PostItemID,
		Type:          domain.DependencyType(p.Type),
		ValidityHours: p.ValidityHours,
		Status:        domain.RuleStatus(p.Status),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func priorityTagFromPO(p po.PriorityTagPO) *domain.PriorityTag {
	return &domain.PriorityTag{
		ID:        p.ID,
		Name:      p.Name,
		Weight:    p.Weight,
		Color:     p.Color,
		IsPreset:  p.IsPreset,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func sortingStrategyFromPO(p po.SortingStrategyPO) *domain.SortingStrategy {
	var campusIDs, deptIDs, deviceIDs []string
	_ = json.Unmarshal([]byte(p.ScopeCampuses), &campusIDs)
	_ = json.Unmarshal([]byte(p.ScopeDepts), &deptIDs)
	_ = json.Unmarshal([]byte(p.ScopeDevices), &deviceIDs)
	return &domain.SortingStrategy{
		ID:   p.ID,
		Type: domain.SortingType(p.Type),
		Scope: domain.EffectiveScope{
			CampusIDs:     campusIDs,
			DepartmentIDs: deptIDs,
			DeviceIDs:     deviceIDs,
		},
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		Status:    domain.RuleStatus(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func patientAdaptRuleFromPO(p po.PatientAdaptRulePO) *domain.PatientAdaptRule {
	params := map[string]string{}
	_ = json.Unmarshal([]byte(p.ActionParams), &params)
	return &domain.PatientAdaptRule{
		ID:             p.ID,
		ConditionType:  domain.AdaptConditionType(p.ConditionType),
		ConditionValue: p.ConditionValue,
		Action:         domain.AdaptAction(p.Action),
		ActionParams:   params,
		Priority:       p.Priority,
		Status:         domain.RuleStatus(p.Status),
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func sourceControlFromPO(p po.SourceControlPO) *domain.SourceControl {
	return &domain.SourceControl{
		ID:                   p.ID,
		SourceType:           domain.OrderSource(p.SourceType),
		SlotPoolID:           p.SlotPoolID,
		AllocationRatio:      p.AllocationRatio,
		OverflowEnabled:      p.OverflowEnabled,
		OverflowTargetPoolID: p.OverflowTargetPoolID,
		Status:               domain.RuleStatus(p.Status),
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}

func equalScope(a, b domain.EffectiveScope) bool {
	return sliceSetEqual(a.CampusIDs, b.CampusIDs) &&
		sliceSetEqual(a.DepartmentIDs, b.DepartmentIDs) &&
		sliceSetEqual(a.DeviceIDs, b.DeviceIDs)
}

func sliceSetEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, x := range a {
		m[x]++
	}
	for _, x := range b {
		if m[x] == 0 {
			return false
		}
		m[x]--
	}
	return true
}
