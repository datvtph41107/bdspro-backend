//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"payment/config"
	payment_config "payment/config"
	payment_external_mod_gate_repo "payment/external/mod/gate/repo"
	payment_external_mod_gate_services "payment/external/mod/gate/services"
	payment_external_mod_wallet_repo "payment/external/mod/wallet/repo"
	payment_external_mod_wallet_router "payment/external/mod/wallet/router"
	payment_external_server "payment/external/server"
	payment_internal_wallet_repo "payment/internal/wallet/repo"
	payment_internal_wallet_usecases "payment/internal/wallet/usecases"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(payment_internal_wallet_repo.AccountRepo), new(*payment_external_mod_wallet_repo.AccountGormRepo)),
	wire.Bind(new(payment_internal_wallet_repo.TransactionRepo), new(*payment_external_mod_wallet_repo.TransactionGormRepo)),
	payment_external_mod_wallet_repo.NewAccountGormRepo,
	payment_external_mod_wallet_router.NewAccountRouter,
	payment_internal_wallet_usecases.NewAccountUsecase,
	payment_external_mod_gate_services.NewApiClient,
	payment_config.NewConfigApp,
	payment_external_server.NewGrpcServer,
	payment_external_mod_gate_repo.NewPaymentRepo,
	payment_external_mod_gate_services.NewPaymentService,
	payment_external_mod_wallet_repo.NewTransactionRepo,
	payment_external_mod_wallet_router.NewTransactionRouter,
	payment_internal_wallet_usecases.NewTransactionUsecase,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}