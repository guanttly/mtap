package resource_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/euler/mtap/internal/domain/resource"
)

// ── ExamItem ──────────────────────────────────────────────────────────────────

func TestNewExamItem(t *testing.T) {
	item := resource.NewExamItem("CT平扫", 30, false, "")
	assert.NotEmpty(t, item.ID)
	assert.Equal(t, "CT平扫", item.Name)
	assert.Equal(t, 30, item.DurationMin)
	assert.False(t, item.IsFasting)
	assert.Empty(t, item.Aliases)
}

func TestNewExamItem_Fasting(t *testing.T) {
	item := resource.NewExamItem("胃镜", 60, true, "检查前8小时禁食禁水")
	assert.True(t, item.IsFasting)
	assert.Equal(t, "检查前8小时禁食禁水", item.FastingDesc)
}

func TestExamItem_AddAlias(t *testing.T) {
	item := resource.NewExamItem("磁共振", 45, false, "")
	item.AddAlias("MRI")
	item.AddAlias("核磁共振")

	assert.Len(t, item.Aliases, 2)
	assert.Equal(t, "MRI", item.Aliases[0].Alias)
	assert.Equal(t, "核磁共振", item.Aliases[1].Alias)
	assert.Equal(t, item.ID, item.Aliases[0].ExamItemID)
}

func TestExamItem_MatchName(t *testing.T) {
	item := resource.NewExamItem("磁共振", 45, false, "")
	item.AddAlias("MRI")
	item.AddAlias("核磁")

	assert.True(t, item.MatchName("磁共振"), "应匹配正式名称")
	assert.True(t, item.MatchName("MRI"), "应匹配别名")
	assert.True(t, item.MatchName("核磁"), "应匹配别名")
	assert.False(t, item.MatchName("CT"), "不应匹配无关名称")
	assert.False(t, item.MatchName(""), "不应匹配空字符串")
}

// ── TimeSlot ──────────────────────────────────────────────────────────────────

func TestNewTimeSlot_DefaultStatus(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)

	assert.Equal(t, resource.TimeSlotAvailable, slot.Status)
	assert.Equal(t, 1, slot.Remaining)
	assert.Equal(t, 30, slot.StandardDuration)
	assert.Equal(t, 30, slot.AdjustedDuration)
}

func TestTimeSlot_Lock(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)

	lockUntil := now.Add(10 * time.Minute)
	ok := slot.Lock("patient-1", lockUntil)

	assert.True(t, ok)
	assert.Equal(t, resource.TimeSlotLocked, slot.Status)
	assert.Equal(t, "patient-1", slot.LockedBy)
	assert.Equal(t, 0, slot.Remaining)
}

func TestTimeSlot_Lock_AlreadyLocked(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)
	slot.Lock("patient-1", now.Add(10*time.Minute))

	// 第二次锁定应失败
	ok := slot.Lock("patient-2", now.Add(10*time.Minute))
	assert.False(t, ok)
}

func TestTimeSlot_Release(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)
	slot.Lock("patient-1", now.Add(10*time.Minute))

	slot.Release()

	assert.Equal(t, resource.TimeSlotAvailable, slot.Status)
	assert.Empty(t, slot.LockedBy)
	assert.Nil(t, slot.LockUntil)
	assert.Equal(t, 1, slot.Remaining)
}

func TestTimeSlot_Book(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)
	slot.Lock("patient-1", now.Add(10*time.Minute))
	slot.Book()

	assert.Equal(t, resource.TimeSlotBooked, slot.Status)
}

func TestTimeSlot_IsExpiredLock(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)

	// 锁定并设置一个过去的 LockUntil
	pastTime := now.Add(-1 * time.Minute)
	slot.Lock("patient-1", pastTime)

	assert.True(t, slot.IsExpiredLock())
}

func TestTimeSlot_IsExpiredLock_NotExpired(t *testing.T) {
	now := time.Now()
	slot := resource.NewTimeSlot("sched-1", "dev-1", "item-1", "public",
		now, now, now.Add(30*time.Minute), 30)
	futureTime := now.Add(10 * time.Minute)
	slot.Lock("patient-1", futureTime)

	assert.False(t, slot.IsExpiredLock())
}

// ── Schedule ──────────────────────────────────────────────────────────────────

func TestNewSchedule_DefaultStatus(t *testing.T) {
	sched := resource.NewSchedule("dev-1", time.Now(), "08:00", "17:00")

	assert.NotEmpty(t, sched.ID)
	assert.Equal(t, resource.ScheduleStatusNormal, sched.Status)
	assert.Equal(t, "08:00", sched.StartTime)
	assert.Equal(t, "17:00", sched.EndTime)
}

