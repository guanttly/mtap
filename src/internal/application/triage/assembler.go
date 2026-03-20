// Package triage 应用层 - 分诊装配器
package triage

import (
	domain "github.com/euler/mtap/internal/domain/triage"
)

// ToCheckInResp 领域结果 → 响应 DTO
func ToCheckInResp(r *domain.CheckInResult) *CheckInResp {
	return &CheckInResp{
		CheckInID:     r.CheckInID,
		QueueNumber:   r.QueueNumber,
		EstimatedWait: r.EstimatedWait,
		RoomLocation:  r.RoomLocation,
		IsLate:        r.IsLate,
	}
}

// ToQueueEntryResp 队列条目 → 响应 DTO
func ToQueueEntryResp(e *domain.QueueEntry) *QueueEntryResp {
	if e == nil {
		return nil
	}
	return &QueueEntryResp{
		ID:                e.ID,
		QueueNumber:       e.QueueNumber,
		PatientNameMasked: e.PatientNameMasked,
		Status:            string(e.Status),
		CallCount:         e.CallCount,
		MissCount:         e.MissCount,
		EnteredAt:         e.EnteredAt,
		CalledAt:          e.CalledAt,
	}
}

// ToQueueStatusResp 队列状态 → 响应 DTO
func ToQueueStatusResp(s *domain.QueueStatus) *QueueStatusResp {
	resp := &QueueStatusResp{
		RoomID:       s.RoomID,
		WaitingCount: s.WaitingCount,
		AverageWait:  s.AverageWait,
	}
	if s.CurrentCalling != nil {
		resp.CurrentCalling = ToQueueEntryResp(s.CurrentCalling)
	}
	for _, e := range s.Entries {
		resp.Entries = append(resp.Entries, QueueEntryResp{
			QueueNumber:       e.QueueNumber,
			PatientNameMasked: e.PatientNameMasked,
			Status:            e.Status,
		})
	}
	return resp
}

// ToExamExecutionResp 检查执行 → 响应 DTO
func ToExamExecutionResp(e *domain.ExamExecution) *ExamExecutionResp {
	return &ExamExecutionResp{
		ID:                e.ID,
		AppointmentItemID: e.AppointmentItemID,
		PatientID:         e.PatientID,
		DeviceID:          e.DeviceID,
		Status:            string(e.Status),
		StartedAt:         e.StartedAt,
		CompletedAt:       e.CompletedAt,
		Duration:          e.Duration,
	}
}
