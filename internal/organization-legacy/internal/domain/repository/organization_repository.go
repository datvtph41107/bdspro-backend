package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationRepository interface {
	CreateOrganizationWithMember(ctx context.Context, organization *entity.Organization, member *entity.OrganizationMember) (*entity.Organization, *entity.OrganizationMember, error)
	FindByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error)
	UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error)
	FindById(ctx context.Context, id uint32) (*entity.Organization, error)
	FindByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error)
	FindByUserIdWithPagination(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*entity.Organization, error)
	GetAllOrganizationsWithPagination(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error)
	GetAllOrganizationsWithMemberCounts(ctx context.Context, offset, limit int) ([]*entity.OrganizationWithMemberCount, uint32, error)
}
