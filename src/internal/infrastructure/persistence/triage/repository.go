// Package triage 基础设施层 - triage 仓储实现（GORM）
package triage

import (
	"context"

	"gorm.io/gorm"

	domain "github.com/euler/mtap/internal/domain/triage"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// Repositories triage 模块仓储集合
type Repositories struct {
	DB *gorm.DB
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{DB: db}
}

// === CheckInRepository ===

type checkInRepo struct{ db *gorm.DB }

func (r *Repositories) CheckInRepo() domain.CheckInRepository {
	return &checkInRepo{db: r.DB}
}

func (r *checkInRepo) Save(ctx context.Context, c *domain.CheckIn) error {
	return r.db.WithContext(ctx).Create(&po.CheckInPO{
		ID:            c.ID,
		AppointmentID: c.AppointmentID,
		PatientID:     c.PatientID,
		Method:        string(c.Method),
		CheckInTime:   c.CheckInTime,
		IsLate:        c.IsLate,
		Remark:        c.Remark,
	}).Error
}

func (r *checkInRepo) FindByAppointmentID(ctx context.Context, appointmentID string) (*domain.CheckIn, error) {
	var p po.CheckInPO
	err := r.db.WithContext(ctx).First(&p, "appointment_id = ?", appointmentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &domain.CheckIn{
		ID:            p.ID,
		AppointmentID: p.AppointmentID,
		PatientID:     p.PatientID,
		Method:        domain.CheckInMethod(p.Method),
		CheckInTime:   p.CheckInTime,
		IsLate:        p.IsLate,
		Remark:        p.Remark,
	}, nil
}

func (r *checkInRepo) FindByPatientAndDate(ctx context.Context, patientID, date string) ([]*domain.CheckIn, error) {
	var ps []po.CheckInPO
	err := r.db.WithContext(ctx).
		Where("patient_id = ? AND DATE(check_in_time) = ?", patientID, date).
		Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.CheckIn, 0, len(ps))
	for _, p := range ps {
		out = append(out, &domain.CheckIn{
			ID:            p.ID,
			AppointmentID: p.AppointmentID,
			PatientID:     p.PatientID,
			Method:        domain.CheckInMethod(p.Method),
			CheckInTime:   p.CheckInTime,
			IsLate:        p.IsLate,
			Remark:        p.Remark,
		})
	}
	return out, nil
}

// === WaitingQueueRepository ===

type waitingQueueRepo struct{ db *gorm.DB }

func (r *Repositories) WaitingQueueRepo() domain.WaitingQueueRepository {
	return &waitingQueueRepo{db: r.DB}
}

func (r *waitingQueueRepo) Save(ctx context.Context, q *domain.WaitingQueue) error {
	return r.db.WithContext(ctx).Create(&po.WaitingQueuePO{
		ID:           q.ID,
		RoomID:       q.RoomID,
		DeviceID:     q.DeviceID,
		DepartmentID: q.DepartmentID,
		Status:       q.Status,
	}).Error
}

func (r *waitingQueueRepo) Update(ctx context.Context, q *domain.WaitingQueue) error {
	return r.db.WithContext(ctx).Model(&po.WaitingQueuePO{}).
		Where("id = ?", q.ID).Update("status", q.Status).Error
}

func (r *waitingQueueRepo) FindByRoomID(ctx context.Context, roomID string) (*domain.WaitingQueue, error) {
	var p po.WaitingQueuePO
	err := r.db.WithContext(ctx).Where("room_id = ? AND status = ?", roomID, "active").First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	q := &domain.WaitingQueue{
		ID:           p.ID,
		RoomID:       p.RoomID,
		DeviceID:     p.DeviceID,
		DepartmentID: p.DepartmentID,
		Status:       p.Status,
	}
	// 加载条目
	var entries []po.QueueEntryPO
	r.db.WithContext(ctx).Where("queue_id = ? AND status IN ?", q.ID,
		[]string{"waiting", "calling"}).Order("entered_at ASC").Find(&entries)
	for _, e := range entries {
		q.Entries = append(q.Entries, entryFromPO(e))
	}
	return q, nil
}

func (r *waitingQueueRepo) FindOrCreateByRoom(ctx context.Context, roomID, deviceID, departmentID string) (*domain.WaitingQueue, error) {
	q, err := r.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if q != nil {
		return q, nil
	}
	newQ := domain.NewWaitingQueue(roomID, deviceID, departmentID)
	if err := r.Save(ctx, newQ); err != nil {
		return nil, err
	}
	return newQ, nil
}

// === QueueEntryRepository ===

type queueEntryRepo struct{ db *gorm.DB }

func (r *Repositories) QueueEntryRepo() domain.QueueEntryRepository {
	return &queueEntryRepo{db: r.DB}
}

func (r *queueEntryRepo) Save(ctx context.Context, e *domain.QueueEntry) error {
	return r.db.WithContext(ctx).Create(&po.QueueEntryPO{
		ID:                e.ID,
		QueueID:           e.QueueID,
		PatientID:         e.PatientID,
		PatientNameMasked: e.PatientNameMasked,
		AppointmentID:     e.AppointmentID,
		CheckInID:         e.CheckInID,
		QueueNumber:       e.QueueNumber,
		Status:            string(e.Status),
		CallCount:         e.CallCount,
		MissCount:         e.MissCount,
		EnteredAt:         e.EnteredAt,
		CalledAt:          e.CalledAt,
		CompletedAt:       e.CompletedAt,
	}).Error
}

func (r *queueEntryRepo) Update(ctx context.Context, e *domain.QueueEntry) error {
	return r.db.WithContext(ctx).Model(&po.QueueEntryPO{}).Where("id = ?", e.ID).Updates(map[string]interface{}{
		"status":       string(e.Status),
		"call_count":   e.CallCount,
		"miss_count":   e.MissCount,
		"entered_at":   e.EnteredAt,
		"called_at":    e.CalledAt,
		"completed_at": e.CompletedAt,
	}).Error
}

func (r *queueEntryRepo) FindByQueueID(ctx context.Context, queueID string, status domain.EntryStatus) ([]*domain.QueueEntry, error) {
	var ps []po.QueueEntryPO
	q := r.db.WithContext(ctx).Where("queue_id = ?", queueID)
	if status != "" {
		q = q.Where("status = ?", string(status))
	}
	if err := q.Order("entered_at ASC").Find(&ps).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.QueueEntry, 0, len(ps))
	for _, p := range ps {
		e := entryFromPO(p)
		out = append(out, &e)
	}
	return out, nil
}

