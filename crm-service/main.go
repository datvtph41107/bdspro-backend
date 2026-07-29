package main

import (
	"fmt"
	"os"

	_ "common/models"
	cmdgrpc "crm/cmd/grpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CRM is an internal gRPC service. gateway-service owns all public/admin HTTP
// routes and invokes CRM through gRPC, including SEO and public content APIs.
var rootCmd = &cobra.Command{
	Use:   "crm",
	Short: "CRM internal gRPC service",
}

func init() {
	rootCmd.AddCommand(cmdgrpc.GrpcCmd)
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
