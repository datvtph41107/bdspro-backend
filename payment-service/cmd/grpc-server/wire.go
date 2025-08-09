//go:build wireinject
// +build wireinject

package grpcserver

import (
	"payment/internal/infrastructure/client"
	"payment/internal/infrastructure/database"
	"payment/internal/infrastructure/repository"
	"payment/internal/interface/handler"
	"payment/internal/server/grpc"
	"payment/internal/usecase"

	"github.com/google/wire"
)

var databaseSet = wire.NewSet(
	database.NewPostgresDB,
)

var repositorySet = wire.NewSet(
	repository.NewPaymentRepository,
	repository.NewWalletRepository,
	repository.NewWalletTransactionRepository,
	repository.NewWalletAuditLogRepository,
	repository.NewWithdrawalRequestRepository,
)

var clientSet = wire.NewSet(
	client.NewNotificationClient,
)

var usecaseSet = wire.NewSet(
	usecase.NewPaymentUsecase,
	usecase.NewNotificationWorker,
)

var handlerSet = wire.NewSet(
	handler.NewPaymentHandler,
)

func InitializeGRPCServer() *grpc.GRPCServer {
	wire.Build(
		databaseSet,
		repositorySet,
		clientSet,
		usecaseSet,
		handlerSet,
		grpc.NewGRPCServer,
	)
	return &grpc.GRPCServer{}
}
