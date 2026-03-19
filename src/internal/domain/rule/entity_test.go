package rule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// === ConflictRule Tests ===

func TestNewConflictRule_Success(t *testing.T) {
	rule, err := NewConflictRule("EXAM_A", "EXAM_B", 48, ConflictLevelWarning, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, rule.ID)
	assert.Equal(t, "EXAM_A", rule.ItemAID)
	assert.Equal(t, "EXAM_B", rule.ItemBID)
	assert.Equal(t, 48, rule.MinInterval)
	assert.Equal(t, ConflictLevelWarning, rule.Level)
	assert.Equal(t, RuleStatusActive, rule.Status)
}

func TestNewConflictRule_SameItem(t *testing.T) {
	_, err := NewConflictRule("EXAM_A", "EXAM_A", 24, ConflictLevelForbid, "admin")
	require.Error(t, err)
	assert.True(t, bizErr.Is(err, bizErr.ErrRuleSameItem))
}

func TestNewConflictRule_IntervalOutOfRange(t *testing.T) {
	_, err := NewConflictRule("A", "B", -1, ConflictLevelWarning, "admin")
	assert.Error(t, err)

	_, err = NewConflictRule("A", "B", 721, ConflictLevelWarning, "admin")
	assert.Error(t, err)
}

func TestNewConflictRule_InvalidLevel(t *testing.T) {
	_, err := NewConflictRule("A", "B", 24, ConflictLevel("invalid"), "admin")
	assert.Error(t, err)
}

func TestNewConflictRule_NormalizesOrder(t *testing.T) {
	rule, _ := NewConflictRule("Z_ITEM", "A_ITEM", 10, ConflictLevelWarning, "admin")
	assert.Equal(t, "A_ITEM", rule.ItemAID)
	assert.Equal(t, "Z_ITEM", rule.ItemBID)
}

func TestConflictRule_Validate(t *testing.T) {
	rule := &ConflictRule{ItemAID: "A", ItemBID: "B", MinInterval: 24, Level: ConflictLevelForbid}
	assert.NoError(t, rule.Validate())

	rule.ItemBID = "A"
	assert.Error(t, rule.Validate())
}

// === ConflictPackage Tests ===

func TestNewConflictPackage_Success(t *testing.T) {
	pkg, err := NewConflictPackage("造影剂互斥", []string{"E1", "E2", "E3"}, 24, ConflictLevelForbid)
	require.NoError(t, err)
	assert.Equal(t, "造影剂互斥", pkg.Name)
	assert.Equal(t, 3, len(pkg.Items))
	assert.True(t, pkg.IsValid())
}

func TestNewConflictPackage_TooFewItems(t *testing.T) {
	_, err := NewConflictPackage("test", []string{"E1"}, 24, ConflictLevelForbid)
	assert.True(t, bizErr.Is(err, bizErr.ErrRulePkgTooFew))
}

func TestNewConflictPackage_EmptyName(t *testing.T) {
	_, err := NewConflictPackage("", []string{"E1", "E2"}, 24, ConflictLevelForbid)
	assert.Error(t, err)
}

func TestConflictPackage_AddItem(t *testing.T) {
	pkg, _ := NewConflictPackage("test", []string{"E1", "E2"}, 24, ConflictLevelForbid)
	err := pkg.AddItem("E3")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(pkg.Items))
}

func TestConflictPackage_AddItem_Duplicate(t *testing.T) {
	pkg, _ := NewConflictPackage("test", []string{"E1", "E2"}, 24, ConflictLevelForbid)
	err := pkg.AddItem("E1")
	assert.Error(t, err)
}

