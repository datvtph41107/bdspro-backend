package initial

import (
	"context"

	"auth/infra/handler"
	"auth/infra/scheduler"
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

	// Scheduler để chạy background jobs
	ZnsScheduler *scheduler.ZnsScheduler
}

func NewApp(authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	internalHandler *handler.InternalHandler,
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	roleHandler *handler.RoleHandler,
	userInfoHandler *handler.AuthUserProfileHandler,
	moduleService *usecase.ModuleService,
	znsScheduler *scheduler.ZnsScheduler,
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
		ZnsScheduler:      znsScheduler,
		Logger:            flogging.MustGetLogger(viper.GetString("server.name")),
	}

	// Load modules khi khởi động
	app.LoadModulesOnStartup()

	// Khởi động ZNS Token Refresh Scheduler
	app.StartZnsScheduler()

	return app
}

// StartZnsScheduler khởi động ZNS scheduler để refresh token tự động
func (app *InitialApp) StartZnsScheduler() {
	if app.ZnsScheduler != nil {
		app.ZnsScheduler.Start()
		app.Logger.Infof("ZNS Token Refresh Scheduler started successfully")
	} else {
		app.Logger.Warnf("ZNS Scheduler is nil, skipping scheduler start")
	}
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
