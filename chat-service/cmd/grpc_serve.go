package cmd

import (
	"chat/internal/server"

	"github.com/spf13/cobra"
)

var grpcServeCmd = &cobra.Command{
	Use:   "grpc-serve",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		grpcPort, err := cmd.Flags().GetString("grpcPort")
		if err != nil {
			panic(err)
		}
		grpcServer := server.NewGRPCServer(grpcPort)
		grpcServer.Start()
	},
}
