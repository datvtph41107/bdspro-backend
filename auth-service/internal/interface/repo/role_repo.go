package repo

import (
	"auth/internal/domain"
	"common/case/crud"
	"context"
)

type RoleRepository interface {
	crud.IBaseRepo[domain.Role]
	GetByKey(ctx context.Context, key string) (*domain.Role, error)
	FindByDomain(ctx context.Context, domain string) ([]*domain.Role, error)
	FindByOrganizationId(ctx context.Context, organizationId uint64) ([]*domain.Role, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint64, key string) (*domain.Role, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint64, page, size int) ([]*domain.Role, uint64, error)
	AddPermissionToRole(ctx context.Context, roleID uint64, permissionID uint64) error
}
