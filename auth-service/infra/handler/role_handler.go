package handler

import (
	_dto "common/domain/dto"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"

	"auth/infra/mapper"
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/usecase"
)

type RoleHandler struct {
	authpb.UnimplementedRoleServiceServer
	roleUsecase     *usecase.RoleUsecase
	roleMapper      *mapper.RoleMapper
	internalHandler *InternalHandler
}

func NewRoleHandler(
	roleUsecase *usecase.RoleUsecase,
	roleMapper *mapper.RoleMapper,
	internalHandler *InternalHandler,
) *RoleHandler {
	return &RoleHandler{
		roleUsecase:     roleUsecase,
		roleMapper:      roleMapper,
		internalHandler: internalHandler,
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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_TAO"}); err != nil {
		return nil, err
	}

	role := &domain.Role{
		RoleName:        req.RoleName,
		RoleDescription: req.RoleDescription,
		Key:             req.Key,
		PermissionIDs:   req.PermissionIds,
		IsDefault:       req.IsDefault,
		DomainType:      req.DomainType,
		ColorID:         req.ColorId,
		RoleGroupID:     req.RoleGroupId,
		RoleKey:         req.RoleKey,
		AllowAssign:     req.AllowAssign,
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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_SUA"}); err != nil {
		return nil, err
	}

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
		RoleKey:         req.RoleKey,
		AllowAssign:     req.AllowAssign,
	}

	if err := h.internalHandler.ResetPermission(ctx, req.Id); err != nil {
		return nil, err
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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_XOA"}); err != nil {
		return nil, err
	}

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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_XEM"}); err != nil {
		return nil, err
	}

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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_XEM"}); err != nil {
		return nil, err
	}

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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_SUA"}); err != nil {
		return nil, err
	}

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
	if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_XEM"}); err != nil {
		return nil, err
	}

	roleSearchRequest := &dto.RoleSearchRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Scope:       req.Scope,
		DomainType:  req.DomainType,
		Name:        req.Name,
		IsDefault:   req.IsDefault,
		RoleGroupID: req.RoleGroupId,
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

// @Summary Lấy danh sách role theo group key
// @Description Lấy danh sách role theo group key
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param groupKey path int true "Group key"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} authpb.RoleListResponse
// @Router /role/by-group/{groupKey} [get]
func (h *RoleHandler) GetRolesByGroupKey(ctx context.Context, req *authpb.GetRolesByGroupKeyRequest) (*authpb.RoleListResponse, error) {
	// if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_XEM"}); err != nil {
	// 	return nil, err
	// }

	roles, total, err := h.roleUsecase.GetListByGroupKey(ctx, req.GroupKey)
	if err != nil {
		return nil, err
	}

	pbRoles := h.roleMapper.MapRoleListToItemPb(roles)

	return &authpb.RoleListResponse{
		Data:          pbRoles,
		TotalElements: uint32(total),
	}, nil
}

// @Summary Lấy danh sách role theo module code
// @Description Lấy danh sách role theo module code
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "Module code"
// @Success 200 {object} authpb.RoleListResponse
// @Router /role/by-module/{code} [get]
func (h *RoleHandler) GetRolesByModuleCode(ctx context.Context, req *authpb.GetRolesByGroupKeyRequest) (*authpb.RoleListResponse, error) {
	roles, err := h.roleUsecase.GetListByModuleCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	pbRoles := h.roleMapper.MapRoleListToItemPb(roles)

	return &authpb.RoleListResponse{
		Data:          pbRoles,
		TotalElements: uint32(len(roles)),
	}, nil
}

// @Summary Gán role cho user
// @Description Gán role cho user dựa trên roleKeys
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body authpb.AssignRoleToUserRequest true "Assign role request"
// @Success 200 {object} authpb.AssignRoleToUserResponse
// @Router /role/assign-user [post]
func (h *RoleHandler) AssignRoleToUser(ctx context.Context, req *authpb.AssignRoleToUserRequest) (*authpb.AssignRoleToUserResponse, error) {
	// if err := h.internalHandler.HasPermissions(ctx, []string{"ADMIN_VT_GAN"}); err != nil {
	// 	return nil, err
	// }

	// Convert request to DTO
	assignRequest := &dto.AssignRoleToUserRequest{
		UserID: req.UserId,
		RoleID: req.RoleId,
	}

	// Call usecase
	result, err := h.roleUsecase.AssignRoleToUser(ctx, assignRequest)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to protobuf
	pbRole := h.convertRoleDTOToPb(result.AssignedRole)

	return &authpb.AssignRoleToUserResponse{
		AssignedRole: pbRole,
		Message:      result.Message,
	}, nil
}

// convertRoleDTOToPb convert dto.RoleDTO to authpb.Role
func (h *RoleHandler) convertRoleDTOToPb(role dto.RoleDTO) *authpb.Role {
	pbRole := &authpb.Role{
		Id:              role.ID,
		RoleName:        role.RoleName,
		RoleDescription: role.RoleDescription,
		Key:             role.Key,
		PermissionIds:   role.PermissionIDs,
		OrganizationId:  role.OrganizationID,
		IsDefault:       role.IsDefault,
		DomainType:      role.DomainType,
		AllowAssign:     role.AllowAssign,
		RoleKey:         role.RoleKey,
		CreatedAt:       role.CreatedAt,
		UpdatedAt:       role.UpdatedAt,
	}

	// Convert permissions
	if role.Permissions != nil {
		permissions := make([]*authpb.Permission, len(role.Permissions))
		for i, perm := range role.Permissions {
			permissions[i] = &authpb.Permission{
				Id:             perm.ID,
				Name:           perm.Name,
				Key:            perm.Key,
				Description:    perm.Description,
				PermissionType: perm.PermissionType,
				Module:         perm.Module,
				ParentId:       perm.ParentID,
				Priority:       perm.Priority,
				CreatedAt:      perm.CreatedAt,
				UpdatedAt:      perm.UpdatedAt,
			}
		}
		pbRole.Permissions = permissions
	}

	// Convert color
	if role.Color != nil {
		pbRole.Color = &authpb.Color{
			Id:              role.Color.ID,
			Name:            role.Color.Name,
			ContentColor:    role.Color.ContentColor,
			BackgroundColor: role.Color.BackgroundColor,
			ColorKey:        role.Color.ColorKey,
			HexCode:         role.Color.HexCode,
			Description:     role.Color.Description,
			IsActive:        role.Color.IsActive,
			CreatedAt:       role.Color.CreatedAt,
			UpdatedAt:       role.Color.UpdatedAt,
			CreatedBy:       role.Color.CreatedBy,
			UpdatedBy:       role.Color.UpdatedBy,
		}
	}

	// Convert role group
	if role.RoleGroupID != nil {
		pbRole.RoleGroupId = role.RoleGroupID
	}
	if role.RoleGroupName != nil {
		pbRole.RoleGroupName = role.RoleGroupName
	}

	return pbRole
}
