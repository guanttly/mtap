package resource

import (
	"context"
	"time"

	"github.com/google/uuid"

	bizErr "github.com/euler/mtap/pkg/errors"
)

type (
	DeviceRepository interface {
		Create(ctx context.Context, d DeviceResp) error
		Get(ctx context.Context, id string) (*DeviceResp, error)
		List(ctx context.Context) ([]DeviceResp, error)
	}
	ExamItemRepository interface {
		Create(ctx context.Context, e ExamItemResp) error
		Get(ctx context.Context, id string) (*ExamItemResp, error)
		List(ctx context.Context) ([]ExamItemResp, error)
		ListFastingIDs(ctx context.Context, ids []string) ([]string, error)
		GetDurationMin(ctx context.Context, id string) (int, error)
	}
	AliasRepository interface {
		Create(ctx context.Context, a AliasResp) error
		List(ctx context.Context, examItemID string) ([]AliasResp, error)
	}
	SlotPoolRepository interface {
		Create(ctx context.Context, p SlotPoolResp) error
		List(ctx context.Context) ([]SlotPoolResp, error)
	}
	ScheduleRepository interface {
		Create(ctx context.Context, deviceID string, date time.Time, startTime, endTime string) (string, error)
		Suspend(ctx context.Context, deviceID string, date time.Time, reason string) error
		Substitute(ctx context.Context, sourceDeviceID, targetDeviceID string, date time.Time) error
	}
	TimeSlotRepository interface {
		BulkCreate(ctx context.Context, slots []TimeSlotResp) error
		ListByDeviceAndDate(ctx context.Context, deviceID string, date time.Time) ([]TimeSlotResp, error)
		QueryAvailable(ctx context.Context, deviceID string, date time.Time, examItemID, poolType string) ([]TimeSlotResp, error)
		Lock(ctx context.Context, slotID string, patientID string, lockUntil time.Time) error
		Release(ctx context.Context, slotID string, patientID string, allowForce bool) error
		SuspendRange(ctx context.Context, deviceID string, date time.Time, startAt, endAt time.Time, reason string) (int64, error)
		UpdateDeviceByDate(ctx context.Context, sourceDeviceID, targetDeviceID string, date time.Time) (int64, error)
		HasOverlap(ctx context.Context, deviceID string, date time.Time, startAt, endAt time.Time) (bool, error)
	}
)

type Service struct {
	deviceRepo   DeviceRepository
	examRepo     ExamItemRepository
	aliasRepo    AliasRepository
	slotPoolRepo SlotPoolRepository
	scheduleRepo ScheduleRepository
	timeSlotRepo TimeSlotRepository
}

func NewService(
	deviceRepo DeviceRepository,
	examRepo ExamItemRepository,
	aliasRepo AliasRepository,
	slotPoolRepo SlotPoolRepository,
	scheduleRepo ScheduleRepository,
	timeSlotRepo TimeSlotRepository,
) *Service {
	return &Service{
		deviceRepo:   deviceRepo,
		examRepo:     examRepo,
		aliasRepo:    aliasRepo,
		slotPoolRepo: slotPoolRepo,
		scheduleRepo: scheduleRepo,
		timeSlotRepo: timeSlotRepo,
	}
}

func (s *Service) CreateDevice(ctx context.Context, req CreateDeviceReq) (*DeviceResp, error) {
	now := time.Now()
	d := DeviceResp{
		ID:           uuid.New().String(),
		Name:         req.Name,
		CampusID:     req.CampusID,
		DepartmentID: req.DepartmentID,
		Status:       "active",
		CreatedAt:    now,
	}
	if err := s.deviceRepo.Create(ctx, d); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &d, nil
}

