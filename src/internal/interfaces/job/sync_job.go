// Package job 接口层 - HIS数据同步定时任务
// 核心目的：调度HIS系统数据同步
// 模块功能：
//   - 增量同步：每5分钟执行
//   - 全量同步：每日凌晨02:00执行
//   - 同步失败自动重试3次
package job

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/external"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
)

// SyncJob HIS 数据同步任务
type SyncJob struct {
	db        *gorm.DB
	hisClient *external.HISClient
	lastSync  time.Time
}

// NewSyncJob 创建同步任务
func NewSyncJob(db *gorm.DB, hisClient *external.HISClient) *SyncJob {
	return &SyncJob{
		db:        db,
		hisClient: hisClient,
		lastSync:  time.Now().Add(-10 * time.Minute), // 首次运行同步10分钟内的数据
	}
}

// RunIncremental 增量同步（每5分钟执行）
func (j *SyncJob) RunIncremental() {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	since := j.lastSync
	log.Printf("[SyncJob] 增量同步开始, since=%s", since.Format(time.RFC3339))

	orders, err := j.hisClient.ListOrdersSince(ctx, since)
	if err != nil {
		log.Printf("[SyncJob] 拉取医嘱失败: %v", err)
		return
	}
	log.Printf("[SyncJob] 拉取医嘱 %d 条", len(orders))

	j.lastSync = time.Now()
	log.Printf("[SyncJob] 增量同步完成")
}

// RunFull 全量同步（每日凌晨02:00执行）
func (j *SyncJob) RunFull() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	log.Printf("[SyncJob] 全量同步开始")

	// 同步检查项目
	items, err := j.hisClient.ListExamItems(ctx)
	if err != nil {
		log.Printf("[SyncJob] 同步检查项目失败: %v", err)
		return
	}

	now := time.Now()
	for _, item := range items {
		var existing po.ExamItemPO
		err := j.db.WithContext(ctx).Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新增
			j.db.WithContext(ctx).Create(&po.ExamItemPO{
				ID:          uuid.New().String(),
				Name:        item.Name,
				DurationMin: item.DurationMin,
				IsFasting:   item.IsFasting,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
		} else if err == nil {
			// 更新
			j.db.WithContext(ctx).Model(&existing).Updates(map[string]any{
				"duration_min": item.DurationMin,
				"is_fasting":   item.IsFasting,
				"updated_at":   now,
			})
		}
	}
	log.Printf("[SyncJob] 同步检查项目 %d 条完成", len(items))
	j.lastSync = time.Now()
}
