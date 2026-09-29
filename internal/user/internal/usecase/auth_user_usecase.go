package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"user/internal"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	_errors "common/errors"
	_jwt "common/jwt"
	_utils "common/utils"
	shared_enum "pb/enums"
	"user/config"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/factory"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
)

// AuthUsecase is User Service's access-token authority.
// @bind: user/internal/usecase.AccessTokenIssuer
type AuthUsecase struct {
	CookieProvider     providers.CookieProvider
	CacheProvider      providers.CacheProvider
	PreePlanID         *uint64
	ProfileProvider    providers.ProfileProvider
	AuthMethodRepo     repo.AuthMethodRepository
	UserInfoRepo       repo.UserInfoRepository
	SessionRepo        repo.SessionRepository
	OTPRepo            repo.OTPRepository
	StatusRepo         repo.StatusRepository
	DeviceRepo         repo.DeviceRepository
	DeviceUsecase      *DeviceUsecase
	OtpUsecase         *OtpUsecase
	OrganizationClient organizationMembershipProvider
	NotificationClient providers.NotificationProvider
	ChatProvider       _provider.ChatProvider
	ZnsProvider        providers.IZnsProvider
	SMSProvider        providers.ISmsProvider
	properties         dto.PropertiesDTO
	CrmProvider        _provider.CrmProvider
}

// organizationMembershipProvider là dependency tối thiểu mà auth flow cần
// từ Organization Service. User Service không tự suy diễn membership từ
// organization_id do client gửi lên.
type organizationMembershipProvider interface {
	GetOrganizationMember(ctx context.Context, organizationID uint64, profileID uint64) (*dto.OrganizationMember, error)
}

const (
	organizationMemberStatusActive uint32 = 10
	// CRM origin is compatibility enrichment, not login authority. A missing CRM
	// must not consume the entire HTTP deadline or prevent token issuance.
	optionalCRMOriginTimeout = 750 * time.Millisecond
	// Login history is best-effort until it has a durable local outbox. Keep the
	// synchronous compatibility call bounded and never delay authentication by
	// the remote provider timeout.
	loginHistoryTimeout = time.Second
)

func NewAuthService(cookieProvider providers.CookieProvider,
	cacheProvider providers.CacheProvider,
	profileProvider providers.ProfileProvider,
	authMethodRepo repo.AuthMethodRepository,
	userInfoRepo repo.UserInfoRepository,
	sessionRepo repo.SessionRepository,
	otpRepo repo.OTPRepository,
	statusRepo repo.StatusRepository,
	deviceRepo repo.DeviceRepository,
	deviceUsecase *DeviceUsecase,
	otpUsecase *OtpUsecase,
	organizationClient providers.OrganizationProvider,
	notificationClient providers.NotificationProvider,
	chatProvider _provider.ChatProvider,
	znsProvider providers.IZnsProvider,
	smsProvider providers.ISmsProvider,
	properties factory.IFactory,
	crmProvider _provider.CrmProvider,
) *AuthUsecase {
	planId := uint64(1)
	return &AuthUsecase{
		CookieProvider:     cookieProvider,
		CacheProvider:      cacheProvider,
		PreePlanID:         &planId,
		ProfileProvider:    profileProvider,
		AuthMethodRepo:     authMethodRepo,
		UserInfoRepo:       userInfoRepo,
		SessionRepo:        sessionRepo,
		OTPRepo:            otpRepo,
		StatusRepo:         statusRepo,
		DeviceRepo:         deviceRepo,
		DeviceUsecase:      deviceUsecase,
		OtpUsecase:         otpUsecase,
		OrganizationClient: organizationClient,
		NotificationClient: notificationClient,
		ChatProvider:       chatProvider,
		ZnsProvider:        znsProvider,
		SMSProvider:        smsProvider,
		properties:         properties.GetProperties(),
		CrmProvider:        crmProvider,
	}
}

func (s *AuthUsecase) RequestOtp(c context.Context, otp dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	if otp.Mode != "auth" && otp.Mode != "verify" {
		return nil, _errors.ReturnError(service.AuthModeInvalid)
	}

	var err error
	var oauth *dto.RequestLoginResponse
	if otp.Mode == "auth" {
		oauth, err = s._requestAuth(c, otp)
	} else if otp.Mode == "verify" {
		oauth, err = s._requestVerify(c, otp)
	}

	if err != nil {
		return nil, err
	}

	return oauth, nil
}

func (s *AuthUsecase) _requestAuth(c context.Context, otp dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	valid := _utils.ValidatePhoneNumber(otp.Phone)
	if !valid {
		return nil, _errors.ReturnError(service.PhoneInvalid)
	}
	existed := false
	var smsChannel string
	oauth, err := s.AuthMethodRepo.FindByPhone(c, otp.Phone)
	if err == nil && oauth != nil {
		// User đã tồn tại - gửi OTP

		param, _ := s.GetParamValidate(c, oauth, "")
		err = s.OtpUsecase.ValidateAuth(c,
			[]enums.AuthCodeEnum{
				enums.LOCKED,
				enums.LIMIT_REQUEST_TIME,
			}, param)
		if err != nil {
			return nil, err
		}
		existed = oauth.UserID != 0
		_, smsChannel = s.SendNewOTP(c, param.OTP, otp.Phone)
	} else {
		// Bắt buộc phải có họ tên khi đăng ký
		// if otp.Fullname == "" {
		// 	return nil, _errors.ReturnError(service.FullNameRequiredForRegistration)
		// }
		oauth, smsChannel, err = s.RequestIsRegister(c, otp)
		if err != nil {
			return nil, err
		}
	}

	return &dto.RequestLoginResponse{
		AuthID:     oauth.ID,
		Provider:   "PHONE",
		OAuthID:    otp.Phone,
		Existed:    existed,
		SmsChannel: smsChannel,
	}, nil
}

func (s *AuthUsecase) _requestVerify(c context.Context, otp dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	if otp.AuthId == 0 {
		return nil, _errors.ReturnError(service.AuthIDRequired)
	}
	var smsChannel string
	oauth, err := s.AuthMethodRepo.FindByID(c, otp.AuthId)
	if err != nil {
		return nil, err
	}
	if oauth == nil {
		return nil, _errors.ReturnError(service.AuthIDInvalid)
	}

	// User đã tồn tại - gửi OTP
	param, _ := s.GetParamValidate(c, oauth, "")
	err = s.OtpUsecase.ValidateAuth(c,
		[]enums.AuthCodeEnum{
			enums.LOCKED,
			enums.LIMIT_REQUEST_TIME,
		}, param)
	if err != nil {
		return nil, err
	}
	_, smsChannel = s.SendNewOTP(c, param.OTP, oauth.AuthName)

	return &dto.RequestLoginResponse{
		AuthID:     oauth.ID,
		Provider:   "PHONE",
		OAuthID:    oauth.AuthName,
		Existed:    oauth.UserID != 0,
		SmsChannel: smsChannel,
	}, nil
}

func (s *AuthUsecase) ResendOTP(c context.Context, otp dto.OtpRequest) (*dto.ResendOTPResponse, error) {
	if otp.AuthId == 0 && otp.Phone == "" {
		return nil, _errors.ReturnError(service.PhoneOrIDRequired)
	}
	var oauth *auth.AuthMethod
	var err error
	if otp.AuthId != 0 {
		oauth, err = s.AuthMethodRepo.FindByID(c, otp.AuthId)
	} else {
		oauth, err = s.AuthMethodRepo.FindByPhone(c, otp.Phone)
	}
	if err != nil {
		return nil, err
	}
	param, err := s.GetParamValidate(c, oauth, "")
	if err != nil {
		return nil, err
	}
	err = s.OtpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED, // check lock trước next time
		enums.LIMIT_NEXT_TIME,
	}, param)
	if err != nil {
		return nil, err
	}

	_, smsChannel := s.SendNewOTP(c, param.OTP, otp.Phone)
	return &dto.ResendOTPResponse{
		AuthId:     oauth.ID,
		Message:    "OTP đã được gửi",
		SmsChannel: smsChannel,
	}, nil
}

