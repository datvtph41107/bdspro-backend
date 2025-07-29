//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	"crm/config"
	crm_config "crm/config"
	crm_infra_client "crm/infra/client"
	crm_infra_impl "crm/infra/impl"
	crm_infra_mapper "crm/infra/mapper"
	crm_infra_postgre "crm/infra/postgre"
	crm_infra_router "crm/infra/router"
	crm_infra_services "crm/infra/services"
	crm_internal_interface "crm/internal/interface"
	crm_internal_repo "crm/internal/repo"
	crm_internal_usecase "crm/internal/usecase"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(crm_internal_interface.INotificationClient), new(*crm_infra_client.NotificationClient)),
	wire.Bind(new(crm_internal_usecase.PermissionUsecase), new(*crm_infra_client.PermissionUsecase)),
	wire.Bind(new(crm_internal_interface.UserClient), new(*crm_infra_client.UserClient)),
	wire.Bind(new(crm_internal_interface.ITransaction), new(*crm_infra_impl.TransactionGorm)),
	wire.Bind(new(crm_internal_repo.BlockRepo), new(*crm_infra_postgre.BlockRepo)),
	wire.Bind(new(crm_internal_repo.ContactRepo), new(*crm_infra_postgre.PostgreContact)),
	wire.Bind(new(crm_internal_repo.DocumentRepo), new(*crm_infra_postgre.DocumentPostgre)),
	wire.Bind(new(crm_internal_repo.FollowRepo), new(*crm_infra_postgre.FollowPostgre)),
	wire.Bind(new(crm_internal_repo.FriendGroupRepo), new(*crm_infra_postgre.FriendGroupRepo)),
	wire.Bind(new(crm_internal_repo.FriendRepo), new(*crm_infra_postgre.FriendPostgre)),
	wire.Bind(new(crm_internal_repo.InvitationInstallRepo), new(*crm_infra_postgre.InvitationInstallPostgre)),
	wire.Bind(new(crm_internal_repo.LeadRepo), new(*crm_infra_postgre.LeadPostgre)),
	wire.Bind(new(crm_internal_repo.PipelineRepo), new(*crm_infra_postgre.PostgrePipeline)),
	wire.Bind(new(crm_internal_repo.ProductCareRepo), new(*crm_infra_postgre.ProductCarePostgre)),
	wire.Bind(new(crm_internal_repo.RuleRepo), new(*crm_infra_postgre.PostgreRule)),
	wire.Bind(new(crm_internal_repo.SharingAccessRepo), new(*crm_infra_postgre.SharingAccessPostgre)),
	wire.Bind(new(crm_internal_repo.StageRepo), new(*crm_infra_postgre.PostgreStage)),
	crm_infra_client.NewBdsproClient,
	crm_infra_postgre.NewBlockRepo,
	crm_infra_services.NewBlockService,
	crm_internal_usecase.NewBlockUsecase,
	crm_config.NewConfigApp,
	crm_infra_mapper.NewContactMapper,
	crm_infra_postgre.NewContactRepo,
	crm_infra_router.NewContactRouter,
	crm_infra_services.NewContactService,
	crm_internal_usecase.NewContactUsecase,
	crm_infra_services.NewCrmInternalService,
	crm_infra_postgre.NewDocumentPostgre,
	crm_infra_postgre.NewFollowPostgre,
	crm_infra_services.NewFollowService,
	crm_internal_usecase.NewFollowUsecase,
	crm_infra_services.NewFriendGroupService,
	crm_internal_usecase.NewFriendGroupUsecase,
	crm_infra_mapper.NewFriendMapper,
	crm_infra_postgre.NewFriendRepo,
	crm_infra_services.NewFriendService,
	crm_internal_usecase.NewFriendUsecase,
	crm_infra_postgre.NewGroupRepo,
	crm_infra_postgre.NewInvitationInstallPostgre,
	crm_infra_services.NewInvitationInstallService,
	crm_internal_usecase.NewInvitationInstallUsecase,
	crm_infra_mapper.NewLeadMapper,
	crm_infra_router.NewLeadRouter,
	crm_infra_services.NewLeadService,
	crm_internal_usecase.NewLeadUsecase,
	crm_infra_client.NewNotificationClient,
	crm_infra_router.NewPermissionRouter,
	crm_infra_client.NewPermissionUsecase,
	crm_infra_router.NewPipelineRouter,
	crm_infra_services.NewPipelineService,
	crm_internal_usecase.NewPipelineUsecase,
	crm_infra_postgre.NewPostgreCustomer,
	crm_infra_postgre.NewPostgrePipeline,
	crm_infra_postgre.NewPostgreRule,
	crm_infra_postgre.NewPostgreStage,
	crm_infra_postgre.NewProductCarePostgre,
	crm_internal_usecase.NewRuleEventUsecase,
	crm_infra_router.NewRuleRouter,
	crm_infra_services.NewRuleService,
	crm_internal_usecase.NewRuleUsecase,
	crm_infra_router.NewRuntimeRouter,
	crm_infra_postgre.NewSharingAccessPostgre,
	crm_infra_services.NewSharingAccessService,
	crm_internal_usecase.NewSharingAccessUsecase,
	crm_infra_router.NewStageRouter,
	crm_infra_services.NewStageService,
	crm_internal_usecase.NewStageUsecase,
	crm_infra_impl.NewTransactionGorm,
	crm_infra_client.NewUserClient,
	crm_infra_client.NewOrganizationClient,
)

func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
