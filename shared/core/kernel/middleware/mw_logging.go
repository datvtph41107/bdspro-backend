package middleware

import (
	"context"
	"fmt"
	"time"

	corectx "core/kernel/ctx"
)

// LoggingMiddleware uses logger from ctx (ctx.Logger). If not present, uses Nop.
// It logs start/end and error with structured fields.
func LoggingMiddleware() Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		start := time.Now()
		l := corectx.Logger(ctx).WithTag("bus")
		// attach trace if available
		if t := corectx.Trace(ctx); t != "" {
			l = l.WithTrace(t)
		}

		err := next(ctx)
		elapsed := time.Since(start)

		if err != nil {
			l.Error("BUS_HANDLE_ERROR", "bus handler error",
				err,
				"payload", fmt.Sprintf("%T", payload),
				"elapsed", elapsed,
			)
		} else {
			l.Info("BUS_HANDLED", "bus handler OK",
				"payload", fmt.Sprintf("%T", payload),
				"elapsed", elapsed,
			)
		}
		return err
	})
}
