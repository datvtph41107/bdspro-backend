package kernel

import "context"

func Event(ctx context.Context, event interface{}) error {
	return Dispatch(ctx, event)
}
