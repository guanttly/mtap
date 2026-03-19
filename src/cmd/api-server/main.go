package main

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	appRule "github.com/euler/mtap/internal/application/rule"
	appRes "github.com/euler/mtap/internal/application/resource"
	domainRule "github.com/euler/mtap/internal/domain/rule"
	pRule "github.com/euler/mtap/internal/infrastructure/persistence/rule"
	pRes "github.com/euler/mtap/internal/infrastructure/persistence/resource"
	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	httpx "github.com/euler/mtap/internal/interfaces/http"
	httpRes "github.com/euler/mtap/internal/interfaces/http/resource"
	httpRule "github.com/euler/mtap/internal/interfaces/http/rule"
	"github.com/euler/mtap/pkg/auth"
	"github.com/euler/mtap/pkg/logger"
)

func main() {
	dbPath := os.Getenv("MTAP_DB")
	if dbPath == "" {
		dbPath = "mtap.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	// AutoMigrate (开发期最小可运行)
	if err := db.AutoMigrate(
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
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	jwtSecret := os.Getenv("MTAP_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me-32-bytes-len!"
	}
	jwtMgr := auth.NewJWTManager(jwtSecret, 2*time.Hour, 7*24*time.Hour)

	// audit logger (开发期默认落文件)
	auditPath := os.Getenv("MTAP_AUDIT_LOG")
	if auditPath == "" {
		auditPath = "audit.log"
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

	engine := httpx.NewEngine(httpx.Deps{
		JWTManager:  jwtMgr,
		RuleHandler: ruleHandler,
		ResourceHandler: resHandler,
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
