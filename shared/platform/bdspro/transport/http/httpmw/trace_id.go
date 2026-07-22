// platform/transport/http/middleware/trace_id.go
package httpmw

import (
	"github.com/gin-gonic/gin"

	pctx "bdspro/core/context"
)

// TraceID đọc X-Trace-ID từ upstream (load balancer, API gateway) hoặc generate mới.
//
// Sự khác biệt với RequestID:
// - Request-ID: unique per HTTP request (kể cả retry — mỗi retry là request mới)
// - Trace-ID: xuyên suốt TOÀN BỘ operation kể cả retry, fan-out → dùng cho distributed tracing
//
// Pattern: Client tạo Trace-ID một lần, pass qua mọi microservice.
// Gateway propagate sang downstream qua gRPC metadata "trace-id".
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Trace-ID")
		if id == "" {
			// Gateway sinh trace-id nếu không có — đảm bảo tracing bắt đầu từ gateway
			id = generateID()
		}

		ctx := pctx.WithTraceID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", id)

		c.Next()
	}
}