func (s *Service) ListDevices(ctx context.Context) ([]DeviceResp, error) {
	list, err := s.deviceRepo.List(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return list, nil
}

func (s *Service) CreateExamItem(ctx context.Context, req CreateExamItemReq) (*ExamItemResp, error) {
	e := ExamItemResp{
		ID:          uuid.New().String(),
		Name:        req.Name,
		DurationMin: req.DurationMin,
		IsFasting:   req.IsFasting,
		FastingDesc: req.FastingDesc,
	}
	if err := s.examRepo.Create(ctx, e); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &e, nil
}

func (s *Service) ListExamItems(ctx context.Context) ([]ExamItemResp, error) {
	list, err := s.examRepo.List(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return list, nil
}

func (s *Service) CreateAlias(ctx context.Context, req CreateAliasReq) (*AliasResp, error) {
	a := AliasResp{
		ID:         uuid.New().String(),
		ExamItemID: req.ExamItemID,
		Alias:      req.Alias,
	}
	if err := s.aliasRepo.Create(ctx, a); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &a, nil
}

func (s *Service) ListAliases(ctx context.Context, examItemID string) ([]AliasResp, error) {
	list, err := s.aliasRepo.List(ctx, examItemID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return list, nil
}

func (s *Service) CreateSlotPool(ctx context.Context, req CreateSlotPoolReq) (*SlotPoolResp, error) {
	p := SlotPoolResp{
		ID:     uuid.New().String(),
		Name:   req.Name,
		Type:   req.Type,
		Status: "active",
	}
	if err := s.slotPoolRepo.Create(ctx, p); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return &p, nil
}

func (s *Service) ListSlotPools(ctx context.Context) ([]SlotPoolResp, error) {
	list, err := s.slotPoolRepo.List(ctx)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return list, nil
}

func (s *Service) GenerateSchedule(ctx context.Context, req GenerateScheduleReq) ([]TimeSlotResp, error) {
	// 支持单日(date) 或批量(start_date~end_date)
	var dates []time.Time
	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
		}
		dates = []time.Time{d}
	} else if req.StartDate != "" && req.EndDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "start_date格式应为YYYY-MM-DD")
		}
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "end_date格式应为YYYY-MM-DD")
		}
		if endDate.Before(startDate) {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "end_date必须不早于start_date")
		}
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			if req.SkipWeekends {
				wd := d.Weekday()
				if wd == time.Saturday || wd == time.Sunday {
					continue
				}
			}
			dates = append(dates, d)
		}
	} else {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "必须提供date，或同时提供start_date与end_date")
	}

	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "start_time格式应为HH:mm")
	}
	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "end_time格式应为HH:mm")
	}

	slotDur := time.Duration(req.SlotMinutes) * time.Minute
	stdMin := int(slotDur.Minutes())
	poolType := req.PoolType
	if poolType == "" {
		poolType = "public"
	}

	var allSlots []TimeSlotResp
	for _, date := range dates {
		startAt := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
		endAt := time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
		if !endAt.After(startAt) {
			return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "结束时间必须晚于开始时间")
		}

		_, err = s.scheduleRepo.Create(ctx, req.DeviceID, date, req.StartTime, req.EndTime)
		if err != nil {
			// 若同设备同日已存在排班，认为是冲突
			if bizErr.Is(err, bizErr.ErrResScheduleConflict) {
				return nil, err
			}
			return nil, bizErr.Wrap(bizErr.ErrInternal, err)
		}

		var slots []TimeSlotResp
		for cur := startAt; !cur.Add(slotDur).After(endAt); cur = cur.Add(slotDur) {
			slots = append(slots, TimeSlotResp{
				ID:              uuid.New().String(),
				DeviceID:         req.DeviceID,
				ExamItemID:       req.ExamItemID,
				PoolType:         poolType,
				StartAt:          cur,
				EndAt:            cur.Add(slotDur),
				Status:           "available",
				StandardDuration: stdMin,
				AdjustedDuration: stdMin,
				Remaining:        1,
			})
		}
		if err := s.timeSlotRepo.BulkCreate(ctx, slots); err != nil {
			return nil, bizErr.Wrap(bizErr.ErrInternal, err)
		}
		allSlots = append(allSlots, slots...)
	}

	return allSlots, nil
}

func (s *Service) ListSlots(ctx context.Context, deviceID, dateStr string) ([]TimeSlotResp, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
	}
	list, err := s.timeSlotRepo.ListByDeviceAndDate(ctx, deviceID, date)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return list, nil
}

