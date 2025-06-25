package grpcserver

import (
	"github.com/spf13/cobra"
	"github.com/tuandq2112/organization-service/appconfig"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		appconfig.InitConfig()
		// repository.AutoMigrate()
		grpcServer := InitializeGRPCServer()
		grpcServer.Start()
	},
}
