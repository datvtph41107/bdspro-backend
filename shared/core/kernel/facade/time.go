package facade

import (
	"context"
	"core/kernel/contract"
	"time"
)

func WithClock(ctx context.Context, c contract.Clock) context.Context {
	return with(ctx, clockKey, c)
}

func Now(ctx context.Context) time.Time {
	if v := get(ctx, clockKey); v != nil {
		return v.(contract.Clock).Now()
	}
	return time.Now()
}
