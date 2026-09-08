package grpcgateway

import (
	"context"
	crmpb "pb/types/crm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterCrmService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		crmpb.RegisterPipelineServiceHandler,
		crmpb.RegisterContactServiceHandler,
		crmpb.RegisterLeadServiceHandler,
		crmpb.RegisterStageServiceHandler,
		crmpb.RegisterRuleServiceHandler,
		crmpb.RegisterBlockServiceHandler,
		crmpb.RegisterFollowServiceHandler,
		crmpb.RegisterFriendGroupServiceHandler,
		crmpb.RegisterFriendServiceHandler,
		crmpb.RegisterInvitationInstallServiceHandler,
		crmpb.RegisterSharingAccessServiceHandler,
		crmpb.RegisterRateServiceHandler,
		crmpb.RegisterReportServiceHandler,
		crmpb.RegisterReportAdminServiceHandler,
		crmpb.RegisterReportReasonServiceHandler,
		crmpb.RegisterSupportTicketServiceHandler,
		crmpb.RegisterAdminSalesServiceHandler,
		crmpb.RegisterRegistryServiceHandler,
		crmpb.RegisterEnumServiceHandler,
		crmpb.RegisterAppointmentServiceHandler,
		crmpb.RegisterSeoDomainServiceHandler,
		crmpb.RegisterCRMPublicContentServiceHandler,
	)
}
