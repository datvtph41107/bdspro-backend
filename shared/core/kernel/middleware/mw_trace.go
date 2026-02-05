package middleware

import (
	"context"

	corectx "core/kernel/ctx"
)

// TraceMiddleware ensures trace id exists in context (generated externally or set here as fallback).
// Prefer HTTP/gRPC interceptor to set trace before pipeline.
// This version leaves generation to caller but ensures key is present.
func TraceMiddleware() Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		if corectx.Trace(ctx) == "" {
			// ensure a trace key exists (empty string ok); you could generate uuid here if desired
			ctx = corectx.WithTrace(ctx, "")
		}
		return next(ctx)
	})
}
