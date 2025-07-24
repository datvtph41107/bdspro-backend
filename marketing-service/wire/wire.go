//go:build wireinject
// +build wireinject

package wire

import (
	"marketing/config"
	"marketing/infrastructure/client"
	"marketing/infrastructure/repository/postgres"
	"marketing/infrastructure/server/grpc"
	"marketing/internal/domain/repository"
	"marketing/internal/handler"
	"marketing/internal/usecase"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(
		// Database
		provideDB,

		// Clients
		client.NewPaymentClient,
		client.NewNotificationClient,

		// Repositories
		postgres.NewCampaignPostgresRepository,
		wire.Bind(new(repository.CampaignRepository), new(*postgres.CampaignPostgresRepository)),
		postgres.NewPackagePostgresRepository,
		postgres.NewBudgetPostgresRepository,
		postgres.NewAdvertisingSettingsPostgresRepository,

		// Usecases
		usecase.NewCampaignUsecase,
		usecase.NewPackageUsecase,
		usecase.NewPaymentUsecase,
		usecase.NewCampaignDetailUsecase,
		usecase.NewBudgetUsecase,
		usecase.NewAdvertisingSettingsUsecase,

		// Handlers
		handler.NewMarketingGrpcHandler,

		// Server
		grpc.NewGRPCServer,

		// Config
		config.NewConfigApp,
	)
	return &config.ConfigApp{}, nil
}

func provideDB() *gorm.DB {
	// TODO: Implement database connection using env.DB_DNS
	return nil
} 