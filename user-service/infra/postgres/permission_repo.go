package postgres

import (
	"context"

	"gorm.io/gorm"
	"user/internal/domain/access"
)

// PermissionRepo implements the global User IAM permission read model.
// @bind: user/internal/interface/repo.PermissionRepository
type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) GetPermissions(ctx context.Context) ([]*access.Permission, error) {
	var permissions []*access.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepo) FindById(ctx context.Context, id uint64) (*access.Permission, error) {
	var permission access.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByModuleWithPagination giữ truy vấn module ở đúng repository owner,
// tránh để handler tự ghép SQL và đảm bảo contract phân trang production.
func (r *PermissionRepo) FindByModuleWithPagination(ctx context.Context, module string, page, size int) ([]*access.Permission, uint64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var permissions []*access.Permission
	var total int64
	query := r.db.WithContext(ctx).Where("module = ?", module)
	if err := query.Model(&access.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * size).Limit(size).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	return permissions, uint64(total), nil
}
