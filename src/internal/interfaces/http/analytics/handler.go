// Package analytics 接口层 - 统计分析 HTTP Handler
package analytics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appAnalytics "github.com/euler/mtap/internal/application/analytics"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/logger"
	"github.com/euler/mtap/pkg/response"
)

// Handler 统计分析 HTTP 处理器
type Handler struct {
	svc *appAnalytics.AnalyticsAppService
}

// NewHandler 创建处理器
func NewHandler(svc *appAnalytics.AnalyticsAppService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	analytics := rg.Group("/analytics")
	analytics.Use(middleware.RequireRoles("admin", "scheduler_admin"))
	{
		analytics.GET("/dashboard", h.GetDashboard)
		analytics.GET("/dashboard/device/:id", h.GetDeviceDetail)
		analytics.GET("/reports", h.ListReports)
		analytics.POST("/reports", h.CreateReport)
		analytics.GET("/reports/:id", h.GetReport)
		analytics.GET("/reports/:id/export", h.ExportReport)
	}
}

// GetDashboard GET /api/v1/analytics/dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	var req appAnalytics.GetDashboardReq
	_ = c.ShouldBindQuery(&req)

	resp, err := h.svc.GetDashboard(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// GetDeviceDetail GET /api/v1/analytics/dashboard/device/:id
func (h *Handler) GetDeviceDetail(c *gin.Context) {
	deviceID := c.Param("id")
	if deviceID == "" {
		response.Fail(c, http.StatusBadRequest, 1004, "设备ID不能为空")
		return
	}
	dateStr := c.Query("date")

	resp, err := h.svc.GetDeviceDetail(c.Request.Context(), deviceID, dateStr)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// ListReports GET /api/v1/analytics/reports
func (h *Handler) ListReports(c *gin.Context) {
	var req appAnalytics.ListReportsReq
	_ = c.ShouldBindQuery(&req)

	reports, total, err := h.svc.ListReports(c.Request.Context(), req.Page, req.Size)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	pageSize := req.Size
	if pageSize <= 0 {
		pageSize = 20
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (total + int64(pageSize) - 1) / int64(pageSize)
	}
	response.OKWithData(c, gin.H{
		"items":       reports,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// CreateReport POST /api/v1/analytics/reports
func (h *Handler) CreateReport(c *gin.Context) {
	var req appAnalytics.CreateReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	operatorID, _ := c.Get("user_id")
	operatorIDStr := fmt.Sprintf("%v", operatorID)

	report, err := h.svc.CreateReport(c.Request.Context(), req, operatorIDStr)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorIDStr, Action: "create_report", Resource: "report",
		ResourceID: report.ID, IP: c.ClientIP(), Timestamp: time.Now(),
	})
	response.Created(c, report)
}

// GetReport GET /api/v1/analytics/reports/:id
func (h *Handler) GetReport(c *gin.Context) {
	id := c.Param("id")
	report, err := h.svc.GetReport(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, report)
}

// ExportReport GET /api/v1/analytics/reports/:id/export
func (h *Handler) ExportReport(c *gin.Context) {
	id := c.Param("id")
	report, data, err := h.svc.ExportReport(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	if report == nil || data == nil {
		response.Fail(c, http.StatusAccepted, 0, "报表生成中，请稍候再试")
		return
	}
	contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if report.Format == "pdf" {
		contentType = "application/pdf"
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="report_%s.%s"`, id, report.Format))
	c.Data(http.StatusOK, contentType, data)
}
