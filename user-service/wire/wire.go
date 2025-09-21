//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"user/initial"
	user_infra_client "user/infra/client"
	user_infra_handler "user/infra/handler"
	user_infra_mapper "user/infra/mapper"
	user_infra_postgres "user/infra/postgres"
	user_infra_repo "user/infra/repo"
	user_infra_routes "user/infra/routes"
	user_infra_services "user/infra/services"
	user_initial "user/initial"
	user_internal_dto "user/internal/dto"
	user_internal_interface "user/internal/interface"
	user_internal_interface_provider "user/internal/interface/provider"
	user_internal_models_auth "user/internal/models/auth"
	user_internal_models_user "user/internal/models/user"
	user_internal_repo "user/internal/repo"
	user_internal_services_auth "user/internal/services/auth"
	user_internal_services_user "user/internal/services/user"
	user_internal_usecases "user/internal/usecases"
	user_validator "user/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(user_internal_interface_provider.AuthProvider), new(*user_infra_client.AuthClient)),
	wire.Bind(new(user_internal_repo.RateRepo), new(*user_infra_postgres.RatePostgres)),
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	wire.Bind(new(user_internal_repo.AdminRepository), new(*user_infra_repo.AdminRepo)),
	wire.Bind(new(user_internal_repo.ReportRepo), new(*user_infra_repo.ReportRepo)),
	user_infra_handler.NewAdminHandler,
	user_infra_repo.NewAdminRepo,
	user_internal_usecases.NewAdminUsecase,
	user_infra_postgres.NewAttachPostgres,
	user_infra_client.NewAuthClient,
	user_infra_routes.NewAuthRoutes,
	user_internal_services_auth.NewAuthService,
	user_validator.NewAuthValidator,
	user_infra_repo.NewBlockRepo,
	user_infra_routes.NewBlockRouter,
	user_internal_services_user.NewBlockService,
	user_infra_repo.NewCertificationRepo,
	user_internal_services_user.NewCertificationService,
	user_infra_repo.NewContactRepo,
	user_internal_services_auth.NewCookieService,
	user_infra_client.NewCrmClient,
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
	user_infra_services.NewGrpcAuthService,
	user_infra_services.NewGrpcProfileService,
	user_infra_services.NewGrpcRateService,
	user_internal_services_auth.NewHttpClient,
	user_initial.NewInitialApp,
	user_infra_routes.NewOAuthRouter,
	user_internal_services_auth.NewOAuthService,
	user_infra_client.NewOrganizationClient,
	user_infra_client.NewOrganizationMemberClient,
	user_infra_client.NewPaymentClient,
	user_infra_repo.NewProfessionRepo,
	user_infra_routes.NewProfessionRouter,
	user_internal_services_user.NewProfessionService,
	user_infra_mapper.NewProfileMapper,
	user_infra_repo.NewProfileRepo,
	user_infra_routes.NewProfileRouter,
	user_internal_services_user.NewProfileService,
	user_infra_mapper.NewRateMapper,
	user_infra_postgres.NewRatePostgres,
	user_infra_repo.NewRateRepo,
	user_infra_routes.NewRateRouter,
	user_internal_services_user.NewRateService,
	user_internal_usecases.NewRateUsecase,
	user_internal_services_auth.NewRedisService,
	user_internal_dto.NewRefreshResponse,
	user_infra_handler.NewReportHandler,
	user_infra_mapper.NewReportMapper,
	user_infra_handler.NewReportProofHandler,
	user_infra_mapper.NewReportProofMapper,
	user_infra_handler.NewReportReasonHandler,
	user_infra_repo.NewReportRepo,
	user_internal_services_user.NewReportService,
	user_internal_usecases.NewReportUsecase,
	user_internal_dto.NewRequestLoginResponse,
	user_infra_postgres.NewTransactionGorm,
	user_internal_models_auth.NewUserOTPRepository,
	user_internal_models_user.NewUserProfileEntityWithRequest,
	user_internal_services_auth.NewZaloAuthService,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}