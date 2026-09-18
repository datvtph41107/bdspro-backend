package main

import (
	"assistant/cmd/grpc"
	"common/logging"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// @title Assistant API
// @version 1.0
// @description API cho assistant service, bao gồm event-queue và các tiện ích khác
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	closeLogger, err := logging.Configure("assistant-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Assistant logging: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = closeLogger() }()

	// Set timezone về UTC (múi 0)
	time.Local = time.UTC

	slog.Info(
		"assistant timezone configured",
		slog.String("timezone", time.Local.String()),
	)
	slog.Info("assistant service starting")

	grpc.RunGRPCServer()
}