func (s *AuthUsecase) VerifyOtp(c context.Context, otp dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	e, err := s.AuthMethodRepo.FindByPhone(c, otp.Phone)
	if err != nil {
		return nil, err
	}
	param, _ := s.GetParamValidate(c, e, otp.Otp)
	err = s.OtpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED,
		enums.CHECK_OTP,
	}, param)
	if err != nil {
		return nil, err
	}

	param.OTP.OTPSendTime = 0
	param.OTP.OTPCheckTime = 0
	param.OTP.Activate = true
	s.OTPRepo.UpdateOTP(c, param.OTP)

	param.Status.LockedUntil = time.Time{}
	param.Status.Verified = true
	s.StatusRepo.UpdateStatus(c, param.Status)

	newSession := s.createSession(c, e, otp)
	s.SessionRepo.CreateSession(c, newSession)

	infoEntity, err := s.ProfileProvider.GetByProfileID(c, e.UserID)
	isNewUser := false
	if err != nil || infoEntity == nil || infoEntity.ProfileID == 0 {
		// Cho phép đăng ký mà không cần họ tên
		// User có thể cập nhật họ tên sau khi đăng ký
		fullname := otp.Fullname
		if fullname == "" {
			fullname = "" // Để trống nếu user không nhập
		}

		// infoEntity = s.ProfileProvider.MakeUserProfileEntity(e.ID, s.PreePlanID, e.FullName, e.Email, e.Phone, e.Avatar)
		infoEntity, err = s.ProfileProvider.CreateProfile(c, e.ID, s.PreePlanID, fullname, e.Email, e.AuthName, e.Avatar)
		if err != nil {
			return nil, err
		}

		e.UserID = infoEntity.ProfileID
		s.AuthMethodRepo.Update(c, e)
		isNewUser = true

		// Tạo nhóm chat mặc định "Hỗ trợ khách hàng" cho user mới
		if s.ChatProvider != nil && len(s.properties.CSKH.TeamIDs) > 0 {
			func() {
				// Tạo context mới cho goroutine để tránh context bị cancel
				ctx := c
				// Set profileId vào context để chat service có thể xác định người tạo
				ctx = context.WithValue(ctx, _enum.ProfileIDKey, infoEntity.ProfileID)

				conversationID, err := s.ChatProvider.CreateConversation(
					ctx,
					"Hỗ trợ khách hàng",
					infoEntity.ProfileID,
					s.properties.CSKH.TeamIDs,
				)
				if err != nil {
					slog.ErrorContext(c, fmt.Sprintf("Failed to create default support chat group for user %d: %v", infoEntity.ProfileID, err))
				} else {
					slog.InfoContext(c, fmt.Sprintf("Created default support chat group (ID: %d) for new user %d", conversationID, infoEntity.ProfileID))
				}
			}()
		}
	} else {
		// User đã có profile - kiểm tra và update authId với profileId nếu chưa có
		if e.UserID == 0 && infoEntity.ProfileID != 0 {
			e.UserID = infoEntity.ProfileID
			s.AuthMethodRepo.Update(c, e)
		}
	}

	// Check account status - nếu tài khoản bị khóa thì logout và báo lỗi
	if e.UserID != 0 {
		_, canLogin, err := s.ProfileProvider.CheckAccountStatus(c, e.UserID)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to check account status: %v", err))
		}
		if !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(c, e.ID)
			return nil, _errors.ReturnError(service.AccountLocked)
		}
	}

	// wallet, err := s.PaymentClient.GetWallet(c.Request.Context(), &paymentpb.GetWalletRequest{
	// 	Id: uint32(infoEntity.ProfileID),
	// })

	// if err != nil {
	// 	return nil, &_routes.Except{
	// 		Code:    500,
	// 		Message: "Lỗi lấy ví",
	// 	}
	// }
	// if wallet == nil || wallet.Id == 0 {
	// 	_, err := s.PaymentClient.CreateWallet(c.Request.Context(), &paymentpb.CreateWalletRequest{
	// 		UserId: uint32(infoEntity.ProfileID),
	// 		Currency: "VND",
	// 		Balance: 0,
	// 	})
	// 	if err != nil {
	// 		return nil, &_routes.Except{
	// 			Code:    500,
	// 			Message: "Lỗi khi tạo ví",
	// 		}
	// 	}
	// }
	// Liên kết tài khoản OAuth với hồ sơ người dùng
	if otp.AuthID != 0 {
		var oauthAcc *auth.AuthMethod
		if oauthAcc, err = s.AuthMethodRepo.FindByID(c, otp.AuthID); err == nil {
			if oauthAcc != nil && oauthAcc.UserID == 0 {
				// var statusOauth *auth.UserStatusEntity
				if statusOauth, err := s.StatusRepo.GetByID(c, oauthAcc.ID); err == nil {
					statusOauth.Verified = true
					s.StatusRepo.UpdateStatus(c, statusOauth)
				}

				oauthAcc.UserID = infoEntity.ProfileID
				s.AuthMethodRepo.Update(c, oauthAcc)

				infoEntity.Avatar = oauthAcc.Avatar
				infoEntity.Email = oauthAcc.Email
				s.ProfileProvider.UpdateProfile(c, infoEntity)
			}
		}
	}

	// s.RedisService.SaveToken(e.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), property.AppProperties.Jwt.RefreshExpMinutes*60)
	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	refreshToken := s.GenRefreshToken(c, e, newSession.SessionID)

	// Gửi thông báo chào mừng nếu là user mới
	if isNewUser {
		func() {
			contextTimeout, cancel := context.WithTimeout(c, 10*time.Second)
			defer cancel()

			err := s.NotificationClient.CreateNotification(
				contextTimeout,
				infoEntity.Avatar,
				"Chào mừng bạn đến với BDSPro",
				[]string{"Thiết lập tài khoản thành công"},
				_enum.NotificationSystem,
				&infoEntity.ProfileID,
				infoEntity.ProfileID,
				_enum.EOwnerOfMember,
				[]string{},
			)
			if err != nil {
				slog.ErrorContext(c, fmt.Sprintf("Failed to send welcome notification: %v", err))
			} else {
				slog.InfoContext(c, fmt.Sprintf("Successfully sent welcome notification to user %d", infoEntity.ProfileID))
			}
		}()
	}

	// Lấy passwordAuthId (auth_method có provider ADMIN)
	response := s.ResponseLogin(c, e, infoEntity, newSession.SessionID)
	response.RefreshToken = refreshToken            // Set refreshToken vào response body
	response.ReferralCode = infoEntity.ReferralCode // Set mã giới thiệu vào response
	if infoEntity.ProfileID != 0 {
		passwordAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(c, infoEntity.ProfileID, "ADMIN")
		if err == nil && passwordAuth != nil {
			response.PasswordAuthId = passwordAuth.ID
		}

		// Lưu app state = "active" khi login thành công
		if err := s.CacheProvider.SaveAppState(c, infoEntity.ProfileID, "active"); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to save app state for user %d: %v", infoEntity.ProfileID, err))
		}
	}

	// Diffie-Hellman key exchange: Mã hóa auth_key nếu client gửi public key
	if otp.ClientPublicKey != "" {
		slog.InfoContext(c, fmt.Sprintf("[VerifyOTP] Starting DH key exchange..."))

		// Decode client public key
		clientPublicKey, err := _utils.DecodePublicKey(otp.ClientPublicKey)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to decode client public key: %v", err))
		} else {

			// Decode server private key
			serverPrivateKey, err := _utils.DecodePrivateKey(e.PrivateKey)
			if err != nil {
				slog.ErrorContext(c, fmt.Sprintf("Failed to decode server private key: %v", err))
			} else {

				// Compute shared secret
				sharedSecret, err := _utils.ComputeSharedSecret(serverPrivateKey, clientPublicKey)
				if err != nil {
					slog.ErrorContext(c, fmt.Sprintf("Failed to compute shared secret: %v", err))
				} else {
					sharedSecretHex := sharedSecret.Text(16)

					// Encrypt auth_key with shared secret
					encryptedAuthKey, err := _utils.EncryptAuthKeyXOR(e.AuthKey, sharedSecretHex)
					if err != nil {
						slog.ErrorContext(c, fmt.Sprintf("Failed to encrypt auth key: %v", err))
					} else {
						response.EncryptedAuthKey = encryptedAuthKey
						response.ServerPublicKey = e.PublicKey
						slog.InfoContext(c, fmt.Sprintf("[VerifyOTP] DH key exchange completed successfully"))
					}
				}
			}
		}
	}

	// Ghi log admin history
	func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(c, loginHistoryTimeout)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Mobile-Client"
		}
		slog.InfoContext(c, fmt.Sprintf("Creating login history for user %d, IP: %s, UserAgent: %s",
			e.UserID, ipAddress, userAgent))

		success := true
		historyDeviceID := otp.DeviceID
		if historyDeviceID == "" {
			historyDeviceID = otp.DeviceName
		}
		deviceName := otp.DeviceName
		if deviceName == "" {
			deviceName = historyDeviceID
		}
		description := "Đăng nhập hệ thống qua OTP"
		if deviceName != "" {
			description = fmt.Sprintf("Đăng nhập hệ thống qua OTP từ %s", deviceName)
		}
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      e.UserID,
			ActionType:  "login",
			ActionName:  "OTP_LOGIN",
			Description: description,
			Success:     &success,
			SessionID:   strconv.FormatUint(newSession.SessionID, 10),
			Channel:     otp.Platform,
			DeviceID:    historyDeviceID,
			Metadata: map[string]interface{}{
				"authId":       e.ID,
				"loginChannel": "otp",
			},
		}
		if deviceName != "" {
			historyPayload.Metadata["deviceName"] = deviceName
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(c, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log auth history (OTP): %v", err))
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			e.UserID,                             // adminID
			e.UserID,                             // targetID (chính user đó)
			int32(shared_enum.TargetHistoryLead), // targetType
			int32(shared_enum.HistoryCreateStep), // actionType
			"Đăng nhập hệ thống bằng OTP",        // title
			[]string{"Đăng nhập thành công vào hệ thống"}, // notes
			"Chưa đăng nhập",                    // preStage
			"Đã đăng nhập",                      // afterStage
			&e.UserID,                           // ownerID
			int32(shared_enum.EOwnerTypeMember), // ownerType
			enums.RoleMap[infoEntity.RoleType],  // adminRole
			ipAddress,                           // ipAddress
			userAgent,                           // userAgent
		)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log login history: %v", err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("Successfully logged login history for user %d", e.UserID))
		}
	}()

	// Cập nhật device với authId và profileId nếu có thông tin thiết bị
	if s.DeviceUsecase != nil {
		deviceID := _utils.GetDeviceIdFromContext(c)
		print("deviceID", deviceID)
		if deviceID != "" {
			func() {
				// Tạo context mới với authId và profileId
				deviceCtx := context.WithValue(c, _enum.AuthIDKey, e.ID)
				deviceCtx = context.WithValue(deviceCtx, _enum.ProfileIDKey, infoEntity.ProfileID)
				// if infoEntity.OrganizationID > 0 {
				// 	deviceCtx = context.WithValue(deviceCtx, _enum.OrganizationIDKey, infoEntity.OrganizationID)
				// }

				err := s.DeviceUsecase.UpdateDeviceWithAuthInfo(deviceCtx, deviceID, e.ID, infoEntity.ProfileID)
				if err != nil {
					slog.ErrorContext(c, fmt.Sprintf("Failed to update device with auth info after verify OTP: %v", err))
				}
			}()
		}
	}

	return response, nil
}

// VerifyRecovery xác minh OTP và gán số điện thoại vào tài khoản hiện tại (không tạo profile mới)
func (s *AuthUsecase) VerifyRecovery(c context.Context, otp dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	// 1. Lấy authID của user đang đăng nhập từ context
	currentAuthID := _utils.GetAuthIdFromContext(c)
	if currentAuthID == 0 {
		return nil, _errors.ReturnError(service.LoginRequired)
	}

	// 2. Lấy thông tin auth_method của user hiện tại
	currentAuthMethod, err := s.AuthMethodRepo.FindByID(c, currentAuthID)
	if err != nil || currentAuthMethod == nil {
		return nil, _errors.ReturnError(service.AccountInfoNotFound)
	}
	if currentAuthMethod.UserID == 0 {
		return nil, _errors.ReturnError(service.CurrentAccountProfileMissing)
	}

	// 3. Verify OTP cho số điện thoại mới
	e, err := s.AuthMethodRepo.FindByPhone(c, otp.Phone)
	if err != nil {
		return nil, err
	}
	param, _ := s.GetParamValidate(c, e, otp.Otp)
	err = s.OtpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED,
		enums.CHECK_OTP,
	}, param)
	if err != nil {
		return nil, err
	}

	// Gán userID
	e.UserID = currentAuthMethod.UserID
	_, err = s.AuthMethodRepo.Update(c, e)
	if err != nil {
		return nil, fmt.Errorf("update account: %w", err)
	}

	// Audit is best-effort, but it remains request-owned and bounded.
	if s.NotificationClient != nil {
		auditCtx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Mobile-Client"
		}

		if err := s.NotificationClient.CreateAdminHistory(auditCtx,
			currentAuthMethod.UserID,
			currentAuthMethod.UserID,
			int32(shared_enum.TargetHistoryLead),
			int32(shared_enum.HistoryCreateStep),
			"Liên kết số điện thoại mới",
			[]string{fmt.Sprintf("Đã liên kết số điện thoại %s vào tài khoản", otp.Phone)},
			"",
			"Đã liên kết",
			&currentAuthMethod.UserID,
			int32(shared_enum.EOwnerTypeMember),
			"",
			ipAddress,
			userAgent,
		); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("failed to log phone recovery audit: %v", err))
		}
	}

	return &dto.AuthLoginResponse{
		AuthID:    currentAuthMethod.ID,
		ProfileID: currentAuthMethod.UserID,
	}, nil
}

// LoginWithKey đăng nhập bằng auth key (Diffie-Hellman)
func (s *AuthUsecase) LoginWithKey(c context.Context, req dto.LoginWithKeyRequest) (*dto.AuthLoginResponse, error) {
	// 1. Tìm auth method bằng phone
	e, err := s.AuthMethodRepo.FindByPhone(c, req.Phone)
	if err != nil {
		return nil, _errors.ReturnError(service.PhoneNotFound)
	}

	// Kiểm tra auth_key đã được khởi tạo chưa
	if e.AuthKey == "" || e.PrivateKey == "" || e.PublicKey == "" {
		return nil, _errors.ReturnError(service.AuthKeyNotInitialized)
	}
	slog.InfoContext(c, fmt.Sprintf("[LoginWithKey] Starting DH key exchange for phone: %s", req.Phone))

	// 2. Decode server private key
	serverPrivateKey, err := _utils.DecodePrivateKey(e.PrivateKey)
	if err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to decode server private key: %v", err))
		return nil, fmt.Errorf("authenticate request: %w", err)
	}

	// 3. Decode client public key
	clientPublicKey, err := _utils.DecodePublicKey(req.ClientPublicKey)
	if err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to decode client public key: %v", err))
		return nil, _errors.ReturnError(service.ClientPublicKeyInvalid)
	}

	// 4. Compute shared secret
	sharedSecret, err := _utils.ComputeSharedSecret(serverPrivateKey, clientPublicKey)
	if err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to compute shared secret: %v", err))
		return nil, fmt.Errorf("authenticate request: %w", err)
	}

	sharedSecretHex := sharedSecret.Text(16)

	// 5. Decrypt auth key từ client gửi lên
	decryptedAuthKey, err := _utils.DecryptAuthKeyXOR(req.EncryptedAuthKey, sharedSecretHex)
	if err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to decrypt auth key: %v", err))
		return nil, _errors.ReturnError(service.AuthKeyInvalid)
	}

	// 6. So sánh với auth_key trong database
	if decryptedAuthKey != e.AuthKey {
		return nil, _errors.ReturnError(service.AuthKeyMismatch)
	}
	slog.InfoContext(c, fmt.Sprintf("[LoginWithKey] Auth key validation successful!"))

	// 7. Kiểm tra account status
	if e.UserID != 0 {
		isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(c, e.UserID)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to check account status: %v", err))
		}
		if isLocked || !canLogin {
			return nil, _errors.ReturnError(service.AccountLockedContactAdmin)
		}
	}

	// 8. Tạo session
	deviceID := _utils.GetDeviceIdFromContext(c)
	if deviceID == "" {
		deviceID = req.DeviceID
	}
	if deviceID == "" {
		deviceID = req.DeviceName
	}
	newSession := &auth.UserSessionEntity{
		AuthID:     e.ID,
		DeviceID:   deviceID,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.OS,
		DeviceName: req.DeviceName,
		IPRequest:  _utils.AuditIPFromContext(c),
		UserAgent:  _utils.GetUserAgentFromContext(c),
	}
	s.SessionRepo.CreateSession(c, newSession)

	// 9. Lấy profile
	infoEntity, err := s.ProfileProvider.GetByProfileID(c, e.UserID)
	if err != nil || infoEntity == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	// 10. Generate tokens
	refreshToken := s.GenRefreshToken(c, e, newSession.SessionID)
	response := s.ResponseLogin(c, e, infoEntity, newSession.SessionID)
	response.RefreshToken = refreshToken
	response.ReferralCode = infoEntity.ReferralCode
	response.ServerPublicKey = ""

	// 11. Get password auth id
	if infoEntity.ProfileID != 0 {
		passwordAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(c, infoEntity.ProfileID, "ADMIN")
		if err == nil && passwordAuth != nil {
			response.PasswordAuthId = passwordAuth.ID
		}

		// Lưu app state = "active" khi login thành công
		if err := s.CacheProvider.SaveAppState(c, infoEntity.ProfileID, "active"); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to save app state for user %d: %v", infoEntity.ProfileID, err))
		}
	}

	// 12. Mã hóa auth_key mới với client public key mới (để client lưu cho lần login sau)
	// encryptedAuthKey, err := _utils.EncryptAuthKeyXOR(e.AuthKey, sharedSecretHex)
	// if err != nil {
	// 	log.Printf("Failed to encrypt auth key: %v", err)
	// } else {
	// 	log.Printf("[LoginWithKey] New EncryptedAuthKey: %s", encryptedAuthKey)
	// 	response.EncryptedAuthKey = encryptedAuthKey
	// 	response.ServerPublicKey = e.PublicKey
	// }

	// 13. Ghi log admin history
	func() {
		contextTimeout, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)

		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Mobile-Client"
		}
		slog.InfoContext(c, fmt.Sprintf("Creating login history for user %d (LoginWithKey), IP: %s, UserAgent: %s",
			e.UserID, ipAddress, userAgent))

		success := true
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      e.UserID,
			ActionType:  "login",
			ActionName:  "AUTH_KEY_LOGIN",
			Description: "Đăng nhập hệ thống bằng sinh trắc học",
			Success:     &success,
			SessionID:   strconv.FormatUint(newSession.SessionID, 10),
			Channel:     req.Platform,
			DeviceID:    deviceID,
			Metadata: map[string]interface{}{
				"authId":       e.ID,
				"loginChannel": "auth_key",
			},
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(c, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log auth key history: %v", err))
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			e.UserID,
			e.UserID,
			int32(shared_enum.TargetHistoryLead),
			int32(shared_enum.HistoryCreateStep),
			"Đăng nhập hệ thống bằng Auth Key",
			[]string{"Đăng nhập thành công vào hệ thống bằng Auth Key (Diffie-Hellman)"},
			"Chưa đăng nhập",
			"Đã đăng nhập",
			&e.UserID,
			int32(shared_enum.EOwnerTypeMember),
			enums.RoleMap[infoEntity.RoleType],
			ipAddress,
			userAgent,
		)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log login history: %v", err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("Successfully logged login history for user %d", e.UserID))
		}
	}()

	// Cập nhật device với authId và profileId nếu có thông tin thiết bị
	if s.DeviceUsecase != nil {
		deviceID := req.DeviceID
		if deviceID == "" {
			deviceID = req.DeviceName
		}
		if deviceID != "" && e.ID > 0 && infoEntity.ProfileID > 0 {
			func() {
				// Tạo context mới với authId và profileId
				deviceCtx := context.WithValue(c, _enum.AuthIDKey, e.ID)
				deviceCtx = context.WithValue(deviceCtx, _enum.ProfileIDKey, infoEntity.ProfileID)
				// if infoEntity.OrganizationID > 0 {
				// 	deviceCtx = context.WithValue(deviceCtx, _enum.OrganizationIDKey, infoEntity.OrganizationID)
				// }

				err := s.DeviceUsecase.UpdateDeviceWithAuthInfo(deviceCtx, deviceID, e.ID, infoEntity.ProfileID)
				if err != nil {
					slog.ErrorContext(c, fmt.Sprintf("Failed to update device with auth info after login with key: %v", err))
				}
			}()
		}
	}

	return response, nil
}

