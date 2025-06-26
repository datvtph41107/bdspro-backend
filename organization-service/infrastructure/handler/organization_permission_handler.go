package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type OrganizationPermissionHandler struct {
	organizationpb.UnimplementedOrganizationPermissionServiceServer
	OrganizationPermissionUsecase     usecase.OrganizationPermissionUsecase
	OrganizationPermissionTransformer transformer.OrganizationPermissionTransformer
	OrganizationPermissionValidator   validator.OrganizationPermissionValidator
}

func NewOrganizationPermissionHandler(
	organizationPermissionUsecase usecase.OrganizationPermissionUsecase,
	organizationPermissionTransformer transformer.OrganizationPermissionTransformer,
	organizationPermissionValidator validator.OrganizationPermissionValidator,
) *OrganizationPermissionHandler {
	return &OrganizationPermissionHandler{
		OrganizationPermissionUsecase:     organizationPermissionUsecase,
		OrganizationPermissionTransformer: organizationPermissionTransformer,
		OrganizationPermissionValidator:   organizationPermissionValidator,
	}
}

// @Summary Tạo quyền
// @Description Tạo quyền
// @Tags Quyền
// @Accept json
// @Produce json
// @Param permission body organizationpb.CreatePermissionRequest true "Thông tin quyền"
// @Security BearerAuth
// @Router /organization/permission [post]
func (h *OrganizationPermissionHandler) CreatePermission(ctx context.Context, req *organizationpb.CreatePermissionRequest) (*organizationpb.CreatePermissionResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateCreatePermissionRequest(req); err != nil {
		return nil, err
	}
	permission := h.OrganizationPermissionTransformer.CreatePermissionRequestToEntity(req)
	permission, err := h.OrganizationPermissionUsecase.CreateNewPermission(ctx, permission)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToCreatePermissionResponse(permission), nil
}

// @Summary Cập nhật quyền
// @Description Cập nhật quyền
// @Tags Quyền
// @Accept json
// @Produce json
// @Param permission body organizationpb.UpdatePermissionRequest true "Thông tin quyền"
// @Security BearerAuth
// @Router /organization/permission/{id} [put]
func (h *OrganizationPermissionHandler) UpdatePermission(ctx context.Context, req *organizationpb.UpdatePermissionRequest) (*organizationpb.UpdatePermissionResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateUpdatePermissionRequest(req); err != nil {
		return nil, err
	}
	permission := h.OrganizationPermissionTransformer.UpdatePermissionRequestToEntity(req)
	permission, err := h.OrganizationPermissionUsecase.UpdatePermission(ctx, permission)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToUpdatePermissionResponse(permission), nil
}

// @Summary Xóa quyền
// @Description Xóa quyền
// @Tags Quyền
// @Accept json
// @Produce json
// @Param id path int true "ID của quyền"
// @Security BearerAuth
// @Router /organization/permission/{id} [delete]
func (h *OrganizationPermissionHandler) DeletePermission(ctx context.Context, req *organizationpb.DeletePermissionRequest) (*organizationpb.DeletePermissionResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateDeletePermissionRequest(req); err != nil {
		return nil, err
	}
	err := h.OrganizationPermissionUsecase.DeletePermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeletePermissionResponse{
		Id: req.Id,
	}, nil
}

// @Summary Lấy danh sách quyền
// @Description Lấy danh sách quyền
// @Tags Quyền
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng mục trên mỗi trang"
// @Security BearerAuth
// @Router /organization/permission [get]
func (h *OrganizationPermissionHandler) GetPermissions(ctx context.Context, req *organizationpb.GetPermissionsRequest) (*organizationpb.GetPermissionsResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateGetPermissionsRequest(req); err != nil {
		return nil, err
	}

	permissions, err := h.OrganizationPermissionUsecase.GetPermissionOrganization(ctx)
	if err != nil {
		return nil, err
	}
	results, err := h.OrganizationPermissionTransformer.EntityToGetPermissionsResponse(permissions, 0), nil
	return results, err
}
