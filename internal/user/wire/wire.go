//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	_redis "common/redis"
	common_tilesession "common/tilesession"
	common_utils "common/utils"
	"github.com/google/wire"
	"gorm.io/gorm"
	user_infra_client "user/infra/client"
	user_infra_client_payment "user/infra/client/payment"
	user_infra_cookie "user/infra/cookie"
	user_infra_handler "user/infra/handler"
	user_infra_handler_grpc "user/infra/handler/grpc"
	user_infra_handler_grpc_entitlement "user/infra/handler/grpc/entitlement"
	user_infra_mapper "user/infra/mapper"
	user_infra_oauth "user/infra/oauth"
	user_infra_postgres "user/infra/postgres"
	user_infra_postgres_bootstrapadmin "user/infra/postgres/bootstrapadmin"
	user_infra_postgres_checkout "user/infra/postgres/checkout"
	user_infra_postgres_entitlement "user/infra/postgres/entitlement"
	user_infra_postgres_iampermission "user/infra/postgres/iampermission"
	user_infra_postgres_organization "user/infra/postgres/organization"
	user_infra_postgres_subscription "user/infra/postgres/subscription"
	user_infra_properties "user/infra/properties"
	user_infra_provider "user/infra/provider"
	user_infra_redis "user/infra/redis"
	user_infra_rpc "user/infra/rpc"
	user_infra_scheduler "user/infra/scheduler"
	user_infra_sms "user/infra/sms"
	user_infra_zns "user/infra/zns"
	user_initial "user/initial"
	user_internal_domain_plan "user/internal/domain/plan"
	user_internal_dto "user/internal/dto"
	user_internal_interface "user/internal/interface"
	user_internal_interface_factory "user/internal/interface/factory"
	user_internal_interface_providers "user/internal/interface/providers"
	user_internal_interface_repo "user/internal/interface/repo"
	user_internal_job "user/internal/job"
	user_internal_models "user/internal/models"
	user_internal_usecase "user/internal/usecase"
	user_internal_usecase_bootstrapadmin "user/internal/usecase/bootstrapadmin"
	user_internal_usecase_commercialprofile "user/internal/usecase/commercialprofile"
	user_internal_usecase_entitlement "user/internal/usecase/entitlement"
	user_internal_usecase_organization "user/internal/usecase/organization"
	user_internal_usecase_plan_admin "user/internal/usecase/plan/admin"
	user_internal_usecase_plan_publish "user/internal/usecase/plan/publish"
	user_internal_usecase_subscription "user/internal/usecase/subscription"
	user_internal_usecase_subscription_checkout "user/internal/usecase/subscription/checkout"
	user_internal_usecase_subscription_settlement "user/internal/usecase/subscription/settlement"
	user_internal_usecase_useradmin "user/internal/usecase/useradmin"
	user_internal_usecases "user/internal/usecases"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewTransactionRepo,
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	wire.Bind(new(user_internal_interface_factory.IFactory), new(*user_infra_properties.RuntimeProperties)),
	wire.Bind(new(user_internal_interface_providers.BdsproProvider), new(*user_infra_client.BdsproClient)),
	wire.Bind(new(user_internal_interface_providers.HubProvider), new(*user_infra_client.HubClient)),
	wire.Bind(new(user_internal_interface_providers.OrganizationProvider), new(*user_infra_client.OrganizationProviderClient)),
	wire.Bind(new(user_internal_interface_providers.ProfileProvider), new(*user_infra_provider.ProfileProvider)),
	wire.Bind(new(user_internal_interface_repo.AdminAccessRepository), new(*user_infra_postgres.AdminAccessPostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.DeviceRepository), new(*user_infra_postgres.DevicePostgresRepo)),
	wire.Bind(new(user_internal_interface_repo.IAuthConfigRepo), new(*user_infra_postgres.AuthConfigPostgres)),
	wire.Bind(new(user_internal_interface_repo.IBlockRepo), new(*user_infra_postgres.BlockPostgres)),
	wire.Bind(new(user_internal_interface_repo.IContactRepo), new(*user_infra_postgres.ContactPostgres)),
	wire.Bind(new(user_internal_interface_repo.IFollowRepo), new(*user_infra_postgres.FollowPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileDeletedRepo), new(*user_infra_postgres.ProfileDeletedPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileRepo), new(*user_infra_postgres.ProfilePostgres)),
	wire.Bind(new(user_internal_interface_repo.OTPRepository), new(*user_infra_postgres.OTPPostgres)),
	wire.Bind(new(user_internal_interface_repo.PINRepository), new(*user_infra_postgres.PINPostgres)),
	wire.Bind(new(user_internal_interface_repo.PermissionRepository), new(*user_infra_postgres.PermissionRepo)),
	wire.Bind(new(user_internal_interface_repo.SessionRepository), new(*user_infra_postgres.SessionPostgres)),
	wire.Bind(new(user_internal_interface_repo.StatusRepository), new(*user_infra_postgres.StatusPostgres)),
	wire.Bind(new(user_internal_usecase.AccessTokenIssuer), new(*user_internal_usecase.AuthUsecase)),
	wire.Bind(new(user_internal_usecase_organization.Repository), new(*user_infra_postgres_organization.Repository)),
	wire.Bind(new(user_internal_usecase_useradmin.PermissionAuthorizer), new(*user_infra_postgres_iampermission.Checker)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	common_tilesession.NewStore,
	user_infra_postgres.NewAdminAccessPostgresRepo,
	user_internal_usecase.NewAdminAccessUsecase,
	user_infra_handler.NewAdminHandler,
	user_internal_usecases.NewAdminKYCUsecase,
	user_internal_usecases.NewAdminMainAreaUsecase,
	user_infra_mapper.NewAdminMapper,
	user_infra_handler_grpc.NewAdminOrganizationHandler,
	user_infra_postgres.NewAdminPostgres,
	user_infra_postgres_subscription.NewAdminProjectionStore,
	user_internal_usecases.NewAdminPurposeUseUsecase,
	user_internal_usecase_subscription.NewAdminService,
	user_internal_usecase.NewAdminUsecase,
	user_internal_usecases.NewAdminUsecase,
	user_infra_handler.NewAdminUserProfileHandler,
	user_infra_client.NewAuthClient,
	user_infra_postgres.NewAuthConfigRepository,
	user_infra_handler.NewAuthHandler,
	user_infra_mapper.NewAuthMethodMapper,
	user_infra_rpc.NewAuthRPCClient,
	user_infra_postgres.NewAuthRepository,
	user_infra_mapper.NewAuthSecurityMapper,
	user_internal_usecase.NewAuthSecurityUsecase,
	user_internal_usecase.NewAuthService,
	user_infra_client.NewBdsproClient,
	user_infra_rpc.NewBdsproRPCClient,
	user_infra_postgres.NewBlockPostgres,
	user_internal_usecases.NewBlockUsecase,
	user_infra_handler.NewBookmarkUserHandler,
	user_infra_postgres.NewBookmarkUserPostgres,
	user_internal_usecases.NewBookmarkUserUsecase,
	user_infra_rpc.NewCRMRPCClient,
	user_infra_postgres.NewCertificationPostgres,
	user_internal_usecases.NewCertificationUsecase,
	user_infra_rpc.NewChatRPCClient,
	user_infra_postgres_iampermission.NewChecker,
	user_infra_postgres_organization.NewCheckoutAuthorizer,
	user_infra_postgres_checkout.NewCheckoutStore,
	user_infra_client_payment.NewCommands,
	user_infra_handler_grpc.NewCommercialProfileHandler,
	user_infra_postgres.NewContactPostgres,
	user_infra_cookie.NewCookieProvider,
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
	user_infra_rpc.NewHubRPCClient,
	user_initial.NewInitialApp,
	user_infra_handler.NewInternalHandler,
	user_infra_handler_grpc.NewInternalOrganizationMembershipHandler,
	user_infra_handler.NewKYCHandler,
	user_infra_mapper.NewKYCMapper,
	user_infra_postgres.NewKYCRepo,
	user_internal_usecases.NewKYCUsecase,
	user_infra_handler.NewMainAreaHandler,
	user_infra_mapper.NewMainAreaMapper,
	user_infra_postgres.NewMainAreaPostgres,
	user_infra_client.NewNotificationClient,
	user_infra_rpc.NewNotificationRPCClient,
	user_infra_handler.NewOAuthRouter,
	user_internal_usecase.NewOAuthUsecase,
	user_infra_postgres.NewOTPRepository,
	user_infra_client.NewOrganizationProviderClient,
	user_infra_rpc.NewOrganizationRPCClient,
	user_internal_usecase.NewOtpUsecase,
	user_infra_postgres.NewPINRepository,
	user_internal_usecase.NewPINUsecase,
	user_infra_handler.NewPermissionHandler,
	user_infra_postgres.NewPermissionRepo,
	user_internal_usecase.NewPermissionUsecase,
	user_infra_postgres_entitlement.NewPlanAccessStore,
	user_infra_handler_grpc.NewPlanVersionAdminHandler,
	user_infra_postgres.NewPlanVersionRepo,
	user_infra_handler.NewPriceTableHandler,
	user_infra_mapper.NewPriceTableMapper,
	user_infra_postgres.NewPriceTableRepo,
	user_internal_usecases.NewPriceTableUsecase,
	user_infra_handler.NewProfessionHandler,
	user_infra_mapper.NewProfessionMapper,
	user_infra_postgres.NewProfessionPostgres,
	user_internal_usecases.NewProfessionUsecase,
	user_internal_models.NewProfileDeleted,
	user_infra_postgres.NewProfileDeletedPostgres,
	user_infra_mapper.NewProfileMapper,
	user_infra_postgres.NewProfileMediaPostgres,
	user_infra_postgres.NewProfilePostgres,
	user_infra_provider.NewProfileProvider,
	user_internal_usecases.NewProfileUsecase,
	user_infra_handler.NewPurposeUseHandler,
	user_infra_mapper.NewPurposeUseMapper,
	user_infra_postgres.NewPurposeUseRepo,
	user_infra_redis.NewRedisProvider,
	user_internal_dto.NewRefreshResponse,
	user_internal_domain_plan.NewRegistry,
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
	user_internal_usecases.NewTagUsecase,
	user_infra_postgres.NewTransactionGorm,
	user_internal_job.NewUserDashboardStatsJob,
	user_infra_handler.NewUserInfoHandler,
	user_infra_mapper.NewUserInfoMapper,
	user_infra_postgres.NewUserInfoRepository,
	user_internal_usecase.NewUserInfoUsecase,
	user_internal_models.NewUserProfileEntityWithRequest,
	user_infra_rpc.NewUserRPCClient,
	user_internal_usecases.NewUserUsecase,
	user_infra_oauth.NewZaloAuthService,
	user_infra_zns.NewZnsProvider,
	user_infra_scheduler.NewZnsScheduler,
)

func InitializeApp(database *gorm.DB, redisService *_redis.RedisService) (*user_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
