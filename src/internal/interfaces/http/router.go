// Package http 接口层 - 路由注册总入口
package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/euler/mtap/internal/interfaces/http/admin"
	"github.com/euler/mtap/internal/interfaces/http/analytics"
	"github.com/euler/mtap/internal/interfaces/http/appointment"
	"github.com/euler/mtap/internal/interfaces/http/auth"
	"github.com/euler/mtap/internal/interfaces/http/middleware"
	"github.com/euler/mtap/internal/interfaces/http/optimization"
	"github.com/euler/mtap/internal/interfaces/http/resource"
	"github.com/euler/mtap/internal/interfaces/http/rule"
	"github.com/euler/mtap/internal/interfaces/http/triage"
	pkgAuth "github.com/euler/mtap/pkg/auth"
)

// Deps Router 依赖集合（由 cmd/api-server 进行注入）
type Deps struct {
	JWTManager *pkgAuth.JWTManager

	AuthHandler         *auth.Handler
	AdminHandler        *admin.Handler
	RuleHandler         *rule.Handler
	ResourceHandler     *resource.Handler
	AppointmentHandler  *appointment.Handler
	TriageHandler       *triage.Handler
	AnalyticsHandler    *analytics.Handler
	OptimizationHandler *optimization.Handler
}

// NewEngine 创建 gin.Engine（挂载全局中间件）
func NewEngine(deps Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware(60, time.Minute))

	// API v1 公开路由（auth 登录/刷新）
	v1 := r.Group("/api/v1")
	if deps.AuthHandler != nil {
		deps.AuthHandler.RegisterRoutes(v1)
	}

	// 需要鉴权的路由
	if deps.JWTManager != nil {
		v1.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
			JWTManager: deps.JWTManager,
			PublicPaths: map[string]struct{}{
				"/api/v1/auth/login":   {},
				"/api/v1/auth/refresh": {},
				"/api/v1/auth/logout":  {},
			},
		}))
	}

	// 注册业务路由
	if deps.AdminHandler != nil {
		deps.AdminHandler.RegisterRoutes(v1)
	}
	if deps.RuleHandler != nil {
		deps.RuleHandler.RegisterRoutes(v1)
	}
	if deps.ResourceHandler != nil {
		deps.ResourceHandler.RegisterRoutes(v1)
	}
	if deps.AppointmentHandler != nil {
		deps.AppointmentHandler.RegisterRoutes(v1)
	}
	if deps.TriageHandler != nil {
		deps.TriageHandler.RegisterRoutes(v1)
	}
	if deps.AnalyticsHandler != nil {
		deps.AnalyticsHandler.RegisterRoutes(v1)
	}
	if deps.OptimizationHandler != nil {
		deps.OptimizationHandler.RegisterRoutes(v1)
	}

	return r
}
