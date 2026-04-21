// platform/transport/http/middleware/recovery.go
package httpmw

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	pctx "bdspro/core/context"
	observability "bdspro/insight"
)

// Recovery bắt panic và trả về 500 JSON response nhất quán.
//
// Tại sao không dùng gin.Default() built-in recovery?
// - Built-in recovery log plain text, không có request_id, không structured
// - Built-in recovery dùng fmt.Fprintf không đúng format ErrorEnvelope
// - Custom recovery cho phép: 1) structured log, 2) metrics.IncPanic(), 3) request_id trong response
//
// Tại sao phải là middleware ĐẦU TIÊN?
// - Nếu RequestID middleware panic → không có recovery → goroutine die → process crash
// - Recovery bao ngoài tất cả → bắt được panic từ bất kỳ middleware nào sau đó
//
// Tại sao write raw JSON string thay vì gin.H?
// - gin.H = map[string]interface{} → allocation + reflection
// - Raw string: zero allocation, compile-time correct
// - Format nhất quán với ErrorEnvelope trong response package
func Recovery(pipe *observability.Pipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return // no panic — happy path
			}

			// Capture stack TRƯỚC khi làm gì khác — defer order matters
			stack := debug.Stack()

			ctx := c.Request.Context()
			pipe.LogPanic(ctx, r, stack)

			// Abort để các middleware/handler sau không chạy tiếp
			c.Abort()

			// Ghi response — nếu header đã được ghi thì không ghi nữa
			if c.Writer.Written() {
				return
			}

			// Lấy request_id để đưa vào response — client có thể debug
			rid := pctx.RequestID(ctx)
			writeInternalError(c.Writer, rid)
		}()

		c.Next()
	}
}

// writeInternalError ghi 500 response theo ErrorEnvelope format.
// Dùng raw JSON string — zero allocation, không reflect, không gin.H.
//
// Tại sao fmt.Sprintf thay vì json.Marshal?
// - json.Marshal cần allocation + reflection
// - Format này KHÔNG BAO GIỜ thay đổi (hardcoded error response)
// - %q escape request_id đúng cách — tránh JSON injection nếu rid có ký tự đặc biệt
func writeInternalError(w http.ResponseWriter, requestID string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	// %q format: escape special chars trong requestID — prevent JSON injection
	body := fmt.Sprintf(
		`{"error":{"code":"INTERNAL","message":"internal server error","request_id":%q}}`,
		requestID,
	)
	_, _ = w.Write([]byte(body))
}
