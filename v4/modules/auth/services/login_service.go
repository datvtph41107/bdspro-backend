package services

import (
	"context"

	"userv4/modules/auth/dto"
)

// LoginService owns the authentication entry flows migrated into modules/auth.
type LoginService interface {
	RequestOTP(ctx context.Context, req dto.OtpRequest) (*dto.RequestLoginResponse, error)
	ResendOTP(ctx context.Context, req dto.OtpRequest) (*dto.ResendOTPResponse, error)
	VerifyOTP(ctx context.Context, req dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error)
	LoginWithPassword(ctx context.Context, req dto.LoginWithPasswordRequest) (*dto.AuthLoginResponse, error)
	LoginWithKey(ctx context.Context, req dto.LoginWithKeyRequest) (*dto.AuthLoginResponse, error)
	AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AuthLoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	PhoneCheck(ctx context.Context, encodedPhone string) (*dto.PhoneCheckResponse, error)
}
