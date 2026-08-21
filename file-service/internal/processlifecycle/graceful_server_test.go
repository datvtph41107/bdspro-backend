package processlifecycle

import (
	"sync"
	"testing"
	"time"
)

type fakeGracefulServer struct {
	gracefulRelease chan struct{}
	stopCalled      chan struct{}
	onceStop        sync.Once
}

func newFakeGracefulServer() *fakeGracefulServer {
	return &fakeGracefulServer{
		gracefulRelease: make(chan struct{}),
		stopCalled:      make(chan struct{}),
	}
}

func (s *fakeGracefulServer) GracefulStop() { <-s.gracefulRelease }
func (s *fakeGracefulServer) Stop() {
	s.onceStop.Do(func() {
		close(s.stopCalled)
		close(s.gracefulRelease)
	})
}

func TestStopGracefullyReturnsWhenGracefulCompletes(t *testing.T) {
	server := newFakeGracefulServer()
	close(server.gracefulRelease)
	StopGracefully(server, 50*time.Millisecond)
	select {
	case <-server.stopCalled:
		t.Fatal("forced Stop called after successful graceful stop")
	default:
	}
}

func TestStopGracefullyForcesStopAtDeadline(t *testing.T) {
	server := newFakeGracefulServer()
	StopGracefully(server, 5*time.Millisecond)
	select {
	case <-server.stopCalled:
	default:
		t.Fatal("forced Stop was not called")
	}
}
