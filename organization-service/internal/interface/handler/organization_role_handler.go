package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type OrganizationRoleHandler struct {
	organizationpb.UnimplementedOrganizationRoleServiceServer
	OrganizationRoleUsecase     usecase.OrganizationRoleUsecase
	OrganizationRoleTransformer transformer.OrganizationRoleTransformer
	OrganizationRoleValidator   validator.OrganizationRoleValidator
}

func NewOrganizationRoleHandler(
	organizationRoleUsecase usecase.OrganizationRoleUsecase,
	organizationRoleTransformer transformer.OrganizationRoleTransformer,
	organizationRoleValidator validator.OrganizationRoleValidator,
) *OrganizationRoleHandler {
	return &OrganizationRoleHandler{
		OrganizationRoleUsecase:     organizationRoleUsecase,
		OrganizationRoleTransformer: organizationRoleTransformer,
		OrganizationRoleValidator:   organizationRoleValidator,
	}
}

func (h *OrganizationRoleHandler) CreateRole(ctx context.Context, req *organizationpb.CreateRoleRequest) (*organizationpb.CreateRoleResponse, error) {
	role := h.OrganizationRoleTransformer.CreateRoleRequestToEntity(req)
	if err := h.OrganizationRoleValidator.ValidateCreateRoleRequest(req); err != nil {
		return nil, err
	}
	role, err := h.OrganizationRoleUsecase.CreateNewRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToCreateRoleResponse(role), nil
}

func (h *OrganizationRoleHandler) UpdateRole(ctx context.Context, req *organizationpb.UpdateRoleRequest) (*organizationpb.UpdateRoleResponse, error) {
	role := h.OrganizationRoleTransformer.UpdateRoleRequestToEntity(req)
	if err := h.OrganizationRoleValidator.ValidateUpdateRoleRequest(req); err != nil {
		return nil, err
	}
	role, err := h.OrganizationRoleUsecase.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToUpdateRoleResponse(role), nil
}

func (h *OrganizationRoleHandler) DeleteRole(ctx context.Context, req *organizationpb.DeleteRoleRequest) (*organizationpb.DeleteRoleResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateDeleteRoleRequest(req); err != nil {
		return nil, err
	}
	err := h.OrganizationRoleUsecase.DeleteRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeleteRoleResponse{
		Id: req.Id,
	}, nil
}

func (h *OrganizationRoleHandler) GetRoles(ctx context.Context, req *organizationpb.GetRolesRequest) (*organizationpb.GetRolesResponse, error) {
	// Set default pagination values
	page := 1
	size := 10

	// Get pagination parameters from request
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Get roles with pagination
	roles, total, err := h.OrganizationRoleUsecase.GetRoles(ctx, req.OrganizationId, page, size)
	if err != nil {
		return nil, err
	}

	// Transform to response
	return h.OrganizationRoleTransformer.EntityToGetRolesResponse(roles, total), nil
}
