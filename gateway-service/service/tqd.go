package service

import (
	"common"
	"context"
	"fmt"
	"log"
	"time"

	tqdpb "pb/types/tqd"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const tqdDialTimeout = 15 * time.Second

// RegisterTqdService registers all TQD service handlers to the gateway mux.
// Khi start gateway sẽ chủ động dial tới TQD (có timeout); nếu thành công dùng 1 conn chung, không thì fallback lazy connect.
func RegisterTqdService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.tqd")
	url := common.GetPrefixProtobufUrl(env, "tqd")
	endpoint := fmt.Sprintf("%s:%s", url, port)

	optsWithBlock := append([]grpc.DialOption{}, opts...)
	optsWithBlock = append(optsWithBlock, grpc.WithBlock())
	dialCtx, cancel := context.WithTimeout(ctx, tqdDialTimeout)
	conn, err := grpc.DialContext(dialCtx, endpoint, optsWithBlock...)
	cancel()

	if err != nil {
		log.Printf("[TQD] Không kết nối được tới %s khi start (sẽ connect khi có request): %v", endpoint, err)
		registerTqdFromEndpoint(ctx, mux, endpoint, opts)
		return
	}
	log.Printf("[TQD] Đã kết nối tới %s", endpoint)
	registerTqdWithConn(ctx, mux, conn)
}

func registerTqdWithConn(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
	_ = tqdpb.RegisterAmenityServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterContactLabelServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterDirectoryCategoryServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterDirectorySourceServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterDirectorySupplierServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterOpenHourServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterPoiCategoryServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterPoiServiceHandler(ctx, mux, conn)
}

func registerTqdFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	_ = tqdpb.RegisterAmenityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterContactLabelServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectoryCategoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectorySourceServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDirectorySupplierServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterOpenHourServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterPoiCategoryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterPoiServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterLocationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
