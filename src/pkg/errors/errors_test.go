package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(ErrNotFound)
	assert.Equal(t, ErrNotFound, err.Code)
	assert.Equal(t, "资源不存在", err.Message)
	assert.Empty(t, err.Detail)
}

func TestNewWithDetail(t *testing.T) {
	err := NewWithDetail(ErrApptSlotTaken, "slot-123")
	assert.Equal(t, ErrApptSlotTaken, err.Code)
	assert.Contains(t, err.Error(), "slot-123")
	assert.Contains(t, err.Error(), "号源已被抢占")
}

func TestWrap(t *testing.T) {
	origErr := assert.AnError
	err := Wrap(ErrInternal, origErr)
	assert.Equal(t, ErrInternal, err.Code)
	assert.Equal(t, origErr.Error(), err.Detail)
}

func TestError_String(t *testing.T) {
	tests := []struct {
		name     string
		err      *BizError
		contains []string
	}{
		{
			name:     "without detail",
			err:      New(ErrUnauthorized),
			contains: []string{"1001", "未认证"},
		},
		{
			name:     "with detail",
			err:      NewWithDetail(ErrRuleSameItem, "item-a == item-b"),
			contains: []string{"2001", "不能相同", "item-a == item-b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.err.Error()
			for _, c := range tt.contains {
				assert.Contains(t, s, c)
			}
		})
	}
}

func TestIs(t *testing.T) {
	err := New(ErrForbidden)
	assert.True(t, Is(err, ErrForbidden))
	assert.False(t, Is(err, ErrNotFound))
	assert.False(t, Is(assert.AnError, ErrNotFound))
}

func TestMessageOf_Unknown(t *testing.T) {
	msg := MessageOf(Code(99999))
	assert.Equal(t, "未知错误", msg)
}

func TestAllCodesHaveMessages(t *testing.T) {
	codes := []Code{
		OK, ErrUnauthorized, ErrForbidden, ErrNotFound, ErrInvalidParam,
		ErrInternal, ErrRateLimit, ErrDuplicate, ErrConflict, ErrTimeout,
		ErrRuleSameItem, ErrRuleDuplicate, ErrRulePkgNameDup, ErrRulePkgTooFew,
		ErrRuleCircularDep, ErrRulePresetNoDelete, ErrRuleTagNameDup,
		ErrRuleSortingConflict, ErrRuleInvalidScope, ErrRuleServiceDown,
		ErrResDeviceNotFound, ErrResScheduleConflict, ErrResSlotLockFail,
		ErrResSlotReleaseFail, ErrResSubstituteIncompat, ErrResExtraSlotOverlap,
		ErrResAliasConflict, ErrResAliasTooMany, ErrResSyncFailed, ErrResSlotOverLimit,
		ErrApptNotPaid, ErrApptConflictForbid, ErrApptDepBlocked, ErrApptSlotTaken,
		ErrApptConfirmTimeout, ErrApptChangeLimitReached, ErrApptTooCloseToExam,
		ErrApptBlacklisted, ErrApptManualForbidden, ErrApptComboTooMany, ErrApptPayTimeout,
		ErrTriageNotFound, ErrTriageOutOfWindow, ErrTriageAlreadyCheckedIn,
		ErrTriageQueueEmpty, ErrTriageRecallLimit, ErrTriageStatusInvalid,
		ErrTriageUndoExpired, ErrTriageInvalidQR,
		ErrOptStrategyLimit, ErrOptRejectReasonReq, ErrOptStatusInvalid,
		ErrOptJointNotComplete, ErrOptTrialActive, ErrOptCooldown,
		ErrOptEmergencyRollback, ErrOptEvalFailed, ErrOptCTypeNoExec, ErrOptCostOverrun,
	}
	for _, c := range codes {
		msg := MessageOf(c)
		assert.NotEqual(t, "未知错误", msg, "code %d missing message", c)
	}
}
