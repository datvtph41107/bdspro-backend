//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/payment-service/internal/infrastructure/client"
	"github.com/tuandq2112/payment-service/internal/infrastructure/database"
	"github.com/tuandq2112/payment-service/internal/infrastructure/repository"
	"github.com/tuandq2112/payment-service/internal/interface/handler"
	"github.com/tuandq2112/payment-service/internal/server/grpc"
	"github.com/tuandq2112/payment-service/internal/usecase"
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
