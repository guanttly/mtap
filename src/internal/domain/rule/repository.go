package rule

import "context"

// ConflictRuleRepository 冲突规则仓储接口
type ConflictRuleRepository interface {
	Save(ctx context.Context, rule *ConflictRule) error
	FindByID(ctx context.Context, id string) (*ConflictRule, error)
	FindByItemPair(ctx context.Context, itemAID, itemBID string) (*ConflictRule, error)
	FindAll(ctx context.Context, status RuleStatus) ([]*ConflictRule, error)
	Update(ctx context.Context, rule *ConflictRule) error
	Delete(ctx context.Context, id string) error
}

// ConflictPackageRepository 冲突包仓储接口
type ConflictPackageRepository interface {
	Save(ctx context.Context, pkg *ConflictPackage) error
	FindByID(ctx context.Context, id string) (*ConflictPackage, error)
	FindByName(ctx context.Context, name string) (*ConflictPackage, error)
	FindAll(ctx context.Context) ([]*ConflictPackage, error)
	Update(ctx context.Context, pkg *ConflictPackage) error
	Delete(ctx context.Context, id string) error
}

// DependencyRuleRepository 依赖规则仓储接口
type DependencyRuleRepository interface {
	Save(ctx context.Context, rule *DependencyRule) error
	FindByID(ctx context.Context, id string) (*DependencyRule, error)
	FindByPostItem(ctx context.Context, postItemID string) ([]*DependencyRule, error)
	FindAll(ctx context.Context, status RuleStatus) ([]*DependencyRule, error)
	Update(ctx context.Context, rule *DependencyRule) error
	Delete(ctx context.Context, id string) error
}

// PriorityTagRepository 优先级标签仓储接口
type PriorityTagRepository interface {
	Save(ctx context.Context, tag *PriorityTag) error
	FindByID(ctx context.Context, id string) (*PriorityTag, error)
	FindByName(ctx context.Context, name string) (*PriorityTag, error)
	FindAll(ctx context.Context) ([]*PriorityTag, error)
	Update(ctx context.Context, tag *PriorityTag) error
	Delete(ctx context.Context, id string) error
}

// SortingStrategyRepository 排序策略仓储接口
type SortingStrategyRepository interface {
	Save(ctx context.Context, strategy *SortingStrategy) error
	FindByID(ctx context.Context, id string) (*SortingStrategy, error)
	FindAll(ctx context.Context) ([]*SortingStrategy, error)
	FindByScope(ctx context.Context, scope EffectiveScope) ([]*SortingStrategy, error)
	Update(ctx context.Context, strategy *SortingStrategy) error
	Delete(ctx context.Context, id string) error
}

// PatientAdaptRuleRepository 患者属性适配规则仓储接口
type PatientAdaptRuleRepository interface {
	SaveAll(ctx context.Context, rules []*PatientAdaptRule) error
	FindAll(ctx context.Context) ([]*PatientAdaptRule, error)
	DeleteAll(ctx context.Context) error
}

// SourceControlRepository 开单来源控制仓储接口
type SourceControlRepository interface {
	SaveAll(ctx context.Context, controls []*SourceControl) error
	FindAll(ctx context.Context) ([]*SourceControl, error)
	DeleteAll(ctx context.Context) error
}
