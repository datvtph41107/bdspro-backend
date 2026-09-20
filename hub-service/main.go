package main

import (
	"fmt"
	"os"

	"common/logging"
	_ "common/models"
	cmd_grpc "hub/cmd/grpc"

	"github.com/spf13/cobra"
)

// Hub is an internal gRPC service. gateway-service owns the public HTTP boundary
// and registers Hub grpc-gateway handlers generated from Proto.
var rootCmd = &cobra.Command{
	Use:   "hub",
	Short: "Hub internal gRPC service",
}

func init() {
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)
}

func main() {
	os.Exit(runProcess(rootCmd.Execute))
}

func runProcess(execute func() error) int {
	closeLogger, err := logging.Configure("hub-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Hub logging: %v\n", err)
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
