package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tqd/config"
	"tqd/infra/db"
	"tqd/wire"

	_middleware "common/middleware"
	"common/quota"
	_redis "common/redis"

	tqdpb "pb/types/tqd"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

// grpcCmd represents the grpc command
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

func RunGRPCServer() error {
	return runGRPCServer()
}

func runGRPCServer() error {
	log.Println("Starting TQD gRPC Server...")

	// Load configuration
	cfg := config.LoadConfig()
	if GrpcPort > 0 {
		cfg.Server.GrpcPort = GrpcPort
	}

	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	app, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}
	if app == nil {
		log.Fatal("app is nil after InitializeApp")
	}

	db.AutoMigrate()

	addr := fmt.Sprintf("%s:%d", GrpcHost, cfg.Server.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// =================== QUOTA INTERCEPTOR INIT ===================
	var quotaInterceptor grpc.UnaryServerInterceptor

	// Khởi tạo quota logger
	if err := quota.InitQuotaLogger("./app/logs/quota", true); err != nil {
		log.Printf("Failed to init quota logger: %v", err)
	}

	quotaCfg, err := quota.LoadQuotaConfig()
	if err == nil && quotaCfg != nil && quotaCfg.Enabled {
		log.Println("[QUOTA] Initializing quota interceptor...")

		redisService := _redis.NewRedisService()
		storage := quota.NewRedisQuotaStorage(redisService)
		policy, err := quota.NewDefaultPolicyEngine(quotaCfg.PolicyFile)
		if err != nil {
			log.Printf("[QUOTA] Failed to load policy: %v, quota disabled", err)
		} else {
			qm := quota.NewQuotaManager(storage, policy)

			extractUserID := func(ctx context.Context) string {
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					// Ưu tiên profileId
					if vals := md.Get("profileId"); len(vals) > 0 && vals[0] != "" {
						return vals[0]
					}
					if vals := md.Get("originId"); len(vals) > 0 && vals[0] != "" {
						return vals[0]
					}
					if vals := md.Get("authId"); len(vals) > 0 && vals[0] != "" {
						return vals[0]
					}
				}
				return ""
			}

			quotaInterceptor = quota.UnaryServerInterceptor(qm, extractUserID)
			log.Println("✅ Quota interceptor enabled for tqd-service")
		}
	} else {
		log.Println("[QUOTA] Quota disabled or config error")
	}

	if quotaInterceptor == nil {
		quotaInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
			quotaInterceptor,
		),
	)

	// Register services
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
	tqdpb.RegisterRegionServiceServer(s, app.RegionGrpcHandler)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC Server listening on %s", addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("gRPC server exited")

	return nil
}
