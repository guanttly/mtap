// Package triage 分诊执行领域 - 值对象
package triage

// CheckInMethod 签到方式
type CheckInMethod string

const (
	CheckInKiosk CheckInMethod = "kiosk" // 自助机扫码
	CheckInNurse CheckInMethod = "nurse" // 护士站手动
	CheckInNFC   CheckInMethod = "nfc"   // NFC读卡
)

func (m CheckInMethod) IsValid() bool {
	switch m {
	case CheckInKiosk, CheckInNurse, CheckInNFC:
		return true
	}
	return false
}

// EntryStatus 队列条目状态
type EntryStatus string

const (
	EntryWaiting   EntryStatus = "waiting"   // 候诊中
	EntryCalling   EntryStatus = "calling"   // 呼叫中
	EntryExamining EntryStatus = "examining" // 检查中
	EntryCompleted EntryStatus = "completed" // 已完成
	EntryMissed    EntryStatus = "missed"    // 已过号
	EntryNoShow    EntryStatus = "no_show"   // 爽约
)

// ExamStatus 检查状态
type ExamStatus string

const (
	ExamCheckedIn ExamStatus = "checked_in" // 已签到
	ExamWaiting   ExamStatus = "waiting"    // 候诊中
	ExamOngoing   ExamStatus = "ongoing"    // 检查中
	ExamDone      ExamStatus = "done"       // 检查完成
)

// MaxCallCount 最大重叫次数
const MaxCallCount = 3

// MaxMissCount 最大过号次数（达到则标记爽约）
const MaxMissCount = 2

// UndoWindowMinutes 误操作撤销窗口（分钟）
const UndoWindowMinutes = 5
