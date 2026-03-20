// Package job 接口层 - 黑名单清理定时任务
// 核心目的：自动解除到期黑名单，校准爽约计数
// 模块功能：
//   - 过期黑名单自动解除（默认180天）
//   - 每日凌晨爽约计数校准（以预约+签到记录为准）
package job

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
)

// BlacklistCleanupJob 黑名单清理定时任务
type BlacklistCleanupJob struct {
	db *gorm.DB
}

// NewBlacklistCleanupJob 创建黑名单清理任务
func NewBlacklistCleanupJob(db *gorm.DB) *BlacklistCleanupJob {
	return &BlacklistCleanupJob{db: db}
}

// Run 执行黑名单清理
func (j *BlacklistCleanupJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	now := time.Now()

	// 1. 解除已到期的黑名单
	result := j.db.WithContext(ctx).
		Model(&po.BlacklistPO{}).
		Where("status = ? AND expire_at < ?", "active", now).
		Updates(map[string]any{
			"status":     "expired",
			"updated_at": now,
		})
	if result.Error != nil {
		log.Printf("[BlacklistCleanupJob] 解除过期黑名单失败: %v", result.Error)
	} else {
		log.Printf("[BlacklistCleanupJob] 解除过期黑名单 %d 条", result.RowsAffected)
	}

	// 2. 校准爽约计数（以90天内的爽约记录为准）
	window := now.AddDate(0, 0, -90)
	if err := j.recalibrateNoShowCounts(ctx, window, now); err != nil {
		log.Printf("[BlacklistCleanupJob] 爽约计数校准失败: %v", err)
	}
}

// recalibrateNoShowCounts 重新统计90天内爽约次数并更新黑名单表
func (j *BlacklistCleanupJob) recalibrateNoShowCounts(ctx context.Context, since, now time.Time) error {
	// 统计每个患者的爽约次数
	type noShowStat struct {
		PatientID string
		Count     int
	}
	var stats []noShowStat
	if err := j.db.WithContext(ctx).Raw(`
		SELECT patient_id, COUNT(*) as count
		FROM no_show_records
		WHERE occurred_at > ?
		GROUP BY patient_id
	`, since).Scan(&stats).Error; err != nil {
		return err
	}

	// 更新黑名单中的爽约计数
	for _, s := range stats {
		j.db.WithContext(ctx).
			Model(&po.BlacklistPO{}).
			Where("patient_id = ? AND status = ?", s.PatientID, "active").
			Update("no_show_count", s.Count)
	}

	log.Printf("[BlacklistCleanupJob] 校准爽约计数，涉及患者 %d 名", len(stats))
	return nil
}
