package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"

	"payment/cmd"
)

// @title Payment Service
// @version 1.0
// @description Payment Service
// @BasePath /v2/payment
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	os.Exit(runProcess(cmd.RootCmd.Execute))
}

func runProcess(execute func() error) int {
	closeLogger, err := logging.Configure("payment-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Payment logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := os.Setenv("TZ", "Asia/Ho_Chi_Minh"); err != nil {
		slog.Error("set Payment timezone", slog.Any("error", err))
		return 1
	}
	if err := execute(); err != nil {
		slog.Error("payment service failed", slog.Any("error", err))
		return 1
	}
	return 0
}
