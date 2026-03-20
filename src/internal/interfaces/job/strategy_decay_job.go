// Package job 接口层 - 策略衰减检测定时任务
// 核心目的：追踪已转正/常态化策略的长期效果
// 模块功能：
//   - 转正后第30/90/180天自动生成长期效果评估报告
//   - 策略效果衰减超10%时触发衰减告警
//   - C类投产后第90天生成投产效果验证报告
package job

import (
	"context"
	"log"
	"time"

	domain "github.com/euler/mtap/internal/domain/optimization"
)

// StrategyDecayJob 策略衰减检测任务
type StrategyDecayJob struct {
	decaySvc *domain.DecayTrackingService
}

// NewStrategyDecayJob 创建策略衰减检测任务
func NewStrategyDecayJob(decaySvc *domain.DecayTrackingService) *StrategyDecayJob {
	return &StrategyDecayJob{decaySvc: decaySvc}
}

// Run 执行策略衰减检测
func (j *StrategyDecayJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	log.Printf("[StrategyDecayJob] 开始策略衰减检测")
	start := time.Now()

	alerts, err := j.decaySvc.CheckDecay(ctx)
	if err != nil {
		log.Printf("[StrategyDecayJob] 衰减检测失败: %v", err)
		return
	}

	log.Printf("[StrategyDecayJob] 衰减检测完成，发现衰减告警 %d 个，耗时 %v", len(alerts), time.Since(start))
}
