//go:build wireinject
// +build wireinject

package main

import (
	_db "common/db"
	_redis "common/redis"
	"noti/config"
	"noti/repo"
	"noti/routes"
	"noti/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_redis.NewRedisService,
	_db.NewDB,
	config.NewFirebaseConfig,
	repo.NewNotificationRepository,
	services.NewNotificationService,
	services.NewFirebaseService,
	routes.NewNotificationRouter,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
