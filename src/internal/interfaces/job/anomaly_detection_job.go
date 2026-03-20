// Package job 接口层 - 异常检测定时任务
// 核心目的：周期性执行全指标异常检测扫描
// 模块功能：
//   - 每小时执行一轮全指标扫描（耗时不超5分钟）
//   - 均值±2σ偏离检测 + 环比突变检测 + 趋势恶化检测
//   - 发现异常自动生成瓶颈告警
package job

import (
	"context"
	"log"
	"time"

	domain "github.com/euler/mtap/internal/domain/optimization"
)

// AnomalyDetectionJob 异常检测定时任务
type AnomalyDetectionJob struct {
	anomalySvc *domain.AnomalyDetectionService
}

// NewAnomalyDetectionJob 创建异常检测任务
func NewAnomalyDetectionJob(anomalySvc *domain.AnomalyDetectionService) *AnomalyDetectionJob {
	return &AnomalyDetectionJob{anomalySvc: anomalySvc}
}

// Run 执行异常检测扫描
func (j *AnomalyDetectionJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Printf("[AnomalyDetectionJob] 开始全指标异常扫描")
	start := time.Now()

	alerts, err := j.anomalySvc.RunFullScan(ctx)
	if err != nil {
		log.Printf("[AnomalyDetectionJob] 扫描失败: %v", err)
		return
	}

	log.Printf("[AnomalyDetectionJob] 扫描完成，发现告警 %d 条，耗时 %v", len(alerts), time.Since(start))
}
