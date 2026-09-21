package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"strings"
	"syscall"

	_db "common/db"
	_middleware "common/middleware"
	process "common/process"
	_redis "common/redis"
	"common/rpcenv"
	assistantpb "pb/types/assistant"
	tqdpb "pb/types/tqd"
	"tqd/config"
	tqddb "tqd/infra/db"
	infra_rpc "tqd/infra/rpc"
	"tqd/initial"
	"tqd/wire"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpcMaxMsgBytes — default gRPC ~4MB; ImportGeoJSON sends large file_content
// from gateway. Keep this aligned with gateway-service TQD client limits.
const grpcMaxMsgBytes = 512 << 20

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGRPCServer()
	},
}

var (
	GrpcPort int
	GrpcHost string
)

func init() {
	grpcCmd.Flags().IntVarP(&GrpcPort, "port", "p", 0, "gRPC server port (overrides config)")
	grpcCmd.Flags().StringVar(&GrpcHost, "host", "", "gRPC server host (default: all interfaces)")
}

func RunGRPCServer() error { return runGRPCServer() }

func runGRPCServer() error {
	slog.Info(strings.TrimSuffix(fmt.Sprintln("Starting TQD gRPC Server..."), "\n"))

	cfg, err := config.LoadRuntimeConfig()
	if err != nil {
		return fmt.Errorf("load TQD runtime config: %w", err)
	}
	if cfg == nil || cfg.Properties == nil {
		return errors.New("TQD runtime config is incomplete")
	}
	if GrpcPort > 0 {
		cfg.Properties.Server.GrpcPort = GrpcPort
	}

	// Process root owns all long-lived resources. Wire composes only business
	// adapters/usecases/handlers from these already-open resources.
	database, err := _db.Open(_db.DatabaseConfig{DSN: cfg.DatabaseDSN})
	if err != nil {
		return fmt.Errorf("open TQD postgres: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("access TQD postgres pool: %w", err)
	}
	defer sqlDB.Close()

	schemaPolicy, err := _db.LoadSchemaPolicy("tqd")
	if err != nil {
		return fmt.Errorf("load TQD schema policy: %w", err)
	}
	slog.Info(fmt.Sprintf("TQD database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source))
	if err := _db.ApplySchemaPolicy(database, schemaPolicy, tqddb.AutoMigrate); err != nil {
		return fmt.Errorf("apply TQD schema policy: %w", err)
	}

	redisService, err := _redis.Open(_redis.Config{
		Address: cfg.Redis.Address, Password: cfg.Redis.Password, DB: cfg.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("open TQD redis: %w", err)
	}
	defer redisService.Close()
	if redisService.Client == nil {
		return errors.New("TQD redis client is unavailable")
	}

	userRPC, closeUser, err := infra_rpc.OpenUserRPCClient(cfg.RPC.User)
	if err != nil {
		return fmt.Errorf("open User gRPC client: %w", err)
	}
	defer closeUser()

	authRPC, closeAuth, err := infra_rpc.OpenAuthRPCClient(cfg.RPC.Auth)
	if err != nil {
		return fmt.Errorf("open Auth gRPC client: %w", err)
	}
	defer closeAuth()

	var assistantRPC assistantpb.AssistantServiceClient
	closeAssistant := func() {}
	if cfg.Classify.Enabled {
		assistantRPC, closeAssistant, err = infra_rpc.OpenAssistantRPCClient(cfg.RPC.Assistant)
		if err != nil {
			return fmt.Errorf("open Assistant gRPC client: %w", err)
		}
	}
	defer closeAssistant()

	app, cleanup, err := wire.InitializeApp(
		database,
		redisService,
		userRPC,
		authRPC,
		assistantRPC,
		cfg,
	)
	if err != nil {
		return fmt.Errorf("initialize TQD app: %w", err)
	}
	defer cleanup()
	if app == nil {
		return errors.New("TQD app is nil after InitializeApp")
	}

	addr := fmt.Sprintf("%s:%d", GrpcHost, cfg.Properties.Server.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen TQD gRPC on %s: %w", addr, err)
	}

	reportTargetOperationInterceptor := buildReportTargetOperationInterceptor()
	transportConfig := rpcenv.LoadTransportConfig()
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(grpcMaxMsgBytes),
		grpc.MaxSendMsgSize(grpcMaxMsgBytes),
		grpc.ChainUnaryInterceptor(targetIngress(transportConfig), _middleware.UnaryErrorInterceptor(), reportTargetOperationInterceptor),
		grpc.ChainStreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware, _middleware.StreamErrorInterceptor()),
	)
	registerTQDServices(s, app)
	reflection.Register(s)

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	closeReportTarget := func() {}
	if userRPC != nil && userRPC.AccessClient != nil {
		if closeTarget, targetErr := startReportTarget(
			processCtx,
			database,
			redisService.Client,
			userRPC.AccessClient,
			app.WorkspaceMapGrpcHandler,
			cfg.Report,
			nil,
		); targetErr != nil {
			slog.
				// Generated Report is capability-local. Planning/GIS/discovery remain
				// available while report creation fails closed.
				Warn(fmt.Sprintf("[REPORT-TARGET] capability unavailable: %v", targetErr))
		} else {
			closeReportTarget = closeTarget
		}
	}

	actors := make([]process.Actor, 0, 1)
	if app.QHPlanningClassifyWorker != nil {
		actors = append(actors, process.ActorFunc(app.QHPlanningClassifyWorker.Run))
	}
	slog.Info(fmt.Sprintf("TQD gRPC listening on %s", addr))
	lifecycleErr := process.Run(
		processCtx,
		processCancel,
		s,
		lis,
		actors,
		process.Config{GracefulStopTimeout: cfg.Shutdown},
	)

	// process.Run cancels and joins registered actors. Report owns additional
	// runners and exposes its own join barrier before process resources close.
	closeReportTarget()
	slog.Info(strings.TrimSuffix(fmt.Sprintln("TQD gRPC server exited"), "\n"))
	if lifecycleErr != nil {
		return fmt.Errorf("TQD gRPC lifecycle: %w", lifecycleErr)
	}
	return nil
}

func registerTQDServices(s *grpc.Server, app *initial.InitialApp) {
	tqdpb.RegisterAmenityServiceServer(s, app.AmenityGrpcHandler)
	tqdpb.RegisterContactLabelServiceServer(s, app.ContactLabelGrpcHandler)
	tqdpb.RegisterDirectoryCategoryServiceServer(s, app.DirectoryCategoryGrpcHandler)
	tqdpb.RegisterDirectorySourceServiceServer(s, app.DirectorySourceGrpcHandler)
	tqdpb.RegisterDirectorySupplierServiceServer(s, app.DirectorySupplierGrpcHandler)
	tqdpb.RegisterOpenHourServiceServer(s, app.OpenHourGrpcHandler)
	tqdpb.RegisterPoiCategoryServiceServer(s, app.PoiCategoryGrpcHandler)
	tqdpb.RegisterPoiServiceServer(s, app.POIGrpcHandler)
	tqdpb.RegisterLocationServiceServer(s, app.LocationGrpcHandler)
	tqdpb.RegisterFeatureServiceServer(s, app.FeatureGrpcHandler)
	tqdpb.RegisterParcelServiceServer(s, app.ParcelGrpcHandler)
	tqdpb.RegisterOneHouseServiceServer(s, app.OneHouseGrpcHandler)
	tqdpb.RegisterLayerServiceServer(s, app.LayerGrpcHandler)
	tqdpb.RegisterQHLabelServiceServer(s, app.QHLabelGrpcHandler)
	tqdpb.RegisterQHAuthorityIssuringServiceServer(s, app.QHAuthorityIssuringGrpcHandler)
	tqdpb.RegisterQHLayerFamilyServiceServer(s, app.QHLayerFamilyGrpcHandler)
	tqdpb.RegisterQHLayerLegendServiceServer(s, app.QHLayerLegendGrpcHandler)
	tqdpb.RegisterQHLandUseServiceServer(s, app.QHLandUseGrpcHandler)
	tqdpb.RegisterRegionServiceServer(s, app.RegionGrpcHandler)
	tqdpb.RegisterRegionExtendServiceServer(s, app.RegionExtendGrpcHandler)
	tqdpb.RegisterImportServiceServer(s, app.ImportGrpcHandler)
	tqdpb.RegisterReportServiceServer(s, app.ReportGrpcHandler)
	tqdpb.RegisterGisReportAdminServiceServer(s, app.GisReportAdminGrpcHandler)
	tqdpb.RegisterAdminUsageServiceServer(s, app.AdminUsageGrpcHandler)
	tqdpb.RegisterProfileUsageServiceServer(s, app.ProfileUsageGrpcHandler)
	tqdpb.RegisterSubscriptionServiceServer(s, app.SubscriptionGrpcHandler)
	tqdpb.RegisterTqdPublicServiceServer(s, app.TqdPublicService)
	tqdpb.RegisterTqdSeoProjectionServiceServer(s, app.SeoProjectionGrpcHandler)
	tqdpb.RegisterPlanningClientServiceServer(s, app.PlanningClientGrpcHandler)
	tqdpb.RegisterLayerLegalServiceServer(s, app.LayerLegalGrpcHandler)
	tqdpb.RegisterMapWorkspaceServiceServer(s, app.WorkspaceMapGrpcHandler)
	tqdpb.RegisterMapPointServiceServer(s, app.MapPointGrpcHandler)
	tqdpb.RegisterQHPlanningServiceServer(s, app.QHPlanningGrpcHandler)
	tqdpb.RegisterDiscoveryServiceServer(s, app.DiscoveryGrpcHandler)
	tqdpb.RegisterRelatedEntityServiceServer(s, app.RelatedEntityGrpcHandler)
}
