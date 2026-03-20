// Package main 分析工作器
// 订阅优化队列事件，驱动瓶颈归因、策略生成与效果评估
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	domainOpt "github.com/euler/mtap/internal/domain/optimization"
	infraConfig "github.com/euler/mtap/internal/infrastructure/config"
	"github.com/euler/mtap/internal/infrastructure/mq"
	pOpt "github.com/euler/mtap/internal/infrastructure/persistence/optimization"
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
		log.Fatalf("[AnalyzerWorker] 数据库连接失败: %v", err)
	}

	// 初始化仓储与领域服务
	optRepos := pOpt.NewRepositories(db)
	attrSvc := domainOpt.NewBottleneckAttributionService(optRepos.MetricRepo(), optRepos.AlertRepo())
	stratSvc := domainOpt.NewStrategyGenerationService(optRepos.AlertRepo(), optRepos.StrategyRepo(), attrSvc)
	evalSvc := domainOpt.NewEvaluationService(
		optRepos.StrategyRepo(), optRepos.TrialRepo(), optRepos.EvalRepo(),
		optRepos.SnapshotRepo(), optRepos.MetricRepo(),
	)

	// 初始化 Kafka 消费者
	brokers := cfg.Kafka.Brokers
	if len(brokers) == 0 {
		brokers = []string{"localhost:9092"}
	}
	consumer, err := mq.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("[AnalyzerWorker] 初始化 Kafka 消费者失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	groupID := cfg.Kafka.GroupID + ".analyzer"
	// 订阅优化 Topic
	err = consumer.Subscribe(ctx, mq.QueueOptimization, groupID,
		makeBottleneckHandler(ctx, attrSvc, stratSvc))
	if err != nil {
		log.Fatalf("[AnalyzerWorker] 订阅优化 Topic 失败: %v", err)
	}

	// 同一 Topic 再订阅 trial/strategy 事件（通过 event.Type 区分）
	err = consumer.Subscribe(ctx, mq.QueueOptimization, groupID+".eval",
		makeTrialHandler(ctx, evalSvc))
	if err != nil {
		log.Fatalf("[AnalyzerWorker] 订阅 trial 事件失败: %v", err)
	}

	err = consumer.Subscribe(ctx, mq.QueueOptimization, groupID+".promote",
		makeStrategyHandler(ctx, evalSvc))
	if err != nil {
		log.Fatalf("[AnalyzerWorker] 订阅 strategy 事件失败: %v", err)
	}

	log.Println("[AnalyzerWorker] 分析工作器已启动，等待消息...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[AnalyzerWorker] 收到退出信号，正在停止...")
	cancel()
	log.Println("[AnalyzerWorker] 已停止")
}

// makeBottleneckHandler 处理瓶颈检测事件：进行归因分析并生成优化策略
func makeBottleneckHandler(
	ctx context.Context,
	attrSvc *domainOpt.BottleneckAttributionService,
	stratSvc *domainOpt.StrategyGenerationService,
) mq.Handler {
	return func(_ context.Context, event mq.Event) error {
		alertID, _ := extractString(event.Payload, "alert_id")
		if alertID == "" {
			log.Printf("[AnalyzerWorker] bottleneck 事件缺少 alert_id，跳过")
			return nil
		}
		log.Printf("[AnalyzerWorker] 开始归因分析: alertID=%s", alertID)
		report, err := attrSvc.Analyze(ctx, alertID)
		if err != nil {
			log.Printf("[AnalyzerWorker] 归因分析失败: %v", err)
			return err
		}
		log.Printf("[AnalyzerWorker] 归因完成，建议分类: %s，假设数: %d", report.SuggestedCategory, len(report.Hypotheses))

		strategy, err := stratSvc.GenerateFromAlert(ctx, alertID)
		if err != nil {
			log.Printf("[AnalyzerWorker] 策略生成失败 (alertID=%s): %v", alertID, err)
			// 策略生成失败不 Nack（可能是待审核策略已满），返回 nil 避免无限重试
			return nil
		}
		log.Printf("[AnalyzerWorker] 策略已生成: strategyID=%s title=%s", strategy.ID, strategy.Title)
		return nil
	}
}

// makeTrialHandler 处理试验完成事件：生成评估报告
func makeTrialHandler(ctx context.Context, evalSvc *domainOpt.EvaluationService) mq.Handler {
	return func(_ context.Context, event mq.Event) error {
		if event.Type != mq.RoutingKeyTrialCompleted {
			return nil
		}
		strategyID, _ := extractString(event.Payload, "strategy_id")
		if strategyID == "" {
			log.Printf("[AnalyzerWorker] trial.completed 事件缺少 strategy_id，跳过")
			return nil
		}
		log.Printf("[AnalyzerWorker] 生成评估报告: strategyID=%s", strategyID)
		report, err := evalSvc.GenerateReport(ctx, strategyID)
		if err != nil {
			log.Printf("[AnalyzerWorker] 评估报告生成失败: %v", err)
			return err
		}
		log.Printf("[AnalyzerWorker] 评估报告已生成: reportID=%s", report.ID)
		return nil
	}
}

// makeStrategyHandler 处理策略审批事件：晋升为生产策略
func makeStrategyHandler(ctx context.Context, evalSvc *domainOpt.EvaluationService) mq.Handler {
	return func(_ context.Context, event mq.Event) error {
		if event.Type != mq.RoutingKeyStrategyPromoted {
			return nil
		}
		strategyID, _ := extractString(event.Payload, "strategy_id")
		if strategyID == "" {
			log.Printf("[AnalyzerWorker] strategy.promoted 事件缺少 strategy_id，跳过")
			return nil
		}
		log.Printf("[AnalyzerWorker] 晋升策略: strategyID=%s", strategyID)
		if err := evalSvc.PromoteStrategy(ctx, strategyID); err != nil {
			log.Printf("[AnalyzerWorker] 策略晋升失败: %v", err)
			return err
		}
		log.Printf("[AnalyzerWorker] 策略已晋升: strategyID=%s", strategyID)
		return nil
	}
}

// extractString 从 map[string]any 中安全提取字符串值
func extractString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	switch s := v.(type) {
	case string:
		return s, true
	default:
		b, _ := json.Marshal(v)
		return string(b), true
	}
}
