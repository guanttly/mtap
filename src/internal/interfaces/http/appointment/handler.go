// Package appointment 接口层 - 预约服务 HTTP Handler
package appointment

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appAppt "github.com/euler/mtap/internal/application/appointment"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/logger"
	"github.com/euler/mtap/pkg/response"
)

// Handler 预约服务 HTTP 处理器
type Handler struct {
	svc *appAppt.AppointmentAppService
}

// NewHandler 创建预约处理器
func NewHandler(svc *appAppt.AppointmentAppService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	appt := rg.Group("/appointments")
	{
		// 所有角色（操作员+）
		opPlus := appt.Group("")
		opPlus.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse", "viewer"))

		opPlus.GET("", h.ListAppointments)
		opPlus.GET("/:id", h.GetAppointment)
		opPlus.GET("/:id/credential", h.GetCredential)
		opPlus.POST("/auto", h.AutoAppointment)
		opPlus.POST("/combo", h.ComboAppointment)
		opPlus.POST("/combo/confirm", h.ComboConfirm)
		opPlus.PUT("/:id/confirm", h.ConfirmAppointment)
		opPlus.PUT("/:id/reschedule", h.RescheduleAppointment)
		opPlus.PUT("/:id/cancel", h.CancelAppointment)

		// 管理员专属
		adminOnly := appt.Group("")
		adminOnly.Use(middleware.RequireRoles("admin", "scheduler_admin"))
		adminOnly.POST("/manual", h.ManualAppointment)

		// 黑名单管理
		bl := appt.Group("/blacklist")
		bl.Use(middleware.RequireRoles("admin", "scheduler_admin"))
		bl.GET("", h.ListBlacklists)
		bl.GET("/:id", h.GetBlacklist)
		bl.POST("/:id/appeal", h.SubmitAppeal)

		// 申诉审核（管理员）
		appeals := appt.Group("")
		appeals.Use(middleware.RequireRoles("admin"))
		appeals.PUT("/blacklist/appeals/:id/review", h.ReviewAppeal)
	}
}

// GetAppointment 获取预约详情
func (h *Handler) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetAppointment(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// ListAppointments 查询预约列表
func (h *Handler) ListAppointments(c *gin.Context) {
	var req appAppt.ListAppointmentReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	list, total, err := h.svc.ListAppointments(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"items": list, "total": total, "page": req.Page, "page_size": req.PageSize})
}

// AutoAppointment 一键自动预约
func (h *Handler) AutoAppointment(c *gin.Context) {
	var req appAppt.AutoAppointmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	resp, err := h.svc.AutoAppointment(c.Request.Context(), req, operatorID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "create", Resource: "appointment",
		ResourceID: resp.AppointmentID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.Created(c, resp)
}

// ComboAppointment 组合预约计算
func (h *Handler) ComboAppointment(c *gin.Context) {
	var req appAppt.ComboAppointmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	// 简化实现：返回空方案（实际需接入号源服务）
	response.OKWithData(c, gin.H{"plans": []interface{}{}, "message": "请配置号源服务"})
}

// ComboConfirm 确认组合预约方案
func (h *Handler) ComboConfirm(c *gin.Context) {
	var req appAppt.ComboConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.ConfirmAppointment(c.Request.Context(), req.AppointmentID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// ManualAppointment 人工干预预约
func (h *Handler) ManualAppointment(c *gin.Context) {
	var req appAppt.ManualAppointmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	resp, err := h.svc.ManualAppointment(c.Request.Context(), req, operatorID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "manual_override", Resource: "appointment",
		ResourceID: resp.ID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.Created(c, resp)
}

// ConfirmAppointment 确认预约
func (h *Handler) ConfirmAppointment(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.ConfirmAppointment(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// RescheduleAppointment 改约
func (h *Handler) RescheduleAppointment(c *gin.Context) {
	id := c.Param("id")
	var req appAppt.RescheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	resp, err := h.svc.RescheduleAppointment(c.Request.Context(), id, req, operatorID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "reschedule", Resource: "appointment",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, resp)
}

// CancelAppointment 取消预约
func (h *Handler) CancelAppointment(c *gin.Context) {
	id := c.Param("id")
	var req appAppt.CancelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	if err := h.svc.CancelAppointment(c.Request.Context(), id, req, operatorID); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "cancel", Resource: "appointment",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// GetCredential 获取预约凭证
func (h *Handler) GetCredential(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCredential(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// ListBlacklists 查询黑名单列表
func (h *Handler) ListBlacklists(c *gin.Context) {
	var req appAppt.ListBlacklistReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	list, total, err := h.svc.ListBlacklists(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"items": list, "total": total})
}

// GetBlacklist 获取黑名单详情
func (h *Handler) GetBlacklist(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetBlacklist(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// SubmitAppeal 提交申诉
func (h *Handler) SubmitAppeal(c *gin.Context) {
	blacklistID := c.Param("id")
	var req appAppt.AppealReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.SubmitAppeal(c.Request.Context(), blacklistID, req.Reason)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

// ReviewAppeal 审核申诉
func (h *Handler) ReviewAppeal(c *gin.Context) {
	appealID := c.Param("id")
	var req appAppt.ReviewAppealReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	reviewerID := c.GetString("user_id")
	if err := h.svc.ReviewAppeal(c.Request.Context(), appealID, reviewerID, req.Approved); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: reviewerID, Action: "review_appeal", Resource: "appeal",
		ResourceID: appealID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}
