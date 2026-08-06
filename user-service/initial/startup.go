package initial

import (
	_provider "common/domain/provider"
	"context"
	"log"
	"os"
	"user/infra/handler"
	"user/infra/scheduler"
	"user/internal/job"
	"user/internal/usecase"
	"user/internal/usecases"

	"github.com/spf13/viper"
)

type InitialApp struct {
	ProfileService *handler.GrpcProfileService
	ProfileUsecase *usecases.ProfileUsecase
	// Admin Services
	AdminHandler            *handler.AdminHandler
	AdminUserProfileHandler *handler.AdminUserProfileHandler
	BookmarkUserHandler     *handler.BookmarkUserHandler
	CrmProvider             _provider.CrmProvider
	// System Config
	// Membership & Plan
	MembershipHandler *handler.MembershipHandler
	PlanHandler       *handler.PlanHandler
	// Internal Service
	AuthInternalHandler   *handler.AuthInternalHandler
	InternalHandler       *handler.InternalHandler
	MainAreaHandler       *handler.MainAreaHandler
	PurposeUseHandler     *handler.PurposeUseHandler
	ProfessionHandler     *handler.ProfessionHandler
	TagHandler            *handler.TagHandler
	KYCHandler            *handler.KYCHandler
	PriceTableHandler     *handler.PriceTableHandler
	UserDashboardStatsJob *job.UserDashboardStatsJob

	AuthHandler     *handler.AuthHandler
	UserInfoHandler *handler.AuthUserProfileHandler
	// InternalHandler   *handler.InternalHandler
	OAuthHandler      *handler.OAuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	RoleHandler       *handler.RoleHandler

	// Module service để cache modules trong memory
	ModuleService *usecase.ModuleService

	// Scheduler để chạy background jobs
	ZnsScheduler *scheduler.ZnsScheduler

	logger *log.Logger
}

func NewInitialApp(
	profileService *handler.GrpcProfileService,
	adminHandler *handler.AdminHandler,
	adminUserProfileHandler *handler.AdminUserProfileHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
	crmProvider _provider.CrmProvider,
	membershipHandler *handler.MembershipHandler,
	planHandler *handler.PlanHandler,
	internalHandler *handler.InternalHandler,
	mainAreaHandler *handler.MainAreaHandler,
	purposeUseHandler *handler.PurposeUseHandler,
	professionHandler *handler.ProfessionHandler,
	tagHandler *handler.TagHandler,
	kycHandler *handler.KYCHandler,
	priceTableHandler *handler.PriceTableHandler,
	userDashboardStatsJob *job.UserDashboardStatsJob,

	authInternalHandler *handler.AuthInternalHandler,
	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	// internalHandler *handler.InternalHandler,
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	roleHandler *handler.RoleHandler,
	userInfoHandler *handler.AuthUserProfileHandler,
	moduleService *usecase.ModuleService,
	znsScheduler *scheduler.ZnsScheduler,
) *InitialApp {
	app := &InitialApp{
		ProfileService:          profileService,
		AdminHandler:            adminHandler,
		AdminUserProfileHandler: adminUserProfileHandler,
		BookmarkUserHandler:     bookmarkUserHandler,
		CrmProvider:             crmProvider,
		MembershipHandler:       membershipHandler,
		PlanHandler:             planHandler,
		InternalHandler:         internalHandler,
		MainAreaHandler:         mainAreaHandler,
		PurposeUseHandler:       purposeUseHandler,
		ProfessionHandler:       professionHandler,
		TagHandler:              tagHandler,
		KYCHandler:              kycHandler,
		PriceTableHandler:       priceTableHandler,
		UserDashboardStatsJob:   userDashboardStatsJob,

		AuthInternalHandler: authInternalHandler,
		AuthHandler:         authHandler,
		UserInfoHandler:     userInfoHandler,
		OAuthHandler:        oauthHandler,
		PermissionHandler:   permissionHandler,
		RoleGroupHandler:    roleGroupHandler,
		RoleHandler:         roleHandler,
		ModuleService:       moduleService,
		ZnsScheduler:        znsScheduler,
	}

	loggerName := viper.GetString("server.name")
	if loggerName == "" {
		loggerName = "user-service"
	}
	app.logger = log.New(os.Stdout, loggerName+": ", log.LstdFlags)

	app.LoadModulesOnStartup()

	app.StartZnsScheduler()

	return app
}

// StartZnsScheduler khởi động ZNS scheduler để refresh token tự động
func (app *InitialApp) StartZnsScheduler() {
	if app.ZnsScheduler != nil {
		app.ZnsScheduler.Start()
		app.logger.Printf("ZNS Token Refresh Scheduler started successfully")
	} else {
		app.logger.Printf("WARN: ZNS Scheduler is nil, skipping scheduler start")
	}
}

// LoadModulesOnStartup load tất cả modules từ database khi khởi động service
func (app *InitialApp) LoadModulesOnStartup() {
	ctx := context.Background()

	err := app.ModuleService.LoadModules(ctx)
	if err != nil {
		app.logger.Printf("ERROR: Failed to load modules on startup: %v", err)
		return
	}

	app.logger.Printf("Successfully loaded modules on startup")
}
