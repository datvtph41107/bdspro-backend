package _utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuditIP lấy địa chỉ IP của client từ request trong Gin
func AuditIP(c *gin.Context) string {
	xfHeader := c.GetHeader("X-Forwarded-For")
	if xfHeader == "" {
		return c.ClientIP()
	}
	// X-Forwarded-For có thể chứa nhiều IP, lấy IP đầu tiên
	clientIP := strings.Split(xfHeader, ",")[0]
	return strings.TrimSpace(clientIP)
}
