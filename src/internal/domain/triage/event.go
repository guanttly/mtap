// Package triage 分诊执行领域 - 领域事件
package triage

import "time"

// EventType 分诊事件类型
type EventType string

const (
	EventPatientCheckedIn EventType = "triage.checked_in"
	EventPatientCalled    EventType = "triage.called"
	EventPatientMissed    EventType = "triage.missed"
	EventExamStarted      EventType = "triage.exam_started"
	EventExamCompleted    EventType = "triage.exam_completed"
)

// ScreenPushMessage 分诊大屏推送消息
type ScreenPushMessage struct {
	Type string      `json:"type"` // call / recall / miss / queue_update
	Data interface{} `json:"data"`
	At   time.Time   `json:"at"`
}

// CallPushData 呼叫信息数据
type CallPushData struct {
	PatientNameMasked string `json:"patient_name"`
	Room              string `json:"room"`
	QueueNumber       int    `json:"queue_number"`
	Message           string `json:"message"`
	CallCount         int    `json:"call_count"`
}

// BuildCallMessage 构建呼叫推送消息
func BuildCallMessage(entry *QueueEntry, roomName string) *ScreenPushMessage {
	msgType := "call"
	if entry.CallCount > 1 {
		msgType = "recall"
	}
	return &ScreenPushMessage{
		Type: msgType,
		Data: CallPushData{
			PatientNameMasked: entry.PatientNameMasked,
			Room:              roomName,
			QueueNumber:       entry.QueueNumber,
			Message:           "请 " + entry.PatientNameMasked + " 到" + roomName,
			CallCount:         entry.CallCount,
		},
		At: time.Now(),
	}
}
