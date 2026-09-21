package usecase

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log/slog"
	"time"
	"user/internal"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/providers"
	"user/internal/interface/repo"

	"golang.org/x/crypto/bcrypt"
)

type AccessTokenIssuer interface {
	GenerateAccessToken(ctx context.Context, authID uint64, sessionID uint64, organizationID *uint64) (string, error)
}

var _ AccessTokenIssuer = (*AuthUsecase)(nil)

type PINUsecase struct {
	pinRepo              repo.PINRepository
	statusRepo           repo.StatusRepository
	authMethodRepo       repo.AuthMethodRepository
	sessionRepo          repo.SessionRepository
	otpRepo              repo.OTPRepository
	otpUsecase           *OtpUsecase
	profileProvider      providers.ProfileProvider
	notificationProvider providers.NotificationProvider
	deviceUsecase        *DeviceUsecase
	cacheProvider        providers.CacheProvider
	tokenIssuer          AccessTokenIssuer
	maxPINAttempts       int
	pinLockDuration      time.Duration
}

func NewPINUsecase(
	pinRepo repo.PINRepository,
	statusRepo repo.StatusRepository,
	authMethodRepo repo.AuthMethodRepository,
	sessionRepo repo.SessionRepository,
	otpRepo repo.OTPRepository,
	otpUsecase *OtpUsecase,
	profileProvider providers.ProfileProvider,
	notificationProvider providers.NotificationProvider,
	deviceUsecase *DeviceUsecase,
	cacheProvider providers.CacheProvider,
	tokenIssuer AccessTokenIssuer,
) *PINUsecase {
	return &PINUsecase{
		pinRepo:              pinRepo,
		statusRepo:           statusRepo,
		authMethodRepo:       authMethodRepo,
		sessionRepo:          sessionRepo,
		otpRepo:              otpRepo,
		otpUsecase:           otpUsecase,
		profileProvider:      profileProvider,
		notificationProvider: notificationProvider,
		deviceUsecase:        deviceUsecase,
		cacheProvider:        cacheProvider,
		tokenIssuer:          tokenIssuer,
		maxPINAttempts:       3,
		pinLockDuration:      15 * time.Minute,
	}
}

// CreatePIN tạo mã PIN mới cho user
func (s *PINUsecase) CreatePIN(c context.Context, request dto.CreatePINRequest) (*dto.CreatePINResponse, error) {
	// Validate PIN (phải là 6 số)
	if len(request.PIN) != 6 {
		return nil, _errors.ReturnError(service.PINLengthInvalid)
	}
	if len(request.OTP) != 6 {
		return nil, _errors.ReturnError(service.OTPInvalid)
	}

	authID := _utils.GetAuthIdFromContext(c)
	// Tìm auth method theo phone
	authEntity, err := s.authMethodRepo.FindByID(c, authID)
	if err != nil {
		return nil, _errors.ReturnError(service.AccountNotFoundGeneric)
	}

	// Validate OTP trước khi tạo PIN
	param, err := s.buildAuthParam(c, authEntity, request.OTP)
	if err != nil {
		return nil, err
	}

	if err := s.otpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED,
		enums.CHECK_OTP,
	}, param); err != nil {
		return nil, err
	}

	// Kiểm tra xem đã có PIN chưa
	exists, err := s.pinRepo.CheckPINExists(c, authEntity.ID)
	if err == nil && exists {
		return nil, _errors.ReturnError(service.PINAlreadyExists)
	}

	// Hash PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(request.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash PIN: %w", err)
	}

	// Tạo PIN entity
	now := time.Now()
	pinEntity := &auth.UserPINEntity{
		IdDomain: auth.IdDomain{
			AuthID: authEntity.ID,
		},
		PIN:          string(hashedPIN),
		PINCheckTime: 0,
		PINDate:      &now,
		IsActive:     true,
	}

	// Lưu vào database
	err = s.pinRepo.Create(c, pinEntity)
	if err != nil {
		return nil, fmt.Errorf("create PIN: %w", err)
	}

	s.logPINActionHistory(c, authEntity, "pin", "PIN_CREATED",
		"Tạo mã PIN mới cho tài khoản", map[string]interface{}{
			"createdAt": now.Format(time.RFC3339),
		})

	return &dto.CreatePINResponse{
		AuthId:  authEntity.ID,
		Message: "Tạo mã PIN thành công",
	}, nil
}

