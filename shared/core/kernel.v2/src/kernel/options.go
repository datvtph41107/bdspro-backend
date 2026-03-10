package kernel

import (
	"context"
	"os"
	"syscall"
	"time"
)

type Option func(*App)

func WithSignals(signals ...os.Signal) Option {
	return func(a *App) {
		if len(signals) == 0 {
			return
		}
		a.signals = signals
	}
}

func WithObserver(obs Observer) Option {
	return func(a *App) {
		if obs != nil {
			a.observers = append(a.observers, obs)
		}
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(a *App) {
		if timeout > 0 {
			a.defaultStopTimeout = timeout
		}
	}
}

func WithStartupTimeout(timeout time.Duration) Option {
	return func(a *App) {
		if timeout > 0 {
			a.defaultStartTimeout = timeout
		}
	}
}

func WithContext(ctx context.Context) Option {
	return func(a *App) {
		if ctx != nil {
			a.baseCtx = ctx
		}
	}
}

func defaultSignals() []os.Signal {
	return []os.Signal{syscall.SIGINT, syscall.SIGTERM}
}
