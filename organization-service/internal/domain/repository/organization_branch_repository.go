package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationBranchRepository interface {
	CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error)
	DeleteOrganizationBranch(ctx context.Context, id uint32) error
	GetOrganizationBranchById(ctx context.Context, id uint32) (*entity.OrganizationBranch, error)
	GetOrganizationBranchByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBranch, error)
	GetOrganizationBranchByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationBranch, uint32, error)
}
