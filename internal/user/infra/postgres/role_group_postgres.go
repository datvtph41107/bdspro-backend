package postgres

import (
	"common/case/crud"
	_db "common/db"
	"context"
	"user/internal/domain/access"
	"user/internal/enums"
	"user/internal/interface/repo"
)

type RoleGroupRepo struct {
	*crud.CrudRepo[access.RoleGroup]
}

func NewRoleGroupRepo(db *_db.TransactionRepo) repo.RoleGroupRepository {
	x := &RoleGroupRepo{
		CrudRepo: &crud.CrudRepo[access.RoleGroup]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *RoleGroupRepo) GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*access.RoleGroup, error) {
	var roleGroup access.RoleGroup
	err := r.GetDB(ctx).Where("key = ?", key).First(&roleGroup).Error
	if err != nil {
		return nil, err
	}
	return &roleGroup, nil
}

// UpdateGroupPermissions cập nhật permission cho group
func (r *RoleGroupRepo) UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error {
	// Lấy group
	var group access.RoleGroup
	err := r.GetDB(ctx).First(&group, groupID).Error
	if err != nil {
		return err
	}

	// Lấy danh sách permission
	var permissions []access.Permission
	if len(permissionIDs) > 0 {
		err = r.GetDB(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error
		if err != nil {
			return err
		}
	}

	// Cập nhật quan hệ many-to-many
	return r.GetDB(ctx).Model(&group).Association("Permissions").Replace(permissions)
}

// GetPermissionsByGroupKey lấy danh sách permission theo group key
func (r *RoleGroupRepo) GetPermissionsByGroupKey(ctx context.Context, groupKey enums.GroupRoleKeyEnum) ([]*access.Permission, error) {
	var permissions []*access.Permission

	err := r.GetDB(ctx).
		Table("permissions").
		Joins("JOIN group_permissions ON permissions.id = group_permissions.permission_id").
		Joins("JOIN role_groups ON group_permissions.role_group_id = role_groups.id").
		Where("role_groups.key = ?", groupKey).
		Order("permissions.priority ASC, permissions.id ASC").
		Find(&permissions).Error

	return permissions, err
}
