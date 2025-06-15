package service

import (
	"context"
	"fmt"
	pb_bdspro "pb/types/bdspro"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterBdsproService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.bdspro")
	url := "localhost"
	if env == "release" {
		url = "bdspro_grpc"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	// _ = pb_bdspro.RegisterBdsproServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterBdsproPublicServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterProductServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterPostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterCostDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterIncomeDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetIncomeTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetExploitationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetCostTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetCostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetLegalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_bdspro.RegisterAssetSplitHistoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
