package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	_db "common/db"
	_middleware "common/middleware"
	"common/process"
	"common/rpcenv"
	authpb "pb/types/auth"
	userpb "pb/types/user"
	"user/config"
	userdb "user/db"
	"user/wire"

	"google.golang.org/grpc"
)

func main() {
	if err := runGRPC(); err != nil {
		log.Fatal(err)
	}
}

func runGRPC() error {
	runtimeConfig, err := config.LoadRuntime()
	if err != nil {
		return fmt.Errorf("load User runtime config: %w", err)
	}
	log.Printf(
		"User database target: %s",
		config.DescribeDatabaseDSN(runtimeConfig.Database.DSN),
	)

	runtimePolicy, err := loadRuntimePolicy()
	if err != nil {
		return fmt.Errorf("load User runtime policy: %w", err)
	}
	log.Printf(
		"User runtime policy: module_warmup=%t dashboard_stats=%t zns_scheduler=%t",
		runtimePolicy.ModuleWarmup,
		runtimePolicy.DashboardStats,
		runtimePolicy.ZNSScheduler,
	)

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize User app: %w", err)
	}
	defer cleanup()
	if app == nil || app.DB == nil {
		return fmt.Errorf("initialize User app: app or database is nil")
	}

	schemaPolicy, err := _db.LoadSchemaPolicy("user")
	if err != nil {
		return fmt.Errorf("load User schema policy: %w", err)
	}
	log.Printf("User database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
	if err := _db.ApplySchemaPolicy(app.DB, schemaPolicy, userdb.AutoMigrate); err != nil {
		return fmt.Errorf("apply User schema policy: %w", err)
	}

	lis, err := net.Listen("tcp", runtimeConfig.GRPCAddress)
	if err != nil {
		return fmt.Errorf("listen User gRPC %s: %w", runtimeConfig.GRPCAddress, err)
	}

	transportConfig := rpcenv.LoadTransportConfig()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			accessIngress(transportConfig),
		),
		grpc.ChainStreamInterceptor(
			_middleware.ParseGrpcMetadataContextStreamMiddleware,
		),
	)

	userpb.RegisterAdminProfileServiceServer(server, app.AdminHandler)
	userpb.RegisterAdminUserProfileServiceServer(server, app.AdminUserProfileHandler)
	userpb.RegisterProfileServiceServer(server, app.ProfileService)
	userpb.RegisterBookmarkUserServiceServer(server, app.BookmarkUserHandler)
	userpb.RegisterInternalUserServiceServer(server, app.InternalHandler)
	userpb.RegisterInternalAccessServiceServer(server, app.InternalAccessHandler)
	userpb.RegisterCommercialProfileServiceServer(server, app.CommercialProfileHandler)
	userpb.RegisterCheckoutServiceServer(server, app.SubscriptionCheckoutHandler)
	userpb.RegisterInternalSubscriptionServiceServer(server, app.SubscriptionSettlementHandler)
	userpb.RegisterAdminCommercialServiceServer(server, app.SubscriptionAdminHandler)
	userpb.RegisterAdminCatalogServiceServer(server, app.PlanVersionAdminHandler)
	userpb.RegisterMainAreaServiceServer(server, app.MainAreaHandler)
	userpb.RegisterPurposeUseServiceServer(server, app.PurposeUseHandler)
	userpb.RegisterProfessionServiceServer(server, app.ProfessionHandler)
	userpb.RegisterTagServiceServer(server, app.TagHandler)
	userpb.RegisterKYCServiceServer(server, app.KYCHandler)
	userpb.RegisterPriceTableServiceServer(server, app.PriceTableHandler)

	authpb.RegisterAuthServiceServer(server, app.AuthHandler)
	authpb.RegisterOAuthServiceServer(server, app.OAuthHandler)
	authpb.RegisterPermissionServiceServer(server, app.PermissionHandler)
	authpb.RegisterRoleServiceServer(server, app.RoleHandler)
	authpb.RegisterRoleGroupServiceServer(server, app.RoleGroupHandler)
	authpb.RegisterUserInfoServiceServer(server, app.UserInfoHandler)

	signalCtx, stopSignals := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	if runtimePolicy.ModuleWarmup {
		if err := app.LoadRoleGroupsOnStartup(processCtx); err != nil {
			log.Printf("User module cache warmup failed; starting degraded: %v", err)
		}
	}

	actors := make([]process.Actor, 0, 2)
	if runtimePolicy.DashboardStats && app.UserDashboardStatsJob != nil {
		actors = append(actors, process.ActorFunc(app.UserDashboardStatsJob.RunDailyAtMidnight))
	}
	if runtimePolicy.ZNSScheduler && app.ZnsScheduler != nil {
		actors = append(actors, process.ActorFunc(app.ZnsScheduler.Run))
	}

	log.Printf("User gRPC listening on %s", runtimeConfig.GRPCAddress)
	return process.Run(
		processCtx,
		processCancel,
		server,
		lis,
		actors,
		process.Config{GracefulStopTimeout: 10 * time.Second},
	)
}
