// platform/transport/http/middleware/logging.go
package httpmw

import (
	observability "bdspro/insight"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging log request start và end với full context fields.
//
// Tại sao chỉ log ở boundary (F3)?
// - Mỗi request chỉ có 2 log entry: "request started" và "request completed"
// - Handler/service KHÔNG tự log request info — tránh duplicate
// - Handler chỉ log business-specific events: "user created", "payment processed"
//
// Tại sao Logging là middleware CUỐI CÙNG trong chain (Layer 0)?
// - Phải sau RequestID, TraceID, ClientInfo để context đã có đủ fields
// - Phải trước Layer 1-6 để wrap toàn bộ processing time
//
// duration_ms là int64 milliseconds (F5):
// - Dễ filter: duration_ms > 500 thay vì duration > 0.5
// - Consistent với industry standard (Datadog, Prometheus đều dùng ms)
func Logging(pipe *observability.Pipeline) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		// LogRequestStart: log + IncActiveRequests + start span
		// Return enriched ctx (có cached logger) và span cleanup
		ctx, spanEnd := pipe.LogRequestStart(ctx, c.Request.Method, c.Request.URL.Path)
		defer spanEnd()

		// Update request context với enriched ctx (có cached logger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// LogRequestEnd: log + ObserveDuration + DecActiveRequests
		pipe.LogRequestEnd(ctx, c.Writer.Status(), time.Since(start))
	}
}
