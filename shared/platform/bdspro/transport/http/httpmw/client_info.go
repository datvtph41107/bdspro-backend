// platform/transport/http/middleware/client_info.go
package httpmw

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	pctx "bdspro/core/context"
)

// ClientInfo extract và inject client IP và User-Agent vào context.
//
// Tại sao middleware riêng thay vì trong Logging?
// - IP và UserAgent được dùng bởi nhiều layer: Rate Limit (Layer 2), Auth (Layer 3), Logging (Layer 0)
// - Single extraction point — không duplicate logic
// - Logging middleware chỉ cần đọc từ context, không cần biết cách extract
func ClientInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ip := extractClientIP(c.Request)
		ua := c.Request.UserAgent()

		ctx = pctx.WithClientIP(ctx, ip)
		ctx = pctx.WithUserAgent(ctx, ua)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// extractClientIP trả về real client IP.
//
// Thứ tự ưu tiên:
// 1. X-Forwarded-For header (khi qua nhiều proxy)
// 2. X-Real-IP header (khi qua một proxy/load balancer)
// 3. RemoteAddr (kết nối trực tiếp)
//
// Tại sao X-Forwarded-For lấy phần đầu tiên?
// Format: "client, proxy1, proxy2" — client IP luôn đứng đầu.
// Proxy cuối (gateway của ta) add IP của mình vào cuối.
//
// Security note: X-Forwarded-For có thể bị spoof nếu không có trusted proxy.
// Layer 2 sẽ handle IP reputation check — đây chỉ extract, không validate.
func extractClientIP(r *http.Request) string {
	// X-Forwarded-For: "1.2.3.4, 5.6.7.8, 9.10.11.12"
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Lấy IP đầu tiên — đó là client IP
		if idx := strings.IndexByte(xff, ','); idx >= 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// X-Real-IP: "1.2.3.4" — set bởi nginx
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// RemoteAddr: "1.2.3.4:56789" — tách port
	addr := r.RemoteAddr
	if idx := strings.LastIndexByte(addr, ':'); idx >= 0 {
		return addr[:idx]
	}
	return addr
}
