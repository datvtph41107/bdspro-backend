//go:build wireinject
// +build wireinject

package main

import (
	"bdspro/mod/district"
	"bdspro/mod/province"
	"bdspro/mod/ward"
	_db "common/db"
	_redis "common/redis"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_redis.NewRedisService,
	_db.NewDB,
	province.NewProvinceRepo,
	province.NewProvinceService,
	province.NewProvinceRouter,
	district.NewDistrictRepo,
	district.NewDistrictService,
	district.NewDistrictRouter,
	ward.NewWardRepo,
	ward.NewWardService,
	ward.NewWardRouter,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
