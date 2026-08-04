// platform/transport/http/middleware/request_id.go
package httpmw

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"

	pctx "bdspro/core/context"
)

// RequestID đọc X-Request-ID từ header hoặc generate mới.
//
// Tại sao cho phép client set request_id?
// - Client SDK có thể tạo ID trước để correlate request với response
// - Request correlation: mỗi retry dùng request_id riêng
//
// Tại sao sinh ID mới khi không có?
// - Log không có request_id = không trace được request → useless log
// - Luôn có ID = luôn có thể grep log theo request
//
// Tại sao hex(16 bytes) = 32 char?
// - 128 bit entropy → collision probability astronomically low
// - Shorter than UUID, no hyphens, URL-safe, case-insensitive
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = generateID()
		}

		// Inject vào request context — tất cả middleware/handler sau đọc từ đây
		ctx := pctx.WithRequestID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)

		// Echo lại trong response header — client dùng để correlate
		c.Header("X-Request-ID", id)

		c.Next()
	}
}

// generateID tạo cryptographically random 16-byte hex string.
// Fallback về timestamp-based ID nếu rand.Read fail (rất hiếm — kernel entropy pool).
func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback: timestamp-based ID — không hoàn hảo nhưng vẫn unique trong ngắn hạn
		return "fb-" + time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b)
}
