package signal

import (
	"context"
	contract_runtime "core/kernel.v2/src/contracts/runtime"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type waiter struct {
	sigChan chan os.Signal
}

func NewWaiter(sigs ...os.Signal) contract_runtime.SignalWaiter {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	return &waiter{sigChan: ch}
}

func (w *waiter) Wait(ctx context.Context, shutdown func(context.Context) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.sigChan:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return shutdown(shutdownCtx)
	}
}

func (w *waiter) Notify() <-chan os.Signal {
	return w.sigChan
}
