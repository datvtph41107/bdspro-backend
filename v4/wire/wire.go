//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	injection "common/injection"
	_provider "common/provider"
	common_utils "common/utils"
	"github.com/google/wire"
	user_infra_client "userv4/infra/client"
	user_infra_client_authpermission "userv4/infra/client/authpermission"
	user_infra_client_payment "userv4/infra/client/payment"
	user_infra_cookie "userv4/infra/cookie"
	user_infra_handler "userv4/infra/handler"
	user_infra_handler_grpc "userv4/infra/handler/grpc"
	user_infra_handler_grpc_entitlement "userv4/infra/handler/grpc/entitlement"
	user_modules_auth_mapper "userv4/modules/auth/mapper"
	user_modules_user_mapper "userv4/modules/user/mapper"
	user_modules_auth_oauth "userv4/modules/auth/oauth"
	user_infra_postgres "userv4/infra/postgres"
	user_modules_auth_postgres "userv4/modules/auth/postgres"
	user_modules_user_postgres "userv4/modules/user/postgres"
	user_modules_crm_postgres_maindb "userv4/modules/crm/postgres/maindb"
	user_modules_user_postgres_bootstrapadmin "userv4/modules/user/postgres/bootstrapadmin"
	user_modules_user_postgres_checkout "userv4/modules/user/postgres/checkout"
	user_modules_user_postgres_entitlement "userv4/modules/user/postgres/entitlement"
	user_modules_user_postgres_iampermission "userv4/modules/user/postgres/iampermission"
	user_modules_user_postgres_organization "userv4/modules/user/postgres/organization"
	user_modules_user_postgres_subscription "userv4/modules/user/postgres/subscription"
	user_infra_properties "userv4/infra/properties"
	user_infra_provider "userv4/infra/provider"
	user_infra_providers "userv4/infra/providers"
	user_infra_redis "userv4/infra/redis"
	user_infra_scheduler "userv4/infra/scheduler"
	user_infra_sms "userv4/infra/sms"
	user_infra_zns "userv4/infra/zns"
	user_initial "userv4/initial"
	user_internal_domain_user "userv4/modules/user/domain/user"
	user_internal_interface "userv4/internal/interface"
	user_internal_interface_factory "userv4/internal/interface/factory"
	user_internal_interface_providers "userv4/internal/interface/providers"
	user_internal_interface_repo "userv4/internal/interface/repo"
	user_internal_job "userv4/internal/job"
	user_internal_usecase "userv4/modules/auth/services"
	user_internal_usecase_bootstrapadmin "userv4/modules/user/services/bootstrapadmin"
	user_internal_usecase_commercialprofile "userv4/modules/user/services/commercialprofile"
	user_internal_usecase_entitlement "userv4/modules/user/services/entitlement"
	user_internal_usecase_organization "userv4/modules/user/services/organization"
	user_internal_usecase_plan_admin "userv4/modules/user/services/plan/admin"
	user_internal_usecase_plan_publish "userv4/modules/user/services/plan/publish"
	user_internal_usecase_subscription "userv4/modules/user/services/subscription"
	user_internal_usecase_subscription_cancel "userv4/modules/user/services/subscription/cancel"
	user_internal_usecase_subscription_checkout "userv4/modules/user/services/subscription/checkout"
	user_internal_usecase_subscription_settlement "userv4/modules/user/services/subscription/settlement"
	user_internal_usecase_user "userv4/modules/user/services/user"
	user_internal_usecase_useradmin "userv4/modules/user/services/useradmin"
	user_internal_usecases "userv4/modules/user/services"
	user_modules_auth_controller "userv4/modules/auth/controller"
	user_modules_auth_services "userv4/modules/auth/services"
	user_modules_bdspro_controller "userv4/modules/bdspro/controller"
	user_modules_bdspro_mapper "userv4/modules/bdspro/mapper"
	user_modules_bdspro_repository "userv4/modules/bdspro/repository"
	user_modules_bdspro_services "userv4/modules/bdspro/services"
	user_modules_crm "userv4/modules/crm"
	user_modules_hub_controller "userv4/modules/hub/controller"
	user_modules_hub_provider "userv4/modules/hub/provider"
	user_modules_seo_controller "userv4/modules/seo/controller"
	user_modules_seo_provider "userv4/modules/seo/provider"
	user_modules_seo_repository "userv4/modules/seo/repository"
	user_modules_seo_services "userv4/modules/seo/services"
	user_modules_hub_repository "userv4/modules/hub/repository"
	user_modules_hub_services "userv4/modules/hub/services"
	user_modules_notification_controller "userv4/modules/notification/controller"
	user_modules_notification_repository "userv4/modules/notification/repository"
	user_modules_notification_services "userv4/modules/notification/services"
	user_modules_payment "userv4/modules/payment"
	user_modules_tqd "userv4/modules/tqd"
	user_modules_tqd_provider "userv4/modules/tqd/provider"
	user_modules_user_controller "userv4/modules/user/controller"
	user_modules_user_services "userv4/modules/user/services"
	user_internal_requeststats "userv4/internal/requeststats"
	user_validator "userv4/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(user_infra_handler_grpc.OrganizationCheckoutAuthorizer), new(*user_modules_user_postgres_organization.CheckoutAuthorizer)),
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	wire.Bind(new(user_internal_interface_factory.IFactory), new(*user_infra_properties.RuntimeProperties)),
	wire.Bind(new(user_internal_interface_providers.BdsproProvider), new(*user_infra_client.BdsproClient)),
	wire.Bind(new(user_internal_interface_providers.HubProvider), new(*user_infra_client.HubClient)),
	wire.Bind(new(user_internal_interface_providers.OrganizationProvider), new(*user_infra_client.OrganizationProviderClient)),
	wire.Bind(new(user_internal_interface_providers.ProfileProvider), new(*user_infra_provider.ProfileProvider)),
	wire.Bind(new(user_internal_interface_repo.AdminAccessRepository), new(*user_modules_auth_postgres.AdminAccessPostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.AuthStatusRepository), new(*user_modules_auth_postgres.StatusPostgres)),
	wire.Bind(new(user_internal_interface_repo.DeviceRepository), new(*user_modules_auth_postgres.DevicePostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.IAuthConfigRepo), new(*user_modules_auth_postgres.AuthConfigPostgres)),
	wire.Bind(new(user_internal_interface_repo.IBlockRepo), new(*user_modules_crm_postgres_maindb.BlockPostgres)),
	wire.Bind(new(user_internal_interface_repo.IContactRepo), new(*user_modules_crm_postgres_maindb.ContactPostgres)),
	wire.Bind(new(user_internal_interface_repo.IDemoProductLogRepo), new(*user_modules_user_postgres.DemoProductLogPostgres)),
	wire.Bind(new(user_internal_interface_repo.IFollowRepo), new(*user_modules_crm_postgres_maindb.FollowPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileDeletedRepo), new(*user_modules_user_postgres.ProfileDeletedPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileRepo), new(*user_modules_user_postgres.ProfilePostgres)),
	wire.Bind(new(user_internal_interface_repo.OTPRepository), new(*user_modules_auth_postgres.OTPPostgres)),
	wire.Bind(new(user_internal_interface_repo.PINRepository), new(*user_modules_auth_postgres.PINPostgres)),
	wire.Bind(new(user_internal_interface_repo.PermissionRepository), new(*user_modules_auth_postgres.PermissionRepo)),
	wire.Bind(new(user_internal_interface_repo.SessionRepository), new(*user_modules_auth_postgres.SessionPostgres)),
	wire.Bind(new(user_internal_usecase.AccessTokenIssuer), new(*user_internal_usecase.AuthUsecase)),
	wire.Bind(new(user_internal_usecase_commercialprofile.CatalogReader), new(*user_internal_usecase_plan_admin.Service)),
	wire.Bind(new(user_internal_usecase_commercialprofile.SubscriptionReader), new(*user_internal_usecase_subscription.AdminService)),
	wire.Bind(new(user_internal_usecase_entitlement.PlanAccessStore), new(*user_modules_user_postgres_entitlement.PlanAccessStore)),
	wire.Bind(new(user_internal_usecase_entitlement.SubscriptionStore), new(*user_modules_user_postgres_entitlement.SubscriptionStore)),
	wire.Bind(new(user_internal_usecase_organization.Repository), new(*user_modules_user_postgres_organization.Repository)),
	wire.Bind(new(user_internal_usecase_plan_admin.PermissionAuthorizer), new(*user_infra_client_authpermission.CatalogAuthorizer)),
	wire.Bind(new(user_internal_usecase_plan_admin.Publisher), new(*user_internal_usecase_plan_publish.Service)),
	wire.Bind(new(user_internal_usecase_plan_admin.Repository), new(*user_modules_user_postgres.PlanVersionRepo)),
	wire.Bind(new(user_internal_usecase_plan_publish.Repository), new(*user_modules_user_postgres.PlanVersionRepo)),
	wire.Bind(new(user_internal_usecase_subscription.AdminRepository), new(*user_modules_user_postgres_subscription.AdminProjectionStore)),
	wire.Bind(new(user_internal_usecase_subscription.PermissionAuthorizer), new(*user_modules_user_postgres_iampermission.Checker)),
	wire.Bind(new(user_internal_usecase_subscription_cancel.Store), new(*user_modules_user_postgres_subscription.CancellationStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.CommandStore), new(*user_modules_user_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.ContactStore), new(*user_modules_user_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentAttemptPort), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentMethodReader), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentOrderPort), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PlanResolver), new(*user_modules_user_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.SubscriptionReader), new(*user_modules_user_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_settlement.Store), new(*user_modules_user_postgres_subscription.SettlementStore)),
	wire.Bind(new(user_internal_usecase_useradmin.PermissionAuthorizer), new(*user_infra_client_authpermission.CatalogAuthorizer)),
	injection.RPCClientSet,
	injection.NewHttpClient,
	common_utils.NewSyncUtil,
	user_modules_auth_postgres.NewAdminAccessPostgresRepo,
	user_internal_usecase.NewAdminAccessUsecase,
	user_infra_handler.NewAdminHandler,
	user_internal_usecase_user.NewAdminKYCUsecase,
	user_internal_usecases.NewAdminMainAreaUsecase,
	user_modules_user_mapper.NewAdminMapper,
	user_infra_handler_grpc.NewAdminOrganizationHandler,
	user_modules_user_postgres.NewAdminPostgres,
	user_modules_user_postgres_subscription.NewAdminProjectionStore,
	user_internal_usecase_user.NewAdminPurposeUseUsecase,
	user_internal_usecase_subscription.NewAdminService,
	user_internal_usecase.NewAdminUsecase,
	user_internal_usecase_user.NewAdminUsecase,
	user_infra_handler.NewAdminUserProfileHandler,
	user_infra_client.NewAuthClient,
	user_modules_auth_postgres.NewAuthConfigRepository,
	user_infra_handler.NewAuthHandler,
	user_modules_auth_controller.NewLoginController,
	user_modules_auth_services.NewAuthUsecaseLoginService,
	wire.Bind(new(user_modules_auth_services.LoginService), new(*user_modules_auth_services.AuthUsecaseLoginService)),
	user_modules_auth_controller.NewPermissionController,
	user_modules_auth_services.NewRoleUsecasePermissionService,
	wire.Bind(new(user_modules_auth_services.PermissionService), new(*user_modules_auth_services.RoleUsecasePermissionService)),
	user_modules_bdspro_controller.NewHandler,
	user_modules_bdspro_services.NewLocalSyncService,
	wire.Bind(new(user_modules_bdspro_services.SyncService), new(*user_modules_bdspro_services.LocalSyncService)),
	user_modules_bdspro_repository.NewDealPostgres,
	user_modules_bdspro_services.NewLocalDealService,
	wire.Bind(new(user_modules_bdspro_services.DealService), new(*user_modules_bdspro_services.LocalDealService)),
	user_modules_bdspro_repository.NewDealMemberPostgres,
	user_modules_bdspro_services.NewLocalDealInvitationService,
	wire.Bind(new(user_modules_bdspro_services.DealInvitationService), new(*user_modules_bdspro_services.LocalDealInvitationService)),
	user_modules_bdspro_repository.NewCostTypePostgres,
	user_modules_bdspro_services.NewLocalCostTypeService,
	wire.Bind(new(user_modules_bdspro_services.CostTypeService), new(*user_modules_bdspro_services.LocalCostTypeService)),
	user_modules_bdspro_repository.NewDealCostPostgres,
	user_modules_bdspro_services.NewLocalDealCostService,
	wire.Bind(new(user_modules_bdspro_services.DealCostService), new(*user_modules_bdspro_services.LocalDealCostService)),
	user_modules_bdspro_repository.NewContractDealPostgres,
	user_modules_bdspro_services.NewLocalContractDealService,
	wire.Bind(new(user_modules_bdspro_services.ContractDealService), new(*user_modules_bdspro_services.LocalContractDealService)),
	user_modules_bdspro_repository.NewInvestmentPostgres,
	user_modules_bdspro_services.NewLocalInvestmentService,
	wire.Bind(new(user_modules_bdspro_services.InvestmentService), new(*user_modules_bdspro_services.LocalInvestmentService)),
	user_modules_bdspro_repository.NewProductPostgres,
	user_modules_bdspro_services.NewLocalProductService,
	wire.Bind(new(user_modules_bdspro_services.ProductService), new(*user_modules_bdspro_services.LocalProductService)),
	user_modules_bdspro_repository.NewDistributePostgres,
	user_modules_bdspro_repository.NewProductPricePostgres,
	user_modules_bdspro_mapper.NewDistributionMapper,
	user_modules_bdspro_services.NewLocalDistributionService,
	wire.Bind(new(user_modules_bdspro_services.DistributionService), new(*user_modules_bdspro_services.LocalDistributionService)),
	user_modules_bdspro_repository.NewPropertyPostgres,
	user_modules_bdspro_services.NewLocalPropertyService,
	wire.Bind(new(user_modules_bdspro_services.PropertyService), new(*user_modules_bdspro_services.LocalPropertyService)),
	user_modules_bdspro_services.NewLocalPublicService,
	wire.Bind(new(user_modules_bdspro_services.PublicService), new(*user_modules_bdspro_services.LocalPublicService)),
	user_modules_bdspro_repository.NewCatalogPostgres,
	user_modules_bdspro_services.NewLocalCatalogService,
	wire.Bind(new(user_modules_bdspro_services.CatalogService), new(*user_modules_bdspro_services.LocalCatalogService)),
	user_modules_bdspro_repository.NewProjectBuildPostgres,
	user_modules_bdspro_services.NewLocalProjectBuildService,
	wire.Bind(new(user_modules_bdspro_services.ProjectBuildService), new(*user_modules_bdspro_services.LocalProjectBuildService)),
	user_modules_bdspro_repository.NewProjectAdminPostgres,
	user_modules_bdspro_services.NewLocalAdminProjectService,
	wire.Bind(new(user_modules_bdspro_services.AdminProjectService), new(*user_modules_bdspro_services.LocalAdminProjectService)),
	user_modules_bdspro_repository.NewProductNotePostgres,
	user_modules_bdspro_services.NewLocalProductNoteService,
	wire.Bind(new(user_modules_bdspro_services.ProductNoteService), new(*user_modules_bdspro_services.LocalProductNoteService)),
	user_modules_bdspro_repository.NewAppointmentPostgres,
	user_modules_bdspro_services.NewLocalAppointmentService,
	wire.Bind(new(user_modules_bdspro_services.AppointmentService), new(*user_modules_bdspro_services.LocalAppointmentService)),
	user_modules_bdspro_repository.NewPropertyHistoryPostgres,
	user_modules_bdspro_services.NewLocalPropertyHistoryService,
	wire.Bind(new(user_modules_bdspro_services.PropertyHistoryService), new(*user_modules_bdspro_services.LocalPropertyHistoryService)),
	user_modules_hub_controller.NewHandler,
	user_modules_hub_repository.NewLocationPostgres,
	user_modules_hub_services.NewLocalLocationService,
	wire.Bind(new(user_modules_hub_services.LocationService), new(*user_modules_hub_services.LocalLocationService)),
	user_modules_hub_repository.NewLocationV2Postgres,
	user_modules_hub_services.NewLocalLocationV2Service,
	wire.Bind(new(user_modules_hub_services.LocationV2Service), new(*user_modules_hub_services.LocalLocationV2Service)),
	user_modules_hub_repository.NewUserGuidePostgres,
	user_modules_hub_services.NewLocalUserGuideService,
	wire.Bind(new(user_modules_hub_services.UserGuideService), new(*user_modules_hub_services.LocalUserGuideService)),
	user_modules_hub_repository.NewFAQPostgres,
	user_modules_hub_services.NewLocalFAQService,
	wire.Bind(new(user_modules_hub_services.FAQService), new(*user_modules_hub_services.LocalFAQService)),
	user_modules_hub_repository.NewVersionPostgres,
	user_modules_hub_services.NewLocalVersionService,
	wire.Bind(new(user_modules_hub_services.VersionService), new(*user_modules_hub_services.LocalVersionService)),
	user_internal_requeststats.NewStore,
	user_internal_requeststats.NewCollector,
	user_modules_hub_services.NewLocalRequestStatsService,
	wire.Bind(new(user_modules_hub_services.RequestStatsService), new(*user_modules_hub_services.LocalRequestStatsService)),
	user_modules_hub_repository.NewErrorLogPostgres,
	user_modules_hub_services.NewLocalErrorLogService,
	wire.Bind(new(user_modules_hub_services.ErrorLogService), new(*user_modules_hub_services.LocalErrorLogService)),
	user_modules_hub_repository.NewApplinkPostgres,
	user_modules_hub_services.NewLocalApplinkService,
	wire.Bind(new(user_modules_hub_services.ApplinkService), new(*user_modules_hub_services.LocalApplinkService)),
	user_modules_hub_repository.NewSystemConfigPostgres,
	user_modules_hub_services.NewLocalSystemConfigService,
	wire.Bind(new(user_modules_hub_services.SystemConfigService), new(*user_modules_hub_services.LocalSystemConfigService)),
	user_modules_hub_repository.NewEventQueuePostgres,
	user_modules_hub_services.NewLocalEventQueueService,
	wire.Bind(new(user_modules_hub_services.EventQueueService), new(*user_modules_hub_services.LocalEventQueueService)),
	user_modules_hub_repository.NewInteractiveEventPostgres,
	user_modules_hub_services.NewLocalInteractiveEventService,
	wire.Bind(new(user_modules_hub_services.InteractiveEventService), new(*user_modules_hub_services.LocalInteractiveEventService)),
	user_modules_hub_repository.NewClickURLPostgres,
	user_modules_hub_services.NewLocalClickURLService,
	wire.Bind(new(user_modules_hub_services.ClickURLService), new(*user_modules_hub_services.LocalClickURLService)),
	user_modules_hub_repository.NewUpdateDataPostgres,
	user_modules_hub_services.NewLocalUpdateDataService,
	wire.Bind(new(user_modules_hub_services.UpdateDataService), new(*user_modules_hub_services.LocalUpdateDataService)),
	user_modules_hub_repository.NewKeywordSearchPostgres,
	user_modules_hub_services.NewLocalKeywordSearchService,
	wire.Bind(new(user_modules_hub_services.KeywordSearchService), new(*user_modules_hub_services.LocalKeywordSearchService)),
	user_modules_hub_repository.NewSearchLatencyPostgres,
	user_modules_hub_services.NewLocalSearchLatencyService,
	wire.Bind(new(user_modules_hub_services.SearchLatencyService), new(*user_modules_hub_services.LocalSearchLatencyService)),
	user_modules_hub_provider.NewLocalSearchAnalytics,
	wire.Bind(new(user_modules_tqd_provider.SearchAnalytics), new(*user_modules_hub_provider.LocalSearchAnalytics)),
	user_modules_hub_repository.NewAPIKeyPostgres,
	user_modules_hub_services.NewLocalAPIKeyService,
	wire.Bind(new(user_modules_hub_services.APIKeyService), new(*user_modules_hub_services.LocalAPIKeyService)),
	user_infra_postgres.NewTqdDB,
	user_modules_tqd.ProviderSet,
	user_infra_postgres.NewCrmDB,
	user_modules_crm.ProviderSet,
	user_infra_postgres.NewPaymentDB,
	user_modules_payment.ProviderSet,
	// modules/seo
	user_modules_seo_controller.NewHandler,
	user_modules_seo_repository.NewNewsPostgres,
	user_modules_seo_services.NewLocalNewsService,
	wire.Bind(new(user_modules_seo_services.NewsService), new(*user_modules_seo_services.LocalNewsService)),
	user_modules_seo_repository.NewSeoDomainPostgres,
	user_modules_seo_services.NewLocalSeoDomainService,
	wire.Bind(new(user_modules_seo_services.SeoDomainService), new(*user_modules_seo_services.LocalSeoDomainService)),
	user_modules_seo_provider.NewLocalCrmSeoProvider,
	wire.Bind(new(user_modules_tqd_provider.CrmSeoProvider), new(*user_modules_seo_provider.LocalCrmSeoProvider)),
	user_modules_notification_controller.NewHandler,
	user_modules_notification_repository.NewNotificationPostgres,
	user_modules_notification_services.NewLocalNotificationService,
	wire.Bind(new(user_modules_notification_services.NotificationService), new(*user_modules_notification_services.LocalNotificationService)),
	user_modules_notification_services.NewLocalNotificationProvider,
	wire.Bind(new(user_internal_interface_providers.NotificationProvider), new(*user_modules_notification_services.LocalNotificationProvider)),
	user_modules_user_controller.NewHandler,
	user_modules_user_controller.NewProfileController,
	user_modules_user_controller.NewCatalogController,
	user_modules_user_controller.NewTagController,
	user_modules_user_controller.NewKYCController,
	user_modules_user_controller.NewBookmarkController,
	user_modules_user_controller.NewAdminUserController,
	user_modules_user_services.NewProfileUsecaseInfoService,
	wire.Bind(new(user_modules_user_services.ProfileService), new(*user_modules_user_services.ProfileUsecaseInfoService)),
	user_modules_user_services.NewLocalDemoService,
	wire.Bind(new(user_modules_user_services.DemoService), new(*user_modules_user_services.LocalDemoService)),
	user_modules_user_services.NewLocalTileSessionService,
	wire.Bind(new(user_modules_user_services.TileSessionService), new(*user_modules_user_services.LocalTileSessionService)),
	user_modules_user_services.NewLocalCatalogService,
	wire.Bind(new(user_modules_user_services.CatalogService), new(*user_modules_user_services.LocalCatalogService)),
	user_modules_user_services.NewLocalTagService,
	wire.Bind(new(user_modules_user_services.TagService), new(*user_modules_user_services.LocalTagService)),
	user_modules_user_services.NewLocalKYCService,
	wire.Bind(new(user_modules_user_services.KYCService), new(*user_modules_user_services.LocalKYCService)),
	user_modules_user_services.NewLocalBookmarkService,
	wire.Bind(new(user_modules_user_services.BookmarkService), new(*user_modules_user_services.LocalBookmarkService)),
	user_modules_user_services.NewLocalAdminUserService,
	wire.Bind(new(user_modules_user_services.AdminUserService), new(*user_modules_user_services.LocalAdminUserService)),
	user_modules_auth_mapper.NewAuthMethodMapper,
	user_modules_auth_postgres.NewAuthRepository,
	user_modules_auth_mapper.NewAuthSecurityMapper,
	user_internal_usecase.NewAuthSecurityUsecase,
	user_internal_usecase.NewAuthService,
	user_infra_client.NewBdsproClient,
	user_modules_crm_postgres_maindb.NewBlockPostgres,
	user_infra_handler.NewBookmarkUserHandler,
	user_modules_user_postgres.NewBookmarkUserPostgres,
	user_internal_usecase_user.NewBookmarkUserUsecase,
	user_infra_providers.NewCancelClock,
	user_modules_user_postgres_subscription.NewCancellationStore,
	user_infra_client_authpermission.NewCatalogAuthorizer,
	user_modules_user_postgres.NewCertificationPostgres,
	user_modules_user_postgres_iampermission.NewChecker,
	user_modules_user_postgres_organization.NewCheckoutAuthorizer,
	user_infra_providers.NewCheckoutClock,
	user_modules_user_postgres_checkout.NewCheckoutStore,
	user_infra_client_payment.NewCommands,
	user_infra_handler_grpc.NewCommercialProfileHandler,
	user_modules_crm_postgres_maindb.NewContactPostgres,
	user_infra_cookie.NewCookieProvider,
	user_modules_user_postgres.NewDemoProductLogPostgres,
	user_modules_auth_mapper.NewDeviceMapper,
	user_modules_auth_postgres.NewDevicePostgresRepo,
	user_internal_usecase.NewDeviceUsecase,
	user_modules_auth_oauth.NewFacebookAuthService,
	user_modules_crm_postgres_maindb.NewFollowPostgres,
	user_modules_crm_postgres_maindb.NewFriendPostgres,
	user_modules_auth_oauth.NewGoogleAuthService,
	user_modules_crm_postgres_maindb.NewGroupPostgres,
	user_infra_handler.NewGrpcProfileService,
	user_infra_handler_grpc_entitlement.NewHandler,
	user_infra_client.NewHubClient,
	user_initial.NewInitialApp,
	user_infra_handler.NewInternalHandler,
	user_infra_handler_grpc.NewInternalOrganizationMembershipHandler,
	user_infra_handler.NewKYCHandler,
	user_modules_user_mapper.NewKYCMapper,
	user_modules_user_postgres.NewKYCRepo,
	user_internal_usecase_user.NewKYCUsecase,
	user_infra_handler.NewMainAreaHandler,
	user_modules_user_mapper.NewMainAreaMapper,
	user_modules_user_postgres.NewMainAreaPostgres,
	user_infra_handler.NewOAuthRouter,
	user_internal_usecase.NewOAuthUsecase,
	user_modules_auth_postgres.NewOTPRepository,
	user_infra_client.NewOrganizationProviderClient,
	user_internal_usecase.NewOtpUsecase,
	user_modules_auth_postgres.NewPINRepository,
	user_internal_usecase.NewPINUsecase,
	user_infra_handler.NewPermissionHandler,
	user_modules_auth_postgres.NewPermissionRepo,
	user_internal_usecase.NewPermissionUsecase,
	user_modules_user_postgres_entitlement.NewPlanAccessStore,
	user_infra_providers.NewPlanPublishTransaction,
	user_infra_handler_grpc.NewPlanVersionAdminHandler,
	user_modules_user_postgres.NewPlanVersionRepo,
	user_infra_handler.NewPriceTableHandler,
	user_modules_user_mapper.NewPriceTableMapper,
	user_modules_user_postgres.NewPriceTableRepo,
	user_internal_usecases.NewPriceTableUsecase,
	user_infra_handler.NewProfessionHandler,
	user_modules_user_mapper.NewProfessionMapper,
	user_modules_user_postgres.NewProfessionPostgres,
	user_internal_usecase_user.NewProfessionUsecase,
	user_internal_domain_user.NewProfileDeleted,
	user_modules_user_postgres.NewProfileDeletedPostgres,
	user_modules_user_mapper.NewProfileMapper,
	user_modules_user_postgres.NewProfileMediaPostgres,
	user_modules_user_postgres.NewProfilePostgres,
	user_infra_provider.NewProfileProvider,
	user_internal_usecase_user.NewProfileUsecase,
	user_infra_handler.NewPurposeUseHandler,
	user_modules_user_mapper.NewPurposeUseMapper,
	user_modules_user_postgres.NewPurposeUseRepo,
	user_infra_redis.NewRedisProvider,
	user_modules_user_postgres_organization.NewRepository,
	user_infra_handler.NewRoleGroupHandler,
	user_modules_auth_mapper.NewRoleGroupMapper,
	user_internal_usecase.NewRoleGroupRegistry,
	user_modules_auth_postgres.NewRoleGroupRepo,
	user_internal_usecase.NewRoleGroupUsecase,
	user_validator.NewRoleGroupValidator,
	user_infra_handler.NewRoleHandler,
	user_modules_auth_mapper.NewRoleMapper,
	user_modules_auth_postgres.NewRoleRepo,
	user_internal_usecase.NewRoleUsecase,
	user_infra_properties.NewRuntimeProperties,
	user_infra_sms.NewSMSProvider,
	user_internal_usecase_bootstrapadmin.NewService,
	user_internal_usecase_commercialprofile.NewService,
	user_internal_usecase_entitlement.NewService,
	user_internal_usecase_organization.NewService,
	user_internal_usecase_plan_admin.NewService,
	user_internal_usecase_plan_publish.NewService,
	user_internal_usecase_subscription_cancel.NewService,
	user_internal_usecase_subscription_checkout.NewService,
	user_internal_usecase_subscription_settlement.NewService,
	user_modules_auth_mapper.NewSessionMapper,
	user_modules_auth_postgres.NewSessionRepository,
	user_modules_user_postgres_subscription.NewSettlementStore,
	user_modules_auth_postgres.NewStatusRepository,
	user_modules_user_postgres_bootstrapadmin.NewStore,
	user_infra_handler_grpc.NewSubscriptionAdminHandler,
	user_infra_handler_grpc.NewSubscriptionCheckoutHandler,
	user_infra_handler_grpc.NewSubscriptionSettlementHandler,
	user_modules_user_postgres_entitlement.NewSubscriptionStore,
	user_infra_handler.NewTagHandler,
	user_modules_user_mapper.NewTagMapper,
	user_modules_user_postgres.NewTagPostgres,
	user_internal_usecase_user.NewTagUsecase,
	user_infra_postgres.NewTransactionGorm,
	user_internal_job.NewUserDashboardStatsJob,
	user_infra_handler.NewUserInfoHandler,
	user_modules_auth_mapper.NewUserInfoMapper,
	user_modules_auth_postgres.NewUserInfoRepository,
	user_internal_usecase.NewUserInfoUsecase,
	user_internal_domain_user.NewUserProfileEntityWithRequest,
	user_internal_usecase_user.NewUserUsecase,
	user_modules_auth_oauth.NewZaloAuthService,
	user_infra_zns.NewZnsProvider,
	user_infra_scheduler.NewZnsScheduler,
)

func InitializeApp() (*user_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
