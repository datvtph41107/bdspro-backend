package main

import (
	"common/identitygrpc"
	"common/requestgrpc"
	"common/transporttrust"
	"context"
	"errors"
	filegrpc "file/infra/handler/filegrpc"
	"file/internal/filecontent"
	"file/internal/processlifecycle"
	"file/services"
	"fmt"
	"net"
	filepb "pb/types/file"
	"time"

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

	stopDone := make(chan struct{})
	go func() {
		defer close(stopDone)
		<-ctx.Done()
		processlifecycle.StopGracefully(server, 10*time.Second)
	}()

	err = server.Serve(listener)
	if ctx.Err() != nil && errors.Is(err, grpc.ErrServerStopped) {
		err = nil
	}
	if ctx.Err() != nil {
		select {
		case <-stopDone:
		case <-time.After(10 * time.Second):
		}
	}
	return err
}
