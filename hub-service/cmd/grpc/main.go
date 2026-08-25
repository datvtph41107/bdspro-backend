package cmd_grpc

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"hub/config"
	"hub/infra/db"
	hubMiddleware "hub/infra/middleware"
	"hub/wire"
	"log"
	"net"
	"os"
	hubpb "pb/types/hub"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()

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
		interceptors := grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
			hubMiddleware.NewAPIKeyUnaryServerInterceptor(app.ApiKeyUsecase),
			_middleware.UnaryRecoveryInterceptor(app.Logger),
		)

		// Load system config vào memory khi startup
		if app.SystemConfigUsecase != nil {
			if err := app.SystemConfigUsecase.InitializeSystemConfig(context.Background()); err != nil {
				log.Printf("Warning: Failed to initialize system config: %v", err)
			}
		}

		s := grpc.NewServer(
			interceptors,
			grpc.ChainStreamInterceptor(
				_middleware.ParseGrpcMetadataContextStreamMiddleware,
			),
		)
		hubpb.RegisterEventQueueServiceServer(s, app.EventQueueService)
		hubpb.RegisterHubInternalServiceServer(s, app.InternalHandler)
		hubpb.RegisterLocationServiceServer(s, app.LocationHandler)
		hubpb.RegisterLocationV2ServiceServer(s, app.LocationV2Handler)
		hubpb.RegisterUserGuideServiceServer(s, app.UserGuideHandler)
		hubpb.RegisterSystemConfigServiceServer(s, app.SystemConfigHandler)
		hubpb.RegisterFAQServiceServer(s, app.FAQHandler)
		hubpb.RegisterVersionServiceServer(s, app.VersionHandler)
		hubpb.RegisterApiKeyServiceServer(s, app.ApiKeyHandler)
		hubpb.RegisterInteractiveEventServiceServer(s, app.InteractiveEventHandler)
		hubpb.RegisterClickURLServiceServer(s, app.ClickURLHandler)
		hubpb.RegisterErrorLogServiceServer(s, app.ErrorLogHandler)
		hubpb.RegisterUpdateDataServiceServer(s, app.UpdateDataHandler)
		hubpb.RegisterApplinkServiceServer(s, app.ApplinkHandler)
		hubpb.RegisterRequestStatsServiceServer(s, app.RequestStatsHandler)

		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	},
}