// VerifyPIN xác minh mã PIN và tạo session
func (s *PINUsecase) VerifyPIN(c context.Context, request dto.VerifyPINRequest) (*dto.AuthLoginResponse, error) {
	// Tìm auth method theo phone
	authEntity, err := s.authMethodRepo.FindByPhone(c, request.Phone)
	if err != nil {
		return nil, _errors.ReturnError(service.AccountNotFoundGeneric)
	}

	_, err = s.verifyPINCode(c, authEntity, request.PIN, &dto.PinVerifyMeta{
		Phone:      request.Phone,
		Platform:   request.Platform,
		OS:         request.OS,
		DeviceID:   request.DeviceID,
		DeviceName: request.DeviceName,
	}, true)
	if err != nil {
		return nil, err
	}

	// Tạo session trước khi phát access token; token luôn bind tới durable session ID.
	newSession := s.createSession(c, authEntity, request)
	if err := s.sessionRepo.CreateSession(c, newSession); err != nil {
		return nil, fmt.Errorf("create login session: %w", err)
	}

	// Lấy thông tin profile
	profile, err := s.profileProvider.GetByProfileID(c, authEntity.UserID)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	if s.tokenIssuer == nil {
		return nil, fmt.Errorf("token issuer is not configured")
	}

	accessToken, err := s.tokenIssuer.GenerateAccessToken(c, authEntity.ID, newSession.SessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	response := &dto.AuthLoginResponse{
		AuthID:       authEntity.ID,
		ProfileID:    profile.ProfileID,
		AccessToken:  accessToken,
		Role:         "ROLE_USER",
		Username:     profile.FullName,
		FullName:     profile.FullName,
		ReferralCode: profile.ReferralCode,
	}

	// Cập nhật device với authId và profileId nếu có thông tin thiết bị
	if s.deviceUsecase != nil {
		deviceID := request.DeviceID
		if deviceID == "" {
			deviceID = request.DeviceName
		}
		if deviceID != "" && authEntity.ID > 0 && profile.ProfileID > 0 {
			deviceCtx := context.WithValue(c, _enum.ProfileIDKey, profile.ProfileID)
			deviceCtx = context.WithValue(deviceCtx, _enum.AuthIDKey, authEntity.ID)
			deviceCtx, cancel := context.WithTimeout(deviceCtx, 5*time.Second)
			defer cancel()

			if err := s.deviceUsecase.UpdateDeviceWithAuthInfo(deviceCtx, deviceID, authEntity.ID, profile.ProfileID); err != nil {
				slog.ErrorContext(c, fmt.Sprintf("Failed to update device with auth info after verify PIN: %v", err))
			}
		}
	}

	// Lưu app state = "active" khi login thành công
	if profile.ProfileID != 0 && s.cacheProvider != nil {
		if err := s.cacheProvider.SaveAppState(c, profile.ProfileID, "active"); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to save app state for user %d: %v", profile.ProfileID, err))
		}
	}

	return response, nil
}

// UpdatePIN cập nhật mã PIN
func (s *PINUsecase) UpdatePIN(c context.Context, authID uint64, request dto.UpdatePINRequest) (*dto.UpdatePINResponse, error) {
	// Validate PIN mới
	if len(request.NewPIN) != 6 {
		return nil, _errors.ReturnError(service.NewPINLengthInvalid)
	}

	authEntity, err := s.authMethodRepo.FindByID(c, authID)
	if err != nil {
		return nil, _errors.ReturnError(service.AccountNotFoundGeneric)
	}

	// Lấy PIN hiện tại
	pinEntity, err := s.pinRepo.GetByAuthID(c, authID)
	if err != nil {
		return nil, _errors.ReturnError(service.PINNotSet)
	}

	// Verify PIN cũ
	err = bcrypt.CompareHashAndPassword([]byte(pinEntity.PIN), []byte(request.OldPIN))
	if err != nil {
		return nil, _errors.ReturnError(service.OldPINInvalid)
	}

	// Hash PIN mới
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(request.NewPIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash new PIN: %w", err)
	}

	// Update PIN
	now := time.Now()
	pinEntity.PIN = string(hashedPIN)
	pinEntity.PINDate = &now
	pinEntity.PINCheckTime = 0
	pinEntity.LockedUntil = nil

	err = s.pinRepo.Update(c, pinEntity)
	if err != nil {
		return nil, fmt.Errorf("save new PIN: %w", err)
	}

	s.logPINActionHistory(c, authEntity, "pin", "PIN_UPDATED",
		"Cập nhật mã PIN thành công", map[string]interface{}{
			"updatedAt": now.Format(time.RFC3339),
		})

	return &dto.UpdatePINResponse{
		Message: "Cập nhật mã PIN thành công",
	}, nil
}

