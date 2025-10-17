//go:build wireinject
// +build wireinject

package grpcserver

import (
	"payment/internal/domain/repository"
	"payment/internal/infrastructure/client"
	"payment/internal/infrastructure/database"
	payment_infrastructure_repository "payment/internal/infrastructure/repository"
	"payment/internal/infrastructure/repository/postgres"
	"payment/internal/interface/handler"
	"payment/internal/job"
	"payment/internal/server/grpc"
	"payment/internal/usecase"

	"github.com/google/wire"
)

var databaseSet = wire.NewSet(
	database.NewPostgresDB,
)

var repositorySet = wire.NewSet(
	payment_infrastructure_repository.NewPaymentRepository,
	payment_infrastructure_repository.NewWalletRepository,
	payment_infrastructure_repository.NewWalletTransactionRepository,
	payment_infrastructure_repository.NewWalletAuditLogRepository,
	payment_infrastructure_repository.NewWithdrawalRequestRepository,
	payment_infrastructure_repository.NewDashboardMetricRepository,
	postgres.NewBankPostgresRepository,
	wire.Bind(new(repository.BankRepository), new(*postgres.BankPostgresRepository)),
)

var clientSet = wire.NewSet(
	client.NewNotificationClient,
)

var usecaseSet = wire.NewSet(
	usecase.NewPaymentUsecase,
	usecase.NewBankUsecase,
	usecase.NewNotificationWorker,
)

var jobSet = wire.NewSet(
	job.NewPaymentDashboardStatsJob,
	job.NewTransactionDashboardJob,
)

var handlerSet = wire.NewSet(
	handler.NewPaymentHandler,
	handler.NewBankHandler,
	handler.NewInternalHandler,
)

func InitializeGRPCServer() *grpc.GRPCServer {
	wire.Build(
		databaseSet,
		repositorySet,
		clientSet,
		usecaseSet,
		jobSet,
		handlerSet,
		grpc.NewGRPCServer,
	)
	return &grpc.GRPCServer{}
}
