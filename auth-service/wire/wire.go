//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"auth/initial"
	auth_infra_client "auth/infra/client"
	auth_infra_cookie "auth/infra/cookie"
	auth_infra_handler "auth/infra/handler"
	auth_infra_mapper "auth/infra/mapper"
	auth_infra_oauth "auth/infra/oauth"
	auth_infra_postgre "auth/infra/postgre"
	auth_infra_properties "auth/infra/properties"
	auth_infra_redis "auth/infra/redis"
	auth_infra_scheduler "auth/infra/scheduler"
	auth_infra_zns "auth/infra/zns"
	auth_initial "auth/initial"
	auth_internal_domain "auth/internal/domain"
	auth_internal_interface_factory "auth/internal/interface/factory"
	auth_internal_interface_provider "auth/internal/interface/provider"
	auth_internal_interface_repo "auth/internal/interface/repo"
	auth_internal_usecase "auth/internal/usecase"
	auth_validator "auth/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(auth_internal_interface_provider.NotificationProvider), new(*auth_infra_client.NotificationClient)),
	wire.Bind(new(auth_internal_interface_provider.OrganizationProvider), new(*auth_infra_client.OrganizationClient)),
	wire.Bind(new(auth_internal_interface_provider.PaymentProvider), new(*auth_infra_client.PaymentClient)),
	wire.Bind(new(auth_internal_interface_provider.ProfileProvider), new(*auth_infra_client.UserClient)),
	wire.Bind(new(auth_internal_interface_provider.CookieProvider), new(*auth_infra_cookie.BroswerCookieProvider)),
	wire.Bind(new(auth_internal_interface_repo.AdminAccessRepository), new(*auth_infra_postgre.AdminAccessPostgresRepo)),
	wire.Bind(new(auth_internal_interface_repo.OTPRepository), new(*auth_infra_postgre.OTPPostgres)),
	wire.Bind(new(auth_internal_interface_repo.SessionRepository), new(*auth_infra_postgre.SessionPostgres)),
	wire.Bind(new(auth_internal_interface_repo.StatusRepository), new(*auth_infra_postgre.StatusPostgres)),
	wire.Bind(new(auth_internal_interface_factory.IFactory), new(*auth_infra_properties.RuntimeProperties)),
	wire.Bind(new(auth_internal_interface_provider.CacheProvider), new(*auth_infra_redis.RedisProvider)),
	wire.Bind(new(auth_internal_interface_provider.IZnsProvider), new(*auth_infra_zns.ZnsProvider)),
	auth_infra_postgre.NewAdminAccessPostgresRepo,
	auth_internal_usecase.NewAdminAccessUsecase,
	auth_internal_usecase.NewAdminUsecase,
	auth_initial.NewApp,
	auth_infra_handler.NewAuthHandler,
	auth_infra_mapper.NewAuthMethodMapper,
	auth_infra_postgre.NewAuthRepository,
	auth_internal_usecase.NewAuthService,
	auth_infra_cookie.NewCookieProvider,
	auth_infra_oauth.NewFacebookAuthService,
	auth_infra_oauth.NewGoogleAuthService,
	auth_infra_client.NewGrpcClient,
	auth_infra_client.NewHttpClient,
	auth_infra_handler.NewInternalHandler,
	auth_infra_postgre.NewModuleRepo,
	auth_internal_usecase.NewModuleService,
	auth_internal_usecase.NewModuleUsecase,
	auth_infra_client.NewNotificationClient,
	auth_infra_handler.NewOAuthRouter,
	auth_internal_usecase.NewOAuthUsecase,
	auth_infra_postgre.NewOTPRepository,
	auth_infra_client.NewOrganizationClient,
	auth_internal_usecase.NewOtpUsecase,
	auth_infra_client.NewPaymentClient,
	auth_infra_handler.NewPermissionHandler,
	auth_infra_postgre.NewPermissionRepo,
	auth_internal_usecase.NewPermissionUsecase,
	auth_infra_redis.NewRedisProvider,
	auth_initial.NewRedisService,
	auth_infra_handler.NewRoleGroupHandler,
	auth_infra_mapper.NewRoleGroupMapper,
	auth_infra_postgre.NewRoleGroupRepo,
	auth_internal_usecase.NewRoleGroupUsecase,
	auth_validator.NewRoleGroupValidator,
	auth_infra_handler.NewRoleHandler,
	auth_infra_mapper.NewRoleMapper,
	auth_infra_postgre.NewRoleRepo,
	auth_internal_usecase.NewRoleUsecase,
	auth_infra_properties.NewRuntimeProperties,
	auth_infra_mapper.NewSessionMapper,
	auth_infra_postgre.NewSessionRepository,
	auth_infra_postgre.NewStatusRepository,
	auth_infra_client.NewUserClient,
	auth_infra_handler.NewUserInfoHandler,
	auth_infra_mapper.NewUserInfoMapper,
	auth_infra_postgre.NewUserInfoRepository,
	auth_internal_usecase.NewUserInfoUsecase,
	auth_internal_domain.NewUserOTPRepository,
	auth_infra_oauth.NewZaloAuthService,
	auth_infra_zns.NewZnsProvider,
	auth_infra_scheduler.NewZnsScheduler,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}