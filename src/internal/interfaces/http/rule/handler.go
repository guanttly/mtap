// Package rule 接口层 - 规则管理 HTTP Handler
package rule

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appRule "github.com/euler/mtap/internal/application/rule"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/logger"
	"github.com/euler/mtap/pkg/response"
)

// Handler 规则管理 HTTP 处理器
type Handler struct {
	appService *appRule.RuleAppService
}

// NewHandler 创建规则处理器
func NewHandler(appService *appRule.RuleAppService) *Handler {
	return &Handler{appService: appService}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rules := rg.Group("/rules")
	{
		adminOnly := rules.Group("")
		adminOnly.Use(middleware.RequireRoles("admin", "scheduler_admin"))

		operatorPlus := rules.Group("")
		operatorPlus.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse", "viewer"))

		// 冲突规则
		adminOnly.POST("/conflicts", h.CreateConflictRule)
		operatorPlus.GET("/conflicts", h.ListConflictRules)
		operatorPlus.GET("/conflicts/:id", h.GetConflictRule)
		adminOnly.DELETE("/conflicts/:id", h.DeleteConflictRule)

		// 冲突包
		adminOnly.POST("/conflict-packages", h.CreateConflictPackage)
		operatorPlus.GET("/conflict-packages", h.ListConflictPackages)
		adminOnly.DELETE("/conflict-packages/:id", h.DeleteConflictPackage)

		// 依赖规则
		adminOnly.POST("/dependencies", h.CreateDependencyRule)
		operatorPlus.GET("/dependencies", h.ListDependencyRules)
		adminOnly.DELETE("/dependencies/:id", h.DeleteDependencyRule)

		// 优先级标签
		adminOnly.POST("/priority-tags", h.CreatePriorityTag)
		operatorPlus.GET("/priority-tags", h.ListPriorityTags)
		adminOnly.DELETE("/priority-tags/:id", h.DeletePriorityTag)

		// 综合规则校验
		operatorPlus.POST("/check", h.CheckRules)

		// 排序策略
		adminOnly.POST("/sorting-strategies", h.SaveSortingStrategy)
		operatorPlus.GET("/sorting-strategies", h.GetSortingStrategy)

		// 患者属性适配
		adminOnly.POST("/patient-adapt", h.SavePatientAdaptRules)
		operatorPlus.GET("/patient-adapt", h.ListPatientAdaptRules)

		// 开单来源控制
		adminOnly.POST("/source-controls", h.SaveSourceControls)
		operatorPlus.GET("/source-controls", h.ListSourceControls)
	}
}

// === 冲突规则 ===

func (h *Handler) CreateConflictRule(c *gin.Context) {
	var req appRule.CreateConflictRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}

	operatorID := c.GetString("user_id")
	resp, err := h.appService.CreateConflictRule(c.Request.Context(), req, operatorID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: operatorID,
		Action:     "create",
		Resource:   "conflict_rule",
		ResourceID: resp.ID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.Created(c, resp)
}

func (h *Handler) GetConflictRule(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.appService.GetConflictRule(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) ListConflictRules(c *gin.Context) {
	var req appRule.ListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.ListConflictRules(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) DeleteConflictRule(c *gin.Context) {
	id := c.Param("id")
	if err := h.appService.DeleteConflictRule(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "delete",
		Resource:   "conflict_rule",
		ResourceID: id,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

// === 冲突包 ===

func (h *Handler) CreateConflictPackage(c *gin.Context) {
	var req appRule.CreateConflictPackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.CreateConflictPackage(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "create",
		Resource:   "conflict_package",
		ResourceID: resp.ID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.Created(c, resp)
}

func (h *Handler) ListConflictPackages(c *gin.Context) {
	resp, err := h.appService.ListConflictPackages(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) DeleteConflictPackage(c *gin.Context) {
	id := c.Param("id")
	if err := h.appService.DeleteConflictPackage(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "delete",
		Resource:   "conflict_package",
		ResourceID: id,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

// === 依赖规则 ===

func (h *Handler) CreateDependencyRule(c *gin.Context) {
	var req appRule.CreateDependencyRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.CreateDependencyRule(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "create",
		Resource:   "dependency_rule",
		ResourceID: resp.ID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.Created(c, resp)
}

func (h *Handler) ListDependencyRules(c *gin.Context) {
	var req appRule.ListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.ListDependencyRules(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) DeleteDependencyRule(c *gin.Context) {
	id := c.Param("id")
	if err := h.appService.DeleteDependencyRule(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "delete",
		Resource:   "dependency_rule",
		ResourceID: id,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

// === 优先级标签 ===

func (h *Handler) CreatePriorityTag(c *gin.Context) {
	var req appRule.CreatePriorityTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.CreatePriorityTag(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "create",
		Resource:   "priority_tag",
		ResourceID: resp.ID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.Created(c, resp)
}

func (h *Handler) ListPriorityTags(c *gin.Context) {
	resp, err := h.appService.ListPriorityTags(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) DeletePriorityTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.appService.DeletePriorityTag(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "delete",
		Resource:   "priority_tag",
		ResourceID: id,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

// === 综合规则校验 ===

func (h *Handler) CheckRules(c *gin.Context) {
	var req appRule.RuleCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.CheckRules(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "check",
		Resource:   "rule_check",
		ResourceID: req.PatientID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OKWithData(c, resp)
}

// === 排序策略 ===

func (h *Handler) SaveSortingStrategy(c *gin.Context) {
	var req appRule.SaveSortingStrategyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.appService.SaveSortingStrategy(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "create",
		Resource:   "sorting_strategy",
		ResourceID: resp.ID,
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.Created(c, resp)
}

func (h *Handler) GetSortingStrategy(c *gin.Context) {
	var scope appRule.EffectiveScopeDTO
	if err := c.ShouldBindQuery(&scope); err != nil {
		// 允许 query 为空时走后续校验
	}
	resp, err := h.appService.GetSortingStrategy(c.Request.Context(), scope)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// === 患者属性适配 ===

func (h *Handler) SavePatientAdaptRules(c *gin.Context) {
	var req []appRule.SavePatientAdaptRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	if err := h.appService.SavePatientAdaptRules(c.Request.Context(), req); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "save",
		Resource:   "patient_adapt_rules",
		ResourceID: "all",
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

func (h *Handler) ListPatientAdaptRules(c *gin.Context) {
	resp, err := h.appService.ListPatientAdaptRules(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

// === 开单来源控制 ===

func (h *Handler) SaveSourceControls(c *gin.Context) {
	var req []appRule.SaveSourceControlReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	if err := h.appService.SaveSourceControls(c.Request.Context(), req); err != nil {
		response.FailWithError(c, err)
		return
	}
	_ = logger.Audit(c.Request.Context(), logger.AuditEntry{
		OperatorID: c.GetString("user_id"),
		Action:     "save",
		Resource:   "source_controls",
		ResourceID: "all",
		IP:         c.ClientIP(),
		Timestamp:  time.Now(),
	})
	response.OK(c)
}

func (h *Handler) ListSourceControls(c *gin.Context) {
	resp, err := h.appService.ListSourceControls(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}
