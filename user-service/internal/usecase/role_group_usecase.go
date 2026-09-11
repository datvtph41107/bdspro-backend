package usecase

import (
	"common/case/crud"
	"context"
	"user/internal/domain/access"
	"user/internal/enums"
	"user/internal/interface/repo"
)

type RoleGroupUsecase struct {
	crud.BaseUsecase[access.RoleGroup, repo.RoleGroupRepository]
	permissionUsecase *PermissionUsecase
}

func NewRoleGroupUsecase(
	roleGroupRepo repo.RoleGroupRepository,
	permissionUsecase *PermissionUsecase,
) *RoleGroupUsecase {
	return &RoleGroupUsecase{
		crud.BaseUsecase[access.RoleGroup, repo.RoleGroupRepository]{
			Repo: roleGroupRepo,
		},
		permissionUsecase,
	}
}

// UpdateGroupPermissions cập nhật toàn bộ permission cho group
func (uc *RoleGroupUsecase) UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error {
	// Kiểm tra group có tồn tại không
	group, err := uc.Repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return access.ErrRoleGroupNotFound
	}

	// Cập nhật permission cho group
	return uc.Repo.UpdateGroupPermissions(ctx, groupID, permissionIDs)
}

// GetPermissionsByGroupKey lấy danh sách permission theo group key
func (uc *RoleGroupUsecase) GetPermissionsByGroupKey(ctx context.Context, groupKey enums.GroupRoleKeyEnum) ([]*access.Permission, error) {
	permissions, err := uc.Repo.GetPermissionsByGroupKey(ctx, groupKey)
	if err != nil {
		return nil, err
	}

	permMap := make(map[uint64]*access.Permission)
	for i := range permissions {
		permMap[permissions[i].ID] = permissions[i]
	}

	var roots []*access.Permission

	// Xây dựng quan hệ cha-con
	for _, p := range permissions {
		permission := uc.permissionUsecase.BuildHierarchy(p, permMap)
		if permission.ParentID == nil {
			roots = append(roots, permission)
		}
	}

	return roots, nil
}
