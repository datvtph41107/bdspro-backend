package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationPermissionRepository interface {
	Create(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	GetPermissionOrganization(ctx context.Context) ([]*entity.OrganizationPermission, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationPermission, error)
	Update(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	Delete(ctx context.Context, id uint32) error
	GetByKey(ctx context.Context, key string) (*entity.OrganizationPermission, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationPermission, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationPermission, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationPermission, uint32, error)
}
