// Package job 接口层 - 未缴费预约自动取消定时任务
// 核心目的：每小时扫描超过24小时未缴费的 pending 预约，自动取消并释放号源
package job

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
)

// PaymentTimeoutReleaseJob 未缴费预约自动取消任务
type PaymentTimeoutReleaseJob struct {
	db *gorm.DB
}

// NewPaymentTimeoutReleaseJob 创建未缴费自动取消任务
func NewPaymentTimeoutReleaseJob(db *gorm.DB) *PaymentTimeoutReleaseJob {
	return &PaymentTimeoutReleaseJob{db: db}
}

// Run 执行未缴费预约自动取消
// 逻辑：status='pending' 且 created_at 超过24小时且 payment_verified=0 的预约自动取消
func (j *PaymentTimeoutReleaseJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	now := time.Now()
	paymentTimeout := now.Add(-24 * time.Hour)

	// 查询超时未缴费的预约（pending + 超过24小时 + payment_verified=0）
	var unpaidAppts []po.AppointmentPO
	if err := j.db.WithContext(ctx).
		Where("status = ? AND payment_verified = ? AND created_at < ?", "pending", false, paymentTimeout).
		Find(&unpaidAppts).Error; err != nil {
		log.Printf("[PaymentTimeoutReleaseJob] 查询超时未缴费预约失败: %v", err)
		return
	}

	if len(unpaidAppts) == 0 {
		log.Printf("[PaymentTimeoutReleaseJob] 无超时未缴费预约")
		return
	}

	successCount := 0
	for _, appt := range unpaidAppts {
		if err := j.cancelAndReleaseSlots(ctx, appt, now); err != nil {
			log.Printf("[PaymentTimeoutReleaseJob] 取消预约 %s 失败: %v", appt.ID, err)
		} else {
			successCount++
		}
	}
	log.Printf("[PaymentTimeoutReleaseJob] 处理超时未缴费预约: 总计 %d 个，成功取消 %d 个", len(unpaidAppts), successCount)
}

// cancelAndReleaseSlots 事务内取消预约并释放关联号源
func (j *PaymentTimeoutReleaseJob) cancelAndReleaseSlots(ctx context.Context, appt po.AppointmentPO, now time.Time) error {
	return j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 取消预约单
		if err := tx.Model(&po.AppointmentPO{}).
			Where("id = ? AND status = ?", appt.ID, "pending").
			Updates(map[string]any{
				"status":        "cancelled",
				"cancel_reason": "超过24小时未缴费自动取消",
				"cancelled_at":  now,
				"updated_at":    now,
			}).Error; err != nil {
			return err
		}

		// 查找关联预约项目
		var items []po.AppointmentItemPO
		if err := tx.Where("appointment_id = ?", appt.ID).Find(&items).Error; err != nil {
			return err
		}

		// 取消预约项目并释放号源
		for _, item := range items {
			// 标记预约项目为已取消
			if err := tx.Model(&po.AppointmentItemPO{}).
				Where("id = ?", item.ID).
				Updates(map[string]any{
					"status":     "cancelled",
					"updated_at": now,
				}).Error; err != nil {
				return err
			}

			// 释放对应号源（设为 available，清除锁信息）
			if err := tx.Model(&po.TimeSlotPO{}).
				Where("id = ? AND status IN ?", item.SlotID, []string{"locked", "booked"}).
				Updates(map[string]any{
					"status":     "available",
					"locked_by":  "",
					"lock_until": nil,
					"remaining":  gorm.Expr("remaining + 1"),
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		}

		log.Printf("[PaymentTimeoutReleaseJob] 已取消预约 %s，释放 %d 个号源", appt.ID, len(items))
		return nil
	})
}
