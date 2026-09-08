package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationBranchRepository interface {
	CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	DeleteOrganizationBranch(ctx context.Context, id uint32) error
	GetOrganizationBranchById(ctx context.Context, id uint32) (*entity.OrganizationBranch, error)
	GetOrganizationBranchByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBranch, error)
	GetOrganizationBranchByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationBranch, uint32, error)
	UpdateIsActive(ctx context.Context, id uint32, isActive bool) error
	CountBranchByOrganizationId(ctx context.Context, organizationId uint32) (uint32, error)
}
