package main

import (
	"fmt"
	"os"
	cmd_grpc "payment/cmd/grpc"
	cmd_http "payment/cmd/http"
	_ "payment/docs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "payment",
	Short: "Một ứng dụng quản lý server gRPC và HTTP",
}

func init() {
	// Đăng ký các subcommand vào rootCmd
	rootCmd.AddCommand(cmd_http.HttpCmd)
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)

	// Cấu hình cho Cobra để lấy thông tin từ các tệp cấu hình (nếu cần)
	viper.AutomaticEnv() // Để lấy thông tin từ biến môi trường (nếu cần)
}

// @title Example API
// @version 1.0
// @description This is a sample API
// @BasePath /v2/payment
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
