package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"

	"notification/cmd"
	_ "notification/docs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title Notification API
// @version 1.0
// @description API phần logic thông báo
// @BasePath /v2/notification
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "notification-service",
	Short: "Service xử lý authentication và authorization",
}

func init() {
	rootCmd.AddCommand(cmd.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	closeLogger, err := logging.Configure("notification-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Notification logging: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = closeLogger() }()

	// go func() {
	// 	log.Println("run go-routine")
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	if err := rootCmd.Execute(); err != nil {
		slog.Error(
			"notification command failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
