//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	utility_infra_handler "utility/infra/handler"
	utility_infra_mapper "utility/infra/mapper"
	utility_infra_postgre "utility/infra/postgre"
	"utility/initial"
	utility_initial "utility/initial"
	utility_internal_repo "utility/internal/repo"
	utility_internal_usecase "utility/internal/usecase"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(utility_internal_repo.IEventQueueRepo), new(*utility_infra_postgre.EventQueueRepo)),
	utility_infra_mapper.NewEventQueueMapper,
	utility_infra_postgre.NewEventQueueRepo,
	utility_infra_handler.NewEventQueueService,
	utility_internal_usecase.NewEventQueueUsecase,
	utility_infra_handler.NewUtilityInternalService,
	utility_initial.NewInitialApp,
)

func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
