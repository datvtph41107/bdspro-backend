package service

import (
	"common"
	"context"
	"fmt"
	hubpb "pb/types/hub"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterHubService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.hub")
	url := common.GetPrefixProtobufUrl(env, "hub")
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = hubpb.RegisterEventQueueServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterHubInternalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterLocationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterLocationV2ServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterUserGuideServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterSystemConfigServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterFAQServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterVersionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterApiKeyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterInteractiveEventServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterErrorLogServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterUpdateDataServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterApplinkServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = hubpb.RegisterRequestStatsServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
