// Package triage 分诊执行领域 - 实体定义
package triage

import (
	"sort"
	"time"

	"github.com/google/uuid"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// CheckIn 签到记录聚合根
type CheckIn struct {
	ID            string        `json:"id"`
	AppointmentID string        `json:"appointment_id"`
	PatientID     string        `json:"patient_id"`
	Method        CheckInMethod `json:"method"`
	CheckInTime   time.Time     `json:"check_in_time"`
	IsLate        bool          `json:"is_late"`
	Remark        string        `json:"remark"`
}

// NewCheckIn 创建签到记录
func NewCheckIn(appointmentID, patientID string, method CheckInMethod, apptStartTime time.Time, remark string) (*CheckIn, error) {
	if !method.IsValid() {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "无效的签到方式")
	}
	now := time.Now()
	// 检查签到时间窗口：预约前30分钟~预约后15分钟
	windowStart := apptStartTime.Add(-30 * time.Minute)
	windowEnd := apptStartTime.Add(15 * time.Minute)
	if now.Before(windowStart) || now.After(windowEnd) {
		if method != CheckInNurse {
			return nil, bizErr.New(bizErr.ErrTriageOutOfWindow)
		}
	}
	return &CheckIn{
		ID:            uuid.New().String(),
		AppointmentID: appointmentID,
		PatientID:     patientID,
		Method:        method,
		CheckInTime:   now,
		IsLate:        now.After(apptStartTime.Add(15 * time.Minute)),
		Remark:        remark,
	}, nil
}

// WaitingQueue 候诊队列聚合根
type WaitingQueue struct {
	ID           string       `json:"id"`
	RoomID       string       `json:"room_id"`
	DeviceID     string       `json:"device_id"`
	DepartmentID string       `json:"department_id"`
	Status       string       `json:"status"` // active
	Entries      []QueueEntry `json:"entries"`
}

// NewWaitingQueue 创建候诊队列
func NewWaitingQueue(roomID, deviceID, departmentID string) *WaitingQueue {
	return &WaitingQueue{
		ID:           uuid.New().String(),
		RoomID:       roomID,
		DeviceID:     deviceID,
		DepartmentID: departmentID,
		Status:       "active",
	}
}

// AddEntry 加入候诊队列
func (q *WaitingQueue) AddEntry(checkIn *CheckIn, patientNameMasked, appointmentID string) *QueueEntry {
	num := q.nextQueueNumber()
	entry := QueueEntry{
		ID:                uuid.New().String(),
		QueueID:           q.ID,
		PatientID:         checkIn.PatientID,
		PatientNameMasked: patientNameMasked,
		AppointmentID:     appointmentID,
		CheckInID:         checkIn.ID,
		QueueNumber:       num,
		Status:            EntryWaiting,
		EnteredAt:         checkIn.CheckInTime,
	}
	q.Entries = append(q.Entries, entry)
	return &q.Entries[len(q.Entries)-1]
}

func (q *WaitingQueue) nextQueueNumber() int {
	max := 0
	for _, e := range q.Entries {
		if e.QueueNumber > max {
			max = e.QueueNumber
		}
	}
	return max + 1
}

// CallNext 呼叫下一位候诊患者
func (q *WaitingQueue) CallNext() (*QueueEntry, error) {
	// 按签到时间排序
	waitingEntries := make([]*QueueEntry, 0)
	for i := range q.Entries {
		if q.Entries[i].Status == EntryWaiting {
			waitingEntries = append(waitingEntries, &q.Entries[i])
		}
	}
	if len(waitingEntries) == 0 {
		return nil, bizErr.New(bizErr.ErrTriageQueueEmpty)
	}
	sort.Slice(waitingEntries, func(i, j int) bool {
		return waitingEntries[i].EnteredAt.Before(waitingEntries[j].EnteredAt)
	})
	next := waitingEntries[0]
	now := time.Now()
	next.Status = EntryCalling
	next.CallCount++
	next.CalledAt = &now
	return next, nil
}

// Recall 重叫当前患者
func (q *WaitingQueue) Recall() (*QueueEntry, error) {
	for i := range q.Entries {
		if q.Entries[i].Status == EntryCalling {
			if q.Entries[i].CallCount >= MaxCallCount {
				return nil, bizErr.New(bizErr.ErrTriageRecallLimit)
			}
			now := time.Now()
			q.Entries[i].CallCount++
			q.Entries[i].CalledAt = &now
			return &q.Entries[i], nil
		}
	}
	return nil, bizErr.NewWithDetail(bizErr.ErrConflict, "当前无正在呼叫中的患者")
}

