package _utils

import (
	_enum "common/domain/enum"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuditIP lấy địa chỉ IP của client từ request trong Gin
// AuditIPFromContext lấy địa chỉ IP của client từ context.Context
func AuditIPFromContext(ctx context.Context) string {
	// Lấy từ context key enum
	if clientIP, ok := ctx.Value(_enum.ClientIPKey).(string); ok && clientIP != "" {
		return clientIP
	}
	// Fallback cho compatibility với code cũ
	if xfHeader, ok := ctx.Value("X-Forwarded-For").(string); ok && xfHeader != "" {
		clientIP := strings.Split(xfHeader, ",")[0]
		return strings.TrimSpace(clientIP)
	}
	if clientIP, ok := ctx.Value("client-ip").(string); ok && clientIP != "" {
		return clientIP
	}
	// Nếu không có, trả về chuỗi rỗng
	return ""
}

// GetUserAgentFromContext lấy User-Agent từ context.Context
func GetUserAgentFromContext(ctx context.Context) string {
	// Lấy từ context key enum
	if userAgent, ok := ctx.Value(_enum.UserAgentKey).(string); ok && userAgent != "" {
		return userAgent
	}
	// Fallback cho compatibility với code cũ
	if userAgent, ok := ctx.Value("User-Agent").(string); ok && userAgent != "" {
		return userAgent
	}
	// Nếu không có, trả về chuỗi rỗng
	return ""
}

func AuditIP(c *gin.Context) string {
	xfHeader := c.GetHeader("X-Forwarded-For")
	if xfHeader == "" {
		return c.ClientIP()
	}
	// X-Forwarded-For có thể chứa nhiều IP, lấy IP đầu tiên
	clientIP := strings.Split(xfHeader, ",")[0]
	return strings.TrimSpace(clientIP)
}

// TrimFirstWords cắt n từ đầu tiên của chuỗi
func TrimFirstWords(text string, n int) string {
	words := strings.Fields(text) // Tách thành mảng các từ, tự động bỏ khoảng trắng thừa

	if len(words) <= n {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:n], " ")
}
