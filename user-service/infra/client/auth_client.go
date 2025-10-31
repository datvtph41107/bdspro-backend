package client

import (
	"context"
	// authpb "pb/types/auth" // TODO: Uncomment sau khi gen proto
	"pb/clients"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"user/env"
	"user/internal/interface/providers"
	"user/internal/models"
)

// @bind: user/internal/interface/providers.AuthProvider
type AuthClient struct {
	*clients.AuthGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) GetUserStatus(ctx context.Context, profileID uint64) (*providers.UserStatusInfo, error) {
	// TODO: Implement gọi auth-service để lấy thông tin user status
	// Hiện tại trả về mock data, sau này sẽ implement thực tế

	// Mock data - sẽ được thay thế bằng gọi API thực tế
	return &providers.UserStatusInfo{
		ProfileID:   profileID,
		Verified:    true,
		VerifiedAt:  nil, // Sẽ lấy từ auth-service
		LockedUntil: nil, // Sẽ lấy từ auth-service
		Active:      true,
	}, nil
}

func (c *AuthClient) GetUserStatusBatch(ctx context.Context, profileIDs []uint64) (map[uint64]*providers.UserStatusInfo, error) {
	// TODO: Implement gọi auth-service để lấy thông tin user status cho nhiều user
	// Hiện tại trả về mock data, sau này sẽ implement thực tế

	result := make(map[uint64]*providers.UserStatusInfo)
	for _, profileID := range profileIDs {
		result[profileID] = &providers.UserStatusInfo{
			ProfileID:   profileID,
			Verified:    true,
			VerifiedAt:  nil,
			LockedUntil: nil,
			Active:      true,
		}
	}

	return result, nil
}

func (c *AuthClient) Create(ctx context.Context, authMethod *providers.AuthMethodDomain) (*providers.AuthMethodDomain, error) {
	// return c.AuthGrpcClient.Create(ctx, authMethod)
	return authMethod, nil
}

func (c *AuthClient) Update(ctx context.Context, authMethod *providers.AuthMethodDomain) (*providers.AuthMethodDomain, error) {
	req := &authpb.AuthMethod{
		Id:       authMethod.ID,
		AuthName: authMethod.AuthName,
		FullName: authMethod.FullName,
		Email:    authMethod.Email,
		Phone:    authMethod.Phone,
		Avatar:   authMethod.Avatar,
		Password: authMethod.Password,
		UserId:   authMethod.UserID,
	}
	resp, err := c.AuthGrpcClient.UpdateAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return &providers.AuthMethodDomain{
		ID:       resp.Id,
		AuthName: resp.AuthName,
		FullName: resp.FullName,
		Email:    resp.Email,
		Phone:    resp.Phone,
		Avatar:   resp.Avatar,
	}, nil
}

func (c *AuthClient) SoftDelete(ctx context.Context, id uint64) error {
	// return c.AuthGrpcClient.SoftDelete(ctx, id)
	return nil
}

// DeleteAllAuthMethodsByUserId xóa tất cả auth_method của 1 user
func (c *AuthClient) DeleteAllAuthMethodsByUserId(ctx context.Context, userId uint64) error {
	// _, err := c.AuthGrpcClient.DeleteAuthMethodsByUserId(ctx, &sharepb.IdRequest{
	// 	Id: userId,
	// })
	// return err
	return nil
}

func (c *AuthClient) FindAdminsByProvider(ctx context.Context, provider string, page, size int) ([]*providers.AuthMethodDomain, int64, error) {
	// return c.AuthGrpcClient.FindAdminsByProvider(ctx, provider, page, size)
	return nil, 0, nil
}

func (c *AuthClient) FindByAuthNameAndProvider(ctx context.Context, authName, provider string) (*providers.AuthMethodDomain, error) {
	// return c.AuthGrpcClient.FindByAuthNameAndProvider(ctx, authName, provider)
	return nil, nil
}

func (c *AuthClient) FindByEmail(ctx context.Context, email string) (*providers.AuthMethodDomain, error) {
	// return c.AuthGrpcClient.FindByEmail(ctx, email)
	return nil, nil
}

func (c *AuthClient) FindByID(ctx context.Context, id uint64) (*providers.AuthMethodDomain, error) {
	// Sử dụng GetAuthAdminByIds với 1 ID
	resp, err := c.AuthGrpcClient.GetAuthAdminByIds(ctx, &authpb.GetAuthAdminByIdsRequest{
		Ids: []uint64{id},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil
	}

	admin := resp.Data[0]
	return &providers.AuthMethodDomain{
		ID:       admin.Id,
		UserID:   admin.UserId,
		AuthName: admin.Username,
		FullName: admin.FullName,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Avatar:   admin.Avatar,
		RoleKey:  admin.RoleKey,
	}, nil
}

// AuthAdminInfo chứa thông tin admin từ auth service
type AuthAdminInfo struct {
	ProfileID uint64
	Username  string
	FullName  string
	Email     string
	Phone     string
	RoleKey   uint32
}

// GetAuthAdminByIds lấy thông tin admin theo danh sách IDs
func (c *AuthClient) GetAuthAdminByIds(ctx context.Context, ids []uint64) ([]*AuthAdminInfo, error) {
	resp, err := c.AuthGrpcClient.GetAuthAdminByIds(ctx, &authpb.GetAuthAdminByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*AuthAdminInfo, 0, len(resp.Data))
	for _, admin := range resp.Data {
		result = append(result, &AuthAdminInfo{
			ProfileID: admin.Id,
			Username:  admin.Username,
			FullName:  admin.FullName,
			Email:     admin.Email,
			Phone:     admin.Phone,
			RoleKey:   admin.RoleKey,
		})
	}

	return result, nil
}

func (c *AuthClient) GetRolesByIds(ctx context.Context, roleIds []uint64) (map[uint64]*models.RoleDTO, error) {
	if len(roleIds) == 0 {
		return make(map[uint64]*models.RoleDTO), nil
	}

	// TODO: Implement sau khi gen proto
	// Hiện tại trả về empty map
	roleMap := make(map[uint64]*models.RoleDTO)

	// NOTE: Code này sẽ được uncomment sau khi gen proto:
	rolesResp, err := c.AuthGrpcClient.GetRolesByIds(ctx, &sharepb.IdRequest{
		Ids: roleIds,
	})
	if err != nil {
		return nil, err
	}
	for _, role := range rolesResp.Data {
		roleMap[role.Id] = &models.RoleDTO{
			ID:              role.Id,
			RoleName:        role.RoleName,
			RoleKey:         role.RoleKey,
			RoleDescription: role.RoleDescription,
		}
	}

	return roleMap, nil
}