// ActivePIN bật/tắt PIN của người dùng hiện tại
func (s *PINUsecase) ActivePIN(c context.Context, authID uint64, request dto.ActivePINRequest) (*dto.UpdatePINResponse, error) {
	if len(request.PIN) != 6 {
		return nil, _errors.ReturnError(service.PINLengthInvalid)
	}

	authEntity, err := s.authMethodRepo.FindByID(c, authID)
	if err != nil {
		return nil, _errors.ReturnError(service.AccountNotFoundGeneric)
	}

	pinEntity, err := s.verifyPINCode(c, authEntity, request.PIN, &dto.PinVerifyMeta{
		Phone: authEntity.AuthName,
	}, !request.Active)
	if err != nil {
		return nil, err
	}

	pinEntity.IsActive = request.Active
	pinEntity.PINCheckTime = 0
	pinEntity.LockedUntil = nil

	if err := s.pinRepo.Update(c, pinEntity); err != nil {
		return nil, fmt.Errorf("update PIN status: %w", err)
	}

	status := "vô hiệu hóa"
	if request.Active {
		status = "kích hoạt"
	}

	actionName := "PIN_DEACTIVATED"
	if request.Active {
		actionName = "PIN_ACTIVATED"
	}
	s.logPINActionHistory(c, authEntity, "pin", actionName,
		fmt.Sprintf("Mã PIN đã được %s", status), map[string]interface{}{
			"active": request.Active,
		})

	return &dto.UpdatePINResponse{
		Message: fmt.Sprintf("Mã PIN đã được %s", status),
	}, nil
}

// CheckPINExists kiểm tra xem user đã có PIN chưa
func (s *PINUsecase) CheckPINExists(c context.Context, authID uint64) (*dto.CheckPINExistsResponse, error) {
	exists, err := s.pinRepo.CheckPINExists(c, authID)
	if err != nil {
		return &dto.CheckPINExistsResponse{
			Exists: false,
			AuthId: authID,
		}, nil
	}

	return &dto.CheckPINExistsResponse{
		Exists: exists,
		AuthId: authID,
	}, nil
}

func (s *PINUsecase) verifyPINCode(ctx context.Context, authEntity *auth.AuthMethod, inputPIN string, meta *dto.PinVerifyMeta, requireActive bool) (*auth.UserPINEntity, error) {
	pinEntity, err := s.pinRepo.GetByAuthID(ctx, authEntity.ID)
	if err != nil {
		return nil, _errors.ReturnError(service.PINNotSet)
	}

	if pinEntity.LockedUntil != nil && pinEntity.LockedUntil.After(time.Now()) {
		minutes := time.Until(*pinEntity.LockedUntil).Minutes()
		return nil, _errors.ReturnError(
			service.PINLocked,
			_errors.WithPublicMessage(fmt.Sprintf("Mã PIN đã bị khóa. Vui lòng thử lại sau %.0f phút", minutes)),
		)
	}

	if requireActive && !pinEntity.IsActive {
		return nil, _errors.ReturnError(service.PINDisabled)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pinEntity.PIN), []byte(inputPIN)); err != nil {
		_ = s.pinRepo.IncrementCheckTime(ctx, authEntity.ID)
		pinEntity.PINCheckTime++

		if pinEntity.PINCheckTime >= s.maxPINAttempts {
			lockedUntil := time.Now().Add(s.pinLockDuration)
			_ = s.pinRepo.LockPIN(ctx, authEntity.ID, lockedUntil)
			s.logPINLockHistory(ctx, authEntity, meta, lockedUntil)
			return nil, _errors.ReturnError(
				service.PINAttemptsExceeded,
				_errors.WithPublicMessage(fmt.Sprintf("Đã nhập sai %d lần. Mã PIN đã bị khóa trong %d phút", s.maxPINAttempts, int(s.pinLockDuration.Minutes()))),
			)
		}

		remaining := s.maxPINAttempts - pinEntity.PINCheckTime
		return nil, _errors.ReturnError(
			service.PINIncorrect,
			_errors.WithPublicMessage(fmt.Sprintf("Mã PIN không đúng. Bạn còn %d lần thử", remaining)),
		)
	}

	if err := s.pinRepo.ResetCheckCounter(ctx, authEntity.ID); err != nil {
		return nil, fmt.Errorf("update PIN status: %w", err)
	}

	pinEntity.PINCheckTime = 0
	pinEntity.LockedUntil = nil

	return pinEntity, nil
}

