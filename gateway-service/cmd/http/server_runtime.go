package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

const defaultHTTPShutdownTimeout = 15 * time.Second

// serveHTTP owns one HTTP serving lifetime after the composition root has
// opened the listener. Readiness flips before drain; graceful shutdown is
// bounded and falls back to force-close before resources borrowed by handlers
// may be released by the composition root.

/**
 * Description
 * "Hãy chạy HTTP server này, trên listener này, trong context này;
 * nếu cần shutdown thì timeout bao lâu;
 * và khi shutdown thì cập nhật readiness state."
 */
func serveHTTP(
	ctx context.Context,
	server *http.Server,
	listener net.Listener,
	shutdownTimeout time.Duration,
	readiness *readinessState,
) error {
	if ctx == nil {
		return errors.New("gateway process context is nil")
	}
	if server == nil {
		return errors.New("gateway HTTP server is nil")
	}
	if listener == nil {
		return errors.New("gateway HTTP listener is nil")
	}
	if shutdownTimeout <= 0 {
		shutdownTimeout = defaultHTTPShutdownTimeout
	}

	serveResult := make(chan error, 1)
	go func() {
		serveResult <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		readiness.beginDrain()
	case err := <-serveResult:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve gateway HTTP: %w", err)
		}
		return nil
	}

	/**
	 * Description
	 * "Hãy shutdown HTTP server một cách graceful;
	 * chờ request đang chạy hoàn tất trong thời gian cho phép;
	 * nếu quá timeout thì force close;
	 * và chỉ kết thúc khi HTTP server thực sự dừng."
	 */
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	// defer cancel()
	shutdownErr := server.Shutdown(shutdownCtx)
	cancel()

	if shutdownErr != nil {
		closeErr := server.Close()
		serveErr := <-serveResult
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			return errors.Join(fmt.Errorf("shutdown gateway HTTP: %w", shutdownErr), fmt.Errorf("force-close gateway HTTP: %w", serveErr), closeErr)
		}
		return errors.Join(fmt.Errorf("shutdown gateway HTTP: %w", shutdownErr), closeErr)
	}

	serveErr := <-serveResult
	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		return fmt.Errorf("serve gateway HTTP during shutdown: %w", serveErr)
	}

	return nil
}
