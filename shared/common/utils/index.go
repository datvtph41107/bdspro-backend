package _utils

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuditIP lấy địa chỉ IP của client từ request trong Gin
// AuditIPFromContext lấy địa chỉ IP của client từ context.Context
func AuditIPFromContext(ctx context.Context) string {
	// Nếu context chứa giá trị "X-Forwarded-For" hoặc "client-ip" thì lấy ra
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
