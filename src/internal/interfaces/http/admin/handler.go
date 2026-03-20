// Package admin 接口层 - 用户与角色管理 HTTP Handler
package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	bizErr "github.com/euler/mtap/pkg/errors"
	"github.com/euler/mtap/pkg/response"
)

// Handler 用户/角色管理处理器
type Handler struct {
	db *gorm.DB
}

// NewHandler 创建处理器
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// RegisterRoutes 注册路由（仅 admin 角色可访问）
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/admin")
	g.Use(middleware.RequireRoles("admin"))

	// 用户管理
	g.GET("/users", h.ListUsers)
	g.POST("/users", h.CreateUser)
	g.PUT("/users/:id", h.UpdateUser)
	g.DELETE("/users/:id", h.DeleteUser)
	g.POST("/users/:id/reset-password", h.ResetPassword)

	// 角色管理
	g.GET("/roles", h.ListRoles)
	g.POST("/roles", h.CreateRole)
	g.PUT("/roles/:id", h.UpdateRole)
	g.DELETE("/roles/:id", h.DeleteRole)
}

// ─── 请求/响应结构体 ─────────────────────────────────────────────

type userResp struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	RealName     string     `json:"real_name"`
	RoleID       string     `json:"role_id"`
	RoleName     string     `json:"role_name"`
	DepartmentID string     `json:"department_id,omitempty"`
	Status       string     `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type createUserReq struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Password     string `json:"password" binding:"required,min=6"`
	RealName     string `json:"real_name"`
	RoleID       string `json:"role_id" binding:"required"`
	DepartmentID string `json:"department_id"`
}

type updateUserReq struct {
	RealName     string `json:"real_name"`
	RoleID       string `json:"role_id"`
	DepartmentID string `json:"department_id"`
	Status       string `json:"status"`
}

type resetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type roleResp struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions"`
	IsPreset    bool      `json:"is_preset"`
	CreatedAt   time.Time `json:"created_at"`
}

type createRoleReq struct {
	Name        string   `json:"name" binding:"required,min=2,max=30"`
	Permissions []string `json:"permissions"`
}

type updateRoleReq struct {
	Permissions []string `json:"permissions"`
}

// ─── 用户管理 ────────────────────────────────────────────────────

// ListUsers GET /api/v1/admin/users
func (h *Handler) ListUsers(c *gin.Context) {
	page := 1
	pageSize := 20
	keyword := c.Query("keyword")
	roleID := c.Query("role_id")
	status := c.Query("status")

	var users []po.UserPO
	query := h.db.WithContext(c.Request.Context()).Model(&po.UserPO{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&users)

	// 批量获取角色名
	roleIDs := make([]string, 0, len(users))
	for _, u := range users {
		roleIDs = append(roleIDs, u.RoleID)
	}
	var roles []po.RolePO
	if len(roleIDs) > 0 {
		h.db.WithContext(c.Request.Context()).Where("id IN ?", roleIDs).Find(&roles)
	}
	roleMap := make(map[string]string, len(roles))
	for _, r := range roles {
		roleMap[r.ID] = r.Name
	}

	resps := make([]userResp, 0, len(users))
	for _, u := range users {
		resps = append(resps, userResp{
			ID:           u.ID,
			Username:     u.Username,
			RealName:     u.RealName,
			RoleID:       u.RoleID,
			RoleName:     roleMap[u.RoleID],
			DepartmentID: u.DepartmentID,
			Status:       u.Status,
			LastLoginAt:  u.LastLoginAt,
			CreatedAt:    u.CreatedAt,
		})
	}

	response.OKWithData(c, gin.H{
		"items":     resps,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateUser POST /api/v1/admin/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	var count int64
	h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).
		Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		response.Fail(c, http.StatusConflict, bizErr.ErrDuplicate, "用户名已存在")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "密码加密失败")
		return
	}

	user := po.UserPO{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		RoleID:       req.RoleID,
		DepartmentID: req.DepartmentID,
		Status:       "active",
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "创建用户失败")
		return
	}

	var role po.RolePO
	h.db.WithContext(c.Request.Context()).Where("id = ?", user.RoleID).First(&role)

	response.OKWithData(c, userResp{
		ID:           user.ID,
		Username:     user.Username,
		RealName:     user.RealName,
		RoleID:       user.RoleID,
		RoleName:     role.Name,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt,
	})
}

