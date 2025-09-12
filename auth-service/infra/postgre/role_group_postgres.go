package postgre

import (
	"auth/internal/domain"
	"auth/internal/enums"
	"auth/internal/interface/repo"
	"common/case/crud"
	_db "common/db"
	"context"
)

type RoleGroupRepo struct {
	*crud.CrudRepo[domain.RoleGroup]
}

func NewRoleGroupRepo(db *_db.TransactionRepo) repo.RoleGroupRepository {
	x := &RoleGroupRepo{
		CrudRepo: &crud.CrudRepo[domain.RoleGroup]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *RoleGroupRepo) GetByKey(ctx context.Context, key enums.GroupRoleKeyEnum) (*domain.RoleGroup, error) {
	var roleGroup domain.RoleGroup
	err := r.GetDB(ctx).Where("key = ?", key).First(&roleGroup).Error
	if err != nil {
		return nil, err
	}
	return &roleGroup, nil
}

// UpdateGroupPermissions cập nhật permission cho group
func (r *RoleGroupRepo) UpdateGroupPermissions(ctx context.Context, groupID uint64, permissionIDs []uint64) error {
	// Lấy group
	var group domain.RoleGroup
	err := r.GetDB(ctx).First(&group, groupID).Error
	if err != nil {
		return err
	}

	// Lấy danh sách permission
	var permissions []domain.Permission
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
func (r *RoleGroupRepo) GetPermissionsByGroupKey(ctx context.Context, groupKey enums.GroupRoleKeyEnum) ([]*domain.Permission, error) {
	var permissions []*domain.Permission

	err := r.GetDB(ctx).
		Table("permissions").
		Joins("JOIN group_permissions ON permissions.id = group_permissions.permission_id").
		Joins("JOIN role_groups ON group_permissions.role_group_id = role_groups.id").
		Where("role_groups.key = ?", groupKey).
		Order("permissions.priority ASC, permissions.id ASC").
		Find(&permissions).Error

	return permissions, err
}
