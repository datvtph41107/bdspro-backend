package grpcserver

import (
	common_db "common/db"
	"fmt"

	config "organization/config"
	"organization/infrastructure/repository"
	"organization/wire"

	"github.com/spf13/cobra"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.InitConfig(); err != nil {
			return err
		}
		app, cleanup, err := wire.InitializeApp()
		if err != nil {
			return fmt.Errorf("initialize organization app: %w", err)
		}
		defer cleanup()

		if err := common_db.EnsureSchema(common_db.DB, "organization", repository.AutoMigrate); err != nil {
			return err
		}
		app.GRPCServer.Start()
		return nil
	},
}
