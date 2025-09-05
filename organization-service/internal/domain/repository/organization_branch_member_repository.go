package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type OrganizationBranchMemberRepository interface {
	CreateOrganizationBranchMember(ctx context.Context, organizationBranchMember *entity.OrganizationBranchMember) (*entity.OrganizationBranchMember, error)
	DeleteOrganizationBranchMember(ctx context.Context, id uint32) error
	GetOrganizationBranchMemberByOrganizationBranchId(ctx context.Context, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error)
	GetOrganizationBranchMemberByOrganizationBranchIdWithPagination(ctx context.Context, organizationBranchId uint32, page, size int) ([]*entity.OrganizationBranchMember, uint32, error)
	GetOrganizationBranchMemberByUserIdAndOrganizationBranchId(ctx context.Context, userId uint32, organizationBranchId uint32) (*entity.OrganizationBranchMember, error)
	GetOrganizationBranchMemberByUserIdsAndOrganizationBranchId(ctx context.Context, userIds []uint64, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error)
}
