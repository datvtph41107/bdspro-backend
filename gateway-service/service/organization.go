package service

import (
	"context"
	"fmt"
	pb_group "pb/types/group"
	pb_organization "pb/types/organization"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterOrganizationService(ctx context.Context, mux *runtime.ServeMux, url string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.organization")
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_organization.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_organization.RegisterOrganizationMemberServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_organization.RegisterOrganizationRoleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_organization.RegisterOrganizationPermissionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_organization.RegisterOrganizationBranchServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_organization.RegisterOrganizationLogActivityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupDealServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupLogActivityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupNotificationServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupSettingServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupMemberServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_group.RegisterGroupDocumentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
