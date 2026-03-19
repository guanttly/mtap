package resource

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appRes "github.com/euler/mtap/internal/application/resource"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/pkg/response"
)

type Handler struct {
	svc *appRes.Service
}

func NewHandler(svc *appRes.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	res := rg.Group("/resources")
	{
		adminOnly := res.Group("")
		adminOnly.Use(middleware.RequireRoles("admin", "scheduler_admin"))

		operatorPlus := res.Group("")
		operatorPlus.Use(middleware.RequireRoles("admin", "scheduler_admin", "operator", "nurse", "viewer"))

		adminOnly.POST("/devices", h.CreateDevice)
		operatorPlus.GET("/devices", h.ListDevices)

		adminOnly.POST("/exam-items", h.CreateExamItem)
		operatorPlus.GET("/exam-items", h.ListExamItems)

		adminOnly.POST("/item-aliases", h.CreateAlias)
		operatorPlus.GET("/item-aliases", h.ListAliases)

		adminOnly.POST("/slot-pools", h.CreateSlotPool)
		operatorPlus.GET("/slot-pools", h.ListSlotPools)

		adminOnly.POST("/schedules/generate", h.GenerateSchedule)
		adminOnly.POST("/schedules/suspend", h.SuspendSchedule)
		adminOnly.POST("/schedules/substitute", h.SubstituteSchedule)
		adminOnly.POST("/schedules/add-slots", h.AddExtraSlots)
		operatorPlus.GET("/slots", h.ListSlots)
		operatorPlus.GET("/slots/available", h.QueryAvailableSlots)
		operatorPlus.POST("/slots/:id/lock", h.LockSlot)
		operatorPlus.POST("/slots/:id/release", h.ReleaseSlot)
	}
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var req appRes.CreateDeviceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.CreateDevice(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

func (h *Handler) ListDevices(c *gin.Context) {
	resp, err := h.svc.ListDevices(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) CreateExamItem(c *gin.Context) {
	var req appRes.CreateExamItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.CreateExamItem(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

func (h *Handler) ListExamItems(c *gin.Context) {
	resp, err := h.svc.ListExamItems(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) CreateAlias(c *gin.Context) {
	var req appRes.CreateAliasReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.CreateAlias(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

func (h *Handler) ListAliases(c *gin.Context) {
	examItemID := c.Query("exam_item_id")
	resp, err := h.svc.ListAliases(c.Request.Context(), examItemID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) CreateSlotPool(c *gin.Context) {
	var req appRes.CreateSlotPoolReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.CreateSlotPool(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

func (h *Handler) ListSlotPools(c *gin.Context) {
	resp, err := h.svc.ListSlotPools(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) GenerateSchedule(c *gin.Context) {
	var req appRes.GenerateScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.GenerateSchedule(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, resp)
}

func (h *Handler) SuspendSchedule(c *gin.Context) {
	var req appRes.SuspendScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	affected, err := h.svc.SuspendSchedule(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"released_slots": affected})
}

func (h *Handler) SubstituteSchedule(c *gin.Context) {
	var req appRes.SubstituteScheduleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	moved, err := h.svc.SubstituteSchedule(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, gin.H{"moved_slots": moved})
}

func (h *Handler) AddExtraSlots(c *gin.Context) {
	var req appRes.AddExtraSlotsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	slots, err := h.svc.AddExtraSlots(c.Request.Context(), req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, slots)
}

func (h *Handler) ListSlots(c *gin.Context) {
	deviceID := c.Query("device_id")
	date := c.Query("date")
	if deviceID == "" || date == "" {
		response.Fail(c, http.StatusBadRequest, 1004, "device_id与date为必填")
		return
	}
	resp, err := h.svc.ListSlots(c.Request.Context(), deviceID, date)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) QueryAvailableSlots(c *gin.Context) {
	var req appRes.QueryAvailableSlotsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	resp, err := h.svc.QueryAvailableSlots(c.Request.Context(), req.DeviceID, req.Date, req.ExamItemID, req.PoolType, req.PatientAge)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithData(c, resp)
}

func (h *Handler) LockSlot(c *gin.Context) {
	slotID := c.Param("id")
	var req appRes.LockSlotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	role := c.GetString("role")
	isAdmin := role == "admin" || role == "scheduler_admin"
	if err := h.svc.LockSlot(c.Request.Context(), slotID, req, isAdmin); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c)
}

func (h *Handler) ReleaseSlot(c *gin.Context) {
	slotID := c.Param("id")
	var req appRes.LockSlotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 1004, err.Error())
		return
	}
	role := c.GetString("role")
	allowForce := role == "admin" || role == "scheduler_admin"
	if err := h.svc.ReleaseSlot(c.Request.Context(), slotID, req.PatientID, allowForce); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c)
}
