package initial

import (
	_provider "common/domain/provider"
	"context"
	"pb/clients"
	"userv4/infra/handler"
	usergrpc "userv4/infra/handler/grpc"
	accessgrpc "userv4/infra/handler/grpc/entitlement"
	"userv4/infra/scheduler"
	"userv4/internal/job"
	usecase "userv4/modules/auth/services"
	usecases "userv4/modules/user/services/user"
	bdsprocontroller "userv4/modules/bdspro/controller"
	crmcontroller "userv4/modules/crm/controller"
	hubcontroller "userv4/modules/hub/controller"
	notificationcontroller "userv4/modules/notification/controller"
	paymentcontroller "userv4/modules/payment/controller"
	seocontroller "userv4/modules/seo/controller"
	tqdcontroller "userv4/modules/tqd/controller"
	usercontroller "userv4/modules/user/controller"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type InitialApp struct {
	DB *gorm.DB

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

	AuthHandler       *handler.AuthHandler
	UserInfoHandler   *handler.AuthUserProfileHandler
	OAuthHandler      *handler.OAuthHandler
	PermissionHandler *handler.PermissionHandler
	RoleGroupHandler  *handler.RoleGroupHandler
	RoleHandler       *handler.RoleHandler
	BdsproHandler       *bdsprocontroller.Handler
	HubHandler          *hubcontroller.Handler
	NotificationHandler *notificationcontroller.Handler
	TqdHandler          *tqdcontroller.Handler
	CrmHandler          *crmcontroller.Handler
	SeoHandler          *seocontroller.Handler
	PaymentHandler      *paymentcontroller.Handler
	UserHandler         *usercontroller.Handler

	// Role-group lookup used by role resolution.
	RoleGroupRegistry *usecase.RoleGroupRegistry

	// Trusted commercial / entitlement gRPC handlers (Wire-owned).
	InternalAccessHandler         *accessgrpc.Handler
	CommercialProfileHandler      *usergrpc.CommercialProfileHandler
	SubscriptionCheckoutHandler   *usergrpc.SubscriptionCheckoutHandler
	SubscriptionSettlementHandler *usergrpc.SubscriptionSettlementHandler
	SubscriptionAdminHandler      *usergrpc.SubscriptionAdminHandler
	PlanVersionAdminHandler       *usergrpc.PlanVersionAdminHandler

	// Scheduler để chạy background jobs
	ZnsScheduler *scheduler.ZnsScheduler

	Logger *flogging.FabricLogger
}

func NewInitialApp(
	db *gorm.DB,
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
	permissionHandler *handler.PermissionHandler,
	roleGroupHandler *handler.RoleGroupHandler,
	roleHandler *handler.RoleHandler,
	userInfoHandler *handler.AuthUserProfileHandler,
	bdsproHandler *bdsprocontroller.Handler,
	hubHandler *hubcontroller.Handler,
	notificationHandler *notificationcontroller.Handler,
	tqdHandler *tqdcontroller.Handler,
	crmHandler *crmcontroller.Handler,
	seoHandler *seocontroller.Handler,
	paymentHandler *paymentcontroller.Handler,
	userHandler *usercontroller.Handler,
	roleGroupRegistry *usecase.RoleGroupRegistry,
	znsScheduler *scheduler.ZnsScheduler,

	internalAccessHandler *accessgrpc.Handler,
	commercialProfileHandler *usergrpc.CommercialProfileHandler,
	subscriptionCheckoutHandler *usergrpc.SubscriptionCheckoutHandler,
	subscriptionSettlementHandler *usergrpc.SubscriptionSettlementHandler,
	subscriptionAdminHandler *usergrpc.SubscriptionAdminHandler,
	planVersionAdminHandler *usergrpc.PlanVersionAdminHandler,
) *InitialApp {
	app := &InitialApp{
		DB:                      db,
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
		BdsproHandler:       bdsproHandler,
		HubHandler:          hubHandler,
		NotificationHandler: notificationHandler,
		TqdHandler:          tqdHandler,
		CrmHandler:          crmHandler,
		SeoHandler:          seoHandler,
		PaymentHandler:      paymentHandler,
		UserHandler:         userHandler,
		RoleGroupRegistry:   roleGroupRegistry,
		ZnsScheduler:      znsScheduler,

		InternalAccessHandler:         internalAccessHandler,
		CommercialProfileHandler:      commercialProfileHandler,
		SubscriptionCheckoutHandler:   subscriptionCheckoutHandler,
		SubscriptionSettlementHandler: subscriptionSettlementHandler,
		SubscriptionAdminHandler:      subscriptionAdminHandler,
		PlanVersionAdminHandler:       planVersionAdminHandler,
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
