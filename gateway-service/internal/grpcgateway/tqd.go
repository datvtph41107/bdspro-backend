package grpcgateway

import (
	"context"
	tqdpb "pb/types/tqd"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterTqdService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		tqdpb.RegisterAmenityServiceHandler,
		tqdpb.RegisterContactLabelServiceHandler,
		tqdpb.RegisterDirectoryCategoryServiceHandler,
		tqdpb.RegisterDirectorySourceServiceHandler,
		tqdpb.RegisterDirectorySupplierServiceHandler,
		tqdpb.RegisterOpenHourServiceHandler,
		tqdpb.RegisterPoiCategoryServiceHandler,
		tqdpb.RegisterPoiServiceHandler,
		tqdpb.RegisterLocationServiceHandler,
		tqdpb.RegisterFeatureServiceHandler,
		tqdpb.RegisterParcelServiceHandler,
		tqdpb.RegisterOneHouseServiceHandler,
		tqdpb.RegisterLayerServiceHandler,
		tqdpb.RegisterGisReportAdminServiceHandler,
		tqdpb.RegisterAdminUsageServiceHandler,
		tqdpb.RegisterProfileUsageServiceHandler,
		tqdpb.RegisterReportServiceHandler,
		tqdpb.RegisterSubscriptionServiceHandler,
		tqdpb.RegisterQHLabelServiceHandler,
		tqdpb.RegisterQHAuthorityIssuringServiceHandler,
		tqdpb.RegisterQHLayerFamilyServiceHandler,
		tqdpb.RegisterQHLayerLegendServiceHandler,
		tqdpb.RegisterImportServiceHandler,
		tqdpb.RegisterRegionServiceHandler,
		tqdpb.RegisterTqdPublicServiceHandler,
		tqdpb.RegisterLayerLegalServiceHandler,
		tqdpb.RegisterQHLandUseServiceHandler,
		tqdpb.RegisterRegionExtendServiceHandler,
		tqdpb.RegisterMapWorkspaceServiceHandler,
		tqdpb.RegisterMapPointServiceHandler,
		tqdpb.RegisterQHPlanningServiceHandler,
		tqdpb.RegisterDiscoveryServiceHandler,
		tqdpb.RegisterRelatedEntityServiceHandler,
		tqdpb.RegisterPlanningClientServiceHandler,
		tqdpb.RegisterTqdSeoProjectionServiceHandler,
	)
}
