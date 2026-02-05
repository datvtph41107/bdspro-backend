package middleware

import (
	"context"
	"core/kernel/infra/redis"
	"fmt"
	"time"
)

type HasIdempotency interface {
	IdempotencyKey() string
}

func IdempotencyMiddleware(rdb *redis.Client, serviceName string) Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		idp, ok := payload.(HasIdempotency)
		if !ok {
			return next(ctx)
		}
		key := fmt.Sprintf("idp:%s:%s", serviceName, idp.IdempotencyKey())
		okSet, err := rdb.SetNX(ctx, key, 1, 24*time.Hour).Result()
		if err != nil {
			return err
		}
		if !okSet {
			// already processed -> idempotent skip
			return nil
		}
		return next(ctx)
	})
}