func (s *AuthUsecase) createSession(c context.Context, e *auth.AuthMethod, otp dto.OtpVerifyRequest) *auth.UserSessionEntity {
	deviceID := _utils.GetDeviceIdFromContext(c)
	if deviceID == "" {
		deviceID = otp.DeviceID
	}
	if deviceID == "" {
		deviceID = otp.DeviceName
	}
	now := time.Now()
	return &auth.UserSessionEntity{
		AuthID:     e.ID,
		DeviceID:   deviceID,
		Version:    otp.Version,
		Platform:   otp.Platform,
		OS:         otp.OS,
		DeviceName: otp.DeviceName,
		IPRequest:  _utils.AuditIPFromContext(c),
		UserAgent:  _utils.GetUserAgentFromContext(c),
		Activate:   false,
		LastLogin:  &now,
		LogoutAt:   nil,
	}
}

func (s *AuthUsecase) RefreshToken(c context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	claims := _jwt.GetPrincipalContext(c)
	props, err := _jwt.GetProperties(refreshToken)
	if err != nil {
		return nil, _errors.ReturnError(service.TokenInvalidOrExpired, _errors.WithCause(err))
	}

	if claims != nil && claims.OrganizationId != nil {
		props.OrganizationID = claims.OrganizationId
	}

	// Lấy timestamp từ cache
	// cachedTimestamp, err := s.CacheProvider.GetToken(c, props.ProfileID)
	// if err != nil {
	// 	log.Printf("⚠️ [RefreshToken] Không tìm thấy token trong cache cho authID %d: %v. Bỏ qua kiểm tra.", props.AuthID, err)
	// 	// Không tìm thấy trong cache thì bỏ qua kiểm tra, tiếp tục xử lý
	// } else {
	// 	// Validate token bằng cách so sánh Unix timestamp (UTC)
	// 	validToken := _utils.ValidTokenByIssueAt(props.IssuedAt, cachedTimestamp)
	// 	if !validToken {
	// 		log.Printf("⚠️ [RefreshToken] Token không hợp lệ cho authID %d. IssuedAt: %s, CachedTime: %s",
	// 			props.AuthID, props.IssuedAt.Format(time.RFC3339), cachedTimestamp)
	// 		return nil, &_routes.Except{
	// 			Code:    401,
	// 			Message: "Phiên làm việc đã hết hạn vui lòng đăng nhập lại",
	// 		}
	// 	}
	// }

	if props.Type == "REFRESH" {
		// Organization context trong access token là snapshot có thời hạn.
		// Khi refresh phải hỏi lại owner membership để không gia hạn
		// quyền dùng quota sau khi thành viên đã bị vô hiệu hoá.
		if props.OrganizationID != nil {
			if err := s.validateOrganizationMembership(c, props.ProfileID, *props.OrganizationID); err != nil {
				return nil, err
			}
		}

		// if props.SessionID != 0 {
		// 	sessionEntity, err := s.SessionRepo.GetBySessionID(c, props.SessionID)
		// 	if err != nil {
		// 		return nil, _errors.ReturnError(service.SessionInvalid)
		// 	}
		// 	if sessionEntity == nil || sessionEntity.LogoutAt != nil {
		// 		return nil, _errors.ReturnError(service.SessionEnded)
		// 	}
		// } else {
		// 	return nil, _errors.ReturnError(service.SessionInvalid)
		// }

		profile, err := s.ProfileProvider.GetByProfileID(c, props.ProfileID)
		if err != nil {
			return nil, err
		}

		// originID, err := s.CrmProvider.GetOriginID(c, props.ProfileID)
		// if err != nil {
		// 	return nil, err
		// }

		newAccessToken, err := s.GenerateAccessToken(c, props.AuthID, props.SessionID, props.OrganizationID)
		if err != nil {
			return nil, err
		}

		// Tạo refresh token mới
		roleIds, _ := s.UserInfoRepo.GetRoleIdsByProfileId(c, props.ProfileID)
		newRefreshToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID: props.AuthID,
			// OriginID:       originID,
			ProfileID:      props.ProfileID,
			Role:           props.Role,
			RoleIds:        roleIds,
			SessionID:      props.SessionID,
			OrganizationID: props.OrganizationID,
			PlanID:         profile.PlanID,
			PlanAt:         profile.PlanAt,
			Type:           "REFRESH",
		})
		slog.

			// Lưu timestamp mới vào cache để invalidate refresh token cũ
			// currentTime := time.Now().Format(time.RFC3339)
			// err = s.CacheProvider.SaveToken(c, props.ProfileID, currentTime, uint64(config.Properties.JWT.RefreshExpMinutes*60))
			// if err != nil {
			// 	log.Printf("⚠️ [RefreshToken] Không thể lưu timestamp vào cache cho authID %d: %v", props.AuthID, err)
			// }
			InfoContext(c, fmt.Sprintf("✅ [RefreshToken] Refresh token thành công cho authID %d", props.AuthID))
		response := &dto.RefreshResponse{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		}

		if props.ProfileID != 0 {
			metadata := map[string]interface{}{
				"authId":    props.AuthID,
				"sessionId": props.SessionID,
			}
			s.logHistoryAuth(c, &dto.HistoryAuthCreateDTO{
				UserID:      props.ProfileID,
				ActionType:  "token",
				ActionName:  "REFRESH_TOKEN",
				Description: "Làm mới phiên",
				Metadata:    metadata,
			})
		}

		return response, nil
	}

	return nil, _errors.ReturnError(service.TokenInvalidOrExpired, _errors.WithPublicMessage("Token không đúng"))
}

func (s *AuthUsecase) Logout(c context.Context, logoutType string, sessionDeviceId uint64) (*dto.LogoutResponse, error) {
	authID := _utils.GetAuthIdFromContext(c)
	if authID == 0 {
		return nil, _errors.ReturnError(service.SessionAuthenticationFailed)
	}

	sessionID := _utils.GetSessionIdFromContext(c)

	logoutType = strings.ToLower(logoutType)
	if logoutType == "" {
		logoutType = "session"
	}

	now := time.Now()
	var sessions []*auth.UserSessionEntity
	var err error
	var currentSession *auth.UserSessionEntity

	switch logoutType {
	case "others":
		sessions, err = s.SessionRepo.GetSessionsByAuthIDWithoutSessionId(c, authID, sessionID)
		if err != nil {
			return nil, err
		}
		if sessionID != 0 {
			currentSession, _ = s.SessionRepo.GetBySessionIDAndAuthID(c, sessionID, authID)
		}
	case "device":
		if sessionDeviceId == 0 {
			return nil, _errors.ReturnError(service.SessionIDRequiredForDeviceLogout)
		}
		session, err := s.SessionRepo.GetBySessionIDAndAuthID(c, sessionDeviceId, authID)
		if err != nil {
			return nil, err
		}
		if session != nil {
			sessions = []*auth.UserSessionEntity{session}
		}
	default: // "session"
		sessionID := _utils.GetSessionIdFromContext(c)
		if sessionID == 0 {
			return nil, _errors.ReturnError(service.SessionAuthenticationFailed)
		}
		session, err := s.SessionRepo.GetBySessionID(c, sessionID)
		if err != nil {
			return nil, err
		}
		if session != nil {
			sessions = []*auth.UserSessionEntity{session}
			currentSession = session
		}
	}

	if len(sessions) == 0 {
		return &dto.LogoutResponse{Message: "Không tìm thấy phiên đăng nhập để đăng xuất"}, nil
	}

	var lastSession *auth.UserSessionEntity
	for _, session := range sessions {
		session.FinishedDate = &now
		session.LogoutAt = &now
		if err := s.SessionRepo.UpdateSession(c, session); err != nil {
			return nil, err
		}
		lastSession = session
	}

	s.CacheProvider.DeleteToken(c, authID)

	if profileID := _utils.GetProfileIdWithContext(c); profileID != 0 && lastSession != nil {
		metadata := map[string]interface{}{
			"authId":     lastSession.AuthID,
			"sessionId":  lastSession.SessionID,
			"platform":   lastSession.Platform,
			"deviceName": lastSession.DeviceName,
			"os":         lastSession.OS,
			"ipAddress":  lastSession.IPRequest,
			"userAgent":  lastSession.UserAgent,
			"type":       logoutType,
		}
		if len(sessions) > 1 {
			metadata["sessionCount"] = len(sessions)
		}
		s.logHistoryAuth(c, &dto.HistoryAuthCreateDTO{
			UserID:      profileID,
			ActionType:  "logout",
			ActionName:  "LOGOUT",
			Description: s.buildLogoutDescription(logoutType, lastSession, currentSession, len(sessions)),
			Metadata:    metadata,
		})
	}

	// Xóa profileId khỏi tất cả devices của user
	if profileID := _utils.GetProfileIdWithContext(c); profileID > 0 {
		s.removeProfileIdFromDevices(c, profileID)

		// Lưu app state = "inactive" khi logout
		if err := s.CacheProvider.SaveAppState(c, profileID, "inactive"); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to save app state for user %d: %v", profileID, err))
		}
	}

	return &dto.LogoutResponse{Message: "Đăng xuất thành công"}, nil
}

func (s *AuthUsecase) buildLogoutDescription(logoutType string, targetSession *auth.UserSessionEntity, currentSession *auth.UserSessionEntity, affected int) string {
	targetDevice := formatSessionDevice(targetSession)

	switch logoutType {
	case "others":
		initiator := formatSessionDevice(currentSession)
		if initiator == "" {
			initiator = "thiết bị hiện tại"
		}
		if affected > 1 {
			return fmt.Sprintf("Đăng xuất %d thiết bị khác từ %s", affected, initiator)
		}
		if targetDevice != "" {
			return fmt.Sprintf("Đăng xuất thiết bị %s từ %s", targetDevice, initiator)
		}
		return fmt.Sprintf("Đăng xuất các thiết bị khác từ %s", initiator)
	case "device":
		if targetDevice != "" {
			return fmt.Sprintf("Đăng xuất khỏi hệ thống trên thiết bị %s", targetDevice)
		}
		return "Đăng xuất khỏi hệ thống trên một thiết bị"
	default:
		if targetDevice != "" {
			return fmt.Sprintf("Đăng xuất khỏi hệ thống trên thiết bị %s", targetDevice)
		}
		return "Đăng xuất khỏi hệ thống"
	}
}

func formatSessionDevice(session *auth.UserSessionEntity) string {
	if session == nil {
		return ""
	}

	deviceName := session.DeviceName
	if deviceName == "" {
		deviceName = session.DeviceID
	}
	if deviceName == "" {
		deviceName = session.Platform
	}
	if deviceName == "" {
		deviceName = "thiết bị không xác định"
	}

	var platformInfo []string
	if session.Platform != "" {
		platformInfo = append(platformInfo, session.Platform)
	}
	if session.OS != "" && session.OS != session.Platform {
		platformInfo = append(platformInfo, session.OS)
	}

	if len(platformInfo) > 0 {
		return fmt.Sprintf("%s (%s)", deviceName, strings.Join(platformInfo, " / "))
	}

	return deviceName
}

