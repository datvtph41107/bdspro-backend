package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationRepository interface {
	CreateOrganizationWithMember(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	FindByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error)
	UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)
	FindById(ctx context.Context, id uint32) (*entity.Organization, error)
	FindByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error)
}
