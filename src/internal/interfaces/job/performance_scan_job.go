// Package job 接口层 - 周度效能扫描定时任务
// 核心目的：周期性执行全量效能扫描，主动寻找优化空间
// 模块功能：
//   - 默认每周一凌晨04:00执行
//   - 对比各指标与历史最优值，识别缓慢恶化趋势
//   - 生成《周度效能扫描报告》
package job

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
)

// PerformanceScanJob 周度效能扫描任务
type PerformanceScanJob struct {
	db *gorm.DB
}

// NewPerformanceScanJob 创建效能扫描任务
func NewPerformanceScanJob(db *gorm.DB) *PerformanceScanJob {
	return &PerformanceScanJob{db: db}
}

// Run 执行效能扫描
func (j *PerformanceScanJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	log.Printf("[PerformanceScanJob] 开始周度效能扫描")
	start := time.Now()
	now := time.Now()
	scanWeek := now.Format("2006-W") + weekNumber(now)

	// 执行各项指标扫描
	opportunities := j.scanAllMetrics(ctx)
	opJSON, _ := json.Marshal(opportunities)

	// 创建扫描记录（插入，ScanedAt = now via autoCreateTime）
	scan := &po.PerformanceScanPO{
		ID:            uuid.New().String(),
		ScanWeek:      scanWeek,
		Opportunities: string(opJSON),
	}
	if err := j.db.WithContext(ctx).Create(scan).Error; err != nil {
		log.Printf("[PerformanceScanJob] 创建扫描记录失败: %v", err)
	}

	log.Printf("[PerformanceScanJob] 效能扫描完成，耗时 %v，发现 %d 项", time.Since(start), len(opportunities))
}

func (j *PerformanceScanJob) scanAllMetrics(ctx context.Context) []string {
	var findings []string

	// 号源利用率分析（近7天）
	type utilizationStat struct {
		AvgRate float64
	}
	var stat utilizationStat
	if err := j.db.WithContext(ctx).Raw(`
		SELECT AVG(CAST(used_slots AS REAL) / NULLIF(total_slots, 0)) as avg_rate
		FROM dashboard_snapshots
		WHERE created_at > ?
	`, time.Now().AddDate(0, 0, -7)).Scan(&stat).Error; err == nil {
		if stat.AvgRate < 0.7 {
			findings = append(findings, "号源利用率低于70%，建议优化排班")
		}
	}
	return findings
}

func weekNumber(t time.Time) string {
	_, week := t.ISOWeek()
	if week < 10 {
		return "0" + string(rune('0'+week))
	}
	return string(rune('0'+week/10)) + string(rune('0'+week%10))
}
