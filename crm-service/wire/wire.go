//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"crm/config"
	crm_config "crm/config"
	crm_infra_impl "crm/infra/impl"
	crm_infra_postgre "crm/infra/postgre"
	crm_infra_router "crm/infra/router"
	crm_internal_repo "crm/internal/repo"
	crm_internal_usecase "crm/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(crm_internal_usecase.PermissionUsecase), new(*crm_infra_impl.PermissionUsecase)),
	wire.Bind(new(crm_internal_repo.CustomerRepo), new(*crm_infra_postgre.PostgreCustomer)),
	wire.Bind(new(crm_internal_repo.HistoryRepo), new(*crm_infra_postgre.HistoryRepo)),
	wire.Bind(new(crm_internal_repo.PipelineRepo), new(*crm_infra_postgre.PostgrePipeline)),
	wire.Bind(new(crm_internal_repo.RuleRepo), new(*crm_infra_postgre.PostgreRule)),
	wire.Bind(new(crm_internal_repo.StageRepo), new(*crm_infra_postgre.PostgreStage)),
	crm_infra_postgre.NewBlockRepo,
	crm_config.NewConfigApp,
	crm_infra_postgre.NewContactRepo,
	crm_infra_router.NewContactRouter,
	crm_internal_usecase.NewContactService,
	crm_infra_router.NewCustomerRouter,
	crm_internal_usecase.NewCustomerUsecase,
	crm_infra_postgre.NewFollowRepo,
	crm_infra_postgre.NewFriendRepo,
	crm_infra_postgre.NewGroupRepo,
	crm_infra_postgre.NewHistoryRepo,
	crm_infra_router.NewHistoryRouter,
	crm_internal_usecase.NewHistoryUsecase,
	crm_infra_impl.NewPermissionUsecase,
	crm_infra_router.NewPipelineRouter,
	crm_internal_usecase.NewPipelineUsecase,
	crm_infra_postgre.NewPostgreCustomer,
	crm_infra_postgre.NewPostgrePipeline,
	crm_infra_postgre.NewPostgreRule,
	crm_infra_postgre.NewPostgreStage,
	crm_infra_router.NewRuleRouter,
	crm_internal_usecase.NewRuleUsecase,
	crm_infra_router.NewRuntimeRouter,
	crm_infra_router.NewStageRouter,
	crm_internal_usecase.NewStageUsecase,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}