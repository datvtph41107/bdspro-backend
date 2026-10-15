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
	user_infra_mapper "userv4/infra/mapper"
	user_infra_oauth "userv4/infra/oauth"
	user_infra_postgres "userv4/infra/postgres"
	user_infra_postgres_bootstrapadmin "userv4/infra/postgres/bootstrapadmin"
	user_infra_postgres_checkout "userv4/infra/postgres/checkout"
	user_infra_postgres_entitlement "userv4/infra/postgres/entitlement"
	user_infra_postgres_iampermission "userv4/infra/postgres/iampermission"
	user_infra_postgres_organization "userv4/infra/postgres/organization"
	user_infra_postgres_subscription "userv4/infra/postgres/subscription"
	user_infra_properties "userv4/infra/properties"
	user_infra_provider "userv4/infra/provider"
	user_infra_providers "userv4/infra/providers"
	user_infra_redis "userv4/infra/redis"
	user_infra_scheduler "userv4/infra/scheduler"
	user_infra_sms "userv4/infra/sms"
	user_infra_zns "userv4/infra/zns"
	user_initial "userv4/initial"
	user_internal_domain_user "userv4/internal/domain/user"
	user_internal_interface "userv4/internal/interface"
	user_internal_interface_factory "userv4/internal/interface/factory"
	user_internal_interface_providers "userv4/internal/interface/providers"
	user_internal_interface_repo "userv4/internal/interface/repo"
	user_internal_job "userv4/internal/job"
	user_internal_usecase "userv4/internal/usecase"
	user_internal_usecase_bootstrapadmin "userv4/internal/usecase/bootstrapadmin"
	user_internal_usecase_commercialprofile "userv4/internal/usecase/commercialprofile"
	user_internal_usecase_entitlement "userv4/internal/usecase/entitlement"
	user_internal_usecase_organization "userv4/internal/usecase/organization"
	user_internal_usecase_plan_admin "userv4/internal/usecase/plan/admin"
	user_internal_usecase_plan_publish "userv4/internal/usecase/plan/publish"
	user_internal_usecase_subscription "userv4/internal/usecase/subscription"
	user_internal_usecase_subscription_cancel "userv4/internal/usecase/subscription/cancel"
	user_internal_usecase_subscription_checkout "userv4/internal/usecase/subscription/checkout"
	user_internal_usecase_subscription_settlement "userv4/internal/usecase/subscription/settlement"
	user_internal_usecase_user "userv4/internal/usecase/user"
	user_internal_usecase_useradmin "userv4/internal/usecase/useradmin"
	user_internal_usecases "userv4/internal/usecases"
	user_modules_auth_controller "userv4/modules/auth/controller"
	user_modules_auth_services "userv4/modules/auth/services"
	user_modules_bdspro_controller "userv4/modules/bdspro/controller"
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
	wire.Bind(new(user_infra_handler_grpc.OrganizationCheckoutAuthorizer), new(*user_infra_postgres_organization.CheckoutAuthorizer)),
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	wire.Bind(new(user_internal_interface_factory.IFactory), new(*user_infra_properties.RuntimeProperties)),
	wire.Bind(new(user_internal_interface_providers.BdsproProvider), new(*user_infra_client.BdsproClient)),
	wire.Bind(new(user_internal_interface_providers.HubProvider), new(*user_infra_client.HubClient)),
	wire.Bind(new(user_internal_interface_providers.OrganizationProvider), new(*user_infra_client.OrganizationProviderClient)),
	wire.Bind(new(user_internal_interface_providers.ProfileProvider), new(*user_infra_provider.ProfileProvider)),
	wire.Bind(new(user_internal_interface_repo.AdminAccessRepository), new(*user_infra_postgres.AdminAccessPostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.AuthStatusRepository), new(*user_infra_postgres.StatusPostgres)),
	wire.Bind(new(user_internal_interface_repo.DeviceRepository), new(*user_infra_postgres.DevicePostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.IAuthConfigRepo), new(*user_infra_postgres.AuthConfigPostgres)),
	wire.Bind(new(user_internal_interface_repo.IBlockRepo), new(*user_infra_postgres.BlockPostgres)),
	wire.Bind(new(user_internal_interface_repo.IContactRepo), new(*user_infra_postgres.ContactPostgres)),
	wire.Bind(new(user_internal_interface_repo.IDemoProductLogRepo), new(*user_infra_postgres.DemoProductLogPostgres)),
	wire.Bind(new(user_internal_interface_repo.IFollowRepo), new(*user_infra_postgres.FollowPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileDeletedRepo), new(*user_infra_postgres.ProfileDeletedPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileRepo), new(*user_infra_postgres.ProfilePostgres)),
	wire.Bind(new(user_internal_interface_repo.OTPRepository), new(*user_infra_postgres.OTPPostgres)),
	wire.Bind(new(user_internal_interface_repo.PINRepository), new(*user_infra_postgres.PINPostgres)),
	wire.Bind(new(user_internal_interface_repo.PermissionRepository), new(*user_infra_postgres.PermissionRepo)),
	wire.Bind(new(user_internal_interface_repo.SessionRepository), new(*user_infra_postgres.SessionPostgres)),
	wire.Bind(new(user_internal_usecase.AccessTokenIssuer), new(*user_internal_usecase.AuthUsecase)),
	wire.Bind(new(user_internal_usecase_commercialprofile.CatalogReader), new(*user_internal_usecase_plan_admin.Service)),
	wire.Bind(new(user_internal_usecase_commercialprofile.SubscriptionReader), new(*user_internal_usecase_subscription.AdminService)),
	wire.Bind(new(user_internal_usecase_entitlement.PlanAccessStore), new(*user_infra_postgres_entitlement.PlanAccessStore)),
	wire.Bind(new(user_internal_usecase_entitlement.SubscriptionStore), new(*user_infra_postgres_entitlement.SubscriptionStore)),
	wire.Bind(new(user_internal_usecase_organization.Repository), new(*user_infra_postgres_organization.Repository)),
	wire.Bind(new(user_internal_usecase_plan_admin.PermissionAuthorizer), new(*user_infra_client_authpermission.CatalogAuthorizer)),
	wire.Bind(new(user_internal_usecase_plan_admin.Publisher), new(*user_internal_usecase_plan_publish.Service)),
	wire.Bind(new(user_internal_usecase_plan_admin.Repository), new(*user_infra_postgres.PlanVersionRepo)),
	wire.Bind(new(user_internal_usecase_plan_publish.Repository), new(*user_infra_postgres.PlanVersionRepo)),
	wire.Bind(new(user_internal_usecase_subscription.AdminRepository), new(*user_infra_postgres_subscription.AdminProjectionStore)),
	wire.Bind(new(user_internal_usecase_subscription.PermissionAuthorizer), new(*user_infra_postgres_iampermission.Checker)),
	wire.Bind(new(user_internal_usecase_subscription_cancel.Store), new(*user_infra_postgres_subscription.CancellationStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.CommandStore), new(*user_infra_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.ContactStore), new(*user_infra_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentAttemptPort), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentMethodReader), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PaymentOrderPort), new(*user_infra_client_payment.Commands)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.PlanResolver), new(*user_infra_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_checkout.SubscriptionReader), new(*user_infra_postgres_checkout.CheckoutStore)),
	wire.Bind(new(user_internal_usecase_subscription_settlement.Store), new(*user_infra_postgres_subscription.SettlementStore)),
	wire.Bind(new(user_internal_usecase_useradmin.PermissionAuthorizer), new(*user_infra_client_authpermission.CatalogAuthorizer)),
	injection.RPCClientSet,
	injection.NewHttpClient,
	common_utils.NewSyncUtil,
	user_infra_postgres.NewAdminAccessPostgresRepo,
	user_internal_usecase.NewAdminAccessUsecase,
	user_infra_handler.NewAdminHandler,
	user_internal_usecase_user.NewAdminKYCUsecase,
	user_internal_usecases.NewAdminMainAreaUsecase,
	user_infra_mapper.NewAdminMapper,
	user_infra_handler_grpc.NewAdminOrganizationHandler,
	user_infra_postgres.NewAdminPostgres,
	user_infra_postgres_subscription.NewAdminProjectionStore,
	user_internal_usecase_user.NewAdminPurposeUseUsecase,
	user_internal_usecase_subscription.NewAdminService,
	user_internal_usecase.NewAdminUsecase,
	user_internal_usecase_user.NewAdminUsecase,
	user_infra_handler.NewAdminUserProfileHandler,
	user_infra_client.NewAuthClient,
	user_infra_postgres.NewAuthConfigRepository,
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
	user_infra_mapper.NewAuthMethodMapper,
	user_infra_postgres.NewAuthRepository,
	user_infra_mapper.NewAuthSecurityMapper,
	user_internal_usecase.NewAuthSecurityUsecase,
	user_internal_usecase.NewAuthService,
	user_infra_client.NewBdsproClient,
	user_infra_postgres.NewBlockPostgres,
	user_internal_usecases.NewBlockUsecase,
	user_infra_handler.NewBookmarkUserHandler,
	user_infra_postgres.NewBookmarkUserPostgres,
	user_internal_usecase_user.NewBookmarkUserUsecase,
	user_infra_providers.NewCancelClock,
	user_infra_postgres_subscription.NewCancellationStore,
	user_infra_client_authpermission.NewCatalogAuthorizer,
	user_infra_postgres.NewCertificationPostgres,
	user_internal_usecases.NewCertificationUsecase,
	user_infra_postgres_iampermission.NewChecker,
	user_infra_postgres_organization.NewCheckoutAuthorizer,
	user_infra_providers.NewCheckoutClock,
	user_infra_postgres_checkout.NewCheckoutStore,
	user_infra_client_payment.NewCommands,
	user_infra_handler_grpc.NewCommercialProfileHandler,
	user_infra_postgres.NewContactPostgres,
	user_infra_cookie.NewCookieProvider,
	user_infra_postgres.NewDemoProductLogPostgres,
	user_infra_mapper.NewDeviceMapper,
	user_infra_postgres.NewDevicePostgresRepo,
	user_internal_usecase.NewDeviceUsecase,
	user_infra_oauth.NewFacebookAuthService,
	user_infra_postgres.NewFollowPostgres,
	user_internal_usecases.NewFollowUsecase,
	user_infra_postgres.NewFriendPostgres,
	user_internal_usecases.NewFriendUsecase,
	user_infra_oauth.NewGoogleAuthService,
	user_infra_postgres.NewGroupPostgres,
	user_infra_handler.NewGrpcProfileService,
	user_infra_handler_grpc_entitlement.NewHandler,
	user_infra_client.NewHubClient,
	user_initial.NewInitialApp,
	user_infra_handler.NewInternalHandler,
	user_infra_handler_grpc.NewInternalOrganizationMembershipHandler,
	user_infra_handler.NewKYCHandler,
	user_infra_mapper.NewKYCMapper,
	user_infra_postgres.NewKYCRepo,
	user_internal_usecase_user.NewKYCUsecase,
	user_infra_handler.NewMainAreaHandler,
	user_infra_mapper.NewMainAreaMapper,
	user_infra_postgres.NewMainAreaPostgres,
	user_infra_handler.NewOAuthRouter,
	user_internal_usecase.NewOAuthUsecase,
	user_infra_postgres.NewOTPRepository,
	user_infra_client.NewOrganizationProviderClient,
	user_internal_usecase.NewOtpUsecase,
	user_infra_postgres.NewPINRepository,
	user_internal_usecase.NewPINUsecase,
	user_infra_handler.NewPermissionHandler,
	user_infra_postgres.NewPermissionRepo,
	user_internal_usecase.NewPermissionUsecase,
	user_infra_postgres_entitlement.NewPlanAccessStore,
	user_infra_providers.NewPlanPublishTransaction,
	user_infra_handler_grpc.NewPlanVersionAdminHandler,
	user_infra_postgres.NewPlanVersionRepo,
	user_infra_handler.NewPriceTableHandler,
	user_infra_mapper.NewPriceTableMapper,
	user_infra_postgres.NewPriceTableRepo,
	user_internal_usecases.NewPriceTableUsecase,
	user_infra_handler.NewProfessionHandler,
	user_infra_mapper.NewProfessionMapper,
	user_infra_postgres.NewProfessionPostgres,
	user_internal_usecase_user.NewProfessionUsecase,
	user_internal_domain_user.NewProfileDeleted,
	user_infra_postgres.NewProfileDeletedPostgres,
	user_infra_mapper.NewProfileMapper,
	user_infra_postgres.NewProfileMediaPostgres,
	user_infra_postgres.NewProfilePostgres,
	user_infra_provider.NewProfileProvider,
	user_internal_usecase_user.NewProfileUsecase,
	user_infra_handler.NewPurposeUseHandler,
	user_infra_mapper.NewPurposeUseMapper,
	user_infra_postgres.NewPurposeUseRepo,
	user_infra_redis.NewRedisProvider,
	user_infra_postgres_organization.NewRepository,
	user_infra_handler.NewRoleGroupHandler,
	user_infra_mapper.NewRoleGroupMapper,
	user_internal_usecase.NewRoleGroupRegistry,
	user_infra_postgres.NewRoleGroupRepo,
	user_internal_usecase.NewRoleGroupUsecase,
	user_validator.NewRoleGroupValidator,
	user_infra_handler.NewRoleHandler,
	user_infra_mapper.NewRoleMapper,
	user_infra_postgres.NewRoleRepo,
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
	user_infra_mapper.NewSessionMapper,
	user_infra_postgres.NewSessionRepository,
	user_infra_postgres_subscription.NewSettlementStore,
	user_infra_postgres.NewStatusRepository,
	user_infra_postgres_bootstrapadmin.NewStore,
	user_infra_handler_grpc.NewSubscriptionAdminHandler,
	user_infra_handler_grpc.NewSubscriptionCheckoutHandler,
	user_infra_handler_grpc.NewSubscriptionSettlementHandler,
	user_infra_postgres_entitlement.NewSubscriptionStore,
	user_infra_handler.NewTagHandler,
	user_infra_mapper.NewTagMapper,
	user_infra_postgres.NewTagPostgres,
	user_internal_usecase_user.NewTagUsecase,
	user_infra_postgres.NewTransactionGorm,
	user_internal_job.NewUserDashboardStatsJob,
	user_infra_handler.NewUserInfoHandler,
	user_infra_mapper.NewUserInfoMapper,
	user_infra_postgres.NewUserInfoRepository,
	user_internal_usecase.NewUserInfoUsecase,
	user_internal_domain_user.NewUserProfileEntityWithRequest,
	user_internal_usecase_user.NewUserUsecase,
	user_infra_oauth.NewZaloAuthService,
	user_infra_zns.NewZnsProvider,
	user_infra_scheduler.NewZnsScheduler,
)

func InitializeApp() (*user_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
