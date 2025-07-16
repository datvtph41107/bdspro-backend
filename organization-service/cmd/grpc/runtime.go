package grpcserver

import (
	"log"
	config "organization/config"

	"organization/infrastructure/repository"
	"organization/wire"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		app, err := wire.InitializeApp()
		repository.AutoMigrate()
		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}
		app.GRPCServer.Start()
	},
}
