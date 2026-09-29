package processlifecycle

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeHTTPServer struct {
	result chan error

	shutdownErr error

	shutdownCalled chan struct{}
	closeCalled    chan struct{}

	shutdownOnce sync.Once
	closeOnce    sync.Once
}

func newFakeHTTPServer() *fakeHTTPServer {
	return &fakeHTTPServer{
		result:         make(chan error, 1),
		shutdownCalled: make(chan struct{}),
		closeCalled:    make(chan struct{}),
	}
}

func (s *fakeHTTPServer) ListenAndServe() error {
	return <-s.result
}

func (s *fakeHTTPServer) Shutdown(
	context.Context,
) error {
	s.shutdownOnce.Do(func() {
		close(s.shutdownCalled)
	})

	if s.shutdownErr != nil {
		return s.shutdownErr
	}

	select {
	case s.result <- http.ErrServerClosed:
	default:
	}

	return nil
}

func (s *fakeHTTPServer) Close() error {
	s.closeOnce.Do(func() {
		close(s.closeCalled)
	})

	select {
	case s.result <- http.ErrServerClosed:
	default:
	}

	return nil
}

type actorFunc func(context.Context)

func (f actorFunc) Run(ctx context.Context) {
	f(ctx)
}

func TestRunStopsIngressBeforeActors(
	t *testing.T,
) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	server := newFakeHTTPServer()

	actorStarted := make(chan struct{})
	actorStopped := make(chan struct{})

	actor := actorFunc(func(ctx context.Context) {
		close(actorStarted)
		<-ctx.Done()
		close(actorStopped)
	})

	result := make(chan error, 1)

	go func() {
		result <- Run(
			ctx,
			cancel,
			server,
			[]Actor{actor},
			Config{
				ShutdownTimeout: 50 * time.Millisecond,
			},
		)
	}()

	<-actorStarted
	cancel()

	select {
	case <-server.shutdownCalled:
	case <-time.After(time.Second):
		t.Fatal("HTTP shutdown was not called")
	}

	select {
	case <-actorStopped:
	case <-time.After(time.Second):
		t.Fatal("actor did not stop")
	}

	if err := <-result; err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunPropagatesHTTPServeFailureAndWaits(
	t *testing.T,
) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	defer cancel()

	server := newFakeHTTPServer()

	actorStopped := make(chan struct{})

	actor := actorFunc(func(ctx context.Context) {
		<-ctx.Done()
		close(actorStopped)
	})

	result := make(chan error, 1)

	go func() {
		result <- Run(
			ctx,
			cancel,
			server,
			[]Actor{actor},
			Config{
				ShutdownTimeout: 50 * time.Millisecond,
			},
		)
	}()

	server.result <- errors.New("serve failed")

	err := <-result
	if err == nil ||
		!strings.Contains(
			err.Error(),
			"http serve: serve failed",
		) {
		t.Fatalf("Run() error = %v", err)
	}

	select {
	case <-actorStopped:
	default:
		t.Fatal(
			"Run returned before actor stopped",
		)
	}
}

func TestRunForcesHTTPCloseWhenGracefulShutdownFails(
	t *testing.T,
) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	server := newFakeHTTPServer()
	server.shutdownErr =
		errors.New("shutdown failed")

	result := make(chan error, 1)

	go func() {
		result <- Run(
			ctx,
			cancel,
			server,
			nil,
			Config{
				ShutdownTimeout: 50 * time.Millisecond,
			},
		)
	}()

	cancel()

	select {
	case <-server.closeCalled:
	case <-time.After(time.Second):
		t.Fatal(
			"HTTP Close was not called after Shutdown failure",
		)
	}

	if err := <-result; err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunRejectsMissingProcessDependencies(
	t *testing.T,
) {
	server := newFakeHTTPServer()

	if err := Run(
		context.Background(),
		func() {},
		nil,
		nil,
		Config{},
	); !errors.Is(err, ErrHTTPServerMissing) {
		t.Fatalf(
			"missing server error = %v",
			err,
		)
	}

	if err := Run(
		context.Background(),
		nil,
		server,
		nil,
		Config{},
	); !errors.Is(err, ErrCancelMissing) {
		t.Fatalf(
			"missing cancel error = %v",
			err,
		)
	}
}
