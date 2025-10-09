package service

import (
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"

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
	// _ = bdspropb.RegisterBdsproServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterBdsproPublicServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterPostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterCostDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterIncomeDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetIncomeTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetExploitationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetCostTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetCostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetLegalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetSplitHistoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterSharingAccessServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterTransactionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAssetLegalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterProjectBuildServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminProductServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminAssetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminPostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminPropertyTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminProjectServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminRegionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = bdspropb.RegisterAdminTransactionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
