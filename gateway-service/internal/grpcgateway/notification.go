package grpcgateway

import (
	"context"
	notificationpb "pb/types/notification"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		notificationpb.RegisterGatewayNotificationServiceHandler,
		notificationpb.RegisterHistoryServiceHandler,
		notificationpb.RegisterPropertyHistoryServiceHandler,
	)
}
