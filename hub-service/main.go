package main

import (
	"fmt"
	"os"

	_ "common/models"
	cmd_grpc "hub/cmd/grpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Hub is an internal gRPC service. gateway-service owns the public HTTP boundary
// and registers Hub grpc-gateway handlers generated from Proto.
var rootCmd = &cobra.Command{
	Use:   "hub",
	Short: "Hub internal gRPC service",
}

func init() {
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
