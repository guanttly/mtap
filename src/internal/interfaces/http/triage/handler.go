// Package triage 接口层 - 分诊管理HTTP处理器
// 核心目的：处理triage模块的HTTP请求，参数校验与响应封装
// 模块功能：
//   - 自助机/护士/NFC 签到
//   - 候诊队列查看、呼叫、重叫、过号
//   - 检查执行状态管理（开始/完成/撤销）
package triage

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appTriage "github.com/euler/mtap/internal/application/triage"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/logger"
	"github.com/euler/mtap/pkg/response"
)

// Handler 分诊服务 HTTP 处理器
type Handler struct {
	svc *appTriage.TriageAppService
}

// NewHandler 创建分诊处理器
func NewHandler(svc *appTriage.TriageAppService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	triage := rg.Group("/triage")
	{
		// 签到接口 — 无需鉴权（患者自助）
		checkin := triage.Group("/checkin")
		checkin.POST("/kiosk", h.KioskCheckIn)
		checkin.POST("/nfc", h.NFCCheckIn)

		// 护士签到 — 需要护士及以上权限
		checkinAuth := triage.Group("/checkin")
		checkinAuth.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse"))
		checkinAuth.POST("/nurse", h.NurseCheckIn)

		// 候诊队列 — 护士及以上
		queue := triage.Group("/queue")
		queue.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse"))
		queue.GET("/:roomId", h.GetQueueStatus)

		// 呼叫操作 — 护士及以上
		call := triage.Group("/call")
		call.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse"))
		call.POST("/:roomId/next", h.CallNext)
		call.POST("/:roomId/recall", h.Recall)
		call.POST("/:roomId/miss", h.MissAndRequeue)

		// 检查状态 — 护士/管理员
		exam := triage.Group("/exam")
		exam.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse"))
		exam.GET("/:id", h.GetExamExecution)
		exam.POST("/:id/start", h.StartExam)
		exam.POST("/:id/complete", h.CompleteExam)

		// 撤销操作 — 管理员或护士均可（时间窗口内）
		examUndo := triage.Group("/exam")
		examUndo.Use(middleware.RequireRoles("admin", "scheduler_admin", "nurse"))
		examUndo.POST("/:id/undo", h.UndoExam)
	}
}

// KioskCheckIn 自助机二维码签到
func (h *Handler) KioskCheckIn(c *gin.Context) {
	var req appTriage.KioskCheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.KioskCheckIn(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		Action: "kiosk_checkin", Resource: "checkin",
		ResourceID: resp.CheckInID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, resp)
}

// NurseCheckIn 护士手动签到
func (h *Handler) NurseCheckIn(c *gin.Context) {
	var req appTriage.NurseCheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	resp, err := h.svc.NurseCheckIn(c.Request.Context(), req, operatorID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "nurse_checkin", Resource: "checkin",
		ResourceID: resp.CheckInID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, resp)
}

// NFCCheckIn NFC 感应签到
func (h *Handler) NFCCheckIn(c *gin.Context) {
	var req appTriage.NFCCheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	// NFC 签到：将 CardID 当作预约凭证处理，调用 Kiosk 路径
	kioskReq := appTriage.KioskCheckInReq{QRCodeData: req.CardID}
	resp, err := h.svc.KioskCheckIn(c.Request.Context(), kioskReq)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		Action: "nfc_checkin", Resource: "checkin",
		ResourceID: resp.CheckInID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, resp)
}

// GetQueueStatus 获取候诊队列状态
func (h *Handler) GetQueueStatus(c *gin.Context) {
	roomID := c.Param("roomId")
	resp, err := h.svc.GetQueueStatus(c.Request.Context(), roomID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// CallNext 呼叫下一位患者
func (h *Handler) CallNext(c *gin.Context) {
	roomID := c.Param("roomId")
	operatorID := c.GetString("user_id")
	entry, err := h.svc.CallNext(c.Request.Context(), roomID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "call_next", Resource: "queue_entry",
		ResourceID: entry.ID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, entry)
}

// Recall 重叫当前患者
func (h *Handler) Recall(c *gin.Context) {
	roomID := c.Param("roomId")
	operatorID := c.GetString("user_id")
	entry, err := h.svc.Recall(c.Request.Context(), roomID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "recall", Resource: "queue_entry",
		ResourceID: entry.ID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, entry)
}

// MissAndRequeue 过号重排
func (h *Handler) MissAndRequeue(c *gin.Context) {
	roomID := c.Param("roomId")
	operatorID := c.GetString("user_id")
	entry, err := h.svc.MissAndRequeue(c.Request.Context(), roomID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "miss_requeue", Resource: "queue_entry",
		ResourceID: entry.ID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, entry)
}

// GetExamExecution 获取检查执行状态
func (h *Handler) GetExamExecution(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetExamExecution(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// StartExam 开始检查
func (h *Handler) StartExam(c *gin.Context) {
	id := c.Param("id")
	operatorID := c.GetString("user_id")
	if err := h.svc.StartExam(c.Request.Context(), id, operatorID); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "start_exam", Resource: "exam_execution",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// CompleteExam 检查完成
func (h *Handler) CompleteExam(c *gin.Context) {
	id := c.Param("id")
	operatorID := c.GetString("user_id")
	if err := h.svc.CompleteExam(c.Request.Context(), id, operatorID); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "complete_exam", Resource: "exam_execution",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// UndoExam 撤销检查状态误操作
func (h *Handler) UndoExam(c *gin.Context) {
	id := c.Param("id")
	var req appTriage.UndoExamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := c.GetString("user_id")
	if err := h.svc.UndoExam(c.Request.Context(), id, operatorID, req.Reason); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "undo_exam", Resource: "exam_execution",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}
