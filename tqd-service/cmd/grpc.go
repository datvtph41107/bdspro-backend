package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	_db "common/db"
	_middleware "common/middleware"
	process "common/process"
	"common/rpcenv"
	tqdpb "pb/types/tqd"
	"tqd/config"
	tqddb "tqd/infra/db"
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
	log.Println("Starting TQD gRPC Server...")

	// Load config first so shared Wire providers (_db.NewDB, Redis, RPCClientSet)
	// read the same Viper values. Wire owns construction of New* dependencies.
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

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize TQD app: %w", err)
	}
	defer cleanup()
	if app == nil {
		return errors.New("TQD app is nil after InitializeApp")
	}

	schemaPolicy, err := _db.LoadSchemaPolicy("tqd")
	if err != nil {
		return fmt.Errorf("load TQD schema policy: %w", err)
	}
	log.Printf("TQD database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
	if err := _db.ApplySchemaPolicy(app.DB, schemaPolicy, tqddb.AutoMigrate); err != nil {
		return fmt.Errorf("apply TQD schema policy: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", GrpcHost, cfg.Properties.Server.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen TQD gRPC on %s: %w", addr, err)
	}

	reportTargetOperationInterceptor := buildReportTargetOperationInterceptor()
	commercialActionInterceptor := buildCommercialActionInterceptor(app.CommercialActionService)
	transportConfig := rpcenv.LoadTransportConfig()
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(grpcMaxMsgBytes),
		grpc.MaxSendMsgSize(grpcMaxMsgBytes),
		grpc.ChainUnaryInterceptor(targetIngress(transportConfig), reportTargetOperationInterceptor, commercialActionInterceptor),
		grpc.ChainStreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)
	registerTQDServices(s, app)
	reflection.Register(s)

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	closeReportTarget := func() {}
	if app.UserRPC != nil && app.UserRPC.AccessClient != nil && app.Redis != nil && app.Redis.Client != nil {
		reportCfg := cfg.Report
		if app.RuntimeConfig != nil {
			reportCfg = app.RuntimeConfig.Report
		}
		if closeTarget, targetErr := startReportTarget(
			processCtx,
			app.DB,
			app.Redis.Client,
			app.UserRPC.AccessClient,
			app.WorkspaceMapGrpcHandler,
			reportCfg,
			nil,
		); targetErr != nil {
			log.Printf("[REPORT-TARGET] capability unavailable: %v", targetErr)
		} else {
			closeReportTarget = closeTarget
		}
	}

	actors := make([]process.Actor, 0, 1)
	if app.QHPlanningClassifyWorker != nil {
		actors = append(actors, process.ActorFunc(app.QHPlanningClassifyWorker.Run))
	}

	log.Printf("TQD gRPC listening on %s", addr)
	lifecycleErr := process.Run(
		processCtx,
		processCancel,
		s,
		lis,
		actors,
		process.Config{GracefulStopTimeout: cfg.Shutdown},
	)

	closeReportTarget()
	log.Println("TQD gRPC server exited")
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
