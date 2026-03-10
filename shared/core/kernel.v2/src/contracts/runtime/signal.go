package contract_runtime

import (
	"context"
	"os"
)

type SignalWaiter interface {
	Wait(ctx context.Context, shutdownFunc func(context.Context) error) error
	Notify() <-chan os.Signal
}
