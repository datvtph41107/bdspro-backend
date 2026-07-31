package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"pb/clients"
	organizationpb "pb/types/organization"
	"user/env"
	"user/internal/dto"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// @bind: user/internal/interface/providers.OrganizationProvider
type OrganizationProviderClient struct {
	*clients.OrganizationClient
	roleClient organizationpb.OrganizationRoleServiceClient
}

func NewOrganizationProviderClient() *OrganizationProviderClient {
	url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, "organization")
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s%s", url, viper.GetString("grpc.organization.host")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("Failed to connect to organization role service: %v", err)
	}

	return &OrganizationProviderClient{
		OrganizationClient: clients.MakeOrganizationClient(env.ENV_RUNTIME),
		roleClient:         organizationpb.NewOrganizationRoleServiceClient(conn),
	}
}

func (c *OrganizationProviderClient) SwitchOrganization(ctx context.Context, organizationID uint64) (*dto.SwitchOrganizationResponse, error) {
	return &dto.SwitchOrganizationResponse{
		OrganizationId: organizationID,
	}, nil
}

func (c *OrganizationProviderClient) GetAdminRoles(ctx context.Context) ([]dto.AdminRoleDTO, error) {
	page := uint32(1)
	size := uint32(100)
	resp, err := c.roleClient.GetRoles(ctx, &organizationpb.GetRolesRequest{
		Page: &page,
		Size: &size,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get admin roles: %v", err)
	}

	adminRoles := make([]dto.AdminRoleDTO, 0)
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
		RoleKey:        resp.RoleKey,
	}, nil
}
