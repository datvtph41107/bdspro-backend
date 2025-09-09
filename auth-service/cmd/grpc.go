package cmd

import (
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"
	authpb "pb/types/auth"

	"auth/config"
	"auth/infra/db"
	"auth/wire"

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
	port := fmt.Sprintf(":%d", config.Properties.Server.TCPPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
	)

	app, _ := wire.InitializeApp()

	db.AutoMigrate()

	authpb.RegisterAuthServiceServer(s, app.AuthHandler)
	authpb.RegisterOAuthServiceServer(s, app.OAuthHandler)
	authpb.RegisterAuthInternalServiceServer(s, app.InternalHandler)
	authpb.RegisterPermissionServiceServer(s, app.PermissionHandler)

	// Enable reflection for debugging
	// reflection.Register(s)

	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
