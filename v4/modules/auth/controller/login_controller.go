package controller

import (
	"context"
	"log"
	"strings"

	_utils "common/utils"
	authpb "pb/types/auth"

	"userv4/modules/auth/dto"
	"userv4/modules/auth/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginController is the modules/auth entry for login RPCs.
// Legacy AuthHandler delegates login methods here.
type LoginController struct {
	login services.LoginService
}

func NewLoginController(login services.LoginService) *LoginController {
	return &LoginController{login: login}
}

func (c *LoginController) PhoneCheckEncoded(ctx context.Context, encodedPhone string) (*dto.PhoneCheckResponse, error) {
	encodedPhone = strings.TrimSpace(encodedPhone)
	encodedPhone = strings.TrimRight(encodedPhone, "|")
	return c.login.PhoneCheck(ctx, encodedPhone)
}

func (c *LoginController) RequestOTP(ctx context.Context, req *authpb.RequestOTPRequest) (*authpb.RequestOTPResponse, error) {
	result, err := c.login.RequestOTP(ctx, dto.OtpRequest{
		Phone:    req.Phone,
		Fullname: req.Fullname,
		Mode:     req.Mode,
		AuthId:   req.AuthId,
	})
	if err != nil {
		return nil, err
	}
	return &authpb.RequestOTPResponse{
		AuthId:     result.AuthID,
		Existed:    result.Existed,
		SmsChannel: result.SmsChannel,
	}, nil
}

func (c *LoginController) VerifyOTP(ctx context.Context, req *authpb.VerifyOTPRequest) (*authpb.AuthSuccessResponse, error) {
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	deviceID := _utils.GetDeviceIdFromContext(ctx)
	if deviceID == "" && req.DeviceName != "" {
		deviceID = req.DeviceName
	}

	result, err := c.login.VerifyOTP(ctx, dto.OtpVerifyRequest{
		OtpRequest: dto.OtpRequest{
			Phone:    req.Phone,
			Fullname: req.Fullname,
		},
		Fullname:        req.Fullname,
		Otp:             req.Otp,
		AuthID:          req.AuthId,
		Version:         req.Version,
		Platform:        req.Platform,
		OS:              req.Os,
		DeviceName:      req.DeviceName,
		DeviceID:        deviceID,
		ClientPublicKey: req.ClientPublicKey,
	})
	if err != nil {
		log.Printf("[VerifyOTP] phone=%s authId=%d failed: %v", req.Phone, req.AuthId, err)
		return nil, err
	}

	return toAuthSuccess(result), nil
}

func (c *LoginController) ResendOTP(ctx context.Context, req *authpb.ResendOTPRequest) (*authpb.ResendOTPResponse, error) {
	result, err := c.login.ResendOTP(ctx, dto.OtpRequest{
		Phone:    req.Phone,
		Fullname: req.Fullname,
		AuthId:   req.AuthId,
	})
	if err != nil {
		return nil, err
	}
	return &authpb.ResendOTPResponse{
		AuthId:  result.AuthId,
		Message: result.Message,
	}, nil
}

func (c *LoginController) LoginWithPassword(ctx context.Context, req *authpb.LoginWithPasswordRequest) (*authpb.AuthSuccessResponse, error) {
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	result, err := c.login.LoginWithPassword(ctx, dto.LoginWithPasswordRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.Os,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	return toAuthSuccess(result), nil
}

func (c *LoginController) LoginWithKey(ctx context.Context, req *authpb.LoginWithKeyRequest) (*authpb.AuthSuccessResponse, error) {
	if req.Phone == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Số điện thoại là bắt buộc")
	}
	if req.EncryptedAuthKey == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Encrypted auth key là bắt buộc")
	}
	if req.ClientPublicKey == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Client public key là bắt buộc")
	}
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	deviceID := _utils.GetDeviceIdFromContext(ctx)
	result, err := c.login.LoginWithKey(ctx, dto.LoginWithKeyRequest{
		Phone:            req.Phone,
		EncryptedAuthKey: req.EncryptedAuthKey,
		ClientPublicKey:  req.ClientPublicKey,
		Version:          req.Version,
		Platform:         req.Platform,
		OS:               req.Os,
		DeviceName:       req.DeviceName,
		DeviceID:         deviceID,
	})
	if err != nil {
		return nil, err
	}
	return toAuthSuccess(result), nil
}

func (c *LoginController) AdminLogin(ctx context.Context, req *authpb.AdminLoginRequest) (*authpb.AuthSuccessResponse, error) {
	result, err := c.login.AdminLogin(ctx, dto.AdminLoginRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.Os,
		DeviceName: req.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	return toAuthSuccess(result), nil
}

func (c *LoginController) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	result, err := c.login.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &authpb.RefreshTokenResponse{}, nil
	}
	return &authpb.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func toAuthSuccess(result *dto.AuthLoginResponse) *authpb.AuthSuccessResponse {
	if result == nil {
		return nil
	}
	return &authpb.AuthSuccessResponse{
		AuthId:           result.AuthID,
		ProfileId:        result.ProfileID,
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		FullName:         firstNonEmpty(result.FullName, result.Username),
		OrganizationId:   result.OrganizationID,
		Role:             result.Role,
		PasswordAuthId:   result.PasswordAuthId,
		ReferralCode:     result.ReferralCode,
		RequireFullname:  result.RequireFullname,
		TempToken:        result.TempToken,
		EncryptedAuthKey: result.EncryptedAuthKey,
		ServerPublicKey:  result.ServerPublicKey,
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
