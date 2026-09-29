package grpcgateway

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterOrganizationService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		organizationpb.RegisterOrganizationServiceHandler,
		organizationpb.RegisterOrganizationMemberServiceHandler,
		organizationpb.RegisterOrganizationRoleServiceHandler,
		organizationpb.RegisterOrganizationPermissionServiceHandler,
		organizationpb.RegisterOrganizationBranchServiceHandler,
		organizationpb.RegisterOrganizationLogActivityServiceHandler,
		organizationpb.RegisterGroupServiceHandler,
		organizationpb.RegisterGroupLogActivityServiceHandler,
		organizationpb.RegisterGroupNotificationServiceHandler,
		organizationpb.RegisterGroupSettingServiceHandler,
		organizationpb.RegisterGroupMemberServiceHandler,
		organizationpb.RegisterGroupDocumentServiceHandler,
		organizationpb.RegisterDealServiceHandler,
		organizationpb.RegisterInvestmentServiceHandler,
		organizationpb.RegisterDealMemberServiceHandler,
		organizationpb.RegisterDealCommissionServiceHandler,
		organizationpb.RegisterBusinessDomainServiceHandler,
		organizationpb.RegisterColorServiceHandler,
		organizationpb.RegisterInternalNoteServiceHandler,
	)
}
