package middleware

import (
	"context"
	core_ctx "core/kernel/ctx"
	"fmt"
)

func RecoverMiddleware() Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) (err error) {
		defer func() {
			if r := recover(); r != nil {
				core_ctx.Logger(ctx).Error(
					"PANIC",
					"panic recovered",
					fmt.Errorf("%v", r),
				)
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		return next(ctx)
	})
}
