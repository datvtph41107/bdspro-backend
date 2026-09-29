package client

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"errors"
	"pb/clients"
	"strings"
	domainauth "user/internal/domain/auth"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"

	"gorm.io/gorm"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	userClient     *clients.UserGrpcClient
	authRepository repo.AuthMethodRepository
}

func NewAuthClient(
	authClient *clients.AuthGrpcClient,
	userClient *clients.UserGrpcClient,
	authRepository repo.AuthMethodRepository,
) providers.AuthProvider {
	return &AuthClient{
		AuthGrpcClient: authClient,
		userClient:     userClient,
		authRepository: authRepository,
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
	if authMethod == nil {
		return nil, errors.New("auth method is nil")
	}
	provider := domainauth.ProviderType(strings.ToUpper(strings.TrimSpace(authMethod.Provider)))
	if provider == "" {
		provider = domainauth.ProviderAdmin
	}
	password := authMethod.Password
	if password != "" {
		var err error
		password, err = _utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
	}
	created, err := c.authRepository.Create(ctx, &domainauth.AuthMethod{
		Provider: provider,
		AuthName: strings.TrimSpace(authMethod.AuthName),
		Password: password,
		FullName: authMethod.FullName,
		Email:    strings.TrimSpace(authMethod.Email),
		Phone:    strings.TrimSpace(authMethod.Phone),
		Avatar:   authMethod.Avatar,
		RoleKey:  authMethod.RoleKey,
		Status:   uint8(domainauth.StatusActive),
		UserID:   authMethod.UserID,
	})
	if err != nil {
		return nil, err
	}
	return authMethodToProvider(created), nil
}

func (c *AuthClient) Update(ctx context.Context, authMethod *providers.AuthMethodDomain) (*providers.AuthMethodDomain, error) {
	if authMethod == nil {
		return nil, errors.New("auth method is nil")
	}

	existing, err := c.findAuthForUpdate(ctx, authMethod)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) && strings.TrimSpace(authMethod.AuthName) != "" {
			return c.Create(ctx, authMethod)
		}
		return nil, err
	}
	if authMethod.UserID > 0 {
		existing.UserID = authMethod.UserID
	}
	if value := strings.TrimSpace(authMethod.AuthName); value != "" {
		existing.AuthName = value
	}
	if value := strings.TrimSpace(authMethod.Provider); value != "" {
		existing.Provider = domainauth.ProviderType(strings.ToUpper(value))
	}
	if authMethod.Password != "" {
		existing.Password, err = _utils.HashPassword(authMethod.Password)
		if err != nil {
			return nil, err
		}
	}
	if authMethod.FullName != "" {
		existing.FullName = authMethod.FullName
	}
	if authMethod.Email != "" {
		existing.Email = strings.TrimSpace(authMethod.Email)
	}
	if authMethod.Phone != "" {
		existing.Phone = strings.TrimSpace(authMethod.Phone)
	}
	if authMethod.Avatar != "" {
		existing.Avatar = authMethod.Avatar
	}
	if authMethod.RoleKey > 0 {
		existing.RoleKey = authMethod.RoleKey
	}
	updated, err := c.authRepository.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return authMethodToProvider(updated), nil
}

func (c *AuthClient) SoftDelete(ctx context.Context, id uint64) error {
	return c.authRepository.SoftDelete(ctx, id)
}

func (c *AuthClient) AssignRoleToUser(ctx context.Context, profileID, roleID uint64) error {
	if profileID == 0 || roleID == 0 {
		return nil
	}
	return c.AuthGrpcClient.AssignRoleToUser(ctx, profileID, roleID)
}

func (c *AuthClient) DeleteAllAuthMethodsByUserId(ctx context.Context, userId uint64) error {
	return c.authRepository.DeleteByUserId(ctx, userId)
}

func (c *AuthClient) FindAdminsByProvider(ctx context.Context, provider string, page, size int) ([]*providers.AuthMethodDomain, int64, error) {
	items, total, err := c.authRepository.FindAdminsByProvider(ctx, provider, page, size, nil, "")
	if err != nil {
		return nil, 0, err
	}
	result := make([]*providers.AuthMethodDomain, 0, len(items))
	for _, item := range items {
		result = append(result, authMethodToProvider(item))
	}
	return result, total, nil
}

func (c *AuthClient) FindByAuthNameAndProvider(ctx context.Context, authName, provider string) (*providers.AuthMethodDomain, error) {
	item, err := c.authRepository.FindByAuthNameAndProvider(ctx, authName, provider)
	return optionalProviderAuth(item, err)
}

func (c *AuthClient) FindByEmail(ctx context.Context, email string) (*providers.AuthMethodDomain, error) {
	item, err := c.authRepository.FindByEmail(ctx, strings.TrimSpace(email))
	return optionalProviderAuth(item, err)
}

func (c *AuthClient) FindByID(ctx context.Context, id uint64) (*providers.AuthMethodDomain, error) {
	item, err := c.authRepository.FindByID(ctx, id)
	return optionalProviderAuth(item, err)
}

func (c *AuthClient) findAuthForUpdate(ctx context.Context, input *providers.AuthMethodDomain) (*domainauth.AuthMethod, error) {
	if input.ID > 0 {
		return c.authRepository.FindByID(ctx, input.ID)
	}
	if input.UserID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if provider := strings.TrimSpace(input.Provider); provider != "" {
		return c.authRepository.GetByUserIdAndProvider(ctx, input.UserID, strings.ToUpper(provider))
	}
	return c.authRepository.GetFirstByUserID(ctx, input.UserID)
}

func optionalProviderAuth(item *domainauth.AuthMethod, err error) (*providers.AuthMethodDomain, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return authMethodToProvider(item), nil
}

func authMethodToProvider(item *domainauth.AuthMethod) *providers.AuthMethodDomain {
	if item == nil {
		return nil
	}
	return &providers.AuthMethodDomain{
		ID:       item.ID,
		UserID:   item.UserID,
		Provider: string(item.Provider),
		AuthName: item.AuthName,
		// Password is write-only at the AuthProvider boundary. Returning even a
		// hash here both leaks credential material into unrelated use cases and
		// makes a later metadata-only Update hash the hash a second time.
		Password:  "",
		FullName:  item.FullName,
		Email:     item.Email,
		Phone:     item.Phone,
		Avatar:    item.Avatar,
		RoleKey:   item.RoleKey,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		DeletedAt: item.DeletedAt,
	}
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
