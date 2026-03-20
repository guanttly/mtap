// Package analytics 应用层 - 统计分析应用服务
package analytics

import (
	"context"
	"fmt"
	"time"

	domain "github.com/euler/mtap/internal/domain/analytics"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// AnalyticsAppService 统计分析应用服务
type AnalyticsAppService struct {
	dashboardSvc *domain.DashboardService
	reportSvc    *domain.ReportService
}

// NewAnalyticsAppService 创建应用服务
func NewAnalyticsAppService(
	dashRepo domain.DashboardRepository,
	reportRepo domain.ReportRepository,
) *AnalyticsAppService {
	return &AnalyticsAppService{
		dashboardSvc: domain.NewDashboardService(dashRepo),
		reportSvc:    domain.NewReportService(reportRepo),
	}
}

// GetDashboard 获取大屏实时数据
func (s *AnalyticsAppService) GetDashboard(ctx context.Context, req GetDashboardReq) (*DashboardResp, error) {
	snap, err := s.dashboardSvc.GetSnapshot(ctx, req.CampusID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return toDashboardResp(snap), nil
}

// GetDeviceDetail 获取设备当日详情
func (s *AnalyticsAppService) GetDeviceDetail(ctx context.Context, deviceID, dateStr string) (*DeviceDetailResp, error) {
	date := time.Now()
	if dateStr != "" {
		var err error
		date, err = time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date 格式应为 YYYY-MM-DD")
		}
	}
	detail, err := s.dashboardSvc.GetDeviceDetail(ctx, deviceID, date)
	if err != nil {
		return nil, err
	}
	return toDeviceDetailResp(detail), nil
}

// CreateReport 生成报表
func (s *AnalyticsAppService) CreateReport(ctx context.Context, req CreateReportReq, operatorID string) (*ReportResp, error) {
	start, err := time.ParseInLocation("2006-01-02", req.DateStart, time.Local)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date_start 格式应为 YYYY-MM-DD")
	}
	end, err := time.ParseInLocation("2006-01-02", req.DateEnd, time.Local)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date_end 格式应为 YYYY-MM-DD")
	}
	if end.Before(start) {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date_end 不能早于 date_start")
	}
	if end.Sub(start) > 730*24*time.Hour {
		return nil, bizErr.New(bizErr.ErrStatsRangeTooLong)
	}
	format := req.Format
	if format == "" {
		format = "xlsx"
	}
	report, err := s.reportSvc.Generate(ctx, domain.ReportInput{
		ReportType: req.ReportType,
		Dimensions: req.Metrics,
		DateRange:  domain.DateRange{Start: start, End: end},
		Format:     format,
		CreatedBy:  operatorID,
	})
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrStatsReportFailed, err)
	}
	return toReportResp(report), nil
}

// GetReport 获取报表详情
func (s *AnalyticsAppService) GetReport(ctx context.Context, id string) (*ReportResp, error) {
	report, err := s.reportSvc.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return toReportResp(report), nil
}

// ExportReport 导出报表文件
func (s *AnalyticsAppService) ExportReport(ctx context.Context, id string) (*ReportResp, []byte, error) {
	report, err := s.reportSvc.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if report == nil {
		return nil, nil, bizErr.New(bizErr.ErrNotFound)
	}
	if report.Status != "ready" {
		return toReportResp(report), nil, nil
	}
	// 简化实现：返回 CSV 内容占位
	data := []byte(fmt.Sprintf("id,type,date_start,date_end,status\n%s,%s,%s,%s,%s\n",
		report.ID, report.ReportType,
		report.DateRange.Start.Format("2006-01-02"),
		report.DateRange.End.Format("2006-01-02"),
		report.Status))
	return toReportResp(report), data, nil
}

// ListReports 获取报表列表
func (s *AnalyticsAppService) ListReports(ctx context.Context, page, size int) ([]*ReportResp, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	reports, total, err := s.reportSvc.List(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]*ReportResp, 0, len(reports))
	for _, r := range reports {
		resps = append(resps, toReportResp(r))
	}
	return resps, total, nil
}