// Restore khôi phục tài khoản
func (s *AuthUsecase) Restore(c context.Context, query dto.RestorePhoneRequest) (*dto.RestoreAccountResponse, error) {
	existingAuth, err := s.AuthMethodRepo.FindByPhone(c, query.Phone)
	ok := existingAuth != nil
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, _errors.ReturnError(service.PhoneAlreadyUsed, _errors.WithPublicMessage("Số điện thoại đã liên kết với tài khoản khác. Vui lòng nhập số điện thoại khác"))
	}

	profileId := _utils.GetProfileIdWithContext(c)

	e := auth.AuthMethod{
		Provider: "PHONE",
		AuthName: query.Phone,
		UserID:   profileId,
	}

	createdAuth, err := s.AuthMethodRepo.Create(c, &e)
	if err != nil {
		return nil, err
	}

	param, _ := s.NewUser(c, createdAuth.ID)
	if err := s.OTPRepo.CreateOTP(c, param.OTP); err != nil {
		return nil, err
	}

	return &dto.RestoreAccountResponse{
		AuthId:   createdAuth.ID,
		Provider: createdAuth.Provider,
		AuthName: createdAuth.AuthName,
	}, nil
}

// RestoreDeletedAccount khôi phục tài khoản đã bị soft delete
func (s *AuthUsecase) RestoreDeletedAccount(c context.Context, req dto.RestoreDeletedAccountRequest) (*dto.RestoreDeletedAccountResponse, error) {
	var deletedAuth *auth.AuthMethod
	var err error

	// Tìm tài khoản đã bị xóa theo phone, email hoặc username
	if req.Phone != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByPhone(c, req.Phone)
	} else if req.Email != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByEmail(c, req.Email)
	} else if req.Username != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByUsername(c, req.Username)
	} else {
		return nil, _errors.ReturnError(service.RecoveryIdentifierRequired)
	}

	if err != nil {
		return nil, _errors.ReturnError(service.DeletedAccountNotFound, _errors.WithCause(err))
	}

	if deletedAuth == nil {
		return nil, _errors.ReturnError(service.DeletedAccountNotFound, _errors.WithCause(err))
	}

	// Kiểm tra xem tài khoản đã được khôi phục chưa
	existingAuth, err := s.AuthMethodRepo.FindByID(c, deletedAuth.ID)
	if err == nil && existingAuth != nil {
		return nil, _errors.ReturnError(service.AccountAlreadyRestored)
	}

	// Khôi phục tài khoản bằng cách set deleted_at = null
	err = s.AuthMethodRepo.RestoreAccount(c, deletedAuth.ID)
	if err != nil {
		return nil, fmt.Errorf("restore deleted account: %w", err)
	}

	return &dto.RestoreDeletedAccountResponse{
		AuthId:   deletedAuth.ID,
		Provider: deletedAuth.Provider,
		AuthName: deletedAuth.AuthName,
		Message:  "Tài khoản đã được khôi phục thành công",
	}, nil
}

// Xóa tài khoản
func (s *AuthUsecase) Delete(c context.Context, otp string) (*dto.DeleteAccountResponse, error) {
	if len(otp) < 6 {
		return nil, _errors.ReturnError(service.OTPLengthInvalid)
	}
	authId := _utils.GetAuthIdFromContext(c)
	e, err := s.AuthMethodRepo.FindByID(c, authId)
	if err != nil {
		return nil, err
	}
	param, _ := s.GetParamValidate(c, e, otp)
	err = s.OtpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED,
		enums.CHECK_OTP,
	}, param)
	if err != nil {
		return nil, err
	}

	// principal := _jwt.GetPrincipalContext(c)
	profileId := _utils.GetProfileIdWithContext(c)
	profile, err := s.ProfileProvider.GetByProfileID(c, profileId)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, _errors.ReturnError(service.AccountNotFound)
	}

	// Xóa vĩnh viễn profile trong user-service
	if err := s.ProfileProvider.HardDeleteProfile(c, profileId); err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to hard delete profile: %v", err))
		return nil, fmt.Errorf("hard delete user profile: %w", err)
	}

	// Xóa vĩnh viễn auth method trong auth-service
	if err := s.AuthMethodRepo.Delete(c, authId); err != nil {
		slog.ErrorContext(c, fmt.Sprintf("Failed to hard delete auth method: %v", err))
		return nil, fmt.Errorf("hard delete auth method: %w", err)
	}

	// Xóa session và cache
	s.CookieProvider.ResetAll(c)
	s.CacheProvider.DeleteToken(c, authId)

	// Xóa profileId khỏi tất cả devices của user
	if profileId > 0 {
		s.removeProfileIdFromDevices(c, profileId)
	}

	return &dto.DeleteAccountResponse{
		Message: "Tài khoản đã được xóa vĩnh viễn",
	}, nil
}

func (s *AuthUsecase) LockAccount(c context.Context, otp string) (*dto.LockAccountResponse, error) {
	if len(otp) < 6 {
		return nil, _errors.ReturnError(service.OTPLengthInvalid)
	}
	authId := _utils.GetAuthIdFromContext(c)
	e, err := s.AuthMethodRepo.FindByID(c, authId)
	if err != nil {
		return nil, err
	}
	param, _ := s.GetParamValidate(c, e, otp)
	err = s.OtpUsecase.ValidateAuth(c, []enums.AuthCodeEnum{
		enums.ACTIVATED,
		enums.EXPIRED,
		enums.LIMIT_OTP_SEND,
		enums.LIMIT_OTP_ENTER,
		enums.LOCKED,
		enums.CHECK_OTP,
	}, param)
	if err != nil {
		return nil, err
	}

	err = s.AuthMethodRepo.LockAccount(c, authId, "temporary", 1, "Khóa tài khoản", 0)
	if err != nil {
		return nil, err
	}

	err = s.ProfileProvider.LockAccount(c, e.UserID)
	if err != nil {
		return nil, err
	}

	metadata := map[string]interface{}{
		"authId":     authId,
		"lockType":   "temporary",
		"durationHr": 24,
	}
	s.logHistoryAuth(c, &dto.HistoryAuthCreateDTO{
		UserID:         e.UserID,
		ActionType:     "account",
		ActionName:     "LOCK_ACCOUNT",
		Description:    "Khóa tài khoản tạm thời theo yêu cầu người dùng",
		AdditionalNote: "User initiated lock account",
		Metadata:       metadata,
	})

	lockedAt := time.Now()
	lockedUntil := lockedAt.Add(time.Hour * 24)
	return &dto.LockAccountResponse{
		Success:     true,
		Message:     "Khóa tài khoản thành công",
		ProfileID:   e.UserID,
		LockType:    auth.AccountStatus(_enum.EUserStatusPermanentlyLocked),
		LockedAt:    _utils.FormatTimeToString(&lockedAt),
		LockedUntil: _utils.FormatTimeToString(&lockedUntil),
	}, nil
}

func (s *AuthUsecase) RequestIsRegister(c context.Context, otpReq dto.OtpRequest) (*auth.AuthMethod, string, error) {
	// Generate Diffie-Hellman key pair
	keyPair, err := _utils.GenerateDHKeyPair()
	if err != nil {
		return nil, "", fmt.Errorf("generate DH key pair: %w", err)
	}

	// Generate auth key
	authKey, err := _utils.GenerateAuthKey()
	if err != nil {
		return nil, "", fmt.Errorf("generate auth key: %w", err)
	}

	// Tạo một UserAuthEntity mới với DH keys và auth_key
	e := &auth.AuthMethod{
		Provider:   "PHONE",
		AuthName:   otpReq.Phone,
		FullName:   otpReq.Fullname,
		PrivateKey: _utils.EncodePrivateKey(keyPair.PrivateKey),
		PublicKey:  _utils.EncodePublicKey(keyPair.PublicKey),
		AuthKey:    authKey,
	}

	// Lưu vào database
	createdAuth, err := s.AuthMethodRepo.Create(c, e)
	if err != nil {
		return nil, "", err
	}

	// Tạo AuthParam mới
	param, err := s.NewUser(c, createdAuth.ID)
	if err != nil {
		return nil, "", err
	}

	// Gửi OTP trong updateNewOTP luôn
	_, smsChannel := s.SendNewOTP(c, param.OTP, otpReq.Phone)

	s.logHistoryAuth(c, &dto.HistoryAuthCreateDTO{
		UserID:      createdAuth.UserID,
		ActionType:  "register",
		ActionName:  "REQUEST_REGISTER",
		Description: "Khởi tạo yêu cầu đăng ký tài khoản",
		Metadata: map[string]interface{}{
			"authId": createdAuth.ID,
			"phone":  otpReq.Phone,
		},
	})

	return createdAuth, smsChannel, nil
}

func (s *AuthUsecase) NewUser(c context.Context, oauthId uint64) (*dto.AuthParam, error) {
	// Tạo UserStatusEntity mới
	statusEntity := &auth.UserStatusEntity{
		Active:   true,
		Verified: false,
		IdDomain: auth.IdDomain{
			AuthID: oauthId,
		},
	}

	// Lưu vào database
	if err := s.StatusRepo.CreateStatus(c, statusEntity); err != nil {
		return nil, err
	}

	// Tạo UserOTPEntity mới
	now := time.Now()
	otpEntity := &auth.UserOTPEntity{
		OTP: _utils.GenerateOTP(),
		IdDomain: auth.IdDomain{
			AuthID: oauthId,
		},
		// AuthID:       oauthId,
		OTPSendTime:  0,
		OTPCheckTime: 0,
		OTPDate:      &now, // Set thời gian tạo OTP
	}

	return &dto.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: "",
		Auth:    nil,
	}, nil
}

func (s *AuthUsecase) SendNewOTP(c context.Context, otpEntity *auth.UserOTPEntity, phone string) (otp string, smsChannel string) {
	ip := _utils.AuditIPFromContext(c)
	s.CheckOtpSpam(c, phone, ip)

	otp = _utils.GenerateOTP()
	otpEntity.OTPSendTime += 1
	otpEntity.OTP = otp
	otpEntity.Activate = false

	// Set thời gian gửi OTP (để validate LIMIT_NEXT_TIME)
	now := time.Now()
	otpEntity.OTPDate = &now

	// Set thời gian hết hạn
	t := time.Now().Add(time.Duration(s.properties.ExpiredAfterSec) * time.Second)
	otpEntity.ExpiredTime = &t

	s.OTPRepo.UpdateOTP(c, otpEntity)

	// Mặc định channel là "sms"
	smsChannel = "sms"

	// Gửi OTP qua ZNS
	if s.ZnsProvider != nil {
		result, err := s.ZnsProvider.SendOTPZNS(c, phone, otp)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("⚠️ [ZNS] Lỗi khi gửi OTP qua ZNS cho phone %s: %v", phone, err))
			// ZNS lỗi → smsChannel = "sms"
			smsChannel = "sms"
		} else {
			if success, ok := result["success"].(bool); ok && success {
				slog.InfoContext(c, fmt.Sprintf("✅ [ZNS] Đã gửi OTP thành công cho phone %s", phone))
				// ZNS thành công → smsChannel = "zalo"
				smsChannel = "zalo"
			} else {
				slog.InfoContext(c, fmt.Sprintf("⚠️ [ZNS] Gửi OTP thất bại cho phone %s: %v", phone, result["message"]))
				// ZNS thất bại → smsChannel = "sms"
				smsChannel = "sms"
			}
		}
	} else {
		slog.InfoContext(c, fmt.Sprintf("⚠️ [ZNS] ZnsProvider chưa được khởi tạo, không gửi OTP qua ZNS"))
		// Không có ZNS provider → smsChannel = "sms"
		smsChannel = "sms"
	}

	if smsChannel == "sms" && s.SMSProvider != nil {
		template := s.properties.SMSTemplate
		if template == "" {
			template = "Mã xác thực của bạn là {OTP}. Vui lòng không chia sẻ OTP cho bất kỳ ai."
		}
		message := strings.ReplaceAll(template, "{OTP}", otp)
		_, err := s.SMSProvider.SendSMS(c, phone, "BDS Pro", message)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("⚠️ [SMS] Lỗi khi gửi OTP qua SMS cho phone %s: %v", phone, err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("✅ [SMS] Đã gửi OTP thành công cho phone %s", phone))
			smsChannel = "sms"
		}
	}

	return otp, smsChannel
}

func (s *AuthUsecase) CheckOtpSpam(c context.Context, phoneNumber, ip string) error {
	redisKeyPhone := fmt.Sprintf("otp:request:%s", phoneNumber)
	redisKeyIp := fmt.Sprintf("otp:ip:%s", ip)

	// Kiểm tra giới hạn gửi OTP theo số điện thoại
	requestCountPhone, err := s.CookieProvider.Get(c, redisKeyPhone)
	if err != nil {
		return err
	}
	countPhone, _ := strconv.ParseInt(requestCountPhone, 10, 32)
	if countPhone >= s.properties.LimitOtpDevice {
		return _errors.ReturnError(service.OTPRequestLimited, _errors.WithPublicMessage("Bạn đã gửi quá nhiều OTP. Vui lòng thử lại sau."), _errors.WithLegacyCode(401))
	}

	// Kiểm tra giới hạn gửi OTP theo IP
	requestCountIp, err := s.CookieProvider.Get(c, redisKeyIp)
	countIp, _ := strconv.ParseInt(requestCountIp, 10, 32)
	if err != nil {
		return err
	}
	if countIp >= s.properties.LimitOtpIP {
		return _errors.ReturnError(service.OTPIPAddressRateLimited)
	}

	// Tăng số lần gửi OTP trong Redis
	if err := s.CacheProvider.Increment(c, redisKeyPhone); err != nil {
		return err
	}
	if err := s.CookieProvider.Expire(c, redisKeyPhone, s.properties.CooldownSeconds); err != nil {
		return err
	}

	if err := s.CacheProvider.Increment(c, redisKeyIp); err != nil {
		return err
	}
	if err := s.CookieProvider.Expire(c, redisKeyIp, s.properties.IpCooldownSeconds); err != nil {
		return err
	}

	return nil
}

