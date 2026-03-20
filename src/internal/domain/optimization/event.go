// Package optimization 智能效能优化领域 - 领域事件
package optimization

import "time"

// EventType 领域事件类型
type EventType string

const (
	EvtBottleneckDetected       EventType = "bottleneck_detected"
	EvtStrategyApproved         EventType = "strategy_approved"
	EvtTrialStarted             EventType = "trial_started"
	EvtTrialEmergencyRollback   EventType = "trial_emergency_rollback"
	EvtTrialCompleted           EventType = "trial_completed"
	EvtStrategyPromoted         EventType = "strategy_promoted"
	EvtStrategyDecayed          EventType = "strategy_decayed"
	EvtPerformanceScanCompleted EventType = "performance_scan_completed"
)

// OptimizationEvent 效能优化领域事件
type OptimizationEvent struct {
	Type       EventType
	StrategyID string
	AlertID    string
	OccurredAt time.Time
	Payload    map[string]interface{}
}
