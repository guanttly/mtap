// Package po 基础设施层 - analytics 持久化对象
package po

import "time"

// DashboardSnapshotPO 大屏快照持久化对象
type DashboardSnapshotPO struct {
	ID        string    `gorm:"primaryKey;column:id;size:36"`
	CampusID  string    `gorm:"column:campus_id;size:36;index:idx_dashboard_campus_time"`
	Snapshot  string    `gorm:"column:snapshot;type:text"` // JSON
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_dashboard_campus_time"`
}

func (DashboardSnapshotPO) TableName() string { return "dashboard_snapshots" }

// ReportPO 报表持久化对象
type ReportPO struct {
	ID          string     `gorm:"primaryKey;column:id;size:36"`
	ReportType  string     `gorm:"column:report_type;size:10;not null"`
	Dimensions  string     `gorm:"column:dimensions;type:text"` // JSON
	DateStart   time.Time  `gorm:"column:date_start;not null"`
	DateEnd     time.Time  `gorm:"column:date_end;not null"`
	Status      string     `gorm:"column:status;size:15;not null;default:generating"`
	Format      string     `gorm:"column:format;size:5;not null;default:xlsx"`
	FilePath    string     `gorm:"column:file_path;size:500"`
	FileSize    int64      `gorm:"column:file_size;default:0"`
	GeneratedAt *time.Time `gorm:"column:generated_at"`
	CreatedBy   string     `gorm:"column:created_by;size:36;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (ReportPO) TableName() string { return "reports" }
