package repo

import (
	"common/case/crud"
	"context"
	"user/internal/domain/access"
	"user/internal/enums"
)

type RoleGroupRepository interface {
	crud.ICrudRepo[access.RoleGroup]
	GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*access.RoleGroup, error)

	// UpdateGroupPermissions cập nhật permission cho group
	UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error

	// GetPermissionsByGroupKey lấy danh sách permission theo group key
	GetPermissionsByGroupKey(ctx context.Context, groupKey enums.GroupRoleKeyEnum) ([]*access.Permission, error)
}
