package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationRoleRepository interface {
	Create(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error)
	GetAll(ctx context.Context) ([]*entity.OrganizationRole, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationRole, error)
	Update(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error)
	Delete(ctx context.Context, id uint32) error
	GetByKey(ctx context.Context, key string) (*entity.OrganizationRole, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationRole, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationRole, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationRole, uint32, error)

}
