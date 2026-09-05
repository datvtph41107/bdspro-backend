package grpcgateway

import (
	"context"
	assistantpb "pb/types/assistant"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterAssistantService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		assistantpb.RegisterAssistantServiceHandler,
	)
}
