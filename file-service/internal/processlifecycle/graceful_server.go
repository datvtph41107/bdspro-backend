package processlifecycle

import "time"

// GracefulServer is the shutdown surface implemented by *grpc.Server.
type GracefulServer interface {
	GracefulStop()
	Stop()
}

// StopGracefully bounds a graceful server stop and falls back to Stop.
func StopGracefully(server GracefulServer, timeout time.Duration) {
	if server == nil {
		return
	}
	if timeout <= 0 {
		timeout = defaultShutdownTimeout
	}

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

	// Stop should unblock GracefulStop. Keep one final bound so a broken
	// implementation cannot hold process shutdown forever.
	forceTimer := time.NewTimer(timeout)
	defer forceTimer.Stop()
	select {
	case <-gracefulDone:
	case <-forceTimer.C:
	}
}
