package postgres

import (
	"common/case/crud"
	_db "common/db"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"errors"
	"time"

	"user/internal/domain/access"

	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/interface/repo"

	"gorm.io/gorm"
)

type RoleRepo struct {
	*crud.CrudRepo[access.Role]
}

func NewRoleRepo(db *_db.TransactionRepo) repo.RoleRepository {
	x := &RoleRepo{
		CrudRepo: &crud.CrudRepo[access.Role]{},
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
            WHERE rp.profile_id = ?
              AND rp.deleted_at IS NULL
              AND r.deleted_at IS NULL
              AND (
                  (
                      r.organization_id = 0
                      AND r.key = 'QHPRO_SYSTEM_ROOT'
                      AND EXISTS (
                          SELECT 1 FROM permissions root_permission
                           WHERE root_permission.deleted_at IS NULL
                             AND root_permission.key IN (?)
                      )
                  )
                  OR EXISTS (
                      SELECT 1
                        FROM role_permissions rp2
                        JOIN permissions p ON p.id = rp2.permission_id
                       WHERE rp2.role_id = r.id
                         AND p.deleted_at IS NULL
                         AND p.key IN (?)
                  )
              )
        )
    `
	tx := r.GetDB(ctx).
		Debug().
		Raw(query, profileId, keys, keys).
		Scan(&exists)
	if tx.Error != nil {
		return false, tx.Error
	}

	return exists, nil
}

func (r *RoleRepo) FindByDomain(ctx context.Context, domainType string) ([]*access.Role, error) {
	var roles []*access.Role
	err := r.GetDB(ctx).Where("domain_type = ?", domainType).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*access.Role, error) {
	var roles []*access.Role
	err := r.GetDB(ctx).Where("organization_id = ?", organizationId).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*access.Role, error) {
	var role access.Role
	err := r.GetDB(ctx).Where("organization_id = ? AND key = ?", organizationId, key).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*access.Role, uint64, error) {
	var roles []*access.Role
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

func (r *RoleRepo) FindByRoleKeys(ctx context.Context, roleKeys []uint32) ([]access.Role, error) {
	var roles []access.Role
	err := r.GetDB(ctx).
		Model(&access.Role{}).
		Preload("Color").
		Where("role_key IN (?) AND deleted_at IS NULL", roleKeys).
		Find(&roles).
		Error
	return roles, err
}

func (r *RoleRepo) AfterSave(ctx context.Context, id *uint64, role *access.Role) error {
	return r.GetDB(ctx).Model(role).Association("Permissions").Replace(role.Permissions)
}

func (r *RoleRepo) GetDetail(ctx context.Context, id uint64) (*access.Role, error) {
	var role access.Role
	err := r.GetDB(ctx).
		Preload("Color").
		Preload("Permissions").
		Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) BeforeSave(ctx context.Context, id *uint64, role *access.Role) error {
	var permissions []*access.Permission
	if len(role.PermissionIDs) > 0 {
		if err := r.GetDB(ctx).
			Where("id IN ? AND deleted_at IS NULL", role.PermissionIDs).
			Find(&permissions).Error; err != nil {
			return err
		}
		if len(permissions) != len(role.PermissionIDs) {
			return errors.New("one or more permissions do not exist")
		}
	}

	role.Permissions = permissions
	return nil
}

func (r *RoleRepo) AddCondition(db *gorm.DB, pagable *dto.RoleSearchRequest) *gorm.DB {
	if pagable.IsDefault != nil {
		db = db.Where("is_default = ?", *pagable.IsDefault)
	}

	if pagable.DomainType != "" {
		db = db.Where("domain_type = ?", pagable.DomainType)
	}

	if pagable.Name != "" {
		db = db.Where("role_name LIKE ?", "%"+pagable.Name+"%")
	}
	if pagable.RoleGroupID != nil {
		db = db.Where("role_group_id = ?", *pagable.RoleGroupID)
	}

	return db
}

func (r *RoleRepo) GetList(ctx context.Context, pagable _dto.IPagable) ([]access.Role, int64, error) {
	param := pagable.(*dto.RoleSearchRequest)

	var roles []access.Role
	var total int64

	query := r.GetDB(ctx).
		Model(&access.Role{}).
		Preload("RoleGroup").
		Preload("Permissions")
	query = r.AddCondition(query, param)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// GetLimitAdmin cho phép size lớn (tới 500) — màn role default gọi size=999 để lấy gần như all
	limit := pagable.GetLimit()
	if p, ok := pagable.(*dto.RoleSearchRequest); ok {
		limit = p.GetLimitAdmin()
	}
	err = query.
		Order("roles.id DESC").
		Offset(pagable.GetOffset()).
		Limit(limit).
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

func (r *RoleRepo) GetPermissionKeysByRoleKey(ctx context.Context, roleKey uint32) ([]string, error) {
	var permissionKeys []string
	err := r.GetDB(ctx).
		Table("role_permissions").
		Select("p.key").
		Joins("JOIN permissions p ON p.id = role_permissions.permission_id AND p.deleted_at IS NULL").
		Joins("JOIN roles r ON r.id = role_permissions.role_id AND r.deleted_at IS NULL").
		Where("r.role_key = ?", roleKey).
		Pluck("p.key", &permissionKeys).
		Error
	return permissionKeys, err
}

// GetPermissionKeysByUserId lấy tất cả permission keys của một user dựa vào userId (1 query duy nhất)
func (r *RoleRepo) GetPermissionKeysByUserId(ctx context.Context, userId uint64) ([]string, error) {
	var permissionKeys []string

	err := r.GetDB(ctx).
		Raw(`
SELECT DISTINCT permission.key
  FROM permissions permission
 WHERE permission.deleted_at IS NULL
   AND permission.key IS NOT NULL
   AND permission.key != ''
   AND EXISTS (
       SELECT 1
         FROM role_profiles assignment
         JOIN roles role ON role.id = assignment.role_id
        WHERE assignment.profile_id = ?
          AND assignment.deleted_at IS NULL
          AND role.deleted_at IS NULL
          AND (
              (role.organization_id = 0 AND role.key = 'QHPRO_SYSTEM_ROOT')
              OR EXISTS (
                  SELECT 1 FROM role_permissions mapping
                   WHERE mapping.role_id = role.id
                     AND mapping.permission_id = permission.id
              )
          )
   )
ORDER BY permission.key
`, userId).
		Scan(&permissionKeys).
		Error
	if err != nil {
		return nil, err
	}

	return permissionKeys, nil
}

// GetAdminRoleIdByGroupKey lấy admin role theo group key, ưu tiên role có allow_assign = true
func (r *RoleRepo) GetAdminRoleIdByGroupKey(ctx context.Context, groupKey uint64) (*access.Role, error) {
	var roles []access.Role

	// Query để lấy 2 role thuộc group key, ưu tiên allow_assign = true
	err := r.GetDB(ctx).
		Debug().
		Joins("JOIN role_groups ON roles.role_group_id = role_groups.id").
		Where("role_groups.key = ?", groupKey).
		Order("roles.allow_assign DESC, roles.id ASC").
		Limit(2).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, nil
	}

	// Trả về role đầu tiên (có allow_assign = true nếu có, hoặc role đầu tiên nếu không có)
	return &roles[0], nil
}

// GetAdminRoleIdByGroupKey lấy admin role theo group key, ưu tiên role có allow_assign = true
func (r *RoleRepo) GetListByGroupKey(ctx context.Context, groupKey uint32) ([]access.Role, int64, error) {
	var roles []access.Role

	// Query để lấy 2 role thuộc group key, ưu tiên allow_assign = true
	err := r.GetDB(ctx).
		Debug().
		Model(&access.Role{}).
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

	// Trả về role đầu tiên (có allow_assign = true nếu có, hoặc role đầu tiên nếu không có)
	return roles, int64(len(roles)), nil
}

// GetRolesByRoleKeysAndOrganization lấy roles theo roleKeys và organizationId
func (r *RoleRepo) GetRolesByRoleKeysAndOrganization(ctx context.Context, roleKeys []uint32, organizationId uint64) ([]*access.Role, error) {
	var roles []*access.Role
	err := r.GetDB(ctx).
		Model(&access.Role{}).
		Preload("Color").
		Preload("Permissions").
		Preload("RoleGroup").
		Where("role_key IN ? AND organization_id = ? AND deleted_at IS NULL", roleKeys, organizationId).
		Find(&roles).Error
	return roles, err
}

// GetRoleByRoleKey lấy role theo roleKey
func (r *RoleRepo) GetRoleByRoleKey(ctx context.Context, roleKey uint32) (*access.Role, error) {
	var role access.Role
	err := r.GetDB(ctx).
		Model(&access.Role{}).
		Preload("Color").
		Preload("Permissions").
		Preload("RoleGroup").
		Where("role_key = ? AND deleted_at IS NULL", roleKey).
		First(&role).Error
	return &role, err
}

// AssignRoleToUser gán role cho user: cập nhật auth_method.role_key + ghi role_profiles
// (permission/me và HasPermissions đều đọc từ role_profiles theo profile_id)
func (r *RoleRepo) AssignRoleToUser(ctx context.Context, userId uint64, role *access.Role) error {
	db := r.GetDB(ctx).WithContext(ctx)

	if err := db.Model(&auth.AuthMethod{}).
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Update("role_key", role.RoleKey).Error; err != nil {
		return err
	}

	// Soft-delete mapping cũ rồi gắn role mới (admin create/update chọn 1 role)
	now := time.Now()
	if err := db.Model(&access.RoleProfile{}).
		Where("profile_id = ? AND deleted_at IS NULL", userId).
		Update("deleted_at", now).Error; err != nil {
		return err
	}

	rp := &access.RoleProfile{
		ProfileID: userId,
		RoleID:    role.ID,
	}
	if err := db.Create(rp).Error; err != nil {
		return err
	}

	return nil
}

// GetAuthMethodsByIds lấy thông tin AuthMethod theo danh sách user IDs
func (r *RoleRepo) GetAuthMethodsByIds(ctx context.Context, userIds []uint64) ([]*auth.AuthMethod, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	var authMethods []*auth.AuthMethod
	err := r.GetDB(ctx).
		Where("user_id IN ? AND deleted_at IS NULL", userIds).
		Find(&authMethods).Error
	if err != nil {
		return nil, err
	}

	return authMethods, nil
}

// GetAuthAdminsWithRolesByIds lấy danh sách admin kèm với thông tin role (theo auth_method.user_id).
func (r *RoleRepo) GetAuthAdminsWithRolesByIds(ctx context.Context, ids []uint64) ([]*access.AuthAdminWithRole, error) {
	if len(ids) == 0 {
		return []*access.AuthAdminWithRole{}, nil
	}
	if r == nil || r.CrudRepo == nil || r.TransactionRepo == nil {
		return nil, errors.New("role repository chưa được khởi tạo")
	}

	var authMethods []*auth.AuthMethod
	err := r.GetDB(ctx).
		Where("user_id IN ? AND provider = ? AND deleted_at IS NULL", ids, auth.ProviderAdmin).
		Find(&authMethods).Error
	if err != nil {
		return nil, err
	}
	if len(authMethods) == 0 {
		return []*access.AuthAdminWithRole{}, nil
	}

	roleKeySet := make(map[uint32]struct{})
	for _, am := range authMethods {
		if am != nil && am.RoleKey != 0 {
			roleKeySet[am.RoleKey] = struct{}{}
		}
	}

	roleByKey := make(map[uint32]access.Role)
	if len(roleKeySet) > 0 {
		roleKeys := make([]uint32, 0, len(roleKeySet))
		for roleKey := range roleKeySet {
			roleKeys = append(roleKeys, roleKey)
		}
		roles, err := r.FindByRoleKeys(ctx, roleKeys)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			roleByKey[role.RoleKey] = role
		}
	}

	admins := make([]*access.AuthAdminWithRole, 0, len(authMethods))
	for _, am := range authMethods {
		if am == nil {
			continue
		}
		admin := &access.AuthAdminWithRole{
			ID:        am.ID,
			UserID:    am.UserID,
			Username:  am.AuthName,
			FullName:  am.FullName,
			Email:     am.Email,
			Phone:     am.Phone,
			Avatar:    am.Avatar,
			RoleKey:   am.RoleKey,
			Status:    uint32(am.Status),
			CreatedAt: am.CreatedAt,
			UpdatedAt: am.UpdatedAt,
		}
		if role, ok := roleByKey[am.RoleKey]; ok {
			roleCopy := role
			admin.Role = &roleCopy
		}
		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *RoleRepo) GetRolesByIds(ctx context.Context, ids []uint64) ([]*access.Role, error) {
	var roles []*access.Role
	err := r.GetDB(ctx).
		Preload("Color").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) GetRoleByGroupId(ctx context.Context, groupId uint64) ([]access.Role, error) {
	var roles []access.Role
	err := r.GetDB(ctx).
		Model(&access.Role{}).
		Preload("Color").
		Where("role_group_id = ? AND deleted_at IS NULL", groupId).
		Find(&roles).Error
	return roles, err
}
