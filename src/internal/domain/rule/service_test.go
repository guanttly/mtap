package rule

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === Mock Repositories ===

type mockConflictRuleRepo struct {
	rules []*ConflictRule
}

func (m *mockConflictRuleRepo) Save(_ context.Context, rule *ConflictRule) error {
	m.rules = append(m.rules, rule)
	return nil
}
func (m *mockConflictRuleRepo) FindByID(_ context.Context, id string) (*ConflictRule, error) {
	for _, r := range m.rules {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, nil
}
func (m *mockConflictRuleRepo) FindByItemPair(_ context.Context, a, b string) (*ConflictRule, error) {
	key := pairKey(a, b)
	for _, r := range m.rules {
		if pairKey(r.ItemAID, r.ItemBID) == key {
			return r, nil
		}
	}
	return nil, nil
}
func (m *mockConflictRuleRepo) FindAll(_ context.Context, status RuleStatus) ([]*ConflictRule, error) {
	var result []*ConflictRule
	for _, r := range m.rules {
		if r.Status == status {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockConflictRuleRepo) Update(_ context.Context, rule *ConflictRule) error { return nil }
func (m *mockConflictRuleRepo) Delete(_ context.Context, id string) error          { return nil }

type mockConflictPkgRepo struct {
	pkgs []*ConflictPackage
}

func (m *mockConflictPkgRepo) Save(_ context.Context, pkg *ConflictPackage) error {
	m.pkgs = append(m.pkgs, pkg)
	return nil
}
func (m *mockConflictPkgRepo) FindByID(_ context.Context, id string) (*ConflictPackage, error) {
	return nil, nil
}
func (m *mockConflictPkgRepo) FindByName(_ context.Context, name string) (*ConflictPackage, error) {
	return nil, nil
}
func (m *mockConflictPkgRepo) FindAll(_ context.Context) ([]*ConflictPackage, error) {
	return m.pkgs, nil
}
func (m *mockConflictPkgRepo) Update(_ context.Context, pkg *ConflictPackage) error { return nil }
func (m *mockConflictPkgRepo) Delete(_ context.Context, id string) error             { return nil }

type mockDependencyRuleRepo struct {
	rules []*DependencyRule
}

func (m *mockDependencyRuleRepo) Save(_ context.Context, rule *DependencyRule) error {
	m.rules = append(m.rules, rule)
	return nil
}
func (m *mockDependencyRuleRepo) FindByID(_ context.Context, id string) (*DependencyRule, error) {
	return nil, nil
}
func (m *mockDependencyRuleRepo) FindByPostItem(_ context.Context, postItemID string) ([]*DependencyRule, error) {
	var result []*DependencyRule
	for _, r := range m.rules {
		if r.PostItemID == postItemID {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockDependencyRuleRepo) FindAll(_ context.Context, status RuleStatus) ([]*DependencyRule, error) {
	var result []*DependencyRule
	for _, r := range m.rules {
		if r.Status == status {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockDependencyRuleRepo) Update(_ context.Context, rule *DependencyRule) error { return nil }
func (m *mockDependencyRuleRepo) Delete(_ context.Context, id string) error             { return nil }

type mockPreItemChecker struct {
	records map[string]*time.Time // key = patientID + "|" + examItemID
}

func (m *mockPreItemChecker) GetCompletedTime(_ context.Context, patientID, examItemID string) (*time.Time, error) {
	key := patientID + "|" + examItemID
	if t, ok := m.records[key]; ok {
		return t, nil
	}
	return nil, nil
}

// === ConflictDetectionService Tests ===

func TestDetectConflicts_NoConflict(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{}
	pkgRepo := &mockConflictPkgRepo{}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1", "E2"}, nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestDetectConflicts_SingleItem(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{}
	pkgRepo := &mockConflictPkgRepo{}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1"}, nil)
	require.NoError(t, err)
	assert.Nil(t, results)
}

func TestDetectConflicts_DirectConflict(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{
		rules: []*ConflictRule{
			{ID: "R1", ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: ConflictLevelForbid, Status: RuleStatusActive},
		},
	}
	pkgRepo := &mockConflictPkgRepo{}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1", "E2"}, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, ConflictLevelForbid, results[0].Level)
	assert.Equal(t, -1, results[0].ActualInterval) // 同次预约
}

func TestDetectConflicts_PackageConflict(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{}
	pkgRepo := &mockConflictPkgRepo{
		pkgs: []*ConflictPackage{
			{
				ID: "PKG1", Name: "造影剂互斥", MinInterval: 24, Level: ConflictLevelForbid, Status: RuleStatusActive,
				Items: []ConflictPackageItem{
					{ExamItemID: "E1"}, {ExamItemID: "E2"}, {ExamItemID: "E3"},
				},
			},
		},
	}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1", "E3"}, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, ConflictLevelForbid, results[0].Level)
}

func TestDetectConflicts_StricterRuleWins(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{
		rules: []*ConflictRule{
			{ID: "R1", ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: ConflictLevelWarning, Status: RuleStatusActive},
		},
	}
	pkgRepo := &mockConflictPkgRepo{
		pkgs: []*ConflictPackage{
			{
				ID: "PKG1", MinInterval: 48, Level: ConflictLevelForbid, Status: RuleStatusActive,
				Items: []ConflictPackageItem{
					{ExamItemID: "E1"}, {ExamItemID: "E2"},
				},
			},
		},
	}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1", "E2"}, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, ConflictLevelForbid, results[0].Level) // forbid > warning
}

func TestDetectConflicts_InactiveRuleIgnored(t *testing.T) {
	ruleRepo := &mockConflictRuleRepo{
		rules: []*ConflictRule{
			{ID: "R1", ItemAID: "E1", ItemBID: "E2", MinInterval: 24, Level: ConflictLevelForbid, Status: RuleStatusInactive},
		},
	}
	pkgRepo := &mockConflictPkgRepo{}
	svc := NewConflictDetectionService(ruleRepo, pkgRepo)

	results, err := svc.DetectConflicts(context.Background(), []string{"E1", "E2"}, nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// === DependencyValidationService Tests ===

func TestValidateDependencies_Passed(t *testing.T) {
	completedAt := time.Now().Add(-24 * time.Hour)
	depRepo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{ID: "D1", PreItemID: "BLOOD", PostItemID: "GASTRO", Type: DependencyTypeMandatory, ValidityHours: 72, Status: RuleStatusActive},
		},
	}
	checker := &mockPreItemChecker{
		records: map[string]*time.Time{"P001|BLOOD": &completedAt},
	}
	svc := NewDependencyValidationService(depRepo, checker)

	results, err := svc.ValidateDependencies(context.Background(), []string{"GASTRO"}, "P001")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, DependencyPassed, results[0].Status)
}

func TestValidateDependencies_Blocked(t *testing.T) {
	depRepo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{ID: "D1", PreItemID: "BLOOD", PostItemID: "GASTRO", Type: DependencyTypeMandatory, ValidityHours: 72, Status: RuleStatusActive},
		},
	}
	checker := &mockPreItemChecker{records: map[string]*time.Time{}}
	svc := NewDependencyValidationService(depRepo, checker)

	results, err := svc.ValidateDependencies(context.Background(), []string{"GASTRO"}, "P001")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, DependencyBlocked, results[0].Status)
}

func TestValidateDependencies_Expired(t *testing.T) {
	completedAt := time.Now().Add(-100 * time.Hour) // 超过72小时
	depRepo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{ID: "D1", PreItemID: "BLOOD", PostItemID: "GASTRO", Type: DependencyTypeMandatory, ValidityHours: 72, Status: RuleStatusActive},
		},
	}
	checker := &mockPreItemChecker{
		records: map[string]*time.Time{"P001|BLOOD": &completedAt},
	}
	svc := NewDependencyValidationService(depRepo, checker)

	results, err := svc.ValidateDependencies(context.Background(), []string{"GASTRO"}, "P001")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, DependencyExpired, results[0].Status)
}

