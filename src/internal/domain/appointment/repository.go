// Package appointment 预约服务领域 - 仓储接口
package appointment

import "context"

// AppointmentRepository 预约单仓储接口
type AppointmentRepository interface {
	Save(ctx context.Context, a *Appointment) error
	Update(ctx context.Context, a *Appointment) error
	FindByID(ctx context.Context, id string) (*Appointment, error)
	FindByPatientID(ctx context.Context, patientID string, page, size int) ([]*Appointment, int64, error)
	FindByStatus(ctx context.Context, status AppointmentStatus, page, size int) ([]*Appointment, int64, error)
	FindConfirmedPendingTimeout(ctx context.Context) ([]*Appointment, error) // 预约确认超时
}

// AppointmentItemRepository 预约项目仓储接口
type AppointmentItemRepository interface {
	FindByAppointmentID(ctx context.Context, appointmentID string) ([]*AppointmentItem, error)
	UpdateStatus(ctx context.Context, itemID string, status ItemStatus) error
}

// CredentialRepository 预约凭证仓储接口
type CredentialRepository interface {
	Save(ctx context.Context, c *Credential) error
	FindByAppointmentID(ctx context.Context, appointmentID string) (*Credential, error)
}

// BlacklistRepository 黑名单仓储接口
type BlacklistRepository interface {
	Save(ctx context.Context, b *Blacklist) error
	Update(ctx context.Context, b *Blacklist) error
	FindByPatientID(ctx context.Context, patientID string) (*Blacklist, error) // 生效中的黑名单
	FindAll(ctx context.Context, page, size int) ([]*Blacklist, int64, error)
	FindExpired(ctx context.Context) ([]*Blacklist, error)
}

// NoShowRecordRepository 爽约记录仓储接口
type NoShowRecordRepository interface {
	Save(ctx context.Context, r *NoShowRecord) error
	CountByPatientIDInWindow(ctx context.Context, patientID string, days int) (int, error)
	FindByPatientID(ctx context.Context, patientID string) ([]*NoShowRecord, error)
}

// AppealRepository 申诉仓储接口
type AppealRepository interface {
	Save(ctx context.Context, a *Appeal) error
	Update(ctx context.Context, a *Appeal) error
	FindByID(ctx context.Context, id string) (*Appeal, error)
	FindByBlacklistID(ctx context.Context, blacklistID string) (*Appeal, error)
}
