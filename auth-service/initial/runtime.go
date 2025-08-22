package initial

import (
	"auth/config"
	"auth/infra/handler"
	"common"
	_redis "common/redis"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	AuthHandler       *handler.AuthHandler
	OAuthHandler      *handler.OAuthHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	AdminHandler      *handler.AdminHandler
	InternalHandler   *handler.InternalHandler
	Logger            *flogging.FabricLogger
}

func NewApp(
	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	adminHandler *handler.AdminHandler,
	internalHandler *handler.InternalHandler,
) *InitialApp {
	common.LoadConfig(config.EmbeddedConfig)
	return &InitialApp{
		AuthHandler:       authHandler,
		OAuthHandler:      oauthHandler,
		RoleHandler:       roleHandler,
		PermissionHandler: permissionHandler,
		RoleGroupHandler:  roleGroupHandler,
		AdminHandler:      adminHandler,
		InternalHandler:   internalHandler,
		Logger:            flogging.MustGetLogger(viper.GetString("server.name")),
	}
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
