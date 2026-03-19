package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/euler/mtap/pkg/logger"
)

// LoggerMiddleware 记录请求方法、路径、耗时、状态码
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		cost := time.Since(start)
		status := c.Writer.Status()

		logger.Info("http_request",
			"method", method,
			"path", path,
			"status", status,
			"cost_ms", cost.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
