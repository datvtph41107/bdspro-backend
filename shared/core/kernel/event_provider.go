package kernel

import (
	"context"
	"core/kernel/contract"
	"core/kernel/facade"
)

func InjectEventBus(ctx context.Context, bus contract.EventBus) context.Context {
	return facade.WithEventBus(ctx, bus)
}
