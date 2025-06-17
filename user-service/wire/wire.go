//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"user/config"
	user_config "user/config"
	user_infra_client "user/infra/client"
	user_infra_repo "user/infra/repo"
	user_infra_routes "user/infra/routes"
	user_infra_services "user/infra/services"
	user_internal_dto "user/internal/dto"
	user_internal_models_auth "user/internal/models/auth"
	user_internal_models_user "user/internal/models/user"
	user_internal_services_auth "user/internal/services/auth"
	user_internal_services_user "user/internal/services/user"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	user_infra_routes.NewAuthRoutes,
	user_internal_services_auth.NewAuthService,
	user_validator.NewAuthValidator,
	user_infra_repo.NewBlockRepo,
	user_infra_routes.NewBlockRouter,
	user_internal_services_user.NewBlockService,
	user_infra_repo.NewCertificationRepo,
	user_internal_services_user.NewCertificationService,
	user_config.NewConfigApp,
	user_infra_repo.NewContactRepo,
	user_internal_services_auth.NewCookieService,
	user_internal_services_auth.NewFacebookAuthService,
	user_infra_repo.NewFollowRepo,
	user_infra_routes.NewFollowRouter,
	user_internal_services_user.NewFollowService,
	user_infra_repo.NewFriendRepo,
	user_infra_routes.NewFriendRouter,
	user_internal_services_user.NewFriendService,
	user_internal_services_auth.NewGoogleAuthService,
	user_infra_client.NewGroupClient,
	user_infra_repo.NewGroupRepo,
	user_infra_routes.NewGroupRouter,
	user_internal_services_user.NewGroupService,
	user_infra_services.NewGrpcProfileService,
	user_internal_services_auth.NewHttpClient,
	user_infra_routes.NewOAuthRouter,
	user_internal_services_auth.NewOAuthService,
	user_infra_client.NewOrganizationClient,
	user_infra_client.NewOrganizationMemberClient,
	user_infra_repo.NewProfessionRepo,
	user_infra_routes.NewProfessionRouter,
	user_internal_services_user.NewProfessionService,
	user_infra_repo.NewProfileRepo,
	user_infra_routes.NewProfileRouter,
	user_internal_services_user.NewProfileService,
	user_infra_repo.NewRateRepo,
	user_infra_routes.NewRateRouter,
	user_internal_services_user.NewRateService,
	user_internal_services_auth.NewRedisService,
	user_internal_dto.NewRefreshResponse,
	user_internal_dto.NewRequestLoginResponse,
	user_internal_models_auth.NewUserOTPRepository,
	user_internal_models_user.NewUserProfileEntityWithRequest,
	user_internal_services_auth.NewZaloAuthService,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}