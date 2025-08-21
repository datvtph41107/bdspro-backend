package handler

import (
	"auth/infra/mapper"
	"auth/internal/usecase"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
)

type PermissionHandler struct {
	authpb.UnimplementedPermissionServiceServer
	permissionUsecase *usecase.PermissionUsecase
	roleMapper        *mapper.RoleMapper
}

func NewPermissionHandler(
	permissionUsecase *usecase.PermissionUsecase,
	roleMapper *mapper.RoleMapper,
) *PermissionHandler {
	return &PermissionHandler{
		permissionUsecase: permissionUsecase,
		roleMapper:        roleMapper,
	}
}

// GetPermissionById lấy permission theo ID
func (h *PermissionHandler) GetPermissionById(ctx context.Context, req *sharepb.IdRequest) (*authpb.Permission, error) {
	permission, err := h.permissionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.roleMapper.MapPermissionToPb(permission), nil
}

// GetPermissionsByModule lấy danh sách permission theo module
func (h *PermissionHandler) GetPermissionsByModule(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionListResponse, error) {
	permissions, total, err := h.permissionUsecase.GetByModuleWithPagination(ctx, uint32(req.Id), int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	pbPermissions := make([]*authpb.Permission, len(permissions))
	for i, permission := range permissions {
		pbPermissions[i] = h.roleMapper.MapPermissionListWithChild(permission)
	}

	return &authpb.PermissionListResponse{
		Data:          pbPermissions,
		TotalElements: total,
	}, nil
}

// GetAllPermissions lấy tất cả permission
func (h *PermissionHandler) GetAllPermissions(ctx context.Context, req *sharepb.IdRequest) (*authpb.PermissionListResponse, error) {
	permissions, err := h.permissionUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	pbPermissions := make([]*authpb.Permission, len(permissions))
	for i, permission := range permissions {
		pbPermissions[i] = h.roleMapper.MapPermissionListWithChild(permission)
	}

	return &authpb.PermissionListResponse{
		Data:          pbPermissions,
		TotalElements: uint64(len(pbPermissions)),
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
