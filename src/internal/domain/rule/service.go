package rule

import (
	"context"
	"time"
)

// ConflictResult 冲突检测结果
type ConflictResult struct {
	ItemAID        string        `json:"item_a_id"`
	ItemAName      string        `json:"item_a_name"`
	ItemBID        string        `json:"item_b_id"`
	ItemBName      string        `json:"item_b_name"`
	Level          ConflictLevel `json:"level"`
	MinInterval    int           `json:"min_interval"`    // 最小间隔(小时)
	ActualInterval int           `json:"actual_interval"` // 实际间隔(小时), -1=同次预约
	Reason         string        `json:"reason"`
	RuleID         string        `json:"rule_id"`
}

// PatientExamRecord 患者检查历史记录
type PatientExamRecord struct {
	ExamItemID  string    `json:"exam_item_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// ConflictDetectionService 冲突检测领域服务
type ConflictDetectionService struct {
	ruleRepo ConflictRuleRepository
	pkgRepo  ConflictPackageRepository
}

// NewConflictDetectionService 创建冲突检测服务
func NewConflictDetectionService(ruleRepo ConflictRuleRepository, pkgRepo ConflictPackageRepository) *ConflictDetectionService {
	return &ConflictDetectionService{ruleRepo: ruleRepo, pkgRepo: pkgRepo}
}

// DetectConflicts 检测一组项目间的冲突
func (s *ConflictDetectionService) DetectConflicts(ctx context.Context, itemIDs []string, history []PatientExamRecord) ([]ConflictResult, error) {
	if len(itemIDs) < 2 {
		return nil, nil
	}

	// 加载所有启用的冲突规则
	rules, err := s.ruleRepo.FindAll(ctx, RuleStatusActive)
	if err != nil {
		return nil, err
	}

	// 加载所有冲突包并展开
	pkgs, err := s.pkgRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 构建冲突映射表 (itemA-itemB → rule)
	conflictMap := make(map[string]*ConflictRule)
	for _, r := range rules {
		key := pairKey(r.ItemAID, r.ItemBID)
		if existing, ok := conflictMap[key]; !ok || isStricter(r, existing) {
			conflictMap[key] = r
		}
	}

	// 展开冲突包到 rule
	for _, pkg := range pkgs {
		if pkg.Status != RuleStatusActive || !pkg.IsValid() {
			continue
		}
		for i := 0; i < len(pkg.Items); i++ {
			for j := i + 1; j < len(pkg.Items); j++ {
				key := pairKey(pkg.Items[i].ExamItemID, pkg.Items[j].ExamItemID)
				pkgRule := &ConflictRule{
					ID:          pkg.ID,
					ItemAID:     pkg.Items[i].ExamItemID,
					ItemBID:     pkg.Items[j].ExamItemID,
					MinInterval: pkg.MinInterval,
					Level:       pkg.Level,
				}
				if existing, ok := conflictMap[key]; !ok || isStricter(pkgRule, existing) {
					conflictMap[key] = pkgRule
				}
			}
		}
	}

	// 构建患者历史映射
	historyMap := make(map[string]time.Time)
	for _, h := range history {
		historyMap[h.ExamItemID] = h.CompletedAt
	}

	// 两两组合检测
	var results []ConflictResult
	for i := 0; i < len(itemIDs); i++ {
		for j := i + 1; j < len(itemIDs); j++ {
			key := pairKey(itemIDs[i], itemIDs[j])
			rule, exists := conflictMap[key]
			if !exists {
				continue
			}

			actualInterval := -1 // -1 表示同次预约
			// 检查是否有历史记录冲突
			if completedAt, ok := historyMap[itemIDs[i]]; ok {
				hours := int(time.Since(completedAt).Hours())
				if hours < actualInterval || actualInterval == -1 {
					actualInterval = hours
				}
			}
			if completedAt, ok := historyMap[itemIDs[j]]; ok {
				hours := int(time.Since(completedAt).Hours())
				if hours < actualInterval || actualInterval == -1 {
					actualInterval = hours
				}
			}

			// 同次预约的两个项目或间隔不够都触发冲突
			if actualInterval == -1 || actualInterval < rule.MinInterval {
				results = append(results, ConflictResult{
					ItemAID:        itemIDs[i],
					ItemBID:        itemIDs[j],
					Level:          rule.Level,
					MinInterval:    rule.MinInterval,
					ActualInterval: actualInterval,
					Reason:         formatConflictReason(rule.Level, rule.MinInterval, actualInterval),
					RuleID:         rule.ID,
				})
			}
		}
	}

	return results, nil
}

// DependencyStatus 依赖校验状态
type DependencyStatus string

const (
	DependencyPassed  DependencyStatus = "passed"
	DependencyBlocked DependencyStatus = "blocked"
	DependencyExpired DependencyStatus = "expired"
	DependencyUnknown DependencyStatus = "unknown"
)

// DependencyResult 依赖校验结果
type DependencyResult struct {
	PostItemID    string           `json:"post_item_id"`
	PreItemID     string           `json:"pre_item_id"`
	Type          DependencyType   `json:"type"`
	Status        DependencyStatus `json:"status"`
	ValidityHours int              `json:"validity_hours"`
	CompletedAt   *time.Time       `json:"completed_at,omitempty"`
	Reason        string           `json:"reason"`
}

// PreItemChecker 前置检查查询接口（由外部HIS适配器实现）
type PreItemChecker interface {
	GetCompletedTime(ctx context.Context, patientID string, examItemID string) (*time.Time, error)
}

// DependencyValidationService 依赖校验领域服务
type DependencyValidationService struct {
	depRepo  DependencyRuleRepository
	checker  PreItemChecker
}

// NewDependencyValidationService 创建依赖校验服务
func NewDependencyValidationService(repo DependencyRuleRepository, checker PreItemChecker) *DependencyValidationService {
	return &DependencyValidationService{depRepo: repo, checker: checker}
}

// ValidateDependencies 校验一组项目的前置依赖
func (s *DependencyValidationService) ValidateDependencies(ctx context.Context, itemIDs []string, patientID string) ([]DependencyResult, error) {
	var results []DependencyResult

	for _, itemID := range itemIDs {
		deps, err := s.depRepo.FindByPostItem(ctx, itemID)
		if err != nil {
			return nil, err
		}

		for _, dep := range deps {
			if dep.Status != RuleStatusActive {
				continue
			}

			result := DependencyResult{
				PostItemID:    dep.PostItemID,
				PreItemID:     dep.PreItemID,
				Type:          dep.Type,
				ValidityHours: dep.ValidityHours,
			}

			completedAt, err := s.checker.GetCompletedTime(ctx, patientID, dep.PreItemID)
			if err != nil {
				result.Status = DependencyUnknown
				result.Reason = "前置检查查询失败"
				results = append(results, result)
				continue
			}

			if completedAt == nil {
				result.Status = DependencyBlocked
				result.Reason = "前置检查未完成"
				results = append(results, result)
				continue
			}

			result.CompletedAt = completedAt
			hoursSince := time.Since(*completedAt).Hours()
			if hoursSince > float64(dep.ValidityHours) {
				result.Status = DependencyExpired
				result.Reason = "前置检查已超过时效"
			} else {
				result.Status = DependencyPassed
				result.Reason = "前置检查已完成且在时效内"
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// CircularDependencyChecker 循环依赖检测
type CircularDependencyChecker struct {
	depRepo DependencyRuleRepository
}

// NewCircularDependencyChecker 创建循环依赖检测器
func NewCircularDependencyChecker(repo DependencyRuleRepository) *CircularDependencyChecker {
	return &CircularDependencyChecker{depRepo: repo}
}

// HasCircularDependency 检测是否存在循环依赖
func (c *CircularDependencyChecker) HasCircularDependency(ctx context.Context, preItemID, postItemID string) (bool, error) {
	allRules, err := c.depRepo.FindAll(ctx, RuleStatusActive)
	if err != nil {
		return false, err
	}

	// 构建邻接表
	graph := make(map[string][]string)
	for _, rule := range allRules {
		graph[rule.PreItemID] = append(graph[rule.PreItemID], rule.PostItemID)
	}
	// 新增待检测的边
	graph[preItemID] = append(graph[preItemID], postItemID)

	// DFS 检测环
	visited := make(map[string]int) // 0=unvisited, 1=visiting, 2=done
	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = 1
		for _, next := range graph[node] {
			if visited[next] == 1 {
				return true // 发现环
			}
			if visited[next] == 0 {
				if dfs(next) {
					return true
				}
			}
		}
		visited[node] = 2
		return false
	}

	return dfs(preItemID), nil
}

// helpers

func pairKey(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "|" + b
}

func isStricter(a, b *ConflictRule) bool {
	if a.Level == ConflictLevelForbid && b.Level == ConflictLevelWarning {
		return true
	}
	if a.Level == b.Level && a.MinInterval > b.MinInterval {
		return true
	}
	return false
}

func formatConflictReason(level ConflictLevel, minInterval, actualInterval int) string {
	levelText := "警告"
	if level == ConflictLevelForbid {
		levelText = "禁止"
	}
	if actualInterval == -1 {
		return levelText + "：同次预约的项目存在冲突"
	}
	return levelText + "：两项目间最小间隔应为" + intToStr(minInterval) + "小时"
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
