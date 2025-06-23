package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
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
	permission := h.OrganizationPermissionTransformer.CreatePermissionRequestToEntity(req)
	if err := h.OrganizationPermissionValidator.ValidateCreatePermissionRequest(req); err != nil {
		return nil, err
	}
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
	permission := h.OrganizationPermissionTransformer.UpdatePermissionRequestToEntity(req)
	if err := h.OrganizationPermissionValidator.ValidateUpdatePermissionRequest(req); err != nil {
		return nil, err
	}
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
// @Param permission body organizationpb.DeletePermissionRequest true "Thông tin quyền"
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
// @Param organizationId query int true "ID của tổ chức"
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng mục trên mỗi trang"
// @Security BearerAuth
// @Router /organization/permission [get]
func (h *OrganizationPermissionHandler) GetPermissions(ctx context.Context, req *organizationpb.GetPermissionsRequest) (*organizationpb.GetPermissionsResponse, error) {
	if err := h.OrganizationPermissionValidator.ValidateGetPermissionsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	permissions, total, err := h.OrganizationPermissionUsecase.GetPermissionsByOrganizationId(ctx, req.OrganizationId, page, size)
	if err != nil {
		return nil, err
	}
	return h.OrganizationPermissionTransformer.EntityToGetPermissionsResponse(permissions, total), nil
}