// GetParamValidate kiểm tra thông tin xác thực và trả về data.AuthParam
func (s *AuthUsecase) GetParamValidate(c context.Context, oauth *auth.AuthMethod, otp string) (*dto.AuthParam, error) {
	if oauth == nil {
		return nil, _errors.ReturnError(service.RequestValidationFailed, _errors.WithPublicMessage("Thông tin không đúng, vui lòng kiểm tra lại"))
	}

	statusEntity, err := s.StatusRepo.GetByID(c, oauth.ID)
	if err != nil {
		return nil, _errors.ReturnError(service.UserStatusNotFound, _errors.WithCause(err))
	}

	otpEntity, err := s.OTPRepo.GetByID(c, oauth.ID)
	if err != nil {
		return nil, _errors.ReturnError(service.ValidOTPNotFound, _errors.WithPublicMessage("Không tìm thấy OTP"), _errors.WithCause(err))
	}

	return &dto.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: otp,
		Auth:    oauth,
	}, nil
}

func (s *AuthUsecase) ResponseLogin(ctx context.Context, authEntity *auth.AuthMethod, profile *dto.ProfileDTO, sessionID uint64) *dto.AuthLoginResponse {
	if authEntity.UserID == 0 {
		return &dto.AuthLoginResponse{
			AuthID:   authEntity.ID,
			FullName: authEntity.FullName,
		}
	}

	// Kiểm tra xem user đã có họ tên chưa
	requireFullname := profile.FullName == ""
	var accessToken string
	accessToken, err := s.GenerateAccessToken(ctx, authEntity.ID, sessionID, nil)
	if err != nil {
		return &dto.AuthLoginResponse{
			AuthID:   authEntity.ID,
			FullName: authEntity.FullName,
		}
	}

	// if requireFullname {
	// 	// User chưa có họ tên, tạo TEMP token
	// 	accessToken, _ = _jwt.GenerateToken(_jwt.JwtTokenProperties{
	// 		AuthID:    authEntity.ID,
	// 		ProfileID: authEntity.UserID,
	// 		SessionID: sessionID,
	// 		Role:      enums.RoleMap[profile.RoleType],
	// 		Type:      _jwt.TempToken, // Token tạm thời, hạn 10 phút
	// 		PlanID:    profile.PlanID,
	// 		PlanAt:    profile.PlanAt,
	// 	})
	// } else {
	// 	// User đã có họ tên, tạo ACCESS token bình thường
	// 	accessToken, _ = _jwt.GenerateToken(_jwt.JwtTokenProperties{
	// 		AuthID:    authEntity.ID,
	// 		ProfileID: authEntity.UserID,
	// 		SessionID: sessionID,
	// 		Role:      enums.RoleMap[profile.RoleType],
	// 		Type:      _jwt.AccessToken,
	// 		PlanID:    profile.PlanID,
	// 		PlanAt:    profile.PlanAt,
	// 	})
	// }

	response := dto.MakeAuthLoginResponse(accessToken, authEntity, enums.RoleMap[profile.RoleType])
	response.RequireFullname = requireFullname
	if requireFullname {
		response.TempToken = accessToken // Đặt vào field tempToken để frontend biết
		response.AccessToken = ""        // Xóa accessToken vì đây là TEMP token
	}

	return response
}

func (s *AuthUsecase) GetRefreshToken(ctx context.Context, authEntity *auth.AuthMethod, sessionID uint64) string {
	roleIds, _ := s.UserInfoRepo.GetRoleIdsByProfileId(ctx, authEntity.UserID)
	refreshToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
		SessionID: sessionID,
		Role:      "ROLE_USER",
		Type:      "REFRESH",
		RoleIds:   roleIds,
	})
	return refreshToken
}

func (s *AuthUsecase) CreateSession(c context.Context, e *auth.AuthMethod, otp *dto.OtpVerifyRequest) *auth.UserSessionEntity {
	var authId uint64
	if e != nil {
		authId = e.ID
	}
	if otp == nil {
		otp = &dto.OtpVerifyRequest{}
	}
	deviceID := _utils.GetDeviceIdFromContext(c)
	if deviceID == "" {
		deviceID = otp.DeviceID
	}
	if deviceID == "" {
		deviceID = otp.DeviceName
	}
	now := time.Now()
	return &auth.UserSessionEntity{
		AuthID:     authId,
		DeviceID:   deviceID,
		Version:    otp.Version,
		Platform:   otp.Platform,
		OS:         otp.OS,
		DeviceName: otp.DeviceName,
		IPRequest:  _utils.AuditIPFromContext(c),
		UserAgent:  _utils.GetUserAgentFromContext(c),
		Activate:   false,
		LastLogin:  &now,
		LogoutAt:   nil,
	}
}

// CreateAndSaveSession tạo session và lưu vào database
func (s *AuthUsecase) CreateAndSaveSession(c context.Context, e *auth.AuthMethod, otp *dto.OtpVerifyRequest) (*auth.UserSessionEntity, error) {
	// Tạo session object
	session := s.CreateSession(c, e, otp)

	// Lưu vào database
	err := s.SessionRepo.CreateSession(c, session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (s *AuthUsecase) CreateQRSession(c context.Context, e *dto.SessionQRRequest, sessionId uint64) (*dto.SessionQRResponse, error) {
	session, _ := s.SessionRepo.GetBySessionID(c, sessionId)
	if session == nil {
		session = s.CreateSession(c, nil, nil)
	} else {
		// copier.Copy(&session, &e)
		// session.Activate
		// e.SessionId = sessionId
		return &dto.SessionQRResponse{
			SessionKey: session.SessionKey,
			SessionId:  session.SessionID,
		}, nil
	}

	key := _utils.RandomString(24)
	session.SessionKey = key

	err := s.SessionRepo.CreateSession(c, session)
	if err != nil {
		return nil, fmt.Errorf("create login session: %w", err)
	}

	s.CookieProvider.Set(c, config.C_SESSION_ID, strconv.FormatUint(session.SessionID, 10))

	return &dto.SessionQRResponse{
		SessionKey: key,
		SessionId:  session.SessionID,
	}, nil
}

const maxConcurrentRequests = 1000 // Giới hạn số request đồng thời

var semaphore = make(chan struct{}, maxConcurrentRequests)

func (s *AuthUsecase) VerifyQRSession(c context.Context, e *dto.SessionQRConfirm, sessionId uint64) (*dto.AuthLoginResponse, error) {
	// Nếu semaphore đã đầy, từ chối request
	select {
	case semaphore <- struct{}{}:
		defer func() { <-semaphore }() // Giải phóng khi xử lý xong
	default:
		return nil, _errors.ReturnError(service.QRSessionCapacityExceeded)
	}

	sessionEntity, err := s.SessionRepo.GetBySessionID(c, sessionId)
	if err != nil {
		return nil, _errors.ReturnError(service.LoginSessionNotFound, _errors.WithCause(err), _errors.WithLegacyCode(400))
	}

	if sessionEntity.FinishedDate != nil ||
		sessionEntity.Activate ||
		sessionEntity.SessionKey != e.SessionKey {
		return nil, _errors.ReturnError(service.LoginSessionIncorrect)
	}

	sessionEntity.Activate = true
	now := time.Now()
	sessionEntity.LastLogin = &now
	sessionEntity.LogoutAt = nil
	s.SessionRepo.UpdateSession(c, sessionEntity)

	authEntity, err := s.AuthMethodRepo.FindByID(c, sessionEntity.AuthID)
	if err != nil {
		return nil, _errors.ReturnError(service.AuthenticationInfoNotFound, _errors.WithCause(err), _errors.WithLegacyCode(400))
	}

	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	refreshToken := s.GenRefreshToken(c, authEntity, sessionId)
	profile, err := s.ProfileProvider.GetByProfileID(c, authEntity.UserID)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound, _errors.WithCause(err))
	}

	// Ghi log admin history
	func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "QR-Client"
		}
		slog.InfoContext(c, fmt.Sprintf("Creating login history for user %d (QR), IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent))

		success := true
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      authEntity.UserID,
			ActionType:  "login",
			ActionName:  "QR_LOGIN",
			Description: "Đăng nhập hệ thống bằng QR Code",
			Success:     &success,
			SessionID:   strconv.FormatUint(sessionEntity.SessionID, 10),
			Channel:     sessionEntity.Platform,
			DeviceID:    sessionEntity.DeviceID,
			Metadata: map[string]interface{}{
				"authId":       authEntity.ID,
				"loginChannel": "qr_code",
			},
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(c, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log QR auth history: %v", err))
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			authEntity.UserID,                             // adminID
			authEntity.UserID,                             // targetID (chính user đó)
			int32(shared_enum.TargetHistoryLead),          // targetType
			int32(shared_enum.HistoryCreateStep),          // actionType
			"Đăng nhập hệ thống bằng QR Code",             // title
			[]string{"Đăng nhập thành công vào hệ thống"}, // notes
			"Chưa đăng nhập",                              // preStage
			"Đã đăng nhập",                                // afterStage
			&authEntity.UserID,                            // ownerID
			int32(shared_enum.EOwnerTypeMember),           // ownerType
			enums.RoleMap[profile.RoleType],               // adminRole
			ipAddress,                                     // ipAddress
			userAgent,                                     // userAgent
		)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log QR login history: %v", err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("Successfully logged QR login history for user %d", authEntity.UserID))
		}
	}()

	response := s.ResponseLogin(c, authEntity, profile, sessionEntity.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}

// ------------------------------------------------------------

func (s *AuthUsecase) GenRefreshToken(c context.Context, authEntity *auth.AuthMethod, sessionId uint64) string {
	// Lưu Unix timestamp (UTC) vào cache để tránh vấn đề timezone
	// Trừ 5 giây để đảm bảo token mới được issue sau thời điểm này sẽ hợp lệ
	validFromTimestamp := time.Now().UTC().Add(-5 * time.Second).Unix()
	s.CacheProvider.SaveToken(c, authEntity.UserID, fmt.Sprintf("%d", validFromTimestamp), s.properties.JwtRefreshExpMinutes*60)

	refreshToken := s.GetRefreshToken(c, authEntity, sessionId)
	slog.
		// Không set cookie nữa, trả refreshToken trong response body
		// s.CookieProvider.SetRefreshToken(c, refreshToken)
		InfoContext(c, fmt.Sprintf("[GenRefreshToken] authID=%d validFrom=%d", authEntity.ID, validFromTimestamp))
	return refreshToken
}

// LogLoginHistory ghi lại lịch sử đăng nhập
func (s *AuthUsecase) LogLoginHistory(c context.Context, authEntity *auth.AuthMethod, profile *dto.ProfileDTO, title string) {
	func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Client"
		}

		// Kiểm tra profile có tồn tại không
		if profile == nil {
			slog.ErrorContext(c, fmt.Sprintf("Cannot log login history: profile is nil for user %d", authEntity.UserID))
			return
		}
		slog.InfoContext(c, fmt.Sprintf("Creating login history for user %d, IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent))

		success := true
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      authEntity.UserID,
			ActionType:  "login",
			ActionName:  title,
			Description: "Đăng nhập hệ thống",
			Success:     &success,
			Metadata: map[string]interface{}{
				"authId":       authEntity.ID,
				"loginChannel": "unknown",
			},
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(c, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log auth history (generic): %v", err))
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			authEntity.UserID,                    // adminID
			authEntity.UserID,                    // targetID (chính user đó)
			int32(shared_enum.TargetHistoryLead), // targetType
			int32(shared_enum.HistoryCreateStep), // actionType
			title,                                // title
			[]string{"Đăng nhập thành công vào hệ thống"}, // notes
			"Chưa đăng nhập",                    // preStage
			"Đã đăng nhập",                      // afterStage
			&authEntity.UserID,                  // ownerID
			int32(shared_enum.EOwnerTypeMember), // ownerType
			enums.RoleMap[profile.RoleType],     // adminRole
			ipAddress,                           // ipAddress
			userAgent,                           // userAgent
		)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log login history: %v", err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("Successfully logged login history for user %d", authEntity.UserID))
		}
	}()
}

func (s *AuthUsecase) prepareHistoryAuthPayload(ctx context.Context, payload *dto.HistoryAuthCreateDTO) dto.HistoryAuthCreateDTO {
	var result dto.HistoryAuthCreateDTO
	if payload != nil {
		result = *payload
	}

	if result.IPAddress == "" {
		result.IPAddress = _utils.AuditIPFromContext(ctx)
	}

	if result.UserAgent == "" {
		result.UserAgent = _utils.GetUserAgentFromContext(ctx)
	}

	if result.OrganizationID == nil {
		if organizationID := _utils.GetOrganizationIdFromContext(ctx); organizationID != 0 {
			result.OrganizationID = &organizationID
		}
	}

	if result.PerformedBy == nil {
		if performer := _utils.GetProfileIdWithContext(ctx); performer != 0 {
			result.PerformedBy = &performer
		}
	}

	if result.SessionID == "" {
		if sessionID := _utils.GetSessionIdFromContext(ctx); sessionID != 0 {
			result.SessionID = strconv.FormatUint(sessionID, 10)
		}
	}

	if result.SourceService == "" {
		result.SourceService = "auth-service"
	}

	if result.Success == nil {
		success := true
		result.Success = &success
	}

	if payload != nil && payload.Metadata != nil {
		metadataCopy := make(map[string]interface{}, len(payload.Metadata))
		for k, v := range payload.Metadata {
			metadataCopy[k] = v
		}
		result.Metadata = metadataCopy
	}

	return result
}

