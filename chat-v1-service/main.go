package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"

	"chat/cmd"
)

// @title Chat Service API
// @version 1.0
// @description API for chat service
// @BasePath /v2/chat
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")

	closeLogger, err := logging.Configure("chat-v1-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Chat V1 logging: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = closeLogger() }()

	if err := cmd.RootCmd.Execute(); err != nil {
		slog.Error(
			"chat-v1 command failed",
			slog.Any("error", err),
		)
		os.Exit(-1)
	}
}
