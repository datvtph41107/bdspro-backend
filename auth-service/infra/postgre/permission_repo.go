package postgre

import (
	"auth/internal/domain"
	"auth/internal/interface/repo"
	"context"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) repo.PermissionRepository {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) Create(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	err := r.db.WithContext(ctx).Create(permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *PermissionRepo) GetPermissions(ctx context.Context) ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) FindById(ctx context.Context, id uint64) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepo) Update(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	err := r.db.WithContext(ctx).Save(permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *PermissionRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Permission{}, id).Error
}

func (r *PermissionRepo) GetByKey(ctx context.Context, key string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepo) FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).Where("organization_id = ? AND key = ?", organizationId, key).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepo) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*domain.Permission, uint64, error) {
	var permissions []*domain.Permission
	var total int64

	offset := (page - 1) * size

	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Offset(offset).Limit(size).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, uint64(total), nil
}

func (r *PermissionRepo) FindByModule(ctx context.Context, module string) ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	err := r.db.WithContext(ctx).Where("module = ?", module).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) FindByModuleWithPagination(ctx context.Context, module uint32, page, size int) ([]*domain.Permission, uint64, error) {
	var permissions []*domain.Permission
	var total int64

	offset := (page - 1) * size

	err := r.db.WithContext(ctx).Where("module = ?", module).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Where("module = ?", module).Offset(offset).Limit(size).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, uint64(total), nil
}