// UpdateUser PUT /api/v1/admin/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.RoleID != "" {
		updates["role_id"] = req.RoleID
	}
	if req.DepartmentID != "" {
		updates["department_id"] = req.DepartmentID
	}
	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" {
			response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, "status 只允许 active/inactive")
			return
		}
		updates["status"] = req.Status
	}

	if err := h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).
		Where("id = ?", id).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "更新失败")
		return
	}
	response.OK(c)
}

// DeleteUser DELETE /api/v1/admin/users/:id  （软删除：置为 inactive）
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	operatorID := c.GetString("user_id")
	if id == operatorID {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, "不能停用自己的账号")
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).
		Where("id = ?", id).
		Update("status", "inactive").Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "操作失败")
		return
	}
	response.OK(c)
}

// ResetPassword POST /api/v1/admin/users/:id/reset-password
func (h *Handler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "密码加密失败")
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).
		Where("id = ?", id).Update("password_hash", string(hash)).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "重置密码失败")
		return
	}
	response.OK(c)
}

// ─── 角色管理 ────────────────────────────────────────────────────

// ListRoles GET /api/v1/admin/roles
func (h *Handler) ListRoles(c *gin.Context) {
	var roles []po.RolePO
	h.db.WithContext(c.Request.Context()).Order("created_at ASC").Find(&roles)

	resps := make([]roleResp, 0, len(roles))
	for _, r := range roles {
		resps = append(resps, toRoleResp(&r))
	}
	response.OKWithData(c, gin.H{"items": resps, "total": len(resps)})
}

// CreateRole POST /api/v1/admin/roles
func (h *Handler) CreateRole(c *gin.Context) {
	var req createRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	var count int64
	h.db.WithContext(c.Request.Context()).Model(&po.RolePO{}).
		Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		response.Fail(c, http.StatusConflict, bizErr.ErrDuplicate, "角色名已存在")
		return
	}

	perms := permissionsToJSON(req.Permissions)
	role := po.RolePO{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Permissions: perms,
		IsPreset:    false,
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&role).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "创建角色失败")
		return
	}
	response.OKWithData(c, toRoleResp(&role))
}

// UpdateRole PUT /api/v1/admin/roles/:id
func (h *Handler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req updateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	var role po.RolePO
	if err := h.db.WithContext(c.Request.Context()).Where("id = ?", id).First(&role).Error; err != nil {
		response.Fail(c, http.StatusNotFound, bizErr.ErrNotFound, "角色不存在")
		return
	}
	if role.IsPreset {
		response.Fail(c, http.StatusForbidden, bizErr.ErrForbidden, "预置角色不允许修改")
		return
	}

	h.db.WithContext(c.Request.Context()).Model(&role).
		Update("permissions", permissionsToJSON(req.Permissions))
	response.OK(c)
}

// DeleteRole DELETE /api/v1/admin/roles/:id
func (h *Handler) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	var role po.RolePO
	if err := h.db.WithContext(c.Request.Context()).Where("id = ?", id).First(&role).Error; err != nil {
		response.Fail(c, http.StatusNotFound, bizErr.ErrNotFound, "角色不存在")
		return
	}
	if role.IsPreset {
		response.Fail(c, http.StatusForbidden, bizErr.ErrForbidden, "预置角色不允许删除")
		return
	}

	var userCount int64
	h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		response.Fail(c, http.StatusConflict, bizErr.ErrDuplicate, "该角色下仍有用户，无法删除")
		return
	}

	h.db.WithContext(c.Request.Context()).Delete(&role)
	response.OK(c)
}

// ─── 内部工具函数 ────────────────────────────────────────────────

func permissionsToJSON(perms []string) string {
	if len(perms) == 0 {
		return "[]"
	}
	return `["` + strings.Join(perms, `","`) + `"]`
}

func toRoleResp(r *po.RolePO) roleResp {
	raw := strings.TrimSpace(r.Permissions)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	var perms []string
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"`)
		if p != "" {
			perms = append(perms, p)
		}
	}
	if perms == nil {
		perms = []string{}
	}
	return roleResp{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: perms,
		IsPreset:    r.IsPreset,
		CreatedAt:   r.CreatedAt,
	}
}
