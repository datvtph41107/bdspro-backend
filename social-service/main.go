package main

import (
	"common/logging"
	"fmt"
	"log/slog"
	"os"
	cmd_grpc "social/cmd/grpc"
	cmd_http "social/cmd/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Một ứng dụng quản lý server gRPC và HTTP",
}

// @title Social Service
// @version 1.0
// @description Social Service
// @BasePath /v2/social
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func init() {
	// Đăng ký các subcommand vào rootCmd
	rootCmd.AddCommand(cmd_http.HttpCmd)
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)

	// Cấu hình cho Cobra để lấy thông tin từ các tệp cấu hình (nếu cần)
	viper.AutomaticEnv() // Để lấy thông tin từ biến môi trường (nếu cần)
}

func main() {
	os.Exit(runProcess(rootCmd.Execute))
}

func runProcess(execute func() error) int {
	closeLogger, err := logging.Configure("social-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Social logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := execute(); err != nil {
		slog.Error("social service failed", slog.Any("error", err))
		return 1
	}

	return 0
}
