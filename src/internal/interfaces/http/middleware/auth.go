// Package middleware 接口层 - 认证与授权中间件
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/euler/mtap/pkg/auth"
	bizErr "github.com/euler/mtap/pkg/errors"
	"github.com/euler/mtap/pkg/response"
)

type (
	// AuthMiddlewareConfig JWT 鉴权中间件配置
	AuthMiddlewareConfig struct {
		JWTManager *auth.JWTManager
		// PublicPaths 以 /api/v1 开头的公开路径（完全匹配）
		PublicPaths map[string]struct{}
	}
)

// AuthMiddleware 验证 JWT，并把 user_id/role/department_id 注入 gin.Context
func AuthMiddleware(cfg AuthMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.PublicPaths != nil {
			if _, ok := cfg.PublicPaths[c.FullPath()]; ok {
				c.Next()
				return
			}
		}

		if cfg.JWTManager == nil {
			response.Fail(c, http.StatusInternalServerError, bizErr.ErrInternal, "JWTManager未初始化")
			c.Abort()
			return
		}

		authz := c.GetHeader("Authorization")
		tokenStr := strings.TrimSpace(authz)
		if tokenStr == "" {
			response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "缺少Authorization Header")
			c.Abort()
			return
		}
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			tokenStr = strings.TrimSpace(tokenStr[7:])
		}
		if tokenStr == "" {
			response.Fail(c, http.StatusUnauthorized, bizErr.ErrUnauthorized, "token为空")
			c.Abort()
			return
		}

		claims, err := cfg.JWTManager.ValidateToken(tokenStr)
		if err != nil {
			response.FailWithError(c, err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("department_id", claims.DepartmentID)
		c.Next()
	}
}

// RequireRoles RBAC：限制允许访问的角色集合
func RequireRoles(allowed ...string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			response.Fail(c, http.StatusForbidden, bizErr.ErrForbidden, "缺少角色信息")
			c.Abort()
			return
		}
		if _, ok := allowedSet[role]; !ok {
			response.Fail(c, http.StatusForbidden, bizErr.ErrForbidden, "无权限访问")
			c.Abort()
			return
		}
		c.Next()
	}
}
