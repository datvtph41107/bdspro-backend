//go:build wireinject
// +build wireinject

package main

import (
	_db "common/db"
	_redis "common/redis"
	"crm/repo"
	"crm/routes"
	"crm/services"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_redis.NewRedisService,
	_db.NewDB,
	repo.NewFriendRepo,
	repo.NewContactRepo,
	repo.NewFollowRepo,
	repo.NewGroupRepo,
	repo.NewBlockRepo,
	services.NewFriendService,
	services.NewFollowService,
	services.NewBlockService,
	services.NewGroupService,
	routes.NewFriendRouter,
	routes.NewGroupRouter,
	routes.NewFollowRouter,
	routes.NewBlockRouter,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
