package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationRepository interface {
	CreateOrganizationWithMember(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	ExistByTaxCode(ctx context.Context, taxCode string) (bool, error)
}
