package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "chat-service",
	Short: "Go Chat Service is a gRPC based chat service",
	Long:  `Go Chat Service is a gRPC based chat service that allows users to communicate in real-time.`,
}

func init() {
	serveCmd.Flags().String("httpPort", "8080", "httpPort")
	RootCmd.AddCommand(serveCmd)

	grpcServeCmd.Flags().String("grpcPort", "8208", "grpcPort")
	RootCmd.AddCommand(grpcServeCmd)

	// wsCmd.Flags().String("wsPort", "3003", "wsPort")
	// RootCmd.AddCommand(wsCmd)
	RootCmd.AddCommand(swaggerCMD)
	RootCmd.AddCommand(prepareSwaggerCMD)
}
