//go:build wireinject
// +build wireinject

package main

import (
	_db "common/db"
	_redis "common/redis"
	"membership/repo"
	"membership/routes"
	"membership/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_redis.NewRedisService,
	_db.NewDB,
	repo.NewPlanRepo,
	repo.NewMembershipRepo,
	services.NewPlanService,
	services.NewMembershipService,
	routes.NewPlanRouter,
	routes.NewMembershipRouter,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
