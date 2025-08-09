//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"payment/config"
	payment_config "payment/config"
	payment_internal_domain_repository "payment/internal/domain/repository"
	payment_internal_infrastructure_client "payment/internal/infrastructure/client"
	payment_internal_infrastructure_database "payment/internal/infrastructure/database"
	payment_internal_infrastructure_repository "payment/internal/infrastructure/repository"
	payment_internal_infrastructure_repository_postgres "payment/internal/infrastructure/repository/postgres"
	payment_internal_interface_handler "payment/internal/interface/handler"
	payment_internal_server_grpc "payment/internal/server/grpc"
	payment_internal_usecase "payment/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(payment_internal_domain_repository.BankRepository), new(*payment_internal_infrastructure_repository_postgres.BankPostgresRepository)),
	payment_internal_interface_handler.NewBankHandler,
	payment_internal_infrastructure_repository_postgres.NewBankPostgresRepository,
	payment_internal_usecase.NewBankUsecase,
	payment_config.NewConfigApp,
	payment_internal_server_grpc.NewGRPCServer,
	payment_internal_infrastructure_client.NewNotificationClient,
	payment_internal_usecase.NewNotificationWorker,
	payment_internal_interface_handler.NewPaymentHandler,
	payment_internal_infrastructure_repository_postgres.NewPaymentMethodRepository,
	payment_internal_infrastructure_repository.NewPaymentRepository,
	payment_internal_usecase.NewPaymentUsecase,
	payment_internal_infrastructure_database.NewPostgresDB,
	payment_internal_infrastructure_repository.NewWalletAuditLogRepository,
	payment_internal_infrastructure_repository.NewWalletRepository,
	payment_internal_infrastructure_repository_postgres.NewWalletTransactionRepository,
	payment_internal_infrastructure_repository_postgres.NewWalletsRepository,
	payment_internal_infrastructure_repository.NewWithdrawalRequestRepository,
	payment_internal_infrastructure_repository_postgres.NewWithdrawalRequestsRepository,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}