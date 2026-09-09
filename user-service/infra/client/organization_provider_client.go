package client

import (
	"context"
	"pb/clients"
	organizationpb "pb/types/organization"
	"user/internal/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// @bind: user/internal/interface/providers.OrganizationProvider
type OrganizationProviderClient struct {
	*clients.OrganizationClient
	database *gorm.DB
}

func NewOrganizationProviderClient(
	rpcClient *clients.OrganizationClient,
	database *gorm.DB,
) *OrganizationProviderClient {
	return &OrganizationProviderClient{
		OrganizationClient: rpcClient,
		database:           database,
	}
}

func (c *OrganizationProviderClient) GetGroupRoleUser(ctx context.Context, groupId uint64, userId uint64) (*dto.GroupMember, error) {
	resp, err := c.OrganizationClient.Client.GetGroupMemberByIds(ctx, &organizationpb.GetMemberByIdsRequest{
		GroupId: groupId,
		Ids:     []uint64{userId},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get group role user: %v", err)
	}
	if resp == nil || len(resp.Data) == 0 {
		return nil, nil
	}

	return &dto.GroupMember{
		ID:      uint64(resp.Data[0].Id),
		GroupId: uint64(resp.Data[0].GroupId),
		UserId:  uint64(resp.Data[0].UserId),
		Role:    resp.Data[0].Role,
		RoleId:  resp.Data[0].RoleId,
	}, nil
}

func (c *OrganizationProviderClient) GetBranchMember(ctx context.Context, branchId uint64, userId uint64) (*dto.BranchMember, error) {
	resp, err := c.OrganizationClient.GetBranchMember(ctx, branchId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get branch member role: %v", err)
	}
	if resp == nil {
		return nil, nil
	}

	return &dto.BranchMember{
		BranchId: branchId,
		UserId:   userId,
		RoleId:   &resp.RoleId,
	}, nil
}

func (c *OrganizationProviderClient) GetOrganizationMember(ctx context.Context, organizationId uint64, userId uint64) (*dto.OrganizationMember, error) {
	if c == nil || c.database == nil {
		return nil, status.Error(codes.Unavailable, "user organization membership store is unavailable")
	}
	var row struct {
		RoleKey string `gorm:"column:role_key"`
	}
	result := c.database.WithContext(ctx).Table("organization_members AS member").
		Select("member.role_key").
		Joins("JOIN organizations organization ON organization.id = member.organization_id").
		Where("member.organization_id = ? AND member.profile_id = ?", organizationId, userId).
		Where("member.status = 'active' AND member.removed_at IS NULL").
		Where("organization.status = 'active' AND organization.archived_at IS NULL").
		Take(&row)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "failed to get organization member: %v", result.Error)
	}
	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		return nil, nil
	}
	roleKey := uint32(530)
	switch row.RoleKey {
	case "owner":
		roleKey = 510
	case "admin", "billing":
		roleKey = 520
	}
	return &dto.OrganizationMember{
		OrganizationId: organizationId,
		UserId:         userId,
		RoleKey:        roleKey,
		Status:         10,
	}, nil
}
