// Package triage 应用层 - 分诊服务应用服务
package triage

import (
	"context"
	"time"

	domain "github.com/euler/mtap/internal/domain/triage"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// TriageAppService 分诊应用服务
type TriageAppService struct {
	checkInSvc *domain.CheckInService
	queueSvc   *domain.QueueManagementService
	examSvc    *domain.ExamStatusService
	execRepo   domain.ExamExecutionRepository
}

// NewTriageAppService 创建分诊应用服务
func NewTriageAppService(
	checkInRepo domain.CheckInRepository,
	queueRepo domain.WaitingQueueRepository,
	entryRepo domain.QueueEntryRepository,
	execRepo domain.ExamExecutionRepository,
) *TriageAppService {
	return &TriageAppService{
		checkInSvc: domain.NewCheckInService(checkInRepo, queueRepo, entryRepo),
		queueSvc:   domain.NewQueueManagementService(queueRepo, entryRepo),
		examSvc:    domain.NewExamStatusService(execRepo),
		execRepo:   execRepo,
	}
}

// KioskCheckIn 自动机签到（通过二维码解析预约信息）
func (s *TriageAppService) KioskCheckIn(ctx context.Context, req KioskCheckInReq) (*CheckInResp, error) {
	// 实际项目中需解析 QRCodeData 获取 appointmentID 等信息
	// 这里简化：直接将 QRCode 作为 appointmentID
	result, err := s.checkInSvc.KioskCheckIn(
		ctx,
		req.QRCodeData, // appointmentID (simplified)
		"",             // patientID
		"",             // patientNameMasked
		"",             // roomID
		"",             // deviceID
		"",             // departmentID
		time.Time{},    // apptStartTime (resolved from QR code in production)
	)
	if err != nil {
		return nil, err
	}
	return ToCheckInResp(result), nil
}

// NurseCheckIn 护士手动签到
func (s *TriageAppService) NurseCheckIn(ctx context.Context, req NurseCheckInReq, operatorID string) (*CheckInResp, error) {
	input := domain.NurseCheckInInput{
		AppointmentID:     req.AppointmentID,
		PatientID:         "", // 实际从预约流程中获取
		PatientNameMasked: "",
		RoomID:            "",
		DeviceID:          "",
		DepartmentID:      "",
		Remark:            req.Remark,
	}
	result, err := s.checkInSvc.NurseCheckIn(ctx, input)
	if err != nil {
		return nil, err
	}
	return ToCheckInResp(result), nil
}

// GetQueueStatus 获取候诊队列状态
func (s *TriageAppService) GetQueueStatus(ctx context.Context, roomID string) (*QueueStatusResp, error) {
	status, err := s.queueSvc.GetQueueStatus(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return ToQueueStatusResp(status), nil
}

// CallNext 呼叫下一位候诊患者
func (s *TriageAppService) CallNext(ctx context.Context, roomID string) (*QueueEntryResp, error) {
	entry, err := s.queueSvc.CallNext(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return ToQueueEntryResp(entry), nil
}

// Recall 重叫当前患者
func (s *TriageAppService) Recall(ctx context.Context, roomID string) (*QueueEntryResp, error) {
	entry, err := s.queueSvc.Recall(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return ToQueueEntryResp(entry), nil
}

// MissAndRequeue 过号重排
func (s *TriageAppService) MissAndRequeue(ctx context.Context, roomID string) (*QueueEntryResp, error) {
	entry, err := s.queueSvc.MissAndRequeue(ctx, roomID)
	if err != nil {
		return nil, err
	}
	return ToQueueEntryResp(entry), nil
}

// StartExam 开始检查
func (s *TriageAppService) StartExam(ctx context.Context, appointmentItemID, operatorID string) error {
	return s.examSvc.StartExam(ctx, appointmentItemID, operatorID)
}

// CompleteExam 检查完成
func (s *TriageAppService) CompleteExam(ctx context.Context, appointmentItemID, operatorID string) error {
	return s.examSvc.CompleteExam(ctx, appointmentItemID, operatorID)
}

// UndoExam 撤销误操作
func (s *TriageAppService) UndoExam(ctx context.Context, appointmentItemID, operatorID, reason string) error {
	return s.examSvc.UndoStatus(ctx, appointmentItemID, operatorID, reason)
}

// GetExamExecution 获取检查执行状态
func (s *TriageAppService) GetExamExecution(ctx context.Context, appointmentItemID string) (*ExamExecutionResp, error) {
	exec, err := s.execRepo.FindByAppointmentItemID(ctx, appointmentItemID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if exec == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToExamExecutionResp(exec), nil
}
