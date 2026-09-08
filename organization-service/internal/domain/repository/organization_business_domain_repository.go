package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationBusinessDomainRepository interface {
	Create(ctx context.Context, organizationBusinessDomain *entity.OrganizationBusinessDomain) (*entity.OrganizationBusinessDomain, error)
	FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBusinessDomain, error)
	FindByBusinessDomainId(ctx context.Context, businessDomainId uint32) ([]*entity.OrganizationBusinessDomain, error)
	DeleteByOrganizationId(ctx context.Context, organizationId uint32) error
	DeleteByOrganizationIdAndBusinessDomainId(ctx context.Context, organizationId uint32, businessDomainId uint32) error
	FindByOrganizationIds(ctx context.Context, organizationIds []uint32) ([]*entity.OrganizationBusinessDomain, error)
}
