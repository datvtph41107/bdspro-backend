package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "common/models"
	cmd_grpc "crm/cmd/grpc"
	cmd_http "crm/cmd/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title CRM API
// @version 1.0
// @description API phần logic chính của CRM. Đã chuyển sang grpc
// @BasePath /v2/crm
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "crm",
	Short: "Một ứng dụng quản lý server gRPC và HTTP",
}

func init() {
	rootCmd.AddCommand(cmd_http.HttpCmd)
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	go func() {
		log.Println("run go-routine")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
