// Package resource 资源管理领域 - 值对象
// 核心目的：定义资源管理领域的不可变值对象
// 模块功能：
//   - AgeFactor:    年龄折算系数（儿童+10%/老年+15%）
//   - SlotPattern:  号源排程模式（项目ID+连续数量）
//   - WorkPeriod:   工作时段（起止时间字符串）
//
// 注：DeviceStatus / ScheduleStatus / TimeSlotStatus 枚举定义在 entity.go
package resource

// AgeFactor 年龄耗时折算系数（值对象）
type AgeFactor struct {
	ChildMaxAge   int     // 儿童年龄上限，默认 14
	ChildFactor   float64 // 儿童耗时折算系数，默认 1.10（+10%）
	ElderlyMinAge int     // 老年年龄下限，默认 70
	ElderlyFactor float64 // 老年耗时折算系数，默认 1.15（+15%）
}

// DefaultAgeFactor 系统默认年龄折算系数
var DefaultAgeFactor = AgeFactor{
	ChildMaxAge:   14,
	ChildFactor:   1.10,
	ElderlyMinAge: 70,
	ElderlyFactor: 1.15,
}

// Apply 按患者年龄计算实际耗时（分钟）
func (a AgeFactor) Apply(baseDurationMin, patientAge int) int {
	d := float64(baseDurationMin)
	switch {
	case patientAge > 0 && patientAge <= a.ChildMaxAge:
		d *= a.ChildFactor
	case patientAge >= a.ElderlyMinAge:
		d *= a.ElderlyFactor
	}
	return int(d)
}

// SlotPattern 号源排程模式（在排班模板中循环使用）
type SlotPattern struct {
	ExamItemID string // 绑定的检查项目ID
	Count      int    // 该模式连续生成的号源数量
}

// WorkPeriod 工作时段（值对象，HH:mm 格式）
type WorkPeriod struct {
	Start string // 如 "08:00"
	End   string // 如 "17:00"
}
