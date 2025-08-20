//go:build wireinject
// +build wireinject

package wire

import (
	auth_infra_client "auth/infra/client"
	auth_infra_cookie "auth/infra/cookie"
	auth_infra_handler "auth/infra/handler"
	auth_infra_oauth "auth/infra/oauth"
	auth_infra_postgre "auth/infra/postgre"
	auth_infra_properties "auth/infra/properties"
	auth_infra_redis "auth/infra/redis"
	"auth/initial"
	auth_initial "auth/initial"
	auth_internal_domain "auth/internal/domain"
	auth_internal_interface_factory "auth/internal/interface/factory"
	auth_internal_interface_provider "auth/internal/interface/provider"
	auth_internal_interface_repo "auth/internal/interface/repo"
	auth_internal_usecase "auth/internal/usecase"
	_db "common/db"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(auth_internal_interface_provider.OrganizationProvider), new(*auth_infra_client.OrganizationClient)),
	wire.Bind(new(auth_internal_interface_provider.PaymentProvider), new(*auth_infra_client.PaymentClient)),
	wire.Bind(new(auth_internal_interface_provider.ProfileProvider), new(*auth_infra_client.UserClient)),
	wire.Bind(new(auth_internal_interface_provider.NotificationProvider), new(*auth_infra_client.NotificationClient)),
	wire.Bind(new(auth_internal_interface_provider.CookieProvider), new(*auth_infra_cookie.BroswerCookieProvider)),
	wire.Bind(new(auth_internal_interface_repo.AuthMethodRepository), new(*auth_infra_postgre.AuthMethodPostgres)),
	wire.Bind(new(auth_internal_interface_repo.OTPRepository), new(*auth_infra_postgre.OTPPostgres)),
	wire.Bind(new(auth_internal_interface_repo.SessionRepository), new(*auth_infra_postgre.SessionPostgres)),
	wire.Bind(new(auth_internal_interface_repo.StatusRepository), new(*auth_infra_postgre.StatusPostgres)),
	wire.Bind(new(auth_internal_interface_factory.IFactory), new(*auth_infra_properties.RuntimeProperties)),
	wire.Bind(new(auth_internal_interface_provider.CacheProvider), new(*auth_infra_redis.RedisProvider)),
	auth_initial.NewApp,
	auth_infra_handler.NewAuthHandler,
	auth_infra_postgre.NewAuthRepository,
	auth_internal_usecase.NewAuthService,
	auth_infra_cookie.NewCookieProvider,
	auth_infra_oauth.NewFacebookAuthService,
	auth_infra_oauth.NewGoogleAuthService,
	auth_infra_client.NewHttpClient,
	auth_infra_handler.NewOAuthRouter,
	auth_infra_postgre.NewOTPRepository,
	auth_infra_client.NewOrganizationClient,
	auth_internal_usecase.NewOtpUsecase,
	auth_infra_client.NewPaymentClient,
	auth_infra_redis.NewRedisProvider,
	auth_infra_properties.NewRuntimeProperties,
	auth_infra_postgre.NewSessionRepository,
	auth_infra_postgre.NewStatusRepository,
	auth_infra_client.NewUserClient,
	auth_infra_client.NewNotificationClient,
	auth_internal_domain.NewUserOTPRepository,
	auth_infra_oauth.NewZaloAuthService,
)

func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
