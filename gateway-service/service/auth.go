package service

import (
	"context"
	"fmt"
	authpb "pb/types/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.auth")
	url := "localhost"
	if env == "release" {
		url = "auth"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = authpb.RegisterOAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
