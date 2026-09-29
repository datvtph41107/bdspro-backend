package operation

import (
	"context"
	"errors"
)

var ErrContextConflict = errors.New(
	"operation conflicts with context",
)

type contextKey struct{}

// Bind stores the current request operation without allowing another operation to overwrite it.
func Bind(
	ctx context.Context,
	code Code,
) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	if !code.IsValid() {
		return ctx, ErrInvalidCode
	}

	if current, ok := FromContext(ctx); ok {
		if current == code {
			return ctx, nil
		}

		return ctx, ErrContextConflict
	}

	return context.WithValue(
		ctx,
		contextKey{},
		code,
	), nil
}

// FromContext returns the operation selected for the current request.
func FromContext(
	ctx context.Context,
) (Code, bool) {
	if ctx == nil {
		return "", false
	}

	code, ok := ctx.Value(
		contextKey{},
	).(Code)

	if !ok || !code.IsValid() {
		return "", false
	}

	return code, true
}
