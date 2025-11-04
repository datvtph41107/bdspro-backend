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
	user_infra_routes "user/infra/routes"
	user_initial "user/initial"
	user_internal_dto "user/internal/dto"
	user_internal_interface "user/internal/interface"
	user_internal_interface_providers "user/internal/interface/providers"
	user_internal_interface_repo "user/internal/interface/repo"
	user_internal_job "user/internal/job"
	user_internal_models "user/internal/models"
	user_internal_usecases "user/internal/usecases"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(user_internal_interface_providers.AuthProvider), new(*user_infra_client.AuthClient)),
	wire.Bind(new(user_internal_interface_providers.FriendProvider), new(*user_infra_client.CrmClient)),
	wire.Bind(new(user_internal_interface_providers.FeedbackProvider), new(*user_infra_client.FeedbackClient)),
	wire.Bind(new(user_internal_interface_providers.HubProvider), new(*user_infra_client.HubClient)),
	wire.Bind(new(user_internal_interface_providers.NotificationProvider), new(*user_infra_client.NotificationClient)),
	wire.Bind(new(user_internal_interface_repo.IBlockRepo), new(*user_infra_postgres.BlockPostgres)),
	wire.Bind(new(user_internal_interface_repo.IContactRepo), new(*user_infra_postgres.ContactPostgres)),
	wire.Bind(new(user_internal_interface_repo.IFollowRepo), new(*user_infra_postgres.FollowPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfessionRepo), new(*user_infra_postgres.ProfessionPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileDeletedRepo), new(*user_infra_postgres.ProfileDeletedPostgres)),
	wire.Bind(new(user_internal_interface_repo.IProfileRepo), new(*user_infra_postgres.ProfilePostgres)),
	wire.Bind(new(user_internal_interface.ITransaction), new(*user_infra_postgres.TransactionGorm)),
	user_infra_handler.NewAdminHandler,
	user_infra_mapper.NewAdminMapper,
	user_infra_postgres.NewAdminPostgres,
	user_internal_usecases.NewAdminUsecase,
	user_infra_handler.NewAdminUserProfileHandler,
	user_infra_postgres.NewAttachPostgres,
	user_infra_client.NewAuthClient,
	user_infra_postgres.NewBlockPostgres,
	user_infra_routes.NewBlockRouter,
	user_internal_usecases.NewBlockUsecase,
	user_infra_handler.NewBookmarkUserHandler,
	user_infra_postgres.NewBookmarkUserPostgres,
	user_internal_usecases.NewBookmarkUserUsecase,
	user_infra_postgres.NewCertificationPostgres,
	user_internal_usecases.NewCertificationUsecase,
	user_infra_postgres.NewContactPostgres,
	user_infra_client.NewCrmClient,
	user_infra_client.NewFeedbackClient,
	user_infra_postgres.NewFollowPostgres,
	user_infra_routes.NewFollowRouter,
	user_internal_usecases.NewFollowUsecase,
	user_infra_postgres.NewFriendPostgres,
	user_infra_routes.NewFriendRouter,
	user_internal_usecases.NewFriendUsecase,
	user_infra_client.NewGroupClient,
	user_infra_postgres.NewGroupPostgres,
	user_infra_handler.NewGrpcProfileService,
	user_infra_client.NewHubClient,
	user_initial.NewInitialApp,
	user_infra_handler.NewInternalHandler,
	user_infra_handler.NewMembershipHandler,
	user_infra_mapper.NewMembershipMapper,
	user_infra_postgres.NewMembershipRepo,
	user_internal_usecases.NewMembershipUsecase,
	user_infra_client.NewNotificationClient,
	user_infra_client.NewOrganizationClient,
	user_infra_client.NewOrganizationMemberClient,
	user_infra_client.NewPaymentClient,
	user_infra_handler.NewPlanHandler,
	user_infra_mapper.NewPlanMapper,
	user_infra_postgres.NewPlanRepo,
	user_internal_usecases.NewPlanUsecase,
	user_infra_postgres.NewProfessionPostgres,
	user_internal_usecases.NewProfessionUsecase,
	user_internal_models.NewProfileDeleted,
	user_infra_postgres.NewProfileDeletedPostgres,
	user_infra_mapper.NewProfileMapper,
	user_infra_postgres.NewProfilePostgres,
	user_infra_routes.NewProfileRouter,
	user_internal_usecases.NewProfileUsecase,
	user_internal_dto.NewRefreshResponse,
	user_infra_postgres.NewTransactionGorm,
	user_internal_job.NewUserDashboardStatsJob,
	user_internal_models.NewUserProfileEntityWithRequest,
	user_internal_usecases.NewUserUsecase,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}