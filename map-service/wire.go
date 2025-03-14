//go:build wireinject
// +build wireinject

package main

import (
	c_config "common/config"
	"map/repo"
	"map/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	c_config.NewDatabaseConfig,
	c_config.NewDB,
	repo.NewMapRepo,
	services.NewMapService,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
