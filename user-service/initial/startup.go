package initial

import (
	"context"
	"user/infra/client"
	"user/infra/handler"
	"user/infra/routes"
	"user/infra/scheduler"
	"user/internal/job"
	"user/internal/usecase"
	"user/internal/usecases"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	ProfileRouter  *routes.ProfileRouter
	BlockRouter    *routes.BlockRouter
	FollowRouter   *routes.FollowRouter
	FriendRouter   *routes.FriendRouter
	ProfileService *handler.GrpcProfileService
	ProfileUsecase *usecases.ProfileUsecase
	// Admin Services
	AdminHandler            *handler.AdminHandler
	AdminUserProfileHandler *handler.AdminUserProfileHandler
	BookmarkUserHandler     *handler.BookmarkUserHandler
	CrmClient               *client.CrmClient
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

	Logger *flogging.FabricLogger
}

func NewInitialApp(
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	profileService *handler.GrpcProfileService,
	adminHandler *handler.AdminHandler,
	adminUserProfileHandler *handler.AdminUserProfileHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
	crmClient *client.CrmClient,
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
		ProfileRouter:           profile,
		BlockRouter:             block,
		FollowRouter:            follow,
		FriendRouter:            friend,
		ProfileService:          profileService,
		AdminHandler:            adminHandler,
		AdminUserProfileHandler: adminUserProfileHandler,
		BookmarkUserHandler:     bookmarkUserHandler,
		CrmClient:               crmClient,
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
	app.Logger = flogging.MustGetLogger(loggerName)

	app.LoadModulesOnStartup()

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
