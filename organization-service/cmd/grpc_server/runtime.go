package grpcserver

import (
	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		grpcServer := InitializeGRPCServer()
		grpcServer.Start()
	},
}