func TestValidateDependencies_NoDeps(t *testing.T) {
	depRepo := &mockDependencyRuleRepo{}
	checker := &mockPreItemChecker{}
	svc := NewDependencyValidationService(depRepo, checker)

	results, err := svc.ValidateDependencies(context.Background(), []string{"NO_DEP_ITEM"}, "P001")
	require.NoError(t, err)
	assert.Empty(t, results)
}

// === CircularDependencyChecker Tests ===

func TestCircularDependency_NoCycle(t *testing.T) {
	repo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{PreItemID: "A", PostItemID: "B", Status: RuleStatusActive},
			{PreItemID: "B", PostItemID: "C", Status: RuleStatusActive},
		},
	}
	checker := NewCircularDependencyChecker(repo)

	hasCycle, err := checker.HasCircularDependency(context.Background(), "C", "D")
	require.NoError(t, err)
	assert.False(t, hasCycle)
}

func TestCircularDependency_DirectCycle(t *testing.T) {
	repo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{PreItemID: "A", PostItemID: "B", Status: RuleStatusActive},
		},
	}
	checker := NewCircularDependencyChecker(repo)

	// Adding B→A creates A→B→A cycle
	hasCycle, err := checker.HasCircularDependency(context.Background(), "B", "A")
	require.NoError(t, err)
	assert.True(t, hasCycle)
}

func TestCircularDependency_IndirectCycle(t *testing.T) {
	repo := &mockDependencyRuleRepo{
		rules: []*DependencyRule{
			{PreItemID: "A", PostItemID: "B", Status: RuleStatusActive},
			{PreItemID: "B", PostItemID: "C", Status: RuleStatusActive},
		},
	}
	checker := NewCircularDependencyChecker(repo)

	// Adding C→A creates A→B→C→A cycle
	hasCycle, err := checker.HasCircularDependency(context.Background(), "C", "A")
	require.NoError(t, err)
	assert.True(t, hasCycle)
}

// === Helper Tests ===

func TestPairKey_Normalized(t *testing.T) {
	assert.Equal(t, pairKey("A", "B"), pairKey("B", "A"))
}

func TestIsStricter(t *testing.T) {
	forbid := &ConflictRule{Level: ConflictLevelForbid, MinInterval: 24}
	warning := &ConflictRule{Level: ConflictLevelWarning, MinInterval: 48}
	warningShort := &ConflictRule{Level: ConflictLevelWarning, MinInterval: 12}

	assert.True(t, isStricter(forbid, warning))
	assert.False(t, isStricter(warning, forbid))
	assert.True(t, isStricter(warning, warningShort)) // same level, longer interval is stricter
}
