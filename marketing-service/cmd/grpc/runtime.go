package grpcserver

import (
	"log"

	"marketing/config"
	"marketing/infrastructure/repository/postgres"
	"marketing/wire"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		app, err := wire.InitializeApp()
		postgres.Migrate()

		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}

		app.GRPCServer.Start()
	},
}
