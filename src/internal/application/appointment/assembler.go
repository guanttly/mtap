// Package appointment 应用层 - 预约服务装配器
package appointment

import (
	domain "github.com/euler/mtap/internal/domain/appointment"
)

// ToAppointmentResp 领域实体 → 响应 DTO
func ToAppointmentResp(a *domain.Appointment) *AppointmentResp {
	items := make([]AppointmentItemResp, 0, len(a.Items))
	for _, it := range a.Items {
		items = append(items, AppointmentItemResp{
			ID:         it.ID,
			ExamItemID: it.ExamItemID,
			SlotID:     it.SlotID,
			DeviceID:   it.DeviceID,
			StartTime:  it.StartTime,
			EndTime:    it.EndTime,
			Status:     string(it.Status),
		})
	}
	return &AppointmentResp{
		ID:              a.ID,
		PatientID:       a.PatientID,
		Mode:            string(a.Mode),
		Status:          string(a.Status),
		PaymentVerified: a.PaymentVerified,
		ChangeCount:     a.ChangeCount,
		Items:           items,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
		ConfirmedAt:     a.ConfirmedAt,
	}
}

// ToBlacklistResp 领域实体 → 黑名单响应 DTO
func ToBlacklistResp(b *domain.Blacklist, noShowCount int, hasPendingAppeal bool) *BlacklistResp {
	return &BlacklistResp{
		ID:               b.ID,
		PatientID:        b.PatientID,
		TriggerTime:      b.TriggerTime,
		ExpiresAt:        b.ExpiresAt,
		Status:           string(b.Status),
		NoShowCount:      noShowCount,
		HasPendingAppeal: hasPendingAppeal,
	}
}

// ToAppealResp 领域实体 → 申诉响应 DTO
func ToAppealResp(a *domain.Appeal) *AppealResp {
	return &AppealResp{
		ID:          a.ID,
		BlacklistID: a.BlacklistID,
		Reason:      a.Reason,
		Status:      string(a.Status),
		ReviewedBy:  a.ReviewedBy,
		ReviewedAt:  a.ReviewedAt,
		CreatedAt:   a.CreatedAt,
	}
}

// ToCredentialResp 领域实体 → 凭证响应 DTO
func ToCredentialResp(c *domain.Credential) *CredentialResp {
	return &CredentialResp{
		ID:                c.ID,
		AppointmentID:     c.AppointmentID,
		QRCodeURL:         c.QRCodeData,
		PatientNameMasked: c.PatientNameMasked,
		ExamSummary:       c.ExamSummary,
		NoticeContent:     c.NoticeContent,
		GeneratedAt:       c.GeneratedAt,
	}
}
