// platform/core/context/keys.go
//
// Package context cung cấp type-safe context key/value cho toàn bộ platform.
// Thiết kế: unexported contextKey type → không package nào ngoài đây tạo được key trùng.
// Đây là nguồn sự thật DUY NHẤT cho context propagation — mọi service đều dùng.
package context

// contextKey là private type để tránh key collision với external packages.
// Dùng string làm underlying type để fmt.Stringer readable trong debug.
//
// Tại sao không dùng int/iota?
// - String key dễ đọc trong log/panic output: "platform.context.request_id"
// - iota đòi hỏi caller biết thứ tự — fragile khi thêm key mới
type contextKey string

// String implement fmt.Stringer — giúp debug output dễ đọc.
func (c contextKey) String() string {
	return "platform.context." + string(c)
}

// Các key chuẩn của platform — không thêm tùy tiện.
// Mỗi key có một getter/setter tương ứng trong values.go.
// Quy tắc đặt tên: snake_case, match với log field name để nhất quán.
const (
	keyRequestID contextKey = "request_id" // ID duy nhất per-request, gắn vào response header
	keyTraceID   contextKey = "trace_id"   // Distributed trace ID, xuyên suốt service chain
	keyUserID    contextKey = "user_id"    // Authenticated user ID (uint64)
	keyClientIP  contextKey = "client_ip"  // Real client IP sau proxy
	keyUserAgent contextKey = "user_agent" // HTTP User-Agent header
	keyRole      contextKey = "role"       // User role string ("admin", "user", ...)
	keyPlanID    contextKey = "plan_id"    // Subscription plan ID
	keyLogger    contextKey = "logger"     // Cached *zap.Logger (internal, không expose ra ngoài)

	// L1 ochestrator config features
	keyRouteService      contextKey = "route_service"
	keyRouteTimeout      contextKey = "route_timeout"
	keyRouteRequiresAuth contextKey = "route_requires_auth"
	keyRouteRequiredRole contextKey = "route_required_role"
	keyServiceInstances  contextKey = "service_instances"
)
