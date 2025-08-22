package handler

import (
	"auth/infra/mapper"
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/usecase"
	_dto "common/dto"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
)

type RoleHandler struct {
	authpb.UnimplementedRoleServiceServer
	roleUsecase *usecase.RoleUsecase
	roleMapper  *mapper.RoleMapper
}

func NewRoleHandler(roleUsecase *usecase.RoleUsecase, roleMapper *mapper.RoleMapper) *RoleHandler {
	return &RoleHandler{
		roleUsecase: roleUsecase,
		roleMapper:  roleMapper,
	}
}

// @Summary Tạo mới role
// @Description Tạo mới role
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body authpb.RoleRequest true "Role to create"
// @Success 200 {object} authpb.Role
// @Router /role [post]
func (h *RoleHandler) CreateRole(ctx context.Context, req *authpb.RoleRequest) (*authpb.Role, error) {
	// if req.Name == "" || req.PermissionIds == nil {
	// 	return nil, _errors.ReturnError(400, "name or permissionIds is required")
	// }

	role := &domain.Role{
		RoleName:        req.RoleName,
		RoleDescription: req.RoleDescription,
		Key:             req.Key,
		PermissionIDs:   req.PermissionIds,
		IsDefault:       req.IsDefault,
		DomainType:      req.DomainType,
		ColorID:         req.ColorId,
		RoleGroupID:     req.RoleGroupId,
		// OrganizationID: 0,
	}

	result, err := h.roleUsecase.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return h.roleMapper.MapRoleToPb(result), nil
}

// @Summary Cập nhật role
// @Description Cập nhật role
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body authpb.RoleRequest true "Role to update"
// @Success 200 {object} authpb.Role
// @Router /role [put]
func (h *RoleHandler) UpdateRole(ctx context.Context, req *authpb.RoleRequest) (*authpb.Role, error) {
	role := &domain.Role{
		RoleName:        req.RoleName,
		RoleDescription: req.RoleDescription,
		Key:             req.Key,
		OrganizationID:  req.OrganizationId,
		IsDefault:       req.IsDefault,
		DomainType:      req.DomainType,
		ColorID:         req.ColorId,
		RoleGroupID:     req.RoleGroupId,
		PermissionIDs:   req.PermissionIds,
	}

	result, err := h.roleUsecase.Update(ctx, req.Id, role)
	if err != nil {
		return nil, err
	}

	return h.roleMapper.MapRoleToPb(result), nil
}

// @Summary Xóa role
// @Description Xóa role
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /role/{id} [delete]
func (h *RoleHandler) DeleteRole(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := h.roleUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Xóa role thành công",
	}, nil
}

// @Summary Lấy role theo ID
// @Description Lấy role theo ID
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} authpb.Role
// @Router /role/{id} [get]
func (h *RoleHandler) GetRoleById(ctx context.Context, req *sharepb.IdRequest) (*authpb.Role, error) {
	role, err := h.roleUsecase.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.roleMapper.MapRoleToPb(role), nil
}

// @Summary Lấy danh sách role theo organization ID
// @Description Lấy danh sách role theo organization ID
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param page query int true "Page"
// @Param size query int true "Size"
// @Success 200 {object} authpb.RoleListResponse
// @Router /role/organization/{id} [get]
func (h *RoleHandler) GetRolesByOrganization(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleListResponse, error) {
	roles, total, err := h.roleUsecase.GetByOrganizationIDWithPagination(ctx, req.Id, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	pbRoles := h.roleMapper.MapRoleListToPb(roles)

	return &authpb.RoleListResponse{
		Data:          pbRoles,
		TotalElements: uint32(total),
	}, nil
}

// @Summary Thêm permission vào role
// @Description Thêm permission vào role
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body authpb.AddPermissionToRoleRequest true "Role to add permission"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /role/permissions [post]
func (h *RoleHandler) AddPermissionToRole(ctx context.Context, req *authpb.AddPermissionToRoleRequest) (*sharepb.SubmitResponse, error) {
	err := h.roleUsecase.AddPermissionsToRole(ctx, req.RoleId, req.PermissionIds)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.RoleId,
		Message: "Thêm permission vào role thành công",
	}, nil
}

// @Summary Lấy danh sách role
// @Description Lấy danh sách role
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body authpb.RoleSearchRequest true "Role search request"
// @Success 200 {object} authpb.RoleListResponse
// @Router /role/list [get]
func (h *RoleHandler) GetListRole(ctx context.Context, req *authpb.RoleSearchRequest) (*authpb.RoleListResponse, error) {
	roleSearchRequest := &dto.RoleSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Scope:      req.Scope,
		DomainType: req.DomainType,
		Name:       req.Name,
		IsDefault:  req.IsDefault,
	}
	roles, total, err := h.roleUsecase.GetList(ctx, roleSearchRequest)
	if err != nil {
		return nil, err
	}

	pbRoles := h.roleMapper.MapRoleListToPb(roles)

	return &authpb.RoleListResponse{
		Data:          pbRoles,
		TotalElements: uint32(total),
	}, nil
}
