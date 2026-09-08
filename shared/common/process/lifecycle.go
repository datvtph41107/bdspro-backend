package process

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const defaultGracefulStopTimeout = 10 * time.Second

var (
	ErrServerMissing             = errors.New("process server is missing")
	ErrCancelMissing             = errors.New("process cancel function is missing")
	ErrServeReturnedWithoutError = errors.New("server Serve returned without an error before shutdown")
)

// Server is the process-owned lifecycle surface required from gRPC.
// *grpc.Server satisfies this interface directly.
type Server interface {
	Serve(net.Listener) error
	GracefulStop()
	Stop()
}

// Actor is a long-running process-owned component. Run must block until the
// actor finishes; it must not hide another goroutine behind this boundary.
type Actor interface {
	Run(context.Context)
}

// ActorFunc adapts a blocking Run function to Actor.
type ActorFunc func(context.Context)

func (f ActorFunc) Run(ctx context.Context) {
	if f != nil {
		f(ctx)
	}
}

type Config struct {
	GracefulStopTimeout time.Duration
}

// Run owns the server goroutine and all supplied actor goroutines.
//
// The caller owns ctx/cancel so other capability-local runners can share the
// same process cancellation domain. Run guarantees that cancellation is sent,
// the gRPC server is stopped, and supplied actors have returned before it
// returns to the composition root.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	server Server,
	listener net.Listener,
	actors []Actor,
	cfg Config,
) error {
	if server == nil {
		return ErrServerMissing
	}
	if cancel == nil {
		return ErrCancelMissing
	}

	timeout := cfg.GracefulStopTimeout
	if timeout <= 0 {
		timeout = defaultGracefulStopTimeout
	}

	var actorsWG sync.WaitGroup
	for _, actor := range actors {
		if actor == nil {
			continue
		}
		actorsWG.Add(1)
		go func(a Actor) {
			defer actorsWG.Done()
			a.Run(ctx)
		}(actor)
	}

	serveResult := make(chan error, 1)
	go func() {
		serveResult <- server.Serve(listener)
	}()

	var rootErr error
	select {
	case <-ctx.Done():
		// Signal/parent cancellation is a normal process shutdown.
	case err := <-serveResult:
		if err == nil {
			rootErr = ErrServeReturnedWithoutError
		} else {
			rootErr = fmt.Errorf("serve: %w", err)
		}
	}

	// One shutdown path for both signal cancellation and Serve failure.
	cancel()
	stopServer(server, timeout)
	actorsWG.Wait()

	return rootErr
}

func stopServer(server Server, timeout time.Duration) {
	gracefulDone := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(gracefulDone)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-gracefulDone:
		return
	case <-timer.C:
		server.Stop()
	}

	// grpc.Server.Stop unblocks GracefulStop. Keep a second bound so a broken
	// implementation cannot hold process shutdown forever.
	forceTimer := time.NewTimer(timeout)
	defer forceTimer.Stop()
	select {
	case <-gracefulDone:
	case <-forceTimer.C:
	}
}
