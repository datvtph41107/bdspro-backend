package handler

import (
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"user/infra/mapper"
	"user/internal/usecase"
)

type PermissionHandler struct {
	authpb.UnimplementedPermissionServiceServer
	permissionUsecase *usecase.PermissionUsecase
	roleMapper        *mapper.RoleMapper
	roleUsecase       *usecase.RoleUsecase
}

func NewPermissionHandler(
	permissionUsecase *usecase.PermissionUsecase,
	roleMapper *mapper.RoleMapper,
	roleUsecase *usecase.RoleUsecase,
) *PermissionHandler {
	return &PermissionHandler{
		permissionUsecase: permissionUsecase,
		roleMapper:        roleMapper,
		roleUsecase:       roleUsecase,
	}
}

// @Summary Get permission by ID
// @Description Get permission by ID
// @Tags Permission
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} sharepb.Permission
// @Router /permission/{id} [get]
func (h *PermissionHandler) GetPermissionById(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Permission, error) {
	permission, err := h.permissionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.roleMapper.MapPermissionToPb(permission), nil
}

// GetPermissionsByModule giữ nguyên endpoint production và chuyển toàn bộ
// truy vấn qua PermissionUsecase của User service.
func (h *PermissionHandler) GetPermissionsByModule(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionListResponse, error) {
	permissions, total, err := h.permissionUsecase.GetByModuleWithPagination(ctx, req.GetId(), int(req.GetPage()), int(req.GetSize()))
	if err != nil {
		return nil, err
	}
	items := make([]*sharepb.Permission, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, h.roleMapper.MapPermissionListWithChild(permission))
	}
	return &authpb.PermissionListResponse{Data: items, TotalElements: total}, nil
}

// @Summary Get all permissions
// @Description Get all permissions
// @Tags Role
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} authpb.PermissionListResponse
// @Router /permission/all [get]
func (h *PermissionHandler) GetAllPermissions(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionListResponse, error) {
	permissions, err := h.permissionUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	pbPermissions := make([]*sharepb.Permission, len(permissions))
	for i, permission := range permissions {
		pbPermissions[i] = h.roleMapper.MapPermissionListWithChild(permission)
	}

	return &authpb.PermissionListResponse{
		Data:          pbPermissions,
		TotalElements: uint64(len(pbPermissions)),
	}, nil
}

// @Summary Get permission keys
// @Description Get all permission keys
// @Tags Permission
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} authpb.PermissionKeysResponse
// @Router /permission/me [get]
func (h *PermissionHandler) GetPermissionKeys(ctx context.Context, req *sharepb.Empty) (*authpb.PermissionKeysResponse, error) {
	keys, err := h.roleUsecase.GetPermissionKeys(ctx)
	if err != nil {
		return nil, err
	}

	return &authpb.PermissionKeysResponse{
		Keys: keys,
	}, nil
}

// // // mapPermissionToPb chuyển domain Permission sang protobuf Permission
// func (h *PermissionHandler) mapPermissionToPb(permission *domain.Permission) *authpb.Permission {
// 	if permission == nil {
// 		return nil
// 	}

// 	pbPermission := &authpb.Permission{
// 		Id:             uint64(permission.ID),
// 		Name:           permission.Name,
// 		Key:            permission.Key,
// 		Description:    permission.Description,
// 		PermissionType: permission.PermissionType,
// 		Module:         permission.Module,
// 		ParentId:       permission.ParentID,
// 		CreatedAt:      _utils.FormatTimeToString(permission.CreatedAt),
// 		UpdatedAt:      _utils.FormatTimeToString(permission.UpdatedAt),
// 	}

// 	if permission.Children != nil {
// 		pbPermission.Children = make([]*authpb.Permission, len(permission.Children))
// 		for i, child := range permission.Children {
// 			pbPermission.Children[i] = h.mapPermissionToPb(child)
// 		}
// 	}

// 	return pbPermission
// }

// @Summary Get permission by group ID
// @Description Get permission by group ID
// @Tags Permission
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} authpb.PermissionKeysResponse
// @Router /permission/group/{id} [get]
func (h *PermissionHandler) GetPermissionByGroupId(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionKeysResponse, error) {
	permissions, err := h.permissionUsecase.GetByGroupId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &authpb.PermissionKeysResponse{
		Keys: permissions,
	}, nil
}

// @Summary Get permission by deal ID
// @Description Get permission by deal ID
// @Tags Permission
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Deal ID"
// @Success 200 {object} authpb.PermissionKeysResponse
// @Router /permission/deal/{id} [get]
func (h *PermissionHandler) GetPermissionByDealId(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionKeysResponse, error) {
	permissions, role, err := h.permissionUsecase.GetByDealId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &authpb.PermissionKeysResponse{
		Keys: permissions,
		Role: h.roleMapper.MapRoleToPb(role),
	}, nil
}

// @Summary Get permission by branch ID
// @Description Get permission by branch ID
// @Tags Permission
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Branch ID"
// @Success 200 {object} authpb.PermissionKeysResponse
// @Router /permission/branch/{id} [get]
func (h *PermissionHandler) GetPermissionByBranchId(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionKeysResponse, error) {
	permissions, err := h.permissionUsecase.GetByBranchId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &authpb.PermissionKeysResponse{
		Keys: permissions,
	}, nil
}

// GetPermissionByOrganizationId trả về union permission của toàn bộ role
// thuộc tổ chức được yêu cầu, tương thích với API production cũ.
func (h *PermissionHandler) GetPermissionByOrganizationId(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionKeysResponse, error) {
	keys, err := h.permissionUsecase.GetPermissionKeysByOrganizationID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &authpb.PermissionKeysResponse{Keys: keys}, nil
}
