package handler

import (
	"context"
	"fmt"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"github.com/tuandq2112/organization-service/pkg/utils"
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
	fmt.Println("go1", h.OrganizationRoleUsecase)

	role := h.OrganizationRoleTransformer.CreateRoleRequestToEntity(req)
	fmt.Println("go2")
	if err := h.OrganizationRoleValidator.ValidateCreateRoleRequest(req); err != nil {
		return nil, err
	}
	fmt.Println("go3")
	role, err := h.OrganizationRoleUsecase.CreateNewRole(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), role)
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
	role, err := h.OrganizationRoleUsecase.UpdateRole(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), role)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToUpdateRoleResponse(role), nil
}

func (h *OrganizationRoleHandler) DeleteRole(ctx context.Context, req *organizationpb.DeleteRoleRequest) (*organizationpb.DeleteRoleResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateDeleteRoleRequest(req); err != nil {
		return nil, err
	}
	err := h.OrganizationRoleUsecase.DeleteRole(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeleteRoleResponse{
		Id: req.Id,
	}, nil
}

func (h *OrganizationRoleHandler) GetRoles(ctx context.Context, req *organizationpb.GetRolesRequest) (*organizationpb.GetRolesResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateGetRolesRequest(req); err != nil {
		return nil, err
	}
	roles, err := h.OrganizationRoleUsecase.GetRoleByOrganizationId(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToGetRolesResponse(roles), nil
}
