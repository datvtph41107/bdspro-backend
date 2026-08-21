package processlifecycle

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"
)

type fakeHTTPServer struct {
	serveStarted chan struct{}
	serveRelease chan struct{}
	shutdown     chan struct{}
	closed       chan struct{}
	serveErr     error
	shutdownErr  error
	onShutdown   func()
	onClose      func()
	onceServe    sync.Once
	onceRelease  sync.Once
	onceShutdown sync.Once
	onceClose    sync.Once
}

func newFakeHTTPServer() *fakeHTTPServer {
	return &fakeHTTPServer{
		serveStarted: make(chan struct{}),
		serveRelease: make(chan struct{}),
		shutdown:     make(chan struct{}),
		closed:       make(chan struct{}),
	}
}

func (s *fakeHTTPServer) ListenAndServe() error {
	s.onceServe.Do(func() { close(s.serveStarted) })
	<-s.serveRelease
	if s.serveErr != nil {
		return s.serveErr
	}
	return http.ErrServerClosed
}

func (s *fakeHTTPServer) Shutdown(context.Context) error {
	if s.onShutdown != nil {
		s.onShutdown()
	}
	if s.shutdownErr != nil {
		return s.shutdownErr
	}
	s.onceShutdown.Do(func() {
		close(s.shutdown)
	})
	s.onceRelease.Do(func() { close(s.serveRelease) })
	return nil
}

func (s *fakeHTTPServer) Close() error {
	if s.onClose != nil {
		s.onClose()
	}
	s.onceClose.Do(func() { close(s.closed) })
	s.onceRelease.Do(func() { close(s.serveRelease) })
	return nil
}

type actorFunc func(context.Context)

func (f actorFunc) Run(ctx context.Context) { f(ctx) }

func TestRunStopsIngressBeforeActors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	httpServer := newFakeHTTPServer()

	internalStarted := make(chan struct{})
	internalStopped := make(chan struct{})
	internal := func(ctx context.Context) error {
		close(internalStarted)
		<-ctx.Done()
		close(internalStopped)
		return nil
	}

	actorStarted := make(chan struct{})
	actorStopped := make(chan struct{})
	actor := actorFunc(func(ctx context.Context) {
		close(actorStarted)
		<-ctx.Done()
		select {
		case <-httpServer.shutdown:
		default:
			t.Error("actor cancelled before HTTP ingress stopped")
		}
		select {
		case <-internalStopped:
		default:
			t.Error("actor cancelled before internal ingress stopped")
		}
		close(actorStopped)
	})

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, httpServer, internal, []Actor{actor}, Config{ShutdownTimeout: 100 * time.Millisecond})
	}()

	<-httpServer.serveStarted
	<-internalStarted
	<-actorStarted
	cancel()

	if err := <-result; err != nil {
		t.Fatalf("Run error = %v, want nil", err)
	}
	select {
	case <-actorStopped:
	default:
		t.Fatal("Run returned before actor stopped")
	}
}

func TestRunPropagatesHTTPServeFailureAndWaits(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	httpServer := newFakeHTTPServer()
	serveErr := errors.New("bind failed")
	httpServer.serveErr = serveErr

	actorStopped := make(chan struct{})
	actor := actorFunc(func(ctx context.Context) {
		<-ctx.Done()
		close(actorStopped)
	})

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, httpServer, nil, []Actor{actor}, Config{ShutdownTimeout: 100 * time.Millisecond})
	}()
	<-httpServer.serveStarted
	httpServer.onceRelease.Do(func() { close(httpServer.serveRelease) })

	if err := <-result; !errors.Is(err, serveErr) {
		t.Fatalf("Run error = %v, want serve error", err)
	}
	select {
	case <-actorStopped:
	default:
		t.Fatal("Run returned before actor stopped")
	}
}

func TestRunPropagatesInternalServeFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	httpServer := newFakeHTTPServer()
	internalErr := errors.New("grpc listener failed")
	internal := func(context.Context) error { return internalErr }

	err := Run(ctx, cancel, httpServer, internal, nil, Config{ShutdownTimeout: 100 * time.Millisecond})
	if !errors.Is(err, internalErr) {
		t.Fatalf("Run error = %v, want internal server error", err)
	}
}

func TestRunForcesHTTPCloseWhenGracefulShutdownFails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	httpServer := newFakeHTTPServer()
	httpServer.shutdownErr = errors.New("shutdown failed")

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, httpServer, nil, nil, Config{ShutdownTimeout: 100 * time.Millisecond})
	}()
	<-httpServer.serveStarted
	cancel()

	if err := <-result; err != nil {
		t.Fatalf("Run error = %v, want nil signal shutdown", err)
	}
	select {
	case <-httpServer.closed:
	default:
		t.Fatal("forced HTTP Close was not called")
	}
}
