// Package triage 分诊执行领域 - 仓储接口
package triage

import "context"

// CheckInRepository 签到仓储接口
type CheckInRepository interface {
	Save(ctx context.Context, c *CheckIn) error
	FindByAppointmentID(ctx context.Context, appointmentID string) (*CheckIn, error)
	FindByPatientAndDate(ctx context.Context, patientID, date string) ([]*CheckIn, error)
}

// WaitingQueueRepository 候诊队列仓储接口
type WaitingQueueRepository interface {
	Save(ctx context.Context, q *WaitingQueue) error
	Update(ctx context.Context, q *WaitingQueue) error
	FindByRoomID(ctx context.Context, roomID string) (*WaitingQueue, error)
	FindOrCreateByRoom(ctx context.Context, roomID, deviceID, departmentID string) (*WaitingQueue, error)
}

// QueueEntryRepository 队列条目仓储接口
type QueueEntryRepository interface {
	Save(ctx context.Context, e *QueueEntry) error
	Update(ctx context.Context, e *QueueEntry) error
	FindByQueueID(ctx context.Context, queueID string, status EntryStatus) ([]*QueueEntry, error)
}

// ExamExecutionRepository 检查执行仓储接口
type ExamExecutionRepository interface {
	Save(ctx context.Context, e *ExamExecution) error
	Update(ctx context.Context, e *ExamExecution) error
	FindByID(ctx context.Context, id string) (*ExamExecution, error)
	FindByAppointmentItemID(ctx context.Context, itemID string) (*ExamExecution, error)
	FindByDevice(ctx context.Context, deviceID string) ([]*ExamExecution, error)
}
