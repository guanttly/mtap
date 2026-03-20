// Package appointment 基础设施层 - appointment 仓储实现（GORM）
package appointment

import (
	"context"
	"time"

	"gorm.io/gorm"

	domain "github.com/euler/mtap/internal/domain/appointment"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// Repositories appointment 模块仓储集合
type Repositories struct {
	DB *gorm.DB
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{DB: db}
}

// === AppointmentRepository ===

type appointmentRepo struct{ db *gorm.DB }

func (r *Repositories) AppointmentRepo() domain.AppointmentRepository {
	return &appointmentRepo{db: r.DB}
}

func (r *appointmentRepo) Save(ctx context.Context, a *domain.Appointment) error {
	p := toPO(a)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		if len(p.Items) > 0 {
			if err := tx.Create(&p.Items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *appointmentRepo) Update(ctx context.Context, a *domain.Appointment) error {
	return r.db.WithContext(ctx).Model(&po.AppointmentPO{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"status":           string(a.Status),
		"payment_verified": a.PaymentVerified,
		"change_count":     a.ChangeCount,
		"override_by":      a.OverrideBy,
		"override_reason":  a.OverrideReason,
		"cancel_reason":    a.CancelReason,
		"confirmed_at":     a.ConfirmedAt,
		"cancelled_at":     a.CancelledAt,
		"updated_at":       a.UpdatedAt,
	}).Error
}

func (r *appointmentRepo) FindByID(ctx context.Context, id string) (*domain.Appointment, error) {
	var p po.AppointmentPO
	err := r.db.WithContext(ctx).Preload("Items").First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return fromPO(p), nil
}

func (r *appointmentRepo) FindByPatientID(ctx context.Context, patientID string, page, size int) ([]*domain.Appointment, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&po.AppointmentPO{}).Where("patient_id = ?", patientID).Count(&total)
	var ps []po.AppointmentPO
	err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).
		Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&ps).Error
	if err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return fromPOs(ps), total, nil
}

func (r *appointmentRepo) FindByStatus(ctx context.Context, status domain.AppointmentStatus, page, size int) ([]*domain.Appointment, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&po.AppointmentPO{}).Where("status = ?", string(status)).Count(&total)
	var ps []po.AppointmentPO
	err := r.db.WithContext(ctx).Where("status = ?", string(status)).
		Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&ps).Error
	if err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return fromPOs(ps), total, nil
}

func (r *appointmentRepo) FindConfirmedPendingTimeout(ctx context.Context) ([]*domain.Appointment, error) {
	deadline := time.Now().Add(-5 * time.Minute)
	var ps []po.AppointmentPO
	err := r.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", string(domain.StatusPending), deadline).
		Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return fromPOs(ps), nil
}

// === AppointmentItemRepository ===

type appointmentItemRepo struct{ db *gorm.DB }

func (r *Repositories) AppointmentItemRepo() domain.AppointmentItemRepository {
	return &appointmentItemRepo{db: r.DB}
}

func (r *appointmentItemRepo) FindByAppointmentID(ctx context.Context, appointmentID string) ([]*domain.AppointmentItem, error) {
	var ps []po.AppointmentItemPO
	err := r.db.WithContext(ctx).Where("appointment_id = ?", appointmentID).Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.AppointmentItem, 0, len(ps))
	for _, p := range ps {
		out = append(out, itemFromPO(p))
	}
	return out, nil
}

func (r *appointmentItemRepo) UpdateStatus(ctx context.Context, itemID string, status domain.ItemStatus) error {
	return r.db.WithContext(ctx).Model(&po.AppointmentItemPO{}).
		Where("id = ?", itemID).Update("status", string(status)).Error
}

// === CredentialRepository ===

type credentialRepo struct{ db *gorm.DB }

func (r *Repositories) CredentialRepo() domain.CredentialRepository {
	return &credentialRepo{db: r.DB}
}

func (r *credentialRepo) Save(ctx context.Context, c *domain.Credential) error {
	return r.db.WithContext(ctx).Create(&po.AppointmentCredentialPO{
		ID:                c.ID,
		AppointmentID:     c.AppointmentID,
		QRCodeData:        c.QRCodeData,
		PatientNameMasked: c.PatientNameMasked,
		ExamSummary:       c.ExamSummary,
		NoticeContent:     c.NoticeContent,
		GeneratedAt:       c.GeneratedAt,
	}).Error
}

