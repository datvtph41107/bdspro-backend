package httpmw

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Lang", "X-Request-ID", "X-Trace-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID", "X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// CORS returns a CORS middleware
func CORS(cfg CORSConfig) gin.HandlerFunc {
	// Pre-validate allowed origins for performance
	allowAll := false
	allowedOrigins := make(map[string]bool)
	for _, origin := range cfg.AllowOrigins {
		if origin == "*" {
			allowAll = true
		} else {
			allowedOrigins[origin] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Handle preflight (OPTIONS) request
		if c.Request.Method == http.MethodOptions {
			if allowAll || allowedOrigins[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods))
				c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders))
				c.Header("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders))
				if cfg.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
				if cfg.MaxAge > 0 {
					c.Header("Access-Control-Max-Age", formatDuration(cfg.MaxAge))
				}
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Handle actual request
		if allowAll || allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			if cfg.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if len(cfg.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders))
			}
			c.Next()
			return
		}

		// Origin not allowed
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func formatDuration(d time.Duration) string {
	return string(strconv.FormatInt(int64(d.Seconds()), 10))
}
