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
	ErrHTTPServerMissing  = errors.New("http server is missing")
	ErrCancelMissing      = errors.New("process cancel function is missing")
	ErrUnexpectedHTTPStop = errors.New(
		"http server stopped before process shutdown",
	)
)

// HTTPServer is the process-owned surface required from net/http.Server.
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
	Close() error
}

// Actor is a long-lived process component. Run must block until the actor has
// fully stopped; implementations must not hide another untracked goroutine.
type Actor interface {
	Run(context.Context)
}

type Config struct {
	ShutdownTimeout time.Duration
}

// Run owns the File-service HTTP ingress and long-lived actors.
//
// Shutdown ordering is deliberate:
//  1. observe process cancellation or HTTP failure;
//  2. stop accepting and drain HTTP;
//  3. cancel actors only after ingress has stopped;
//  4. wait for actors before returning.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	httpServer HTTPServer,
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

	actorCtx, actorCancel :=
		context.WithCancel(context.Background())
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

	var rootErr error
	httpDone := false

	select {
	case <-ctx.Done():
		// Parent/process cancellation is normal shutdown.

	case err := <-httpResult:
		httpDone = true

		if !errors.Is(err, http.ErrServerClosed) {
			if err == nil {
				rootErr = ErrUnexpectedHTTPStop
			} else {
				rootErr = fmt.Errorf(
					"http serve: %w",
					err,
				)
			}
		} else if ctx.Err() == nil {
			rootErr = ErrUnexpectedHTTPStop
		}
	}

	cancel()
	shutdownHTTP(httpServer, timeout)

	if !httpDone {
		select {
		case <-httpResult:

		case <-time.After(timeout):
			if rootErr == nil {
				rootErr = fmt.Errorf(
					"http server shutdown exceeded %s",
					timeout,
				)
			}
		}
	}

	// Ingress is closed. Long-lived actors may now stop without a draining
	// request creating work after the actor disappeared.
	actorCancel()
	actorWG.Wait()

	return rootErr
}

func shutdownHTTP(
	server HTTPServer,
	timeout time.Duration,
) {
	// Cleanup needs an independent deadline after process cancellation.
	shutdownCtx, cancel :=
		context.WithTimeout(
			context.Background(),
			timeout,
		)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
	}
}
