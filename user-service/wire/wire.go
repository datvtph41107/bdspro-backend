//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"user/config"
	user_config "user/config"
	user_data "user/data"
	user_models_auth "user/models/auth"
	user_models_user "user/models/user"
	user_repo "user/repo"
	user_routes "user/routes"
	user_services_auth "user/services/auth"
	user_services_user "user/services/user"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	user_routes.NewAuthRoutes,
	user_services_auth.NewAuthService,
	user_validator.NewAuthValidator,
	user_repo.NewBlockRepo,
	user_routes.NewBlockRouter,
	user_services_user.NewBlockService,
	user_config.NewConfigApp,
	user_repo.NewContactRepo,
	user_services_auth.NewCookieService,
	user_services_auth.NewFacebookAuthService,
	user_repo.NewFollowRepo,
	user_routes.NewFollowRouter,
	user_services_user.NewFollowService,
	user_repo.NewFriendRepo,
	user_routes.NewFriendRouter,
	user_services_user.NewFriendService,
	user_services_auth.NewGoogleAuthService,
	user_repo.NewGroupRepo,
	user_routes.NewGroupRouter,
	user_services_user.NewGroupService,
	user_services_auth.NewHttpClient,
	user_routes.NewOAuthRouter,
	user_services_auth.NewOAuthService,
	user_repo.NewProfileRepo,
	user_routes.NewProfileRouter,
	user_services_auth.NewProfileService,
	user_services_auth.NewRedisService,
	user_data.NewRefreshResponse,
	user_data.NewRequestLoginResponse,
	user_data.NewUserInfoResponse,
	user_models_auth.NewUserOTPRepository,
	user_models_user.NewUserProfileEntityWithRequest,
	user_services_auth.NewZaloAuthService,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}