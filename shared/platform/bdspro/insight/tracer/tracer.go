// platform/observability/tracer/tracer.go
package tracer

import "context"

// Tracer interface cho distributed tracing.
// NoopTracer là default — swap với OpenTelemetry khi cần.
type Tracer interface {
	// Start bắt đầu một span mới. Trả về context mới và cleanup function.
	// Caller phải gọi cleanup() khi xong: defer cleanup()
	Start(ctx context.Context, spanName string) (context.Context, func())
}

var _ Tracer = (*NoopTracer)(nil)

// NoopTracer không làm gì — zero overhead.
type NoopTracer struct{}

func (t *NoopTracer) Start(ctx context.Context, _ string) (context.Context, func()) {
	return ctx, func() {}
}
