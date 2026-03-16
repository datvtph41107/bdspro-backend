package service

import (
	"common"
	"context"
	"fmt"
	tqdpb "pb/types/tqd"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// RegisterTqdService registers all TQD service handlers to the gateway mux
func RegisterTqdService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.tqd")
	url := common.GetPrefixProtobufUrl(env, "tqd")

	endpoint := fmt.Sprintf("%s:%s", url, port)
	
	// Register all TQD service handlers
	_ = tqdpb.RegisterAmenityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterContactLabelServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectoryCategoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectorySourceServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectorySupplierServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterOpenHourServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterPoiCategoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterPoiServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
