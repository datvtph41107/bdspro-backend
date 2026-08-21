package requestcontext

import (
	canonical "common/request"
	"context"
)

var ErrIdempotencyKeyContextConflict = canonical.ErrIdempotencyKeyContextConflict

func BindIdempotencyKey(ctx context.Context, key string) (context.Context, error) {
	return canonical.BindIdempotencyKey(ctx, key)
}

func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	return canonical.IdempotencyKeyFromContext(ctx)
}
