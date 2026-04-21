// platform/transport/http/middleware/chain.go
//
// Chain định nghĩa thứ tự middleware chuẩn cho Layer 0.
// Tập trung ở một file để dễ review và đảm bảo thứ tự không bị thay đổi nhầm.
package httpmw

import (
	observability "bdspro/insight"

	"github.com/gin-gonic/gin"
)

// Layer0Chain trả về middleware chain cho Layer 0 (Observability).
// Thứ tự CỐ ĐỊNH và có ý nghĩa — không thay đổi tùy tiện.
//
// Recovery  → phải đứng đầu, bắt panic từ mọi middleware sau
// RequestID → phải trước Logging (Logging cần request_id)
// TraceID   → phải trước Logging (Logging cần trace_id)
// ClientInfo → phải trước Logging (Logging cần client_ip)
// Logging   → phải cuối Layer 0, wrap toàn bộ processing
func Layer0Chain(pipe *observability.Pipeline) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Recovery(pipe),
		RequestID(),
		TraceID(),
		ClientInfo(),
		Logging(pipe),
	}
}
