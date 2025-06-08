//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"noti/config"
	noti_internal_repo "noti/internal/repo"
	notification_config "noti/config"
	notification_infra_firebase "noti/infra/firebase"
	notification_infra_postgre "noti/infra/postgre"
	notification_infra_service "noti/infra/service"
	notification_internal_dto "noti/internal/dto"
	notification_internal_usecase "noti/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(noti_internal_repo.NotificationRepo), new(*notification_infra_postgre.PostgreNotification)),
	notification_config.NewConfigApp,
	notification_config.NewFirebaseConfig,
	notification_infra_firebase.NewFirebaseService,
	notification_infra_service.NewGatewayNotiService,
	notification_infra_postgre.NewHistoryRepository,
	notification_infra_service.NewInternalNotiService,
	notification_internal_usecase.NewNotificationService,
	notification_infra_postgre.NewPostgreNotification,
	notification_internal_dto.NewSendNotiResponse,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}