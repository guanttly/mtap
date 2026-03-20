// Package rule 应用层 - 规则引擎应用服务
package rule

import (
	"context"
	"strconv"
	"strings"
	"time"

	domain "github.com/euler/mtap/internal/domain/rule"
	bizErr "github.com/euler/mtap/pkg/errors"
	"github.com/google/uuid"
)

// RuleAppService 规则管理应用服务
type RuleAppService struct {
	conflictRuleRepo domain.ConflictRuleRepository
	conflictPkgRepo  domain.ConflictPackageRepository
	depRuleRepo      domain.DependencyRuleRepository
	tagRepo          domain.PriorityTagRepository
	sortingRepo      domain.SortingStrategyRepository
	adaptRepo        domain.PatientAdaptRuleRepository
	sourceRepo       domain.SourceControlRepository

	patientHistory PatientHistoryProvider
	examItemMeta   ExamItemMetaProvider

	conflictSvc     *domain.ConflictDetectionService
	depSvc          *domain.DependencyValidationService
	circularChecker *domain.CircularDependencyChecker
}

// PatientHistoryProvider 查询患者历史检查记录（供冲突时间间隔判断）
type PatientHistoryProvider interface {
	ListRecent(ctx context.Context, patientID string, examItemIDs []string) ([]domain.PatientExamRecord, error)
}

// ExamItemMetaProvider 查询项目元数据（空腹标记等）
type ExamItemMetaProvider interface {
	GetFastingItemIDs(ctx context.Context, examItemIDs []string) ([]string, error)
}

// NewRuleAppService 创建规则应用服务
func NewRuleAppService(
	conflictRuleRepo domain.ConflictRuleRepository,
	conflictPkgRepo domain.ConflictPackageRepository,
	depRuleRepo domain.DependencyRuleRepository,
	tagRepo domain.PriorityTagRepository,
	sortingRepo domain.SortingStrategyRepository,
	adaptRepo domain.PatientAdaptRuleRepository,
	sourceRepo domain.SourceControlRepository,
	conflictSvc *domain.ConflictDetectionService,
	depSvc *domain.DependencyValidationService,
	circularChecker *domain.CircularDependencyChecker,
) *RuleAppService {
	return &RuleAppService{
		conflictRuleRepo: conflictRuleRepo,
		conflictPkgRepo:  conflictPkgRepo,
		depRuleRepo:      depRuleRepo,
		tagRepo:          tagRepo,
		sortingRepo:      sortingRepo,
		adaptRepo:        adaptRepo,
		sourceRepo:       sourceRepo,
		conflictSvc:      conflictSvc,
		depSvc:           depSvc,
		circularChecker:  circularChecker,
	}
}

func (s *RuleAppService) WithPatientHistoryProvider(p PatientHistoryProvider) *RuleAppService {
	s.patientHistory = p
	return s
}

func (s *RuleAppService) WithExamItemMetaProvider(p ExamItemMetaProvider) *RuleAppService {
	s.examItemMeta = p
	return s
}

// === 冲突规则 CRUD ===

