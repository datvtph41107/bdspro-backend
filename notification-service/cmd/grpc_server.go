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
	RunE: func(cmd *cobra.Command, args []string) error {
		return startGRPCServer()
	},
}

func startGRPCServer() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load notification config: %w", err)
	}
	// Start gRPC server
	port := fmt.Sprintf(":%d", config.AppProperties.Server.TCPPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen notification gRPC on %s: %w", port, err)
	}
	defer func() { _ = lis.Close() }()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		grpc.StreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize notification app: %w", err)
	}
	defer cleanup()

	db.AutoMigrate()

	notificationpb.RegisterGatewayNotificationServiceServer(s, app.NotificationHandler)
	notificationpb.RegisterInternalNotificationServiceServer(s, app.InternalHandler)
	notificationpb.RegisterHistoryServiceServer(s, app.HistoryHandler)
	notificationpb.RegisterDealHistoryServiceServer(s, app.DealHistoryHandler)
	notificationpb.RegisterPropertyHistoryServiceServer(s, app.PropertyHistoryHandler)

	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("serve notification gRPC: %w", err)
	}

	return nil
}
