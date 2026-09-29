package main

import (
	"fmt"
	"os"

	"common/logging"
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
	os.Exit(runProcess(rootCmd.Execute))
}

func runProcess(execute func() error) int {
	closeLogger, err := logging.Configure("crm-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure CRM logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := execute(); err != nil {
		// Preserve Cobra's existing functional command-error output.
		fmt.Println(err)
		return 1
	}
	return 0
}
