package postgre

import (
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/interface/repo"
	"common/case/crud"
	_db "common/db"
	_dto "common/dto"
	"context"

	"gorm.io/gorm"
)

type RoleRepo struct {
	crud.BaseRepo[domain.Role]
}

func NewRoleRepo(db *_db.DBTransaction) repo.RoleRepository {
	return crud.MakeCrudRepo(&RoleRepo{}, db).(*RoleRepo)
}

func (r *RoleRepo) GetByKey(ctx context.Context, key string) (*domain.Role, error) {
	var role domain.Role
	err := r.GetDB(ctx).Where("key = ?", key).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
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
		Preload("Permissions").
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
