package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"

	"relay/config"
	"relay/server"
)

func main() {
	closeLogger, err := logging.Configure("relay-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Relay logging: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = closeLogger() }()

	if err := run(); err != nil {
		slog.Error(
			"relay startup failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}

func run() error {
	if err := os.Setenv("TZ", "Asia/Ho_Chi_Minh"); err != nil {
		return fmt.Errorf("set relay timezone: %w", err)
	}
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load relay config: %w", err)
	}

	wsServer, err := server.NewWSServer(config.AppProperties.Server.Port)
	if err != nil {
		return err
	}
	return wsServer.Start()
}
