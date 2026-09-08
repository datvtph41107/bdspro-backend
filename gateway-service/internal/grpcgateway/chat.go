package grpcgateway

import (
	"context"
	chatpb "pb/types/chat"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterChatService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		chatpb.RegisterChatServiceHandler,
		chatpb.RegisterBackgroundImageServiceHandler,
	)
}
