package rule

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
}

// ConflictRuleUpdated 冲突规则变更事件
type ConflictRuleUpdated struct {
	RuleID  string `json:"rule_id"`
	Action  string `json:"action"` // created / updated / deleted
	ItemAID string `json:"item_a_id"`
	ItemBID string `json:"item_b_id"`
}

func (e ConflictRuleUpdated) EventName() string { return "conflict_rule.updated" }

// DependencyRuleUpdated 依赖规则变更事件
type DependencyRuleUpdated struct {
	RuleID     string `json:"rule_id"`
	Action     string `json:"action"`
	PreItemID  string `json:"pre_item_id"`
	PostItemID string `json:"post_item_id"`
}

func (e DependencyRuleUpdated) EventName() string { return "dependency_rule.updated" }
