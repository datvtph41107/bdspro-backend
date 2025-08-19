package handler

import (
	"auth/internal/dto"
	"auth/internal/usecase"
	"context"
	authpb "pb/types/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	AuthUsecase *usecase.AuthUsecase
}

func NewAuthHandler(service *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: service,
	}
}

// @Summary Gửi OTP
// @Description Gửi OTP để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param phone body authpb.RequestOTPRequest true "Số điện thoại"
// @Success 200 {object} authpb.RequestOTPResponse
// @Failure 400 {object} authpb.RequestOTPResponse
// @Failure 500 {object} authpb.RequestOTPResponse
// @Router /otp/request [post]
func (h *AuthHandler) RequestOTP(ctx context.Context, req *authpb.RequestOTPRequest) (*authpb.RequestOTPResponse, error) {
	// Convert proto request to DTO
	dto := &dto.OtpRequest{
		Phone:    req.Phone,
		Fullname: req.Fullname,
	}

	result, err := h.AuthUsecase.RequestOtp(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	response := &authpb.RequestOTPResponse{
		AuthId: result.AuthID,
		// Message: result.Message,
	}

	return response, nil
}

// @Summary Xác minh OTP
// @Description Xác minh OTP để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.VerifyOTPRequest true "Thông tin OTP"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /otp/verify [post]
func (h *AuthHandler) VerifyOTP(ctx context.Context, req *authpb.VerifyOTPRequest) (*authpb.AuthSuccessResponse, error) {
	// Convert proto request to DTO
	dto := &dto.OtpVerifyRequest{
		OtpRequest: dto.OtpRequest{
			Phone:    req.Phone,
			Fullname: req.Fullname,
		},
		Otp:      req.Otp,
		AuthID:   req.AuthId,
		Version:  req.Version,
		Platform: req.Platform,
		OS:       req.Os,
		Fullname: req.Fullname,
	}

	result, err := h.AuthUsecase.VerifyOtp(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	response := &authpb.AuthSuccessResponse{
		AuthId:       result.AuthID,
		ProfileId:    result.ProfileID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		FullName:     result.Username,
		// Avatar:         result.Avatar,
		// Phone:          result.Phone,
		// Email:          result.Email,
		OrganizationId: result.OrganizationID,
		Role:           result.Role,
	}

	return response, nil
}

// @Summary Gửi lại OTP
// @Description Gửi lại OTP để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.ResendOTPRequest true "Thông tin OTP"
// @Success 200 {object} authpb.ResendOTPResponse
// @Failure 400 {object} authpb.ResendOTPResponse
// @Failure 500 {object} authpb.ResendOTPResponse
// @Router /otp/resend [post]
func (h *AuthHandler) ResendOTP(ctx context.Context, req *authpb.ResendOTPRequest) (*authpb.ResendOTPResponse, error) {
	// Convert proto request to DTO
	dto := &dto.OtpRequest{
		Phone:    req.Phone,
		Fullname: req.Fullname,
	}

	result, err := h.AuthUsecase.ResendOTP(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi gửi lại OTP: %v", err)
	}

	// Convert DTO response to proto response
	response := &authpb.ResendOTPResponse{
		AuthId:  result.AuthId,
		Message: result.Message,
	}

	return response, nil
}

// @Summary Refresh token
// @Description Refresh token để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.RefreshTokenRequest true "Thông tin refresh token"
// @Success 200 {object} authpb.RefreshTokenResponse
// @Failure 400 {object} authpb.RefreshTokenResponse
// @Failure 500 {object} authpb.RefreshTokenResponse
// @Router /token/refresh [post]
func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	result, err := h.AuthUsecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi refresh token: %v", err)
	}

	return &authpb.RefreshTokenResponse{
		AccessToken: result.AccessToken,
	}, nil
}

// @Summary Logout
// @Description Logout để đăng xuất
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} authpb.LogoutResponse
// @Failure 400 {object} authpb.LogoutResponse
// @Failure 500 {object} authpb.LogoutResponse
// @Router /logout [delete]
func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	result, err := h.AuthUsecase.Logout(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi logout: %v", err)
	}

	return &authpb.LogoutResponse{
		Message: result.Message,
	}, nil
}

