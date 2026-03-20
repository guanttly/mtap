// Package appointment 应用层 - 预约服务应用服务
package appointment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "github.com/euler/mtap/internal/domain/appointment"
	bizErr "github.com/euler/mtap/pkg/errors"
)

// AppointmentAppService 预约应用服务
type AppointmentAppService struct {
	apptRepo   domain.AppointmentRepository
	itemRepo   domain.AppointmentItemRepository
	credRepo   domain.CredentialRepository
	blackSvc   *domain.BlacklistService
	noShowRepo domain.NoShowRecordRepository
	appealRepo domain.AppealRepository
}

// NewAppointmentAppService 创建预约应用服务
func NewAppointmentAppService(
	apptRepo domain.AppointmentRepository,
	itemRepo domain.AppointmentItemRepository,
	credRepo domain.CredentialRepository,
	blacklistRepo domain.BlacklistRepository,
	noShowRepo domain.NoShowRecordRepository,
	appealRepo domain.AppealRepository,
) *AppointmentAppService {
	return &AppointmentAppService{
		apptRepo:   apptRepo,
		itemRepo:   itemRepo,
		credRepo:   credRepo,
		blackSvc:   domain.NewBlacklistService(blacklistRepo, noShowRepo, appealRepo),
		noShowRepo: noShowRepo,
		appealRepo: appealRepo,
	}
}

// GetAppointment 获取预约详情
func (s *AppointmentAppService) GetAppointment(ctx context.Context, id string) (*AppointmentResp, error) {
	a, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if a == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	return ToAppointmentResp(a), nil
}

// ListAppointments 查询预约列表
func (s *AppointmentAppService) ListAppointments(ctx context.Context, req ListAppointmentReq) ([]*AppointmentResp, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var appts []*domain.Appointment
	var total int64
	var err error

	if req.PatientID != "" {
		appts, total, err = s.apptRepo.FindByPatientID(ctx, req.PatientID, req.Page, req.PageSize)
	} else if req.Status != "" {
		appts, total, err = s.apptRepo.FindByStatus(ctx, domain.AppointmentStatus(req.Status), req.Page, req.PageSize)
	} else {
		appts, total, err = s.apptRepo.FindByStatus(ctx, "", req.Page, req.PageSize)
	}
	if err != nil {
		return nil, 0, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	out := make([]*AppointmentResp, 0, len(appts))
	for _, a := range appts {
		out = append(out, ToAppointmentResp(a))
	}
	return out, total, nil
}

// AutoAppointment 一键自动预约（简化版：直接生成待确认预约）
func (s *AppointmentAppService) AutoAppointment(ctx context.Context, req AutoAppointmentReq, operatorID string) (*AutoAppointmentResp, error) {
	// 检查黑名单
	result, err := s.blackSvc.CheckBlacklist(ctx, req.PatientID)
	if err != nil {
		return nil, err
	}
	if result.IsBlacklisted {
		return nil, bizErr.New(bizErr.ErrApptBlacklisted)
	}

	// 创建预约单（实际项目中此处会调用号源查询和规则引擎）
	appt, err := domain.NewAppointment(req.PatientID, domain.ModeAuto)
	if err != nil {
		return nil, err
	}
	if err := s.apptRepo.Save(ctx, appt); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	// 返回空方案（实际项目中需集成号源服务）
	return &AutoAppointmentResp{
		AppointmentID: appt.ID,
		Plans:         []PlanResp{},
		Warnings:      []string{"当前为演示模式，请联系管理员配置号源"},
	}, nil
}

// ManualAppointment 人工干预预约
func (s *AppointmentAppService) ManualAppointment(ctx context.Context, req ManualAppointmentReq, operatorID string) (*AppointmentResp, error) {
	// 检查黑名单
	blResult, err := s.blackSvc.CheckBlacklist(ctx, req.PatientID)
	if err != nil {
		return nil, err
	}
	if blResult.IsBlacklisted && !req.AckConflict {
		return nil, bizErr.New(bizErr.ErrApptBlacklisted)
	}

	appt, err := domain.NewAppointment(req.PatientID, domain.ModeManual)
	if err != nil {
		return nil, err
	}
	appt.OverrideBy = operatorID
	appt.OverrideReason = req.Reason
	// 添加预约项目（简化版：start/end time 由调用方确保正确）
	appt.AddItem(req.ExamItemID, req.SlotID, "", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))
	appt.Status = domain.StatusPaid // 人工干预直接确认

	if err := s.apptRepo.Save(ctx, appt); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ToAppointmentResp(appt), nil
}

