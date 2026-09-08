package process

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

type fakeServer struct {
	serveErr error

	serveStarted chan struct{}
	serveRelease chan struct{}
	gracefulDone chan struct{}
	forceDone    chan struct{}

	onGraceful func()
	onStop     func()

	onceServe    sync.Once
	onceGraceful sync.Once
	onceStop     sync.Once
}

func newFakeServer() *fakeServer {
	return &fakeServer{
		serveStarted: make(chan struct{}),
		serveRelease: make(chan struct{}),
		gracefulDone: make(chan struct{}),
		forceDone:    make(chan struct{}),
	}
}

func (s *fakeServer) Serve(net.Listener) error {
	s.onceServe.Do(func() { close(s.serveStarted) })
	<-s.serveRelease
	return s.serveErr
}

func (s *fakeServer) GracefulStop() {
	if s.onGraceful != nil {
		s.onGraceful()
	}
	<-s.gracefulDone
	s.onceGraceful.Do(func() {})
}

func (s *fakeServer) Stop() {
	if s.onStop != nil {
		s.onStop()
	}
	s.onceStop.Do(func() {
		close(s.forceDone)
		close(s.gracefulDone)
		select {
		case <-s.serveRelease:
		default:
			close(s.serveRelease)
		}
	})
}

func (s *fakeServer) allowGracefulStop() {
	s.onceGraceful.Do(func() {
		close(s.gracefulDone)
		select {
		case <-s.serveRelease:
		default:
			close(s.serveRelease)
		}
	})
}

type actorFunc func(context.Context)

func (f actorFunc) Run(ctx context.Context) { f(ctx) }

func TestRunPropagatesServeFailureAndWaitsForActor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := newFakeServer()
	serveFailure := errors.New("listener failed")
	server.serveErr = serveFailure
	server.onGraceful = server.allowGracefulStop

	actorStarted := make(chan struct{})
	actorStopped := make(chan struct{})
	actor := actorFunc(func(ctx context.Context) {
		close(actorStarted)
		<-ctx.Done()
		close(actorStopped)
	})

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, server, nil, []Actor{actor}, Config{
			GracefulStopTimeout: 100 * time.Millisecond,
		})
	}()

	<-server.serveStarted
	<-actorStarted
	close(server.serveRelease)

	err := <-result
	if !errors.Is(err, serveFailure) {
		t.Fatalf("Run error = %v, want serve failure", err)
	}
	select {
	case <-actorStopped:
	default:
		t.Fatal("Run returned before actor observed cancellation")
	}
}

func TestRunTreatsContextCancellationAsNormalShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server := newFakeServer()
	server.onGraceful = server.allowGracefulStop

	actorStarted := make(chan struct{})
	actorStopped := make(chan struct{})
	actor := actorFunc(func(ctx context.Context) {
		close(actorStarted)
		<-ctx.Done()
		close(actorStopped)
	})

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, server, nil, []Actor{actor}, Config{
			GracefulStopTimeout: 100 * time.Millisecond,
		})
	}()

	<-server.serveStarted
	<-actorStarted
	cancel()

	if err := <-result; err != nil {
		t.Fatalf("Run error = %v, want nil on context shutdown", err)
	}
	select {
	case <-actorStopped:
	default:
		t.Fatal("Run returned before actor stopped")
	}
}

func TestRunForcesStopWhenGracefulStopExceedsDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	server := newFakeServer()

	stopCalled := make(chan struct{})
	server.onStop = func() {
		select {
		case <-stopCalled:
		default:
			close(stopCalled)
		}
	}

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, server, nil, nil, Config{
			GracefulStopTimeout: 20 * time.Millisecond,
		})
	}()

	<-server.serveStarted
	cancel()

	select {
	case <-stopCalled:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("forced Stop was not called")
	}

	if err := <-result; err != nil {
		t.Fatalf("Run error = %v, want nil on context shutdown", err)
	}
}

func TestRunRejectsUnexpectedNilServeResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := newFakeServer()
	server.onGraceful = server.allowGracefulStop

	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, cancel, server, nil, nil, Config{
			GracefulStopTimeout: 100 * time.Millisecond,
		})
	}()

	<-server.serveStarted
	close(server.serveRelease)

	if err := <-result; !errors.Is(err, ErrServeReturnedWithoutError) {
		t.Fatalf("Run error = %v, want ErrServeReturnedWithoutError", err)
	}
}
