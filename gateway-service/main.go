package main

import (
	"fmt"
	cmd_grpc "gateway/cmd/grpc"
	cmd_http "gateway/cmd/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "API Gateway cho microservices",
}

func init() {
	// Đăng ký các subcommand
	rootCmd.AddCommand(cmd_http.HttpCmd)
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)

	// Cấu hình viper
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
