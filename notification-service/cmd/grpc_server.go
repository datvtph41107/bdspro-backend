package cmd

import (
	"fmt"
	"log"
	"net"
	"notification/config"
	"notification/db"
	"notification/wire"
	notificationpb "pb/types/notification"

	_middleware "common/middleware"

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

	app, _ := wire.InitializeApp()

	db.AutoMigrate()

	notificationpb.RegisterGatewayNotificationServiceServer(s, app.NotificationHandler)
	notificationpb.RegisterInternalNotificationServiceServer(s, app.InternalHandler)
	notificationpb.RegisterHistoryServiceServer(s, app.HistoryHandler)
	notificationpb.RegisterDealHistoryServiceServer(s, app.DealHistoryHandler)
	notificationpb.RegisterPropertyHistoryServiceServer(s, app.PropertyHistoryHandler)

	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
