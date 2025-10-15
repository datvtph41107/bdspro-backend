package service

import (
	"common"
	"context"
	"fmt"
	marketingpb "pb/types/marketing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterMarketingService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.marketing")
	url := common.GetPrefixProtobufUrl(env, "marketing")

	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = marketingpb.RegisterMarketingServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
