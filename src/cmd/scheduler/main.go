package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	domainOpt "github.com/euler/mtap/internal/domain/optimization"
	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	pOpt "github.com/euler/mtap/internal/infrastructure/persistence/optimization"
	"github.com/euler/mtap/internal/interfaces/job"
)

func main() {
	cfgPath := os.Getenv("MTAP_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/config.yaml"
	}
	cfg, err := infraConfig.Load(cfgPath)
	if err != nil {
		log.Printf("load config failed (using defaults): %v", err)
		cfg, _ = infraConfig.Load("")
	}

	// 初始化数据库
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(cfg.Database.DSNString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	// 初始化仓储
	optRepos := pOpt.NewRepositories(db)

	// 构建领域服务
	anomalySvc := domainOpt.NewAnomalyDetectionService(
		optRepos.MetricRepo(), optRepos.SnapshotRepo(), optRepos.AlertRepo(),
	)
	evalSvc := domainOpt.NewEvaluationService(
		optRepos.StrategyRepo(), optRepos.TrialRepo(), optRepos.EvalRepo(),
		optRepos.SnapshotRepo(), optRepos.MetricRepo(),
	)
	decaySvc := domainOpt.NewDecayTrackingService(
		optRepos.StrategyRepo(), optRepos.EvalRepo(), optRepos.DecayRepo(),
		optRepos.SnapshotRepo(), optRepos.MetricRepo(),
	)

	// 初始化 Job 实例
	slotReleaseJob := job.NewSlotReleaseJob(db)
	blacklistJob := job.NewBlacklistCleanupJob(db)
	perfScanJob := job.NewPerformanceScanJob(db)
	anomalyJob := job.NewAnomalyDetectionJob(anomalySvc)
	decayJob := job.NewStrategyDecayJob(decaySvc)
	trialMonitorJob := job.NewTrialMonitorJob(evalSvc, optRepos.StrategyRepo())

	// 注册定时任务（使用 cron with seconds 支持精确到秒的表达式）
	c := cron.New(cron.WithSeconds())

	// 号源释放：每5分钟
	c.AddFunc("0 */5 * * * *", slotReleaseJob.Run)
	// 黑名单清理：每日凌晨01:00
	c.AddFunc("0 0 1 * * *", blacklistJob.Run)
	// 异常检测：每小时整点
	c.AddFunc("0 0 * * * *", anomalyJob.Run)
	// 试运行监控：每小时05分
	c.AddFunc("0 5 * * * *", trialMonitorJob.Run)
	// 策略衰减检测：每日凌晨03:00
	c.AddFunc("0 0 3 * * *", decayJob.Run)
	// 周度效能扫描：每周一凌晨04:00
	c.AddFunc("0 0 4 * * 1", perfScanJob.Run)

	c.Start()
	log.Println("[Scheduler] 定时任务调度器已启动")
	for _, entry := range c.Entries() {
		log.Printf("  任务 ID=%d 下次执行: %v", entry.ID, entry.Next)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Scheduler] 收到退出信号，正在停止...")
	<-c.Stop().Done()
	log.Println("[Scheduler] 已停止")
}
