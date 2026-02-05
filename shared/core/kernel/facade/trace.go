// core/facade/trace.go
package facade

import "context"

type traceKey struct{}

func WithTrace(ctx context.Context, tid string) context.Context {
	return context.WithValue(ctx, traceKey{}, tid)
}

func TraceID(ctx context.Context) string {
	if v := ctx.Value(traceKey{}); v != nil {
		return v.(string)
	}
	return ""
}
