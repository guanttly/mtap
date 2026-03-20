package appointment_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/euler/mtap/internal/domain/appointment"
)

// ── NewAppointment ────────────────────────────────────────────────────────────

func TestNewAppointment_Success(t *testing.T) {
	appt, err := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	require.NoError(t, err)
	assert.NotEmpty(t, appt.ID)
	assert.Equal(t, appointment.StatusPending, appt.Status)
	assert.Equal(t, "patient-1", appt.PatientID)
}

func TestNewAppointment_EmptyPatientID(t *testing.T) {
	_, err := appointment.NewAppointment("", appointment.ModeAuto)
	assert.Error(t, err)
}

func TestNewAppointment_InvalidMode(t *testing.T) {
	_, err := appointment.NewAppointment("patient-1", appointment.AppointmentMode("invalid"))
	assert.Error(t, err)
}

// ── AddItem ───────────────────────────────────────────────────────────────────

func TestAppointment_AddItem(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	start := time.Now().Add(2 * time.Hour)
	appt.AddItem("exam-1", "slot-1", "dev-1", start, start.Add(30*time.Minute))

	assert.Len(t, appt.Items, 1)
	assert.Equal(t, "exam-1", appt.Items[0].ExamItemID)
	assert.Equal(t, appointment.ItemStatusPending, appt.Items[0].Status)
}

// ── Confirm ───────────────────────────────────────────────────────────────────

func TestAppointment_Confirm(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	err := appt.Confirm()
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusConfirmed, appt.Status)
	assert.NotNil(t, appt.ConfirmedAt)
}

func TestAppointment_Confirm_WrongStatus(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.Confirm()
	err := appt.Confirm()
	assert.Error(t, err)
}

// ── MarkPaid ──────────────────────────────────────────────────────────────────

func TestAppointment_MarkPaid(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.Confirm()
	err := appt.MarkPaid()
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusPaid, appt.Status)
	assert.True(t, appt.PaymentVerified)
}

func TestAppointment_MarkPaid_WrongStatus(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	err := appt.MarkPaid()
	assert.Error(t, err)
}

// ── Cancel ────────────────────────────────────────────────────────────────────

func TestAppointment_Cancel_Paid(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	start := time.Now().Add(3 * time.Hour)
	appt.AddItem("exam-1", "slot-1", "dev-1", start, start.Add(30*time.Minute))
	appt.Confirm()
	appt.MarkPaid()

	err := appt.Cancel("患者主动取消")
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusCancelled, appt.Status)
	assert.NotNil(t, appt.CancelledAt)
	assert.Equal(t, "患者主动取消", appt.CancelReason)
}

func TestAppointment_Cancel_TooCloseToExam(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	start := time.Now().Add(1 * time.Hour)
	appt.AddItem("exam-1", "slot-1", "dev-1", start, start.Add(30*time.Minute))
	appt.Confirm()
	appt.MarkPaid()

	err := appt.Cancel("临时取消")
	assert.Error(t, err, "距检查不足2小时不允许取消")
}

func TestAppointment_Cancel_WrongStatus(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	err := appt.Cancel("原因")
	assert.Error(t, err)
}

// ── Reschedule ────────────────────────────────────────────────────────────────

func TestAppointment_Reschedule_Success(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	start := time.Now().Add(3 * time.Hour)
	appt.AddItem("exam-1", "slot-1", "dev-1", start, start.Add(30*time.Minute))
	appt.Confirm()
	appt.MarkPaid()

	err := appt.Reschedule()
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusRescheduling, appt.Status)
}

func TestAppointment_Reschedule_LimitReached(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.ChangeCount = 3
	appt.Status = appointment.StatusPaid

	err := appt.Reschedule()
	assert.Error(t, err, "改约次数达到上限应报错")
}

func TestAppointment_Reschedule_WrongStatus(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	err := appt.Reschedule()
	assert.Error(t, err)
}

func TestAppointment_CompleteReschedule(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	start := time.Now().Add(3 * time.Hour)
	appt.AddItem("exam-1", "slot-1", "dev-1", start, start.Add(30*time.Minute))
	appt.Confirm()
	appt.MarkPaid()
	appt.Reschedule()

	appt.CompleteReschedule()
	assert.Equal(t, appointment.StatusPaid, appt.Status)
	assert.Equal(t, 1, appt.ChangeCount)
}

