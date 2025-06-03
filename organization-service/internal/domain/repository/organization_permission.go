package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationPermissionRepository interface {
	Create(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationPermission, error)
	Update(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error)
	Delete(ctx context.Context, id uint32) error
	FindByKeyAndOrganizationId(ctx context.Context, key string, organizationId uint32) (*entity.OrganizationPermission, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationPermission, error)
}