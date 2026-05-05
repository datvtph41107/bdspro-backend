package service

import (
	"common"
	"context"
	"fmt"
	"log"

	tqdpb "pb/types/tqd"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// RegisterTqdService registers all TQD service handlers to the gateway mux.
// Dùng grpc.NewClient (không block khi start): kết nối thật sự và reconnect diễn ra
// khi có RPC; backend tạm tắt thì request sau vẫn thử lại theo ConnectParams ở main.
// Nếu NewClient trả lỗi (cấu hình/target), fallback FromEndpoint; gateway vẫn khởi động bình thường.
func RegisterTqdService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.tqd")
	url := common.GetPrefixProtobufUrl(env, "tqd")
	endpoint := fmt.Sprintf("%s:%s", url, port)

	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		log.Printf("[TQD] Không tạo được client tới %s (đăng ký lazy theo endpoint): %v", endpoint, err)
		registerTqdFromEndpoint(ctx, mux, endpoint, opts)
		return
	}
	log.Printf("[TQD] Đã tạo client tới %s (kết nối khi có request)", endpoint)
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
	_ = tqdpb.RegisterLocationServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterParcelServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterOneHouseServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterLayerServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterReportServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterSubscriptionServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHLabelServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterRegionServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterTqdPublicServiceHandler(ctx, mux, conn)
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
	_ = tqdpb.RegisterOneHouseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterLayerServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterSubscriptionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterTqdPublicServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
