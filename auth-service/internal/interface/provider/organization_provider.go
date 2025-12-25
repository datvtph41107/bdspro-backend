package provider

import (
	"auth/internal/dto"
	"context"
)

type OrganizationProvider interface {
	SwitchOrganization(ctx context.Context, organizationID uint64) (*dto.SwitchOrganizationResponse, error)

	// Admin role management
	GetAdminRoles(ctx context.Context) ([]dto.AdminRoleDTO, error)
	GetGroupRoleUser(ctx context.Context, groupId uint64, userId uint64) (*dto.GroupMember, error)
	GetBranchMember(ctx context.Context, branchId uint64, userId uint64) (*dto.BranchMember, error)
	GetOrganizationMember(ctx context.Context, organizationId uint64, userId uint64) (*dto.OrganizationMember, error)
}
