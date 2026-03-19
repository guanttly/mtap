// Package errors 提供统一的错误码体系和自定义错误类型
package errors

import "fmt"

// Code 错误码定义
type Code int

// 系统级错误码 1xxx
const (
	OK               Code = 0
	ErrUnauthorized  Code = 1001
	ErrForbidden     Code = 1002
	ErrNotFound      Code = 1003
	ErrInvalidParam  Code = 1004
	ErrInternal      Code = 1005
	ErrRateLimit     Code = 1006
	ErrDuplicate     Code = 1007
	ErrConflict      Code = 1008
	ErrTimeout       Code = 1009
)

// 规则引擎 RULE_0xx → 2xxx
const (
	ErrRuleSameItem        Code = 2001
	ErrRuleDuplicate       Code = 2002
	ErrRulePkgNameDup      Code = 2003
	ErrRulePkgTooFew       Code = 2004
	ErrRuleCircularDep     Code = 2005
	ErrRulePresetNoDelete  Code = 2006
	ErrRuleTagNameDup      Code = 2007
	ErrRuleSortingConflict Code = 2008
	ErrRuleInvalidScope    Code = 2009
	ErrRuleServiceDown     Code = 2010
)

// 资源管理 RES_0xx → 3xxx
const (
	ErrResDeviceNotFound    Code = 3001
	ErrResScheduleConflict  Code = 3002
	ErrResSlotLockFail      Code = 3003
	ErrResSlotReleaseFail   Code = 3004
	ErrResSubstituteIncompat Code = 3005
	ErrResExtraSlotOverlap  Code = 3006
	ErrResAliasConflict     Code = 3007
	ErrResAliasTooMany      Code = 3008
	ErrResSyncFailed        Code = 3009
	ErrResSlotOverLimit     Code = 3010
)

// 预约服务 APPT_0xx → 4xxx
const (
	ErrApptNotPaid          Code = 4001
	ErrApptConflictForbid   Code = 4002
	ErrApptDepBlocked       Code = 4003
	ErrApptSlotTaken        Code = 4004
	ErrApptConfirmTimeout   Code = 4005
	ErrApptChangeLimitReached Code = 4006
	ErrApptTooCloseToExam   Code = 4007
	ErrApptBlacklisted      Code = 4008
	ErrApptManualForbidden  Code = 4009
	ErrApptComboTooMany     Code = 4010
	ErrApptPayTimeout       Code = 4011
)

// 分诊管理 TRIAGE_0xx → 5xxx
const (
	ErrTriageNotFound       Code = 5001
	ErrTriageOutOfWindow    Code = 5002
	ErrTriageAlreadyCheckedIn Code = 5003
	ErrTriageQueueEmpty     Code = 5004
	ErrTriageRecallLimit    Code = 5005
	ErrTriageStatusInvalid  Code = 5006
	ErrTriageUndoExpired    Code = 5007
	ErrTriageInvalidQR      Code = 5008
)

// 统计分析 STATS_0xx → 6xxx
const (
	ErrStatsQueryTimeout    Code = 6001
	ErrStatsRangeTooLong    Code = 6002
	ErrStatsReportFailed    Code = 6003
	ErrStatsFileTooLarge    Code = 6004
)

// 效能优化 OPT_0xx → 7xxx
const (
	ErrOptStrategyLimit     Code = 7001
	ErrOptRejectReasonReq   Code = 7002
	ErrOptStatusInvalid     Code = 7003
	ErrOptJointNotComplete  Code = 7004
	ErrOptTrialActive       Code = 7005
	ErrOptCooldown          Code = 7006
	ErrOptEmergencyRollback Code = 7007
	ErrOptEvalFailed        Code = 7008
	ErrOptCTypeNoExec       Code = 7009
	ErrOptCostOverrun       Code = 7010
)

