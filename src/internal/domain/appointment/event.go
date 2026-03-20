// Package appointment 预约服务领域 - 领域事件
package appointment

import "time"

// EventType 预约事件类型
type EventType string

const (
	EventAppointmentCreated     EventType = "appointment.created"
	EventAppointmentConfirmed   EventType = "appointment.confirmed"
	EventAppointmentCancelled   EventType = "appointment.cancelled"
	EventAppointmentRescheduled EventType = "appointment.rescheduled"
	EventAppointmentCompleted   EventType = "appointment.completed"
	EventNoShowRecorded         EventType = "appointment.no_show"
	EventBlacklistCreated       EventType = "appointment.blacklist_created"
)

// DomainEvent 领域事件基类
type DomainEvent struct {
	Type       EventType   `json:"type"`
	Payload    interface{} `json:"payload"`
	OccurredAt time.Time   `json:"occurred_at"`
}

// AppointmentCreatedEvent 预约创建事件
type AppointmentCreatedEvent struct {
	AppointmentID string `json:"appointment_id"`
	PatientID     string `json:"patient_id"`
	Mode          string `json:"mode"`
}

// AppointmentCancelledEvent 预约取消事件
type AppointmentCancelledEvent struct {
	AppointmentID string `json:"appointment_id"`
	PatientID     string `json:"patient_id"`
	Reason        string `json:"reason"`
}

// NoShowRecordedEvent 爽约记录事件
type NoShowRecordedEvent struct {
	PatientID     string `json:"patient_id"`
	AppointmentID string `json:"appointment_id"`
	NoShowCount   int    `json:"no_show_count"`
	Blacklisted   bool   `json:"blacklisted"`
}
