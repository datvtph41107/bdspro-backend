package grpcserver

import (
	"payment/appconfig"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()
		grpcServer := InitializeGRPCServer()
		grpcServer.Start()
	},
}
