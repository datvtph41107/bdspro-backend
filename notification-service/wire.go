//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	_db "common/db"
	noti_config "noti/config"
	noti_data "noti/data"
	noti_repo "noti/repo"
	noti_routes "noti/routes"
	noti_server "noti/server"
	noti_services "noti/services"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	NewConfigApp,
	noti_config.NewFirebaseConfig,
	noti_services.NewFirebaseService,
	noti_repo.NewHistoryRepository,
	noti_repo.NewNotificationRepository,
	noti_routes.NewNotificationRouter,
	noti_services.NewNotificationService,
	noti_services.NewNotificationService2,
	noti_data.NewSendNotiResponse,
	noti_server.NewServerHttp,
)
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}