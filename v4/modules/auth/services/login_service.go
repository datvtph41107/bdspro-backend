package services

import (
	"context"

	_errors "common/errors"
	_utils "common/utils"
	legacydto "userv4/internal/dto"
	"userv4/internal/usecase"
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

// AuthUsecaseLoginService adapts the legacy AuthUsecase into modules/auth LoginService.
// Login ownership lives here; AuthUsecase remains the implementation backend until
// bodies are moved into this package.
type AuthUsecaseLoginService struct {
	auth *usecase.AuthUsecase
}

func NewAuthUsecaseLoginService(auth *usecase.AuthUsecase) *AuthUsecaseLoginService {
	return &AuthUsecaseLoginService{auth: auth}
}

func (s *AuthUsecaseLoginService) RequestOTP(ctx context.Context, req dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	result, err := s.auth.RequestOtp(ctx, toLegacyOtpRequest(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyRequestLogin(result), nil
}

func (s *AuthUsecaseLoginService) ResendOTP(ctx context.Context, req dto.OtpRequest) (*dto.ResendOTPResponse, error) {
	result, err := s.auth.ResendOTP(ctx, toLegacyOtpRequest(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyResendOTP(result), nil
}

func (s *AuthUsecaseLoginService) VerifyOTP(ctx context.Context, req dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	result, err := s.auth.VerifyOtp(ctx, toLegacyOtpVerify(req))
	if err != nil {
		return nil, err
	}
	return fromLegacyAuthLogin(result), nil
}

func (s *AuthUsecaseLoginService) LoginWithPassword(ctx context.Context, req dto.LoginWithPasswordRequest) (*dto.AuthLoginResponse, error) {
	result, err := s.auth.LoginWithPassword(ctx, legacydto.LoginWithPasswordRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.OS,
		DeviceName: req.DeviceName,
		DeviceID:   req.DeviceID,
	})
	if err != nil {
		return nil, err
	}
	return fromLegacyAuthLogin(result), nil
}

func (s *AuthUsecaseLoginService) LoginWithKey(ctx context.Context, req dto.LoginWithKeyRequest) (*dto.AuthLoginResponse, error) {
	result, err := s.auth.LoginWithKey(ctx, legacydto.LoginWithKeyRequest{
		Phone:            req.Phone,
		EncryptedAuthKey: req.EncryptedAuthKey,
		ClientPublicKey:  req.ClientPublicKey,
		Version:          req.Version,
		Platform:         req.Platform,
		OS:               req.OS,
		DeviceName:       req.DeviceName,
		DeviceID:         req.DeviceID,
	})
	if err != nil {
		return nil, err
	}
	return fromLegacyAuthLogin(result), nil
}

func (s *AuthUsecaseLoginService) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AuthLoginResponse, error) {
	result, err := s.auth.AdminLogin(ctx, legacydto.AdminLoginRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.OS,
		DeviceName: req.DeviceName,
		DeviceID:   req.DeviceID,
	})
	if err != nil {
		return nil, err
	}
	return fromLegacyAuthLogin(result), nil
}

func (s *AuthUsecaseLoginService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	result, err := s.auth.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &dto.RefreshResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
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

func toLegacyOtpRequest(req dto.OtpRequest) legacydto.OtpRequest {
	return legacydto.OtpRequest{
		Phone:    req.Phone,
		Fullname: req.Fullname,
		Mode:     req.Mode,
		AuthId:   req.AuthId,
	}
}

func toLegacyOtpVerify(req dto.OtpVerifyRequest) legacydto.OtpVerifyRequest {
	return legacydto.OtpVerifyRequest{
		OtpRequest:      toLegacyOtpRequest(req.OtpRequest),
		Otp:             req.Otp,
		AuthID:          req.AuthID,
		Version:         req.Version,
		Platform:        req.Platform,
		OS:              req.OS,
		DeviceName:      req.DeviceName,
		DeviceID:        req.DeviceID,
		Fullname:        req.Fullname,
		ClientPublicKey: req.ClientPublicKey,
	}
}

func fromLegacyRequestLogin(result *legacydto.RequestLoginResponse) *dto.RequestLoginResponse {
	if result == nil {
		return nil
	}
	return &dto.RequestLoginResponse{
		AuthID:     result.AuthID,
		Provider:   result.Provider,
		OAuthID:    result.OAuthID,
		Existed:    result.Existed,
		SmsChannel: result.SmsChannel,
	}
}

func fromLegacyResendOTP(result *legacydto.ResendOTPResponse) *dto.ResendOTPResponse {
	if result == nil {
		return nil
	}
	return &dto.ResendOTPResponse{
		AuthId:     result.AuthId,
		Message:    result.Message,
		SmsChannel: result.SmsChannel,
	}
}

func fromLegacyAuthLogin(result *legacydto.AuthLoginResponse) *dto.AuthLoginResponse {
	if result == nil {
		return nil
	}
	return &dto.AuthLoginResponse{
		AccessToken:      result.AccessToken,
		AuthID:           result.AuthID,
		Role:             result.Role,
		Username:         result.Username,
		ProfileID:        result.ProfileID,
		RefreshToken:     result.RefreshToken,
		OrganizationID:   result.OrganizationID,
		FullName:         result.FullName,
		PasswordAuthId:   result.PasswordAuthId,
		ReferralCode:     result.ReferralCode,
		RequireFullname:  result.RequireFullname,
		TempToken:        result.TempToken,
		EncryptedAuthKey: result.EncryptedAuthKey,
		ServerPublicKey:  result.ServerPublicKey,
	}
}
