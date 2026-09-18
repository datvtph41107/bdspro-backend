package grpc

import (
	"assistant/wire"
	"common/configloader"
	_middleware "common/middleware"
	"fmt"
	"log/slog"
	"net"
	"os"

	assistantpb "pb/types/assistant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer() {
	// Composition root resolves environment/topology and loads runtime.yml once.
	// The Wire graph then materializes typed Assistant config from that loaded state.
	selection, err := configloader.LoadRuntimeYML()
	if err != nil {
		slog.Error(
			"load assistant runtime config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	assistantHandler, cleanup, err := wire.InitializeApp()
	if err != nil {
		slog.Error(
			"initialize assistant dependencies",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer cleanup()

	grpcPort := assistantHandler.Runtime.GRPCPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		slog.Error(
			"listen for assistant gRPC",
			slog.String("port", grpcPort),
			slog.Any("error", err),
		)
		os.Exit(1)
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

	slog.Info(
		"assistant runtime config loaded",
		slog.String("config.path", "config/runtime.yml"),
		slog.String("runtime.environment", selection.Environment),
	)
	slog.Info(
		"assistant gRPC listening",
		slog.String("port", grpcPort),
	)
	slog.Info("assistant service ready")

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error(
			"serve assistant gRPC",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