func (s *AuthUsecase) logHistoryAuth(ctx context.Context, payload *dto.HistoryAuthCreateDTO) {
	if s.NotificationClient == nil || payload == nil {
		return
	}

	prepared := s.prepareHistoryAuthPayload(ctx, payload)

	func(p dto.HistoryAuthCreateDTO) {
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.NotificationClient.CreateHistoryAuth(timeoutCtx, &p); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to create auth history: %v", err))
		}
	}(prepared)
}

func (s *AuthUsecase) SwitchOrganization(c context.Context, organizationID uint64) (*dto.AuthLoginResponse, error) {
	authId := _utils.GetAuthIdFromContext(c)
	sessionId := _utils.GetSessionIdFromContext(c)
	profileId := _utils.GetProfileIdWithContext(c)
	if profileId == 0 {
		return nil, _errors.ReturnError(service.SessionAuthenticationFailed, _errors.WithLegacyCode(400))
	}
	if organizationID == 0 {
		return nil, _errors.ReturnError(service.OrganizationIDRequired, _errors.WithPublicMessage("Không thể chuyển đổi tổ chức"))
	}
	if err := s.validateOrganizationMembership(c, profileId, organizationID); err != nil {
		return nil, err
	}

	profile, err := s.ProfileProvider.GetByProfileID(c, profileId)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound, _errors.WithCause(err))
	}

	// payload := _jwt.JwtTokenProperties{
	// 	AuthID:    authId,
	// 	ProfileID: profileId,
	// 	SessionID: sessionId,
	// 	Role:      "ROLE_USER",
	// 	Type:      "ACCESS",
	// 	PlanID:    profile.PlanID,
	// 	PlanAt:    profile.PlanAt,
	// }
	// payload.AdditionalInfo = map[string]interface{}{
	// 	"organization": clientResponse.OrganizationId,
	// }
	// newAccessToken, err := _jwt.GenerateToken(payload)
	// if err != nil {
	// 	return nil, err
	// }
	newAccessToken, err := s.GenerateAccessToken(c, authId, sessionId, &organizationID)
	if err != nil {
		return nil, err
	}
	response := &dto.AuthLoginResponse{
		AuthID:      authId,
		AccessToken: newAccessToken,
		Role:        "ROLE_USER",
		Username:    profile.FullName,
		ProfileID:   profile.ProfileID,
		FullName:    profile.FullName,
		// Avatar:         profile.Avatar,
		// Email:          profile.Email,
		// Phone:          profile.Phone,
		OrganizationID: organizationID,
		ReferralCode:   profile.ReferralCode, // Mã giới thiệu
	}

	return response, nil
}

// validateOrganizationMembership giữ một owner duy nhất cho quy tắc phát
// token có organization context: profile phải là thành viên active theo
// dữ liệu của Organization Service. Role/permission chi tiết vẫn do
// Organization Service quyết định; auth flow không tự sao chép catalog đó.
func (s *AuthUsecase) validateOrganizationMembership(ctx context.Context, profileID, organizationID uint64) error {
	if profileID == 0 || organizationID == 0 || s.OrganizationClient == nil {
		return _errors.ReturnError(service.OrganizationMembershipAuthenticationFailed)
	}

	member, err := s.OrganizationClient.GetOrganizationMember(ctx, organizationID, profileID)
	if err != nil {
		return err
	}
	if member == nil || member.OrganizationId != organizationID || member.UserId != profileID || member.Status != organizationMemberStatusActive {
		return _errors.ReturnError(service.OrganizationMembershipInactive)
	}
	return nil
}

// AdminLogin xử lý đăng nhập admin bằng username và password
func (s *AuthUsecase) AdminLogin(c context.Context, adminLogin dto.AdminLoginRequest) (*dto.AuthLoginResponse, error) {
	// Validate input
	if adminLogin.Username == "" || adminLogin.Password == "" {
		return nil, _errors.ReturnError(service.UsernamePasswordRequired)
	}

	// Tìm theo username (ADMIN); có thể có bản ghi legacy trùng auth_name → thử password từng candidate
	candidates, err := s.AuthMethodRepo.FindAllByAuthNameAndProvider(c, adminLogin.Username, "ADMIN")
	if err != nil || len(candidates) == 0 {
		candidates, err = s.AuthMethodRepo.FindAllByEmailAndProvider(c, adminLogin.Username, "ADMIN")
	}
	if err != nil || len(candidates) == 0 {
		return nil, _errors.ReturnError(service.UsernamePasswordInvalid)
	}

	var adminAuthEntity *auth.AuthMethod
	for _, candidate := range candidates {
		if candidate != nil && _utils.CheckPasswordHash(adminLogin.Password, candidate.Password) {
			adminAuthEntity = candidate
			break
		}
	}
	if adminAuthEntity == nil {
		return nil, _errors.ReturnError(service.UsernamePasswordInvalid)
	}

	// Lấy profile của admin
	adminProfile, err := s.ProfileProvider.GetByProfileID(c, adminAuthEntity.UserID)
	if err != nil || adminProfile == nil || adminProfile.ProfileID == 0 {
		// Tạo profile mới nếu chưa có
		adminProfile = s.ProfileProvider.MakeUserProfileEntity(
			adminAuthEntity.ID,
			s.PreePlanID,
			adminAuthEntity.FullName,
			adminAuthEntity.Email,
			adminAuthEntity.Phone,
			adminAuthEntity.Avatar,
		)
		s.ProfileProvider.SaveProfile(c, adminProfile)

		// Cập nhật UserID cho admin auth entity
		adminAuthEntity.UserID = adminProfile.ProfileID
		s.AuthMethodRepo.Update(c, adminAuthEntity)
	}

	// Check account status - nếu tài khoản bị khóa thì logout và báo lỗi
	if adminAuthEntity.UserID != 0 {
		isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(c, adminAuthEntity.UserID)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to check account status: %v", err))
		}
		if isLocked || !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(c, adminAuthEntity.ID)
			return nil, _errors.ReturnError(service.AccountLocked)
		}
	}

	// Tạo session cho admin
	adminDeviceID := _utils.GetDeviceIdFromContext(c)
	if adminDeviceID == "" {
		adminDeviceID = adminLogin.DeviceID
	}
	if adminDeviceID == "" {
		adminDeviceID = adminLogin.DeviceName
	}

	adminOtpRequest := dto.OtpVerifyRequest{
		OtpRequest: dto.OtpRequest{
			Phone:    adminAuthEntity.Phone,
			Fullname: adminAuthEntity.FullName,
		},
		Otp:        "",
		AuthID:     adminAuthEntity.ID,
		Version:    adminLogin.Version,
		Platform:   adminLogin.Platform,
		OS:         adminLogin.OS,
		DeviceName: adminLogin.DeviceName,
		DeviceID:   adminDeviceID,
		Fullname:   adminAuthEntity.FullName,
	}

	newSession := s.createSession(c, adminAuthEntity, adminOtpRequest)
	s.SessionRepo.CreateSession(c, newSession)

	// Tạo refresh token
	refreshToken := s.GenRefreshToken(c, adminAuthEntity, newSession.SessionID)

	adminProfile.RoleType = enums.ERoleAdmin
	// Tạo response với đầy đủ thông tin
	response := s.ResponseLogin(c, adminAuthEntity, adminProfile, newSession.SessionID)
	response.RefreshToken = refreshToken              // Set refreshToken vào response body
	response.ReferralCode = adminProfile.ReferralCode // Set mã giới thiệu vào response
	response.Role = enums.RoleMap[adminProfile.RoleType]
	response.Username = adminAuthEntity.AuthName

	// Ghi log admin history
	func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(c, loginHistoryTimeout)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(c)
		userAgent := _utils.GetUserAgentFromContext(c)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Admin-Client"
		}
		slog.InfoContext(c, fmt.Sprintf("Creating admin history for user %d, IP: %s, UserAgent: %s",
			adminAuthEntity.UserID, ipAddress, userAgent))

		success := true
		historyDeviceID := adminDeviceID
		if historyDeviceID == "" {
			historyDeviceID = adminLogin.DeviceName
		}
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      adminAuthEntity.UserID,
			ActionType:  "login",
			ActionName:  "ADMIN_LOGIN",
			Description: "Admin đăng nhập hệ thống",
			Success:     &success,
			SessionID:   strconv.FormatUint(newSession.SessionID, 10),
			Channel:     adminLogin.Platform,
			DeviceID:    historyDeviceID,
			Metadata: map[string]interface{}{
				"authId":       adminAuthEntity.ID,
				"loginChannel": "admin_password",
			},
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(c, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log admin auth history: %v", err))
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			adminAuthEntity.UserID,                        // adminID
			adminAuthEntity.UserID,                        // targetID (chính admin đó)
			int32(shared_enum.TargetHistoryLead),          // targetType
			int32(shared_enum.HistoryCreateStep),          // actionType
			"Admin đăng nhập hệ thống",                    // title
			[]string{"Đăng nhập thành công vào hệ thống"}, // notes
			"Chưa đăng nhập",                              // preStage
			"Đã đăng nhập",                                // afterStage
			&adminAuthEntity.UserID,                       // ownerID
			int32(shared_enum.EOwnerTypeMember),           // ownerType
			enums.RoleMap[adminProfile.RoleType],          // adminRole
			ipAddress,                                     // ipAddress
			userAgent,                                     // userAgent
		)
		if err != nil {
			slog.ErrorContext(c, fmt.Sprintf("Failed to log admin login: %v", err))
		} else {
			slog.InfoContext(c, fmt.Sprintf("Successfully logged admin login for user %d", adminAuthEntity.UserID))
		}
	}()

	return response, nil
}

// LockAccount khóa tài khoản (tạm thời hoặc vĩnh viễn)
// func (s *AuthUsecase) LockAccount(ctx context.Context, req *dto.LockAccountRequest) (*dto.LockAccountResponse, error) {
// 	// Validate duration cho khóa tạm thời
// 	if req.LockType == auth.StatusTemporarilyLocked && req.Duration <= 0 {
// 		return &dto.LockAccountResponse{
// 			Success: false,
// 			Message: "Thời gian khóa tạm thời phải lớn hơn 0",
// 		}, fmt.Errorf("invalid duration for temporary lock")
// 	}

// 	// Kiểm tra user đã bị khóa chưa
// 	userInfo, err := s.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
// 	if err != nil {
// 		return &dto.LockAccountResponse{
// 			Success: false,
// 			Message: "Lỗi khi kiểm tra trạng thái user: " + err.Error(),
// 		}, err
// 	}

// 	if userInfo.IsLocked() {
// 		return &dto.LockAccountResponse{
// 			Success: false,
// 			Message: "User đã bị khóa",
// 		}, fmt.Errorf("user already locked")
// 	}

// 	// Thực hiện khóa user
// 	err = s.UserInfoRepo.LockUser(ctx, req.ProfileID, req.LockType, req.Duration, req.Reason, req.LockedBy)
// 	if err != nil {
// 		return &dto.LockAccountResponse{
// 			Success: false,
// 			Message: "Lỗi khi khóa user: " + err.Error(),
// 		}, err
// 	}

// 	// Tạo response
// 	response := &dto.LockAccountResponse{
// 		Success:   true,
// 		Message:   "Khóa user thành công",
// 		ProfileID: req.ProfileID,
// 		LockType:  req.LockType,
// 		LockedAt:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
// 	}

// 	// Thêm thông tin thời gian khóa đến cho khóa tạm thời
// 	if req.LockType == auth.StatusTemporarilyLocked {
// 		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
// 		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
// 	}

// 	return response, nil
// }

