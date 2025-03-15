//go:build wireinject
// +build wireinject

package main

import (
	c_config "common/config"
	"file/repositories"
	"file/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	c_config.NewDB,
	repositories.NewFileRepository,
	repositories.NewFileAccessRepository,
	services.NewStorageService,
	services.NewFileService,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
