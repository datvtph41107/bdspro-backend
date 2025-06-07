package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tuandq3/go-chat-service/internal/server"
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