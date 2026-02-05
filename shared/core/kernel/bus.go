package kernel

import (
	"context"
	"core/kernel/container"
	"core/kernel/event/bus"
)

func Dispatch(ctx context.Context, payload interface{}) error {
	b, err := container.ResolveTyped[*bus.Bus](appC, "bus")
	if err != nil {
		panic(err)
	}
	return b.Dispatch(ctx, payload)
}

func DispatchIfNeed(ctx context.Context, payload interface{}) error {
	b, err := container.ResolveTyped[*bus.Bus](appC, "bus")
	if err != nil {
		panic(err)
	}
	return b.DispatchIfNeed(ctx, payload)
}
