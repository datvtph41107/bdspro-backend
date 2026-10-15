package services

import (
	"context"

	_errors "common/errors"
	_utils "common/utils"
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

// AuthUsecaseLoginService adapts AuthUsecase into modules/auth LoginService.
type AuthUsecaseLoginService struct {
	auth *AuthUsecase
}

func NewAuthUsecaseLoginService(auth *AuthUsecase) *AuthUsecaseLoginService {
	return &AuthUsecaseLoginService{auth: auth}
}

func (s *AuthUsecaseLoginService) RequestOTP(ctx context.Context, req dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	return s.auth.RequestOtp(ctx, req)
}

func (s *AuthUsecaseLoginService) ResendOTP(ctx context.Context, req dto.OtpRequest) (*dto.ResendOTPResponse, error) {
	return s.auth.ResendOTP(ctx, req)
}

func (s *AuthUsecaseLoginService) VerifyOTP(ctx context.Context, req dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	return s.auth.VerifyOtp(ctx, req)
}

func (s *AuthUsecaseLoginService) LoginWithPassword(ctx context.Context, req dto.LoginWithPasswordRequest) (*dto.AuthLoginResponse, error) {
	return s.auth.LoginWithPassword(ctx, req)
}

func (s *AuthUsecaseLoginService) LoginWithKey(ctx context.Context, req dto.LoginWithKeyRequest) (*dto.AuthLoginResponse, error) {
	return s.auth.LoginWithKey(ctx, req)
}

func (s *AuthUsecaseLoginService) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AuthLoginResponse, error) {
	return s.auth.AdminLogin(ctx, req)
}

func (s *AuthUsecaseLoginService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	return s.auth.RefreshToken(ctx, refreshToken)
}

func (s *AuthUsecaseLoginService) PhoneCheck(ctx context.Context, encodedPhone string) (*dto.PhoneCheckResponse, error) {
	phone, err := _utils.DecodePhone(encodedPhone)
	if err != nil {
		return nil, _errors.ReturnError(400, "Mã số điện thoại không hợp lệ")
	}
	existed, userID, err := s.auth.PhoneCheck(ctx, phone)
	if err != nil {
		return nil, err
	}
	return &dto.PhoneCheckResponse{
		Existed: existed,
		UserId:  userID,
	}, nil
}
