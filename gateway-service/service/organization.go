package service

import (
	"context"
	"fmt"
	organizationpb "pb/types/organization"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterOrganizationService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.organization")
	url := "localhost"
	if env == "release" {
		url = "organization"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = organizationpb.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterOrganizationMemberServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterOrganizationRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterOrganizationPermissionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterOrganizationBranchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterOrganizationLogActivityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupDealServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupLogActivityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupSettingServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupMemberServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterGroupDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterInvestmentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterDealInvitationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterDealTransactionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterDealCommissionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterBusinessDomainServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterColorServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterDealCostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterCostTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterInternalNoteServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = organizationpb.RegisterDealContractServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
