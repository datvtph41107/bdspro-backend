//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"crm/initial"
	crm_infra_client "crm/infra/client"
	crm_infra_handler "crm/infra/handler"
	crm_infra_impl "crm/infra/impl"
	crm_infra_mapper "crm/infra/mapper"
	crm_infra_postgre "crm/infra/postgre"
	crm_initial "crm/initial"
	crm_internal_interface_provider "crm/internal/interface/provider"
	crm_internal_repo "crm/internal/repo"
	crm_internal_usecase "crm/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(crm_internal_usecase.PermissionUsecase), new(*crm_infra_client.PermissionUsecase)),
	wire.Bind(new(crm_internal_interface_provider.UserClient), new(*crm_infra_client.UserClient)),
	wire.Bind(new(crm_internal_interface_provider.ITransaction), new(*crm_infra_impl.TransactionGorm)),
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
	crm_infra_handler.NewBlockService,
	crm_internal_usecase.NewBlockUsecase,
	crm_infra_mapper.NewContactMapper,
	crm_infra_postgre.NewContactRepo,
	crm_infra_handler.NewContactService,
	crm_internal_usecase.NewContactUsecase,
	crm_infra_handler.NewCrmInternalService,
	crm_infra_postgre.NewDocumentPostgre,
	crm_infra_postgre.NewFollowPostgre,
	crm_infra_handler.NewFollowService,
	crm_internal_usecase.NewFollowUsecase,
	crm_infra_handler.NewFriendGroupService,
	crm_internal_usecase.NewFriendGroupUsecase,
	crm_infra_mapper.NewFriendMapper,
	crm_infra_postgre.NewFriendRepo,
	crm_infra_handler.NewFriendService,
	crm_internal_usecase.NewFriendUsecase,
	crm_infra_postgre.NewGroupRepo,
	crm_initial.NewInitialApp,
	crm_infra_postgre.NewInvitationInstallPostgre,
	crm_infra_handler.NewInvitationInstallService,
	crm_internal_usecase.NewInvitationInstallUsecase,
	crm_infra_mapper.NewLeadMapper,
	crm_infra_handler.NewLeadService,
	crm_internal_usecase.NewLeadUsecase,
	crm_infra_client.NewNotificationClient,
	crm_infra_client.NewOrganizationClient,
	crm_internal_usecase.NewOwnerUsecase,
	crm_infra_client.NewPermissionUsecase,
	crm_infra_handler.NewPipelineService,
	crm_internal_usecase.NewPipelineUsecase,
	crm_infra_postgre.NewPostgreCustomer,
	crm_infra_postgre.NewPostgrePipeline,
	crm_infra_postgre.NewPostgreRule,
	crm_infra_postgre.NewPostgreStage,
	crm_infra_postgre.NewProductCarePostgre,
	crm_infra_handler.NewRateHandler,
	crm_infra_mapper.NewRateMapper,
	crm_infra_postgre.NewRatePostgres,
	crm_internal_usecase.NewRateUsecase,
	crm_infra_handler.NewRegistryHandler,
	crm_infra_postgre.NewRegistryPostgres,
	crm_internal_usecase.NewRegistryUsecase,
	crm_infra_handler.NewReportAdminHandler,
	crm_infra_handler.NewReportHandler,
	crm_infra_mapper.NewReportMapper,
	crm_infra_postgre.NewReportPostgres,
	crm_infra_postgre.NewReportProofPostgres,
	crm_infra_handler.NewReportReasonHandler,
	crm_infra_postgre.NewReportReasonPostgres,
	crm_internal_usecase.NewReportReasonUsecase,
	crm_internal_usecase.NewReportUsecase,
	crm_internal_usecase.NewRuleEventUsecase,
	crm_infra_handler.NewRuleService,
	crm_internal_usecase.NewRuleUsecase,
	crm_infra_postgre.NewSharingAccessPostgre,
	crm_infra_handler.NewSharingAccessService,
	crm_internal_usecase.NewSharingAccessUsecase,
	crm_infra_handler.NewStageService,
	crm_internal_usecase.NewStageUsecase,
	crm_infra_mapper.NewTagMapper,
	crm_infra_postgre.NewTagPostgres,
	crm_infra_handler.NewTagService,
	crm_internal_usecase.NewTagUsecase,
	crm_infra_impl.NewTransactionGorm,
	crm_infra_client.NewUserClient,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}