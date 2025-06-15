package service

import (
	"context"
	"fmt"

	pb_user "pb/types/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.user")
	url := "localhost"
	if env == "release" {
		url = "user_grpc"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_user.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
