package processlifecycle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const defaultShutdownTimeout = 10 * time.Second

var (
	ErrHTTPServerMissing    = errors.New("http server is missing")
	ErrCancelMissing        = errors.New("process cancel function is missing")
	ErrUnexpectedHTTPStop   = errors.New("http server stopped before process shutdown")
	ErrUnexpectedServerStop = errors.New("internal server stopped before process shutdown")
)

// HTTPServer is the process-owned surface required from net/http.Server.
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	Close() error
}

// InternalServerFunc is a blocking capability-local server such as the
// internal gRPC target. It must return after ctx cancellation.
type InternalServerFunc func(context.Context) error

// Actor is a long-lived process component. Run must block until the actor has
// fully stopped; implementations must not hide another untracked goroutine.
type Actor interface {
	Run(context.Context)
}

type Config struct {
	ShutdownTimeout time.Duration
}

// Run owns the File-service HTTP server, optional internal server and actors.
//
// Shutdown ordering is deliberate:
//  1. observe signal/server failure;
//  2. stop accepting/drain HTTP and cancel the internal server;
//  3. cancel actors only after ingress has stopped;
//  4. wait for all actors/server goroutines before returning.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	httpServer HTTPServer,
	internalServer InternalServerFunc,
	actors []Actor,
	cfg Config,
) error {
	if httpServer == nil {
		return ErrHTTPServerMissing
	}
	if cancel == nil {
		return ErrCancelMissing
	}

	timeout := cfg.ShutdownTimeout
	if timeout <= 0 {
		timeout = defaultShutdownTimeout
	}

	actorCtx, actorCancel := context.WithCancel(context.Background())
	defer actorCancel()

	var actorWG sync.WaitGroup
	for _, actor := range actors {
		if actor == nil {
			continue
		}
		actorWG.Add(1)
		go func(a Actor) {
			defer actorWG.Done()
			a.Run(actorCtx)
		}(actor)
	}

	httpResult := make(chan error, 1)
	go func() {
		httpResult <- httpServer.ListenAndServe()
	}()

	internalCtx, internalCancel := context.WithCancel(context.Background())
	defer internalCancel()
	var internalResult chan error
	if internalServer != nil {
		internalResult = make(chan error, 1)
		go func() {
			internalResult <- internalServer(internalCtx)
		}()
	}

	var rootErr error
	httpDone := false
	internalDone := false
	select {
	case <-ctx.Done():
		// signal/parent cancellation is normal shutdown
	case err := <-httpResult:
		httpDone = true
		if !errors.Is(err, http.ErrServerClosed) {
			if err == nil {
				rootErr = ErrUnexpectedHTTPStop
			} else {
				rootErr = fmt.Errorf("http serve: %w", err)
			}
		} else if ctx.Err() == nil {
			rootErr = ErrUnexpectedHTTPStop
		}
	case err := <-internalResult:
		internalDone = true
		if err == nil {
			rootErr = ErrUnexpectedServerStop
		} else {
			rootErr = fmt.Errorf("internal serve: %w", err)
		}
	}

	// One process shutdown path for signals and either server failure.
	cancel()
	internalCancel()
	shutdownHTTP(httpServer, timeout)

	// Both ingress runners must return before background actors are cancelled.
	// This prevents a draining request from enqueueing work after its workers
	// have already disappeared.
	if internalResult != nil && !internalDone {
		select {
		case <-internalResult:
		case <-time.After(timeout):
			if rootErr == nil {
				rootErr = fmt.Errorf("internal server shutdown exceeded %s", timeout)
			}
		}
	}
	if !httpDone {
		select {
		case <-httpResult:
		case <-time.After(timeout):
			if rootErr == nil {
				rootErr = fmt.Errorf("http server shutdown exceeded %s", timeout)
			}
		}
	}

	// Ingress is now closed. Stop background actors and wait for in-flight
	// jobs before composition may close shared dependencies.
	actorCancel()
	actorWG.Wait()

	return rootErr
}

func shutdownHTTP(server HTTPServer, timeout time.Duration) {
	// Cleanup needs an independent deadline after the process context is
	// cancelled; this context is lifecycle control, not business I/O.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
	}
}