func (r *credentialRepo) FindByAppointmentID(ctx context.Context, appointmentID string) (*domain.Credential, error) {
	var p po.AppointmentCredentialPO
	err := r.db.WithContext(ctx).First(&p, "appointment_id = ?", appointmentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &domain.Credential{
		ID:                p.ID,
		AppointmentID:     p.AppointmentID,
		QRCodeData:        p.QRCodeData,
		PatientNameMasked: p.PatientNameMasked,
		ExamSummary:       p.ExamSummary,
		NoticeContent:     p.NoticeContent,
		GeneratedAt:       p.GeneratedAt,
	}, nil
}

// === BlacklistRepository ===

type blacklistRepo struct{ db *gorm.DB }

func (r *Repositories) BlacklistRepo() domain.BlacklistRepository {
	return &blacklistRepo{db: r.DB}
}

func (r *blacklistRepo) Save(ctx context.Context, b *domain.Blacklist) error {
	return r.db.WithContext(ctx).Create(&po.BlacklistPO{
		ID:          b.ID,
		PatientID:   b.PatientID,
		TriggerTime: b.TriggerTime,
		ExpiresAt:   b.ExpiresAt,
		Status:      string(b.Status),
	}).Error
}

func (r *blacklistRepo) Update(ctx context.Context, b *domain.Blacklist) error {
	return r.db.WithContext(ctx).Model(&po.BlacklistPO{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
		"status":         string(b.Status),
		"released_at":    b.ReleasedAt,
		"release_reason": b.ReleaseReason,
	}).Error
}

func (r *blacklistRepo) FindByPatientID(ctx context.Context, patientID string) (*domain.Blacklist, error) {
	var p po.BlacklistPO
	err := r.db.WithContext(ctx).
		Where("patient_id = ? AND status = ?", patientID, string(domain.BlacklistActive)).
		First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return blacklistFromPO(p), nil
}

func (r *blacklistRepo) FindAll(ctx context.Context, page, size int) ([]*domain.Blacklist, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&po.BlacklistPO{}).Count(&total)
	var ps []po.BlacklistPO
	err := r.db.WithContext(ctx).Offset((page - 1) * size).Limit(size).
		Order("created_at DESC").Find(&ps).Error
	if err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.Blacklist, 0, len(ps))
	for _, p := range ps {
		out = append(out, blacklistFromPO(p))
	}
	return out, total, nil
}

func (r *blacklistRepo) FindExpired(ctx context.Context) ([]*domain.Blacklist, error) {
	var ps []po.BlacklistPO
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < ?", string(domain.BlacklistActive), time.Now()).
		Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.Blacklist, 0, len(ps))
	for _, p := range ps {
		out = append(out, blacklistFromPO(p))
	}
	return out, nil
}

// === NoShowRecordRepository ===

type noShowRepo struct{ db *gorm.DB }

func (r *Repositories) NoShowRepo() domain.NoShowRecordRepository {
	return &noShowRepo{db: r.DB}
}

func (r *noShowRepo) Save(ctx context.Context, rec *domain.NoShowRecord) error {
	return r.db.WithContext(ctx).Create(&po.NoShowRecordPO{
		ID:            rec.ID,
		PatientID:     rec.PatientID,
		AppointmentID: rec.AppointmentID,
		OccurredAt:    rec.OccurredAt,
	}).Error
}

func (r *noShowRepo) CountByPatientIDInWindow(ctx context.Context, patientID string, days int) (int, error) {
	var count int64
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.WithContext(ctx).Model(&po.NoShowRecordPO{}).
		Where("patient_id = ? AND occurred_at >= ?", patientID, since).
		Count(&count).Error
	return int(count), err
}

func (r *noShowRepo) FindByPatientID(ctx context.Context, patientID string) ([]*domain.NoShowRecord, error) {
	var ps []po.NoShowRecordPO
	err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).
		Order("occurred_at DESC").Find(&ps).Error
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	out := make([]*domain.NoShowRecord, 0, len(ps))
	for _, p := range ps {
		out = append(out, &domain.NoShowRecord{
			ID:            p.ID,
			PatientID:     p.PatientID,
			AppointmentID: p.AppointmentID,
			OccurredAt:    p.OccurredAt,
		})
	}
	return out, nil
}

// === AppealRepository ===

type appealRepo struct{ db *gorm.DB }

func (r *Repositories) AppealRepo() domain.AppealRepository {
	return &appealRepo{db: r.DB}
}

