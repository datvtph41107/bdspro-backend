//go:build wireinject
// +build wireinject

package wire

import (
	common_db "common/db"
	common_injection "common/injection"
	sharedredis "common/redis"
	common_utils "common/utils"
	"github.com/google/wire"
	"gorm.io/gorm"
	notification_config "notification/config"
	notification_infra_cache "notification/infra/cache"
	notification_infra_client "notification/infra/client"
	notification_infra_firebase "notification/infra/firebase"
	notification_infra_handler "notification/infra/handler"
	notification_infra_mapper "notification/infra/mapper"
	notification_infra_postgres "notification/infra/postgres"
	notification_initial "notification/initial"
	notification_internal_dto "notification/internal/dto"
	notification_internal_usecase "notification/internal/usecase"
	"pb/clients"
)

var wireSet = wire.NewSet(
	common_db.NewTransactionRepo,
	wire.Bind(new(notification_internal_usecase.PushTokenResolver), new(*notification_infra_client.AuthClient)),
	wire.Bind(new(notification_internal_usecase.PushSender), new(*notification_infra_firebase.FirebaseProvider)),
	wire.Bind(new(notification_internal_usecase.AccountWarningStore), new(*notification_infra_postgres.AccountWarningPostgresRepo)),
	wire.Bind(new(notification_internal_usecase.HistoryAuthStore), new(*notification_infra_postgres.HistoryAuthPostgres)),
	wire.Bind(new(notification_internal_usecase.PersonConfigStore), new(*notification_infra_postgres.PersonConfigPostgres)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	notification_infra_cache.NewRedisClient,
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
	notification_internal_dto.NewSendNotiResponse,
	notification_infra_client.NewUserClient,
)

func InitializeApp(
	database *gorm.DB,
	redisService *sharedredis.RedisService,
	userRPC *clients.UserGrpcClient,
	firebaseConfig *notification_config.FirebaseConfig,
) (*notification_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
