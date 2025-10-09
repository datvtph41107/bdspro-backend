//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"tx/initial"
	infra_mapper "tx/mapper"
	organization_internal_interface "tx/internal/interface"
	transaction_infra_client "tx/infra/client"
	transaction_infra_handler "tx/infra/handler"
	transaction_infra_mapper "tx/infra/mapper"
	transaction_infra_postgres "tx/infra/postgres"
	transaction_infra_validator "tx/infra/validator"
	transaction_initial "tx/initial"
	transaction_internal_usecases "tx/internal/usecases"
	tx_internal_interface "tx/internal/interface"
	tx_internal_repo "tx/internal/repo"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(organization_internal_interface.IOrganizationClient), new(*transaction_infra_client.OrganizationClient)),
	wire.Bind(new(infra_mapper.TransactionMapper), new(*transaction_infra_handler.AdminTransactionHandler)),
	wire.Bind(new(tx_internal_repo.DealTransactionRepository), new(*transaction_infra_postgres.DealTransactionPostgresRepository)),
	wire.Bind(new(tx_internal_interface.ITransaction), new(*transaction_infra_postgres.TransactionGorm)),
	transaction_infra_handler.NewAdminTransactionHandler,
	transaction_infra_client.NewBdsproClient,
	transaction_infra_handler.NewContractDealHandler,
	transaction_infra_mapper.NewContractDealMapper,
	transaction_infra_validator.NewContractDealValidator,
	transaction_infra_handler.NewCostTypeHandler,
	transaction_infra_mapper.NewCostTypeMapper,
	transaction_infra_postgres.NewCostTypeRepository,
	transaction_internal_usecases.NewCostTypeUsecase,
	transaction_internal_usecases.NewDealContractUsecase,
	transaction_infra_handler.NewDealCostHandler,
	transaction_infra_mapper.NewDealCostMapper,
	transaction_infra_postgres.NewDealCostRepository,
	transaction_internal_usecases.NewDealCostUsecase,
	transaction_infra_validator.NewDealCostValidator,
	transaction_infra_postgres.NewDealTransactionPostgresRepository,
	transaction_initial.NewInitalApp,
	transaction_infra_handler.NewInternalTransactionHandler,
	transaction_infra_client.NewOrganizationClient,
	transaction_infra_postgres.NewTransactionActionPostgresRepository,
	transaction_infra_validator.NewTransactionActionValidator,
	transaction_infra_postgres.NewTransactionGorm,
	transaction_infra_handler.NewTransactionHandler,
	transaction_infra_mapper.NewTransactionMapper,
	transaction_infra_postgres.NewTransactionPostgresRepository,
	transaction_internal_usecases.NewTransactionUsecase,
	transaction_infra_client.NewUserClient,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}