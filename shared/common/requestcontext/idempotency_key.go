package requestcontext

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	IdempotencyKeyHeader      = "Idempotency-Key"
	IdempotencyKeyMetadataKey = "idempotency-key"
	MaxIdempotencyKeyLength   = maxContextIDLength
)

var ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")

const idempotencyKeyContextKey contextKey = "idempotency-key"

/**
 * IdempotencyKeyFromRequest validates an optional client key.
 */
func IdempotencyKeyFromRequest(r *http.Request) (string, bool, error) {
	if r == nil || r.Header == nil {
		return "", false, nil
	}

	values := r.Header.Values(IdempotencyKeyHeader)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 {
		return "", false, ErrInvalidIdempotencyKey
	}

	value := strings.TrimSpace(values[0])
	if !IsValidIdempotencyKey(value) {
		return "", false, ErrInvalidIdempotencyKey
	}

	r.Header.Set(IdempotencyKeyHeader, value)
	return value, true, nil
}

/**
 * WithIdempotencyKey stores a validated key in service context.
 */
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if !IsValidIdempotencyKey(key) {
		return ctx
	}

	return context.WithValue(ctx, idempotencyKeyContextKey, key)
}

/**
 * IdempotencyKeyFromContext reads the canonical service key.
 */
func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	key, ok := ctx.Value(idempotencyKeyContextKey).(string)
	return key, ok && IsValidIdempotencyKey(key)
}

/**
 * IsValidIdempotencyKey validates a client-owned key.
 */
func IsValidIdempotencyKey(value string) bool {
	return isValidContextID(value, MaxIdempotencyKeyLength)
}

/**
 * NewIdempotencyKey creates a key for an owned command.
 */
func NewIdempotencyKey() string {
	return newContextID("idem")
}
