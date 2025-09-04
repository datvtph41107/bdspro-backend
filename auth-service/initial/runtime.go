package initial

import (
	"common"
	_redis "common/redis"

	"auth/config"
	"auth/infra/handler"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	Logger *flogging.FabricLogger

	AuthHandler       *handler.AuthHandler
	InternalHandler   *handler.InternalHandler
	OAuthHandler      *handler.OAuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	RoleHandler       *handler.RoleHandler
}

func NewApp(authHandler *handler.AuthHandler, oauthHandler *handler.OAuthHandler, internalHandler *handler.InternalHandler, permissionHandler *handler.PermissionHandler, roleGroupHandler *handler.RoleGroupHandler, roleHandler *handler.RoleHandler) *InitialApp {
	common.LoadConfig(config.EmbeddedConfig)
	return &InitialApp{
		AuthHandler:       authHandler,
		InternalHandler:   internalHandler,
		OAuthHandler:      oauthHandler,
		PermissionHandler: permissionHandler,
		RoleGroupHandler:  roleGroupHandler,
		RoleHandler:       roleHandler,
		Logger:            flogging.MustGetLogger(viper.GetString("server.name")),
	}
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