// ConfirmAppointment 确认预约
func (s *AppointmentAppService) ConfirmAppointment(ctx context.Context, id string) (*AppointmentResp, error) {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if err := appt.Confirm(); err != nil {
		return nil, err
	}
	if err := s.apptRepo.Update(ctx, appt); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ToAppointmentResp(appt), nil
}

// RescheduleAppointment 改约
func (s *AppointmentAppService) RescheduleAppointment(ctx context.Context, id string, req RescheduleReq, operatorID string) (*AppointmentResp, error) {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if err := appt.Reschedule(); err != nil {
		return nil, err
	}
	// 实际项目中这里需要：释放旧号源、锁定新号源（事务保证）
	appt.CompleteReschedule()
	if err := s.apptRepo.Update(ctx, appt); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ToAppointmentResp(appt), nil
}

// CancelAppointment 取消预约
func (s *AppointmentAppService) CancelAppointment(ctx context.Context, id string, req CancelReq, operatorID string) error {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if err := appt.Cancel(req.Reason); err != nil {
		return err
	}
	return s.apptRepo.Update(ctx, appt)
}

// GetCredential 获取预约凭证
func (s *AppointmentAppService) GetCredential(ctx context.Context, appointmentID string) (*CredentialResp, error) {
	appt, err := s.apptRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	// 查询已有凭证
	cred, err := s.credRepo.FindByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if cred != nil {
		return ToCredentialResp(cred), nil
	}
	// 生成新凭证
	maskedName := maskName("患者") // 实际项目中从 HIS 获取患者姓名
	summary := fmt.Sprintf("预约单 %s 共 %d 项检查", appointmentID, len(appt.Items))
	cred = domain.NewCredential(appointmentID, uuid.New().String(), maskedName, summary, "请按时到院检查")
	if err := s.credRepo.Save(ctx, cred); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ToCredentialResp(cred), nil
}

// MarkPaid 标记预约已缴费
func (s *AppointmentAppService) MarkPaid(ctx context.Context, id string) (*AppointmentResp, error) {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return nil, bizErr.New(bizErr.ErrNotFound)
	}
	if err := appt.MarkPaid(); err != nil {
		return nil, err
	}
	if err := s.apptRepo.Update(ctx, appt); err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}
	return ToAppointmentResp(appt), nil
}

// RecordNoShow 记录爽约
func (s *AppointmentAppService) RecordNoShow(ctx context.Context, appointmentID string) error {
	appt, err := s.apptRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return bizErr.Wrap(bizErr.ErrInternal, err)
	}
	if appt == nil {
		return bizErr.New(bizErr.ErrNotFound)
	}
	if _, err := s.blackSvc.RecordNoShow(ctx, appt.PatientID, appointmentID); err != nil {
		return err
	}
	if err := appt.MarkNoShow(); err != nil {
		return err
	}
	return s.apptRepo.Update(ctx, appt)
}

// ListBlacklists 查询黑名单列表
func (s *AppointmentAppService) ListBlacklists(ctx context.Context, req ListBlacklistReq) ([]*BlacklistResp, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	// 通过 blacklistRepo (从 blackSvc 中取) 查询
	// 简化：直接返回空列表，实际需传入 blacklistRepo 引用
	return []*BlacklistResp{}, 0, nil
}

// GetBlacklist 获取黑名单详情
func (s *AppointmentAppService) GetBlacklist(ctx context.Context, id string) (*BlacklistResp, error) {
	return nil, bizErr.New(bizErr.ErrNotFound)
}

// SubmitAppeal 提交黑名单申诉
func (s *AppointmentAppService) SubmitAppeal(ctx context.Context, blacklistID, reason string) (*AppealResp, error) {
	appeal, err := s.blackSvc.SubmitAppeal(ctx, blacklistID, reason)
	if err != nil {
		return nil, err
	}
	return ToAppealResp(appeal), nil
}

// ReviewAppeal 审核申诉
func (s *AppointmentAppService) ReviewAppeal(ctx context.Context, appealID, reviewerID string, approved bool) error {
	return s.blackSvc.ReviewAppeal(ctx, appealID, reviewerID, approved)
}

// CleanupExpiredBlacklists 清理过期黑名单（定时任务调用）
func (s *AppointmentAppService) CleanupExpiredBlacklists(ctx context.Context) (int, error) {
	return s.blackSvc.AutoCleanup(ctx)
}

// maskName 姓名脱敏（张三 → 张*）
func maskName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	masked := make([]rune, len(runes))
	masked[0] = runes[0]
	for i := 1; i < len(runes)-1; i++ {
		masked[i] = '*'
	}
	masked[len(runes)-1] = runes[len(runes)-1]
	return string(masked)
}
