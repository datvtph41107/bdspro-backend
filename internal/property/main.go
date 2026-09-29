package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	cmd_grpc "bdspro/cmd/grpc"
	cmd_http "bdspro/cmd/http"
	"common/logging"
	_ "common/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title BDSPRO API
// @version 1.0
// @description API phần logic chính của BĐSPro
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Một ứng dụng quản lý server gRPC và HTTP",
}

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
	os.Exit(runProcess(rootCmd.Execute))
}

func runProcess(execute func() error) int {
	closeLogger, err := logging.Configure("bdspro-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure BDSPro logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := execute(); err != nil {
		// Preserve the existing Cobra functional output contract.
		fmt.Println(err)
		return 1
	}

	return 0
}

// userProductService.UserAssetService = userAssetService
// 	userProductService.UserPostService = userPostService
// 	uPostUsecase.UProductUsecase = uProductUsecase
// 	uAssetUsecase.ProductUc = uProductUsecase
