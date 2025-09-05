package cmd

import (
	"fmt"
	"log"
	"net"
	"noti/config"
	"noti/db"
	"noti/wire"
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
	config.LoadProperties()
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

	notificationpb.RegisterGatewayNotificationServiceServer(s, app.ServerHttp)
	notificationpb.RegisterNotificationServiceServer(s, app.InternalNotiService)
	notificationpb.RegisterHistoryServiceServer(s, app.HistoryService)
	// Enable reflection for debugging
	// reflection.Register(s)

	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
