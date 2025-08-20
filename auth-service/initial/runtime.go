package initial

import (
	"auth/config"
	"auth/infra/handler"
	"common"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	AuthHandler  *handler.AuthHandler
	AdminHandler *handler.AdminHandler
	OAuthHandler *handler.OAuthHandler
	Logger       *flogging.FabricLogger
}

func NewApp(authHandler *handler.AuthHandler, adminHandler *handler.AdminHandler, oauthHandler *handler.OAuthHandler) *InitialApp {
	common.LoadConfig(config.EmbeddedConfig)
	return &InitialApp{
		AuthHandler:  authHandler,
		AdminHandler: adminHandler,
		OAuthHandler: oauthHandler,
		Logger:       flogging.MustGetLogger(viper.GetString("server.name")),
	}
}
