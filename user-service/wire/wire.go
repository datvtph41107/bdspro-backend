//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"user/config"
	user_config "user/config"
	user_data "user/data"
	user_infra_routes "user/infra/routes"
	user_infra_services "user/infra/services"
	user_models_auth "user/models/auth"
	user_models_user "user/models/user"
	user_repo "user/repo"
	user_services_auth "user/services/auth"
	user_services_user "user/services/user"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	user_infra_routes.NewAuthRoutes,
	user_services_auth.NewAuthService,
	user_validator.NewAuthValidator,
	user_repo.NewBlockRepo,
	user_infra_routes.NewBlockRouter,
	user_services_user.NewBlockService,
	user_repo.NewCertificationRepo,
	user_services_user.NewCertificationService,
	user_config.NewConfigApp,
	user_repo.NewContactRepo,
	user_services_auth.NewCookieService,
	user_services_auth.NewFacebookAuthService,
	user_repo.NewFollowRepo,
	user_infra_routes.NewFollowRouter,
	user_services_user.NewFollowService,
	user_repo.NewFriendRepo,
	user_infra_routes.NewFriendRouter,
	user_services_user.NewFriendService,
	user_services_auth.NewGoogleAuthService,
	user_repo.NewGroupRepo,
	user_infra_routes.NewGroupRouter,
	user_services_user.NewGroupService,
	user_infra_services.NewGrpcProfileService,
	user_services_auth.NewHttpClient,
	user_infra_routes.NewOAuthRouter,
	user_services_auth.NewOAuthService,
	user_services_auth.NewOrganizationClient,
	user_repo.NewProfessionRepo,
	user_infra_routes.NewProfessionRouter,
	user_services_user.NewProfessionService,
	user_repo.NewProfileRepo,
	user_infra_routes.NewProfileRouter,
	user_services_user.NewProfileService,
	user_repo.NewRateRepo,
	user_infra_routes.NewRateRouter,
	user_services_user.NewRateService,
	user_services_auth.NewRedisService,
	user_data.NewRefreshResponse,
	user_data.NewRequestLoginResponse,
	user_models_auth.NewUserOTPRepository,
	user_models_user.NewUserProfileEntityWithRequest,
	user_services_auth.NewZaloAuthService,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}