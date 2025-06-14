package service

import (
	"context"
	"fmt"
	crmpb "pb/types/crm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterCrmService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.crm")
	url := "localhost"
	if env == "release" {
		url = "crm"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = crmpb.RegisterPipelineServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterContactServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterLeadServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterStageServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterHistoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterRuleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterBlockServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFollowServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFriendGroupServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFriendServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterInvitationInstallServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
