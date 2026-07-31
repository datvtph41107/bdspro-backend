package repo

import (
	"context"
	"user/internal/domain"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *domain.Permission) (*domain.Permission, error)
	GetPermissions(ctx context.Context) ([]*domain.Permission, error)
	FindById(ctx context.Context, id uint64) (*domain.Permission, error)
	Update(ctx context.Context, permission *domain.Permission) (*domain.Permission, error)
	Delete(ctx context.Context, id uint64) error
	GetByKey(ctx context.Context, key string) (*domain.Permission, error)
	FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*domain.Permission, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*domain.Permission, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*domain.Permission, uint64, error)
	FindByModule(ctx context.Context, module string) ([]*domain.Permission, error)
	FindByModuleWithPagination(ctx context.Context, module uint32, page, size int) ([]*domain.Permission, uint64, error)
}
