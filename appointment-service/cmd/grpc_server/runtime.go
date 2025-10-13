package grpcserver

import (
	"appointment/config"
	"appointment/wire"
	"log"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()

		app, err := wire.InitializeApp()
		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}
		grpcServer := app.GRPCServer
		grpcServer.Start()
	},
}
