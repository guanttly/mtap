package rule

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/euler/mtap/internal/domain/rule"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// === Mock Repositories (application-layer tests) ===

type mockConflictRuleRepo struct {
	rules map[string]*domain.ConflictRule
}

func newMockConflictRuleRepo() *mockConflictRuleRepo {
	return &mockConflictRuleRepo{rules: make(map[string]*domain.ConflictRule)}
}

func (m *mockConflictRuleRepo) Save(_ context.Context, rule *domain.ConflictRule) error {
	m.rules[rule.ID] = rule
	return nil
}
func (m *mockConflictRuleRepo) FindByID(_ context.Context, id string) (*domain.ConflictRule, error) {
	return m.rules[id], nil
}
func (m *mockConflictRuleRepo) FindByItemPair(_ context.Context, a, b string) (*domain.ConflictRule, error) {
	if a > b {
		a, b = b, a
	}
	for _, r := range m.rules {
		if r.ItemAID == a && r.ItemBID == b {
			return r, nil
		}
	}
	return nil, nil
}
func (m *mockConflictRuleRepo) FindAll(_ context.Context, status domain.RuleStatus) ([]*domain.ConflictRule, error) {
	var result []*domain.ConflictRule
	for _, r := range m.rules {
		if r.Status == status {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockConflictRuleRepo) Update(_ context.Context, _ *domain.ConflictRule) error { return nil }
func (m *mockConflictRuleRepo) Delete(_ context.Context, id string) error {
	delete(m.rules, id)
	return nil
}

type mockConflictPkgRepo struct {
	pkgs map[string]*domain.ConflictPackage
}

func newMockConflictPkgRepo() *mockConflictPkgRepo {
	return &mockConflictPkgRepo{pkgs: make(map[string]*domain.ConflictPackage)}
}

func (m *mockConflictPkgRepo) Save(_ context.Context, pkg *domain.ConflictPackage) error {
	m.pkgs[pkg.ID] = pkg
	return nil
}
func (m *mockConflictPkgRepo) FindByID(_ context.Context, id string) (*domain.ConflictPackage, error) {
	return m.pkgs[id], nil
}
func (m *mockConflictPkgRepo) FindByName(_ context.Context, name string) (*domain.ConflictPackage, error) {
	for _, p := range m.pkgs {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, nil
}
func (m *mockConflictPkgRepo) FindAll(_ context.Context) ([]*domain.ConflictPackage, error) {
	var result []*domain.ConflictPackage
	for _, p := range m.pkgs {
		result = append(result, p)
	}
	return result, nil
}
func (m *mockConflictPkgRepo) Update(_ context.Context, _ *domain.ConflictPackage) error { return nil }
func (m *mockConflictPkgRepo) Delete(_ context.Context, id string) error {
	delete(m.pkgs, id)
	return nil
}

type mockDepRuleRepo struct {
	rules map[string]*domain.DependencyRule
}

func newMockDepRuleRepo() *mockDepRuleRepo {
	return &mockDepRuleRepo{rules: make(map[string]*domain.DependencyRule)}
}

func (m *mockDepRuleRepo) Save(_ context.Context, rule *domain.DependencyRule) error {
	m.rules[rule.ID] = rule
	return nil
}
func (m *mockDepRuleRepo) FindByID(_ context.Context, id string) (*domain.DependencyRule, error) {
	return m.rules[id], nil
}
func (m *mockDepRuleRepo) FindByPostItem(_ context.Context, postItemID string) ([]*domain.DependencyRule, error) {
	var result []*domain.DependencyRule
	for _, r := range m.rules {
		if r.PostItemID == postItemID {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockDepRuleRepo) FindAll(_ context.Context, status domain.RuleStatus) ([]*domain.DependencyRule, error) {
	var result []*domain.DependencyRule
	for _, r := range m.rules {
		if r.Status == status {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockDepRuleRepo) Update(_ context.Context, _ *domain.DependencyRule) error { return nil }
func (m *mockDepRuleRepo) Delete(_ context.Context, id string) error {
	delete(m.rules, id)
	return nil
}

type mockTagRepo struct {
	tags map[string]*domain.PriorityTag
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{tags: make(map[string]*domain.PriorityTag)}
}

func (m *mockTagRepo) Save(_ context.Context, tag *domain.PriorityTag) error {
	m.tags[tag.ID] = tag
	return nil
}
func (m *mockTagRepo) FindByID(_ context.Context, id string) (*domain.PriorityTag, error) {
	return m.tags[id], nil
}
func (m *mockTagRepo) FindByName(_ context.Context, name string) (*domain.PriorityTag, error) {
	for _, t := range m.tags {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, nil
}
func (m *mockTagRepo) FindAll(_ context.Context) ([]*domain.PriorityTag, error) {
	var result []*domain.PriorityTag
	for _, t := range m.tags {
		result = append(result, t)
	}
	return result, nil
}
func (m *mockTagRepo) Update(_ context.Context, _ *domain.PriorityTag) error { return nil }
func (m *mockTagRepo) Delete(_ context.Context, id string) error {
	delete(m.tags, id)
	return nil
}

type mockPreItemChecker struct{}

func (m *mockPreItemChecker) GetCompletedTime(_ context.Context, _, _ string) (*time.Time, error) {
	return nil, nil
}

type mockSortingRepo struct {
	items []*domain.SortingStrategy
}

func (m *mockSortingRepo) Save(_ context.Context, s *domain.SortingStrategy) error {
	m.items = append(m.items, s)
	return nil
}
func (m *mockSortingRepo) FindByID(_ context.Context, _ string) (*domain.SortingStrategy, error) {
	return nil, nil
}
func (m *mockSortingRepo) FindAll(_ context.Context) ([]*domain.SortingStrategy, error) {
	return m.items, nil
}
func (m *mockSortingRepo) FindByScope(_ context.Context, _ domain.EffectiveScope) ([]*domain.SortingStrategy, error) {
	return m.items, nil
}
func (m *mockSortingRepo) Update(_ context.Context, _ *domain.SortingStrategy) error { return nil }
func (m *mockSortingRepo) Delete(_ context.Context, _ string) error                  { return nil }

type mockAdaptRepo struct {
	rules []*domain.PatientAdaptRule
}

func (m *mockAdaptRepo) SaveAll(_ context.Context, rules []*domain.PatientAdaptRule) error {
	m.rules = rules
	return nil
}
func (m *mockAdaptRepo) FindAll(_ context.Context) ([]*domain.PatientAdaptRule, error) {
	return m.rules, nil
}
func (m *mockAdaptRepo) DeleteAll(_ context.Context) error {
	m.rules = nil
	return nil
}

type mockSourceRepo struct {
	controls []*domain.SourceControl
}

func (m *mockSourceRepo) SaveAll(_ context.Context, controls []*domain.SourceControl) error {
	m.controls = controls
	return nil
}
func (m *mockSourceRepo) FindAll(_ context.Context) ([]*domain.SourceControl, error) {
	return m.controls, nil
}
func (m *mockSourceRepo) DeleteAll(_ context.Context) error {
	m.controls = nil
	return nil
}

// === Helper ===

func newTestAppService() *RuleAppService {
	crRepo := newMockConflictRuleRepo()
	cpRepo := newMockConflictPkgRepo()
	drRepo := newMockDepRuleRepo()
	tRepo := newMockTagRepo()
	sortRepo := &mockSortingRepo{}
	adaptRepo := &mockAdaptRepo{}
	sourceRepo := &mockSourceRepo{}
	conflictSvc := domain.NewConflictDetectionService(crRepo, cpRepo)
	depSvc := domain.NewDependencyValidationService(drRepo, &mockPreItemChecker{})
	checker := domain.NewCircularDependencyChecker(drRepo)
	return NewRuleAppService(crRepo, cpRepo, drRepo, tRepo, sortRepo, adaptRepo, sourceRepo, conflictSvc, depSvc, checker)
}

// === ConflictRule Application Service Tests ===

func TestCreateConflictRule_OK(t *testing.T) {
	svc := newTestAppService()
	resp, err := svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: "warning",
	}, "admin")

	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "warning", resp.Level)
}

func TestCreateConflictRule_Duplicate(t *testing.T) {
	svc := newTestAppService()
	_, _ = svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: "warning",
	}, "admin")

	_, err := svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E2", MinInterval: 48, Level: "forbid",
	}, "admin")
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleDuplicate))
}

