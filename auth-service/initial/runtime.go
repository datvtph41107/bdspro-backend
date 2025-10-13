package initial

import (
	_redis "common/redis"
	"context"

	"auth/infra/handler"
	"auth/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	Logger *flogging.FabricLogger

	AuthHandler       *handler.AuthHandler
	UserInfoHandler   *handler.AuthUserProfileHandler
	InternalHandler   *handler.InternalHandler
	OAuthHandler      *handler.OAuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	RoleHandler       *handler.RoleHandler

	// Module service để cache modules trong memory
	ModuleService *usecase.ModuleService
}

func NewApp(authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	internalHandler *handler.InternalHandler,
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	roleHandler *handler.RoleHandler,
	userInfoHandler *handler.AuthUserProfileHandler,
	moduleService *usecase.ModuleService,
) *InitialApp {
	// Config đã được load trong main.go rồi, không cần load lại
	app := &InitialApp{
		AuthHandler:       authHandler,
		UserInfoHandler:   userInfoHandler,
		InternalHandler:   internalHandler,
		OAuthHandler:      oauthHandler,
		PermissionHandler: permissionHandler,
		RoleGroupHandler:  roleGroupHandler,
		RoleHandler:       roleHandler,
		ModuleService:     moduleService,
		Logger:            flogging.MustGetLogger(viper.GetString("server.name")),
	}

	// Load modules khi khởi động
	app.LoadModulesOnStartup()

	return app
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}

// LoadModulesOnStartup load tất cả modules từ database khi khởi động service
func (app *InitialApp) LoadModulesOnStartup() {
	ctx := context.Background()

	err := app.ModuleService.LoadModules(ctx)
	if err != nil {
		app.Logger.Errorf("Failed to load modules on startup: %v", err)
		return
	}

	app.Logger.Infof("Successfully loaded modules on startup")
}
