package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/euler/mtap/internal/domain/analytics"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// ============================================================
// Mock Repositories
// ============================================================

type mockDashboardRepo struct {
	slotUsage    domain.SlotUsageData
	devices      []domain.DeviceStatusData
	waitTrend    []domain.WaitTrendPoint
	deviceDetail *domain.DeviceDetail
}

func newMockDashboardRepo() *mockDashboardRepo {
	return &mockDashboardRepo{
		slotUsage: domain.SlotUsageData{
			TotalSlots:     100,
			UsedSlots:      60,
			ExpiredSlots:   5,
			AvailableSlots: 35,
			UsageRate:      0.60,
		},
		devices: []domain.DeviceStatusData{
			{DeviceID: "DEV001", DeviceName: "CT-1", Status: "in_use", QueueCount: 5},
			{DeviceID: "DEV002", DeviceName: "MRI-1", Status: "idle", QueueCount: 0},
		},
		waitTrend: []domain.WaitTrendPoint{
			{Time: time.Now().Add(-30 * time.Minute), AvgWaitMin: 12.5},
			{Time: time.Now(), AvgWaitMin: 10.0},
		},
		deviceDetail: &domain.DeviceDetail{
			DeviceID:   "DEV001",
			DeviceName: "CT-1",
			Date:       time.Now(),
			TimeSlots: []domain.SlotSummary{
				{Hour: 8, Total: 10, Used: 8, UsageRate: 0.8},
				{Hour: 9, Total: 10, Used: 9, UsageRate: 0.9},
			},
			QueueSummary: domain.QueueSummary{
				TotalCheckedIn: 17,
				AvgWaitMin:     11.0,
				MaxWaitMin:     25.0,
				NoShowCount:    1,
			},
		},
	}
}

func (m *mockDashboardRepo) GetSlotUsage(_ context.Context, _ string, _ time.Time) (domain.SlotUsageData, error) {
	return m.slotUsage, nil
}
func (m *mockDashboardRepo) GetDeviceStatus(_ context.Context, _ string) ([]domain.DeviceStatusData, error) {
	return m.devices, nil
}
func (m *mockDashboardRepo) GetWaitTrend(_ context.Context, _ string, _ int) ([]domain.WaitTrendPoint, error) {
	return m.waitTrend, nil
}
func (m *mockDashboardRepo) SaveSnapshot(_ context.Context, _ *domain.DashboardSnapshot) error {
	return nil
}
func (m *mockDashboardRepo) GetDeviceDetail(_ context.Context, _ string, _ time.Time) (*domain.DeviceDetail, error) {
	return m.deviceDetail, nil
}

type mockReportRepo struct {
	items map[string]*domain.Report
}

func newMockReportRepo() *mockReportRepo {
	return &mockReportRepo{items: make(map[string]*domain.Report)}
}

func (m *mockReportRepo) Save(_ context.Context, r *domain.Report) error {
	m.items[r.ID] = r
	return nil
}
func (m *mockReportRepo) FindByID(_ context.Context, id string) (*domain.Report, error) {
	return m.items[id], nil
}
func (m *mockReportRepo) List(_ context.Context, page, size int) ([]*domain.Report, int64, error) {
	all := make([]*domain.Report, 0, len(m.items))
	for _, r := range m.items {
		all = append(all, r)
	}
	total := int64(len(all))
	start := (page - 1) * size
	if start >= len(all) {
		return []*domain.Report{}, total, nil
	}
	end := start + size
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}
func (m *mockReportRepo) Update(_ context.Context, r *domain.Report) error {
	m.items[r.ID] = r
	return nil
}

// ============================================================
// Helper
// ============================================================

func newTestAnalyticsService() *AnalyticsAppService {
	return NewAnalyticsAppService(
		newMockDashboardRepo(),
		newMockReportRepo(),
	)
}

// ============================================================
// GetDashboard Tests
// ============================================================

func TestGetDashboard_OK(t *testing.T) {
	svc := newTestAnalyticsService()
	resp, err := svc.GetDashboard(context.Background(), GetDashboardReq{CampusID: "CAMPUS001"})
	require.NoError(t, err)
	assert.Equal(t, "CAMPUS001", resp.CampusID)
	assert.Equal(t, 100, resp.SlotUsage.TotalSlots)
	assert.Equal(t, 60, resp.SlotUsage.UsedSlots)
	assert.Len(t, resp.DeviceStatus, 2)
	assert.Len(t, resp.WaitTrend, 2)
}

