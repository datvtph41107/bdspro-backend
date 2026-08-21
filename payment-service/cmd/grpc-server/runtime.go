package grpcserver

import (
	"payment/appconfig"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		appconfig.InitConfig()
		grpcServer, cleanup, err := InitializeGRPCServer()
		if err != nil {
			return err
		}
		defer cleanup()
		grpcServer.Start()
		return nil
	},
}