func (s *PINUsecase) logPINLockHistory(ctx context.Context, authEntity *auth.AuthMethod, meta *dto.PinVerifyMeta, lockedUntil time.Time) {
	if s.notificationProvider == nil || authEntity == nil || authEntity.UserID == 0 {
		return
	}

	ipAddress := _utils.AuditIPFromContext(ctx)
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}

	userAgent := _utils.GetUserAgentFromContext(ctx)
	if userAgent == "" {
		userAgent = "Mobile-Client"
	}

	deviceID := _utils.GetDeviceIdFromContext(ctx)
	if meta != nil {
		if deviceID == "" {
			deviceID = meta.DeviceID
		}
		if deviceID == "" {
			deviceID = meta.DeviceName
		}
	}

	phone := ""
	platform := ""
	osName := ""
	if meta != nil {
		phone = meta.Phone
		platform = meta.Platform
		osName = meta.OS
	}

	metadata := map[string]interface{}{
		"authId":      authEntity.ID,
		"phone":       phone,
		"lockedUntil": lockedUntil.Format(time.RFC3339),
		"maxAttempts": s.maxPINAttempts,
	}
	if platform != "" {
		metadata["platform"] = platform
	}
	if osName != "" {
		metadata["os"] = osName
	}
	if deviceID != "" {
		metadata["deviceId"] = deviceID
	}

	success := false
	description := fmt.Sprintf("Nhập sai mã PIN quá %d lần. Mã PIN đã bị khóa trong %d phút", s.maxPINAttempts, int(s.pinLockDuration.Minutes()))
	payload := dto.HistoryAuthCreateDTO{
		UserID:        authEntity.UserID,
		ActionType:    "pin_lock",
		ActionName:    "PIN_LOCKED",
		Description:   description,
		Success:       &success,
		Reason:        "PIN_ATTEMPTS_EXCEEDED",
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Channel:       platform,
		DeviceID:      deviceID,
		SourceService: "auth-service",
		Metadata:      metadata,
	}

	s.dispatchHistoryAuth(ctx, payload, "failed to log PIN lock history")
}

func (s *PINUsecase) logPINActionHistory(ctx context.Context, authEntity *auth.AuthMethod, actionType, actionName, description string, metadata map[string]interface{}) {
	if s.notificationProvider == nil || authEntity == nil || authEntity.UserID == 0 {
		return
	}

	ipAddress := _utils.AuditIPFromContext(ctx)
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}

	userAgent := _utils.GetUserAgentFromContext(ctx)
	if userAgent == "" {
		userAgent = "Mobile-Client"
	}

	deviceID := _utils.GetDeviceIdFromContext(ctx)

	mergedMetadata := map[string]interface{}{
		"authId": authEntity.ID,
		"phone":  authEntity.AuthName,
	}
	for k, v := range metadata {
		if mergedMetadata == nil {
			mergedMetadata = map[string]interface{}{}
		}
		mergedMetadata[k] = v
	}

	success := true
	payload := dto.HistoryAuthCreateDTO{
		UserID:        authEntity.UserID,
		ActionType:    actionType,
		ActionName:    actionName,
		Description:   description,
		Success:       &success,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		DeviceID:      deviceID,
		Metadata:      mergedMetadata,
		SourceService: "auth-service",
	}

	s.dispatchHistoryAuth(ctx, payload, fmt.Sprintf("failed to log %s history", actionName))
}

func (s *PINUsecase) dispatchHistoryAuth(ctx context.Context, payload dto.HistoryAuthCreateDTO, logPrefix string) {
	if s.notificationProvider == nil {
		return
	}

	if payload.SourceService == "" {
		payload.SourceService = "auth-service"
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.notificationProvider.CreateHistoryAuth(timeoutCtx, &payload); err != nil {
		slog.InfoContext(ctx, fmt.Sprintf("%s: %v", logPrefix, err))
	}
}

// createSession tạo session mới cho user
func (s *PINUsecase) createSession(c context.Context, authEntity *auth.AuthMethod, request dto.VerifyPINRequest) *auth.UserSessionEntity {
	userAgent := _utils.GetUserAgentFromContext(c)
	ip := _utils.AuditIPFromContext(c)
	deviceID := _utils.GetDeviceIdFromContext(c)
	if deviceID == "" {
		deviceID = request.DeviceID
	}
	if deviceID == "" {
		deviceID = request.DeviceName
	}

	now := time.Now()

	return &auth.UserSessionEntity{
		AuthID:     authEntity.ID,
		DeviceID:   deviceID,
		Version:    request.Version,
		Platform:   request.Platform,
		OS:         request.OS,
		DeviceName: request.DeviceName,
		UserAgent:  userAgent,
		IPRequest:  ip,
		LastLogin:  &now,
		Activate:   false,
		LogoutAt:   nil,
	}
}

func (s *PINUsecase) buildAuthParam(c context.Context, authEntity *auth.AuthMethod, otp string) (*dto.AuthParam, error) {
	if s.otpRepo == nil || s.otpUsecase == nil {
		return nil, fmt.Errorf("OTP service is not configured")
	}

	statusEntity, err := s.statusRepo.GetByID(c, authEntity.ID)
	if err != nil {
		return nil, _errors.ReturnError(service.UserStatusNotFound)
	}

	otpEntity, err := s.otpRepo.GetByID(c, authEntity.ID)
	if err != nil {
		return nil, _errors.ReturnError(service.ValidOTPNotFound)
	}

	return &dto.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: otp,
		Auth:    authEntity,
	}, nil
}
