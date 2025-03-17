//go:build wireinject
// +build wireinject

package main

import (
	"auth/routes"
	"auth/services"
	_redis "common/redis"

	"github.com/google/wire"
)

// WireSet tập hợp các dependencies để inject
var wireSet = wire.NewSet(
	_redis.NewRedisService,
	services.NewCookieService,
	services.NewAuthService,
	services.NewHttpClient,
	services.NewFacebookAuthService,
	services.NewGoogleAuthService,
	services.NewZaloAuthService,
	services.NewProfileService,
	routes.NewAuthRoutes,
	routes.NewProfileRouter,
	routes.NewOAuthRouter,
	NewConfigApp,
)

// InitializeApp khởi tạo toàn bộ hệ thống với Wire
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
