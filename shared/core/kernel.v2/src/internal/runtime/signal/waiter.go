package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Signal interface {
	Wait(ctx context.Context, stop func(context.Context) error) error
}

type signalHandler struct {
	ch chan os.Signal
}

// New tạo signal handler, mặc định SIGINT + SIGTERM
func New(sigs ...os.Signal) Signal {
	if len(sigs) == 0 {
		sigs = []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	return &signalHandler{ch: ch}
}

// Wait block tới khi nhận signal hoặc ctx bị cancel
func (s *signalHandler) Wait(ctx context.Context, stop func(context.Context) error) error {
	defer signal.Stop(s.ch)

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-s.ch:
		// graceful shutdown với timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		return stop(shutdownCtx)
	}
}
