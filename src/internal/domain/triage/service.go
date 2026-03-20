// Package triage 分诊执行领域 - 领域服务
package triage

import (
	"context"
	"time"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// QueueStatus 队列状态快照
type QueueStatus struct {
	RoomID         string            `json:"room_id"`
	WaitingCount   int               `json:"waiting_count"`
	AverageWait    int               `json:"average_wait"`
	CurrentCalling *QueueEntry       `json:"current_calling"`
	Entries        []QueueEntryBrief `json:"entries"`
}

// QueueEntryBrief 队列条目简要信息
type QueueEntryBrief struct {
	QueueNumber       int    `json:"queue_number"`
	PatientNameMasked string `json:"patient_name_masked"`
	Status            string `json:"status"`
}

// CheckInService 签到领域服务
type CheckInService struct {
	checkInRepo CheckInRepository
	queueRepo   WaitingQueueRepository
	entryRepo   QueueEntryRepository
}

// NewCheckInService 创建签到服务
func NewCheckInService(
	checkInRepo CheckInRepository,
	queueRepo WaitingQueueRepository,
	entryRepo QueueEntryRepository,
) *CheckInService {
	return &CheckInService{
		checkInRepo: checkInRepo,
		queueRepo:   queueRepo,
		entryRepo:   entryRepo,
	}
}

// CheckInResult 签到结果
type CheckInResult struct {
	CheckInID     string `json:"check_in_id"`
	QueueNumber   int    `json:"queue_number"`
	EstimatedWait int    `json:"estimated_wait"`
	RoomLocation  string `json:"room_location"`
	IsLate        bool   `json:"is_late"`
}

// NurseCheckInInput 护士手动签到输入
type NurseCheckInInput struct {
	AppointmentID     string
	PatientID         string
	PatientNameMasked string
	RoomID            string
	DeviceID          string
	DepartmentID      string
	ApptStartTime     time.Time
	Remark            string
}

// NurseCheckIn 护士站手动签到
func (s *CheckInService) NurseCheckIn(ctx context.Context, input NurseCheckInInput) (*CheckInResult, error) {
	// 检查是否已签到
	existing, err := s.checkInRepo.FindByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if existing != nil {
		return nil, bizErr.New(bizErr.ErrTriageAlreadyCheckedIn)
	}

	checkIn, err := NewCheckIn(input.AppointmentID, input.PatientID, CheckInNurse, input.ApptStartTime, input.Remark)
	if err != nil {
		return nil, err
	}
	if err := s.checkInRepo.Save(ctx, checkIn); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	// 入队
	queue, err := s.queueRepo.FindOrCreateByRoom(ctx, input.RoomID, input.DeviceID, input.DepartmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	entry := queue.AddEntry(checkIn, input.PatientNameMasked, input.AppointmentID)
	if err := s.entryRepo.Save(ctx, entry); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return &CheckInResult{
		CheckInID:     checkIn.ID,
		QueueNumber:   entry.QueueNumber,
		EstimatedWait: queue.EstimateWaitTime(),
		IsLate:        checkIn.IsLate,
	}, nil
}

// KioskCheckIn 自动机扫码签到
func (s *CheckInService) KioskCheckIn(ctx context.Context, appointmentID, patientID, patientNameMasked, roomID, deviceID, departmentID string, apptStartTime time.Time) (*CheckInResult, error) {
	existing, err := s.checkInRepo.FindByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if existing != nil {
		return nil, bizErr.New(bizErr.ErrTriageAlreadyCheckedIn)
	}
	checkIn, err := NewCheckIn(appointmentID, patientID, CheckInKiosk, apptStartTime, "")
	if err != nil {
		return nil, err
	}
	if err := s.checkInRepo.Save(ctx, checkIn); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	queue, err := s.queueRepo.FindOrCreateByRoom(ctx, roomID, deviceID, departmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	entry := queue.AddEntry(checkIn, patientNameMasked, appointmentID)
	if err := s.entryRepo.Save(ctx, entry); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &CheckInResult{
		CheckInID:     checkIn.ID,
		QueueNumber:   entry.QueueNumber,
		EstimatedWait: queue.EstimateWaitTime(),
		IsLate:        checkIn.IsLate,
	}, nil
}

// QueueManagementService 队列管理领域服务
type QueueManagementService struct {
	queueRepo WaitingQueueRepository
	entryRepo QueueEntryRepository
}

// NewQueueManagementService 创建队列管理服务
func NewQueueManagementService(queueRepo WaitingQueueRepository, entryRepo QueueEntryRepository) *QueueManagementService {
	return &QueueManagementService{queueRepo: queueRepo, entryRepo: entryRepo}
}

// CallNext 呼叫下一位候诊患者
func (s *QueueManagementService) CallNext(ctx context.Context, roomID string) (*QueueEntry, error) {
	queue, err := s.queueRepo.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if queue == nil {
		return nil, bizErr.New(bizErr.ErrTriageQueueEmpty)
	}
	entry, err := queue.CallNext()
	if err != nil {
		return nil, err
	}
	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return entry, nil
}

// Recall 重叫当前患者
func (s *QueueManagementService) Recall(ctx context.Context, roomID string) (*QueueEntry, error) {
	queue, err := s.queueRepo.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if queue == nil {
		return nil, bizErr.New(bizErr.ErrTriageQueueEmpty)
	}
	entry, err := queue.Recall()
	if err != nil {
		return nil, err
	}
	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return entry, nil
}

// MissAndRequeue 过号重排
func (s *QueueManagementService) MissAndRequeue(ctx context.Context, roomID string) (*QueueEntry, error) {
	queue, err := s.queueRepo.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if queue == nil {
		return nil, bizErr.New(bizErr.ErrTriageQueueEmpty)
	}
	entry, err := queue.MissAndRequeue()
	if err != nil {
		return nil, err
	}
	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return entry, nil
}

// GetQueueStatus 获取候诊队列状态
func (s *QueueManagementService) GetQueueStatus(ctx context.Context, roomID string) (*QueueStatus, error) {
	queue, err := s.queueRepo.FindByRoomID(ctx, roomID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if queue == nil {
		return &QueueStatus{RoomID: roomID}, nil
	}

	status := &QueueStatus{
		RoomID:       roomID,
		WaitingCount: queue.GetWaitCount(),
		AverageWait:  queue.EstimateWaitTime(),
	}
	for i := range queue.Entries {
		e := &queue.Entries[i]
		if e.Status == EntryCalling {
			status.CurrentCalling = e
		}
		if e.Status == EntryWaiting || e.Status == EntryCalling {
			status.Entries = append(status.Entries, QueueEntryBrief{
				QueueNumber:       e.QueueNumber,
				PatientNameMasked: e.PatientNameMasked,
				Status:            string(e.Status),
			})
		}
	}
	return status, nil
}

// ExamStatusService 检查状态领域服务
type ExamStatusService struct {
	execRepo ExamExecutionRepository
}

// NewExamStatusService 创建检查状态服务
func NewExamStatusService(execRepo ExamExecutionRepository) *ExamStatusService {
	return &ExamStatusService{execRepo: execRepo}
}

// StartExam 开始检查
func (s *ExamStatusService) StartExam(ctx context.Context, appointmentItemID, operatorID string) error {
	exec, err := s.execRepo.FindByAppointmentItemID(ctx, appointmentItemID)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if exec == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := exec.Start(operatorID); err != nil {
		return err
	}
	return s.execRepo.Update(ctx, exec)
}

// CompleteExam 完成检查
func (s *ExamStatusService) CompleteExam(ctx context.Context, appointmentItemID, operatorID string) error {
	exec, err := s.execRepo.FindByAppointmentItemID(ctx, appointmentItemID)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if exec == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := exec.Complete(operatorID); err != nil {
		return err
	}
	return s.execRepo.Update(ctx, exec)
}

// UndoStatus 撤销误操作
func (s *ExamStatusService) UndoStatus(ctx context.Context, appointmentItemID, operatorID, reason string) error {
	exec, err := s.execRepo.FindByAppointmentItemID(ctx, appointmentItemID)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if exec == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := exec.Undo(operatorID, reason); err != nil {
		return err
	}
	return s.execRepo.Update(ctx, exec)
}
