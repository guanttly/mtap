// Package analytics 统计分析领域 - 仓储接口
package analytics

import (
	"context"
	"time"
)

// DashboardRepository 大屏数据仓储（跨上下文只读查询由基础设施层实现）
type DashboardRepository interface {
	GetSlotUsage(ctx context.Context, campusID string, date time.Time) (SlotUsageData, error)
	GetDeviceStatus(ctx context.Context, campusID string) ([]DeviceStatusData, error)
	GetWaitTrend(ctx context.Context, campusID string, points int) ([]WaitTrendPoint, error)
	SaveSnapshot(ctx context.Context, snap *DashboardSnapshot) error
	GetDeviceDetail(ctx context.Context, deviceID string, date time.Time) (*DeviceDetail, error)
}

// ReportRepository 报表仓储
type ReportRepository interface {
	Save(ctx context.Context, r *Report) error
	FindByID(ctx context.Context, id string) (*Report, error)
	List(ctx context.Context, page, size int) ([]*Report, int64, error)
	Update(ctx context.Context, r *Report) error
}