// UnlockAccount mở khóa tài khoản
func (s *AuthUsecase) UnlockAccount(ctx context.Context, req *dto.UnlockAccountRequest) (*dto.UnlockAccountResponse, error) {
	// Kiểm tra user có bị khóa không
	userInfo, err := s.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "Lỗi khi kiểm tra trạng thái user: " + err.Error(),
		}, err
	}

	if !userInfo.IsLocked() {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "User không bị khóa",
		}, fmt.Errorf("user is not locked")
	}

	// Thực hiện mở khóa user
	err = s.UserInfoRepo.UnlockUser(ctx, req.ProfileID, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockAccountResponse{
			Success: false,
			Message: "Lỗi khi mở khóa user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.UnlockAccountResponse{
		Success:    true,
		Message:    "Mở khóa user thành công",
		ProfileID:  req.ProfileID,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// GetAccountStatus lấy trạng thái tài khoản
func (s *AuthUsecase) GetAccountStatus(ctx context.Context, req *dto.GetAccountStatusRequest) (*dto.GetAccountStatusResponse, error) {
	// Lấy thông tin user
	userInfo, err := s.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.GetAccountStatusResponse{}, err
	}

	// Xác định trạng thái text
	var statusText string
	switch userInfo.Status {
	case auth.StatusActive:
		statusText = "Hoạt động"
	case (auth.StatusInactive):
		statusText = "Vô hiệu hóa"
	case auth.StatusTemporarilyLocked:
		statusText = "Khóa tạm thời"
	case auth.StatusPermanentlyLocked:
		statusText = "Khóa vĩnh viễn"
	default:
		statusText = "Không xác định"
	}

	// Xác định loại khóa
	var lockType auth.AccountStatus
	if userInfo.IsTemporarilyLocked() {
		lockType = auth.StatusTemporarilyLocked
	} else if userInfo.IsPermanentlyLocked() {
		lockType = auth.StatusPermanentlyLocked
	}

	// Kiểm tra khóa có hết hạn không
	var isExpired bool
	if userInfo.IsTemporarilyLocked() {
		isExpired = userInfo.IsLockExpired()
	}

	// Tạo response
	response := &dto.GetAccountStatusResponse{
		ProfileID:  userInfo.ProfileID,
		Status:     uint32(userInfo.Status),
		StatusText: statusText,
		IsLocked:   userInfo.IsLocked(),
		LockType:   lockType,
		CanLogin:   userInfo.CanLogin(),
		IsExpired:  isExpired,
	}

	// Thêm thông tin khóa nếu có
	if userInfo.IsLocked() {
		response.LockedAt = userInfo.LockedAt.Format("2006-01-02T15:04:05Z07:00")
		response.LockReason = userInfo.LockReason
		response.LockedBy = userInfo.LockedBy

		if userInfo.LockedUntil != nil {
			response.LockedUntil = userInfo.LockedUntil.Format("2006-01-02T15:04:05Z07:00")
		}
	}

	return response, nil
}

// CheckAndUnlockExpiredAccounts kiểm tra và tự động mở khóa các tài khoản hết hạn
func (s *AuthUsecase) CheckAndUnlockExpiredAccounts(ctx context.Context) error {
	return s.UserInfoRepo.CheckLockExpiration(ctx)
}

// LockMultipleUsers khóa nhiều user cùng lúc
func (s *AuthUsecase) LockMultipleUsers(ctx context.Context, req *dto.LockMultipleUsersRequest) (*dto.LockMultipleUsersResponse, error) {
	// Validate duration cho khóa tạm thời
	if req.LockType == auth.StatusTemporarilyLocked && req.Duration <= 0 {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Thời gian khóa tạm thời phải lớn hơn 0",
		}, fmt.Errorf("invalid duration for temporary lock")
	}

	// Thực hiện khóa nhiều user
	err := s.UserInfoRepo.LockMultipleUsers(ctx, req.ProfileIDs, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi khóa nhiều user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.LockMultipleUsersResponse{
		Success:    true,
		Message:    "Khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		LockType:   req.LockType,
		LockedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Thêm thông tin thời gian khóa đến cho khóa tạm thời
	if req.LockType == auth.StatusTemporarilyLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

// UnlockMultipleUsers mở khóa nhiều user cùng lúc
func (s *AuthUsecase) UnlockMultipleUsers(ctx context.Context, req *dto.UnlockMultipleUsersRequest) (*dto.UnlockMultipleUsersResponse, error) {
	// Thực hiện mở khóa nhiều user
	err := s.UserInfoRepo.UnlockMultipleUsers(ctx, req.ProfileIDs, req.Reason, req.UnlockedBy)
	if err != nil {
		return &dto.UnlockMultipleUsersResponse{
			Success: false,
			Message: "Lỗi khi mở khóa nhiều user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.UnlockMultipleUsersResponse{
		Success:    true,
		Message:    "Mở khóa nhiều user thành công",
		ProfileIDs: req.ProfileIDs,
		UnlockedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// CreatePassword tạo mật khẩu cho user chưa có mật khẩu
func (s *AuthUsecase) CreatePassword(ctx context.Context, req dto.CreatePasswordRequest) (*dto.CreatePasswordResponse, error) {
	// Lấy profileId từ context
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, _errors.ReturnError(service.SessionAuthenticationFailed)
	}

	// Kiểm tra password và confirmPassword có khớp không
	if req.Password != req.ConfirmPassword {
		return nil, _errors.ReturnError(service.PasswordConfirmationMismatch)
	}

	// Kiểm tra độ dài password (tối thiểu 8 ký tự)
	if len(req.Password) < 8 {
		return nil, _errors.ReturnError(service.PasswordTooShort)
	}

	// Kiểm tra user đã có auth_method với provider ADMIN chưa
	existingAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profileId, "ADMIN")
	if err == nil && existingAuth != nil {
		return nil, _errors.ReturnError(service.PasswordAlreadySet)
	}

	// Hash password
	hashedPassword, err := _utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Lấy thông tin profile để có username
	profile, err := s.ProfileProvider.GetByProfileID(ctx, profileId)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	// Tạo auth_method mới với provider ADMIN
	newAuthMethod := &auth.AuthMethod{
		Provider: "ADMIN",
		AuthName: profile.Phone, // Sử dụng phone làm username
		Password: hashedPassword,
		UserID:   profileId,
		FullName: profile.FullName,
		Email:    profile.Email,
		Phone:    profile.Phone,
		Avatar:   profile.Avatar,
	}

	// Lưu vào database
	createdAuth, err := s.AuthMethodRepo.Create(ctx, newAuthMethod)
	if err != nil {
		return nil, fmt.Errorf("create password auth method: %w", err)
	}

	// Tạo status và OTP cho auth_method mới (nếu cần)
	statusEntity := &auth.UserStatusEntity{
		Active:   true,
		Verified: true, // Đã verified vì user đã đăng nhập
		IdDomain: auth.IdDomain{
			AuthID: createdAuth.ID,
		},
	}
	s.StatusRepo.CreateStatus(ctx, statusEntity)

	// Trả về response
	return &dto.CreatePasswordResponse{
		AuthId:  createdAuth.ID,
		Message: "Tạo mật khẩu thành công",
		Success: true,
	}, nil
}

// LoginWithPassword đăng nhập bằng username/phone và password
func (s *AuthUsecase) LoginWithPassword(ctx context.Context, req dto.LoginWithPasswordRequest) (*dto.AuthLoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, _errors.ReturnError(service.UsernamePasswordRequired)
	}

	// Tìm auth_method theo username (có thể là phone)
	// Thử tìm theo phone trước
	var authEntity *auth.AuthMethod
	var err error

	// Kiểm tra xem username có phải là số điện thoại không
	if _utils.ValidatePhoneNumber(req.Username) {
		// Resolve the PHONE identity first, then its canonical PASSWORD
		// credential. ADMIN remains a compatibility fallback for accounts whose
		// password credential predates provider separation.
		authEntity, err = s.AuthMethodRepo.FindByPhone(ctx, req.Username)
		if err != nil || authEntity == nil {
			return nil, _errors.ReturnError(service.PhonePasswordInvalid)
		}
		if authEntity.Provider != "PASSWORD" && authEntity.Provider != "ADMIN" {
			if authEntity.UserID != 0 {
				profileID := authEntity.UserID
				authEntity, err = s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profileID, "PASSWORD")
				if err != nil || authEntity == nil {
					authEntity, err = s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profileID, "ADMIN")
				}
				if err != nil || authEntity == nil {
					return nil, _errors.ReturnError(service.PasswordNotSetUseOTP)
				}
			} else {
				return nil, _errors.ReturnError(service.PasswordNotSetUseOTP)
			}
		}
	} else {
		// Prefer canonical end-user password credentials, then accept the
		// existing ADMIN provider for rolling compatibility.
		authEntity, err = s.AuthMethodRepo.FindByAuthNameAndProvider(ctx, req.Username, "PASSWORD")
		if err != nil || authEntity == nil {
			authEntity, err = s.AuthMethodRepo.FindByAuthNameAndProvider(ctx, req.Username, "ADMIN")
		}
		if err != nil || authEntity == nil {
			return nil, _errors.ReturnError(service.UsernameOrPasswordInvalid)
		}
	}

	if authEntity.Provider != "PASSWORD" && authEntity.Provider != "ADMIN" {
		return nil, _errors.ReturnError(service.LoginMethodInvalid)
	}

	// Kiểm tra password
	if !_utils.CheckPasswordHash(req.Password, authEntity.Password) {
		return nil, _errors.ReturnError(service.UsernameOrPasswordInvalid)
	}

	// Lấy profile của user
	profile, err := s.ProfileProvider.GetByProfileID(ctx, authEntity.UserID)
	if err != nil || profile == nil || profile.ProfileID == 0 {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	// Check account status - nếu tài khoản bị khóa thì logout và báo lỗi
	if authEntity.UserID != 0 {
		isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(ctx, authEntity.UserID)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to check account status: %v", err))
		}
		if isLocked || !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(ctx, authEntity.ID)
			return nil, _errors.ReturnError(service.AccountLocked)
		}
	}

	// Tạo session
	deviceID := _utils.GetDeviceIdFromContext(ctx)
	if deviceID == "" {
		deviceID = req.DeviceID
	}
	if deviceID == "" {
		deviceID = req.DeviceName
	}

	otpRequest := dto.OtpVerifyRequest{
		OtpRequest: dto.OtpRequest{
			Phone:    authEntity.Phone,
			Fullname: authEntity.FullName,
		},
		Otp:        "",
		AuthID:     authEntity.ID,
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.OS,
		DeviceName: req.DeviceName,
		DeviceID:   deviceID,
		Fullname:   authEntity.FullName,
	}

	newSession := s.createSession(ctx, authEntity, otpRequest)
	s.SessionRepo.CreateSession(ctx, newSession)

	// Tạo refresh token
	refreshToken := s.GenRefreshToken(ctx, authEntity, newSession.SessionID)

	// Lấy passwordAuthId (chính là auth_method hiện tại)
	response := s.ResponseLogin(ctx, authEntity, profile, newSession.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	response.PasswordAuthId = authEntity.ID

	// Lưu app state = "active" khi login thành công
	if profile.ProfileID != 0 {
		if err := s.CacheProvider.SaveAppState(ctx, profile.ProfileID, "active"); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to save app state for user %d: %v", profile.ProfileID, err))
		}
	}

	// Ghi log admin history
	func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(ctx, loginHistoryTimeout)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(ctx)
		userAgent := _utils.GetUserAgentFromContext(ctx)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Password-Client"
		}
		slog.InfoContext(ctx, fmt.Sprintf("Creating password login history for user %d, IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent))

		success := true
		historyDeviceID := deviceID
		if historyDeviceID == "" {
			historyDeviceID = req.DeviceName
		}
		historyPayload := dto.HistoryAuthCreateDTO{
			UserID:      authEntity.UserID,
			ActionType:  "login",
			ActionName:  "PASSWORD_LOGIN",
			Description: "Đăng nhập hệ thống bằng mật khẩu",
			Success:     &success,
			SessionID:   strconv.FormatUint(newSession.SessionID, 10),
			Channel:     req.Platform,
			DeviceID:    historyDeviceID,
			Metadata: map[string]interface{}{
				"authId":       authEntity.ID,
				"loginChannel": "password",
			},
		}
		historyPayload.IPAddress = ipAddress
		historyPayload.UserAgent = userAgent
		preparedHistory := s.prepareHistoryAuthPayload(ctx, &historyPayload)
		if err := s.NotificationClient.CreateHistoryAuth(contextTimeout, &preparedHistory); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to log password auth history: %v", err))
		}
		if authEntity.Provider != "ADMIN" {
			return
		}

		err := s.NotificationClient.CreateAdminHistory(contextTimeout,
			authEntity.UserID,                             // adminID
			authEntity.UserID,                             // targetID (chính admin đó)
			int32(shared_enum.TargetHistoryLead),          // targetType
			int32(shared_enum.HistoryCreateStep),          // actionType
			"Admin đăng nhập hệ thống",                    // title
			[]string{"Đăng nhập thành công vào hệ thống"}, // notes
			"Chưa đăng nhập",                              // preStage
			"Đã đăng nhập",                                // afterStage
			&authEntity.UserID,                            // ownerID
			int32(shared_enum.EOwnerTypeMember),           // ownerType
			enums.RoleMap[profile.RoleType],               // adminRole
			ipAddress,                                     // ipAddress
			userAgent,                                     // userAgent
		)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to log admin login: %v", err))
		} else {
			slog.InfoContext(ctx, fmt.Sprintf("Successfully logged admin login for user %d", authEntity.UserID))
		}
	}()

	return response, nil
}

