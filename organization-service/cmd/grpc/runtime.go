package grpcserver

import (
	common_db "common/db"
	"fmt"
	"log"

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

		schemaPolicy, err := common_db.LoadSchemaPolicy("organization")
		if err != nil {
			return fmt.Errorf("load Organization schema policy: %w", err)
		}
		log.Printf("Organization database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
		if err := common_db.ApplySchemaPolicy(common_db.DB, schemaPolicy, repository.AutoMigrate); err != nil {
			return fmt.Errorf("apply Organization schema policy: %w", err)
		}
		app.GRPCServer.Start()
		return nil
	},
}
