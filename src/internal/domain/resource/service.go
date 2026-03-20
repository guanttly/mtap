// Package resource 资源管理领域 - 领域服务
// 核心目的：实现号源生成与排班管理的核心业务逻辑
// 模块功能：
//   - SlotGenerationService: 动态号源生成（按设备切分时段、年龄折算）
//   - ScheduleService: 排班管理（批量生成、临时停诊、替班、追加号源）
//   - SlotPoolService: 号源池管理（分配、锁定、释放、溢出）
//   - ItemMappingService: 项目别名映射解析
package resource

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ─── SlotGenerationService ──────────────────────────────────────────────────

// SlotGenerationInput 号源生成输入
type SlotGenerationInput struct {
	DeviceID      string    // 设备ID
	ScheduleID    string    // 排班ID
	WorkDate      time.Time // 工作日期
	StartTime     string    // HH:mm 开始时间
	EndTime       string    // HH:mm 结束时间
	ExamItemID    string    // 绑定检查项目ID
	SlotMinutes   int       // 单个号源时长（分钟）
	PoolType      string    // 号源池类型（public/department/doctor）
	MaxDailySlots int       // 单日最大号源数，0=不限
}

// SlotGenerationService 号源生成领域服务（纯域逻辑，不涉及持久化）
type SlotGenerationService struct {
	ageFactor AgeFactor
}

// NewSlotGenerationService 创建号源生成服务
func NewSlotGenerationService() *SlotGenerationService {
	return &SlotGenerationService{ageFactor: DefaultAgeFactor}
}

// Generate 按排班参数生成号源时段列表
func (s *SlotGenerationService) Generate(input SlotGenerationInput) ([]*TimeSlot, error) {
	if input.SlotMinutes <= 0 {
		return nil, fmt.Errorf("slot_minutes must be positive")
	}

	startH, startM, err := parseHHMM(input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time %q: %w", input.StartTime, err)
	}
	endH, endM, err := parseHHMM(input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time %q: %w", input.EndTime, err)
	}

	d := input.WorkDate
	startAt := time.Date(d.Year(), d.Month(), d.Day(), startH, startM, 0, 0, d.Location())
	endAt := time.Date(d.Year(), d.Month(), d.Day(), endH, endM, 0, 0, d.Location())
	if !endAt.After(startAt) {
		return nil, fmt.Errorf("end_time must be after start_time")
	}

	poolType := input.PoolType
	if poolType == "" {
		poolType = "public"
	}

	slotDur := time.Duration(input.SlotMinutes) * time.Minute
	var slots []*TimeSlot
	for cur := startAt; !cur.Add(slotDur).After(endAt); cur = cur.Add(slotDur) {
		if input.MaxDailySlots > 0 && len(slots) >= input.MaxDailySlots {
			break
		}
		slots = append(slots, NewTimeSlot(
			input.ScheduleID,
			input.DeviceID,
			input.ExamItemID,
			poolType,
			d,
			cur,
			cur.Add(slotDur),
			input.SlotMinutes,
		))
	}
	return slots, nil
}

// AdjustDurationForAge 按患者年龄折算实际检查耗时（分钟）
func (s *SlotGenerationService) AdjustDurationForAge(baseDurationMin, patientAge int) int {
	return s.ageFactor.Apply(baseDurationMin, patientAge)
}

// ─── ScheduleService ────────────────────────────────────────────────────────

// ScheduleConflictResult 排班重叠检测结果
type ScheduleConflictResult struct {
	HasConflict bool
	Reason      string
}

// ScheduleService 排班管理领域服务（纯域逻辑）
type ScheduleService struct{}

// NewScheduleService 创建排班管理服务
func NewScheduleService() *ScheduleService { return &ScheduleService{} }

// ValidateNoOverlap 校验新排班与已有排班不存在时间重叠
func (s *ScheduleService) ValidateNoOverlap(existing []Schedule, newStart, newEnd string) ScheduleConflictResult {
	for _, sch := range existing {
		if timesOverlap(sch.StartTime, sch.EndTime, newStart, newEnd) {
			return ScheduleConflictResult{
				HasConflict: true,
				Reason:      fmt.Sprintf("与已有排班 %s—%s 时间重叠", sch.StartTime, sch.EndTime),
			}
		}
	}
	return ScheduleConflictResult{}
}

// GenerateDates 在 [startDate, endDate] 范围内生成工作日期列表
func (s *ScheduleService) GenerateDates(startDate, endDate time.Time, skipWeekends bool) []time.Time {
	var dates []time.Time
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if skipWeekends {
			wd := d.Weekday()
			if wd == time.Saturday || wd == time.Sunday {
				continue
			}
		}
		dates = append(dates, d)
	}
	return dates
}

// BatchGenerateInput 批量排班生成输入
type BatchGenerateInput struct {
	DeviceID     string
	StartDate    time.Time
	EndDate      time.Time
	StartTime    string // HH:mm
	EndTime      string // HH:mm
	SlotMinutes  int
	ExamItemID   string
	PoolType     string
	SkipWeekends bool
	MaxPerDay    int
}

