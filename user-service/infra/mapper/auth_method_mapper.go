package mapper

import (
	authpb "pb/types/auth"
	"user/internal/domain/auth"
)

type AuthMethodMapper struct {
}

func NewAuthMethodMapper() *AuthMethodMapper {
	return &AuthMethodMapper{}
}

func (m *AuthMethodMapper) MapAuthMethodPb(authMethod *auth.AuthMethod) *authpb.AuthMethod {
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

func (m *AuthMethodMapper) MapAuthMethodPbToDomain(authMethod *authpb.AuthMethod) *auth.AuthMethod {
	return &auth.AuthMethod{
		AuthName: authMethod.AuthName,
		Provider: auth.ProviderType(authMethod.Provider),
		Email:    authMethod.Email,
		Phone:    authMethod.Phone,
		Avatar:   authMethod.Avatar,
		FullName: authMethod.FullName,
		UserID:   authMethod.UserId,
		Password: authMethod.Password,
	}
}
