// Package analytics 基础设施层 - analytics 仓储实现（GORM）
package analytics

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	domain "github.com/euler/mtap/internal/domain/analytics"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// Repositories analytics 模块仓储集合
type Repositories struct {
	db *gorm.DB
}

// NewRepositories 创建仓储集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{db: db}
}

// DashboardRepo 返回大屏数据仓储
func (r *Repositories) DashboardRepo() domain.DashboardRepository {
	return &dashboardRepo{db: r.db}
}

// ReportRepo 返回报表仓储
func (r *Repositories) ReportRepo() domain.ReportRepository {
	return &reportRepo{db: r.db}
}

// ─── dashboardRepo ────────────────────────────────────────────────────────────

type dashboardRepo struct{ db *gorm.DB }

func (r *dashboardRepo) GetSlotUsage(ctx context.Context, campusID string, date time.Time) (domain.SlotUsageData, error) {
	type row struct {
		Status string
		Count  int
	}
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var rows []row
	if err := r.db.WithContext(ctx).
		Raw(`SELECT status, COUNT(*) as count FROM time_slots
             WHERE date >= ? AND date < ? GROUP BY status`, dayStart, dayEnd).
		Scan(&rows).Error; err != nil {
		return domain.SlotUsageData{}, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	var data domain.SlotUsageData
	for _, r := range rows {
		data.TotalSlots += r.Count
		switch r.Status {
		case "booked", "locked", "checked_in", "examining", "completed":
			data.UsedSlots += r.Count
		case "expired", "no_show":
			data.ExpiredSlots += r.Count
		case "available":
			data.AvailableSlots += r.Count
		}
	}
	if data.TotalSlots > 0 {
		data.UsageRate = float64(data.UsedSlots) / float64(data.TotalSlots)
	}
	return data, nil
}

func (r *dashboardRepo) GetDeviceStatus(ctx context.Context, campusID string) ([]domain.DeviceStatusData, error) {
	type deviceRow struct {
		ID         string
		Name       string
		Status     string
		QueueCount int
	}

	baseSQL := `SELECT d.id, d.name, d.status,
		(SELECT COUNT(*) FROM queue_entries qe
		 JOIN waiting_queues wq ON qe.queue_id = wq.id
		 WHERE wq.device_id = d.id AND qe.status IN ('waiting','calling')) AS queue_count
		FROM devices d WHERE d.status != 'inactive'`

	var rows []deviceRow
	var q *gorm.DB
	if campusID != "" {
		q = r.db.WithContext(ctx).Raw(baseSQL+` AND d.campus_id = ?`, campusID)
	} else {
		q = r.db.WithContext(ctx).Raw(baseSQL)
	}
	if err := q.Scan(&rows).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	result := make([]domain.DeviceStatusData, 0, len(rows))
	for _, row := range rows {
		status := "idle"
		if row.Status == "maintenance" {
			status = "maintenance"
		} else if row.QueueCount > 0 {
			status = "in_use"
		}
		result = append(result, domain.DeviceStatusData{
			DeviceID: row.ID, DeviceName: row.Name,
			Status: status, QueueCount: row.QueueCount,
		})
	}
	return result, nil
}

func (r *dashboardRepo) GetWaitTrend(ctx context.Context, campusID string, points int) ([]domain.WaitTrendPoint, error) {
	type trendRow struct {
		Hour       int
		AvgWaitMin float64
	}
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var rows []trendRow
	if err := r.db.WithContext(ctx).
		Raw(`SELECT HOUR(qe.entered_at) AS hour,
				 AVG(TIMESTAMPDIFF(SECOND, qe.entered_at, qe.called_at)/60.0) AS avg_wait_min
			 FROM queue_entries qe
			 WHERE qe.entered_at >= ? AND qe.called_at IS NOT NULL
			 GROUP BY hour ORDER BY hour DESC LIMIT ?`, dayStart, points).
		Scan(&rows).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	result := make([]domain.WaitTrendPoint, 0, len(rows))
	for _, row := range rows {
		t := time.Date(now.Year(), now.Month(), now.Day(), row.Hour, 0, 0, 0, now.Location())
		result = append(result, domain.WaitTrendPoint{Time: t, AvgWaitMin: row.AvgWaitMin})
	}
	return result, nil
}

func (r *dashboardRepo) SaveSnapshot(ctx context.Context, snap *domain.DashboardSnapshot) error {
	type payload struct {
		SlotUsage    domain.SlotUsageData      `json:"slot_usage"`
		DeviceStatus []domain.DeviceStatusData `json:"device_status"`
		WaitTrend    []domain.WaitTrendPoint   `json:"wait_trend"`
		Alerts       []domain.AlertItem        `json:"alerts"`
	}
	b, _ := json.Marshal(payload{
		SlotUsage: snap.SlotUsage, DeviceStatus: snap.DeviceStatus,
		WaitTrend: snap.WaitTrend, Alerts: snap.Alerts,
	})
	p := &po.DashboardSnapshotPO{
		ID: snap.ID, CampusID: snap.CampusID,
		Snapshot: string(b), CreatedAt: snap.Timestamp,
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *dashboardRepo) GetDeviceDetail(ctx context.Context, deviceID string, date time.Time) (*domain.DeviceDetail, error) {
	type devInfo struct {
		ID   string
		Name string
	}
	var dev devInfo
	if err := r.db.WithContext(ctx).
		Raw("SELECT id, name FROM devices WHERE id = ?", deviceID).Scan(&dev).Error; err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if dev.ID == "" {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	type slotRow struct {
		Hour  int
		Total int
		Used  int
	}
	var slotRows []slotRow
	_ = r.db.WithContext(ctx).
		Raw(`SELECT HOUR(start_at) AS hour,
				 COUNT(*) AS total,
				 SUM(CASE WHEN status NOT IN ('available') THEN 1 ELSE 0 END) AS used
			 FROM time_slots
			 WHERE device_id = ? AND date >= ? AND date < ?
			 GROUP BY hour ORDER BY hour`, deviceID, dayStart, dayEnd).
		Scan(&slotRows).Error

	slots := make([]domain.SlotSummary, 0, len(slotRows))
	for _, s := range slotRows {
		rate := 0.0
		if s.Total > 0 {
			rate = float64(s.Used) / float64(s.Total)
		}
		slots = append(slots, domain.SlotSummary{Hour: s.Hour, Total: s.Total, Used: s.Used, UsageRate: rate})
	}

	type qRow struct {
		TotalCheckedIn int
		AvgWaitMin     float64
		MaxWaitMin     float64
		NoShowCount    int
	}
	var qr qRow
	_ = r.db.WithContext(ctx).
		Raw(`SELECT COUNT(*) AS total_checked_in,
				 COALESCE(AVG(CASE WHEN qe.called_at IS NOT NULL
					 THEN TIMESTAMPDIFF(SECOND, qe.entered_at, qe.called_at)/60.0 END),0) AS avg_wait_min,
				 COALESCE(MAX(CASE WHEN qe.called_at IS NOT NULL
					 THEN TIMESTAMPDIFF(SECOND, qe.entered_at, qe.called_at)/60.0 END),0) AS max_wait_min,
				 SUM(CASE WHEN qe.status='no_show' THEN 1 ELSE 0 END) AS no_show_count
			 FROM queue_entries qe
			 JOIN waiting_queues wq ON qe.queue_id = wq.id
			 WHERE wq.device_id = ? AND qe.entered_at >= ? AND qe.entered_at < ?`,
			deviceID, dayStart, dayEnd).Scan(&qr).Error

	return &domain.DeviceDetail{
		DeviceID: dev.ID, DeviceName: dev.Name, Date: date,
		TimeSlots: slots,
		QueueSummary: domain.QueueSummary{
			TotalCheckedIn: qr.TotalCheckedIn,
			AvgWaitMin:     qr.AvgWaitMin, MaxWaitMin: qr.MaxWaitMin,
			NoShowCount: qr.NoShowCount,
		},
	}, nil
}

// ─── reportRepo ───────────────────────────────────────────────────────────────

type reportRepo struct{ db *gorm.DB }

func (r *reportRepo) Save(ctx context.Context, rep *domain.Report) error {
	dims, _ := json.Marshal(rep.Dimensions)
	p := &po.ReportPO{
		ID: rep.ID, ReportType: rep.ReportType,
		Dimensions: string(dims),
		DateStart:  rep.DateRange.Start, DateEnd: rep.DateRange.End,
		Status: rep.Status, Format: rep.Format,
		FilePath: rep.FilePath, FileSize: rep.FileSize,
		GeneratedAt: rep.GeneratedAt,
		CreatedBy:   rep.CreatedBy, CreatedAt: rep.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *reportRepo) FindByID(ctx context.Context, id string) (*domain.Report, error) {
	var p po.ReportPO
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return fromReportPO(&p), nil
}

func (r *reportRepo) List(ctx context.Context, page, size int) ([]*domain.Report, int64, error) {
	var total int64
	var pos []po.ReportPO
	db := r.db.WithContext(ctx).Model(&po.ReportPO{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if err := db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&pos).Error; err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	results := make([]*domain.Report, 0, len(pos))
	for i := range pos {
		results = append(results, fromReportPO(&pos[i]))
	}
	return results, total, nil
}

func (r *reportRepo) Update(ctx context.Context, rep *domain.Report) error {
	dims, _ := json.Marshal(rep.Dimensions)
	return r.db.WithContext(ctx).Model(&po.ReportPO{}).Where("id = ?", rep.ID).
		Updates(map[string]interface{}{
			"status":       rep.Status,
			"file_path":    rep.FilePath,
			"file_size":    rep.FileSize,
			"generated_at": rep.GeneratedAt,
			"dimensions":   string(dims),
		}).Error
}

func fromReportPO(p *po.ReportPO) *domain.Report {
	var dims []string
	_ = json.Unmarshal([]byte(p.Dimensions), &dims)
	return &domain.Report{
		ID: p.ID, ReportType: p.ReportType,
		Dimensions: dims,
		DateRange:  domain.DateRange{Start: p.DateStart, End: p.DateEnd},
		Status:     p.Status, Format: p.Format,
		FilePath: p.FilePath, FileSize: p.FileSize,
		GeneratedAt: p.GeneratedAt,
		CreatedBy:   p.CreatedBy, CreatedAt: p.CreatedAt,
	}
}