// SwitchAccount chuyển đổi tài khoản sang profile khác hoặc tổ chức khác
func (s *AuthUsecase) SwitchAccount(c context.Context, req dto.SwitchAccountRequest) (*dto.SwitchAccountResponse, error) {
	// Lấy profileId hiện tại từ context
	sessionId := _utils.GetSessionIdFromContext(c)
	currentProfileId := _utils.GetProfileIdWithContext(c)
	authId := _utils.GetAuthIdFromContext(c)
	if currentProfileId == 0 {
		return nil, _errors.ReturnError(service.SessionAuthenticationFailed, _errors.WithLegacyCode(400))
	}

	// Kiểm tra quyền truy cập profile mới
	if req.ProfileID != currentProfileId {
		// TODO: Kiểm tra quyền truy cập profile khác (có thể là admin hoặc có quyền đặc biệt)
		// Hiện tại chỉ cho phép chuyển đổi trong cùng profile
		return nil, _errors.ReturnError(service.ProfileAccessDenied)
	}

	// Lấy thông tin profile
	profile, err := s.ProfileProvider.GetByProfileID(c, req.ProfileID)
	if err != nil {
		return nil, _errors.ReturnError(service.ProfileNotFound, _errors.WithCause(err))
	}

	// Lấy thông tin auth hiện tại (sử dụng provider OTP làm mặc định)
	authEntity, err := s.AuthMethodRepo.GetByUserIdAndProvider(c, req.ProfileID, "PHONE")
	if err != nil {
		return nil, _errors.ReturnError(service.AuthenticationInfoNotFound, _errors.WithCause(err))
	}

	// Tạo JWT token mới với thông tin chuyển đổi
	roleIds, _ := s.UserInfoRepo.GetRoleIdsByProfileId(c, currentProfileId)
	payload := _jwt.JwtTokenProperties{
		AuthID:    authId,
		ProfileID: currentProfileId,
		SessionID: sessionId,
		Role:      "ROLE_USER",
		RoleIds:   roleIds,
		Type:      "ACCESS",
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	}

	// Thêm thông tin organization/group nếu có
	payload.AdditionalInfo = make(map[string]interface{})
	if req.OrganizationID != nil {
		payload.AdditionalInfo["organization"] = *req.OrganizationID
	}
	if req.GroupID != nil {
		payload.AdditionalInfo["group"] = *req.GroupID
	}
	if req.OwnerType != nil {
		payload.AdditionalInfo["ownerType"] = *req.OwnerType
	}

	// Tạo access token mới
	accessToken, err := _jwt.GenerateToken(payload)
	if err != nil {
		return nil, fmt.Errorf("generate switch-account access token: %w", err)
	}

	// Tạo response
	response := &dto.SwitchAccountResponse{
		AuthID:         authEntity.ID,
		ProfileID:      req.ProfileID,
		OrganizationID: req.OrganizationID,
		GroupID:        req.GroupID,
		Role:           "ROLE_USER",
		FullName:       profile.FullName,
		Avatar:         profile.Avatar,
		Email:          profile.Email,
		Phone:          profile.Phone,
		AccessToken:    accessToken,
	}

	return response, nil
}

// LoginWithToken đăng nhập bằng refreshToken và trả về thông tin đầy đủ
func (s *AuthUsecase) LoginWithToken(ctx context.Context, refreshToken string) (*dto.AuthLoginResponse, error) {
	// Parse refreshToken để lấy thông tin
	props, err := _jwt.GetProperties(refreshToken)
	if err != nil {
		return nil, _errors.ReturnError(service.TokenInvalidOrExpired)
	}

	// Kiểm tra xem có phải là REFRESH token không
	if props.Type != "REFRESH" {
		return nil, _errors.ReturnError(service.RefreshTokenRequired)
	}

	// Kiểm tra token có trong cache không (validate bằng issuedAt)
	cachedTimestamp, err := s.CacheProvider.GetToken(ctx, props.ProfileID)
	if err != nil {
		slog.InfoContext(ctx, fmt.Sprintf("⚠️ [LoginWithToken] Không tìm thấy token trong cache cho authID %d: %v", props.AuthID, err))
		// Không tìm thấy thì vẫn cho login tiếp
	} else {
		invalidToken := _utils.ValidTokenByIssueAt(props.IssuedAt, cachedTimestamp)
		if invalidToken {
			return nil, _errors.ReturnError(service.SessionExpired)
		}
	}

	// Lấy auth entity
	authEntity, err := s.AuthMethodRepo.FindByID(ctx, props.AuthID)
	if err != nil || authEntity == nil {
		return nil, _errors.ReturnError(service.LoginInfoNotFound)
	}

	// Lấy profile
	profile, err := s.ProfileProvider.GetByProfileID(ctx, props.ProfileID)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	// Check account status
	isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(ctx, props.ProfileID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to check account status: %v", err))
	}
	if isLocked || !canLogin {
		return nil, _errors.ReturnError(service.AccountLocked)
	}

	// Tạo access token mới
	roleIds := props.RoleIds
	if len(roleIds) == 0 {
		roleIds, _ = s.UserInfoRepo.GetRoleIdsByProfileId(ctx, props.ProfileID)
	}

	var accessToken string
	if profile.FullName == "" {
		// Chưa có fullname, trả TEMP token
		accessToken, _ = _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID:         props.AuthID,
			ProfileID:      props.ProfileID,
			SessionID:      props.SessionID,
			Role:           props.Role,
			RoleIds:        roleIds,
			Type:           _jwt.TempToken,
			OrganizationID: props.OrganizationID,
			PlanID:         profile.PlanID,
			PlanAt:         profile.PlanAt,
		})
	} else {
		// Đã có fullname, trả ACCESS token bình thường
		accessToken, _ = _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID:         props.AuthID,
			ProfileID:      props.ProfileID,
			SessionID:      props.SessionID,
			Role:           props.Role,
			RoleIds:        roleIds,
			Type:           _jwt.AccessToken,
			OrganizationID: props.OrganizationID,
			PlanID:         profile.PlanID,
			PlanAt:         profile.PlanAt,
		})
	}

	// Tạo response
	response := dto.MakeAuthLoginResponse(accessToken, authEntity, props.Role)
	response.RefreshToken = refreshToken // Giữ nguyên refresh token
	response.RequireFullname = (profile.FullName == "")
	response.ReferralCode = profile.ReferralCode
	if props.OrganizationID != nil {
		response.OrganizationID = *props.OrganizationID
	} else {
		response.OrganizationID = 0
	}

	// Lấy passwordAuthId nếu có
	if profile.ProfileID != 0 {
		passwordAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profile.ProfileID, "ADMIN")
		if err == nil && passwordAuth != nil {
			response.PasswordAuthId = passwordAuth.ID
		}
	}

	// Nếu chưa có fullname thì set tempToken
	if response.RequireFullname {
		response.TempToken = accessToken
		response.AccessToken = ""
	} else {
		response.TempToken = ""
	}

	return response, nil
}

// GetLoginHistory lấy danh sách lịch sử đăng nhập của user hiện tại
func (s *AuthUsecase) GetLoginHistory(ctx context.Context, page, size int32) ([]*auth.UserSessionEntity, int64, error) {
	authID := _utils.GetAuthIdFromContext(ctx)
	if authID == 0 {
		return nil, 0, _errors.ReturnError(service.LoginInfoNotFound, _errors.WithLegacyCode(401))
	}

	sessions, total, err := s.SessionRepo.GetLoginHistoryByAuthID(ctx, authID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("get login history: %w", err)
	}

	return sessions, total, nil
}

// GetQRStatus kiểm tra trạng thái phiên QR
func (s *AuthUsecase) GetQRStatus(ctx context.Context, req *dto.QRStatusRequest) (*dto.QRStatusResponse, error) {
	if req == nil || req.SessionId == 0 || req.SessionKey == "" {
		return nil, _errors.ReturnError(service.QRSessionInvalid)
	}

	session, err := s.SessionRepo.GetBySessionID(ctx, req.SessionId)
	if err != nil || session == nil {
		return &dto.QRStatusResponse{
			Status:  "expired",
			Message: "Phiên QR không tồn tại hoặc đã hết hạn",
		}, nil
	}

	if session.SessionKey != req.SessionKey {
		return nil, _errors.ReturnError(service.LoginSessionIncorrect)
	}

	if session.FinishedDate != nil || session.LogoutAt != nil {
		return &dto.QRStatusResponse{
			Status:  "expired",
			Message: "Phiên QR đã hết hạn",
		}, nil
	}

	if session.Activate {
		return &dto.QRStatusResponse{
			Status:  "confirmed",
			Message: "Đăng nhập QR thành công",
		}, nil
	}

	return &dto.QRStatusResponse{
		Status:  "pending",
		Message: "Đang chờ xác nhận QR",
	}, nil
}

// UpdateFullname cập nhật họ tên cho user và trả về thông tin login
func (s *AuthUsecase) UpdateFullname(ctx context.Context, fullname string) (*dto.AuthLoginResponse, error) {
	// Validate fullname
	if fullname == "" {
		return nil, _errors.ReturnError(service.FullNameRequired)
	}

	// Lấy profileId từ context
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, _errors.ReturnError(service.UserAuthenticationFailed)
	}

	// Lấy profile hiện tại
	profile, err := s.ProfileProvider.GetByProfileID(ctx, profileId)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(service.UserInfoNotFound)
	}

	// Update fullname
	profile.FullName = fullname
	err = s.ProfileProvider.UpdateProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("update profile full name: %w", err)
	}

	// Lấy authId từ context
	authId := _utils.GetAuthIdFromContext(ctx)
	if authId == 0 {
		return nil, _errors.ReturnError(service.UserAuthenticationFailed)
	}

	// Lấy sessionId từ context
	sessionId := _utils.GetSessionIdFromContext(ctx)
	if sessionId == 0 {
		return nil, _errors.ReturnError(service.LoginSessionNotFound)
	}

	// Lấy auth entity
	authEntity, err := s.AuthMethodRepo.FindByID(ctx, authId)
	if err != nil || authEntity == nil {
		return nil, _errors.ReturnError(service.LoginInfoNotFound)
	}

	// Tạo ACCESS token mới (không phải TEMP nữa vì đã có fullname)
	roleIds, _ := s.UserInfoRepo.GetRoleIdsByProfileId(ctx, profileId)
	accessToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authId,
		ProfileID: profileId,
		SessionID: sessionId,
		Role:      enums.RoleMap[profile.RoleType],
		RoleIds:   roleIds,
		Type:      _jwt.AccessToken, // Token bình thường, không phải TEMP
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	})

	// Refresh token
	refreshToken := s.GenRefreshToken(ctx, authEntity, sessionId)

	// Tạo response giống như login thành công
	response := dto.MakeAuthLoginResponse(accessToken, authEntity, enums.RoleMap[profile.RoleType])
	response.RefreshToken = refreshToken
	response.RequireFullname = false // Không cần fullname nữa
	response.TempToken = ""          // Xóa temp token
	response.ReferralCode = profile.ReferralCode

	// Lấy passwordAuthId nếu có
	if profileId != 0 {
		passwordAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profileId, "ADMIN")
		if err == nil && passwordAuth != nil {
			response.PasswordAuthId = passwordAuth.ID
		}
	}

	return response, nil
}

// GetSessionsByProfile lấy danh sách session theo profileId từ request
func (s *AuthUsecase) GetSessionsByProfile(ctx context.Context, pagable *_dto.Pagable) ([]*auth.UserSessionEntity, int64, error) {
	if pagable == nil {
		pagable = &_dto.Pagable{}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, 0, _errors.ReturnError(service.UserInfoNotFound, _errors.WithLegacyCode(401))
	}

	sessions, total, err := s.SessionRepo.GetSessionsByProfile(ctx, profileID, pagable)
	if err != nil {
		return nil, 0, fmt.Errorf("get profile sessions: %w", err)
	}

	return sessions, total, nil
}

// removeProfileIdFromDevices xóa profileId khỏi tất cả devices của user
func (s *AuthUsecase) removeProfileIdFromDevices(ctx context.Context, profileID uint64) {
	if s.DeviceRepo == nil {
		slog.WarnContext(ctx, fmt.Sprintf("DeviceRepo is nil, skipping remove profileId from devices"))
		return
	}

	// Lấy tất cả devices của profileId
	devices, err := s.DeviceRepo.ListByProfileID(ctx, profileID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to list devices for profileId %d: %v", profileID, err))
		return
	}

	// Xóa profileId khỏi mỗi device
	for _, device := range devices {
		device.ProfileID = nil
		if _, err := s.DeviceRepo.Update(ctx, device); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to remove profileId from device %s: %v", device.DeviceID, err))
		}
	}

	if len(devices) > 0 {
		slog.InfoContext(ctx, fmt.Sprintf("Removed profileId from %d device(s) for profileId %d", len(devices), profileID))
	}
}

func (s *AuthUsecase) PhoneCheck(ctx context.Context, phone string) (bool, *uint64, error) {
	// Check if phone exists
	auth, err := s.AuthMethodRepo.PhoneCheck(ctx, phone)
	if err != nil {
		return false, nil, fmt.Errorf("check phone existence: %w", err)
	}
	if auth == nil {
		return false, nil, nil
	}

	return true, &auth.ID, nil
}

func (s *AuthUsecase) GenerateAccessToken(ctx context.Context, authId uint64, sessionId uint64, organizationId *uint64) (string, error) {
	authMethod, err := s.AuthMethodRepo.FindByID(ctx, authId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("GenerateAccessToken error - FindByID failed: %v", err))
		return "", err
	}

	var originID uint64
	originCtx, cancelOriginLookup := context.WithTimeout(ctx, optionalCRMOriginTimeout)
	originID, err = s.CrmProvider.GetOriginID(originCtx, authMethod.UserID)
	cancelOriginLookup()
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Warning: Cannot get origin ID from CRM (service may be down), using 0: %v", err))
		originID = 0
	}

	profile, err := s.ProfileProvider.GetByProfileID(ctx, authMethod.UserID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("GenerateAccessToken error - GetByProfileID failed: %v", err))
		return "", err
	}

	roleIds, _ := s.UserInfoRepo.GetRoleIdsByProfileId(ctx, authMethod.UserID)

	accessToken, err := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:         authId,
		OriginID:       originID,
		ProfileID:      authMethod.UserID,
		Role:           enums.RoleMap[enums.ERole(authMethod.RoleKey)],
		RoleIds:        roleIds,
		SessionID:      sessionId,
		OrganizationID: organizationId,
		PlanID:         profile.PlanID,
		PlanAt:         profile.PlanAt,
		Type:           _jwt.AccessToken,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("GenerateAccessToken error - GenerateToken failed: %v", err))
		return "", err
	}
	slog.InfoContext(ctx, fmt.Sprintf("GenerateAccessToken success: authId=%d, originID=%d", authId, originID))
	return accessToken, nil
}