func TestConflictPackage_RemoveItem(t *testing.T) {
	pkg, _ := NewConflictPackage("test", []string{"E1", "E2", "E3"}, 24, ConflictLevelForbid)
	err := pkg.RemoveItem("E2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(pkg.Items))
	assert.True(t, pkg.IsValid())
}

func TestConflictPackage_RemoveItem_NotFound(t *testing.T) {
	pkg, _ := NewConflictPackage("test", []string{"E1", "E2"}, 24, ConflictLevelForbid)
	err := pkg.RemoveItem("E99")
	assert.Error(t, err)
}

func TestConflictPackage_IsValid_AfterRemove(t *testing.T) {
	pkg, _ := NewConflictPackage("test", []string{"E1", "E2"}, 24, ConflictLevelForbid)
	_ = pkg.RemoveItem("E1")
	assert.False(t, pkg.IsValid())
}

// === DependencyRule Tests ===

func TestNewDependencyRule_Success(t *testing.T) {
	rule, err := NewDependencyRule("PRE", "POST", DependencyTypeMandatory, 72)
	require.NoError(t, err)
	assert.Equal(t, "PRE", rule.PreItemID)
	assert.Equal(t, "POST", rule.PostItemID)
	assert.Equal(t, DependencyTypeMandatory, rule.Type)
	assert.Equal(t, 72, rule.ValidityHours)
}

func TestNewDependencyRule_SameItem(t *testing.T) {
	_, err := NewDependencyRule("A", "A", DependencyTypeMandatory, 72)
	assert.Error(t, err)
}

func TestNewDependencyRule_InvalidValidity(t *testing.T) {
	_, err := NewDependencyRule("A", "B", DependencyTypeMandatory, 0)
	assert.Error(t, err)

	_, err = NewDependencyRule("A", "B", DependencyTypeMandatory, -1)
	assert.Error(t, err)
}

func TestNewDependencyRule_InvalidType(t *testing.T) {
	_, err := NewDependencyRule("A", "B", DependencyType("invalid"), 72)
	assert.Error(t, err)
}

// === PriorityTag Tests ===

func TestNewPriorityTag_Success(t *testing.T) {
	tag, err := NewPriorityTag("急诊", 90, "#FF4D4F")
	require.NoError(t, err)
	assert.Equal(t, "急诊", tag.Name)
	assert.Equal(t, 90, tag.Weight)
	assert.False(t, tag.IsPreset)
	assert.True(t, tag.CanDelete())
}

func TestNewPriorityTag_InvalidWeight(t *testing.T) {
	_, err := NewPriorityTag("test", 0, "#000")
	assert.Error(t, err)

	_, err = NewPriorityTag("test", 101, "#000")
	assert.Error(t, err)
}

func TestPriorityTag_PresetCannotDelete(t *testing.T) {
	tag := &PriorityTag{IsPreset: true}
	assert.False(t, tag.CanDelete())
}

// === SortingStrategy Tests ===

func TestNewSortingStrategy_Success(t *testing.T) {
	start := time.Now()
	end := start.Add(30 * 24 * time.Hour)
	s, err := NewSortingStrategy(SortingTypeShortestWait, EffectiveScope{}, start, end)
	require.NoError(t, err)
	assert.True(t, s.IsEffective(start.Add(24*time.Hour)))
}

func TestNewSortingStrategy_InvalidDateRange(t *testing.T) {
	now := time.Now()
	_, err := NewSortingStrategy(SortingTypeShortestWait, EffectiveScope{}, now, now)
	assert.Error(t, err)

	_, err = NewSortingStrategy(SortingTypeShortestWait, EffectiveScope{}, now, now.Add(-1*time.Hour))
	assert.Error(t, err)
}

func TestSortingStrategy_IsEffective(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	s, _ := NewSortingStrategy(SortingTypeNearest, EffectiveScope{}, start, end)

	assert.True(t, s.IsEffective(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)))
	assert.False(t, s.IsEffective(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)))
	assert.False(t, s.IsEffective(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))

	s.Status = RuleStatusInactive
	assert.False(t, s.IsEffective(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)))
}

// === Value Object Tests ===

func TestConflictLevel_IsValid(t *testing.T) {
	assert.True(t, ConflictLevelForbid.IsValid())
	assert.True(t, ConflictLevelWarning.IsValid())
	assert.False(t, ConflictLevel("invalid").IsValid())
}

func TestDependencyType_IsValid(t *testing.T) {
	assert.True(t, DependencyTypeMandatory.IsValid())
	assert.True(t, DependencyTypeRecommended.IsValid())
	assert.False(t, DependencyType("invalid").IsValid())
}

func TestEffectiveScope_IsEmpty(t *testing.T) {
	assert.True(t, EffectiveScope{}.IsEmpty())
	assert.False(t, EffectiveScope{CampusIDs: []string{"C1"}}.IsEmpty())
}
