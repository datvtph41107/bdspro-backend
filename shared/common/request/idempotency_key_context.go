package request

import (
	"context"
	"errors"
)

type idempotencyKeyContextKey struct{}

var ErrIdempotencyKeyContextConflict = errors.New(
	"idempotency key context conflict",
)

/**
 * BindIdempotencyKey binds one immutable command identity.
 * Rebinding the same key is idempotent; a different key is a conflict.
 */
func BindIdempotencyKey(
	ctx context.Context,
	key string,
) (
	context.Context,
	error,
) {
	if ctx == nil {
		panic("request: nil context")
	}

	if !IsValidIdempotencyKey(key) {
		return ctx, ErrInvalidIdempotencyKey
	}

	if existing, ok := IdempotencyKeyFromContext(ctx); ok {
		if existing == key {
			return ctx, nil
		}

		return ctx, ErrIdempotencyKeyContextConflict
	}

	return context.WithValue(
		ctx,
		idempotencyKeyContextKey{},
		key,
	), nil
}

/**
 * IdempotencyKeyFromContext reads the canonical optional command identity.
 */
func IdempotencyKeyFromContext(
	ctx context.Context,
) (
	string,
	bool,
) {
	if ctx == nil {
		return "", false
	}

	key, ok := ctx.Value(idempotencyKeyContextKey{}).(string)
	if !ok || !IsValidIdempotencyKey(key) {
		return "", false
	}

	return key, true
}
