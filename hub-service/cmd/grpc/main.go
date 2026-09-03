package cmd_grpc

import (
	common_db "common/db"
	_middleware "common/middleware"
	process "common/process"
	"context"
	"fmt"
	"hub/config"
	"hub/infra/db"
	hubMiddleware "hub/infra/middleware"
	"hub/wire"
	"log"
	"net"
	"os"
	"os/signal"
	hubpb "pb/types/hub"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.LoadConfig(); err != nil {
			return fmt.Errorf("load hub config: %w", err)
		}

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		app, cleanup, err := wire.InitializeApp()
		if err != nil {
			return fmt.Errorf("initialize hub app: %w", err)
		}
		defer cleanup()
		schemaPolicy, err := common_db.LoadSchemaPolicy("hub")
		if err != nil {
			return fmt.Errorf("load Hub schema policy: %w", err)
		}
		log.Printf("Hub database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
		if err := common_db.ApplySchemaPolicy(common_db.DB, schemaPolicy, db.AutoMigrate); err != nil {
			return fmt.Errorf("apply Hub schema policy: %w", err)
		}

		signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stopSignals()
		processCtx, processCancel := context.WithCancel(signalCtx)
		defer processCancel()

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			return fmt.Errorf("listen hub gRPC on %s: %w", port, err)
		}
		defer func() { _ = lis.Close() }()
		interceptors := grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
			hubMiddleware.NewAPIKeyUnaryServerInterceptor(app.ApiKeyUsecase),
			_middleware.UnaryRecoveryInterceptor(app.Logger),
		)

		// Load system config vào memory khi startup
		if app.SystemConfigUsecase != nil {
			if err := app.SystemConfigUsecase.InitializeSystemConfig(processCtx); err != nil {
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
		hubpb.RegisterErrorLogServiceServer(s, app.ErrorLogHandler)
		hubpb.RegisterUpdateDataServiceServer(s, app.UpdateDataHandler)
		hubpb.RegisterApplinkServiceServer(s, app.ApplinkHandler)

		log.Printf("Listen: %v", port)

		return process.Run(
			processCtx, processCancel, s, lis, nil,
			process.Config{GracefulStopTimeout: 10 * time.Second},
		)
	},
}
