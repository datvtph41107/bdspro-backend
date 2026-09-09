package grpc

import (
	"assistant/wire"
	"common/configloader"
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"

	assistantpb "pb/types/assistant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer() {
	// Composition root resolves environment/topology and loads runtime.yml once.
	// The Wire graph then materializes typed Assistant config from that loaded state.
	selection, err := configloader.LoadRuntimeYML()
	if err != nil {
		log.Fatalf("load assistant runtime config: %v", err)
	}

	assistantHandler, cleanup, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cleanup()

	grpcPort := assistantHandler.Runtime.GRPCPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

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
