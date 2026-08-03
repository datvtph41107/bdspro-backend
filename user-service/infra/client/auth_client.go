package client

import (
	_dto "common/domain/dto"
	"context"
	"pb/clients"
	"user/env"
	"user/internal/interface/providers"
	"user/internal/models"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	userClient *clients.UserGrpcClient
}

func NewAuthClient() providers.AuthProvider {
	return &AuthClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
		userClient:     clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) GetUserStatus(ctx context.Context, profileID uint64) (*providers.UserStatusInfo, error) {
	return &providers.UserStatusInfo{
		ProfileID:   profileID,
		Verified:    true,
		VerifiedAt:  nil,
		LockedUntil: nil,
		Active:      true,
	}, nil
}

func (c *AuthClient) GetUserStatusBatch(ctx context.Context, profileIDs []uint64) (map[uint64]*providers.UserStatusInfo, error) {
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
	return authMethod, nil
}

func (c *AuthClient) Update(ctx context.Context, authMethod *providers.AuthMethodDomain) (*providers.AuthMethodDomain, error) {
	return authMethod, nil
}

func (c *AuthClient) SoftDelete(ctx context.Context, id uint64) error {
	return nil
}

func (c *AuthClient) AssignRoleToUser(ctx context.Context, profileID, roleID uint64) error {
	if profileID == 0 || roleID == 0 {
		return nil
	}
	return c.AuthGrpcClient.AssignRoleToUser(ctx, profileID, roleID)
}

func (c *AuthClient) DeleteAllAuthMethodsByUserId(ctx context.Context, userId uint64) error {
	return nil
}

func (c *AuthClient) FindAdminsByProvider(ctx context.Context, provider string, page, size int) ([]*providers.AuthMethodDomain, int64, error) {
	return nil, 0, nil
}

func (c *AuthClient) FindByAuthNameAndProvider(ctx context.Context, authName, provider string) (*providers.AuthMethodDomain, error) {
	return nil, nil
}

func (c *AuthClient) FindByEmail(ctx context.Context, email string) (*providers.AuthMethodDomain, error) {
	return nil, nil
}

func (c *AuthClient) FindByID(ctx context.Context, id uint64) (*providers.AuthMethodDomain, error) {
	resp, err := c.userClient.GetAuthAdminByIds(ctx, []uint64{id})
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

type AuthAdminInfo struct {
	ID        uint64
	UserID    uint64
	ProfileID uint64
	Username  string
	FullName  string
	Email     string
	Phone     string
	RoleKey   uint32
}

func (c *AuthClient) GetAuthAdminByIds(ctx context.Context, ids []uint64) ([]*AuthAdminInfo, error) {
	resp, err := c.userClient.GetAuthAdminByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]*AuthAdminInfo, 0, len(resp.Data))
	for _, admin := range resp.Data {
		result = append(result, &AuthAdminInfo{
			ID:        admin.Id,
			UserID:    admin.UserId,
			ProfileID: admin.UserId,
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

	roleMap := make(map[uint64]*models.RoleDTO)

	rolesResp, err := c.userClient.GetRolesByIds(ctx, roleIds)
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

func (c *AuthClient) GetAuthUsersByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*_dto.AuthUserProfileDTO, error) {
	return c.userClient.GetAuthUsersByProfileIDs(ctx, profileIDs)
}
