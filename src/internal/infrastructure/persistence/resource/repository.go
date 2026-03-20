package resource

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/euler/mtap/internal/application/resource"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

type Repositories struct {
	DB *gorm.DB
}

func NewRepositories(db *gorm.DB) *Repositories { return &Repositories{DB: db} }

type DeviceRepo struct{ db *gorm.DB }
type ExamItemRepo struct{ db *gorm.DB }
type AliasRepo struct{ db *gorm.DB }
type SlotPoolRepo struct{ db *gorm.DB }
type ScheduleRepo struct{ db *gorm.DB }
type TimeSlotRepo struct{ db *gorm.DB }
type CampusRepo struct{ db *gorm.DB }
type DepartmentRepo struct{ db *gorm.DB }

func (r *Repositories) DeviceRepo() *DeviceRepo         { return &DeviceRepo{db: r.DB} }
func (r *Repositories) ExamItemRepo() *ExamItemRepo     { return &ExamItemRepo{db: r.DB} }
func (r *Repositories) AliasRepo() *AliasRepo           { return &AliasRepo{db: r.DB} }
func (r *Repositories) SlotPoolRepo() *SlotPoolRepo     { return &SlotPoolRepo{db: r.DB} }
func (r *Repositories) ScheduleRepo() *ScheduleRepo     { return &ScheduleRepo{db: r.DB} }
func (r *Repositories) TimeSlotRepo() *TimeSlotRepo     { return &TimeSlotRepo{db: r.DB} }
func (r *Repositories) CampusRepo() *CampusRepo         { return &CampusRepo{db: r.DB} }
func (r *Repositories) DepartmentRepo() *DepartmentRepo { return &DepartmentRepo{db: r.DB} }

func (r *DeviceRepo) Create(ctx context.Context, d resource.DeviceResp) error {
	return r.db.WithContext(ctx).Create(&po.DevicePO{
		ID:           d.ID,
		Name:         d.Name,
		CampusID:     d.CampusID,
		DepartmentID: d.DepartmentID,
		Status:       d.Status,
		CreatedAt:    d.CreatedAt,
	}).Error
}
func (r *DeviceRepo) Get(ctx context.Context, id string) (*resource.DeviceResp, error) {
	var p po.DevicePO
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &resource.DeviceResp{
		ID:           p.ID,
		Name:         p.Name,
		CampusID:     p.CampusID,
		DepartmentID: p.DepartmentID,
		Status:       p.Status,
		CreatedAt:    p.CreatedAt,
	}, nil
}
func (r *DeviceRepo) List(ctx context.Context) ([]resource.DeviceResp, error) {
	var ps []po.DevicePO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.DeviceResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.DeviceResp{
			ID:           p.ID,
			Name:         p.Name,
			CampusID:     p.CampusID,
			DepartmentID: p.DepartmentID,
			Status:       p.Status,
			CreatedAt:    p.CreatedAt,
		})
	}
	return out, nil
}

func (r *DeviceRepo) Update(ctx context.Context, id string, d resource.DeviceResp) error {
	return r.db.WithContext(ctx).Model(&po.DevicePO{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":          d.Name,
		"campus_id":     d.CampusID,
		"department_id": d.DepartmentID,
		"status":        d.Status,
	}).Error
}

func (r *DeviceRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&po.DevicePO{}).Error
}

func (r *ExamItemRepo) Create(ctx context.Context, e resource.ExamItemResp) error {
	return r.db.WithContext(ctx).Create(&po.ExamItemPO{
		ID:          e.ID,
		Name:        e.Name,
		DurationMin: e.DurationMin,
		IsFasting:   e.IsFasting,
		FastingDesc: e.FastingDesc,
	}).Error
}
func (r *ExamItemRepo) Get(ctx context.Context, id string) (*resource.ExamItemResp, error) {
	var p po.ExamItemPO
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &resource.ExamItemResp{
		ID:          p.ID,
		Name:        p.Name,
		DurationMin: p.DurationMin,
		IsFasting:   p.IsFasting,
		FastingDesc: p.FastingDesc,
	}, nil
}
func (r *ExamItemRepo) List(ctx context.Context) ([]resource.ExamItemResp, error) {
	var ps []po.ExamItemPO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.ExamItemResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.ExamItemResp{
			ID:          p.ID,
			Name:        p.Name,
			DurationMin: p.DurationMin,
			IsFasting:   p.IsFasting,
			FastingDesc: p.FastingDesc,
		})
	}
	return out, nil
}

