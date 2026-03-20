// Package resource 资源管理领域 - 仓储接口
// 核心目的：定义资源管理聚合根的持久化接口
// 模块功能：
//   - DeviceRepository / ScheduleRepository / TimeSlotRepository
//   - SlotPoolRepository / ExamItemRepository
//   - CampusRepository / DepartmentRepository / DoctorRepository
package resource

import (
	"context"
	"time"
)

// DeviceRepository 设备聚合根仓储接口
type DeviceRepository interface {
	FindByID(ctx context.Context, id string) (*Device, error)
	FindAll(ctx context.Context) ([]*Device, error)
	FindByDepartment(ctx context.Context, deptID string) ([]*Device, error)
	Save(ctx context.Context, d *Device) error
	Delete(ctx context.Context, id string) error
}

// ExamItemRepository 检查项目聚合根仓储接口
type ExamItemRepository interface {
	FindByID(ctx context.Context, id string) (*ExamItem, error)
	FindAll(ctx context.Context) ([]*ExamItem, error)
	FindByIDs(ctx context.Context, ids []string) ([]*ExamItem, error)
	Save(ctx context.Context, e *ExamItem) error
	Delete(ctx context.Context, id string) error
}

// ScheduleRepository 排班聚合根仓储接口
type ScheduleRepository interface {
	FindByID(ctx context.Context, id string) (*Schedule, error)
	FindByDeviceAndDate(ctx context.Context, deviceID string, date time.Time) ([]*Schedule, error)
	FindByDateRange(ctx context.Context, deviceID string, start, end time.Time) ([]*Schedule, error)
	Save(ctx context.Context, s *Schedule) error
	Suspend(ctx context.Context, id string, reason string) error
	UpdateDevice(ctx context.Context, id string, newDeviceID string) error
}

// TimeSlotRepository 号源时段仓储接口
type TimeSlotRepository interface {
	FindByID(ctx context.Context, id string) (*TimeSlot, error)
	FindBySchedule(ctx context.Context, scheduleID string) ([]*TimeSlot, error)
	QueryAvailable(ctx context.Context, deviceID string, date time.Time, examItemID, poolType string) ([]*TimeSlot, error)
	BulkSave(ctx context.Context, slots []*TimeSlot) error
	Lock(ctx context.Context, slotID, patientID string, until time.Time) error
	Release(ctx context.Context, slotID, patientID string, allowForce bool) error
	SuspendRange(ctx context.Context, deviceID string, date time.Time, startAt, endAt time.Time, reason string) (int, error)
	UpdateDevice(ctx context.Context, sourceDeviceID, targetDeviceID string, date time.Time) (int, error)
}

// SlotPoolRepository 号源池聚合根仓储接口
type SlotPoolRepository interface {
	FindByID(ctx context.Context, id string) (*SlotPool, error)
	FindAll(ctx context.Context) ([]*SlotPool, error)
	FindByDepartment(ctx context.Context, deptID string) ([]*SlotPool, error)
	Save(ctx context.Context, p *SlotPool) error
	Delete(ctx context.Context, id string) error
}

// CampusRepository 院区实体仓储接口
type CampusRepository interface {
	FindAll(ctx context.Context) ([]*Campus, error)
	FindByID(ctx context.Context, id string) (*Campus, error)
	Save(ctx context.Context, c *Campus) error
}

// DepartmentRepository 科室实体仓储接口
type DepartmentRepository interface {
	FindAll(ctx context.Context) ([]*Department, error)
	FindByCampus(ctx context.Context, campusID string) ([]*Department, error)
	FindByID(ctx context.Context, id string) (*Department, error)
	Save(ctx context.Context, d *Department) error
}

// DoctorRepository 医生实体仓储接口
type DoctorRepository interface {
	FindByID(ctx context.Context, id string) (*Doctor, error)
	FindByDepartment(ctx context.Context, deptID string) ([]*Doctor, error)
	Save(ctx context.Context, d *Doctor) error
}
