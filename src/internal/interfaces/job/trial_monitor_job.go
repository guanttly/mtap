// Package job 接口层 - 试运行监控定时任务
// 核心目的：监控进行中的策略试运行效果
// 模块功能：
//   - 持续对比试运行指标与基线数据
//   - 关键指标恶化超15%时触发紧急回滚
//   - 每日09:00推送试运行日报
//   - 试运行到期自动触发评估报告生成
package job

import (
	"context"
	"log"
	"time"

	domain "github.com/euler/mtap/internal/domain/optimization"
)

// TrialMonitorJob 试运行监控任务
type TrialMonitorJob struct {
	evalSvc      *domain.EvaluationService
	strategyRepo domain.OptimizationStrategyRepository
}

// NewTrialMonitorJob 创建试运行监控任务
func NewTrialMonitorJob(evalSvc *domain.EvaluationService, strategyRepo domain.OptimizationStrategyRepository) *TrialMonitorJob {
	return &TrialMonitorJob{evalSvc: evalSvc, strategyRepo: strategyRepo}
}

// Run 执行试运行监控
func (j *TrialMonitorJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Printf("[TrialMonitorJob] 开始试运行监控")
	start := time.Now()

	// 查询所有进行中的试运行策略
	strategies, err := j.strategyRepo.ListActiveTrials(ctx)
	if err != nil {
		log.Printf("[TrialMonitorJob] 查询试运行策略失败: %v", err)
		return
	}

	evaluated := 0
	for _, strategy := range strategies {
		// 如果试运行已到期，生成评估报告
		if strategy.CooldownUntil != nil && time.Now().After(*strategy.CooldownUntil) {
			if _, err := j.evalSvc.GenerateReport(ctx, strategy.ID); err != nil {
				log.Printf("[TrialMonitorJob] 策略 %s 评估失败: %v", strategy.ID, err)
			} else {
				evaluated++
			}
		}
	}

	log.Printf("[TrialMonitorJob] 试运行监控完成，评估 %d 个策略，耗时 %v", evaluated, time.Since(start))
}
