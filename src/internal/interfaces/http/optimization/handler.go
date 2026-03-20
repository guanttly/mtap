// Package optimization 接口层 - 效能优化 HTTP Handler
package optimization

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appOpt "github.com/euler/mtap/internal/application/optimization"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/logger"
	"github.com/euler/mtap/pkg/response"
)

// Handler 效能优化 HTTP 处理器
type Handler struct {
	svc *appOpt.OptimizationAppService
}

// NewHandler 创建处理器
func NewHandler(svc *appOpt.OptimizationAppService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	opt := rg.Group("/optimization")
	opt.Use(middleware.RequireRoles("admin", "scheduler_admin"))
	{
		// 指标
		opt.GET("/metrics", h.ListMetrics)
		opt.GET("/metrics/:code/trend", h.GetMetricTrend)
		// 告警
		opt.GET("/alerts", h.ListAlerts)
		opt.GET("/alerts/:id", h.GetAlert)
		opt.PUT("/alerts/:id/dismiss", h.DismissAlert)
		// 策略
		opt.GET("/strategies", h.ListStrategies)
		opt.GET("/strategies/:id", h.GetStrategy)
		opt.POST("/strategies/:id/approve", h.ApproveStrategy)
		opt.POST("/strategies/:id/reject", h.RejectStrategy)
		opt.POST("/strategies/:id/rollback", h.RollbackStrategy)
		opt.POST("/strategies/:id/promote", h.PromoteStrategy)
		// 试运行与评估
		opt.GET("/trials/:id/monitor", h.GetTrialMonitor)
		opt.GET("/evaluations/:id", h.GetEvaluation)
		// ROI 报告与扫描
		opt.GET("/roi-reports/:id", h.GetROIReport)
		opt.GET("/roi-reports/:id/export", h.ExportROIReport)
		opt.POST("/roi-reports/:id/result", h.SubmitROIResult)
		opt.GET("/scans", h.ListScans)
		opt.GET("/scans/:id", h.GetScan)
	}
}

// ListMetrics GET /api/v1/optimization/metrics
func (h *Handler) ListMetrics(c *gin.Context) {
	resps, err := h.svc.ListMetrics(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resps)
}

// GetMetricTrend GET /api/v1/optimization/metrics/:code/trend
func (h *Handler) GetMetricTrend(c *gin.Context) {
	code := c.Param("code")
	var req appOpt.GetMetricTrendReq
	_ = c.ShouldBindQuery(&req)

	resps, err := h.svc.GetMetricTrend(c.Request.Context(), code, req.Days)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resps)
}

// ListAlerts GET /api/v1/optimization/alerts
func (h *Handler) ListAlerts(c *gin.Context) {
	var req appOpt.ListAlertsReq
	_ = c.ShouldBindQuery(&req)

	alerts, total, err := h.svc.ListAlerts(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"list": alerts, "total": total})
}

// GetAlert GET /api/v1/optimization/alerts/:id
func (h *Handler) GetAlert(c *gin.Context) {
	id := c.Param("id")
	alert, err := h.svc.GetAlert(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, alert)
}

// DismissAlert PUT /api/v1/optimization/alerts/:id/dismiss
func (h *Handler) DismissAlert(c *gin.Context) {
	id := c.Param("id")
	var req appOpt.DismissAlertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	if err := h.svc.DismissAlert(c.Request.Context(), id, req); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c)
}

// ListStrategies GET /api/v1/optimization/strategies
func (h *Handler) ListStrategies(c *gin.Context) {
	var req appOpt.ListStrategiesReq
	_ = c.ShouldBindQuery(&req)

	strategies, total, err := h.svc.ListStrategies(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"list": strategies, "total": total})
}

// GetStrategy GET /api/v1/optimization/strategies/:id
func (h *Handler) GetStrategy(c *gin.Context) {
	id := c.Param("id")
	st, err := h.svc.GetStrategy(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, st)
}

// ApproveStrategy POST /api/v1/optimization/strategies/:id/approve
func (h *Handler) ApproveStrategy(c *gin.Context) {
	id := c.Param("id")
	var req appOpt.ApproveStrategyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%v", c.MustGet("user_id"))

	st, err := h.svc.ApproveStrategy(c.Request.Context(), id, operatorID, req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "approve_strategy", Resource: "strategy",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OKWithData(c, st)
}

// RejectStrategy POST /api/v1/optimization/strategies/:id/reject
func (h *Handler) RejectStrategy(c *gin.Context) {
	id := c.Param("id")
	var req appOpt.RejectStrategyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%v", c.MustGet("user_id"))

	if err := h.svc.RejectStrategy(c.Request.Context(), id, req); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "reject_strategy", Resource: "strategy",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// RollbackStrategy POST /api/v1/optimization/strategies/:id/rollback
func (h *Handler) RollbackStrategy(c *gin.Context) {
	id := c.Param("id")
	operatorID := fmt.Sprintf("%v", c.MustGet("user_id"))

	if err := h.svc.RollbackStrategy(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "rollback_strategy", Resource: "strategy",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// PromoteStrategy POST /api/v1/optimization/strategies/:id/promote
func (h *Handler) PromoteStrategy(c *gin.Context) {
	id := c.Param("id")
	operatorID := fmt.Sprintf("%v", c.MustGet("user_id"))

	if err := h.svc.PromoteStrategy(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID, Action: "promote_strategy", Resource: "strategy",
		ResourceID: id, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.OK(c)
}

// GetTrialMonitor GET /api/v1/optimization/trials/:id/monitor
func (h *Handler) GetTrialMonitor(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetTrialMonitor(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// GetEvaluation GET /api/v1/optimization/evaluations/:id
func (h *Handler) GetEvaluation(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEvaluation(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// GetROIReport GET /api/v1/optimization/roi-reports/:id
func (h *Handler) GetROIReport(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetROIReport(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// SubmitROIResult POST /api/v1/optimization/roi-reports/:id/result
func (h *Handler) SubmitROIResult(c *gin.Context) {
	id := c.Param("id")
	var req appOpt.ROIApprovalResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID := fmt.Sprintf("%v", c.MustGet("user_id"))

	if err := h.svc.SubmitROIApprovalResult(c.Request.Context(), id, operatorID, req); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c)
}

// ListScans GET /api/v1/optimization/scans
func (h *Handler) ListScans(c *gin.Context) {
	var req appOpt.ListScansReq
	_ = c.ShouldBindQuery(&req)

	scans, total, err := h.svc.ListScans(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"list": scans, "total": total})
}

// GetScan GET /api/v1/optimization/scans/:id
func (h *Handler) GetScan(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetScan(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// ExportROIReport GET /api/v1/optimization/roi-reports/:id/export
func (h *Handler) ExportROIReport(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.ExportROIReport(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	// 返回报告数据，客户端可据此生成PDF
	c.Header("Content-Disposition", "attachment; filename=\"roi-report-"+id+".json\"")
	response.OKWithData(c, resp)
}
