package requestcontext

import (
	canonical "common/request"
	"context"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return canonical.WithRequestID(ctx, requestID)
}

func FromContext(ctx context.Context) (string, bool) {
	return canonical.RequestIDFromContext(ctx)
}
