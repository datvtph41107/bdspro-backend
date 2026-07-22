// platform/core/context/values.go
//
// Getter/setter type-safe cho từng context key.
//
// Pattern nhất quán:
//
//	With<Field>(ctx, value) context.Context  — setter
//	<Field>(ctx) T                           — getter, zero value nếu không có
//	Must<Field>(ctx) T                       — getter + panic nếu không có (dùng cho required field)
//
// Tại sao không dùng generic WithValue[T]?
// - Go generics không support interface constraint với context.WithValue
// - Từng field có validation riêng (ví dụ MustUserID check != 0)
// - Explicit function name → caller biết chính xác field nào đang set/get
package context

import (
	"context"
)

// ─── Request ID ──────────────────────────────────────────────────────────────

// WithRequestID gắn request ID vào context.
// Được gọi bởi RequestID middleware — một lần duy nhất per request.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

// RequestID đọc request ID từ context.
// Trả về "" nếu không có — không panic, vì health check có thể không có ID.
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestID).(string)
	return v
}

// ─── Trace ID ────────────────────────────────────────────────────────────────

// WithTraceID gắn trace ID vào context — dùng cho distributed tracing.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}

// TraceID đọc trace ID.
func TraceID(ctx context.Context) string {
	v, _ := ctx.Value(keyTraceID).(string)
	return v
}

// ─── User ID ─────────────────────────────────────────────────────────────────

// WithUserID gắn authenticated user ID.
// Chỉ được gọi sau khi JWT validated thành công (Layer 3).
func WithUserID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

// UserID đọc user ID. Trả về 0 nếu unauthenticated.
func UserID(ctx context.Context) uint64 {
	v, _ := ctx.Value(keyUserID).(uint64)
	return v
}

// MustUserID đọc user ID và panic nếu == 0.
// Dùng trong handler đã được bảo vệ bởi auth middleware — nếu đến đây mà không có ID thì là bug.
// Panic message có context để debug nhanh.
func MustUserID(ctx context.Context) uint64 {
	id := UserID(ctx)
	if id == 0 {
		panic("platform/context: user_id required but not set in context — missing auth middleware?")
	}
	return id
}

// ─── Client IP ───────────────────────────────────────────────────────────────

// WithClientIP gắn real client IP (đã extract từ X-Forwarded-For hoặc RemoteAddr).
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, keyClientIP, ip)
}

// ClientIP đọc client IP.
func ClientIP(ctx context.Context) string {
	v, _ := ctx.Value(keyClientIP).(string)
	return v
}

// ─── User Agent ──────────────────────────────────────────────────────────────

// WithUserAgent gắn HTTP User-Agent header.
func WithUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, keyUserAgent, ua)
}

// UserAgent đọc User-Agent.
func UserAgent(ctx context.Context) string {
	v, _ := ctx.Value(keyUserAgent).(string)
	return v
}

// ─── Role ────────────────────────────────────────────────────────────────────

// WithRole gắn user role.
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, keyRole, role)
}

// Role đọc user role. Trả về "" nếu unauthenticated.
func Role(ctx context.Context) string {
	v, _ := ctx.Value(keyRole).(string)
	return v
}

// ─── Plan ID ─────────────────────────────────────────────────────────────────

// WithPlanID gắn subscription plan ID.
func WithPlanID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, keyPlanID, id)
}

// PlanID đọc plan ID. Trả về 0 nếu không có plan.
func PlanID(ctx context.Context) uint64 {
	v, _ := ctx.Value(keyPlanID).(uint64)
	return v
}
