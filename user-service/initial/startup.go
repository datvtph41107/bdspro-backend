package initial

import (
	_provider "common/domain/provider"
	"context"
	"pb/clients"
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
	// Internal Service
	InternalHandler       *handler.InternalHandler
	MainAreaHandler       *handler.MainAreaHandler
	PurposeUseHandler     *handler.PurposeUseHandler
	ProfessionHandler     *handler.ProfessionHandler
	TagHandler            *handler.TagHandler
	KYCHandler            *handler.KYCHandler
	PriceTableHandler     *handler.PriceTableHandler
	UserDashboardStatsJob *job.UserDashboardStatsJob
	OrganizationClient    *clients.OrganizationClient

	AuthHandler     *handler.AuthHandler
	UserInfoHandler *handler.AuthUserProfileHandler
	// InternalHandler   *handler.InternalHandler
	OAuthHandler      *handler.OAuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	RoleHandler       *handler.RoleHandler

	// Role-group lookup used by role resolution.
	RoleGroupRegistry *usecase.RoleGroupRegistry

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
	internalHandler *handler.InternalHandler,
	mainAreaHandler *handler.MainAreaHandler,
	purposeUseHandler *handler.PurposeUseHandler,
	professionHandler *handler.ProfessionHandler,
	tagHandler *handler.TagHandler,
	kycHandler *handler.KYCHandler,
	priceTableHandler *handler.PriceTableHandler,
	userDashboardStatsJob *job.UserDashboardStatsJob,
	organizationClient *clients.OrganizationClient,

	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	// internalHandler *handler.InternalHandler,
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	roleHandler *handler.RoleHandler,
	userInfoHandler *handler.AuthUserProfileHandler,
	roleGroupRegistry *usecase.RoleGroupRegistry,
	znsScheduler *scheduler.ZnsScheduler,
) *InitialApp {
	app := &InitialApp{
		ProfileService:          profileService,
		AdminHandler:            adminHandler,
		AdminUserProfileHandler: adminUserProfileHandler,
		BookmarkUserHandler:     bookmarkUserHandler,
		CrmProvider:             crmProvider,
		InternalHandler:         internalHandler,
		MainAreaHandler:         mainAreaHandler,
		PurposeUseHandler:       purposeUseHandler,
		ProfessionHandler:       professionHandler,
		TagHandler:              tagHandler,
		KYCHandler:              kycHandler,
		PriceTableHandler:       priceTableHandler,
		UserDashboardStatsJob:   userDashboardStatsJob,
		OrganizationClient:      organizationClient,

		AuthHandler:       authHandler,
		UserInfoHandler:   userInfoHandler,
		OAuthHandler:      oauthHandler,
		PermissionHandler: permissionHandler,
		RoleGroupHandler:  roleGroupHandler,
		RoleHandler:       roleHandler,
		RoleGroupRegistry: roleGroupRegistry,
		ZnsScheduler:      znsScheduler,
	}

	loggerName := viper.GetString("server.name")
	if loggerName == "" {
		loggerName = "user-service"
	}
	app.Logger = flogging.MustGetLogger(loggerName)

	return app
}

// LoadRoleGroupsOnStartup performs the one-shot role-group lookup warmup under the process context.
// Failure preserves the historical degraded-start behavior; the process root decides
// whether to log/alert rather than constructors hiding background work.
func (app *InitialApp) LoadRoleGroupsOnStartup(ctx context.Context) error {
	if app == nil || app.RoleGroupRegistry == nil {
		return nil
	}
	if err := app.RoleGroupRegistry.Load(ctx); err != nil {
		return err
	}
	app.Logger.Infof("Successfully loaded role groups on startup")
	return nil
}
