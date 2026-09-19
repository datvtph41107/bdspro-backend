package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_db "common/db"
	"common/logging"
	_middleware "common/middleware"
	"common/process"
	_redis "common/redis"
	"common/rpcenv"
	authpb "pb/types/auth"
	userpb "pb/types/user"
	"user/config"
	userdb "user/db"
	"user/wire"

	"google.golang.org/grpc"
)

func main() {
	os.Exit(runProcess())
}

func runProcess() int {
	closeLogger, err := logging.Configure("user-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure User logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := runGRPC(); err != nil {
		slog.Error(
			"user service failed",
			slog.Any("error", err),
		)
		return 1
	}

	return 0
}

func runGRPC() error {
	runtimeConfig, err := config.LoadRuntime()
	if err != nil {
		return fmt.Errorf("load User runtime config: %w", err)
	}

	runtimePolicy, err := loadRuntimePolicy()
	if err != nil {
		return fmt.Errorf("load User runtime policy: %w", err)
	}
	slog.Info(
		"user runtime policy",
		slog.Bool("module_warmup", runtimePolicy.ModuleWarmup),
		slog.Bool("dashboard_stats", runtimePolicy.DashboardStats),
		slog.Bool("zns_scheduler", runtimePolicy.ZNSScheduler),
	)

	database, err := _db.Open(_db.DatabaseConfig{
		DSN:             runtimeConfig.Database.DSN,
		MaxOpenConns:    runtimeConfig.Database.MaxOpenConns,
		MaxIdleConns:    runtimeConfig.Database.MaxIdleConns,
		ConnMaxLifetime: runtimeConfig.Database.ConnMaxLifetime,
	})
	if err != nil {
		return fmt.Errorf("open User database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("get User database pool: %w", err)
	}
	defer sqlDB.Close()

	schemaPolicy, err := _db.LoadSchemaPolicy("user")
	if err != nil {
		return fmt.Errorf("load User schema policy: %w", err)
	}
	slog.Info(
		"user database schema policy",
		slog.String("mode", string(schemaPolicy.Mode)),
		slog.String("source", schemaPolicy.Source),
	)
	if err := _db.ApplySchemaPolicy(database, schemaPolicy, userdb.AutoMigrate); err != nil {
		return fmt.Errorf("apply User schema policy: %w", err)
	}

	redisService, err := _redis.Open(_redis.Config{
		Address:  runtimeConfig.Redis.Address,
		Password: runtimeConfig.Redis.Password,
		DB:       runtimeConfig.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("open User Redis: %w", err)
	}
	defer redisService.Close()

	app, cleanup, err := wire.InitializeApp(database, redisService)
	if err != nil {
		return fmt.Errorf("initialize User app: %w", err)
	}
	defer cleanup()

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
	targetCleanup, err := registerTargetServices(server, database)
	if err != nil {
		return targetRegistrationError(err)
	}
	defer targetCleanup()
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
			slog.WarnContext(
				processCtx,
				"user module cache warmup failed",
				slog.Any("error", err),
			)
		}
	}

	actors := make([]process.Actor, 0, 2)
	if runtimePolicy.DashboardStats && app.UserDashboardStatsJob != nil {
		actors = append(actors, process.ActorFunc(app.UserDashboardStatsJob.RunDailyAtMidnight))
	}
	if runtimePolicy.ZNSScheduler && app.ZnsScheduler != nil {
		actors = append(actors, process.ActorFunc(app.ZnsScheduler.Run))
	}

	slog.InfoContext(
		processCtx,
		"user gRPC listening",
		slog.String("grpc.address", runtimeConfig.GRPCAddress),
	)
	return process.Run(
		processCtx,
		processCancel,
		server,
		lis,
		actors,
		process.Config{GracefulStopTimeout: 10 * time.Second},
	)
}
