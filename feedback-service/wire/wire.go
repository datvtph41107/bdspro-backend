//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"feedback/initial"
	feedback_infra_client "feedback/infra/client"
	feedback_infra_handler "feedback/infra/handler"
	feedback_infra_mapper "feedback/infra/mapper"
	feedback_infra_postgre "feedback/infra/postgre"
	feedback_initial "feedback/initial"
	feedback_internal_usecase "feedback/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	feedback_initial.NewApp,
	feedback_infra_handler.NewInternalHandler,
	feedback_infra_handler.NewRateHandler,
	feedback_infra_mapper.NewRateMapper,
	feedback_infra_postgre.NewRatePostgres,
	feedback_internal_usecase.NewRateUsecase,
	feedback_initial.NewRedisService,
	feedback_infra_handler.NewRegistryHandler,
	feedback_infra_postgre.NewRegistryPostgres,
	feedback_internal_usecase.NewRegistryUsecase,
	feedback_infra_handler.NewReportAdminHandler,
	feedback_infra_handler.NewReportHandler,
	feedback_infra_mapper.NewReportMapper,
	feedback_infra_postgre.NewReportPostgres,
	feedback_infra_postgre.NewReportProofPostgres,
	feedback_infra_handler.NewReportReasonHandler,
	feedback_infra_postgre.NewReportReasonPostgres,
	feedback_internal_usecase.NewReportReasonUsecase,
	feedback_internal_usecase.NewReportUsecase,
	feedback_infra_postgre.NewTransactionPostgres,
	feedback_infra_client.NewUserClient,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}