func TestCreateConflictRule_SameItem(t *testing.T) {
	svc := newTestAppService()
	_, err := svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E1", MinInterval: 24, Level: "warning",
	}, "admin")
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleSameItem))
}

func TestGetConflictRule_NotFound(t *testing.T) {
	svc := newTestAppService()
	_, err := svc.GetConflictRule(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

func TestDeleteConflictRule_OK(t *testing.T) {
	svc := newTestAppService()
	resp, _ := svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: "warning",
	}, "admin")

	err := svc.DeleteConflictRule(context.Background(), resp.ID)
	assert.NoError(t, err)
}

// === ConflictPackage Application Service Tests ===

func TestCreateConflictPackage_OK(t *testing.T) {
	svc := newTestAppService()
	resp, err := svc.CreateConflictPackage(context.Background(), CreateConflictPackageReq{
		Name: "造影剂互斥", ItemIDs: []string{"E1", "E2", "E3"}, MinInterval: 24, Level: "forbid",
	})

	require.NoError(t, err)
	assert.Equal(t, "造影剂互斥", resp.Name)
	assert.Equal(t, 3, len(resp.Items))
}

func TestCreateConflictPackage_DuplicateName(t *testing.T) {
	svc := newTestAppService()
	_, _ = svc.CreateConflictPackage(context.Background(), CreateConflictPackageReq{
		Name: "造影剂互斥", ItemIDs: []string{"E1", "E2"}, MinInterval: 24, Level: "forbid",
	})
	_, err := svc.CreateConflictPackage(context.Background(), CreateConflictPackageReq{
		Name: "造影剂互斥", ItemIDs: []string{"E3", "E4"}, MinInterval: 24, Level: "forbid",
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrRulePkgNameDup))
}

