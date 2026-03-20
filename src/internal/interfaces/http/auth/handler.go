// Package auth 接口层 - 认证鉴权 HTTP Handler
package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/euler/mtap/internal/infrastructure/persistence/po"
	pkgAuth "github.com/euler/mtap/pkg/auth"
	bizErr "github.com/euler/mtap/pkg/errors"
	"github.com/euler/mtap/pkg/response"
)

// Handler 认证 HTTP 处理器
type Handler struct {
	db         *gorm.DB
	jwtManager *pkgAuth.JWTManager
}

// NewHandler 创建认证处理器
func NewHandler(db *gorm.DB, jwtManager *pkgAuth.JWTManager) *Handler {
	return &Handler{db: db, jwtManager: jwtManager}
}

// RegisterRoutes 注册认证路由（无需 JWT 鉴权的公开接口）
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/auth")
	g.POST("/login", h.Login)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
	g.GET("/profile", h.Profile)
}

// ── 请求 / 响应结构体 ────────────────────────────────────────────

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	RealName     string `json:"real_name"`
	Role         string `json:"role"`
	DepartmentID string `json:"department_id,omitempty"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type profileResp struct {
	UserID       string     `json:"user_id"`
	Username     string     `json:"username"`
	RealName     string     `json:"real_name"`
	Role         string     `json:"role"`
	DepartmentID string     `json:"department_id,omitempty"`
	Status       string     `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// ── Handler 实现 ─────────────────────────────────────────────────

// Login POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	var user po.UserPO
	if err := h.db.WithContext(c.Request.Context()).
		Where("username = ? AND status = ?", req.Username, "active").
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "用户名或密码错误")
			return
		}
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "服务器内部错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "用户名或密码错误")
		return
	}

	var role po.RolePO
	h.db.WithContext(c.Request.Context()).Where("id = ?", user.RoleID).First(&role)

	pair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, role.Name, user.DepartmentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "Token 生成失败")
		return
	}

	now := time.Now()
	h.db.WithContext(c.Request.Context()).Model(&po.UserPO{}).
		Where("id = ?", user.ID).Update("last_login_at", now)

	response.OKWithData(c, loginResp{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.ExpiresAt,
		UserID:       user.ID,
		Username:     user.Username,
		RealName:     user.RealName,
		Role:         role.Name,
		DepartmentID: user.DepartmentID,
	})
}

// Refresh POST /api/v1/auth/refresh
func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, bizErr.ErrInvalidParam, err.Error())
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "refresh_token 无效或已过期")
		return
	}

	var user po.UserPO
	if err := h.db.WithContext(c.Request.Context()).
		Where("id = ? AND status = ?", claims.UserID, "active").
		First(&user).Error; err != nil {
		response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "用户不存在或已停用")
		return
	}

	var role po.RolePO
	h.db.WithContext(c.Request.Context()).Where("id = ?", user.RoleID).First(&role)

	pair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, role.Name, user.DepartmentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "Token 生成失败")
		return
	}

	response.OKWithData(c, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"expires_at":    pair.ExpiresAt,
	})
}

// Logout POST /api/v1/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	// 无状态 JWT：客户端清除本地 token 即可
	// 若后续接入 Redis 黑名单，在此处将 jti 写入黑名单
	response.OKWithData(c, gin.H{"message": "登出成功"})
}

// Profile GET /api/v1/auth/profile
func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "未认证")
		return
	}

	var user po.UserPO
	if err := h.db.WithContext(c.Request.Context()).
		Where("id = ?", userID).First(&user).Error; err != nil {
		response.Fail(c, http.StatusNotFound, bizErr.ErrNotFound, "用户不存在")
		return
	}

	var role po.RolePO
	h.db.WithContext(c.Request.Context()).Where("id = ?", user.RoleID).First(&role)

	response.OKWithData(c, profileResp{
		UserID:       user.ID,
		Username:     user.Username,
		RealName:     user.RealName,
		Role:         role.Name,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		LastLoginAt:  user.LastLoginAt,
	})
}

// SeedAdminUser 初始化预置角色和默认管理员账号（开发/首次部署时调用）
func SeedAdminUser(db *gorm.DB) error {
	const adminRoleID = "00000000-0000-0000-0000-000000000001"

	var roleCount int64
	db.Model(&po.RolePO{}).Where("id = ?", adminRoleID).Count(&roleCount)
	if roleCount == 0 {
		presets := []po.RolePO{
			{ID: adminRoleID, Name: "admin", Permissions: `["*"]`, IsPreset: true},
			{ID: "00000000-0000-0000-0000-000000000002", Name: "scheduler_admin", Permissions: `["rule:*","resource:*","appt:*"]`, IsPreset: true},
			{ID: "00000000-0000-0000-0000-000000000003", Name: "operator", Permissions: `["appt:*","triage:*"]`, IsPreset: true},
			{ID: "00000000-0000-0000-0000-000000000004", Name: "nurse", Permissions: `["triage:*"]`, IsPreset: true},
			{ID: "00000000-0000-0000-0000-000000000005", Name: "viewer", Permissions: `["*:read"]`, IsPreset: true},
		}
		for i := range presets {
			db.Create(&presets[i])
		}
	}

	var userCount int64
	db.Model(&po.UserPO{}).Where("username = ?", "admin").Count(&userCount)
	if userCount == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("Admin@1234"), bcrypt.DefaultCost)
		db.Create(&po.UserPO{
			ID:           uuid.New().String(),
			Username:     "admin",
			PasswordHash: string(hash),
			RealName:     "系统管理员",
			RoleID:       adminRoleID,
			Status:       "active",
		})
	}
	return nil
}
