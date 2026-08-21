package initial

import (
	_provider "common/domain/provider"
	"context"
	"user/infra/handler"
	"user/infra/scheduler"
	"user/internal/job"
	"user/internal/usecase"
	"user/internal/usecases"

	"github.com/hyperledger/fabric/common/flogging"
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

	Logger *flogging.FabricLogger
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
	app.Logger = flogging.MustGetLogger(loggerName)

	return app
}

// LoadModulesOnStartup performs the one-shot cache warmup under the process context.
// Failure preserves the historical degraded-start behavior; the process root decides
// whether to log/alert rather than constructors hiding background work.
func (app *InitialApp) LoadModulesOnStartup(ctx context.Context) error {
	if app == nil || app.ModuleService == nil {
		return nil
	}
	if err := app.ModuleService.LoadModules(ctx); err != nil {
		return err
	}
	app.Logger.Infof("Successfully loaded modules on startup")
	return nil
}
