package grpc

import (
	"context"
	"errors"
	"time"

	userpb "pb/types/user"
	domain "user/internal/domain/organization"
	organization "user/internal/usecase/organization"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AdminOrganizationHandler struct {
	userpb.UnimplementedAdminOrganizationServiceServer
	service    *organization.Service
	authorizer organization.PermissionAuthorizer
}

func NewAdminOrganizationHandler(service *organization.Service, authorizer organization.PermissionAuthorizer) *AdminOrganizationHandler {
	return &AdminOrganizationHandler{service: service, authorizer: authorizer}
}

func (h *AdminOrganizationHandler) authorize(ctx context.Context, permission string) (uint64, error) {
	actorID, err := organization.ActorIDFromContext(ctx)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, err.Error())
	}
	if h == nil || h.authorizer == nil {
		return 0, status.Error(codes.Internal, "organization permission authority is unavailable")
	}
	allowed, err := h.authorizer.HasPermission(ctx, actorID, permission)
	if err != nil {
		return 0, status.Error(codes.Unavailable, "organization permission authority is unavailable")
	}
	if !allowed {
		return 0, status.Error(codes.PermissionDenied, "organization permission denied")
	}
	return actorID, nil
}

func organizationError(err error) error {
	switch {
	case errors.Is(err, organization.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, organization.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, organization.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, organization.ErrActiveSubscription):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "organization operation failed")
	}
}

func organizationProjection(value domain.Organization) *userpb.AdminOrganizationProjection {
	createdAt := ""
	updatedAt := ""
	if !value.CreatedAt.IsZero() {
		createdAt = value.CreatedAt.UTC().Format(time.RFC3339Nano)
	}
	if !value.UpdatedAt.IsZero() {
		updatedAt = value.UpdatedAt.UTC().Format(time.RFC3339Nano)
	}
	return &userpb.AdminOrganizationProjection{
		Organization: &userpb.AdminOrganization{
			Id: value.ID, Code: value.Code, Name: value.Name, Type: value.Type,
			TaxCode: value.TaxCode, Phone: value.Phone, Email: value.Email,
			Address: value.Address, Website: value.Website, Description: value.Description,
			Status: value.Status, VerificationStatus: value.VerificationStatus,
			WarningLevel: value.WarningLevel, WarningCount: value.WarningCount,
			OwnerProfileId: value.OwnerProfileID, CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
		TotalMembers: value.ActiveMemberCount,
	}
}

func (h *AdminOrganizationHandler) ListOrganizations(ctx context.Context, request *userpb.ListAdminOrganizationsRequest) (*userpb.ListAdminOrganizationsResponse, error) {
	if _, err := h.authorize(ctx, organization.PermissionView); err != nil {
		return nil, err
	}
	if request == nil {
		request = &userpb.ListAdminOrganizationsRequest{}
	}
	values, total, err := h.service.List(ctx, domain.ListQuery{
		Page: request.Page, Size: request.Size, SearchText: request.SearchText,
		Type: request.Type, Status: request.Status,
		VerificationStatus: request.VerificationStatus, WarningLevel: request.WarningLevel,
	})
	if err != nil {
		return nil, organizationError(err)
	}
	response := &userpb.ListAdminOrganizationsResponse{Data: make([]*userpb.AdminOrganizationProjection, 0, len(values)), Total: total, Page: request.Page, Size: request.Size}
	for _, value := range values {
		response.Data = append(response.Data, organizationProjection(value))
	}
	return response, nil
}

func (h *AdminOrganizationHandler) GetOrganization(ctx context.Context, request *userpb.GetAdminOrganizationRequest) (*userpb.AdminOrganizationProjection, error) {
	if _, err := h.authorize(ctx, organization.PermissionView); err != nil {
		return nil, err
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	value, err := h.service.Get(ctx, request.Id)
	if err != nil {
		return nil, organizationError(err)
	}
	return organizationProjection(value), nil
}

func (h *AdminOrganizationHandler) CreateOrganization(ctx context.Context, request *userpb.CreateAdminOrganizationRequest) (*userpb.AdminOrganizationProjection, error) {
	actorID, err := h.authorize(ctx, organization.PermissionManage)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	value, err := h.service.Create(ctx, domain.Organization{
		Code: request.Code, Name: request.Name, Type: request.Type, TaxCode: request.TaxCode,
		Phone: request.Phone, Email: request.Email, Address: request.Address,
		Website: request.Website, Description: request.Description,
		OwnerProfileID: actorID, CreatedBy: actorID, UpdatedBy: actorID,
	})
	if err != nil {
		return nil, organizationError(err)
	}
	return organizationProjection(value), nil
}

func (h *AdminOrganizationHandler) UpdateOrganization(ctx context.Context, request *userpb.UpdateAdminOrganizationRequest) (*userpb.AdminOrganizationProjection, error) {
	actorID, err := h.authorize(ctx, organization.PermissionManage)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	value, err := h.service.Update(ctx, domain.Update{
		ID: request.Id, ActorID: actorID, Code: request.Code, Name: request.Name,
		Type: request.Type, TaxCode: request.TaxCode, Phone: request.Phone,
		Email: request.Email, Address: request.Address, Website: request.Website,
		Description: request.Description, Status: request.Status,
		VerificationStatus: request.VerificationStatus, WarningLevel: request.WarningLevel,
	})
	if err != nil {
		return nil, organizationError(err)
	}
	return organizationProjection(value), nil
}

func (h *AdminOrganizationHandler) ArchiveOrganization(ctx context.Context, request *userpb.ArchiveAdminOrganizationRequest) (*emptypb.Empty, error) {
	actorID, err := h.authorize(ctx, organization.PermissionManage)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := h.service.Archive(ctx, request.Id, actorID); err != nil {
		return nil, organizationError(err)
	}
	return &emptypb.Empty{}, nil
}
