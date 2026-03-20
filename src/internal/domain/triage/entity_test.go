package triage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/euler/mtap/internal/domain/triage"
)

// ── CheckIn ───────────────────────────────────────────────────────────────────

func TestNewCheckIn_KioskInWindow(t *testing.T) {
	apptStart := time.Now().Add(20 * time.Minute)
	ci, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInKiosk, apptStart, "")
	require.NoError(t, err)
	assert.Equal(t, triage.CheckInKiosk, ci.Method)
	assert.False(t, ci.IsLate)
}

func TestNewCheckIn_KioskOutOfWindow_TooEarly(t *testing.T) {
	apptStart := time.Now().Add(60 * time.Minute)
	_, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInKiosk, apptStart, "")
	assert.Error(t, err, "自助机提前超过30分钟不允许签到")
}

func TestNewCheckIn_KioskOutOfWindow_TooLate(t *testing.T) {
	apptStart := time.Now().Add(-30 * time.Minute)
	_, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInKiosk, apptStart, "")
	assert.Error(t, err, "自助机超时15分钟后不允许签到")
}

func TestNewCheckIn_NurseBypassWindow(t *testing.T) {
	apptStart := time.Now().Add(60 * time.Minute)
	ci, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInNurse, apptStart, "特殊情况")
	require.NoError(t, err, "护士签到不受时间窗口限制")
	assert.Equal(t, "特殊情况", ci.Remark)
}

func TestNewCheckIn_IsLate(t *testing.T) {
	apptStart := time.Now().Add(-20 * time.Minute)
	ci, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInNurse, apptStart, "")
	require.NoError(t, err)
	assert.True(t, ci.IsLate)
}

func TestNewCheckIn_InvalidMethod(t *testing.T) {
	_, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInMethod("invalid"), time.Now(), "")
	assert.Error(t, err)
}

// ── WaitingQueue ──────────────────────────────────────────────────────────────

func makeCheckIn(t *testing.T) *triage.CheckIn {
	t.Helper()
	ci, err := triage.NewCheckIn("appt-1", "patient-1", triage.CheckInKiosk, time.Now().Add(20*time.Minute), "")
	require.NoError(t, err)
	return ci
}

func TestWaitingQueue_AddEntry(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	entry := q.AddEntry(makeCheckIn(t), "张**", "appt-1")
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, 1, entry.QueueNumber)
	assert.Equal(t, triage.EntryWaiting, entry.Status)
}

func TestWaitingQueue_AddEntry_IncrementsNumber(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	ci := makeCheckIn(t)
	q.AddEntry(ci, "张**", "appt-1")
	e2 := q.AddEntry(ci, "李**", "appt-2")
	assert.Equal(t, 2, e2.QueueNumber)
}

func TestWaitingQueue_GetWaitCount(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	ci := makeCheckIn(t)
	q.AddEntry(ci, "张**", "appt-1")
	q.AddEntry(ci, "李**", "appt-2")
	assert.Equal(t, 2, q.GetWaitCount())
}

func TestWaitingQueue_EstimateWaitTime(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	ci := makeCheckIn(t)
	q.AddEntry(ci, "张**", "appt-1")
	q.AddEntry(ci, "李**", "appt-2")
	q.AddEntry(ci, "王**", "appt-3")
	assert.Equal(t, 30, q.EstimateWaitTime())
}

func TestWaitingQueue_CallNext(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	q.AddEntry(makeCheckIn(t), "张**", "appt-1")

	entry, err := q.CallNext()
	require.NoError(t, err)
	assert.Equal(t, triage.EntryCalling, entry.Status)
	assert.Equal(t, 1, entry.CallCount)
	assert.NotNil(t, entry.CalledAt)
}

func TestWaitingQueue_CallNext_EmptyQueue(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	_, err := q.CallNext()
	assert.Error(t, err)
}

func TestWaitingQueue_Recall(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	q.AddEntry(makeCheckIn(t), "张**", "appt-1")
	q.CallNext()

	entry, err := q.Recall()
	require.NoError(t, err)
	assert.Equal(t, 2, entry.CallCount)
}

func TestWaitingQueue_Recall_ExceedLimit(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	q.AddEntry(makeCheckIn(t), "张**", "appt-1")
	q.CallNext()
	for i := 0; i < triage.MaxCallCount-1; i++ {
		q.Recall()
	}
	_, err := q.Recall()
	assert.Error(t, err)
}

func TestWaitingQueue_Recall_NoCalling(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	_, err := q.Recall()
	assert.Error(t, err)
}

func TestWaitingQueue_MissAndRequeue_FirstMiss(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	q.AddEntry(makeCheckIn(t), "张**", "appt-1")
	q.CallNext()

	entry, err := q.MissAndRequeue()
	require.NoError(t, err)
	assert.Equal(t, 1, entry.MissCount)
	assert.Equal(t, triage.EntryWaiting, entry.Status)
}

func TestWaitingQueue_MissAndRequeue_MaxMiss(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	q.AddEntry(makeCheckIn(t), "张**", "appt-1")

	var lastEntry *triage.QueueEntry
	for i := 0; i < triage.MaxMissCount; i++ {
		q.CallNext()
		e, err := q.MissAndRequeue()
		require.NoError(t, err)
		lastEntry = e
	}
	assert.Equal(t, triage.EntryNoShow, lastEntry.Status)
}

func TestWaitingQueue_MissAndRequeue_NoCalling(t *testing.T) {
	q := triage.NewWaitingQueue("room-001", "dev-1", "dept-1")
	_, err := q.MissAndRequeue()
	assert.Error(t, err)
}

// ── CheckInMethod ─────────────────────────────────────────────────────────────

func TestCheckInMethod_IsValid(t *testing.T) {
	assert.True(t, triage.CheckInKiosk.IsValid())
	assert.True(t, triage.CheckInNurse.IsValid())
	assert.True(t, triage.CheckInNFC.IsValid())
	assert.False(t, triage.CheckInMethod("unknown").IsValid())
}
