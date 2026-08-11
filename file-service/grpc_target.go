package main

import (
	"common/identitygrpc"
	"common/requestgrpc"
	"common/transporttrust"
	"context"
	filegrpc "file/infra/handler/filegrpc"
	"file/internal/filecontent"
	"file/services"
	"fmt"
	"net"
	"os"
	filepb "pb/types/file"
	"strings"

	"google.golang.org/grpc"
)

const fileGRPCMaxMessageBytes = filecontent.MaxContentBytes + (1 << 20)

/**
 * serveFileGRPCTarget chạy API nội bộ của file-service bằng gRPC.
 *
 * FileService ở tham số là bridge sang storage/repository hiện tại.
 * Sau migration, HTTP business routes của file-service sẽ được bỏ khỏi runtime
 * và external HTTP chỉ đi qua gateway-service.
 */
func serveFileGRPCTarget(
	ctx context.Context,
	port string,
	legacy *services.FileService,
) error {
	if legacy == nil {
		return fmt.Errorf("file service is nil")
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	files := filecontent.NewService(
		legacy.FileRepository,
		legacy.AccessRepository,
		legacy.Storage,
		legacy.XorCryptKey,
	)
	trust := transporttrust.ConfigFromEnv()
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(fileGRPCMaxMessageBytes),
		grpc.MaxSendMsgSize(fileGRPCMaxMessageBytes),
		grpc.ChainUnaryInterceptor(
			requestgrpc.UnaryServerInterceptor,
			identitygrpc.UnaryServerInterceptor(trust),
		),
	)
	filepb.RegisterInternalFileServiceServer(server, filegrpc.NewHandler(files))

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	return server.Serve(listener)
}

func fileGRPCTargetEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_FILE_GRPC_ENABLED")))
	return value == "1" || value == "true" || value == "yes"
}

func fileGRPCTargetPort() string {
	if port := strings.TrimSpace(os.Getenv("QHPRO_FILE_GRPC_PORT")); port != "" {
		return port
	}
	return "8203"
}
