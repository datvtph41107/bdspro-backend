package repository

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type OrganizationRoleRepository interface {
	Create(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error)
	GetAll(ctx context.Context) ([]*entity.OrganizationRole, error)
	FindById(ctx context.Context, id uint32) (*entity.OrganizationRole, error)
	Update(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error)
	Delete(ctx context.Context, id uint32) error
	GetByKey(ctx context.Context, key string) (*entity.OrganizationRole, error)
	FindByDomain(ctx context.Context, domain enums.DomainType) ([]*entity.OrganizationRole, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationRole, error)
	FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationRole, error)
	FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationRole, uint32, error)
	AddPermissionToRole(ctx context.Context, roleId uint32, permissionId uint32) error
	GetPermissionKeysByRoleId(ctx context.Context, roleId uint64) ([]string, error)
}