// messages 错误码到中文消息的映射
var messages = map[Code]string{
	OK:              "成功",
	ErrUnauthorized: "未认证，请先登录",
	ErrForbidden:    "权限不足",
	ErrNotFound:     "资源不存在",
	ErrInvalidParam: "参数校验失败",
	ErrInternal:     "服务器内部错误",
	ErrRateLimit:    "请求频率过高，请稍后再试",
	ErrDuplicate:    "重复操作",
	ErrConflict:     "资源冲突",
	ErrTimeout:      "请求超时",

	ErrRuleSameItem:        "冲突规则中项目A与项目B不能相同",
	ErrRuleDuplicate:       "同一项目对的冲突规则已存在",
	ErrRulePkgNameDup:      "冲突包名称已存在",
	ErrRulePkgTooFew:       "冲突包内项目至少需要2个",
	ErrRuleCircularDep:     "存在循环依赖关系",
	ErrRulePresetNoDelete:  "预置标签不可删除",
	ErrRuleTagNameDup:      "优先级标签名称已存在",
	ErrRuleSortingConflict: "同一范围同一时段已存在排序策略",
	ErrRuleInvalidScope:    "生效范围包含无效ID",
	ErrRuleServiceDown:     "规则引擎服务暂时不可用",

	ErrResDeviceNotFound:     "设备不存在或已离线",
	ErrResScheduleConflict:   "排班日期冲突",
	ErrResSlotLockFail:       "号源锁定失败，已被他人锁定",
	ErrResSlotReleaseFail:    "号源释放失败",
	ErrResSubstituteIncompat: "替班目标设备不兼容",
	ErrResExtraSlotOverlap:   "追加号源时段重叠",
	ErrResAliasConflict:      "别名与现有名称冲突",
	ErrResAliasTooMany:       "别名数量超过上限",
	ErrResSyncFailed:         "HIS数据同步失败",
	ErrResSlotOverLimit:      "号源超出设备单日上限",

	ErrApptNotPaid:            "缴费未完成",
	ErrApptConflictForbid:     "存在禁止级冲突",
	ErrApptDepBlocked:         "强制前置依赖未满足",
	ErrApptSlotTaken:          "号源已被抢占",
	ErrApptConfirmTimeout:     "确认超时，号源已释放",
	ErrApptChangeLimitReached: "改约次数已达上限",
	ErrApptTooCloseToExam:     "距检查不足2小时，仅支持取消",
	ErrApptBlacklisted:        "患者处于黑名单限制期",
	ErrApptManualForbidden:    "人工干预需管理员权限",
	ErrApptComboTooMany:       "组合预约项目数量超过上限",
	ErrApptPayTimeout:         "缴费接口超时",

	ErrTriageNotFound:         "预约不存在或非当日预约",
	ErrTriageOutOfWindow:      "不在签到时间窗口内",
	ErrTriageAlreadyCheckedIn: "已签到，不可重复签到",
	ErrTriageQueueEmpty:       "候诊队列为空",
	ErrTriageRecallLimit:      "重叫次数超限",
	ErrTriageStatusInvalid:    "状态流转异常",
	ErrTriageUndoExpired:      "超过撤销时限",
	ErrTriageInvalidQR:        "二维码无效",

	ErrOptStrategyLimit:     "待审核策略数量已达上限",
	ErrOptRejectReasonReq:   "驳回原因为必填项",
	ErrOptStatusInvalid:     "策略当前状态不允许此操作",
	ErrOptJointNotComplete:  "联合审批未全部通过",
	ErrOptTrialActive:       "试运行期间不可生成同类新建议",
	ErrOptCooldown:          "驳回冷却期内",
	ErrOptEmergencyRollback: "紧急回滚已触发",
	ErrOptEvalFailed:        "评估报告生成失败",
	ErrOptCTypeNoExec:       "C类策略仅生成报告",
	ErrOptCostOverrun:       "成本超出预估",
}

// BizError 业务错误
type BizError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *BizError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建业务错误
func New(code Code) *BizError {
	return &BizError{
		Code:    code,
		Message: MessageOf(code),
	}
}

// NewWithDetail 创建带详情的业务错误
func NewWithDetail(code Code, detail string) *BizError {
	return &BizError{
		Code:    code,
		Message: MessageOf(code),
		Detail:  detail,
	}
}

// Wrap 包装标准错误为业务错误
func Wrap(code Code, err error) *BizError {
	return &BizError{
		Code:    code,
		Message: MessageOf(code),
		Detail:  err.Error(),
	}
}

// MessageOf 根据错误码获取消息
func MessageOf(code Code) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Is 判断 error 是否为指定 Code
func Is(err error, code Code) bool {
	if bizErr, ok := err.(*BizError); ok {
		return bizErr.Code == code
	}
	return false
}
