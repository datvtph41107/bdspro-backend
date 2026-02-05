package middleware

import "context"

// Middleware signature matches Laravel pipeline idea
type Middleware interface {
	Handle(ctx context.Context, payload interface{}, next func(ctx context.Context) error) error
}

// MiddlewareFunc helper
type MiddlewareFunc func(
	ctx context.Context,
	payload interface{},
	next func(ctx context.Context) error,
) error

func (f MiddlewareFunc) Handle(ctx context.Context,
	payload interface{},
	next func(ctx context.Context) error,
) error {
	return f(ctx, payload, next)
}
