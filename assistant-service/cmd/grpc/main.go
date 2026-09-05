package grpc

import (
	"assistant/wire"
	"common/configloader"
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"

	assistantpb "pb/types/assistant"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer() {
	// Composition root resolve environment/topology đúng một lần. Business
	// provider chỉ đọc giá trị đã được Viper nạp từ profile này.
	selection, err := configloader.LoadRuntimeYML()
	if err != nil {
		log.Fatalf("load assistant runtime config: %v", err)
	}
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
	assistantHandler, cleanup, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	// Create gRPC server — max msg đủ cho ClassifyDocumentWithGemini (file inline ~15MB).
	const grpcMaxMsgBytes = 32 << 20
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(grpcMaxMsgBytes),
		grpc.MaxSendMsgSize(grpcMaxMsgBytes),
		grpc.ChainUnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		grpc.ChainStreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)

	// Register services
	assistantpb.RegisterAssistantServiceServer(grpcServer, assistantHandler.AssistantHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	log.Printf("✓ Runtime config loaded from config/runtime.yml (environment=%s)", selection.Environment)
	log.Printf("✓ gRPC Server listening on port %s", grpcPort)
	log.Printf("✓ Assistant Service is ready to serve requests")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
