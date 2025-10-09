package mapper

import (
	"auth/internal/domain"
	authpb "pb/types/auth"
)

type AuthMethodMapper struct {
}

func NewAuthMethodMapper() *AuthMethodMapper {
	return &AuthMethodMapper{}
}

func (m *AuthMethodMapper) MapAuthMethodPb(authMethod *domain.AuthMethod) *authpb.AuthMethod {
	return &authpb.AuthMethod{
		AuthName: authMethod.AuthName,
		Provider: string(authMethod.Provider),
		Email:    authMethod.Email,
		Phone:    authMethod.Phone,
		Avatar:   authMethod.Avatar,
		FullName: authMethod.FullName,
		UserId:   authMethod.UserID,
		Password: authMethod.Password,
	}
}

func (m *AuthMethodMapper) MapAuthMethodPbToDomain(authMethod *authpb.AuthMethod) *domain.AuthMethod {
	return &domain.AuthMethod{
		AuthName: authMethod.AuthName,
		Provider: domain.ProviderType(authMethod.Provider),
		Email:    authMethod.Email,
		Phone:    authMethod.Phone,
		Avatar:   authMethod.Avatar,
		FullName: authMethod.FullName,
		UserID:   authMethod.UserId,
		Password: authMethod.Password,
	}
}