func (r *ExamItemRepo) Update(ctx context.Context, id string, e resource.ExamItemResp) error {
	return r.db.WithContext(ctx).Model(&po.ExamItemPO{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":         e.Name,
		"duration_min": e.DurationMin,
		"is_fasting":   e.IsFasting,
		"fasting_desc": e.FastingDesc,
	}).Error
}

func (r *ExamItemRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&po.ExamItemPO{}, "id = ?", id).Error
}
func (r *ExamItemRepo) ListFastingIDs(ctx context.Context, ids []string) ([]string, error) {
	var ps []po.ExamItemPO
	if err := r.db.WithContext(ctx).Where("id IN ? AND is_fasting = ?", ids, true).Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.ID)
	}
	return out, nil
}

func (r *ExamItemRepo) GetDurationMin(ctx context.Context, id string) (int, error) {
	var p po.ExamItemPO
	if err := r.db.WithContext(ctx).Select("duration_min").First(&p, "id = ?", id).Error; err != nil {
		return 0, err
	}
	return p.DurationMin, nil
}

// GetFastingItemIDs implements rule.ExamItemMetaProvider (跨模块端口)。
func (r *ExamItemRepo) GetFastingItemIDs(ctx context.Context, examItemIDs []string) ([]string, error) {
	return r.ListFastingIDs(ctx, examItemIDs)
}

func (r *AliasRepo) Create(ctx context.Context, a resource.AliasResp) error {
	return r.db.WithContext(ctx).Create(&po.ItemAliasPO{
		ID:         a.ID,
		ExamItemID: a.ExamItemID,
		Alias:      a.Alias,
	}).Error
}
func (r *AliasRepo) List(ctx context.Context, examItemID string) ([]resource.AliasResp, error) {
	var ps []po.ItemAliasPO
	q := r.db.WithContext(ctx)
	if examItemID != "" {
		q = q.Where("exam_item_id = ?", examItemID)
	}
	if err := q.Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.AliasResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.AliasResp{
			ID:         p.ID,
			ExamItemID: p.ExamItemID,
			Alias:      p.Alias,
		})
	}
	return out, nil
}

func (r *SlotPoolRepo) Create(ctx context.Context, p resource.SlotPoolResp) error {
	return r.db.WithContext(ctx).Create(&po.SlotPoolPO{
		ID:     p.ID,
		Name:   p.Name,
		Type:   p.Type,
		Status: p.Status,
	}).Error
}
func (r *SlotPoolRepo) List(ctx context.Context) ([]resource.SlotPoolResp, error) {
	var ps []po.SlotPoolPO
	if err := r.db.WithContext(ctx).Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.SlotPoolResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.SlotPoolResp{
			ID:     p.ID,
			Name:   p.Name,
			Type:   p.Type,
			Status: p.Status,
		})
	}
	return out, nil
}

func (r *ScheduleRepo) Create(ctx context.Context, deviceID string, date time.Time, startTime, endTime string) (string, error) {
	id := uuid.New().String()
	err := r.db.WithContext(ctx).Create(&po.SchedulePO{
		ID:        id,
		DeviceID:  deviceID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}).Error
	if err != nil {
		// sqlite unique constraint / 其他 db 的唯一约束错误统一映射
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			return "", bizErr.New(bizErr.ErrResScheduleConflict)
		}
	}
	return id, err
}

func (r *ScheduleRepo) Suspend(ctx context.Context, deviceID string, date time.Time, reason string) error {
	return r.db.WithContext(ctx).Model(&po.SchedulePO{}).
		Where("device_id = ? AND date = ?", deviceID, date).
		Updates(map[string]interface{}{
			"status":         "suspended",
			"suspend_reason": reason,
			"updated_at":     time.Now(),
		}).Error
}

func (r *ScheduleRepo) Substitute(ctx context.Context, sourceDeviceID, targetDeviceID string, date time.Time) error {
	return r.db.WithContext(ctx).Model(&po.SchedulePO{}).
		Where("device_id = ? AND date = ?", sourceDeviceID, date).
		Updates(map[string]interface{}{
			"device_id":  targetDeviceID,
			"updated_at": time.Now(),
		}).Error
}

