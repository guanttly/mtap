// Package analytics 应用层 - DTO与领域对象装配器
package analytics

import (
	domain "github.com/euler/mtap/internal/domain/analytics"
)

// toDashboardResp 将领域快照对象转换为响应DTO
func toDashboardResp(snap *domain.DashboardSnapshot) *DashboardResp {
	devices := make([]DeviceStatusResp, 0, len(snap.DeviceStatus))
	for _, d := range snap.DeviceStatus {
		devices = append(devices, DeviceStatusResp{
			DeviceID: d.DeviceID, DeviceName: d.DeviceName,
			Status: d.Status, QueueCount: d.QueueCount,
		})
	}
	trend := make([]WaitTrendResp, 0, len(snap.WaitTrend))
	for _, t := range snap.WaitTrend {
		trend = append(trend, WaitTrendResp{Time: t.Time, AvgWaitMin: t.AvgWaitMin})
	}
	alerts := make([]AlertItemResp, 0, len(snap.Alerts))
	for _, a := range snap.Alerts {
		alerts = append(alerts, AlertItemResp{
			Type: a.Type, Message: a.Message, DeviceID: a.DeviceID, Value: a.Value,
		})
	}
	return &DashboardResp{
		Timestamp: snap.Timestamp,
		CampusID:  snap.CampusID,
		SlotUsage: SlotUsageResp{
			TotalSlots:     snap.SlotUsage.TotalSlots,
			UsedSlots:      snap.SlotUsage.UsedSlots,
			ExpiredSlots:   snap.SlotUsage.ExpiredSlots,
			AvailableSlots: snap.SlotUsage.AvailableSlots,
			UsageRate:      snap.SlotUsage.UsageRate,
		},
		DeviceStatus: devices,
		WaitTrend:    trend,
		Alerts:       alerts,
	}
}

// toDeviceDetailResp 将领域设备详情对象转换为响应DTO
func toDeviceDetailResp(d *domain.DeviceDetail) *DeviceDetailResp {
	slots := make([]SlotSummaryResp, 0, len(d.TimeSlots))
	for _, s := range d.TimeSlots {
		slots = append(slots, SlotSummaryResp{
			Hour: s.Hour, Total: s.Total, Used: s.Used, UsageRate: s.UsageRate,
		})
	}
	return &DeviceDetailResp{
		DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		Date:       d.Date.Format("2006-01-02"),
		TimeSlots:  slots,
		QueueSummary: QueueSummaryResp{
			TotalCheckedIn: d.QueueSummary.TotalCheckedIn,
			AvgWaitMin:     d.QueueSummary.AvgWaitMin,
			MaxWaitMin:     d.QueueSummary.MaxWaitMin,
			NoShowCount:    d.QueueSummary.NoShowCount,
		},
	}
}

// toReportResp 将领域报表对象转换为响应DTO
func toReportResp(r *domain.Report) *ReportResp {
	return &ReportResp{
		ID:          r.ID,
		ReportType:  r.ReportType,
		Status:      r.Status,
		Format:      r.Format,
		FilePath:    r.FilePath,
		FileSize:    r.FileSize,
		GeneratedAt: r.GeneratedAt,
		CreatedAt:   r.CreatedAt,
	}
}
