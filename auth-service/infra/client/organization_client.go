package client

import (
	"auth/internal/dto"
	"context"
	organizationpb "pb/types/organization"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @bind: auth/internal/interface/provider.OrganizationProvider
type OrganizationClient struct {
	client organizationpb.OrganizationRoleServiceClient
}

func NewOrganizationClient(conn *grpc.ClientConn) *OrganizationClient {
	return &OrganizationClient{
		client: organizationpb.NewOrganizationRoleServiceClient(conn),
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