// MissAndRequeue 过号重排
func (q *WaitingQueue) MissAndRequeue() (*QueueEntry, error) {
	for i := range q.Entries {
		if q.Entries[i].Status == EntryCalling {
			q.Entries[i].MissCount++
			if q.Entries[i].MissCount >= MaxMissCount {
				q.Entries[i].Status = EntryNoShow
			} else {
				q.Entries[i].Status = EntryMissed
				// 重排：加入队列末尾
				q.Entries[i].EnteredAt = time.Now()
				q.Entries[i].Status = EntryWaiting
			}
			return &q.Entries[i], nil
		}
	}
	return nil, bizErr.NewWithDetail(bizErr.ErrConflict, "当前无正在呼叫中的患者")
}

// GetWaitCount 获取当前等候人数
func (q *WaitingQueue) GetWaitCount() int {
	count := 0
	for _, e := range q.Entries {
		if e.Status == EntryWaiting {
			count++
		}
	}
	return count
}

// EstimateWaitTime 预计等候分钟（简化：每人平10分钟）
func (q *WaitingQueue) EstimateWaitTime() int {
	return q.GetWaitCount() * 10
}

// QueueEntry 队列条目实体
type QueueEntry struct {
	ID                string      `json:"id"`
	QueueID           string      `json:"queue_id"`
	PatientID         string      `json:"patient_id"`
	PatientNameMasked string      `json:"patient_name_masked"`
	AppointmentID     string      `json:"appointment_id"`
	CheckInID         string      `json:"check_in_id"`
	QueueNumber       int         `json:"queue_number"`
	Status            EntryStatus `json:"status"`
	CallCount         int         `json:"call_count"`
	MissCount         int         `json:"miss_count"`
	EnteredAt         time.Time   `json:"entered_at"`
	CalledAt          *time.Time  `json:"called_at"`
	CompletedAt       *time.Time  `json:"completed_at"`
}

// ExamExecution 检查执行聚合根
type ExamExecution struct {
	ID                string     `json:"id"`
	AppointmentItemID string     `json:"appointment_item_id"`
	PatientID         string     `json:"patient_id"`
	DeviceID          string     `json:"device_id"`
	Status            ExamStatus `json:"status"`
	StartedAt         *time.Time `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	Duration          int        `json:"duration"` // 实际耗时分钟
	OperatorID        string     `json:"operator_id"`
	UndoDeadline      *time.Time `json:"undo_deadline"` // 撤销截止时间
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// NewExamExecution 创建检查执行记录
func NewExamExecution(appointmentItemID, patientID, deviceID string) *ExamExecution {
	now := time.Now()
	return &ExamExecution{
		ID:                uuid.New().String(),
		AppointmentItemID: appointmentItemID,
		PatientID:         patientID,
		DeviceID:          deviceID,
		Status:            ExamCheckedIn,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// Start 开始检查
func (e *ExamExecution) Start(operatorID string) error {
	if e.Status != ExamCheckedIn && e.Status != ExamWaiting {
		return bizErr.New(bizErr.ErrTriageStatusInvalid)
	}
	now := time.Now()
	deadline := now.Add(UndoWindowMinutes * time.Minute)
	e.Status = ExamOngoing
	e.StartedAt = &now
	e.OperatorID = operatorID
	e.UndoDeadline = &deadline
	e.UpdatedAt = now
	return nil
}

// Complete 完成检查
func (e *ExamExecution) Complete(operatorID string) error {
	if e.Status != ExamOngoing {
		return bizErr.New(bizErr.ErrTriageStatusInvalid)
	}
	now := time.Now()
	deadline := now.Add(UndoWindowMinutes * time.Minute)
	e.Status = ExamDone
	e.CompletedAt = &now
	if e.StartedAt != nil {
		e.Duration = int(now.Sub(*e.StartedAt).Minutes())
	}
	e.UndoDeadline = &deadline
	e.UpdatedAt = now
	return nil
}

// Undo 撤销误操作
func (e *ExamExecution) Undo(operatorID, reason string) error {
	now := time.Now()
	if e.UndoDeadline != nil && now.After(*e.UndoDeadline) {
		return bizErr.New(bizErr.ErrTriageUndoExpired)
	}
	switch e.Status {
	case ExamOngoing:
		e.Status = ExamCheckedIn
		e.StartedAt = nil
		e.UndoDeadline = nil
	case ExamDone:
		e.Status = ExamOngoing
		e.CompletedAt = nil
		e.Duration = 0
	default:
		return bizErr.New(bizErr.ErrTriageStatusInvalid)
	}
	e.UpdatedAt = now
	return nil
}
