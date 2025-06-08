package service

import (
	"context"
	"fmt"

	pb_social "pb/types/social"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterSocialService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.social")
	url := "localhost"
	if env == "release" {
		url = "social"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_social.RegisterNewsFeedServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_social.RegisterCommentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = pb_social.RegisterReportReasonServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_social.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = pb_social.RegisterLikeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = pb_social.RegisterShareServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = pb_social.RegisterViewServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
