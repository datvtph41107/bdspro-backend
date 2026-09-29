package providers

import (
	"context"
	"user/internal/dto"
)

type OrganizationProvider interface {
	GetGroupRoleUser(ctx context.Context, groupId uint64, userId uint64) (*dto.GroupMember, error)
	GetBranchMember(ctx context.Context, branchId uint64, userId uint64) (*dto.BranchMember, error)
	GetOrganizationMember(ctx context.Context, organizationId uint64, userId uint64) (*dto.OrganizationMember, error)
}
