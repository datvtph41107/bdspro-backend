package api

import (
	"context"
	"core/kernel.v2/src/runtime"
	"os/signal"
	"syscall"
	"time"
)

type Kernel struct {
	runtime *runtime.Manager
	name    string
	version string
	env     string
}

type Option func(*Kernel)

func New(opts ...Option) *Kernel {
	k := &Kernel{
		runtime: runtime.NewManager(),
		name:    "unknown",
		version: "0.0.0",
		env:     "development",
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

func (k *Kernel) Run(ctx context.Context, app Application) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return app.Stop(shutdownCtx)
}

type Application interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