func (r *TimeSlotRepo) BulkCreate(ctx context.Context, slots []resource.TimeSlotResp) error {
	ps := make([]po.TimeSlotPO, 0, len(slots))
	for _, s := range slots {
		ps = append(ps, po.TimeSlotPO{
			ID:               s.ID,
			DeviceID:         s.DeviceID,
			Date:             time.Date(s.StartAt.Year(), s.StartAt.Month(), s.StartAt.Day(), 0, 0, 0, 0, s.StartAt.Location()),
			ExamItemID:       s.ExamItemID,
			PoolType:         s.PoolType,
			StartAt:          s.StartAt,
			EndAt:            s.EndAt,
			StandardDuration: s.StandardDuration,
			AdjustedDuration: s.AdjustedDuration,
			Status:           s.Status,
			LockedBy:         s.LockedBy,
			LockUntil:        s.LockUntil,
			Remaining:        s.Remaining,
		})
	}
	if len(ps) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&ps).Error
}

func (r *TimeSlotRepo) ListByDeviceAndDate(ctx context.Context, deviceID string, date time.Time) ([]resource.TimeSlotResp, error) {
	var ps []po.TimeSlotPO
	if err := r.db.WithContext(ctx).Where("device_id = ? AND date = ?", deviceID, date).Order("start_at asc").Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.TimeSlotResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.TimeSlotResp{
			ID:               p.ID,
			DeviceID:         p.DeviceID,
			ExamItemID:       p.ExamItemID,
			PoolType:         p.PoolType,
			StartAt:          p.StartAt,
			EndAt:            p.EndAt,
			Status:           p.Status,
			StandardDuration: p.StandardDuration,
			AdjustedDuration: p.AdjustedDuration,
			LockedBy:         p.LockedBy,
			LockUntil:        p.LockUntil,
			Remaining:        p.Remaining,
		})
	}
	return out, nil
}

func (r *TimeSlotRepo) QueryAvailable(ctx context.Context, deviceID string, date time.Time, examItemID, poolType string) ([]resource.TimeSlotResp, error) {
	now := time.Now()
	q := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).
		Where("device_id = ? AND date = ?", deviceID, date).
		Where("status = ?", "available").
		Or("status = ? AND lock_until IS NOT NULL AND lock_until < ?", "locked", now)
	if examItemID != "" {
		q = q.Where("exam_item_id = ?", examItemID)
	}
	if poolType != "" {
		q = q.Where("pool_type = ?", poolType)
	}
	var ps []po.TimeSlotPO
	if err := q.Order("start_at asc").Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.TimeSlotResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.TimeSlotResp{
			ID:               p.ID,
			DeviceID:         p.DeviceID,
			ExamItemID:       p.ExamItemID,
			PoolType:         p.PoolType,
			StartAt:          p.StartAt,
			EndAt:            p.EndAt,
			Status:           p.Status,
			StandardDuration: p.StandardDuration,
			AdjustedDuration: p.AdjustedDuration,
			LockedBy:         p.LockedBy,
			LockUntil:        p.LockUntil,
			Remaining:        p.Remaining,
		})
	}
	return out, nil
}

func (r *TimeSlotRepo) Lock(ctx context.Context, slotID string, patientID string, lockUntil time.Time) error {
	now := time.Now()
	// 原子更新：只有 available 或已过期锁 才能锁定
	tx := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).
		Where("id = ?", slotID).
		Where("(status = ?) OR (status = ? AND lock_until IS NOT NULL AND lock_until < ?)", "available", "locked", now).
		Updates(map[string]interface{}{
			"status":     "locked",
			"locked_by":  patientID,
			"lock_until": lockUntil,
			"updated_at": now,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return bizErr.New(bizErr.ErrResSlotLockFail)
	}
	return nil
}

func (r *TimeSlotRepo) Release(ctx context.Context, slotID string, patientID string, allowForce bool) error {
	now := time.Now()
	q := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).Where("id = ?", slotID).Where("status = ?", "locked")
	if !allowForce {
		q = q.Where("locked_by = ?", patientID)
	}
	tx := q.Updates(map[string]interface{}{
		"status":     "available",
		"locked_by":  "",
		"lock_until": nil,
		"updated_at": now,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return bizErr.New(bizErr.ErrResSlotReleaseFail)
	}
	return nil
}

