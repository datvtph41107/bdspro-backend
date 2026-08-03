//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	_provider "common/provider"
	common_utils "common/utils"
	"github.com/google/wire"
	notification_config "notification/config"
	notification_infra_cache "notification/infra/cache"
	notification_infra_client "notification/infra/client"
	notification_infra_firebase "notification/infra/firebase"
	notification_infra_handler "notification/infra/handler"
	notification_infra_mapper "notification/infra/mapper"
	notification_infra_postgres "notification/infra/postgres"
	notification_initial "notification/initial"
	notification_internal_dto "notification/internal/dto"
	notification_internal_provider "notification/internal/provider"
	notification_internal_repo "notification/internal/repo"
	notification_internal_usecase "notification/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(notification_internal_provider.IAuthClient), new(*notification_infra_client.AuthClient)),
	wire.Bind(new(notification_internal_provider.IFirebaseProvider), new(*notification_infra_firebase.FirebaseProvider)),
	wire.Bind(new(notification_internal_repo.AccountWarningRepository), new(*notification_infra_postgres.AccountWarningPostgresRepo)),
	wire.Bind(new(notification_internal_repo.HistoryAuthRepo), new(*notification_infra_postgres.HistoryAuthPostgres)),
	wire.Bind(new(notification_internal_repo.PersonConfigRepo), new(*notification_infra_postgres.PersonConfigPostgres)),
	common_injection.NewAuthClient,
	common_injection.NewBDSProClient,
	common_injection.NewCRMClient,
	common_injection.NewChatClient,
	common_injection.NewHubClient,
	common_injection.NewNotification,
	common_injection.NewOrganizationClient,
	common_injection.NewPaymentClient,
	common_injection.NewTQClient,
	common_injection.NewTransactionClient,
	common_injection.NewUserClient,
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	notification_infra_postgres.NewAccountWarningPostgresRepo,
	notification_internal_usecase.NewAccountWarningUsecase,
	notification_infra_postgres.NewActivityHistoryPostgres,
	notification_infra_mapper.NewAdminHistoryMapper,
	notification_infra_postgres.NewAdminHistoryPostgres,
	notification_internal_usecase.NewAdminHistoryUsecase,
	notification_infra_client.NewAuthClient,
	notification_infra_handler.NewDealHistoryHandler,
	notification_infra_postgres.NewDealHistoryPostgres,
	notification_internal_usecase.NewDealHistoryUsecase,
	notification_config.NewFirebaseConfig,
	notification_infra_firebase.NewFirebaseService,
	notification_infra_mapper.NewHistoryAuthMapper,
	notification_infra_postgres.NewHistoryAuthPostgres,
	notification_internal_usecase.NewHistoryAuthUsecase,
	notification_infra_handler.NewHistoryHandler,
	notification_infra_mapper.NewHistoryMapper,
	notification_infra_postgres.NewHistoryRepo,
	notification_infra_postgres.NewHistoryRepository,
	notification_internal_usecase.NewHistoryUsecase,
	notification_initial.NewInitialApp,
	notification_infra_handler.NewInternalHandler,
	notification_infra_handler.NewNotificationHandler,
	notification_infra_mapper.NewNotificationMapper,
	notification_internal_usecase.NewNotificationService,
	notification_infra_postgres.NewPersonConfigPostgres,
	notification_internal_usecase.NewPersonConfigUsecase,
	notification_infra_postgres.NewPostgreNotification,
	notification_infra_handler.NewPropertyHistoryHandler,
	notification_infra_mapper.NewPropertyHistoryMapper,
	notification_internal_usecase.NewPropertyHistoryUseCase,
	notification_infra_cache.NewRedisClient,
	notification_internal_dto.NewSendNotiResponse,
	notification_infra_client.NewUserClient,
)

func InitializeApp() (*notification_initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
