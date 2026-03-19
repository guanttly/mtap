package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "github.com/euler/mtap/pkg/errors"
	"github.com/euler/mtap/pkg/response"
)

type rateCounter struct {
	windowStart time.Time
	count       int
}

// RateLimitMiddleware 固定窗口限流（默认：60 次/分钟/用户）
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	if maxRequests <= 0 {
		maxRequests = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	var (
		mu       sync.Mutex
		counters = make(map[string]*rateCounter)
	)

	return func(c *gin.Context) {
		key := c.GetString("user_id")
		if key == "" {
			key = c.ClientIP()
		}

		now := time.Now()
		mu.Lock()
		rc, ok := counters[key]
		if !ok {
			rc = &rateCounter{windowStart: now, count: 0}
			counters[key] = rc
		}
		if now.Sub(rc.windowStart) >= window {
			rc.windowStart = now
			rc.count = 0
		}
		rc.count++
		over := rc.count > maxRequests
		mu.Unlock()

		if over {
			response.Fail(c, http.StatusTooManyRequests, bizErr.ErrRateLimit, "请求过于频繁")
			c.Abort()
			return
		}

		c.Next()
	}
}
