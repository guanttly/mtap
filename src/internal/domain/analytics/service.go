// Package analytics 统计分析领域 - 领域服务
package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ─── DashboardService ────────────────────────────────────────────────────────

// DashboardService 大屏数据服务
type DashboardService struct {
	dashRepo DashboardRepository
}

// NewDashboardService 创建大屏服务
func NewDashboardService(dashRepo DashboardRepository) *DashboardService {
	return &DashboardService{dashRepo: dashRepo}
}

// GetSnapshot 获取实时大屏快照（同时落库用于历史对比）
func (s *DashboardService) GetSnapshot(ctx context.Context, campusID string) (*DashboardSnapshot, error) {
	now := time.Now()
	usage, err := s.dashRepo.GetSlotUsage(ctx, campusID, now)
	if err != nil {
		return nil, err
	}
	devices, err := s.dashRepo.GetDeviceStatus(ctx, campusID)
	if err != nil {
		return nil, err
	}
	trend, err := s.dashRepo.GetWaitTrend(ctx, campusID, 12)
	if err != nil {
		return nil, err
	}

	var alerts []AlertItem
	for _, d := range devices {
		if d.QueueCount > 20 {
			alerts = append(alerts, AlertItem{
				Type:     "device_queue_overflow",
				Message:  fmt.Sprintf("设备 %s 等待队列已有 %d 人，超出阈值", d.DeviceName, d.QueueCount),
				DeviceID: d.DeviceID,
				Value:    d.QueueCount,
			})
		}
	}
	if usage.TotalSlots > 0 && float64(usage.AvailableSlots)/float64(usage.TotalSlots) < 0.05 {
		alerts = append(alerts, AlertItem{
			Type:    "slot_exhausted",
			Message: fmt.Sprintf("今日号源剩余不足5%%（剩余 %d 个）", usage.AvailableSlots),
		})
	}

	snap := &DashboardSnapshot{
		ID:           uuid.New().String(),
		CampusID:     campusID,
		Timestamp:    now,
		SlotUsage:    usage,
		DeviceStatus: devices,
		WaitTrend:    trend,
		Alerts:       alerts,
	}
	_ = s.dashRepo.SaveSnapshot(ctx, snap)
	return snap, nil
}

// GetDeviceDetail 获取设备当日详情
func (s *DashboardService) GetDeviceDetail(ctx context.Context, deviceID string, date time.Time) (*DeviceDetail, error) {
	return s.dashRepo.GetDeviceDetail(ctx, deviceID, date)
}

// ─── ReportService ────────────────────────────────────────────────────────────

// ReportService 报表服务
type ReportService struct {
	reportRepo ReportRepository
}

// NewReportService 创建报表服务
func NewReportService(reportRepo ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

// ReportInput 报表生成输入
type ReportInput struct {
	ReportType string
	Dimensions []string
	DateRange  DateRange
	Format     string
	CreatedBy  string
}

// Generate 生成报表
func (s *ReportService) Generate(ctx context.Context, input ReportInput) (*Report, error) {
	report := &Report{
		ID:         uuid.New().String(),
		ReportType: input.ReportType,
		Dimensions: input.Dimensions,
		DateRange:  input.DateRange,
		Status:     "generating",
		Format:     input.Format,
		CreatedBy:  input.CreatedBy,
		CreatedAt:  time.Now(),
	}
	if err := s.reportRepo.Save(ctx, report); err != nil {
		return nil, err
	}
	// 开发阶段同步标记 ready；生产应走异步队列
	now := time.Now()
	report.Status = "ready"
	report.GeneratedAt = &now
	_ = s.reportRepo.Update(ctx, report)
	return report, nil
}

// FindByID 查询报表详情
func (s *ReportService) FindByID(ctx context.Context, id string) (*Report, error) {
	return s.reportRepo.FindByID(ctx, id)
}

// List 查询报表列表
func (s *ReportService) List(ctx context.Context, page, size int) ([]*Report, int64, error) {
	return s.reportRepo.List(ctx, page, size)
}