func (s *RuleAppService) CreateConflictRule(ctx context.Context, req CreateConflictRuleReq, operatorID string) (*ConflictRuleResp, error) {
	// 检查是否已存在
	existing, err := s.conflictRuleRepo.FindByItemPair(ctx, req.ItemAID, req.ItemBID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if existing != nil {
		return nil, bizErr.New(bizErr.ErrRuleDuplicate)
	}

	rule, err := domain.NewConflictRule(req.ItemAID, req.ItemBID, req.MinInterval, domain.ConflictLevel(req.Level), operatorID)
	if err != nil {
		return nil, err
	}

	if err := s.conflictRuleRepo.Save(ctx, rule); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return ConflictRuleToResp(rule), nil
}

func (s *RuleAppService) GetConflictRule(ctx context.Context, id string) (*ConflictRuleResp, error) {
	rule, err := s.conflictRuleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if rule == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ConflictRuleToResp(rule), nil
}

func (s *RuleAppService) ListConflictRules(ctx context.Context, req ListReq) ([]ConflictRuleResp, error) {
	req.SetDefaults()
	status := domain.RuleStatusActive
	if req.Status != "" {
		status = domain.RuleStatus(req.Status)
	}
	rules, err := s.conflictRuleRepo.FindAll(ctx, status)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ConflictRulesToResp(rules), nil
}

func (s *RuleAppService) DeleteConflictRule(ctx context.Context, id string) error {
	rule, err := s.conflictRuleRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if rule == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	return s.conflictRuleRepo.Delete(ctx, id)
}

// === 冲突包 CRUD ===

func (s *RuleAppService) CreateConflictPackage(ctx context.Context, req CreateConflictPackageReq) (*ConflictPackageResp, error) {
	existing, err := s.conflictPkgRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if existing != nil {
		return nil, bizErr.New(bizErr.ErrRulePkgNameDup)
	}

	pkg, err := domain.NewConflictPackage(req.Name, req.ItemIDs, req.MinInterval, domain.ConflictLevel(req.Level))
	if err != nil {
		return nil, err
	}

	if err := s.conflictPkgRepo.Save(ctx, pkg); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return ConflictPackageToResp(pkg), nil
}

func (s *RuleAppService) ListConflictPackages(ctx context.Context) ([]ConflictPackageResp, error) {
	pkgs, err := s.conflictPkgRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	result := make([]ConflictPackageResp, len(pkgs))
	for i, pkg := range pkgs {
		result[i] = *ConflictPackageToResp(pkg)
	}
	return result, nil
}

func (s *RuleAppService) DeleteConflictPackage(ctx context.Context, id string) error {
	pkg, err := s.conflictPkgRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if pkg == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	return s.conflictPkgRepo.Delete(ctx, id)
}

// === 依赖规则 CRUD ===

func (s *RuleAppService) CreateDependencyRule(ctx context.Context, req CreateDependencyRuleReq) (*DependencyRuleResp, error) {
	// 循环依赖检测
	hasCycle, err := s.circularChecker.HasCircularDependency(ctx, req.PreItemID, req.PostItemID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if hasCycle {
		return nil, bizErr.New(bizErr.ErrRuleCircularDep)
	}

	rule, err := domain.NewDependencyRule(req.PreItemID, req.PostItemID, domain.DependencyType(req.Type), req.ValidityHours)
	if err != nil {
		return nil, err
	}

	if err := s.depRuleRepo.Save(ctx, rule); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return DependencyRuleToResp(rule), nil
}

func (s *RuleAppService) ListDependencyRules(ctx context.Context, req ListReq) ([]DependencyRuleResp, error) {
	req.SetDefaults()
	status := domain.RuleStatusActive
	if req.Status != "" {
		status = domain.RuleStatus(req.Status)
	}
	rules, err := s.depRuleRepo.FindAll(ctx, status)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return DependencyRulesToResp(rules), nil
}

func (s *RuleAppService) DeleteDependencyRule(ctx context.Context, id string) error {
	rule, err := s.depRuleRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if rule == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	return s.depRuleRepo.Delete(ctx, id)
}

// === 优先级标签 CRUD ===

func (s *RuleAppService) CreatePriorityTag(ctx context.Context, req CreatePriorityTagReq) (*PriorityTagResp, error) {
	existing, err := s.tagRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if existing != nil {
		return nil, bizErr.New(bizErr.ErrRuleTagNameDup)
	}

	tag, err := domain.NewPriorityTag(req.Name, req.Weight, req.Color)
	if err != nil {
		return nil, err
	}

	if err := s.tagRepo.Save(ctx, tag); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return PriorityTagToResp(tag), nil
}

func (s *RuleAppService) ListPriorityTags(ctx context.Context) ([]PriorityTagResp, error) {
	tags, err := s.tagRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return PriorityTagsToResp(tags), nil
}

func (s *RuleAppService) DeletePriorityTag(ctx context.Context, id string) error {
	tag, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if tag == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if !tag.CanDelete() {
		return bizErr.New(bizErr.ErrRulePresetNoDelete)
	}
	return s.tagRepo.Delete(ctx, id)
}

// === 更新操作 ===

func (s *RuleAppService) UpdateConflictRule(ctx context.Context, id string, req UpdateConflictRuleReq) (*ConflictRuleResp, error) {
	rule, err := s.conflictRuleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if rule == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if req.MinInterval != nil {
		rule.MinInterval = *req.MinInterval
	}
	if req.Level != nil {
		rule.Level = domain.ConflictLevel(*req.Level)
	}
	if req.Status != nil {
		rule.Status = domain.RuleStatus(*req.Status)
	}
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	rule.UpdatedAt = time.Now()
	if err := s.conflictRuleRepo.Update(ctx, rule); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ConflictRuleToResp(rule), nil
}

func (s *RuleAppService) UpdateConflictPackage(ctx context.Context, id string, req UpdateConflictPackageReq) (*ConflictPackageResp, error) {
	pkg, err := s.conflictPkgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if pkg == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if req.Name != nil {
		pkg.Name = *req.Name
	}
	if req.MinInterval != nil {
		pkg.MinInterval = *req.MinInterval
	}
	if req.Level != nil {
		pkg.Level = domain.ConflictLevel(*req.Level)
	}
	if len(req.ItemIDs) > 0 {
		if len(req.ItemIDs) < 2 {
			return nil, bizErr.New(bizErr.ErrRulePkgTooFew)
		}
		items := make([]domain.ConflictPackageItem, len(req.ItemIDs))
		for i, itemID := range req.ItemIDs {
			items[i] = domain.ConflictPackageItem{PackageID: pkg.ID, ExamItemID: itemID}
		}
		pkg.Items = items
	}
	pkg.UpdatedAt = time.Now()
	if err := s.conflictPkgRepo.Update(ctx, pkg); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ConflictPackageToResp(pkg), nil
}

func (s *RuleAppService) UpdateDependencyRule(ctx context.Context, id string, req UpdateDependencyRuleReq) (*DependencyRuleResp, error) {
	rule, err := s.depRuleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if rule == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if req.Type != nil {
		rule.Type = domain.DependencyType(*req.Type)
	}
	if req.ValidityHours != nil {
		rule.ValidityHours = *req.ValidityHours
	}
	if req.Status != nil {
		rule.Status = domain.RuleStatus(*req.Status)
	}
	rule.UpdatedAt = time.Now()
	if err := s.depRuleRepo.Update(ctx, rule); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return DependencyRuleToResp(rule), nil
}

func (s *RuleAppService) UpdatePriorityTag(ctx context.Context, id string, req UpdatePriorityTagReq) (*PriorityTagResp, error) {
	tag, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if tag == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Weight != nil {
		tag.Weight = *req.Weight
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}
	tag.UpdatedAt = time.Now()
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return PriorityTagToResp(tag), nil
}

// === 排序策略 ===

func (s *RuleAppService) SaveSortingStrategy(ctx context.Context, req SaveSortingStrategyReq) (*SortingStrategyResp, error) {
	scope := ScopeFromDTO(req.Scope)
	if scope.IsEmpty() {
		return nil, bizErr.New(bizErr.ErrRuleInvalidScope)
	}

	existing, err := s.sortingRepo.FindByScope(ctx, scope)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	for _, e := range existing {
		if e.Status != domain.RuleStatusActive {
			continue
		}
		// 日期区间有交集则冲突
		if !(req.EndDate.Before(e.StartDate) || req.StartDate.After(e.EndDate)) {
			return nil, bizErr.New(bizErr.ErrRuleSortingConflict)
		}
	}

	strategy, err := domain.NewSortingStrategy(domain.SortingType(req.Type), scope, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	if err := s.sortingRepo.Save(ctx, strategy); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return SortingStrategyToResp(strategy), nil
}

func (s *RuleAppService) GetSortingStrategy(ctx context.Context, scopeDTO EffectiveScopeDTO) (*SortingStrategyResp, error) {
	scope := ScopeFromDTO(scopeDTO)
	if scope.IsEmpty() {
		return nil, bizErr.New(bizErr.ErrRuleInvalidScope)
	}
	list, err := s.sortingRepo.FindByScope(ctx, scope)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	now := time.Now()
	var best *domain.SortingStrategy
	for _, st := range list {
		if st.Status != domain.RuleStatusActive || !st.IsEffective(now) {
			continue
		}
		if best == nil || st.EndDate.After(best.EndDate) {
			best = st
		}
	}
	if best == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return SortingStrategyToResp(best), nil
}

// === 患者属性适配 ===

func (s *RuleAppService) SavePatientAdaptRules(ctx context.Context, req []SavePatientAdaptRuleReq) error {
	now := time.Now()
	rules := make([]*domain.PatientAdaptRule, 0, len(req))
	for _, r := range req {
		rules = append(rules, &domain.PatientAdaptRule{
			ID:             uuid.New().String(),
			ConditionType:  domain.AdaptConditionType(r.ConditionType),
			ConditionValue: r.ConditionValue,
			Action:         domain.AdaptAction(r.Action),
			ActionParams:   r.ActionParams,
			Priority:       r.Priority,
			Status:         domain.RuleStatusActive,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
	if err := s.adaptRepo.SaveAll(ctx, rules); err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return nil
}

func (s *RuleAppService) ListPatientAdaptRules(ctx context.Context) ([]PatientAdaptRuleResp, error) {
	rules, err := s.adaptRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return PatientAdaptRulesToResp(rules), nil
}

// === 开单来源控制 ===

func (s *RuleAppService) SaveSourceControls(ctx context.Context, req []SaveSourceControlReq) error {
	now := time.Now()
	controls := make([]*domain.SourceControl, 0, len(req))
	for _, c0 := range req {
		controls = append(controls, &domain.SourceControl{
			ID:                   uuid.New().String(),
			SourceType:           domain.OrderSource(c0.SourceType),
			SlotPoolID:           c0.SlotPoolID,
			AllocationRatio:      c0.AllocationRatio,
			OverflowEnabled:      c0.OverflowEnabled,
			OverflowTargetPoolID: c0.OverflowTargetPoolID,
			Status:               domain.RuleStatusActive,
			CreatedAt:            now,
			UpdatedAt:            now,
		})
	}
	if err := s.sourceRepo.SaveAll(ctx, controls); err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return nil
}

func (s *RuleAppService) ListSourceControls(ctx context.Context) ([]SourceControlResp, error) {
	controls, err := s.sourceRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return SourceControlsToResp(controls), nil
}

// === 综合规则校验 ===

func (s *RuleAppService) CheckRules(ctx context.Context, req RuleCheckReq) (*RuleCheckResp, error) {
	// 0. 获取患者历史（用于冲突间隔判断）
	var history []domain.PatientExamRecord
	if s.patientHistory != nil {
		h, err := s.patientHistory.ListRecent(ctx, req.PatientID, req.ExamItemIDs)
		if err != nil {
			return nil, bizErr.Wrap(bizErr.ErrInternal, err)
		}
		history = h
	}

	// 1. 冲突检测
	conflicts, err := s.conflictSvc.DetectConflicts(ctx, req.ExamItemIDs, history)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	// 2. 依赖校验
	deps, err := s.depSvc.ValidateDependencies(ctx, req.ExamItemIDs, req.PatientID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	// 3. 构建结果
	resp := &RuleCheckResp{
		Conflicts:    ConflictResultsToResp(conflicts),
		Dependencies: DependencyResultsToResp(deps),
		Warnings:     []string{},
	}

	for _, c := range conflicts {
		if c.Level == domain.ConflictLevelForbid {
			resp.HasForbidden = true
			break
		}
	}

	for _, d := range deps {
		if d.Status == domain.DependencyBlocked && d.Type == domain.DependencyTypeMandatory {
			resp.HasBlocked = true
			break
		}
	}

	// 3.1 衍生 warnings（冲突 warning、依赖 unknown/expired）
	for _, c := range conflicts {
		if c.Level == domain.ConflictLevelWarning {
			resp.Warnings = append(resp.Warnings, c.Reason)
		}
	}
	for _, d := range deps {
		if d.Status == domain.DependencyUnknown {
			resp.Warnings = append(resp.Warnings, "前置检查状态未知（HIS查询失败），已标记待人工复核")
		}
		if d.Status == domain.DependencyExpired {
			resp.Warnings = append(resp.Warnings, "存在前置检查已过时效的项目，请确认是否需要紧急放行")
		}
	}

	// 4. 空腹项目前置识别
	if s.examItemMeta != nil {
		fastingIDs, err := s.examItemMeta.GetFastingItemIDs(ctx, req.ExamItemIDs)
		if err != nil {
			return nil, bizErr.Wrap(bizErr.ErrInternal, err)
		}
		resp.FastingItems = fastingIDs
		if len(fastingIDs) > 0 {
			resp.Warnings = append(resp.Warnings, "存在空腹项目，请安排在当日最早时段")
		}
	}

	// 5. 优先级评分（基于来源的最小可解释实现；后续可叠加标签/策略）
	resp.PriorityScore = 50
	switch req.OrderSource {
	case "inpatient":
		resp.PriorityScore += 10
	case "referral":
		resp.PriorityScore += 5
	}

	// 6. 开单来源控制 → 号源池过滤
	controls, err := s.sourceRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	for _, c0 := range controls {
		if c0.Status != domain.RuleStatusActive {
			continue
		}
		if string(c0.SourceType) == req.OrderSource {
			resp.FilteredPoolIDs = append(resp.FilteredPoolIDs, c0.SlotPoolID)
			if c0.OverflowEnabled && c0.OverflowTargetPoolID != "" {
				resp.FilteredPoolIDs = append(resp.FilteredPoolIDs, c0.OverflowTargetPoolID)
			}
		}
	}

	// 7. 患者属性适配（对 filtered_pool_ids 做补充过滤建议 + warnings）
	// 说明：适配动作在完整系统中应作用于“可用号源集合”；此处先把规则匹配结果体现在 filtered_pool_ids / warnings 上，便于验收与后续与资源模块对接。
	rules, err := s.adaptRepo.FindAll(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	for _, r := range rules {
		if r.Status != domain.RuleStatusActive {
			continue
		}
		if !matchPatientAdaptRule(req.PatientAttr, r) {
			continue
		}
		resp.Warnings = append(resp.Warnings, "命中患者属性适配规则："+string(r.ConditionType))
		if poolID, ok := r.ActionParams["slot_pool_id"]; ok && poolID != "" {
			resp.FilteredPoolIDs = append(resp.FilteredPoolIDs, poolID)
		}
	}

	return resp, nil
}

func matchPatientAdaptRule(attr PatientAttr, r *domain.PatientAdaptRule) bool {
	switch r.ConditionType {
	case domain.AdaptConditionGender:
		return strings.EqualFold(attr.Gender, r.ConditionValue)
	case domain.AdaptConditionPregnancy:
		v := strings.ToLower(strings.TrimSpace(r.ConditionValue))
		want := v == "true" || v == "1" || v == "yes"
		return attr.IsPregnant == want
	case domain.AdaptConditionAge:
		return matchAgeCondition(attr.Age, r.ConditionValue)
	default:
		return false
	}
}

func matchAgeCondition(age int, cond string) bool {
	cond = strings.TrimSpace(cond)
	if cond == "" {
		return false
	}
	// 支持："<14"、">70"、"8-14"
	if strings.HasPrefix(cond, "<") {
		n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(cond, "<")))
		return err == nil && age < n
	}
	if strings.HasPrefix(cond, ">") {
		n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(cond, ">")))
		return err == nil && age > n
	}
	if strings.Contains(cond, "-") {
		parts := strings.SplitN(cond, "-", 2)
		minV, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		maxV, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		return err1 == nil && err2 == nil && age >= minV && age <= maxV
	}
	n, err := strconv.Atoi(cond)
	return err == nil && age == n
}
