package main

import (
	"fmt"
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
	// go func() {
	// 	log.Println("run go-routine")
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
