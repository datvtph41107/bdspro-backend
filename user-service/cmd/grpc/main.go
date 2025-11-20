package cmd_grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	userpb "pb/types/user"
	"time"

	"user/config"
	"user/wire"

	_middleware "common/middleware"
	admin_middleware "user/infra/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadProperties()
		// common.LoadConfig(config.EmbeddedConfig)
		// _db.NewDB()

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		app, _ := wire.InitializeApp()
		// db.AutoMigrate()

		go app.UserDashboardStatsJob.RunDailyAtMidnight(context.Background())
		// Đăng ký service vào crm registry
		app.CrmClient.RegisterService(context.Background())

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				_middleware.ParseGrpcMetadataContextMiddleware,
				admin_middleware.AdminAuthInterceptor,
			),
		)

		userpb.RegisterAdminProfileServiceServer(s, app.AdminHandler)
		userpb.RegisterAdminUserProfileServiceServer(s, app.AdminUserProfileHandler)
		userpb.RegisterProfileServiceServer(s, app.ProfileService)
		// V2 Services (protobuf-based)
		userpb.RegisterBookmarkUserServiceServer(s, app.BookmarkUserHandler)
		userpb.RegisterMembershipServiceServer(s, app.MembershipHandler)
		userpb.RegisterPlanServiceServer(s, app.PlanHandler)
		userpb.RegisterInternalUserServiceServer(s, app.InternalHandler)
		userpb.RegisterMainAreaServiceServer(s, app.MainAreaHandler)
		userpb.RegisterPurposeUseServiceServer(s, app.PurposeUseHandler)
		userpb.RegisterProfessionServiceServer(s, app.ProfessionHandler)
		userpb.RegisterTagServiceServer(s, app.TagHandler)
		userpb.RegisterKYCServiceServer(s, app.KYCHandler)
		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
