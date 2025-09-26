package cmd

import (
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"

	"feedback/config"
	"feedback/infra/db"
	"feedback/wire"
	feedbackpb "pb/types/feedback"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Khởi chạy gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		startGRPCServer()
	},
}

func startGRPCServer() {
	config.LoadConfig()
	// Start gRPC server
	port := fmt.Sprintf(":%d", config.AppProperties.Server.TCPPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
	)

	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	db.AutoMigrate()

	// Register services
	feedbackpb.RegisterInternalServiceServer(s, app.InternalHandler)
	feedbackpb.RegisterRateServiceServer(s, app.RateHandler)
	feedbackpb.RegisterReportServiceServer(s, app.ReportHandler)
	feedbackpb.RegisterReportReasonServiceServer(s, app.ReportReasonHandler)
	// feedbackpb.RegisterRegistryServiceServer(s, app.RegistryHandler)

	// Enable reflection for debugging
	// reflection.Register(s)

	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
