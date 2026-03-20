package main

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
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

	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(cfg.Database.DSNString()), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	log.Printf("数据库迁移开始，驱动: mysql")

	if err := db.AutoMigrate(
		// 用户权限
		&po.RolePO{},
		&po.UserPO{},
		// 规则引擎
		&po.ConflictRulePO{},
		&po.ConflictPackagePO{},
		&po.ConflictPackageItemPO{},
		&po.DependencyRulePO{},
		&po.PriorityTagPO{},
		&po.SortingStrategyPO{},
		&po.PatientAdaptRulePO{},
		&po.SourceControlPO{},
		// 资源管理
		&po.DevicePO{},
		&po.ExamItemPO{},
		&po.ItemAliasPO{},
		&po.SlotPoolPO{},
		&po.SchedulePO{},
		&po.TimeSlotPO{},
		// 预约服务
		&po.AppointmentPO{},
		&po.AppointmentItemPO{},
		&po.AppointmentCredentialPO{},
		&po.AppointmentChangeLogPO{},
		&po.BlacklistPO{},
		&po.NoShowRecordPO{},
		&po.AppealPO{},
		// 分诊执行
		&po.CheckInPO{},
		&po.WaitingQueuePO{},
		&po.QueueEntryPO{},
		&po.ExamExecutionPO{},
		// 统计分析
		&po.DashboardSnapshotPO{},
		&po.ReportPO{},
		// 效能优化
		&po.EfficiencyMetricPO{},
		&po.MetricSnapshotPO{},
		&po.BottleneckAlertPO{},
		&po.OptimizationStrategyPO{},
		&po.TrialRunPO{},
		&po.EvaluationReportPO{},
		&po.ROIReportPO{},
		&po.PerformanceScanPO{},
		&po.StrategyDecayAlertPO{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	log.Println("数据库迁移完成")
}
