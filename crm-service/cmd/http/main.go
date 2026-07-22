package cmd_http

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the CRM public content HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := NewPublicHTTPServer()
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() { errCh <- server.ListenAndServe() }()

		select {
		case err := <-errCh:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		}
	},
}
