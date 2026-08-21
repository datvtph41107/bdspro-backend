package cmd_grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	authpb "pb/types/auth"
	userpb "pb/types/user"
	"user/config"
	"user/db"
	"user/env"
	"user/wire"

	_db "common/db"
	quota "common/legacyquota"
	_middleware "common/middleware"
	"common/processlifecycle"
	_redis "common/redis"
	"common/transporttrust"
	admin_middleware "user/infra/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGRPC()
	},
}

func runGRPC() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := config.LoadProperties(); err != nil {
		return fmt.Errorf("load User config: %w", err)
	}
	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	if _, err := _db.NewDB(); err != nil {
		return fmt.Errorf("connect User database: %w", err)
	}
	db.AutoMigrate()

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize User app: %w", err)
	}
	defer cleanup()
	// Long-lived actors are started only after a process cancellation domain exists.
	// app.CrmProvider.RegisterService remains disabled.

	// === QUOTA INIT ===
	var quotaInterceptor grpc.UnaryServerInterceptor
	quotaCfg, err := quota.LoadQuotaConfig()
	if err == nil && quotaCfg.Enabled {
		redisService := _redis.NewRedisService()
		storage := quota.NewRedisQuotaStorage(redisService)
		policy, err := quota.NewDefaultPolicyEngine(quotaCfg.PolicyFile)
		if err != nil {
			log.Printf("Failed to load quota policy: %v, quota disabled", err)
		} else {
			qm := quota.NewQuotaManager(storage, policy)
			extractUserID := func(ctx context.Context) string {
				if uid := ctx.Value("user_id"); uid != nil {
					switch v := uid.(type) {
					case uint64:
						return strconv.FormatUint(v, 10)
					case string:
						return v
					}
				}
				return ""
			}
			quotaInterceptor = quota.UnaryServerInterceptor(qm, extractUserID)
			log.Println("Quota interceptor enabled for user-service")
		}
	}
	if quotaInterceptor == nil {
		quotaInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	legacyQuotaInterceptor := quotaInterceptor
	quotaInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info != nil && trustedIdentityMethod(info.FullMethod) {
			return handler(ctx, req)
		}
		return legacyQuotaInterceptor(ctx, req, info, handler)
	}
	// =========================

	port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen User gRPC %s: %w", port, err)
	}

	trustConfig := transporttrust.ConfigFromEnv()
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			accessIngress(trustConfig),
			admin_middleware.AdminAuthInterceptor,
			quotaInterceptor,
		),
		grpc.ChainStreamInterceptor(
			_middleware.ParseGrpcMetadataContextStreamMiddleware,
		),
	)

	// Register services (unchanged)
	userpb.RegisterAdminProfileServiceServer(s, app.AdminHandler)
	userpb.RegisterAdminUserProfileServiceServer(s, app.AdminUserProfileHandler)
	userpb.RegisterProfileServiceServer(s, app.ProfileService)
	userpb.RegisterBookmarkUserServiceServer(s, app.BookmarkUserHandler)
	userpb.RegisterMembershipServiceServer(s, app.MembershipHandler)
	userpb.RegisterPlanServiceServer(s, app.PlanHandler)
	userpb.RegisterInternalUserServiceServer(s, app.InternalHandler)
	if err := registerTargetServices(s); err != nil {
		return targetRegistrationError(err)
	}
	userpb.RegisterMainAreaServiceServer(s, app.MainAreaHandler)
	userpb.RegisterPurposeUseServiceServer(s, app.PurposeUseHandler)
	userpb.RegisterProfessionServiceServer(s, app.ProfessionHandler)
	userpb.RegisterTagServiceServer(s, app.TagHandler)
	userpb.RegisterKYCServiceServer(s, app.KYCHandler)
	userpb.RegisterPriceTableServiceServer(s, app.PriceTableHandler)

	authpb.RegisterAuthServiceServer(s, app.AuthHandler)
	authpb.RegisterOAuthServiceServer(s, app.OAuthHandler)
	// authpb.RegisterAuthInternalServiceServer(s, app.InternalHandler)
	authpb.RegisterPermissionServiceServer(s, app.PermissionHandler)
	authpb.RegisterRoleServiceServer(s, app.RoleHandler)
	authpb.RegisterRoleGroupServiceServer(s, app.RoleGroupHandler)
	authpb.RegisterUserInfoServiceServer(s, app.UserInfoHandler)

	log.Printf("Listen: %v", port)
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	if err := app.LoadModulesOnStartup(processCtx); err != nil {
		log.Printf("User module cache warmup failed; starting degraded: %v", err)
	}

	actors := make([]processlifecycle.Actor, 0, 2)
	if app.UserDashboardStatsJob != nil {
		actors = append(actors, processlifecycle.ActorFunc(app.UserDashboardStatsJob.RunDailyAtMidnight))
	}
	if app.ZnsScheduler != nil {
		actors = append(actors, processlifecycle.ActorFunc(app.ZnsScheduler.Run))
	}

	return processlifecycle.Run(
		processCtx,
		processCancel,
		s,
		lis,
		actors,
		processlifecycle.Config{GracefulStopTimeout: 10 * time.Second},
	)
}
