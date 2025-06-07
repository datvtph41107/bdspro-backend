package service

import (
	"context"
	"fmt"
	pb_noti "pb/types/notification"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterNotificationService(ctx context.Context, mux *runtime.ServeMux, url string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.notification")
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_noti.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
