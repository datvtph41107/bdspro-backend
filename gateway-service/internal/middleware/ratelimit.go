package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.RWMutex
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(limit), limit),
	}
}

func (rl *RateLimiter) UpdateLimit(limit int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limiter = rate.NewLimiter(rate.Limit(limit), limit)
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.limiter.Allow()
}

func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
