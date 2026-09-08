//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	_provider "common/provider"
	common_utils "common/utils"
	crm_infra_client "crm/infra/client"
	crm_infra_email_sender "crm/infra/email_sender"
	crm_infra_handler "crm/infra/handler"
	crm_infra_handler_grpc "crm/infra/handler/grpc"
	crm_infra_handler_http_seopublic "crm/infra/handler/http/seopublic"
	crm_infra_impl "crm/infra/impl"
	crm_infra_job "crm/infra/job"
	crm_infra_job_appointment_reminder "crm/infra/job/appointment_reminder"
	crm_infra_mapper "crm/infra/mapper"
	crm_infra_postgre "crm/infra/postgre"
	crm_infra_postgre_contentcatalog "crm/infra/postgre/contentcatalog"
	crm_infra_rpc "crm/infra/rpc"
	crm_infra_validator "crm/infra/validator"
	crm_initial "crm/initial"
	crm_internal_domain_seo "crm/internal/domain/seo"
	crm_internal_integrations_paymentcompleted "crm/internal/integrations/paymentcompleted"
	crm_internal_interface_provider "crm/internal/interface/provider"
	crm_internal_modules_contentcatalog_application "crm/internal/modules/contentcatalog/application"
	crm_internal_repo "crm/internal/repo"
	crm_internal_usecase "crm/internal/usecase"
	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	OpenDatabase,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(crm_internal_interface_provider.BdsproProvider), new(*crm_infra_client.BdsproClient)),
	wire.Bind(new(crm_internal_interface_provider.ITransaction), new(*crm_infra_impl.TransactionGorm)),
	wire.Bind(new(crm_internal_interface_provider.TqdProvider), new(*crm_infra_client.TqdClient)),
	wire.Bind(new(crm_internal_interface_provider.UserClient), new(*crm_infra_client.UserClient)),
	wire.Bind(new(crm_internal_modules_contentcatalog_application.Repository), new(*crm_infra_postgre_contentcatalog.Repository)),
	wire.Bind(new(crm_internal_repo.BlockRepo), new(*crm_infra_postgre.BlockRepo)),
	wire.Bind(new(crm_internal_repo.ContactProductRepo), new(*crm_infra_postgre.ContactProductPostgre)),
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
	wire.Bind(new(crm_internal_usecase.PermissionUsecase), new(*crm_infra_client.PermissionUsecase)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	crm_infra_postgre.NewAdminOpportunityEventPostgres,
	crm_infra_handler.NewAdminSalesHandler,
	crm_internal_usecase.NewAdminSalesUsecase,
	crm_infra_postgre.NewAdvertisingSettingsPostgresRepository,
	crm_internal_usecase.NewAdvertisingSettingsUsecase,
	crm_infra_handler.NewAppointmentHandler,
	crm_infra_postgre.NewAppointmentPostgresRepository,
	crm_infra_job_appointment_reminder.NewAppointmentReminderJob,
	crm_infra_postgre.NewAppointmentReminderPostgresRepository,
	crm_internal_usecase.NewAppointmentReminderUsecase,
	crm_infra_mapper.NewAppointmentTransformer,
	crm_internal_usecase.NewAppointmentUsecase,
	crm_infra_validator.NewAppointmentValidator,
	crm_infra_client.NewAuthClient,
	crm_infra_rpc.NewAuthRPCClient,
	crm_infra_client.NewBdsproClient,
	crm_infra_rpc.NewBdsproRPCClient,
	crm_infra_postgre.NewBlockRepo,
	crm_infra_handler.NewBlockService,
	crm_internal_usecase.NewBlockUsecase,
	crm_infra_postgre.NewBudgetPostgresRepository,
	crm_internal_usecase.NewBudgetUsecase,
	crm_internal_usecase.NewCampaignDetailUsecase,
	crm_infra_postgre.NewCampaignPostgresRepository,
	crm_internal_usecase.NewCampaignUsecase,
	crm_infra_rpc.NewChatRPCClient,
	crm_internal_integrations_paymentcompleted.NewConsumer,
	crm_infra_mapper.NewContactMapper,
	crm_infra_postgre.NewContactProductPostgre,
	crm_infra_postgre.NewContactRepo,
	crm_infra_handler.NewContactService,
	crm_internal_usecase.NewContactUsecase,
	crm_infra_handler.NewCrmInternalService,
	crm_infra_job.NewDashboardStatsJob,
	crm_infra_postgre.NewDocumentPostgre,
	crm_infra_email_sender.NewEmailSender,
	crm_infra_handler.NewEnumService,
	crm_internal_usecase.NewEnumUsecase,
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
	crm_infra_handler_http_seopublic.NewHandler,
	crm_initial.NewInitialApp,
	crm_infra_postgre.NewInvitationInstallPostgre,
	crm_infra_handler.NewInvitationInstallService,
	crm_internal_usecase.NewInvitationInstallUsecase,
	crm_infra_mapper.NewLeadMapper,
	crm_infra_handler.NewLeadService,
	crm_internal_usecase.NewLeadUsecase,
	crm_infra_handler.NewMarketingGrpcHandler,
	crm_infra_client.NewNotificationClient,
	crm_infra_rpc.NewNotificationRPCClient,
	crm_infra_client.NewOrganizationClient,
	crm_infra_rpc.NewOrganizationRPCClient,
	crm_internal_usecase.NewOwnerUsecase,
	crm_infra_postgre.NewPackagePostgresRepository,
	crm_internal_usecase.NewPackageUsecase,
	crm_infra_client.NewPaymentClient,
	crm_infra_rpc.NewPaymentRPCClient,
	crm_internal_usecase.NewPaymentUsecase,
	crm_infra_client.NewPermissionUsecase,
	crm_infra_handler.NewPipelineService,
	crm_internal_usecase.NewPipelineUsecase,
	crm_infra_postgre.NewPostgreCustomer,
	crm_infra_postgre.NewPostgreOriginProfileRepo,
	crm_infra_postgre.NewPostgrePipeline,
	crm_infra_postgre.NewPostgreRule,
	crm_infra_postgre.NewPostgreStage,
	crm_infra_postgre.NewProductCarePostgre,
	crm_infra_handler_grpc.NewPublicContentGrpcHandler,
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
	crm_infra_postgre_contentcatalog.NewRepository,
	crm_internal_usecase.NewRuleEventUsecase,
	crm_infra_handler.NewRuleService,
	crm_internal_usecase.NewRuleUsecase,
	crm_infra_handler.NewSeoDomainHandler,
	crm_infra_mapper.NewSeoDomainMapper,
	crm_infra_postgre.NewSeoDomainPostgres,
	crm_internal_usecase.NewSeoDomainUsecase,
	crm_internal_domain_seo.NewSeoGenerationLog,
	crm_infra_postgre.NewSeoGenerationLogPostgres,
	crm_infra_postgre.NewSeoInternalLinkPostgres,
	crm_internal_usecase.NewSeoPublicUsecase,
	crm_infra_postgre.NewSeoRelativePostgres,
	crm_internal_usecase.NewSeoRenderUsecase,
	crm_internal_usecase.NewSeoSitemapUsecase,
	crm_internal_usecase.NewSeoWorkerUsecase,
	crm_internal_modules_contentcatalog_application.NewService,
	crm_infra_postgre.NewSharingAccessPostgre,
	crm_infra_handler.NewSharingAccessService,
	crm_internal_usecase.NewSharingAccessUsecase,
	crm_infra_handler_http_seopublic.NewSitemapHandler,
	crm_infra_handler.NewStageService,
	crm_internal_usecase.NewStageUsecase,
	crm_internal_integrations_paymentcompleted.NewStore,
	crm_infra_handler.NewSupportTicketHandler,
	crm_infra_postgre.NewSupportTicketPostgres,
	crm_internal_usecase.NewSupportTicketUsecase,
	crm_infra_rpc.NewTQDRPCClient,
	crm_infra_mapper.NewTagMapper,
	crm_infra_postgre.NewTagPostgres,
	crm_infra_handler.NewTagService,
	crm_internal_usecase.NewTagUsecase,
	crm_infra_client.NewTqdClient,
	crm_infra_impl.NewTransactionGorm,
	crm_infra_client.NewUserClient,
	crm_infra_rpc.NewUserRPCClient,
)

func InitializeApp() (*crm_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
