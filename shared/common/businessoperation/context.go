package businessoperation

import (
	"context"
	"errors"
)

var (
	ErrInvalidOperation = errors.New("business operation is invalid")
	ErrContextConflict  = errors.New("business operation conflicts with context")
)

type contextKey struct{}

/**
 * Bind gắn business operation vào context.
 *
 * Không silent overwrite nếu context đã có một operation khác.
 */
func Bind(ctx context.Context, operation Code) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !operation.IsValid() {
		return ctx, ErrInvalidOperation
	}

	if current, ok := FromContext(ctx); ok {
		if current == operation {
			return ctx, nil
		}
		return ctx, ErrContextConflict
	}

	return context.WithValue(ctx, contextKey{}, operation), nil
}

/**
 * FromContext lấy business operation đã được chọn cho request hiện tại.
 */
func FromContext(ctx context.Context) (Code, bool) {
	if ctx == nil {
		return "", false
	}

	operation, ok := ctx.Value(contextKey{}).(Code)
	if !ok || !operation.IsValid() {
		return "", false
	}

	return operation, true
}
