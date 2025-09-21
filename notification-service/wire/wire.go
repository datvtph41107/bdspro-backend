//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	notification_config "notification/config"
	notification_infra_client "notification/infra/client"
	notification_infra_firebase "notification/infra/firebase"
	notification_infra_handler "notification/infra/handler"
	notification_infra_mapper "notification/infra/mapper"
	notification_infra_postgres "notification/infra/postgres"
	"notification/initial"
	notification_initial "notification/initial"
	notification_internal_dto "notification/internal/dto"
	crm_internal_repo "notification/internal/repo"
	noti_internal_repo "notification/internal/repo"
	notification_internal_repo "notification/internal/repo"
	notification_internal_usecase "notification/internal/usecase"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(notification_internal_repo.AccountWarningRepository), new(*notification_infra_postgres.AccountWarningPostgresRepo)),
	wire.Bind(new(noti_internal_repo.AdminHistoryRepo), new(*notification_infra_postgres.AdminHistoryPostgres)),
	wire.Bind(new(crm_internal_repo.HistoryRepo), new(*notification_infra_postgres.HistoryRepo)),
	wire.Bind(new(noti_internal_repo.NotificationRepo), new(*notification_infra_postgres.PostgreNotification)),
	wire.Bind(new(notification_internal_usecase.IAuthClient), new(*notification_infra_client.AuthClient)),
	notification_infra_postgres.NewAccountWarningPostgresRepo,
	notification_internal_usecase.NewAccountWarningUsecase,
	notification_infra_mapper.NewAdminHistoryMapper,
	notification_infra_postgres.NewAdminHistoryPostgres,
	notification_internal_usecase.NewAdminHistoryUsecase,
	notification_infra_client.NewAuthClient,
	notification_infra_handler.NewDealHistoryHandler,
	notification_infra_postgres.NewDealHistoryPostgres,
	notification_internal_usecase.NewDealHistoryUsecase,
	notification_config.NewFirebaseConfig,
	notification_infra_firebase.NewFirebaseService,
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
	notification_infra_postgres.NewPostgreNotification,
	notification_internal_dto.NewSendNotiResponse,
)

func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
