package grpcgateway

import (
	"context"
	userpb "pb/types/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		userpb.RegisterProfileServiceHandler,
		userpb.RegisterReportServiceHandler,
		userpb.RegisterReportReasonServiceHandler,
		userpb.RegisterReportProofServiceHandler,
		userpb.RegisterBookmarkUserServiceHandler,
		userpb.RegisterAdminProfileServiceHandler,
		userpb.RegisterAdminUserProfileServiceHandler,
		userpb.RegisterAdminCatalogServiceHandler,
		userpb.RegisterAdminCommercialServiceHandler,
		userpb.RegisterAdminOrganizationServiceHandler,
		userpb.RegisterCommercialProfileServiceHandler,
		userpb.RegisterMainAreaServiceHandler,
		userpb.RegisterPurposeUseServiceHandler,
		userpb.RegisterProfessionServiceHandler,
		userpb.RegisterTagServiceHandler,
		userpb.RegisterKYCServiceHandler,
		userpb.RegisterPriceTableServiceHandler,
		userpb.RegisterCheckoutServiceHandler,
	)
}
