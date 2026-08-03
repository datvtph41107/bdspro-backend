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

// TqdGRPCMaxMsgBytes — vượt default gRPC 4MB cho ImportGeoJSON (file_content) và response lớn.
// Giữ cùng mức với multipart ParseMultipartForm ở tqd_import_multipart.go (512 MiB).
const TqdGRPCMaxMsgBytes = 512 << 20

func mergeTqdDialOpts(opts []grpc.DialOption) []grpc.DialOption {
	return append([]grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(TqdGRPCMaxMsgBytes),
			grpc.MaxCallSendMsgSize(TqdGRPCMaxMsgBytes),
		),
	}, opts...)
}

// DialTQDGRPC tạo client tới tqd-service (dùng chung cho mux và route multipart).
func DialTQDGRPC(env string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	port := viper.GetString("service.port.tqd")
	url := common.GetPrefixProtobufUrl(env, "tqd")
	endpoint := fmt.Sprintf("%s:%s", url, port)
	return grpc.NewClient(endpoint, mergeTqdDialOpts(opts)...)
}

// RegisterTqdService registers all TQD service handlers to the gateway mux.
// Dùng grpc.NewClient (không block khi start): kết nối thật sự và reconnect diễn ra
// khi có RPC; backend tạm tắt thì request sau vẫn thử lại theo ConnectParams ở main.
// Nếu NewClient trả lỗi (cấu hình/target), fallback FromEndpoint; gateway vẫn khởi động bình thường.
func RegisterTqdService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) *grpc.ClientConn {
	port := viper.GetString("service.port.tqd")
	url := common.GetPrefixProtobufUrl(env, "tqd")
	endpoint := fmt.Sprintf("%s:%s", url, port)

	conn, err := DialTQDGRPC(env, opts)
	if err != nil {
		log.Printf("[TQD] Không tạo được client tới %s (đăng ký lazy theo endpoint): %v", endpoint, err)
		registerTqdFromEndpoint(ctx, mux, endpoint, mergeTqdDialOpts(opts))
		return nil
	}
	log.Printf("[TQD] Đã tạo client tới %s (kết nối khi có request)", endpoint)
	registerTqdWithConn(ctx, mux, conn)
	return conn
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
	_ = tqdpb.RegisterFeatureServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterParcelServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterOneHouseServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterLayerServiceHandler(ctx, mux, conn)
	// GisReportAdmin trước ReportService — admin HTTP paths chỉ còn trên GisReportAdmin.
	_ = tqdpb.RegisterGisReportAdminServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterReportServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterGisReportAdminServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterSubscriptionServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHLabelServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHAuthorityIssuringServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHLayerFamilyServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHLayerLegendServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterImportServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterRegionServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterTqdPublicServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterLayerLegalServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHLandUseServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterRegionExtendServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterMapWorkspaceServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterQHPlanningServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterDiscoveryServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterPlanningClientServiceHandler(ctx, mux, conn)
	_ = tqdpb.RegisterTqdSeoProjectionServiceHandler(ctx, mux, conn)
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
	_ = tqdpb.RegisterFeatureServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterParcelServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterOneHouseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterLayerServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHLayerFamilyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHLayerLegendServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterImportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterGisReportAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterRegionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterGisReportAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterSubscriptionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHLabelServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHAuthorityIssuringServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterTqdPublicServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterLayerLegalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHLandUseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterRegionExtendServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterMapWorkspaceServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterQHPlanningServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterDiscoveryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterPlanningClientServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = tqdpb.RegisterTqdSeoProjectionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
