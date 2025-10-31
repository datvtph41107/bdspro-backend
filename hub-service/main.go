package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "common/models"
	cmd_grpc "hub/cmd/grpc"
	cmd_http "hub/cmd/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// @title Utility API
// @version 1.0
// @description API cho utility service, bao gồm event-queue và các tiện ích khác
// @BasePath /v2/utility
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
var rootCmd = &cobra.Command{
	Use:   "utility",
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
