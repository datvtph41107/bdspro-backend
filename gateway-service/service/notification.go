package service

import (
	"context"
	"fmt"
	pb_noti "pb/types/notification"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.notification")
	url := "localhost"
	if env == "release" {
		url = "notification"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	fmt.Println("url", url, "endpoint", endpoint)
	_ = pb_noti.RegisterGatewayNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
