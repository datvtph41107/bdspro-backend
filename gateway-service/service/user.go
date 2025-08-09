package service

import (
	"context"
	"fmt"

	userpb "pb/types/user"

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
	_ = userpb.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterRateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportReasonServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportProofServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