// === ExamExecutionRepository ===

type examExecutionRepo struct{ db *gorm.DB }

func (r *Repositories) ExamExecutionRepo() domain.ExamExecutionRepository {
	return &examExecutionRepo{db: r.DB}
}

func (r *examExecutionRepo) Save(ctx context.Context, e *domain.ExamExecution) error {
	return r.db.WithContext(ctx).Create(&po.ExamExecutionPO{
		ID:                e.ID,
		AppointmentItemID: e.AppointmentItemID,
		PatientID:         e.PatientID,
		DeviceID:          e.DeviceID,
		Status:            string(e.Status),
		StartedAt:         e.StartedAt,
		CompletedAt:       e.CompletedAt,
		Duration:          e.Duration,
		OperatorID:        e.OperatorID,
		UndoDeadline:      e.UndoDeadline,
	}).Error
}

func (r *examExecutionRepo) Update(ctx context.Context, e *domain.ExamExecution) error {
	return r.db.WithContext(ctx).Model(&po.ExamExecutionPO{}).Where("id = ?", e.ID).Updates(map[string]interface{}{
		"status":        string(e.Status),
		"started_at":    e.StartedAt,
		"completed_at":  e.CompletedAt,
		"duration":      e.Duration,
		"operator_id":   e.OperatorID,
		"undo_deadline": e.UndoDeadline,
		"updated_at":    e.UpdatedAt,
	}).Error
}

func (r *examExecutionRepo) FindByID(ctx context.Context, id string) (*domain.ExamExecution, error) {
	var p po.ExamExecutionPO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return execFromPO(p), nil
}

func (r *examExecutionRepo) FindByAppointmentItemID(ctx context.Context, itemID string) (*domain.ExamExecution, error) {
	var p po.ExamExecutionPO
	err := r.db.WithContext(ctx).First(&p, "appointment_item_id = ?", itemID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return execFromPO(p), nil
}

func (r *examExecutionRepo) FindByDevice(ctx context.Context, deviceID string) ([]*domain.ExamExecution, error) {
	var ps []po.ExamExecutionPO
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).
		Order("created_at DESC").Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.ExamExecution, 0, len(ps))
	for _, p := range ps {
		out = append(out, execFromPO(p))
	}
	return out, nil
}

// === 转换函数 ===

func entryFromPO(p po.QueueEntryPO) domain.QueueEntry {
	return domain.QueueEntry{
		ID:                p.ID,
		QueueID:           p.QueueID,
		PatientID:         p.PatientID,
		PatientNameMasked: p.PatientNameMasked,
		AppointmentID:     p.AppointmentID,
		CheckInID:         p.CheckInID,
		QueueNumber:       p.QueueNumber,
		Status:            domain.EntryStatus(p.Status),
		CallCount:         p.CallCount,
		MissCount:         p.MissCount,
		EnteredAt:         p.EnteredAt,
		CalledAt:          p.CalledAt,
		CompletedAt:       p.CompletedAt,
	}
}

func execFromPO(p po.ExamExecutionPO) *domain.ExamExecution {
	return &domain.ExamExecution{
		ID:                p.ID,
		AppointmentItemID: p.AppointmentItemID,
		PatientID:         p.PatientID,
		DeviceID:          p.DeviceID,
		Status:            domain.ExamStatus(p.Status),
		StartedAt:         p.StartedAt,
		CompletedAt:       p.CompletedAt,
		Duration:          p.Duration,
		OperatorID:        p.OperatorID,
		UndoDeadline:      p.UndoDeadline,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}