func (s *Service) QueryAvailableSlots(ctx context.Context, deviceID, dateStr, examItemID, poolType string, patientAge int) ([]TimeSlotResp, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
	}
	list, err := s.timeSlotRepo.QueryAvailable(ctx, deviceID, date, examItemID, poolType)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	// 年龄折算：在资源模块内计算 adjusted_duration（不改动slot本身，只影响返回）
	// 规则：儿童(<14)+10%，老年(>70)+15%
	// 如果 slot.StandardDuration 为 0 且绑定了 exam_item_id，则从 exam_items 补齐
	for i := range list {
		std := list[i].StandardDuration
		if std == 0 && list[i].ExamItemID != "" && s.examRepo != nil {
			if d, err := s.examRepo.GetDurationMin(ctx, list[i].ExamItemID); err == nil && d > 0 {
				std = d
				list[i].StandardDuration = d
			}
		}
		factor := 1.0
		if patientAge > 0 {
			if patientAge < 14 {
				factor = 1.10
			} else if patientAge > 70 {
				factor = 1.15
			}
		}
		list[i].AdjustedDuration = int(float64(std)*factor + 0.5)
	}
	return list, nil
}

func (s *Service) LockSlot(ctx context.Context, slotID string, req LockSlotReq, isAdmin bool) error {
	_ = isAdmin // 保留扩展：管理员可强制锁定
	lockUntil := time.Now().Add(5 * time.Minute)
	if err := s.timeSlotRepo.Lock(ctx, slotID, req.PatientID, lockUntil); err != nil {
		return err
	}
	return nil
}

func (s *Service) ReleaseSlot(ctx context.Context, slotID string, patientID string, allowForce bool) error {
	if err := s.timeSlotRepo.Release(ctx, slotID, patientID, allowForce); err != nil {
		return err
	}
	return nil
}

func (s *Service) SuspendSchedule(ctx context.Context, req SuspendScheduleReq) (int64, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return 0, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
	}
	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return 0, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "start_time格式应为HH:mm")
	}
	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return 0, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "end_time格式应为HH:mm")
	}
	startAt := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
	endAt := time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
	if !endAt.After(startAt) {
		return 0, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "结束时间必须晚于开始时间")
	}

	if err := s.scheduleRepo.Suspend(ctx, req.DeviceID, date, req.Reason); err != nil {
		return 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	affected, err := s.timeSlotRepo.SuspendRange(ctx, req.DeviceID, date, startAt, endAt, req.Reason)
	if err != nil {
		return 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return affected, nil
}

func (s *Service) SubstituteSchedule(ctx context.Context, req SubstituteScheduleReq) (int64, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return 0, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
	}
	if err := s.scheduleRepo.Substitute(ctx, req.SourceDeviceID, req.TargetDeviceID, date); err != nil {
		return 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	moved, err := s.timeSlotRepo.UpdateDeviceByDate(ctx, req.SourceDeviceID, req.TargetDeviceID, date)
	if err != nil {
		return 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return moved, nil
}

func (s *Service) AddExtraSlots(ctx context.Context, req AddExtraSlotsReq) ([]TimeSlotResp, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "date格式应为YYYY-MM-DD")
	}
	start, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "start_time格式应为HH:mm")
	}
	end, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "end_time格式应为HH:mm")
	}
	startAt := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
	endAt := time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
	if !endAt.After(startAt) {
		return nil, bizErr.NewWithDetail(bizErr.ErrInvalidParam, "结束时间必须晚于开始时间")
	}
	overlap, err := s.timeSlotRepo.HasOverlap(ctx, req.DeviceID, date, startAt, endAt)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if overlap {
		return nil, bizErr.New(bizErr.ErrResExtraSlotOverlap)
	}

	slotDur := time.Duration(req.SlotMinutes) * time.Minute
	stdMin := int(slotDur.Minutes())
	poolType := req.PoolType
	if poolType == "" {
		poolType = "public"
	}
	var slots []TimeSlotResp
	for cur := startAt; !cur.Add(slotDur).After(endAt); cur = cur.Add(slotDur) {
		slots = append(slots, TimeSlotResp{
			ID:              uuid.New().String(),
			DeviceID:         req.DeviceID,
			ExamItemID:       req.ExamItemID,
			PoolType:         poolType,
			StartAt:          cur,
			EndAt:            cur.Add(slotDur),
			Status:           "available",
			StandardDuration: stdMin,
			AdjustedDuration: stdMin,
			Remaining:        1,
		})
	}
	if err := s.timeSlotRepo.BulkCreate(ctx, slots); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return slots, nil
}
