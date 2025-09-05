package client

import (
	"auth/env"
	"auth/internal/dto"
	"context"
	clients "pb/clients"
	organizationpb "pb/types/organization"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @bind: auth/internal/interface/provider.OrganizationProvider
type OrganizationClient struct {
	*clients.OrganizationClient
	client organizationpb.OrganizationRoleServiceClient
}

func NewOrganizationClient(conn *grpc.ClientConn) *OrganizationClient {
	return &OrganizationClient{
		OrganizationClient: clients.MakeOrganizationClient(env.ENV_RELEASE),
		client:             organizationpb.NewOrganizationRoleServiceClient(conn),
	}
}

// SwitchOrganization chuyển đổi organization
func (c *OrganizationClient) SwitchOrganization(ctx context.Context, organizationID uint64) (*dto.SwitchOrganizationResponse, error) {
	// Implementation sẽ được thêm sau
	return &dto.SwitchOrganizationResponse{
		OrganizationId: organizationID,
	}, nil
}

// GetAdminRoles lấy danh sách role admin từ organization service
func (c *OrganizationClient) GetAdminRoles(ctx context.Context) ([]dto.AdminRoleDTO, error) {
	// Gọi API lấy danh sách role từ organization service
	page := uint32(1)
	size := uint32(100)
	resp, err := c.client.GetRoles(ctx, &organizationpb.GetRolesRequest{
		Page: &page,
		Size: &size, // Lấy tất cả role
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get admin roles: %v", err)
	}

	// Lọc chỉ lấy role có key là "admin"
	var adminRoles []dto.AdminRoleDTO
	for _, role := range resp.Data {
		if role.Key == "admin" {
			adminRoles = append(adminRoles, dto.AdminRoleDTO{
				ID:   role.Id,
				Name: role.Name,
				Key:  role.Key,
			})
		}
	}

	return adminRoles, nil
}

func (c *OrganizationClient) GetGroupRoleUser(ctx context.Context, groupId uint64, userId uint64) (*dto.GroupMember, error) {
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

func (c *OrganizationClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error) {
	resp, err := c.OrganizationClient.GetDealMember(ctx, dealId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get deal member role: %v", err)
	}
	if resp == nil {
		return nil, nil
	}

	return &dto.DealMember{
		DealId: resp.DealId,
		UserId: resp.MemberId,
		RoleId: resp.RoleId,
	}, nil
}

func (c *OrganizationClient) GetBranchMember(ctx context.Context, branchId uint64, userId uint64) (*dto.BranchMember, error) {
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

func (c *OrganizationClient) GetOrganizationMember(ctx context.Context, organizationId uint64, userId uint64) (*dto.OrganizationMember, error) {
	resp, err := c.OrganizationClient.GetOrganizationMember(ctx, organizationId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get organization member role: %v", err)
	}
	if resp == nil {
		return nil, nil
	}

	return &dto.OrganizationMember{
		OrganizationId: organizationId,
		UserId:         userId,
		RoleId:         &resp.RoleId,
	}, nil
}