// @Summary Delete account
// @Description Xóa tài khoản
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} authpb.DeleteAccountResponse
// @Failure 400 {object} authpb.DeleteAccountResponse
// @Failure 500 {object} authpb.DeleteAccountResponse
// @Router /delete [delete]
func (h *AuthHandler) DeleteAccount(ctx context.Context, req *authpb.DeleteAccountRequest) (*authpb.DeleteAccountResponse, error) {
	result, err := h.AuthUsecase.Delete(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi xóa tài khoản: %v", err)
	}

	return &authpb.DeleteAccountResponse{
		Message: result.Message,
	}, nil
}

// @Summary Restore account
// @Description Khôi phục tài khoản
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.RestoreAccountRequest true "Thông tin khôi phục tài khoản"
// @Success 200 {object} authpb.RestoreAccountResponse
// @Failure 400 {object} authpb.RestoreAccountResponse
// @Failure 500 {object} authpb.RestoreAccountResponse
// @Router /restore [put]
func (h *AuthHandler) RestoreAccount(ctx context.Context, req *authpb.RestoreAccountRequest) (*authpb.RestoreAccountResponse, error) {
	// Convert proto request to DTO
	dto := &dto.RestorePhoneRequest{
		Phone: req.Phone,
	}

	result, err := h.AuthUsecase.Restore(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi khôi phục tài khoản: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.RestoreAccountResponse{
		AuthId:   result.AuthId,
		Provider: result.Provider,
		OauthId:  result.OAuthId,
	}, nil
}

// @Summary QR Init
// @Description Tạo QR code để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.QRInitRequest true "Thông tin QR"
// @Success 200 {object} authpb.QRInitResponse
// @Failure 400 {object} authpb.QRInitResponse
// @Failure 500 {object} authpb.QRInitResponse
// @Router /qr/init [post]
func (h *AuthHandler) QRInit(ctx context.Context, req *authpb.QRInitRequest) (*authpb.QRInitResponse, error) {
	// Convert proto request to DTO
	dto := &dto.SessionQRRequest{
		ClientID:   req.ClientId,
		UserAgent:  req.UserAgent,
		CsrfToken:  req.CsrfToken,
		Expiration: req.Expiration,
		SessionKey: req.SessionKey,
		SessionId:  req.SessionId,
	}

	result, err := h.AuthUsecase.CreateQRSession(ctx, dto, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo QR session: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.QRInitResponse{
		SessionKey: result.SessionKey,
		SessionId:  result.SessionId,
	}, nil
}

// @Summary QR Confirm
// @Description Xác nhận QR code để đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.QRConfirmRequest true "Thông tin QR"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /qr/confirm [post]
func (h *AuthHandler) QRConfirm(ctx context.Context, req *authpb.QRConfirmRequest) (*authpb.AuthSuccessResponse, error) {
	// Convert proto request to DTO
	dto := &dto.SessionQRConfirm{
		SessionId:  req.SessionId,
		SessionKey: req.SessionKey,
	}

	result, err := h.AuthUsecase.VerifyQRSession(ctx, dto, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi kiểm tra trạng thái QR: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:       result.AuthID,
		ProfileId:    result.ProfileID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

// @Summary Switch Organization
// @Description Chuyển đổi tổ chức
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.SwitchOrganizationRequest true "Thông tin chuyển đổi tổ chức"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /switch/organization [post]
func (h *AuthHandler) SwitchOrganization(ctx context.Context, req *authpb.SwitchOrganizationRequest) (*authpb.AuthSuccessResponse, error) {
	result, err := h.AuthUsecase.SwitchOrganization(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi chuyển đổi tổ chức: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:         result.AuthID,
		ProfileId:      result.ProfileID,
		OrganizationId: result.OrganizationID,
		Role:           result.Role,
		// FullName:       result.FullName,
		// Avatar:         result.Avatar,
		// Email:          result.Email,
		// Phone:          result.Phone,
		AccessToken: result.AccessToken,
	}, nil
}

// @Summary Admin Login
// @Description Đăng nhập admin bằng username và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.AdminLoginRequest true "Thông tin đăng nhập admin"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /admin/login [post]
func (h *AuthHandler) AdminLogin(ctx context.Context, req *authpb.AdminLoginRequest) (*authpb.AuthSuccessResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.AuthUsecase.AdminLogin(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:         result.AuthID,
		ProfileId:      result.ProfileID,
		OrganizationId: result.OrganizationID,
		Role:           result.Role,
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
	}, nil
}
