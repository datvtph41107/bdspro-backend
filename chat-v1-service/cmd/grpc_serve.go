package cmd

import (
	"chat/infrastructure/server"

	"github.com/spf13/cobra"
)

var grpcServeCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		grpcPort, err := cmd.Flags().GetString("grpcPort")
		if err != nil {
			return err
		}
		grpcServer, err := server.NewGRPCServer(grpcPort)
		if err != nil {
			return err
		}
		return grpcServer.Start()
	},
}