func (r *appealRepo) Save(ctx context.Context, a *domain.Appeal) error {
	return r.db.WithContext(ctx).Create(&po.AppealPO{
		ID:          a.ID,
		BlacklistID: a.BlacklistID,
		Reason:      a.Reason,
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
	}).Error
}

func (r *appealRepo) Update(ctx context.Context, a *domain.Appeal) error {
	return r.db.WithContext(ctx).Model(&po.AppealPO{}).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"status":      string(a.Status),
		"reviewed_by": a.ReviewedBy,
		"reviewed_at": a.ReviewedAt,
	}).Error
}

func (r *appealRepo) FindByID(ctx context.Context, id string) (*domain.Appeal, error) {
	var p po.AppealPO
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return appealFromPO(p), nil
}

func (r *appealRepo) FindByBlacklistID(ctx context.Context, blacklistID string) (*domain.Appeal, error) {
	var p po.AppealPO
	err := r.db.WithContext(ctx).Where("blacklist_id = ?", blacklistID).
		Order("created_at DESC").First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return appealFromPO(p), nil
}

// === PO <-> Domain 转换 ===

func toPO(a *domain.Appointment) po.AppointmentPO {
	items := make([]po.AppointmentItemPO, 0, len(a.Items))
	for _, it := range a.Items {
		items = append(items, po.AppointmentItemPO{
			ID:            it.ID,
			AppointmentID: it.AppointmentID,
			ExamItemID:    it.ExamItemID,
			SlotID:        it.SlotID,
			DeviceID:      it.DeviceID,
			StartTime:     it.StartTime,
			EndTime:       it.EndTime,
			Status:        string(it.Status),
		})
	}
	return po.AppointmentPO{
		ID:              a.ID,
		PatientID:       a.PatientID,
		Mode:            string(a.Mode),
		Status:          string(a.Status),
		OverrideBy:      a.OverrideBy,
		OverrideReason:  a.OverrideReason,
		PaymentVerified: a.PaymentVerified,
		ChangeCount:     a.ChangeCount,
		CancelReason:    a.CancelReason,
		ConfirmedAt:     a.ConfirmedAt,
		CancelledAt:     a.CancelledAt,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
		Items:           items,
	}
}

func fromPO(p po.AppointmentPO) *domain.Appointment {
	a := &domain.Appointment{
		ID:              p.ID,
		PatientID:       p.PatientID,
		Mode:            domain.AppointmentMode(p.Mode),
		Status:          domain.AppointmentStatus(p.Status),
		OverrideBy:      p.OverrideBy,
		OverrideReason:  p.OverrideReason,
		PaymentVerified: p.PaymentVerified,
		ChangeCount:     p.ChangeCount,
		CancelReason:    p.CancelReason,
		ConfirmedAt:     p.ConfirmedAt,
		CancelledAt:     p.CancelledAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
	for _, it := range p.Items {
		a.Items = append(a.Items, *itemFromPO(it))
	}
	return a
}

func fromPOs(ps []po.AppointmentPO) []*domain.Appointment {
	out := make([]*domain.Appointment, 0, len(ps))
	for _, p := range ps {
		out = append(out, fromPO(p))
	}
	return out
}

func itemFromPO(p po.AppointmentItemPO) *domain.AppointmentItem {
	return &domain.AppointmentItem{
		ID:            p.ID,
		AppointmentID: p.AppointmentID,
		ExamItemID:    p.ExamItemID,
		SlotID:        p.SlotID,
		DeviceID:      p.DeviceID,
		StartTime:     p.StartTime,
		EndTime:       p.EndTime,
		Status:        domain.ItemStatus(p.Status),
	}
}

func blacklistFromPO(p po.BlacklistPO) *domain.Blacklist {
	return &domain.Blacklist{
		ID:            p.ID,
		PatientID:     p.PatientID,
		TriggerTime:   p.TriggerTime,
		ExpiresAt:     p.ExpiresAt,
		Status:        domain.BlacklistStatus(p.Status),
		ReleasedAt:    p.ReleasedAt,
		ReleaseReason: p.ReleaseReason,
	}
}

func appealFromPO(p po.AppealPO) *domain.Appeal {
	return &domain.Appeal{
		ID:          p.ID,
		BlacklistID: p.BlacklistID,
		Reason:      p.Reason,
		Status:      domain.AppealStatus(p.Status),
		ReviewedBy:  p.ReviewedBy,
		ReviewedAt:  p.ReviewedAt,
		CreatedAt:   p.CreatedAt,
	}
}
