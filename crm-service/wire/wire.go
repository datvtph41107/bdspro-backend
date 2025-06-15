//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"crm/config"
	crm_config "crm/config"
	crm_infra_client "crm/infra/client"
	crm_infra_grpc "crm/infra/grpc"
	crm_infra_postgre "crm/infra/postgre"
	crm_infra_router "crm/infra/router"
	crm_internal_interface "crm/internal/interface"
	crm_internal_repo "crm/internal/repo"
	crm_internal_usecase "crm/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(crm_internal_interface.NotiClient), new(*crm_infra_client.NotiClient)),
	wire.Bind(new(crm_internal_usecase.PermissionUsecase), new(*crm_infra_client.PermissionUsecase)),
	wire.Bind(new(crm_internal_repo.BlockRepo), new(*crm_infra_postgre.BlockRepo)),
	wire.Bind(new(crm_internal_repo.ContactRepo), new(*crm_infra_postgre.PostgreContact)),
	wire.Bind(new(crm_internal_repo.CustomerRepo), new(*crm_infra_postgre.PostgreCustomer)),
	wire.Bind(new(crm_internal_repo.FollowRepo), new(*crm_infra_postgre.FollowPostgre)),
	wire.Bind(new(crm_internal_repo.FriendRepo), new(*crm_infra_postgre.FriendRepo)),
	wire.Bind(new(crm_internal_repo.FriendGroupRepo), new(*crm_infra_postgre.FriendGroupRepo)),
	wire.Bind(new(crm_internal_repo.HistoryRepo), new(*crm_infra_postgre.HistoryRepo)),
	wire.Bind(new(crm_internal_repo.PipelineRepo), new(*crm_infra_postgre.PostgrePipeline)),
	wire.Bind(new(crm_internal_repo.RuleRepo), new(*crm_infra_postgre.PostgreRule)),
	wire.Bind(new(crm_internal_repo.SharingAccessRepo), new(*crm_infra_postgre.SharingAccessPostgre)),
	wire.Bind(new(crm_internal_repo.StageRepo), new(*crm_infra_postgre.PostgreStage)),
	crm_infra_postgre.NewBlockRepo,
	crm_infra_grpc.NewBlockService,
	crm_internal_usecase.NewBlockUsecase,
	crm_config.NewConfigApp,
	crm_infra_postgre.NewContactRepo,
	crm_infra_router.NewContactRouter,
	crm_infra_grpc.NewContactService,
	crm_internal_usecase.NewContactUsecase,
	crm_infra_router.NewCustomerRouter,
	crm_infra_grpc.NewCustomerService,
	crm_internal_usecase.NewCustomerUsecase,
	crm_infra_postgre.NewFollowPostgre,
	crm_infra_grpc.NewFollowService,
	crm_internal_usecase.NewFollowUsecase,
	crm_infra_grpc.NewFriendGroupService,
	crm_internal_usecase.NewFriendGroupUsecase,
	crm_infra_postgre.NewFriendRepo,
	crm_infra_grpc.NewFriendService,
	crm_internal_usecase.NewFriendUsecase,
	crm_infra_postgre.NewGroupRepo,
	crm_infra_postgre.NewHistoryRepo,
	crm_infra_router.NewHistoryRouter,
	crm_infra_grpc.NewHistoryService,
	crm_internal_usecase.NewHistoryUsecase,
	crm_infra_client.NewNotiClient,
	crm_infra_router.NewPermissionRouter,
	crm_infra_client.NewPermissionUsecase,
	crm_infra_router.NewPipelineRouter,
	crm_infra_grpc.NewPipelineService,
	crm_internal_usecase.NewPipelineUsecase,
	crm_infra_postgre.NewPostgreCustomer,
	crm_infra_postgre.NewPostgrePipeline,
	crm_infra_postgre.NewPostgreRule,
	crm_infra_postgre.NewPostgreStage,
	crm_internal_usecase.NewRuleEventUsecase,
	crm_infra_router.NewRuleRouter,
	crm_infra_grpc.NewRuleService,
	crm_internal_usecase.NewRuleUsecase,
	crm_infra_router.NewRuntimeRouter,
	crm_infra_postgre.NewSharingAccessPostgre,
	crm_infra_grpc.NewSharingAccessService,
	crm_internal_usecase.NewSharingAccessUsecase,
	crm_infra_router.NewStageRouter,
	crm_infra_grpc.NewStageService,
	crm_internal_usecase.NewStageUsecase,
	crm_infra_client.NewUserClient,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}