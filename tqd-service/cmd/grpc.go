package cmd

import (
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	tqdpb "pb/types/tqd"
	"syscall"
	"time"
	"tqd/config"
	"tqd/infra/db"
	"tqd/wire"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	Long: `Start the gRPC server for TQD Service.
	
The server will listen on the port specified in the configuration file.
Default port is 8230 for gRPC (instead of HTTP).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGRPCServer()
	},
}

var (
	GrpcPort int
	GrpcHost string
)

func init() {
	// gRPC specific flags
	grpcCmd.Flags().IntVarP(&GrpcPort, "port", "p", 0, "gRPC server port (overrides config)")
	grpcCmd.Flags().StringVar(&GrpcHost, "host", "", "gRPC server host (default: all interfaces)")
}

// RunGRPCServer starts the gRPC server (exported function for main.go)
func RunGRPCServer() error {
	return runGRPCServer()
}

func runGRPCServer() error {
	log.Println("Starting TQD gRPC Server...")

	// Load configuration
	cfg := config.LoadConfig()

	// Override port if specified via flag
	if GrpcPort > 0 {
		cfg.Server.GrpcPort = GrpcPort
	}

	// Set timezone
	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	// Initialize dependencies
	app, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}
	if app == nil {
		log.Fatal("app is nil after InitializeApp")
	}

	// Run database migration if needed
	db.AutoMigrate()

	// Create listener
	addr := fmt.Sprintf("%s:%d", GrpcHost, cfg.Server.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create gRPC server with interceptors
	s := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
	)

	// Register all gRPC services
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

	// Register reflection service on gRPC server (useful for grpcurl)
	reflection.Register(s)

	// Start server in a goroutine
	go func() {
		log.Printf("gRPC Server listening on %s", addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")

	// Graceful stop
	s.GracefulStop()
	log.Println("gRPC server exited")

	return nil
}
