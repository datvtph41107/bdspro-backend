package grpc

import (
	"assistant/wire"
	"fmt"
	"log"
	"net"

	assistantpb "pb/types/assistant"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer() {
	// Get port from config
	grpcPort := viper.GetString("app.port.grpc")
	if grpcPort == "" {
		grpcPort = "50061"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	// Initialize dependencies
	assistantHandler, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	assistantpb.RegisterAssistantServiceServer(grpcServer, assistantHandler.AssistantHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Printf("✓ gRPC Server listening on port %s", grpcPort)
	log.Printf("✓ Assistant Service is ready to serve requests")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
