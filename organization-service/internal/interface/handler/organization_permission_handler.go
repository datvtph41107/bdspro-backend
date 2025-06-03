package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type OrganizationPermissionHandler struct {
	organizationpb.UnimplementedOrganizationPermissionServiceServer
	OrganizationPermissionUsecase usecase.OrganizationPermissionUsecase
	OrganizationPermissionTransformer transformer.OrganizationPermissionTransformer
	OrganizationPermissionValidator validator.OrganizationPermissionValidator
}

func NewOrganizationPermissionHandler(
	organizationPermissionUsecase usecase.OrganizationPermissionUsecase, 
	organizationPermissionTransformer transformer.OrganizationPermissionTransformer,
	organizationPermissionValidator validator.OrganizationPermissionValidator,
) *OrganizationPermissionHandler {

	return &OrganizationPermissionHandler{
		OrganizationPermissionUsecase: organizationPermissionUsecase,
		OrganizationPermissionTransformer: organizationPermissionTransformer,
		OrganizationPermissionValidator: organizationPermissionValidator,
	}
}

func (h *OrganizationPermissionHandler) CreatePermission(ctx context.Context, req *organizationpb.CreatePermissionRequest) (*organizationpb.CreatePermissionResponse, error) {
	permission := h.OrganizationPermissionTransformer.CreatePermissionRequestToEntity(req)
	if err := h.OrganizationPermissionValidator.ValidateCreatePermissionRequest(req); err != nil {
		return nil, err
	}
	permission, err := h.OrganizationPermissionUsecase.CreateNewPermission(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), permission)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToCreatePermissionResponse(permission), nil
}

func (h *OrganizationPermissionHandler) UpdatePermission(ctx context.Context, req *organizationpb.UpdatePermissionRequest) (*organizationpb.UpdatePermissionResponse, error) {
	permission := h.OrganizationPermissionTransformer.UpdatePermissionRequestToEntity(req)
	if err := h.OrganizationPermissionValidator.ValidateUpdatePermissionRequest(req); err != nil {
		return nil, err
	}
	permission, err := h.OrganizationPermissionUsecase.UpdatePermission(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), permission)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToUpdatePermissionResponse(permission), nil
}

func (h *OrganizationPermissionHandler) DeletePermission(ctx context.Context, req *organizationpb.DeletePermissionRequest) (*organizationpb.DeletePermissionResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateDeletePermissionRequest(req); err != nil {
		return nil, err
	}
	err := h.OrganizationPermissionUsecase.DeletePermission(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeletePermissionResponse{
		Id: req.Id,
	}, nil
}

func (h *OrganizationPermissionHandler) GetPermissions(ctx context.Context, req *organizationpb.GetPermissionsRequest) (*organizationpb.GetPermissionsResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateGetPermissionsRequest(req); err != nil {
		return nil, err
	}
	permissions, err := h.OrganizationPermissionUsecase.GetPermissionsByOrganizationId(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToGetPermissionsResponse(permissions), nil
}