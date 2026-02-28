package middleware

import (
	"context"
	"file/internal/config"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimitMiddleware struct {
	ips          map[string]*RateLimiter
	mu           sync.RWMutex
	rateLimit    rate.Limit
	burst        int
	cleanupEvery time.Duration
	lastCleanup  time.Time
	redisEnabled bool
	redisClient  *redis.Client
}

func NewRateLimitMiddleware(cfg *config.Config) *RateLimitMiddleware {
	rl := &RateLimitMiddleware{
		ips:          make(map[string]*RateLimiter),
		rateLimit:    rate.Limit(cfg.RateLimit.RequestsPerSecond),
		burst:        cfg.RateLimit.Burst,
		cleanupEvery: time.Minute * 5,
		lastCleanup:  time.Now(),
		redisEnabled: cfg.Redis.Enabled,
	}

	if rl.redisEnabled {
		rl.redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	// Start cleanup goroutine
	go rl.cleanupStaleEntries()

	return rl
}

func (m *RateLimitMiddleware) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		// Get endpoint
		endpoint := c.FullPath()

		// Check if endpoint is excluded
		if m.isExcludedEndpoint(endpoint) {
			c.Next()
			return
		}

		// Apply rate limiting
		if !m.allow(ip, endpoint) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":        http.StatusTooManyRequests,
				"message":     "Rate limit exceeded",
				"retry_after": m.retryAfter(ip, endpoint),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *RateLimitMiddleware) allow(ip, endpoint string) bool {
	if m.redisEnabled {
		return m.allowRedis(ip, endpoint)
	}
	return m.allowMemory(ip)
}

func (m *RateLimitMiddleware) allowMemory(ip string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Cleanup stale entries if needed
	if time.Since(m.lastCleanup) > m.cleanupEvery {
		m.cleanupMemory()
		m.lastCleanup = time.Now()
	}

	limiter, exists := m.ips[ip]
	if !exists {
		limiter = &RateLimiter{
			limiter:  rate.NewLimiter(m.rateLimit, m.burst),
			lastSeen: time.Now(),
		}
		m.ips[ip] = limiter
	} else {
		limiter.lastSeen = time.Now()
	}

	return limiter.limiter.Allow()
}

func (m *RateLimitMiddleware) allowRedis(ip, endpoint string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("rate:%s:%s", ip, endpoint)

	// Use Redis sorted sets for rate limiting
	now := time.Now().Unix()
	windowStart := now - 60 // 60 second window

	pipe := m.redisClient.Pipeline()

	// Remove old entries
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count requests in window
	pipe.ZCard(ctx, key)

	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})

	// Set expiry on key
	pipe.Expire(ctx, key, 65*time.Second) // Slightly longer than window

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return true // Fail open
	}

	// Get count
	countCmd := cmds[1].(*redis.IntCmd)
	count, err := countCmd.Result()
	if err != nil {
		return true
	}

	// Check if within limit (e.g., 100 requests per minute)
	return count <= int64(m.burst)
}

func (m *RateLimitMiddleware) retryAfter(ip, endpoint string) int {
	if m.redisEnabled {
		return m.retryAfterRedis(ip, endpoint)
	}
	return m.retryAfterMemory(ip)
}

func (m *RateLimitMiddleware) retryAfterMemory(ip string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limiter, exists := m.ips[ip]
	if !exists {
		return 1
	}

	reservation := limiter.limiter.Reserve()
	if !reservation.OK() {
		return 60 // Default 60 seconds
	}

	delay := reservation.Delay()
	return int(delay.Seconds()) + 1
}

func (m *RateLimitMiddleware) retryAfterRedis(ip, endpoint string) int {
	ctx := context.Background()
	key := fmt.Sprintf("rate:%s:%s", ip, endpoint)

	// Get oldest request timestamp
	result, err := m.redisClient.ZRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil || len(result) == 0 {
		return 1
	}

	oldest := int64(result[0].Score)
	now := time.Now().Unix()
	windowEnd := oldest + 60

	if now >= windowEnd {
		return 1
	}

	return int(windowEnd - now)
}

func (m *RateLimitMiddleware) cleanupMemory() {
	now := time.Now()
	for ip, limiter := range m.ips {
		if now.Sub(limiter.lastSeen) > 10*time.Minute {
			delete(m.ips, ip)
		}
	}
}

func (m *RateLimitMiddleware) cleanupStaleEntries() {
	ticker := time.NewTicker(m.cleanupEvery)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		m.cleanupMemory()
		m.mu.Unlock()
	}
}

func (m *RateLimitMiddleware) isExcludedEndpoint(endpoint string) bool {
	// List of endpoints that don't require rate limiting
	excluded := []string{
		"/health",
		"/metrics",
		"/v1/files/:file_id/download", // Large downloads need different handling
	}

	for _, excl := range excluded {
		if strings.HasPrefix(endpoint, excl) {
			return true
		}
	}

	return false
}

// Per-user rate limiting
func (m *RateLimitMiddleware) LimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		ip := c.ClientIP()
		endpoint := c.FullPath()
		key := fmt.Sprintf("%s:%v", ip, userID)

		if !m.allow(key, endpoint) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Dynamic rate limiting based on endpoint
func (m *RateLimitMiddleware) DynamicLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()
		// ip := c.ClientIP()

		// Apply different limits based on endpoint
		var limit rate.Limit
		var burst int

		switch {
		case strings.Contains(endpoint, "/upload"):
			limit = rate.Limit(5) // 5 requests per second for uploads
			burst = 10
		case strings.Contains(endpoint, "/download"):
			limit = rate.Limit(20) // 20 requests per second for downloads
			burst = 50
		default:
			limit = m.rateLimit
			burst = m.burst
		}

		// Create temporary limiter for this request
		limiter := rate.NewLimiter(limit, burst)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Rate limit exceeded for this endpoint",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