// BatchGenerateResult 批量生成结果摘要
type BatchGenerateResult struct {
	DateCount    int
	SlotCount    int
	SkippedDates []time.Time // 因冲突跳过的日期
}

// PlanBatchGenerate 规划批量排班（仅做计算，不持久化）
func (s *ScheduleService) PlanBatchGenerate(input BatchGenerateInput) BatchGenerateResult {
	dates := s.GenerateDates(input.StartDate, input.EndDate, input.SkipWeekends)
	svc := NewSlotGenerationService()
	totalSlots := 0
	for _, date := range dates {
		slots, err := svc.Generate(SlotGenerationInput{
			DeviceID:      input.DeviceID,
			WorkDate:      date,
			StartTime:     input.StartTime,
			EndTime:       input.EndTime,
			ExamItemID:    input.ExamItemID,
			SlotMinutes:   input.SlotMinutes,
			PoolType:      input.PoolType,
			MaxDailySlots: input.MaxPerDay,
		})
		if err == nil {
			totalSlots += len(slots)
		}
	}
	return BatchGenerateResult{
		DateCount: len(dates),
		SlotCount: totalSlots,
	}
}

// ─── SlotPoolService ─────────────────────────────────────────────────────────

// SlotPoolService 号源池领域服务（纯域逻辑）
type SlotPoolService struct{}

// NewSlotPoolService 创建号源池服务
func NewSlotPoolService() *SlotPoolService { return &SlotPoolService{} }

// ValidatePoolQuota 校验各池配额之和不超过100%
// 规则：各号源池 AllocationRatio 之和 ≤ 1.0
func (s *SlotPoolService) ValidatePoolQuota(pools []*SlotPool) error {
	var total float64
	for _, p := range pools {
		if p.Status == "active" {
			total += p.AllocationRatio
		}
	}
	if total > 1.0001 { // 允许浮点精度误差
		return fmt.Errorf("号源池配额之和(%.1f%%)超过100%%，请调整各池分配比例", total*100)
	}
	return nil
}

// IsSlotAvailable 判断号源是否可用
func (s *SlotPoolService) IsSlotAvailable(slot *TimeSlot) bool {
	if slot.Remaining <= 0 {
		return false
	}
	if slot.Status == TimeSlotAvailable {
		return true
	}
	// 锁定超时自动视为可用
	if slot.IsExpiredLock() {
		return true
	}
	return false
}

// FilterByPoolType 按号源池类型过滤号源
func (s *SlotPoolService) FilterByPoolType(slots []*TimeSlot, poolType string) []*TimeSlot {
	if poolType == "" {
		return slots
	}
	result := make([]*TimeSlot, 0)
	for _, slot := range slots {
		if slot.PoolType == poolType {
			result = append(result, slot)
		}
	}
	return result
}

// SelectBestSlots 从可用号源中选出最优分配（最早时间优先）
func (s *SlotPoolService) SelectBestSlots(slots []*TimeSlot, count int) []*TimeSlot {
	var available []*TimeSlot
	for _, slot := range slots {
		if s.IsSlotAvailable(slot) {
			available = append(available, slot)
		}
	}
	if count <= 0 || count >= len(available) {
		return available
	}
	return available[:count]
}

// ─── ItemMappingService ──────────────────────────────────────────────────────

// ItemMappingService 检查项目别名映射领域服务
type ItemMappingService struct{}

// NewItemMappingService 创建项目别名映射服务
func NewItemMappingService() *ItemMappingService { return &ItemMappingService{} }

// ResolveAlias 将输入名称（标准名或别名）解析为标准检查项目
func (s *ItemMappingService) ResolveAlias(input string, items []ExamItem) *ExamItem {
	input = strings.TrimSpace(input)
	for i := range items {
		if items[i].MatchName(input) {
			return &items[i]
		}
	}
	return nil
}

// BatchResolve 批量解析，分别返回匹配成功与未匹配的项目名称
func (s *ItemMappingService) BatchResolve(inputs []string, items []ExamItem) (found []*ExamItem, notFound []string) {
	for _, inp := range inputs {
		item := s.ResolveAlias(inp, items)
		if item != nil {
			found = append(found, item)
		} else {
			notFound = append(notFound, inp)
		}
	}
	return
}

// ─── 内部辅助 ────────────────────────────────────────────────────────────────

// parseHHMM 解析 "HH:mm" 格式返回 hour, minute
func parseHHMM(s string) (int, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected HH:mm")
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return 0, 0, fmt.Errorf("invalid hour")
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("invalid minute")
	}
	return h, m, nil
}

// timesOverlap 判断两个 HH:mm 时段是否重叠
func timesOverlap(startA, endA, startB, endB string) bool {
	toMin := func(t string) int {
		parts := strings.SplitN(t, ":", 2)
		if len(parts) != 2 {
			return 0
		}
		h, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		return h*60 + m
	}
	aStart, aEnd := toMin(startA), toMin(endA)
	bStart, bEnd := toMin(startB), toMin(endB)
	return aStart < bEnd && bStart < aEnd
}
