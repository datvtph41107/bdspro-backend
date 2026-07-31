package cmd_grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	authpb "pb/types/auth"
	userpb "pb/types/user"
	"user/config"
	"user/db"
	"user/env"
	"user/wire"

	_db "common/db"
	_middleware "common/middleware"
	"common/quota"
	_redis "common/redis"
	admin_middleware "user/infra/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
		config.LoadProperties()
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		if _, err := _db.NewDB(); err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		db.AutoMigrate()

		app, _ := wire.InitializeApp()
		go app.UserDashboardStatsJob.RunDailyAtMidnight(context.Background())
		app.CrmClient.RegisterService(context.Background())

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
		// =========================

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				_middleware.ParseGrpcMetadataContextMiddleware,
				admin_middleware.AdminAuthInterceptor,
				quotaInterceptor,
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
		userpb.RegisterMainAreaServiceServer(s, app.MainAreaHandler)
		userpb.RegisterPurposeUseServiceServer(s, app.PurposeUseHandler)
		userpb.RegisterProfessionServiceServer(s, app.ProfessionHandler)
		userpb.RegisterTagServiceServer(s, app.TagHandler)
		userpb.RegisterKYCServiceServer(s, app.KYCHandler)
		userpb.RegisterPriceTableServiceServer(s, app.PriceTableHandler)

		authpb.RegisterAuthInternalServiceServer(s, app.AuthInternalHandler)
		authpb.RegisterAuthServiceServer(s, app.AuthHandler)
		authpb.RegisterOAuthServiceServer(s, app.OAuthHandler)
		// authpb.RegisterAuthInternalServiceServer(s, app.InternalHandler)
		authpb.RegisterPermissionServiceServer(s, app.PermissionHandler)
		authpb.RegisterRoleServiceServer(s, app.RoleHandler)
		authpb.RegisterRoleGroupServiceServer(s, app.RoleGroupHandler)
		authpb.RegisterUserInfoServiceServer(s, app.UserInfoHandler)

		log.Printf("Listen: %v", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
