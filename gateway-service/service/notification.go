package service

import (
	"common"
	"context"
	"fmt"
	notificationpb "pb/types/notification"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.notification")
	url := common.GetPrefixProtobufUrl(env, "notification")

	endpoint := fmt.Sprintf("%s:%s", url, port)
	fmt.Println("url", url, "endpoint", endpoint)
	_ = notificationpb.RegisterGatewayNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = notificationpb.RegisterHistoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = notificationpb.RegisterPropertyHistoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
