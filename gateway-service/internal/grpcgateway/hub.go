package grpcgateway

import (
	"context"
	hubpb "pb/types/hub"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterHubService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		hubpb.RegisterEventQueueServiceHandler,
		hubpb.RegisterHubInternalServiceHandler,
		hubpb.RegisterLocationServiceHandler,
		hubpb.RegisterLocationV2ServiceHandler,
		hubpb.RegisterUserGuideServiceHandler,
		hubpb.RegisterSystemConfigServiceHandler,
		hubpb.RegisterFAQServiceHandler,
		hubpb.RegisterVersionServiceHandler,
		hubpb.RegisterApiKeyServiceHandler,
		hubpb.RegisterInteractiveEventServiceHandler,
		hubpb.RegisterErrorLogServiceHandler,
		hubpb.RegisterUpdateDataServiceHandler,
		hubpb.RegisterApplinkServiceHandler,
	)
}
