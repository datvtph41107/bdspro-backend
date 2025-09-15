package postgre

import (
	"common/case/crud"
	_db "common/db"
	_dto "common/dto"
	_errors "common/errors"
	"context"

	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/interface/repo"

	"gorm.io/gorm"
)

type RoleRepo struct {
	*crud.CrudRepo[domain.Role]
}

func NewRoleRepo(db *_db.TransactionRepo) repo.RoleRepository {
	x := &RoleRepo{
		CrudRepo: &crud.CrudRepo[domain.Role]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *RoleRepo) ExistProfileIdWithKeys(ctx context.Context, profileId uint64, keys []string) (bool, error) {
	var exists bool

	// raw SQL dùng EXISTS để check
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM role_profiles rp
            JOIN roles r ON r.id = rp.role_id
            JOIN role_permissions rp2 ON rp2.role_id = r.id
            JOIN permissions p ON p.id = rp2.permission_id
            WHERE rp.profile_id = ?
              AND p.key IN (?)
        )
    `
	tx := r.GetDB(ctx).
		Debug().
		Raw(query, profileId, keys).
		Scan(&exists)
	if tx.Error != nil {
		return false, tx.Error
	}

	return exists, nil
}

func (r *RoleRepo) FindByDomain(ctx context.Context, domainType string) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.GetDB(ctx).Where("domain_type = ?", domainType).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.GetDB(ctx).Where("organization_id = ?", organizationId).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*domain.Role, error) {
	var role domain.Role
	err := r.GetDB(ctx).Where("organization_id = ? AND key = ?", organizationId, key).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*domain.Role, uint64, error) {
	var roles []*domain.Role
	var total int64

	offset := (page - 1) * size

	err := r.GetDB(ctx).Where("organization_id = ?", organizationId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.GetDB(ctx).Where("organization_id = ?", organizationId).Offset(offset).Limit(size).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, uint64(total), nil
}

func (r *RoleRepo) FindByRoleKeys(ctx context.Context, roleKeys []uint32) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.GetDB(ctx).
		Model(&domain.Role{}).
		Preload("Color").
		Where("role_key IN (?)", roleKeys).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) AddPermissionToRole(ctx context.Context, roleID uint64, permissionID uint64) error {
	// Implementation for adding permission to role
	// This might involve a junction table or updating the role's permission list
	return nil
}

func (r *RoleRepo) AfterSave(ctx context.Context, id *uint64, role *domain.Role) error {
	return r.GetDB(ctx).Model(role).Association("Permissions").Replace(role.Permissions)
}

func (r *RoleRepo) GetDetail(ctx context.Context, id uint64) (*domain.Role, error) {
	var role domain.Role
	err := r.GetDB(ctx).
		Preload("Color").
		Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) BeforeSave(ctx context.Context, id *uint64, role *domain.Role) error {
	var permissions []*domain.Permission
	if len(role.PermissionIDs) > 0 {
		if err := r.GetDB(ctx).
			Where("id IN ?", role.PermissionIDs).
			Find(&permissions).Error; err != nil {
			return err
		}
	}

	// Gán lại vào field Permissions để gorm xử lý many2many
	role.Permissions = permissions
	return nil
}

func (r *RoleRepo) AddCondition(db *gorm.DB, pagable *dto.RoleSearchRequest) *gorm.DB {
	if pagable.IsDefault != nil {
		db = db.Where("is_default = ?", *pagable.IsDefault)
	}

	if pagable.DomainType != 0 {
		db = db.Where("domain_type = ?", pagable.DomainType)
	}

	if pagable.Name != "" {
		db = db.Where("name LIKE ?", "%"+pagable.Name+"%")
	}

	return db
}

func (r *RoleRepo) GetList(ctx context.Context, pagable _dto.IPagable) ([]domain.Role, int64, error) {
	param := pagable.(*dto.RoleSearchRequest)

	var roles []domain.Role
	var total int64

	query := r.GetDB(ctx).
		Model(&domain.Role{}).
		Preload("RoleGroup")
	query = r.AddCondition(query, param)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// 1. Lấy roleIds của user
func (r *RoleRepo) GetRoleIdsByProfileId(ctx context.Context, profileId uint64) ([]uint64, error) {
	var roleIds []uint64
	err := r.GetDB(ctx).
		Debug().
		Table("role_profiles").
		Select("role_id").
		Where("profile_id = ? AND deleted_at IS NULL", profileId).
		Pluck("role_id", &roleIds).Error
	return roleIds, err
}

// 2. Lấy permission_ids từ roleIds
func (r *RoleRepo) GetPermissionIdsByRoleIds(ctx context.Context, roleIds []uint64) (map[uint64][]uint64, error) {
	result := make(map[uint64][]uint64)
	if len(roleIds) == 0 {
		return result, nil
	}

	rows, err := r.GetDB(ctx).
		Table("role_permissions").
		Select("role_id, permission_id").
		Where("role_id IN ?", roleIds).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roleId, permId uint64
		if err := rows.Scan(&roleId, &permId); err != nil {
			return nil, err
		}
		result[roleId] = append(result[roleId], permId)
	}
	return result, nil
}

// 3. Lấy permission keys từ permission ids
func (r *RoleRepo) GetPermissionKeysByIds(ctx context.Context, permIds []uint64) ([]string, error) {
	if len(permIds) == 0 {
		return nil, nil
	}
	var keys []string
	err := r.GetDB(ctx).
		Table("permissions").
		Select("key").
		Where("id IN ?", permIds).
		Pluck("key", &keys).Error
	return keys, err
}

func (r *RoleRepo) GetPermissionKeysByRoleId(ctx context.Context, roleId uint64) ([]string, error) {
	var permIds []string
	err := r.GetDB(ctx).
		Table("role_permissions").
		Select("key").
		Where("role_id = ?", roleId).
		Joins("INNER JOIN permissions p ON p.id = role_permissions.permission_id AND p.deleted_at IS NULL").
		Pluck("p.key", &permIds).Error
	return permIds, err
}

// GetPermissionKeysByUserId lấy tất cả permission keys của một user dựa vào userId (1 query duy nhất)
func (r *RoleRepo) GetPermissionKeysByUserId(ctx context.Context, userId uint64) ([]string, error) {
	var permissionKeys []string

	err := r.GetDB(ctx).
		Debug().
		Table("role_profiles p").
		Select("DISTINCT perm.key").
		Joins("INNER JOIN role_permissions rperm ON p.role_id = rperm.role_id").
		Joins("INNER JOIN permissions perm ON rperm.permission_id = perm.id").
		Where("p.profile_id = ? AND p.deleted_at IS NULL AND perm.key IS NOT NULL AND perm.key != ''", userId).
		Pluck("perm.key", &permissionKeys).
		Error
	if err != nil {
		return nil, err
	}

	return permissionKeys, nil
}

// GetAdminRoleIdByGroupKey lấy admin role theo group key, ưu tiên role có auto_assign = true
func (r *RoleRepo) GetAdminRoleIdByGroupKey(ctx context.Context, groupKey uint64) (*domain.Role, error) {
	var roles []domain.Role

	// Query để lấy 2 role thuộc group key, ưu tiên auto_assign = true
	err := r.GetDB(ctx).
		Debug().
		Joins("JOIN role_groups ON roles.role_group_id = role_groups.id").
		Where("role_groups.key = ?", groupKey).
		Order("roles.auto_assign DESC, roles.id ASC").
		Limit(2).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}

	// Trả về role đầu tiên (có auto_assign = true nếu có, hoặc role đầu tiên nếu không có)
	return &roles[0], nil
}

// GetAdminRoleIdByGroupKey lấy admin role theo group key, ưu tiên role có auto_assign = true
func (r *RoleRepo) GetListByGroupKey(ctx context.Context, groupKey uint32) ([]domain.Role, int64, error) {
	var roles []domain.Role

	// Query để lấy 2 role thuộc group key, ưu tiên auto_assign = true
	err := r.GetDB(ctx).
		Debug().
		Model(&domain.Role{}).
		// Preload("Color").
		Joins("JOIN role_groups ON roles.role_group_id = role_groups.id").
		Where("role_groups.key = ? and roles.deleted_at IS NULL and roles.allow_assign = true", groupKey).
		Order("roles.allow_assign DESC, roles.id ASC").
		Limit(100).
		Find(&roles).Error
	if err != nil {
		return nil, 0, _errors.ReturnError(500, err.Error())
	}

	// if len(roles) == 0 {
	// 	return nil, 0, _errors.ReturnError(404, "no roles found")
	// }

	// Trả về role đầu tiên (có auto_assign = true nếu có, hoặc role đầu tiên nếu không có)
	return roles, int64(len(roles)), nil
}

// GetRolesByRoleKeysAndOrganization lấy roles theo roleKeys và organizationId
func (r *RoleRepo) GetRolesByRoleKeysAndOrganization(ctx context.Context, roleKeys []uint32, organizationId uint64) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.GetDB(ctx).
		Model(&domain.Role{}).
		Preload("Color").
		Preload("Permissions").
		Preload("RoleGroup").
		Where("role_key IN ? AND organization_id = ? AND deleted_at IS NULL", roleKeys, organizationId).
		Find(&roles).Error
	return roles, err
}

// GetRoleByRoleKey lấy role theo roleKey
func (r *RoleRepo) GetRoleByRoleKey(ctx context.Context, roleKey uint32) (*domain.Role, error) {
	var role domain.Role
	err := r.GetDB(ctx).
		Model(&domain.Role{}).
		Preload("Color").
		Preload("Permissions").
		Preload("RoleGroup").
		Where("role_key = ? AND deleted_at IS NULL", roleKey).
		First(&role).Error
	return &role, err
}

// AssignRoleToUser gán role cho user bằng cách cập nhật role_key trong auth method
func (r *RoleRepo) AssignRoleToUser(ctx context.Context, userId uint64, role *domain.Role) error {
	// Cập nhật role_key trong bảng user_auth
	err := r.GetDB(ctx).
		Model(&domain.AuthMethodDomain{}).
		Where("user_id = ?", userId).
		Update("role_key", role.RoleKey).Error
	if err != nil {
		return err
	}

	return nil
}
