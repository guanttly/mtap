package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	"github.com/euler/mtap/internal/infrastructure/external"
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
	switch cfg.Database.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSNString()), &gorm.Config{})
	default:
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSNString()), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	// 初始化 HIS 客户端
	hisClient := external.NewHISClient(
		cfg.HIS.BaseURL,
		cfg.HIS.APIKey,
		cfg.HIS.Timeout,
		cfg.HIS.Retry,
	)

	syncJob := job.NewSyncJob(db, hisClient)

	c := cron.New(cron.WithSeconds())
	// 增量同步：每5分钟
	c.AddFunc("0 */5 * * * *", syncJob.RunIncremental)
	// 全量同步：每日凌晨02:00
	c.AddFunc("0 0 2 * * *", syncJob.RunFull)

	// 启动时执行一次增量同步
	go func() {
		time.Sleep(3 * time.Second)
		log.Println("[SyncWorker] 启动增量同步...")
		syncJob.RunIncremental()
	}()

	c.Start()
	log.Println("[SyncWorker] HIS同步工作器已启动")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[SyncWorker] 收到退出信号，正在停止...")
	<-c.Stop().Done()
	log.Println("[SyncWorker] 已停止")
}
