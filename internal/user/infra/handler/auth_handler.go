package handler

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"strings"
	"time"
	"user/internal"

	"user/infra/mapper"
	"user/internal/dto"
	"user/internal/usecase"
	"user/internal/usecase/useradmin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	AuthUsecase         *usecase.AuthUsecase
	AdminAccessUsecase  *usecase.AdminAccessUsecase
	InternalHandler     *InternalHandler
	AdminUsecase        *usecase.AuthAdminUsecase
	SessionMapper       *mapper.SessionMapper
	PINUsecase          *usecase.PINUsecase
	DeviceUsecase       *usecase.DeviceUsecase
	DeviceMapper        *mapper.DeviceMapper
	AuthSecurityUsecase *usecase.AuthSecurityUsecase
	AuthSecurityMapper  *mapper.AuthSecurityMapper
	UserAdminAuthorizer useradmin.PermissionAuthorizer
}

func NewAuthHandler(service *usecase.AuthUsecase,
	adminAccessUsecase *usecase.AdminAccessUsecase,
	internalHandler *InternalHandler,
	adminUsecase *usecase.AuthAdminUsecase,
	sessionMapper *mapper.SessionMapper,
	pinUsecase *usecase.PINUsecase,
	deviceUsecase *usecase.DeviceUsecase,
	deviceMapper *mapper.DeviceMapper,
	authSecurityUsecase *usecase.AuthSecurityUsecase,
	authSecurityMapper *mapper.AuthSecurityMapper,
	userAdminAuthorizer useradmin.PermissionAuthorizer,
) *AuthHandler {
	return &AuthHandler{
		AuthUsecase:         service,
		AdminAccessUsecase:  adminAccessUsecase,
		InternalHandler:     internalHandler,
		AdminUsecase:        adminUsecase,
		SessionMapper:       sessionMapper,
		PINUsecase:          pinUsecase,
		DeviceUsecase:       deviceUsecase,
		DeviceMapper:        deviceMapper,
		AuthSecurityUsecase: authSecurityUsecase,
		AuthSecurityMapper:  authSecurityMapper,
		UserAdminAuthorizer: userAdminAuthorizer,
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
		Mode:     req.Mode,
		AuthId:   req.AuthId,
	}

	result, err := h.AuthUsecase.RequestOtp(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	response := &authpb.RequestOTPResponse{
		AuthId:     result.AuthID,
		Existed:    result.Existed,
		SmsChannel: result.SmsChannel,
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
	// Validate required fields for session creation
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	// Convert proto request to DTO
	// Lấy DeviceID từ context hoặc dùng deviceName làm fallback
	deviceID := _utils.GetDeviceIdFromContext(ctx)
	if deviceID == "" && req.DeviceName != "" {
		deviceID = req.DeviceName
	}

	otpDto := &dto.OtpVerifyRequest{
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
		DeviceID:        deviceID,            // Lấy từ context hoặc request
		ClientPublicKey: req.ClientPublicKey, // Diffie-Hellman public key từ client
		// Fullname:   req.Fullname,
	}

	result, err := h.AuthUsecase.VerifyOtp(ctx, *otpDto)
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
		OrganizationId:   result.OrganizationID,
		Role:             result.Role,
		PasswordAuthId:   result.PasswordAuthId,
		ReferralCode:     result.ReferralCode,     // Mã giới thiệu
		RequireFullname:  result.RequireFullname,  // User cần cập nhật họ tên
		TempToken:        result.TempToken,        // TEMP token (nếu chưa có họ tên)
		EncryptedAuthKey: result.EncryptedAuthKey, // Auth key đã mã hóa (Diffie-Hellman)
		ServerPublicKey:  result.ServerPublicKey,  // Server public key (Diffie-Hellman)
	}

	return response, nil
}

// @Summary Check phone
// @Description Check phone đã tồn tại hay chưa
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.PhoneCheckRequest true "Thông tin phone"
// @Success 200 {object} authpb.PhoneCheckResponse
// @Failure 400 {object} authpb.PhoneCheckResponse
// @Failure 500 {object} authpb.PhoneCheckResponse
// @Router /phone/check [get]
func (h *AuthHandler) PhoneCheck(ctx context.Context, req *authpb.PhoneCheckRequest) (*authpb.PhoneCheckResponse, error) {
	// Convert proto request to DTO
	existed, userId, err := h.AuthUsecase.PhoneCheck(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.PhoneCheckResponse{
		Existed: existed,
		UserId:  userId,
	}, nil
}

// @Summary Tạo mã PIN
// @Description Thiết lập mã PIN mới sau khi xác thực OTP
// @Tags PIN
// @Accept json
// @Produce json
// @Param body body authpb.CreatePINRequest true "Thông tin tạo PIN"
// @Success 200 {object} authpb.CreatePINResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/pin/create [post]
func (h *AuthHandler) CreatePIN(ctx context.Context, req *authpb.CreatePINRequest) (*authpb.CreatePINResponse, error) {
	dtoReq := dto.CreatePINRequest{
		PIN:   strings.TrimSpace(req.Pin),
		Phone: strings.TrimSpace(req.Phone),
		OTP:   strings.TrimSpace(req.Otp),
	}

	result, err := h.PINUsecase.CreatePIN(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.CreatePINResponse{
		AuthId:  result.AuthId,
		Message: result.Message,
	}, nil
}

// @Summary Xác minh mã PIN
// @Description Đăng nhập thông qua mã PIN đã thiết lập
// @Tags PIN
// @Accept json
// @Produce json
// @Param body body authpb.VerifyPINRequest true "Thông tin xác minh PIN"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 423 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/pin/verify [post]
func (h *AuthHandler) VerifyPIN(ctx context.Context, req *authpb.VerifyPINRequest) (*authpb.AuthSuccessResponse, error) {
	dtoReq := dto.VerifyPINRequest{
		Phone:      strings.TrimSpace(req.Phone),
		PIN:        strings.TrimSpace(req.Pin),
		Version:    strings.TrimSpace(req.Version),
		Platform:   strings.TrimSpace(req.Platform),
		OS:         strings.TrimSpace(req.Os),
		DeviceName: strings.TrimSpace(req.DeviceName),
		DeviceID:   _utils.GetDeviceIdFromContext(ctx),
	}

	result, err := h.PINUsecase.VerifyPIN(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return mapAuthSuccessResponse(result), nil
}

// @Summary Cập nhật mã PIN
// @Description Thay đổi mã PIN bằng cách xác thực mã PIN hiện tại
// @Tags PIN
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body authpb.UpdatePINRequest true "Thông tin cập nhật PIN"
// @Success 200 {object} authpb.UpdatePINResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/pin/update [put]
func (h *AuthHandler) UpdatePIN(ctx context.Context, req *authpb.UpdatePINRequest) (*authpb.UpdatePINResponse, error) {
	authID := _utils.GetAuthIdFromContext(ctx)
	if authID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Không xác định được người dùng hiện tại")
	}

	dtoReq := dto.UpdatePINRequest{
		OldPIN: strings.TrimSpace(req.OldPin),
		NewPIN: strings.TrimSpace(req.NewPin),
	}

	result, err := h.PINUsecase.UpdatePIN(ctx, authID, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.UpdatePINResponse{
		Message: result.Message,
	}, nil
}

// @Summary Cập nhật trạng thái PIN
// @Description Bật hoặc tắt đăng nhập bằng PIN sau khi xác thực mã PIN hiện tại
// @Tags PIN
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body authpb.ActivePINRequest true "Thông tin trạng thái PIN"
// @Success 200 {object} authpb.UpdatePINResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 423 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/pin/active [put]
func (h *AuthHandler) ActivePIN(ctx context.Context, req *authpb.ActivePINRequest) (*authpb.UpdatePINResponse, error) {
	authID := _utils.GetAuthIdFromContext(ctx)
	if authID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Không xác định được người dùng hiện tại")
	}

	dtoReq := dto.ActivePINRequest{
		Active: req.Active,
		PIN:    strings.TrimSpace(req.Pin),
	}

	result, err := h.PINUsecase.ActivePIN(ctx, authID, dtoReq)
	// result, err -> handle
	if err != nil {
		return nil, err
	}

	return &authpb.UpdatePINResponse{
		Message: result.Message,
	}, nil
}

// @Summary Kiểm tra trạng thái PIN
// @Description Kiểm tra người dùng hiện tại đã thiết lập PIN hay chưa
// @Tags PIN
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authpb.CheckPINExistsResponse
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/pin/check [get]
func (h *AuthHandler) CheckPINExists(ctx context.Context, _ *authpb.CheckPINExistsRequest) (*authpb.CheckPINExistsResponse, error) {
	authID := _utils.GetAuthIdFromContext(ctx)
	if authID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Không xác định được người dùng hiện tại")
	}

	result, err := h.PINUsecase.CheckPINExists(ctx, authID)
	if err != nil {
		return nil, err
	}

	return &authpb.CheckPINExistsResponse{
		Exists: result.Exists,
		AuthId: result.AuthId,
	}, nil
}

// @Summary Xác minh OTP Recovery
// @Description Xác minh OTP và gán số điện thoại vào tài khoản hiện tại (không tạo profile mới)
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.VerifyOTPRequest true "Thông tin OTP"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /otp/verify-recovery [post]
func (h *AuthHandler) VerifyRecovery(ctx context.Context, req *authpb.VerifyOTPRequest) (*authpb.AuthSuccessResponse, error) {
	// Validate required fields for session creation
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	// Convert proto request to DTO
	dto := &dto.OtpVerifyRequest{
		OtpRequest: dto.OtpRequest{
			Phone:    req.Phone,
			Fullname: req.Fullname,
		},
		Otp:             req.Otp,
		AuthID:          req.AuthId,
		Version:         req.Version,
		Platform:        req.Platform,
		OS:              req.Os,
		DeviceName:      req.DeviceName,
		ClientPublicKey: req.ClientPublicKey,
	}

	result, err := h.AuthUsecase.VerifyRecovery(ctx, *dto)
	if err != nil {
		return nil, err
	}

	return mapAuthSuccessResponse(result), nil
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
		AuthId:   req.AuthId,
	}

	result, err := h.AuthUsecase.ResendOTP(ctx, *dto)
	if err != nil {
		return nil, err
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
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
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
// @Param type path string true "Loại logout" Enums(session,device,others)
// @Param request body authpb.LogoutRequest false "Thông tin logout"
// @Router /v2/auth/logout/{type} [post]
func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	logoutType := strings.ToLower(req.GetType())
	switch logoutType {
	case "", "session", "device", "others":
		if logoutType == "" {
			logoutType = "session"
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Type logout không hợp lệ")
	}

	if logoutType == "device" && req.GetSessionId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "deviceId là bắt buộc khi type=device")
	}

	result, err := h.AuthUsecase.Logout(ctx, logoutType, req.GetSessionId())
	if err != nil {
		return nil, err
	}

	return &authpb.LogoutResponse{
		Message: result.Message,
	}, nil
}

// @Summary Đăng nhập bằng Auth Key
// @Description Đăng nhập bằng auth key đã được lưu (Diffie-Hellman key exchange)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authpb.LoginWithKeyRequest true "Login request"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 401 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /login/key [post]
func (h *AuthHandler) LoginWithKey(ctx context.Context, req *authpb.LoginWithKeyRequest) (*authpb.AuthSuccessResponse, error) {
	// Validate required fields
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

	// Lấy DeviceID từ context hoặc dùng deviceName làm fallback
	deviceID := _utils.GetDeviceIdFromContext(ctx)
	// if deviceID == "" && req.DeviceName != "" {
	// 	deviceID = req.DeviceName
	// }

	// Convert proto request to DTO
	dtoReq := dto.LoginWithKeyRequest{
		Phone:            req.Phone,
		EncryptedAuthKey: req.EncryptedAuthKey,
		ClientPublicKey:  req.ClientPublicKey,
		Version:          req.Version,
		Platform:         req.Platform,
		OS:               req.Os,
		DeviceName:       req.DeviceName,
		DeviceID:         deviceID,
	}

	result, err := h.AuthUsecase.LoginWithKey(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	response := &authpb.AuthSuccessResponse{
		AuthId:           result.AuthID,
		ProfileId:        result.ProfileID,
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		FullName:         result.Username,
		OrganizationId:   result.OrganizationID,
		Role:             result.Role,
		PasswordAuthId:   result.PasswordAuthId,
		ReferralCode:     result.ReferralCode,
		RequireFullname:  result.RequireFullname,
		TempToken:        result.TempToken,
		EncryptedAuthKey: result.EncryptedAuthKey,
		ServerPublicKey:  result.ServerPublicKey,
	}

	return response, nil
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
	result, err := h.AuthUsecase.Delete(ctx, req.Otp)
	if err != nil {
		return nil, err
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
		Provider: string(result.Provider),
		AuthName: result.AuthName,
	}, nil
}

// @Summary Khôi phục tài khoản đã bị xóa
// @Description Khôi phục tài khoản đã bị soft delete bằng cách bỏ trường deleted_at
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.RestoreDeletedAccountRequest true "Thông tin khôi phục tài khoản đã bị xóa"
// @Success 200 {object} authpb.RestoreDeletedAccountResponse
// @Failure 400 {object} authpb.RestoreDeletedAccountResponse
// @Failure 404 {object} authpb.RestoreDeletedAccountResponse
// @Failure 500 {object} authpb.RestoreDeletedAccountResponse
// @Router /restore-deleted [put]
func (h *AuthHandler) RestoreDeletedAccount(ctx context.Context, req *authpb.RestoreDeletedAccountRequest) (*authpb.RestoreDeletedAccountResponse, error) {
	// Convert proto request to DTO
	dto := &dto.RestoreDeletedAccountRequest{
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
	}

	result, err := h.AuthUsecase.RestoreDeletedAccount(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi khôi phục tài khoản đã bị xóa: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.RestoreDeletedAccountResponse{
		AuthId:   result.AuthId,
		Provider: string(result.Provider),
		AuthName: result.AuthName,
		Message:  result.Message,
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
		ReferralCode: result.ReferralCode, // Mã giới thiệu

		RequireFullname: result.RequireFullname, // User cần cập nhật họ tên
		TempToken:       result.TempToken,       // TEMP token (nếu chưa có họ tên)
	}, nil
}

// @Summary QR Status
// @Description Kiểm tra trạng thái phiên QR
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.QRStatusRequest true "Thông tin QR"
// @Success 200 {object} authpb.QRStatusResponse
// @Failure 400 {object} authpb.QRStatusResponse
// @Failure 500 {object} authpb.QRStatusResponse
// @Router /v2/auth/qr/status [post]
func (h *AuthHandler) QRStatus(ctx context.Context, req *authpb.QRStatusRequest) (*authpb.QRStatusResponse, error) {
	result, err := h.AuthUsecase.GetQRStatus(ctx, &dto.QRStatusRequest{
		SessionKey: req.SessionKey,
		SessionId:  req.SessionId,
	})
	if err != nil {
		return nil, err
	}

	return &authpb.QRStatusResponse{
		Status:  result.Status,
		Message: result.Message,
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
		FullName:       result.FullName,
		// Avatar:         result.Avatar,
		// Email:          result.Email,
		// Phone:          result.Phone,
		AccessToken:     result.AccessToken,
		ReferralCode:    result.ReferralCode,    // Mã giới thiệu
		RequireFullname: result.RequireFullname, // User cần cập nhật họ tên
		TempToken:       result.TempToken,       // TEMP token (nếu chưa có họ tên)
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
	// Validate required fields for session creation
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	// Convert proto request to DTO
	dto := &dto.AdminLoginRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.Os,
		DeviceName: req.DeviceName,
	}

	result, err := h.AuthUsecase.AdminLogin(ctx, *dto)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:          result.AuthID,
		ProfileId:       result.ProfileID,
		OrganizationId:  result.OrganizationID,
		Role:            result.Role,
		AccessToken:     result.AccessToken,
		RefreshToken:    result.RefreshToken,
		ReferralCode:    result.ReferralCode,    // Mã giới thiệu
		RequireFullname: result.RequireFullname, // User cần cập nhật họ tên
		TempToken:       result.TempToken,       // TEMP token (nếu chưa có họ tên)
	}, nil
}

// @Summary List Admins
// @Description Lấy danh sách admin từ auth_method (authId, username, fullName)
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng item trên mỗi trang"
// @Param roleId query int false "Lọc theo role ID"
// @Param name query string false "Lọc theo tên"
// @Success 200 {object} authpb.AdminListResponse
// @Failure 400 {object} authpb.AdminListResponse
// @Failure 500 {object} authpb.AdminListResponse
// @Router /v2/auth/admin/list [get]
func (h *AuthHandler) ListAdmins(ctx context.Context, req *authpb.AdminListRequest) (*authpb.AdminListResponse, error) {
	if err := requireUserAdminPermission(ctx, h.UserAdminAuthorizer, useradmin.PermissionView); err != nil {
		return nil, err
	}
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 100
	}
	if size > 100 {
		size = 100
	}

	listReq := dto.AdminListRequest{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		RoleID: req.RoleId,
		Name:   req.Name,
	}

	result, err := h.AdminUsecase.ListAdmins(ctx, listReq)
	if err != nil {
		return nil, err
	}

	adminItems := make([]*authpb.AdminItem, 0, len(result.Data))
	for _, item := range result.Data {
		adminItems = append(adminItems, &authpb.AdminItem{
			AuthId:    item.AuthID,
			ProfileId: item.ProfileID,
			Username:  item.Username,
			FullName:  item.FullName,
			Email:     item.Email,
			Phone:     item.Phone,
			Avatar:    item.Avatar,
			Role:      item.Role,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &authpb.AdminListResponse{
		Data:  adminItems,
		Total: result.Total,
	}, nil
}

func (h *AuthHandler) CreateAdminAccess(ctx context.Context, req *authpb.CreateAdminAccessRequest) (*authpb.AdminAccessResponse, error) {
	// Convert proto request to DTO
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveFrom); err == nil {
			effectiveFrom = &t
		}
	}
	if req.EffectiveTo != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveTo); err == nil {
			effectiveTo = &t
		}
	}

	dto := &dto.AdminAccessRequest{
		UserID:        req.UserId,
		IPAddress:     req.IpAddress,
		IPRange:       req.IpRange,
		DeviceID:      req.DeviceId,
		DeviceName:    req.DeviceName,
		DeviceType:    req.DeviceType,
		AuthType:      req.AuthType,
		Status:        req.Status,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		MaxDevices:    int(req.MaxDevices),
		Require2FA:    req.Require2Fa,
	}

	result, err := h.AdminAccessUsecase.CreateAdminAccess(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo quyền truy cập: %v", err)
	}

	effectiveFrom = result.EffectiveFrom
	effectiveTo = result.EffectiveTo

	resultProto := &authpb.AdminAccessResponse{
		Id:         result.ID,
		UserId:     result.UserID,
		IpAddress:  result.IPAddress,
		IpRange:    result.IPRange,
		DeviceId:   result.DeviceID,
		DeviceName: result.DeviceName,
		DeviceType: result.DeviceType,
		AuthType:   result.AuthType,
		Status:     result.Status,

		MaxDevices: int32(result.MaxDevices),
		Require2Fa: result.Require2FA,
		CreatedBy:  result.CreatedBy,
		UpdatedBy:  result.UpdatedBy,
		CreatedAt:  result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if effectiveFrom != nil {
		resultProto.EffectiveFrom = effectiveFrom.Format("2006-01-02T15:04:05Z07:00")
	}
	if effectiveTo != nil {
		resultProto.EffectiveTo = effectiveTo.Format("2006-01-02T15:04:05Z07:00")
	}
	return resultProto, nil
}

// @Summary Cập nhật quyền truy cập admin
// @Description Cập nhật quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param id path int true "ID quyền truy cập"
// @Param body body authpb.UpdateAdminAccessRequest true "Thông tin cập nhật"
// @Success 200 {object} authpb.AdminAccessResponse
// @Failure 400 {object} authpb.AdminAccessResponse
// @Failure 401 {object} authpb.AdminAccessResponse
// @Failure 404 {object} authpb.AdminAccessResponse
// @Failure 500 {object} authpb.AdminAccessResponse
// @Router /admin-access/{id} [put]
func (h *AuthHandler) UpdateAdminAccess(ctx context.Context, req *authpb.UpdateAdminAccessRequest) (*authpb.AdminAccessResponse, error) {
	// Convert proto request to DTO
	var effectiveFrom, effectiveTo *time.Time
	if req.EffectiveFrom != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveFrom); err == nil {
			effectiveFrom = &t
		}
	}
	if req.EffectiveTo != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", req.EffectiveTo); err == nil {
			effectiveTo = &t
		}
	}

	dto := &dto.AdminAccessRequest{
		UserID:        req.UserId,
		IPAddress:     req.IpAddress,
		IPRange:       req.IpRange,
		DeviceID:      req.DeviceId,
		DeviceName:    req.DeviceName,
		DeviceType:    req.DeviceType,
		AuthType:      req.AuthType,
		Status:        req.Status,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		MaxDevices:    int(req.MaxDevices),
		Require2FA:    req.Require2Fa,
	}

	result, err := h.AdminAccessUsecase.UpdateAdminAccess(ctx, req.Id, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi cập nhật quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	response := &authpb.AdminAccessResponse{
		Id:         result.ID,
		UserId:     result.UserID,
		IpAddress:  result.IPAddress,
		IpRange:    result.IPRange,
		DeviceId:   result.DeviceID,
		DeviceName: result.DeviceName,
		DeviceType: result.DeviceType,
		AuthType:   result.AuthType,
		Status:     result.Status,
		MaxDevices: int32(result.MaxDevices),
		Require2Fa: result.Require2FA,
		CreatedBy:  result.CreatedBy,
		UpdatedBy:  result.UpdatedBy,
		CreatedAt:  result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle nullable time fields
	if result.EffectiveFrom != nil {
		response.EffectiveFrom = result.EffectiveFrom.Format("2006-01-02T15:04:05Z07:00")
	}
	if result.EffectiveTo != nil {
		response.EffectiveTo = result.EffectiveTo.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// @Summary Xóa quyền truy cập admin
// @Description Xóa quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param id path int true "ID quyền truy cập"
// @Success 200 {object} authpb.DeleteAdminAccessResponse
// @Failure 400 {object} authpb.DeleteAdminAccessResponse
// @Failure 401 {object} authpb.DeleteAdminAccessResponse
// @Failure 404 {object} authpb.DeleteAdminAccessResponse
// @Failure 500 {object} authpb.DeleteAdminAccessResponse
// @Router /admin-access/{id} [delete]
func (h *AuthHandler) DeleteAdminAccess(ctx context.Context, req *authpb.DeleteAdminAccessRequest) (*authpb.DeleteAdminAccessResponse, error) {
	err := h.AdminAccessUsecase.DeleteAdminAccess(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi xóa quyền truy cập: %v", err)
	}

	return &authpb.DeleteAdminAccessResponse{
		Message: "Quyền truy cập đã được xóa thành công",
	}, nil
}

// @Summary Lấy danh sách quyền truy cập admin
// @Description Lấy danh sách quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param userId query int false "ID người dùng"
// @Param status query string false "Trạng thái"
// @Param authType query string false "Loại xác thực"
// @Param ipAddress query string false "Địa chỉ IP"
// @Success 200 {object} authpb.GetAdminAccessListResponse
// @Failure 400 {object} authpb.GetAdminAccessListResponse
// @Failure 401 {object} authpb.GetAdminAccessListResponse
// @Failure 500 {object} authpb.GetAdminAccessListResponse
// @Router /admin-access [get]
func (h *AuthHandler) GetAdminAccessList(ctx context.Context, req *authpb.GetAdminAccessListRequest) (*authpb.GetAdminAccessListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		Status:    req.Status,
		AuthType:  req.AuthType,
		IPAddress: req.IpAddress,
	}

	if req.UserId > 0 {
		dto.UserID = &req.UserId
	}

	result, err := h.AdminAccessUsecase.GetAdminAccessList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*authpb.AdminAccessItem, len(result.Data))
	for i, item := range result.Data {
		protoItem := &authpb.AdminAccessItem{
			Id:     item.ID,
			UserId: item.UserID,
			// Username:   item.Username,
			IpAddress:  item.IPAddress,
			IpRange:    item.IPRange,
			DeviceId:   item.DeviceID,
			DeviceName: item.DeviceName,
			DeviceType: item.DeviceType,
			AuthType:   item.AuthType,
			Status:     item.Status,
			MaxDevices: int32(item.MaxDevices),
			Require2Fa: item.Require2FA,
			CreatedBy:  item.CreatedBy,
			UpdatedBy:  item.UpdatedBy,
			CreatedAt:  _utils.FormatTimeToString(&item.CreatedAt),
			UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Handle nullable time fields
		if item.EffectiveFrom != nil {
			protoItem.EffectiveFrom = _utils.FormatTimeToString(item.EffectiveFrom)
		}
		if item.EffectiveTo != nil {
			protoItem.EffectiveTo = _utils.FormatTimeToString(item.EffectiveTo)
		}

		items[i] = protoItem
	}

	return &authpb.GetAdminAccessListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// @Summary Kiểm tra quyền truy cập admin
// @Description Kiểm tra quyền truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param body body authpb.ValidateAdminAccessRequest true "Thông tin kiểm tra"
// @Success 200 {object} authpb.ValidateAdminAccessResponse
// @Failure 400 {object} authpb.ValidateAdminAccessResponse
// @Failure 500 {object} authpb.ValidateAdminAccessResponse
// @Router /admin-access/validate [post]
func (h *AuthHandler) ValidateAdminAccess(ctx context.Context, req *authpb.ValidateAdminAccessRequest) (*authpb.ValidateAdminAccessResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessValidationRequest{
		UserID:    req.UserId,
		IPAddress: req.IpAddress,
		DeviceID:  req.DeviceId,
	}

	result, err := h.AdminAccessUsecase.ValidateAdminAccess(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi kiểm tra quyền truy cập: %v", err)
	}

	// Convert DTO response to proto response
	return &authpb.ValidateAdminAccessResponse{
		IsAllowed:  result.IsAllowed,
		Reason:     result.Reason,
		Require2Fa: result.Require2FA,
		AccessId:   result.AccessID,
		DeviceId:   result.DeviceID,
		IpAddress:  result.IPAddress,
	}, nil
}

// @Summary Tạo log truy cập admin
// @Description Tạo log truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param body body authpb.CreateAdminAccessLogRequest true "Thông tin log"
// @Success 200 {object} authpb.CreateAdminAccessLogResponse
// @Failure 400 {object} authpb.CreateAdminAccessLogResponse
// @Failure 500 {object} authpb.CreateAdminAccessLogResponse
// @Router /admin-access/log [post]
func (h *AuthHandler) CreateAdminAccessLog(ctx context.Context, req *authpb.CreateAdminAccessLogRequest) (*authpb.CreateAdminAccessLogResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessLogRequest{
		UserID:     req.UserId,
		IPAddress:  req.IpAddress,
		DeviceID:   req.DeviceId,
		DeviceName: req.DeviceName,
		Action:     req.Action,
		Status:     req.Status,
		Reason:     req.Reason,
		UserAgent:  req.UserAgent,
	}

	err := h.AdminAccessUsecase.CreateAdminAccessLog(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi tạo log truy cập: %v", err)
	}

	return &authpb.CreateAdminAccessLogResponse{
		Message: "Log truy cập đã được tạo thành công",
	}, nil
}

// @Summary Lấy danh sách log truy cập admin
// @Description Lấy danh sách log truy cập hệ thống nội bộ
// @Tags Admin Access Control
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param userId query int false "ID người dùng"
// @Param ipAddress query string false "Địa chỉ IP"
// @Success 200 {object} authpb.GetAdminAccessLogListResponse
// @Failure 400 {object} authpb.GetAdminAccessLogListResponse
// @Failure 401 {object} authpb.GetAdminAccessLogListResponse
// @Failure 500 {object} authpb.GetAdminAccessLogListResponse
// @Router /admin-access/log [get]
func (h *AuthHandler) GetAdminAccessLogList(ctx context.Context, req *authpb.GetAdminAccessLogListRequest) (*authpb.GetAdminAccessLogListResponse, error) {
	// Convert proto request to DTO
	dto := &dto.AdminAccessLogListRequest{
		Page:      int(req.Page),
		Size:      int(req.Size),
		UserID:    req.UserId,
		IPAddress: req.IpAddress,
	}

	result, err := h.AdminAccessUsecase.GetAdminAccessLogList(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi lấy danh sách log truy cập: %v", err)
	}

	// Convert DTO response to proto response
	items := make([]*authpb.AdminAccessLogItem, len(result.Data))
	for i, item := range result.Data {
		items[i] = &authpb.AdminAccessLogItem{
			Id:         item.ID,
			UserId:     item.UserID,
			IpAddress:  item.IPAddress,
			DeviceId:   item.DeviceID,
			DeviceName: item.DeviceName,
			Action:     item.Action,
			Status:     item.Status,
			Reason:     item.Reason,
			UserAgent:  item.UserAgent,
			CreatedAt:  _utils.FormatTimeToString(&item.CreatedAt),
		}
	}

	return &authpb.GetAdminAccessLogListResponse{
		Data:  items,
		Total: int32(result.Total),
		Page:  int32(result.Page),
		Size:  int32(result.Size),
	}, nil
}

// @Summary Tạo mật khẩu
// @Description Tạo mật khẩu cho user chưa có mật khẩu
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.CreatePasswordRequest true "Thông tin mật khẩu"
// @Success 200 {object} authpb.CreatePasswordResponse
// @Failure 400 {object} authpb.CreatePasswordResponse
// @Failure 401 {object} authpb.CreatePasswordResponse
// @Failure 500 {object} authpb.CreatePasswordResponse
// @Router /password/create [post]
func (h *AuthHandler) CreatePassword(ctx context.Context, req *authpb.CreatePasswordRequest) (*authpb.CreatePasswordResponse, error) {
	// Convert proto request to DTO
	dtoReq := dto.CreatePasswordRequest{
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	result, err := h.AuthUsecase.CreatePassword(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.CreatePasswordResponse{
		AuthId:  result.AuthId,
		Message: result.Message,
		Success: result.Success,
	}, nil
}

// @Summary Đăng nhập bằng mật khẩu
// @Description Đăng nhập bằng username/phone và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.LoginWithPasswordRequest true "Thông tin đăng nhập"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 401 {object} authpb.AuthSuccessResponse
// @Failure 403 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /password/login [post]
func (h *AuthHandler) LoginWithPassword(ctx context.Context, req *authpb.LoginWithPasswordRequest) (*authpb.AuthSuccessResponse, error) {
	// Validate required fields for session creation
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	// Convert proto request to DTO
	dtoReq := dto.LoginWithPasswordRequest{
		Username:   req.Username,
		Password:   req.Password,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.Os,
		DeviceName: req.DeviceName,
	}

	result, err := h.AuthUsecase.LoginWithPassword(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	return &authpb.AuthSuccessResponse{
		AuthId:         result.AuthID,
		ProfileId:      result.ProfileID,
		AccessToken:    result.AccessToken,
		RefreshToken:   result.RefreshToken,
		FullName:       result.FullName,
		OrganizationId: result.OrganizationID,
		Role:           result.Role,
		PasswordAuthId: result.PasswordAuthId,
		ReferralCode:   result.ReferralCode, // Mã giới thiệu

		RequireFullname: result.RequireFullname, // User cần cập nhật họ tên
		TempToken:       result.TempToken,       // TEMP token (nếu chưa có họ tên)
	}, nil
}

// @Summary Đăng ký thiết bị
// @Description Lưu thông tin thiết bị của ứng dụng khi khởi động
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.CreateDeviceRequest true "Thông tin thiết bị"
// @Success 200 {object} authpb.CreateDeviceResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /device/check [post]
func (h *AuthHandler) RegisterDevice(ctx context.Context, req *authpb.CreateDeviceRequest) (*authpb.CreateDeviceResponse, error) {
	dtoReq := h.DeviceMapper.MapCreateDeviceRequestPbToDTO(req)
	result, err := h.DeviceUsecase.RegisterDevice(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &authpb.CreateDeviceResponse{
		Data: h.DeviceMapper.MapDeviceToPb(result),
	}, nil
}

// validateAppState kiểm tra app state có hợp lệ không
// Các giá trị hợp lệ: 'active', 'background', 'inactive', 'unknown', 'extension'
func validateAppState(state string) error {
	validStates := map[string]bool{
		"active":     true,
		"background": true,
		"inactive":   true,
		"unknown":    true,
		"extension":  true,
	}

	if state == "" {
		return status.Errorf(codes.InvalidArgument, "App state không được để trống")
	}

	if !validStates[state] {
		return status.Errorf(codes.InvalidArgument, "App state không hợp lệ. Giá trị hợp lệ: 'active', 'background', 'inactive', 'unknown', 'extension'")
	}

	return nil
}

// @Summary Cập nhật trạng thái thiết bị
// @Description Cập nhật trạng thái app (foreground/background) của thiết bị
// @Tags Device
// @Produce json
// @Security BearerAuth
// @Param state path string true "Trạng thái thiết bị (active/background/inactive/unknown/extension)"
// @Success 200 {object} authpb.AppStateResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/device/state/{state} [put]
func (h *AuthHandler) UpdateDeviceMode(ctx context.Context, req *authpb.AppStateRequest) (*authpb.AppStateResponse, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(service.CurrentUserUnknown)
	}

	// Validate app state
	state := req.State
	if err := validateAppState(state); err != nil {
		return nil, err
	}

	// Lưu app state vào Redis
	err := h.AuthUsecase.CacheProvider.SaveAppState(ctx, profileID, state)
	if err != nil {
		return nil, err
	}

	return &authpb.AppStateResponse{
		Message: "Cập nhật trạng thái thiết bị thành công",
	}, nil
}

// @Summary Lấy thông tin bảo mật tài khoản hiện tại
// @Description Trả về trạng thái, PIN và danh sách thiết bị của người dùng đang đăng nhập
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} sharepb.AuthDataV3Proto
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/security/me [get]
func (h *AuthHandler) GetCurrentAuthSetting(ctx context.Context, _ *sharepb.Empty) (*sharepb.AuthDataV3Proto, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Không xác định được người dùng hiện tại")
	}

	result, err := h.AuthSecurityUsecase.GetAuthSecurityDataByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return h.AuthSecurityMapper.MapToProto(result), nil
}

// @Summary Chuyển đổi tài khoản
// @Description Chuyển đổi sang tài khoản cá nhân hoặc tổ chức khác
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body authpb.SwitchAccountRequest true "Thông tin chuyển đổi tài khoản"
// @Success 200 {object} authpb.SwitchAccountResponse
// @Failure 400 {object} authpb.SwitchAccountResponse
// @Failure 500 {object} authpb.SwitchAccountResponse
// @Router /v2/auth/switch/user [post]
func (h *AuthHandler) SwitchAccount(ctx context.Context, req *authpb.SwitchAccountRequest) (*authpb.SwitchAccountResponse, error) {
	// Convert proto request to DTO
	dto := &dto.SwitchAccountRequest{
		ProfileID:      req.ProfileId,
		OrganizationID: req.OrganizationId,
		// GroupID:        req.GroupId,
		OwnerType: req.OwnerType,
	}

	result, err := h.AuthUsecase.SwitchAccount(ctx, *dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "lỗi khi chuyển đổi tài khoản: %v", err)
	}

	// Convert DTO response to proto response
	// var groupId *uint64
	// if result.GroupID != nil {
	// 	groupIdVal := uint64(*result.GroupID)
	// 	groupId = &groupIdVal
	// }

	return &authpb.SwitchAccountResponse{
		AuthId:         result.AuthID,
		ProfileId:      result.ProfileID,
		OrganizationId: result.OrganizationID,
		Role:           result.Role,
		FullName:       result.FullName,
		Avatar:         result.Avatar,
		Email:          result.Email,
		Phone:          result.Phone,
		AccessToken:    result.AccessToken,

		// GroupId:        groupId,
	}, nil
}

// @Summary Lấy lịch sử đăng nhập
// @Description Lấy danh sách lịch sử đăng nhập của user hiện tại
// @Tags Auth
// @Accept json
// @Produce json
// @Param page query int false "Số trang" default(1)
// @Param size query int false "Số lượng mỗi trang" default(20)
// @Success 200 {object} authpb.GetLoginHistoryResponse
// @Failure 401 {object} authpb.GetLoginHistoryResponse
// @Failure 500 {object} authpb.GetLoginHistoryResponse
// @Router /v2/auth/login-history [get]
// @Security BearerAuth
func (h *AuthHandler) GetLoginHistory(ctx context.Context, req *authpb.GetLoginHistoryRequest) (*authpb.GetLoginHistoryResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	size := req.Size
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	sessions, total, err := h.AuthUsecase.GetLoginHistory(ctx, page, size)
	if err != nil {
		return nil, err
	}

	items := h.SessionMapper.MapSessionsToLoginHistoryListPb(sessions)

	return &authpb.GetLoginHistoryResponse{
		Data:  items,
		Total: total,
	}, nil
}

// @Summary Lấy danh sách session
// @Description Lấy danh sách session của người dùng hiện tại
// @Tags Auth
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param size query int false "Số lượng mỗi trang" default(20)
// @Success 200 {object} authpb.SessionPageResponse
// @Failure 400 {object} sharepb.ErrorResponse
// @Failure 401 {object} sharepb.ErrorResponse
// @Failure 500 {object} sharepb.ErrorResponse
// @Router /v2/auth/sessions [get]
// @Security BearerAuth
func (h *AuthHandler) GetSessionsByProfile(ctx context.Context, req *sharepb.RequestV3Proto) (*authpb.SessionPageResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Request không hợp lệ")
	}

	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}

	sessions, total, err := h.AuthUsecase.GetSessionsByProfile(ctx, pagable)
	if err != nil {
		return nil, err
	}

	items := h.SessionMapper.MapSessionsToSessionV3ProtoList(sessions)

	return &authpb.SessionPageResponse{
		Data:  items,
		Total: total,
	}, nil
}

// @Summary Update Fullname
// @Description Cập nhật họ tên cho user đã đăng nhập và trả về thông tin login mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authpb.UpdateFullnameRequest true "Họ tên mới"
// @Success 200 {object} authpb.AuthSuccessResponse
// @Failure 400 {object} authpb.AuthSuccessResponse
// @Failure 401 {object} authpb.AuthSuccessResponse
// @Failure 500 {object} authpb.AuthSuccessResponse
// @Router /v2/auth/fullname [put]
// @Security BearerAuth
func (h *AuthHandler) UpdateFullname(ctx context.Context, req *authpb.UpdateFullnameRequest) (*authpb.AuthSuccessResponse, error) {
	// Validate request
	if req.Fullname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Họ tên không được để trống")
	}

	// Gọi usecase để cập nhật fullname
	result, err := h.AuthUsecase.UpdateFullname(ctx, req.Fullname)
	if err != nil {
		return nil, err
	}

	// Convert DTO response to proto response
	response := &authpb.AuthSuccessResponse{
		AuthId:          result.AuthID,
		ProfileId:       result.ProfileID,
		AccessToken:     result.AccessToken,
		RefreshToken:    result.RefreshToken,
		FullName:        result.FullName,
		OrganizationId:  result.OrganizationID,
		Role:            result.Role,
		PasswordAuthId:  result.PasswordAuthId,
		ReferralCode:    result.ReferralCode,
		RequireFullname: result.RequireFullname,
		TempToken:       result.TempToken,
	}

	return response, nil
}

// LoginWithToken is deprecated - use LoginWithKey instead

// @Summary Khóa tài khoản của người dùng
// @Description Khóa tài khoản của người dùng
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.OTPV3Proto true "Thông tin khóa tài khoản"
// @Success 200 {object} sharepb.Empty
// @Router /v2/auth/account/lock [put]
func (h *AuthHandler) OwnerLockAccount(ctx context.Context, req *sharepb.OTPV3Proto) (*sharepb.Empty, error) {
	_, err := h.AuthUsecase.LockAccount(ctx, req.Otp)
	if err != nil {
		return nil, err
	}

	return &sharepb.Empty{}, nil
}

func mapAuthSuccessResponse(result *dto.AuthLoginResponse) *authpb.AuthSuccessResponse {
	if result == nil {
		return &authpb.AuthSuccessResponse{}
	}

	fullName := result.FullName
	if fullName == "" {
		fullName = result.Username
	}

	return &authpb.AuthSuccessResponse{
		AuthId:           result.AuthID,
		ProfileId:        result.ProfileID,
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		FullName:         fullName,
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
