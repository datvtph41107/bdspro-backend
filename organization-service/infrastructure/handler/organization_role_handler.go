package handler

import (
	_dto "common/domain/dto"
	"context"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
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

// @Summary Tạo vai trò
// @Description Tạo vai trò
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateRoleRequest true "Thông tin vai trò"
// @Security BearerAuth
// @Router /organization/role [post]
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

// @Summary Lấy danh sách thành viên của vai trò thương vụ
// @Description Lấy danh sách thành viên của vai trò thương vụ
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param id path int true "ID của vai trò"
// @Security BearerAuth
// @Router /organization/role/deal-members [get]
func (h *OrganizationRoleHandler) GetRoleDealMembers(ctx context.Context, req *organizationpb.GetRoleRequest) (*organizationpb.GetRolesResponse, error) {
	roles, err := h.OrganizationRoleUsecase.GetRoleDealMembers(ctx)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToGetRolesResponse(roles, 0), nil
}

// @Summary Cập nhật vai trò
// @Description Cập nhật vai trò
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param request body organizationpb.UpdateRoleRequest true "Thông tin vai trò"
// @Security BearerAuth
// @Param id path int true "ID của vai trò"
// @Router /organization/role/{id} [put]
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

// @Summary Xóa vai trò
// @Description Xóa vai trò
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param id path int true "ID của vai trò"
// @Security BearerAuth
// @Router /organization/role/{id} [delete]
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

// @Summary Lấy danh sách vai trò
// @Description Lấy danh sách vai trò
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng mục trên mỗi trang"
// @Security BearerAuth
// @Router /organization/role [get]
func (h *OrganizationRoleHandler) GetRoles(ctx context.Context, req *organizationpb.GetRolesRequest) (*organizationpb.GetRolesResponse, error) {
	// Set default pagination values
	page := 0
	size := 10

	// Get pagination parameters from request
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Get roles with pagination
	roles, total, err := h.OrganizationRoleUsecase.GetRoles(ctx, page, size)
	if err != nil {
		return nil, err
	}

	// Transform to response
	return h.OrganizationRoleTransformer.EntityToGetRolesResponse(roles, total), nil
}

// @Summary Lấy chi tiết vai trò
// @Description Lấy chi tiết vai trò theo ID
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param id path int true "ID của vai trò"
// @Security BearerAuth
// @Router /organization/role/{id} [get]
func (h *OrganizationRoleHandler) GetRole(ctx context.Context, req *organizationpb.GetRoleRequest) (*organizationpb.GetRoleResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateGetRoleRequest(req); err != nil {
		return nil, err
	}

	role, err := h.OrganizationRoleUsecase.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.OrganizationRoleTransformer.EntityToGetRoleResponse(role), nil
}

func (h *OrganizationRoleHandler) AddPermissionToRole(ctx context.Context, req *organizationpb.AddPermissionToRoleRequest) (*organizationpb.AddPermissionToRoleResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateAddPermissionToRoleRequest(req); err != nil {
		return nil, err
	}
	err := h.OrganizationRoleUsecase.AddPermissionToRole(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.AddPermissionToRoleResponse{
		Id: req.RoleId,
	}, nil
}

// @Summary Tạo vai trò với thành viên và quyền
// @Description Tạo vai trò mới và gán thành viên cùng quyền
// @Tags Vai trò
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateNewRoleWithMembersRequest true "Thông tin vai trò với thành viên và quyền"
// @Security BearerAuth
// @Router /organization/role/with-members [post]
func (h *OrganizationRoleHandler) CreateNewRoleWithMembers(ctx context.Context, req *organizationpb.CreateNewRoleWithMembersRequest) (*organizationpb.CreateNewRoleWithMembersResponse, error) {
	if err := h.OrganizationRoleValidator.ValidateCreateNewRoleWithMembersRequest(req); err != nil {
		return nil, err
	}

	role := h.OrganizationRoleTransformer.CreateNewRoleWithMembersRequestToEntity(req)
	role, err := h.OrganizationRoleUsecase.CreateNewRoleWithMembers(ctx, role, req.MemberIds, req.PermissionIds)
	if err != nil {
		return nil, err
	}

	return h.OrganizationRoleTransformer.EntityToCreateNewRoleWithMembersResponse(role, req.MemberIds, req.PermissionIds), nil
}

// @Summary Lấy danh sách vai trò của tổ chức hiện tại user đang request
// @Description Lấy danh sách vai trò của tổ chức hiện tại user đang request
// @Tags Vai trò
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /organization/roles/current [get]
func (h *OrganizationRoleHandler) GetRolesCurrent(ctx context.Context, req *sharepb.IdRequest) (*organizationpb.GetRolesResponse, error) {
	dto := _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}

	roles, total, err := h.OrganizationRoleUsecase.GetRoleCurrent(ctx, dto)
	if err != nil {
		return nil, err
	}
	return h.OrganizationRoleTransformer.EntityToGetRolesResponse(roles, total), nil
}