// ── MarkCheckedIn ─────────────────────────────────────────────────────────────

func TestAppointment_MarkCheckedIn(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.Confirm()
	appt.MarkPaid()

	err := appt.MarkCheckedIn()
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusCheckedIn, appt.Status)
}

func TestAppointment_MarkCheckedIn_WrongStatus(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.Confirm()
	err := appt.MarkCheckedIn()
	assert.Error(t, err)
}

// ── MarkNoShow ────────────────────────────────────────────────────────────────

func TestAppointment_MarkNoShow(t *testing.T) {
	appt, _ := appointment.NewAppointment("patient-1", appointment.ModeAuto)
	appt.Confirm()
	appt.MarkPaid()

	err := appt.MarkNoShow()
	require.NoError(t, err)
	assert.Equal(t, appointment.StatusNoShow, appt.Status)
}

// ── Blacklist ─────────────────────────────────────────────────────────────────

func TestNewBlacklist(t *testing.T) {
	bl := appointment.NewBlacklist("patient-1", 90)
	assert.NotEmpty(t, bl.ID)
	assert.Equal(t, appointment.BlacklistActive, bl.Status)
	assert.True(t, bl.ExpiresAt.After(time.Now()))
}

func TestBlacklist_IsExpired(t *testing.T) {
	bl := appointment.NewBlacklist("patient-1", 90)
	assert.False(t, bl.IsExpired())
}

func TestBlacklist_CanAppointOnline_Active(t *testing.T) {
	bl := appointment.NewBlacklist("patient-1", 90)
	assert.False(t, bl.CanAppointOnline())
}

func TestBlacklist_Release(t *testing.T) {
	bl := appointment.NewBlacklist("patient-1", 90)
	err := bl.Release("申诉通过")
	require.NoError(t, err)
	assert.Equal(t, appointment.BlacklistReleased, bl.Status)
	assert.NotNil(t, bl.ReleasedAt)
}

func TestBlacklist_Release_AlreadyReleased(t *testing.T) {
	bl := appointment.NewBlacklist("patient-1", 90)
	bl.Release("申诉通过")
	err := bl.Release("再次解除")
	assert.Error(t, err)
}

// ── Appeal ────────────────────────────────────────────────────────────────────

func TestNewAppeal_Success(t *testing.T) {
	appeal, err := appointment.NewAppeal("blacklist-1", "确实有急事无法赴约")
	require.NoError(t, err)
	assert.Equal(t, appointment.AppealPending, appeal.Status)
}

func TestNewAppeal_EmptyReason(t *testing.T) {
	_, err := appointment.NewAppeal("blacklist-1", "")
	assert.Error(t, err)
}

func TestNewAppeal_TooLongReason(t *testing.T) {
	longReason := make([]byte, 501)
	for i := range longReason {
		longReason[i] = 'a'
	}
	_, err := appointment.NewAppeal("blacklist-1", string(longReason))
	assert.Error(t, err)
}

func TestAppeal_Review_Approved(t *testing.T) {
	appeal, _ := appointment.NewAppeal("blacklist-1", "有效理由")
	err := appeal.Review("admin-1", true)
	require.NoError(t, err)
	assert.Equal(t, appointment.AppealApproved, appeal.Status)
	assert.Equal(t, "admin-1", appeal.ReviewedBy)
}

func TestAppeal_Review_Rejected(t *testing.T) {
	appeal, _ := appointment.NewAppeal("blacklist-1", "有效理由")
	err := appeal.Review("admin-1", false)
	require.NoError(t, err)
	assert.Equal(t, appointment.AppealRejected, appeal.Status)
}

func TestAppeal_Review_AlreadyReviewed(t *testing.T) {
	appeal, _ := appointment.NewAppeal("blacklist-1", "有效理由")
	appeal.Review("admin-1", true)
	err := appeal.Review("admin-2", false)
	assert.Error(t, err, "已审核的申诉不可重复审核")
}

// ── AppointmentMode ───────────────────────────────────────────────────────────

func TestAppointmentMode_IsValid(t *testing.T) {
	assert.True(t, appointment.ModeAuto.IsValid())
	assert.True(t, appointment.ModeCombo.IsValid())
	assert.True(t, appointment.ModeManual.IsValid())
	assert.False(t, appointment.AppointmentMode("unknown").IsValid())
}
