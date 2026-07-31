//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	_provider "common/provider"
	"github.com/google/wire"
	user_infra_client "user/infra/client"
	user_infra_cookie "user/infra/cookie"
	user_infra_handler "user/infra/handler"
	user_infra_mapper "user/infra/mapper"
	user_infra_oauth "user/infra/oauth"
	user_infra_postgres "user/infra/postgres"
	user_infra_properties "user/infra/properties"
	user_infra_provider "user/infra/provider"
	user_infra_redis "user/infra/redis"
	user_infra_routes "user/infra/routes"
	user_infra_scheduler "user/infra/scheduler"
	user_infra_sms "user/infra/sms"
	user_infra_zns "user/infra/zns"
	user_initial "user/initial"
	user_internal_domain "user/internal/domain"
	user_internal_dto "user/internal/dto"
	user_internal_interface "user/internal/interface"
	auth_internal_interface_factory "user/internal/interface/factory"
	user_internal_interface_providers "user/internal/interface/providers"
	auth_internal_interface_repo "user/internal/interface/repo"
	user_internal_interface_repo "user/internal/interface/repo"
	user_internal_job "user/internal/job"
	user_internal_models "user/internal/models"
	user_internal_usecase "user/internal/usecase"
	user_internal_usecases "user/internal/usecases"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(auth_internal_interface_factory.IFactory), new(*user_infra_properties.RuntimeProperties)),
	wire.Bind(new(auth_internal_interface_repo.AdminAccessRepository), new(*user_infra_postgres.AdminAccessPostgresRepo)),
	wire.Bind(new(auth_internal_interface_repo.DeviceRepository), new(*user_infra_postgres.DevicePostgresRepo)),
	wire.Bind(new(auth_internal_interface_repo.IAuthConfigRepo), new(*user_infra_postgres.AuthConfigPostgres)),
	wire.Bind(new(auth_internal_interface_repo.OTPRepository), new(*user_infra_postgres.OTPPostgres)),
	wire.Bind(new(auth_internal_interface_repo.PINRepository), new(*user_infra_postgres.PINPostgres)),
	wire.Bind(new(auth_internal_interface_repo.SessionRepository), new(*user_infra_postgres.SessionPostgres)),
	wire.Bind(new(auth_internal_interface_repo.StatusRepository), new(*user_infra_postgres.StatusPostgres)),
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	wire.Bind(new(user_internal_interface_providers.BdsproProvider), new(*user_infra_client.BdsproClient)),
	wire.Bind(new(user_internal_interface_providers.CrmProvider), new(*user_infra_client.CrmClient)),
	wire.Bind(new(user_internal_interface_providers.HubProvider), new(*user_infra_client.HubClient)),
	wire.Bind(new(user_internal_interface_providers.OrganizationProvider), new(*user_infra_client.OrganizationProviderClient)),
	wire.Bind(new(user_internal_interface_providers.ProfileProvider), new(*user_infra_provider.ProfileProvider)),
	wire.Bind(new(user_internal_interface_repo.IBlockRepo), new(*user_infra_postgres.BlockPostgres)),
	wire.Bind(new(user_internal_interface_repo.IContactRepo), new(*user_infra_postgres.ContactPostgres)),
	wire.Bind(new(user_internal_interface_repo.IFollowRepo), new(*user_infra_postgres.FollowPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileDeletedRepo), new(*user_infra_postgres.ProfileDeletedPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileRepo), new(*user_infra_postgres.ProfilePostgres)),
	user_infra_postgres.NewAdminAccessPostgresRepo,
	user_internal_usecase.NewAdminAccessUsecase,
	user_infra_handler.NewAdminHandler,
	user_internal_usecases.NewAdminKYCUsecase,
	user_internal_usecases.NewAdminMainAreaUsecase,
	user_infra_mapper.NewAdminMapper,
	user_infra_postgres.NewAdminPostgres,
	user_internal_usecases.NewAdminPurposeUseUsecase,
	user_internal_usecase.NewAdminUsecase,
	user_internal_usecases.NewAdminUsecase,
	user_infra_handler.NewAdminUserProfileHandler,
	user_infra_postgres.NewAttachPostgres,
	user_infra_client.NewAuthClient,
	user_infra_postgres.NewAuthConfigRepository,
	user_infra_handler.NewAuthHandler,
	user_infra_handler.NewAuthInternalHandler,
	user_infra_mapper.NewAuthMethodMapper,
	user_infra_postgres.NewAuthRepository,
	user_infra_mapper.NewAuthSecurityMapper,
	user_internal_usecase.NewAuthSecurityUsecase,
	user_internal_usecase.NewAuthService,
	user_infra_client.NewBdsproClient,
	user_infra_postgres.NewBlockPostgres,
	user_infra_routes.NewBlockRouter,
	user_internal_usecases.NewBlockUsecase,
	user_infra_handler.NewBookmarkUserHandler,
	user_infra_postgres.NewBookmarkUserPostgres,
	user_internal_usecases.NewBookmarkUserUsecase,
	user_infra_postgres.NewCertificationPostgres,
	user_internal_usecases.NewCertificationUsecase,
	user_infra_client.NewChatClient,
	user_infra_postgres.NewContactPostgres,
	user_infra_cookie.NewCookieProvider,
	user_infra_client.NewCrmClient,
	user_infra_mapper.NewDeviceMapper,
	user_infra_postgres.NewDevicePostgresRepo,
	user_internal_usecase.NewDeviceUsecase,
	user_infra_oauth.NewFacebookAuthService,
	user_infra_client.NewFeedbackClient,
	user_infra_postgres.NewFollowPostgres,
	user_infra_routes.NewFollowRouter,
	user_internal_usecases.NewFollowUsecase,
	user_infra_postgres.NewFriendPostgres,
	user_infra_routes.NewFriendRouter,
	user_internal_usecases.NewFriendUsecase,
	user_infra_oauth.NewGoogleAuthService,
	user_infra_client.NewGroupClient,
	user_infra_postgres.NewGroupPostgres,
	user_infra_handler.NewGrpcProfileService,
	user_infra_client.NewHttpClient,
	user_infra_client.NewHubClient,
	user_initial.NewInitialApp,
	user_infra_handler.NewInternalHandler,
	user_infra_handler.NewKYCHandler,
	user_infra_mapper.NewKYCMapper,
	user_infra_postgres.NewKYCRepo,
	user_internal_usecases.NewKYCUsecase,
	user_infra_handler.NewMainAreaHandler,
	user_infra_mapper.NewMainAreaMapper,
	user_infra_postgres.NewMainAreaPostgres,
	user_infra_handler.NewMembershipHandler,
	user_infra_mapper.NewMembershipMapper,
	user_infra_postgres.NewMembershipRepo,
	user_internal_usecases.NewMembershipUsecase,
	user_infra_postgres.NewModuleRepo,
	user_internal_usecase.NewModuleService,
	user_internal_usecase.NewModuleUsecase,
	user_infra_client.NewNotificationClient,
	user_infra_handler.NewOAuthRouter,
	user_internal_usecase.NewOAuthUsecase,
	user_infra_postgres.NewOTPRepository,
	user_infra_client.NewOrganizationClient,
	user_infra_client.NewOrganizationMemberClient,
	user_infra_client.NewOrganizationProviderClient,
	user_internal_usecase.NewOtpUsecase,
	user_infra_postgres.NewPINRepository,
	user_internal_usecase.NewPINUsecase,
	user_infra_client.NewPaymentClient,
	user_infra_handler.NewPermissionHandler,
	user_infra_postgres.NewPermissionRepo,
	user_internal_usecase.NewPermissionUsecase,
	user_infra_handler.NewPlanHandler,
	user_infra_mapper.NewPlanMapper,
	user_infra_postgres.NewPlanRepo,
	user_internal_usecases.NewPlanUsecase,
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
	user_infra_routes.NewProfileRouter,
	user_internal_usecases.NewProfileUsecase,
	user_infra_handler.NewPurposeUseHandler,
	user_infra_mapper.NewPurposeUseMapper,
	user_infra_postgres.NewPurposeUseRepo,
	user_infra_redis.NewRedisProvider,
	user_internal_dto.NewRefreshResponse,
	user_infra_handler.NewRoleGroupHandler,
	user_infra_mapper.NewRoleGroupMapper,
	user_infra_postgres.NewRoleGroupRepo,
	user_internal_usecase.NewRoleGroupUsecase,
	user_validator.NewRoleGroupValidator,
	user_infra_handler.NewRoleHandler,
	user_infra_mapper.NewRoleMapper,
	user_infra_postgres.NewRoleRepo,
	user_internal_usecase.NewRoleUsecase,
	user_infra_properties.NewRuntimeProperties,
	user_infra_sms.NewSMSProvider,
	user_infra_mapper.NewSessionMapper,
	user_infra_postgres.NewSessionRepository,
	user_infra_postgres.NewStatusRepository,
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
	user_internal_domain.NewUserOTPRepository,
	user_internal_models.NewUserProfileEntityWithRequest,
	user_internal_usecases.NewUserUsecase,
	user_infra_oauth.NewZaloAuthService,
	user_infra_zns.NewZnsProvider,
	user_infra_scheduler.NewZnsScheduler,
)

func InitializeApp() (*user_initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