func (r *TimeSlotRepo) SuspendRange(ctx context.Context, deviceID string, date time.Time, startAt, endAt time.Time, reason string) (int64, error) {
	// 将范围内 slots 标记为 suspended，并清除锁信息
	tx := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).
		Where("device_id = ? AND date = ?", deviceID, date).
		Where("start_at >= ? AND end_at <= ?", startAt, endAt).
		Where("status IN ?", []string{"available", "locked"}).
		Updates(map[string]interface{}{
			"status":     "suspended",
			"locked_by":  "",
			"lock_until": nil,
			"updated_at": time.Now(),
		})
	return tx.RowsAffected, tx.Error
}

func (r *TimeSlotRepo) UpdateDeviceByDate(ctx context.Context, sourceDeviceID, targetDeviceID string, date time.Time) (int64, error) {
	tx := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).
		Where("device_id = ? AND date = ?", sourceDeviceID, date).
		Updates(map[string]interface{}{
			"device_id":  targetDeviceID,
			"updated_at": time.Now(),
		})
	return tx.RowsAffected, tx.Error
}

func (r *TimeSlotRepo) HasOverlap(ctx context.Context, deviceID string, date time.Time, startAt, endAt time.Time) (bool, error) {
	var count int64
	// overlap condition: existing.start < end && existing.end > start
	if err := r.db.WithContext(ctx).Model(&po.TimeSlotPO{}).
		Where("device_id = ? AND date = ?", deviceID, date).
		Where("start_at < ? AND end_at > ?", endAt, startAt).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// === AliasRepo 扩展 ===

func (r *AliasRepo) Delete(ctx context.Context, aliasID string) error {
	res := r.db.WithContext(ctx).Delete(&po.ItemAliasPO{}, "id = ?", aliasID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return bizErr.New(bizErr.ErrNotFound)
	}
	return nil
}

// === ScheduleRepo 扩展 ===

func (r *ScheduleRepo) List(ctx context.Context, deviceID string, startDate, endDate time.Time) ([]resource.ScheduleResp, error) {
	var ps []po.SchedulePO
	q := r.db.WithContext(ctx).Model(&po.SchedulePO{})
	if deviceID != "" {
		q = q.Where("device_id = ?", deviceID)
	}
	if !startDate.IsZero() {
		q = q.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		q = q.Where("date <= ?", endDate)
	}
	if err := q.Order("date ASC, start_time ASC").Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.ScheduleResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.ScheduleResp{
			ID:        p.ID,
			DeviceID:  p.DeviceID,
			Date:      p.Date.Format("2006-01-02"),
			StartTime: p.StartTime,
			EndTime:   p.EndTime,
			Status:    p.Status,
		})
	}
	return out, nil
}

// === CampusRepo ===

func (r *CampusRepo) List(ctx context.Context) ([]resource.CampusResp, error) {
	var ps []po.CampusPO
	if err := r.db.WithContext(ctx).Where("status = ?", "active").Order("name ASC").Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.CampusResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.CampusResp{
			ID:      p.ID,
			Name:    p.Name,
			Code:    p.Code,
			Address: p.Address,
			Status:  p.Status,
		})
	}
	return out, nil
}

// === DepartmentRepo ===

func (r *DepartmentRepo) List(ctx context.Context, campusID string) ([]resource.DepartmentResp, error) {
	var ps []po.DepartmentPO
	q := r.db.WithContext(ctx).Where("status = ?", "active")
	if campusID != "" {
		q = q.Where("campus_id = ?", campusID)
	}
	if err := q.Order("name ASC").Find(&ps).Error; err != nil {
		return nil, err
	}
	out := make([]resource.DepartmentResp, 0, len(ps))
	for _, p := range ps {
		out = append(out, resource.DepartmentResp{
			ID:       p.ID,
			CampusID: p.CampusID,
			Name:     p.Name,
			Code:     p.Code,
			Floor:    p.Floor,
			Status:   p.Status,
		})
	}
	return out, nil
}
