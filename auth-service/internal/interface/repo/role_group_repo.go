package repo

import (
	"auth/internal/domain"
	"auth/internal/enums"
	"common/case/crud"
	"context"
)

type RoleGroupRepository interface {
	crud.ICrudRepo[domain.RoleGroup]
	GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*domain.RoleGroup, error)

	// UpdateGroupPermissions cập nhật permission cho group
	UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error

	// GetPermissionsByGroupKey lấy danh sách permission theo group key
	GetPermissionsByGroupKey(ctx context.Context, groupKey enums.GroupRoleKeyEnum) ([]*domain.Permission, error)
}