func TestSchedule_Suspend(t *testing.T) {
	sched := resource.NewSchedule("dev-1", time.Now(), "08:00", "17:00")
	sched.Suspend("设备故障")

	assert.Equal(t, resource.ScheduleStatusSuspended, sched.Status)
	assert.Equal(t, "设备故障", sched.SuspendReason)
}

func TestSchedule_SubstituteTo(t *testing.T) {
	sched := resource.NewSchedule("dev-1", time.Now(), "08:00", "17:00")
	sched.SubstituteTo("dev-2")

	assert.Equal(t, "dev-2", sched.DeviceID)
	assert.Equal(t, resource.ScheduleStatusSubstitute, sched.Status)
}

// ── SlotPool ──────────────────────────────────────────────────────────────────

func TestNewSlotPool(t *testing.T) {
	pool := resource.NewSlotPool("公共号源池", resource.SlotPoolPublic, 0.6)

	assert.NotEmpty(t, pool.ID)
	assert.Equal(t, resource.SlotPoolPublic, pool.Type)
	assert.InDelta(t, 0.6, pool.AllocationRatio, 0.001)
	assert.Equal(t, "active", pool.Status)
}

// ── Device ────────────────────────────────────────────────────────────────────

func TestNewDevice(t *testing.T) {
	types := []string{"CT", "CT增强"}
	dev := resource.NewDevice("dept-1", "campus-1", "CT机-01", "GE Revolution", "GE", 40, types)

	assert.NotEmpty(t, dev.ID)
	assert.Equal(t, resource.DeviceStatusActive, dev.Status)
	assert.Equal(t, 40, dev.MaxDailySlots)
	assert.Equal(t, types, dev.SupportedExamTypes)
}

// ── SlotGenerationService ─────────────────────────────────────────────────────

func TestSlotGenerationService_Generate(t *testing.T) {
	svc := resource.NewSlotGenerationService()
	workDate := time.Date(2026, 3, 20, 0, 0, 0, 0, time.Local)

	slots, err := svc.Generate(resource.SlotGenerationInput{
		DeviceID:    "dev-1",
		ScheduleID:  "sched-1",
		WorkDate:    workDate,
		StartTime:   "08:00",
		EndTime:     "12:00",
		SlotMinutes: 30,
		PoolType:    "public",
	})

	require.NoError(t, err)
	assert.Len(t, slots, 8) // 4h / 30min = 8 个时段
	assert.Equal(t, resource.TimeSlotAvailable, slots[0].Status)
}

func TestSlotGenerationService_Generate_MaxDailySlots(t *testing.T) {
	svc := resource.NewSlotGenerationService()
	workDate := time.Date(2026, 3, 20, 0, 0, 0, 0, time.Local)

	slots, err := svc.Generate(resource.SlotGenerationInput{
		DeviceID:      "dev-1",
		ScheduleID:    "sched-1",
		WorkDate:      workDate,
		StartTime:     "08:00",
		EndTime:       "17:00",
		SlotMinutes:   30,
		MaxDailySlots: 5,
	})

	require.NoError(t, err)
	assert.Len(t, slots, 5)
}

func TestSlotGenerationService_Generate_InvalidSlotMinutes(t *testing.T) {
	svc := resource.NewSlotGenerationService()
	_, err := svc.Generate(resource.SlotGenerationInput{
		DeviceID:    "dev-1",
		ScheduleID:  "sched-1",
		WorkDate:    time.Now(),
		StartTime:   "08:00",
		EndTime:     "12:00",
		SlotMinutes: 0, // 无效
	})
	assert.Error(t, err)
}

func TestSlotGenerationService_Generate_EndBeforeStart(t *testing.T) {
	svc := resource.NewSlotGenerationService()
	_, err := svc.Generate(resource.SlotGenerationInput{
		DeviceID:    "dev-1",
		ScheduleID:  "sched-1",
		WorkDate:    time.Now(),
		StartTime:   "17:00",
		EndTime:     "08:00", // 结束 < 开始
		SlotMinutes: 30,
	})
	assert.Error(t, err)
}

func TestSlotGenerationService_AdjustDurationForAge(t *testing.T) {
	svc := resource.NewSlotGenerationService()

	// 正常成人
	assert.Equal(t, 30, svc.AdjustDurationForAge(30, 35))
	// 儿童 +10%
	assert.Equal(t, 33, svc.AdjustDurationForAge(30, 10))
	// 老年 +15%
	assert.Equal(t, 34, svc.AdjustDurationForAge(30, 75))
}
