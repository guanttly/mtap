// Package http 接口层 - 路由注册总入口
package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/internal/interfaces/http/resource"
	"github.com/euler/mtap/internal/interfaces/http/rule"
	"github.com/euler/mtap/pkg/auth"
)

// Deps Router 依赖集合（由 cmd/api-server 进行注入）
type Deps struct {
	JWTManager *auth.JWTManager

	RuleHandler *rule.Handler
	ResourceHandler *resource.Handler
	// 后续模块（resource/appointment/triage/...）按需补齐
}

// NewEngine 创建 gin.Engine（挂载全局中间件）
func NewEngine(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	// 认证后才能按 user_id 限流（公开接口按 ip 限流）
	r.Use(middleware.RateLimitMiddleware(60, time.Minute))

	// API v1
	v1 := r.Group("/api/v1")
	if deps.JWTManager != nil {
		v1.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
			JWTManager: deps.JWTManager,
			PublicPaths: map[string]struct{}{
				// 认证模块完成后在这里放开 /auth/login 等
			},
		}))
	}

	// 注册业务路由
	if deps.RuleHandler != nil {
		deps.RuleHandler.RegisterRoutes(v1)
	}
	if deps.ResourceHandler != nil {
		deps.ResourceHandler.RegisterRoutes(v1)
	}

	return r
}
