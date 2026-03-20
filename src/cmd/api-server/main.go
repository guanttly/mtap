package main

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	appAnalytics "github.com/euler/mtap/internal/application/analytics"
	appAppt "github.com/euler/mtap/internal/application/appointment"
	appOpt "github.com/euler/mtap/internal/application/optimization"
	appRes "github.com/euler/mtap/internal/application/resource"
	appRule "github.com/euler/mtap/internal/application/rule"
	appTriage "github.com/euler/mtap/internal/application/triage"
	domainRule "github.com/euler/mtap/internal/domain/rule"
	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	pAnalytics "github.com/euler/mtap/internal/infrastructure/persistence/analytics"
	pAppt "github.com/euler/mtap/internal/infrastructure/persistence/appointment"
	pOpt "github.com/euler/mtap/internal/infrastructure/persistence/optimization"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	pRes "github.com/euler/mtap/internal/infrastructure/persistence/resource"
	pRule "github.com/euler/mtap/internal/infrastructure/persistence/rule"
	pTriage "github.com/euler/mtap/internal/infrastructure/persistence/triage"
	httpx "github.com/euler/mtap/internal/interfaces/http"
	httpAdmin "github.com/euler/mtap/internal/interfaces/http/admin"
	httpAnalytics "github.com/euler/mtap/internal/interfaces/http/analytics"
	httpAppt "github.com/euler/mtap/internal/interfaces/http/appointment"
	httpAuth "github.com/euler/mtap/internal/interfaces/http/auth"
	httpOpt "github.com/euler/mtap/internal/interfaces/http/optimization"
	httpRes "github.com/euler/mtap/internal/interfaces/http/resource"
	httpRule "github.com/euler/mtap/internal/interfaces/http/rule"
	httpTriage "github.com/euler/mtap/internal/interfaces/http/triage"
	"github.com/euler/mtap/pkg/auth"
	"github.com/euler/mtap/pkg/logger"
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
		dbPath := cfg.Database.DSNString()
		if v := os.Getenv("MTAP_DB"); v != "" {
			dbPath = v
		}
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	// AutoMigrate (开发期最小可运行)
	if err := db.AutoMigrate(
		&po.RolePO{},
		&po.UserPO{},
		&po.ConflictRulePO{},
		&po.ConflictPackagePO{},
		&po.ConflictPackageItemPO{},
		&po.DependencyRulePO{},
		&po.PriorityTagPO{},
		&po.SortingStrategyPO{},
		&po.PatientAdaptRulePO{},
		&po.SourceControlPO{},

		&po.DevicePO{},
		&po.ExamItemPO{},
		&po.ItemAliasPO{},
		&po.SlotPoolPO{},
		&po.SchedulePO{},
		&po.TimeSlotPO{},

		&po.AppointmentPO{},
		&po.AppointmentItemPO{},
		&po.AppointmentCredentialPO{},
		&po.AppointmentChangeLogPO{},
		&po.BlacklistPO{},
		&po.NoShowRecordPO{},
		&po.AppealPO{},

		&po.CheckInPO{},
		&po.WaitingQueuePO{},
		&po.QueueEntryPO{},
		&po.ExamExecutionPO{},

		&po.DashboardSnapshotPO{},
		&po.ReportPO{},

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

	jwtSecret := cfg.JWT.Secret
	if v := os.Getenv("MTAP_JWT_SECRET"); v != "" {
		jwtSecret = v
	}
	jwtMgr := auth.NewJWTManager(jwtSecret, cfg.JWT.AccessExpire, cfg.JWT.RefreshExpire)

	// audit logger
	auditPath := cfg.Log.AuditPath
	if v := os.Getenv("MTAP_AUDIT_LOG"); v != "" {
		auditPath = v
	}
	auditLogger, err := logger.NewFileAuditLogger(auditPath)
	if err != nil {
		log.Fatalf("init audit logger failed: %v", err)
	}
	defer auditLogger.Close()
	logger.SetAuditLogger(auditLogger)

	// rule module DI
	ruleRepos := pRule.NewRepositories(db)
	conflictRuleRepo := ruleRepos.ConflictRuleRepo()
	conflictPkgRepo := ruleRepos.ConflictPackageRepo()
	depRuleRepo := ruleRepos.DependencyRuleRepo()
	tagRepo := ruleRepos.PriorityTagRepo()
	sortRepo := ruleRepos.SortingStrategyRepo()
	adaptRepo := ruleRepos.PatientAdaptRuleRepo()
	sourceRepo := ruleRepos.SourceControlRepo()

	conflictSvc := domainRule.NewConflictDetectionService(conflictRuleRepo, conflictPkgRepo)
	depSvc := domainRule.NewDependencyValidationService(depRuleRepo, &noopPreItemChecker{})
	circular := domainRule.NewCircularDependencyChecker(depRuleRepo)

	ruleApp := appRule.NewRuleAppService(
		conflictRuleRepo, conflictPkgRepo, depRuleRepo, tagRepo,
		sortRepo, adaptRepo, sourceRepo,
		conflictSvc, depSvc, circular,
	)
	ruleHandler := httpRule.NewHandler(ruleApp)

	// resource module DI
	resRepos := pRes.NewRepositories(db)
	resSvc := appRes.NewService(
		resRepos.CampusRepo(),
		resRepos.DepartmentRepo(),
		resRepos.DeviceRepo(),
		resRepos.ExamItemRepo(),
		resRepos.AliasRepo(),
		resRepos.SlotPoolRepo(),
		resRepos.ScheduleRepo(),
		resRepos.TimeSlotRepo(),
	)
	resHandler := httpRes.NewHandler(resSvc)

	// connect rule-check fasting -> resource exam items
	ruleApp.WithExamItemMetaProvider(resRepos.ExamItemRepo())

	// appointment module DI
	apptRepos := pAppt.NewRepositories(db)
	apptSvc := appAppt.NewAppointmentAppService(
		apptRepos.AppointmentRepo(),
		apptRepos.AppointmentItemRepo(),
		apptRepos.CredentialRepo(),
		apptRepos.BlacklistRepo(),
		apptRepos.NoShowRepo(),
		apptRepos.AppealRepo(),
	)
	apptHandler := httpAppt.NewHandler(apptSvc)

	// triage module DI
	triageRepos := pTriage.NewRepositories(db)
	triageSvc := appTriage.NewTriageAppService(
		triageRepos.CheckInRepo(),
		triageRepos.WaitingQueueRepo(),
		triageRepos.QueueEntryRepo(),
		triageRepos.ExamExecutionRepo(),
	)
	triageHandler := httpTriage.NewHandler(triageSvc)

	// analytics module DI
	analyticsRepos := pAnalytics.NewRepositories(db)
	analyticsSvc := appAnalytics.NewAnalyticsAppService(
		analyticsRepos.DashboardRepo(),
		analyticsRepos.ReportRepo(),
	)
	analyticsHandler := httpAnalytics.NewHandler(analyticsSvc)

	// optimization module DI
	optRepos := pOpt.NewRepositories(db)
	optSvc := appOpt.NewOptimizationAppService(
		optRepos.MetricRepo(),
		optRepos.SnapshotRepo(),
		optRepos.AlertRepo(),
		optRepos.StrategyRepo(),
		optRepos.TrialRepo(),
		optRepos.EvalRepo(),
		optRepos.ROIRepo(),
		optRepos.ScanRepo(),
		optRepos.DecayRepo(),
	)
	optHandler := httpOpt.NewHandler(optSvc)

	authHandler := httpAuth.NewHandler(db, jwtMgr)
	adminHandler := httpAdmin.NewHandler(db)
	// seed 预置角色和默认 admin 账号（幂等）
	if err := httpAuth.SeedAdminUser(db); err != nil {
		log.Printf("seed admin user failed: %v", err)
	}

	engine := httpx.NewEngine(httpx.Deps{
		JWTManager:          jwtMgr,
		AuthHandler:         authHandler,
		AdminHandler:        adminHandler,
		RuleHandler:         ruleHandler,
		ResourceHandler:     resHandler,
		AppointmentHandler:  apptHandler,
		TriageHandler:       triageHandler,
		AnalyticsHandler:    analyticsHandler,
		OptimizationHandler: optHandler,
	})

	port := os.Getenv("MTAP_PORT")
	if port == "" {
		port = "8080"
	}
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}

type noopPreItemChecker struct{}

func (n *noopPreItemChecker) GetCompletedTime(_ context.Context, _, _ string) (*time.Time, error) {
	return nil, nil
}
