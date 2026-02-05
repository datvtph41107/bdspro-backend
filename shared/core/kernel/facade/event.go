package facade

import (
	"context"
	"core/kernel/contract"
)

func WithEventBus(ctx context.Context, b contract.EventBus) context.Context {
	return with(ctx, "facade_event_bus", b)
}

func Dispatch(ctx context.Context, event interface{}) error {
	if v := ctx.Value("facade_event_bus"); v != nil {
		return v.(contract.EventBus).Dispatch(ctx, event)
	}
	panic("Event bus not bound to context")
}
