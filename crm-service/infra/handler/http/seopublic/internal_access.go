package seopublic

import (
	"crypto/hmac"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const GovernanceInternalKeyHeader = "X-QHPro-Internal-Operation-Key"

// RequireGovernanceInternalKey keeps governance and operations off the public
// CRM port. Gateway authenticates the user, then injects this server-only key.
func RequireGovernanceInternalKey(expectedKey string) gin.HandlerFunc {
	expectedKey = strings.TrimSpace(expectedKey)
	return func(c *gin.Context) {
		if expectedKey == "" {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":    http.StatusServiceUnavailable,
				"message": "CRM SEO governance internal key chưa được cấu hình",
			})
			return
		}
		provided := strings.TrimSpace(c.GetHeader(GovernanceInternalKeyHeader))
		if provided == "" || !hmac.Equal([]byte(provided), []byte(expectedKey)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "CRM SEO governance internal key không hợp lệ",
			})
			return
		}
		c.Next()
	}
}
