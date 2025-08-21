package cmd_grpc

import (
	"common"
	"fmt"
	"log"
	"net"
	"os"
	userpb "pb/types/user"
	"time"
	"user/config"
	"user/db"
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
		common.LoadConfig(config.ConfigFile)
		// _db.NewDB()

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		app, _ := wire.InitializeApp()
		db.AutoMigrate()

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
		userpb.RegisterProfileServiceServer(s, app.ProfileService)
		userpb.RegisterAuthServiceServer(s, app.AuthService)
		userpb.RegisterRateServiceServer(s, app.RateService)
		// V2 Services (protobuf-based)
		userpb.RegisterReportServiceServer(s, app.ReportHandler)
		userpb.RegisterReportReasonServiceServer(s, app.ReportReasonHandler)
		userpb.RegisterReportProofServiceServer(s, app.ReportProofHandler)
		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
