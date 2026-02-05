package facade

import "context"

type ctxKey string

const (
	authKey  ctxKey = "facade_auth"
	clockKey ctxKey = "facade_clock"
	uuidKey  ctxKey = "facade_uuid"
)

func with(ctx context.Context, k ctxKey, v interface{}) context.Context {
	return context.WithValue(ctx, k, v)
}

func get(ctx context.Context, k ctxKey) interface{} {
	if v := ctx.Value(k); v != nil {
		return v
	}
	return nil
}