// === DependencyRule Application Service Tests ===

func TestCreateDependencyRule_OK(t *testing.T) {
	svc := newTestAppService()
	resp, err := svc.CreateDependencyRule(context.Background(), CreateDependencyRuleReq{
		PreItemID: "BLOOD", PostItemID: "GASTRO", Type: "mandatory", ValidityHours: 72,
	})
	require.NoError(t, err)
	assert.Equal(t, "BLOOD", resp.PreItemID)
}

func TestCreateDependencyRule_CircularDep(t *testing.T) {
	svc := newTestAppService()
	// A → B
	_, _ = svc.CreateDependencyRule(context.Background(), CreateDependencyRuleReq{
		PreItemID: "A", PostItemID: "B", Type: "mandatory", ValidityHours: 72,
	})
	// try B → A (creates cycle)
	_, err := svc.CreateDependencyRule(context.Background(), CreateDependencyRuleReq{
		PreItemID: "B", PostItemID: "A", Type: "mandatory", ValidityHours: 72,
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleCircularDep))
}

// === PriorityTag Application Service Tests ===

func TestCreatePriorityTag_OK(t *testing.T) {
	svc := newTestAppService()
	resp, err := svc.CreatePriorityTag(context.Background(), CreatePriorityTagReq{
		Name: "急诊", Weight: 90, Color: "#FF4D4F",
	})
	require.NoError(t, err)
	assert.Equal(t, "急诊", resp.Name)
}

func TestCreatePriorityTag_DuplicateName(t *testing.T) {
	svc := newTestAppService()
	_, _ = svc.CreatePriorityTag(context.Background(), CreatePriorityTagReq{
		Name: "急诊", Weight: 90, Color: "#FF4D4F",
	})
	_, err := svc.CreatePriorityTag(context.Background(), CreatePriorityTagReq{
		Name: "急诊", Weight: 80, Color: "#000",
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleTagNameDup))
}

func TestDeletePriorityTag_PresetForbidden(t *testing.T) {
	svc := newTestAppService()
	// Manually insert a preset tag
	presetTag := &domain.PriorityTag{ID: "PT_PRESET", Name: "预置", Weight: 50, IsPreset: true}
	svc.tagRepo.Save(context.Background(), presetTag)

	err := svc.DeletePriorityTag(context.Background(), "PT_PRESET")
	assert.True(t, bizErr.Is(err, bizErr.ErrRulePresetNoDelete))
}

