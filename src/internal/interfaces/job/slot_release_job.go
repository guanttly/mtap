// Package job 接口层 - 号源释放定时任务
// 核心目的：自动释放超时未确认/未缴费的号源
// 模块功能：
//   - 5分钟未确认号源自动释放
//   - 24小时未缴费订单号源自动释放
//   - 释放号源回流公共池
package job

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
)

// SlotReleaseJob 号源释放定时任务
type SlotReleaseJob struct {
	db *gorm.DB
}

// NewSlotReleaseJob 创建号源释放任务
func NewSlotReleaseJob(db *gorm.DB) *SlotReleaseJob {
	return &SlotReleaseJob{db: db}
}

// Run 执行号源释放
func (j *SlotReleaseJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	now := time.Now()

	// 1. 释放锁定超时的号源（lock_until < now）
	result := j.db.WithContext(ctx).
		Model(&po.TimeSlotPO{}).
		Where("status = ? AND lock_until < ?", "locked", now).
		Updates(map[string]any{
			"status":     "available",
			"locked_by":  "",
			"lock_until": nil,
			"remaining":  gorm.Expr("remaining + 1"),
			"updated_at": now,
		})
	if result.Error != nil {
		log.Printf("[SlotReleaseJob] 释放锁定超时号源失败: %v", result.Error)
	} else {
		log.Printf("[SlotReleaseJob] 释放锁定超时号源 %d 个", result.RowsAffected)
	}

	// 2. 取消24小时未缴费的预约，释放号源
	unpaidTimeout := now.Add(-24 * time.Hour)
	var unpaidAppts []po.AppointmentPO
	if err := j.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", "confirmed", unpaidTimeout).
		Find(&unpaidAppts).Error; err != nil {
		log.Printf("[SlotReleaseJob] 查询未缴费预约失败: %v", err)
		return
	}

	for _, appt := range unpaidAppts {
		if err := j.cancelAndReleaseSlots(ctx, appt, now); err != nil {
			log.Printf("[SlotReleaseJob] 取消预约 %s 失败: %v", appt.ID, err)
		}
	}
	log.Printf("[SlotReleaseJob] 处理超时未缴费预约 %d 个", len(unpaidAppts))
}

func (j *SlotReleaseJob) cancelAndReleaseSlots(ctx context.Context, appt po.AppointmentPO, now time.Time) error {
	return j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 取消预约
		if err := tx.Model(&po.AppointmentPO{}).
			Where("id = ?", appt.ID).
			Updates(map[string]any{
				"status":        "cancelled",
				"cancel_reason": "超时未缴费自动取消",
				"cancelled_at":  now,
				"updated_at":    now,
			}).Error; err != nil {
			return err
		}
		// 释放关联号源
		return tx.Model(&po.TimeSlotPO{}).
			Where("id IN (SELECT slot_id FROM appointment_items WHERE appointment_id = ?)", appt.ID).
			Updates(map[string]any{
				"status":     "available",
				"remaining":  gorm.Expr("remaining + 1"),
				"updated_at": now,
			}).Error
	})
}