func TestGetDashboard_AlertsGenerated(t *testing.T) {
	// 构造设备队列超过阈值的场景
	dashRepo := newMockDashboardRepo()
	dashRepo.devices = []domain.DeviceStatusData{
		{DeviceID: "DEV-OVF", DeviceName: "CT-OVF", Status: "in_use", QueueCount: 25},
	}
	// 号源剩余不足5%
	dashRepo.slotUsage = domain.SlotUsageData{
		TotalSlots:     100,
		UsedSlots:      96,
		AvailableSlots: 4,
		UsageRate:      0.96,
	}
	svc := NewAnalyticsAppService(dashRepo, newMockReportRepo())

	resp, err := svc.GetDashboard(context.Background(), GetDashboardReq{})
	require.NoError(t, err)
	// 应包含队列超限和号源不足两条告警
	assert.GreaterOrEqual(t, len(resp.Alerts), 2)

	types := make(map[string]bool)
	for _, a := range resp.Alerts {
		types[a.Type] = true
	}
	assert.True(t, types["device_queue_overflow"])
	assert.True(t, types["slot_exhausted"])
}

// ============================================================
// GetDeviceDetail Tests
// ============================================================

func TestGetDeviceDetail_OK(t *testing.T) {
	svc := newTestAnalyticsService()
	detail, err := svc.GetDeviceDetail(context.Background(), "DEV001", "")
	require.NoError(t, err)
	assert.Equal(t, "DEV001", detail.DeviceID)
	assert.Equal(t, "CT-1", detail.DeviceName)
	assert.Len(t, detail.TimeSlots, 2)
}

func TestGetDeviceDetail_InvalidDate(t *testing.T) {
	svc := newTestAnalyticsService()
	_, err := svc.GetDeviceDetail(context.Background(), "DEV001", "not-a-date")
	assert.True(t, bizErr.Is(err, bizErr.ErrInvalidParam))
}

// ============================================================
// CreateReport Tests
// ============================================================

func TestCreateReport_OK(t *testing.T) {
	svc := newTestAnalyticsService()
	resp, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "daily",
		DateStart:  "2026-03-01",
		DateEnd:    "2026-03-15",
		Format:     "xlsx",
	}, "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "daily", resp.ReportType)
	assert.Equal(t, "xlsx", resp.Format)
}

func TestCreateReport_InvalidDateFormat(t *testing.T) {
	svc := newTestAnalyticsService()
	_, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "daily",
		DateStart:  "20260301",
		DateEnd:    "2026-03-15",
	}, "admin")
	assert.True(t, bizErr.Is(err, bizErr.ErrInvalidParam))
}

func TestCreateReport_EndBeforeStart(t *testing.T) {
	svc := newTestAnalyticsService()
	_, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "daily",
		DateStart:  "2026-03-15",
		DateEnd:    "2026-03-01",
	}, "admin")
	assert.True(t, bizErr.Is(err, bizErr.ErrInvalidParam))
}

func TestCreateReport_RangeTooLong(t *testing.T) {
	svc := newTestAnalyticsService()
	_, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "monthly",
		DateStart:  "2024-01-01",
		DateEnd:    "2026-03-20",
	}, "admin")
	assert.True(t, bizErr.Is(err, bizErr.ErrStatsRangeTooLong))
}

// ============================================================
// GetReport Tests
// ============================================================

func TestGetReport_OK(t *testing.T) {
	svc := newTestAnalyticsService()
	created, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "weekly",
		DateStart:  "2026-03-01",
		DateEnd:    "2026-03-07",
	}, "admin")
	require.NoError(t, err)

	got, err := svc.GetReport(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "weekly", got.ReportType)
}

func TestGetReport_NotFound(t *testing.T) {
	svc := newTestAnalyticsService()
	_, err := svc.GetReport(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}

// ============================================================
// ListReports Tests
// ============================================================

func TestListReports_OK(t *testing.T) {
	svc := newTestAnalyticsService()
	_, _ = svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "daily", DateStart: "2026-03-01", DateEnd: "2026-03-02",
	}, "admin")
	_, _ = svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "weekly", DateStart: "2026-03-01", DateEnd: "2026-03-07",
	}, "admin")

	list, total, err := svc.ListReports(context.Background(), 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

// ============================================================
// ExportReport Tests
// ============================================================

func TestExportReport_Ready(t *testing.T) {
	svc := newTestAnalyticsService()
	created, err := svc.CreateReport(context.Background(), CreateReportReq{
		ReportType: "daily", DateStart: "2026-03-01", DateEnd: "2026-03-02",
	}, "admin")
	require.NoError(t, err)

	// 开发阶段 Generate 同步标记 ready，导出应返回 CSV 数据
	resp, data, err := svc.ExportReport(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
	assert.Equal(t, "ready", resp.Status)
	assert.NotEmpty(t, data) // ready 状态返回文件内容
}

func TestExportReport_NotFound(t *testing.T) {
	svc := newTestAnalyticsService()
	_, _, err := svc.ExportReport(context.Background(), "nonexistent")
	assert.True(t, bizErr.Is(err, bizErr.ErrNotFound))
}