// === CheckRules Tests ===

func TestCheckRules_NoConflicts(t *testing.T) {
	svc := newTestAppService()
	resp, err := svc.CheckRules(context.Background(), RuleCheckReq{
		PatientID:   "P001",
		ExamItemIDs: []string{"E1", "E2"},
		OrderSource: "outpatient",
	})
	require.NoError(t, err)
	assert.False(t, resp.HasForbidden)
	assert.False(t, resp.HasBlocked)
}

func TestCheckRules_WithForbiddenConflict(t *testing.T) {
	svc := newTestAppService()
	// Create a forbid-level rule
	_, _ = svc.CreateConflictRule(context.Background(), CreateConflictRuleReq{
		ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: "forbid",
	}, "admin")

	resp, err := svc.CheckRules(context.Background(), RuleCheckReq{
		PatientID:   "P001",
		ExamItemIDs: []string{"E1", "E2"},
		OrderSource: "outpatient",
	})
	require.NoError(t, err)
	assert.True(t, resp.HasForbidden)
	assert.Len(t, resp.Conflicts, 1)
}

// === SortingStrategy Tests ===

func TestSaveSortingStrategy_Conflict(t *testing.T) {
	svc := newTestAppService()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	_, err := svc.SaveSortingStrategy(context.Background(), SaveSortingStrategyReq{
		Type: "shortest_wait",
		Scope: EffectiveScopeDTO{
			DepartmentIDs: []string{"D1"},
		},
		StartDate: start,
		EndDate:   end,
	})
	require.NoError(t, err)

	_, err = svc.SaveSortingStrategy(context.Background(), SaveSortingStrategyReq{
		Type: "nearest",
		Scope: EffectiveScopeDTO{
			DepartmentIDs: []string{"D1"},
		},
		StartDate: start,
		EndDate:   end,
	})
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleSortingConflict))
}

// === PatientAdaptRule Tests ===

func TestSaveAndListPatientAdaptRules(t *testing.T) {
	svc := newTestAppService()
	err := svc.SavePatientAdaptRules(context.Background(), []SavePatientAdaptRuleReq{
		{
			ConditionType:  "gender",
			ConditionValue: "female",
			Action:         "filter_device",
			ActionParams:   map[string]string{"device_id": "DEV1"},
			Priority:       10,
		},
	})
	require.NoError(t, err)

	list, err := svc.ListPatientAdaptRules(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "gender", list[0].ConditionType)
}

// === SourceControl Tests ===

func TestSaveAndListSourceControls(t *testing.T) {
	svc := newTestAppService()
	err := svc.SaveSourceControls(context.Background(), []SaveSourceControlReq{
		{
			SourceType:      "outpatient",
			SlotPoolID:      "POOL1",
			AllocationRatio: 0.6,
		},
	})
	require.NoError(t, err)

	list, err := svc.ListSourceControls(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "outpatient", list[0].SourceType)
	assert.Equal(t, "POOL1", list[0].SlotPoolID)
}

type mockMeta struct{}

func (m mockMeta) GetFastingItemIDs(_ context.Context, examItemIDs []string) ([]string, error) {
	// 固定返回第一个项目为空腹
	if len(examItemIDs) == 0 {
		return nil, nil
	}
	return []string{examItemIDs[0]}, nil
}

func TestCheckRules_FastingAndPools(t *testing.T) {
	svc := newTestAppService().WithExamItemMetaProvider(mockMeta{})
	_ = svc.SaveSourceControls(context.Background(), []SaveSourceControlReq{
		{
			SourceType:           "outpatient",
			SlotPoolID:           "POOL_OUT",
			AllocationRatio:      0.6,
			OverflowEnabled:      true,
			OverflowTargetPoolID: "POOL_PUBLIC",
		},
	})

	resp, err := svc.CheckRules(context.Background(), RuleCheckReq{
		PatientID:   "P001",
		ExamItemIDs: []string{"E_FAST", "E2"},
		OrderSource: "outpatient",
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"E_FAST"}, resp.FastingItems)
	assert.ElementsMatch(t, []string{"POOL_OUT", "POOL_PUBLIC"}, resp.FilteredPoolIDs)
	assert.NotEmpty(t, resp.Warnings)
}
