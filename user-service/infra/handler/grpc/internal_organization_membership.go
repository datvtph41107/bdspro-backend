package grpc

import (
	"context"

	"common/identity"
	userpb "pb/types/user"
	organization "user/internal/usecase/organization"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalOrganizationMembershipHandler struct {
	userpb.UnimplementedInternalOrganizationMembershipServiceServer
	service *organization.Service
}

func NewInternalOrganizationMembershipHandler(service *organization.Service) *InternalOrganizationMembershipHandler {
	return &InternalOrganizationMembershipHandler{service: service}
}

func (h *InternalOrganizationMembershipHandler) CheckMembers(ctx context.Context, request *userpb.CheckOrganizationMembersRequest) (*userpb.CheckOrganizationMembersResponse, error) {
	caller, ok := identity.ServiceCallerFromContext(ctx)
	if !ok || caller.ServiceID != "crm-service" {
		return nil, status.Error(codes.PermissionDenied, "organization membership caller is not allowed")
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	checks, err := h.service.CheckMembers(ctx, request.OrganizationId, request.ProfileIds)
	if err != nil {
		return nil, organizationError(err)
	}
	response := &userpb.CheckOrganizationMembersResponse{Data: make([]*userpb.OrganizationMemberCheck, 0, len(checks))}
	for _, check := range checks {
		joinedAt := ""
		if !check.JoinedAt.IsZero() {
			joinedAt = check.JoinedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
		}
		response.Data = append(response.Data, &userpb.OrganizationMemberCheck{
			ProfileId: check.ProfileID, IsMember: check.IsMember, JoinedAt: joinedAt,
			MemberId: check.MemberID, RoleKey: check.RoleKey, Status: check.Status,
		})
	}
	return response, nil
}
