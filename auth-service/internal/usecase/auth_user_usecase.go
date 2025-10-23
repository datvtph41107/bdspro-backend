package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"auth/config"
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/enums"
	"auth/internal/interface/factory"
	"auth/internal/interface/provider"
	"auth/internal/interface/repo"
	_db "common/db"
	_errors "common/errors"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	shared_enum "pb/enums"
)

type AuthUsecase struct {
	CookieProvider     provider.CookieProvider
	CacheProvider      provider.CacheProvider
	PreePlanID         *uint64
	ProfileProvider    provider.ProfileProvider
	AuthMethodRepo     repo.AuthMethodRepository
	UserInfoRepo       repo.UserInfoRepository
	SessionRepo        repo.SessionRepository
	OTPRepo            repo.OTPRepository
	StatusRepo         repo.StatusRepository
	OtpUsecase         *OtpUsecase
	OrganizationClient provider.OrganizationProvider
	PaymentProvider    provider.PaymentProvider
	NotificationClient provider.NotificationProvider
	ZnsProvider        provider.IZnsProvider
	properties         dto.PropertiesDTO
}

func NewAuthService(cookieProvider provider.CookieProvider,
	cacheProvider provider.CacheProvider,
	profileProvider provider.ProfileProvider,
	authMethodRepo repo.AuthMethodRepository,
	userInfoRepo repo.UserInfoRepository,
	sessionRepo repo.SessionRepository,
	otpRepo repo.OTPRepository,
	statusRepo repo.StatusRepository,
	otpUsecase *OtpUsecase,
	organizationClient provider.OrganizationProvider,
	paymentProvider provider.PaymentProvider,
	notificationClient provider.NotificationProvider,
	znsProvider provider.IZnsProvider,
	properties factory.IFactory,
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
		OtpUsecase:         otpUsecase,
		OrganizationClient: organizationClient,
		PaymentProvider:    paymentProvider,
		NotificationClient: notificationClient,
		ZnsProvider:        znsProvider,
		properties:         properties.GetProperties(),
	}
}

func (s *AuthUsecase) RequestOtp(c context.Context, otp dto.OtpRequest) (*dto.RequestLoginResponse, error) {
	valid := _utils.ValidatePhoneNumber(otp.Phone)
	if !valid {
		return nil, _errors.ReturnError(
			int32(400),
			"Số điện thoại không đúng",
		)
	}
	existed := false
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
		s.SendNewOTP(c, param.OTP, otp.Phone) // todo
	} else {
		// Bắt buộc phải có họ tên khi đăng ký
		// if otp.Fullname == "" {
		// 	return nil, _errors.ReturnError(400, "Vui lòng nhập họ tên khi đăng ký")
		// }
		oauth, err = s.RequestIsRegister(c, otp)
		if err != nil {
			return nil, err
		}
	}

	return &dto.RequestLoginResponse{
		AuthID:   oauth.ID,
		Provider: "PHONE",
		OAuthID:  otp.Phone,
		Existed:  existed,
	}, nil
}

func (s *AuthUsecase) ResendOTP(c context.Context, otp dto.OtpRequest) (*dto.ResendOTPResponse, error) {
	oauth, err := s.AuthMethodRepo.FindByPhone(c, otp.Phone)
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

	s.SendNewOTP(c, param.OTP, otp.Phone)
	return &dto.ResendOTPResponse{
		AuthId:  oauth.ID,
		Message: "OTP đã được gửi",
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
		// Kiểm tra nếu là user mới thì bắt buộc phải có họ tên
		if otp.Fullname == "" {
			return nil, _errors.ReturnError(400, "Vui lòng nhập họ tên khi đăng ký tài khoản mới")
		}

		// infoEntity = s.ProfileProvider.MakeUserProfileEntity(e.ID, s.PreePlanID, e.FullName, e.Email, e.Phone, e.Avatar)
		infoEntity, err = s.ProfileProvider.CreateProfile(c, e.ID, s.PreePlanID, otp.Fullname, e.Email, e.AuthName, e.Avatar)
		if err != nil {
			return nil, err
		}

		e.UserID = infoEntity.ProfileID
		s.AuthMethodRepo.Update(c, e)
		isNewUser = true
	}

	// Check account status - nếu tài khoản bị khóa thì logout và báo lỗi
	if e.UserID != 0 {
		isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(c, e.UserID)
		if err != nil {
			log.Printf("Failed to check account status: %v", err)
		}
		if isLocked || !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(c, e.ID)
			return nil, _errors.ReturnError(
				int32(403),
				"Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên để được hỗ trợ",
			)
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
		var oauthAcc *domain.AuthMethod
		if err := _db.DB.First(&oauthAcc, otp.AuthID).Error; err == nil {
			if oauthAcc != nil && oauthAcc.UserID == 0 {
				// var statusOauth *domain.UserStatusEntity
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
		go func() {
			contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := s.NotificationClient.CreateNotification(
				contextTimeout,
				"Chào mừng bạn đến với BDSPro",
				"Thiết lập tài khoản thành công",
				1, // type: 1 - thông báo hệ thống
				"",
				infoEntity.ProfileID,
			)
			if err != nil {
				log.Printf("Failed to send welcome notification: %v", err)
			} else {
				log.Printf("Successfully sent welcome notification to user %d", infoEntity.ProfileID)
			}
		}()
	}

	// Lấy passwordAuthId (auth_method có provider ADMIN)
	response := s.ResponseLogin(e, infoEntity, newSession.SessionID)
	response.RefreshToken = refreshToken            // Set refreshToken vào response body
	response.ReferralCode = infoEntity.ReferralCode // Set mã giới thiệu vào response
	if infoEntity.ProfileID != 0 {
		passwordAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(c, infoEntity.ProfileID, "ADMIN")
		if err == nil && passwordAuth != nil {
			response.PasswordAuthId = passwordAuth.ID
		}
	}

	// Ghi log admin history
	go func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

		log.Printf("Creating login history for user %d, IP: %s, UserAgent: %s",
			e.UserID, ipAddress, userAgent)

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
			log.Printf("Failed to log login history: %v", err)
		} else {
			log.Printf("Successfully logged login history for user %d", e.UserID)
		}
	}()

	return response, nil
}

func (s *AuthUsecase) createSession(c context.Context, e *domain.AuthMethod, otp dto.OtpVerifyRequest) *domain.UserSessionEntity {
	return &domain.UserSessionEntity{
		AuthID:     e.ID,
		Version:    otp.Version,
		Platform:   otp.Platform,
		OS:         otp.OS,
		DeviceName: otp.DeviceName,
		IPRequest:  _utils.AuditIPFromContext(c),
		UserAgent:  _utils.GetUserAgentFromContext(c),
	}
}

func (s *AuthUsecase) RefreshToken(c context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	claims := _jwt.GetPrincipalContext(c)
	fmt.Println("claims", claims)

	log.Printf("refreshToken: %s :???", refreshToken)
	props, err := _jwt.GetProperties(refreshToken)
	if err != nil {
		return nil, &_routes.Except{
			Code:    401,
			Message: err.Error(),
		}
	}

	if claims != nil && claims.OrganizationId != nil {
		props.OrganizationID = claims.OrganizationId
	}

	// Lấy timestamp từ cache
	cachedTimestamp, err := s.CacheProvider.GetToken(c, props.AuthID)
	if err != nil {
		log.Printf("⚠️ [RefreshToken] Không tìm thấy token trong cache cho authID %d: %v. Bỏ qua kiểm tra.", props.AuthID, err)
		// Không tìm thấy trong cache thì bỏ qua kiểm tra, tiếp tục xử lý
	} else {
		// Validate token bằng cách so sánh Unix timestamp (UTC)
		invalidToken := _utils.ValidTokenByIssueAt(props.IssuedAt, cachedTimestamp)
		if invalidToken {
			log.Printf("⚠️ [RefreshToken] Token không hợp lệ cho authID %d. IssuedAt: %s, CachedTime: %s",
				props.AuthID, props.IssuedAt.Format(time.RFC3339), cachedTimestamp)
			return nil, &_routes.Except{
				Code:    401,
				Message: "Phiên làm việc đã hết hạn vui lòng đăng nhập lại",
			}
		}
	}

	if props.Type == "REFRESH" {
		profile, err := s.ProfileProvider.GetByProfileID(c, props.ProfileID)
		if err != nil {
			return nil, err
		}
		newToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID:         props.AuthID,
			ProfileID:      props.ProfileID,
			Role:           props.Role,
			SessionID:      props.SessionID,
			OrganizationID: props.OrganizationID,
			PlanID:         profile.PlanID,
			PlanAt:         profile.PlanAt,
			Type:           "ACCESS",
		})
		log.Printf("✅ [RefreshToken] Refresh token thành công cho authID %d", props.AuthID)
		return &dto.RefreshResponse{AccessToken: newToken}, nil
	}

	return nil, &_routes.Except{
		Code:    401,
		Message: "Token không đúng",
	}
}

func (s *AuthUsecase) Logout(c context.Context) (*dto.LogoutResponse, error) {
	sessionID := _utils.GetSessionIdFromContext(c)
	currentSession, err := s.SessionRepo.GetBySessionID(c, sessionID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	currentSession.FinishedDate = &now
	s.SessionRepo.UpdateSession(c, currentSession)
	s.CacheProvider.DeleteToken(c, currentSession.AuthID)

	return &dto.LogoutResponse{Message: "Đăng xuất thành công"}, nil
}

// Restore khôi phục tài khoản
func (s *AuthUsecase) Restore(c context.Context, query dto.RestorePhoneRequest) (*dto.RestoreAccountResponse, error) {
	existingAuth, err := s.AuthMethodRepo.FindByPhone(c, query.Phone)
	ok := existingAuth != nil
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số điện thoại đã liên kết với tài khoản khác. Vui lòng nhập số điện thoại khác",
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)

	e := domain.AuthMethod{
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
	var deletedAuth *domain.AuthMethod
	var err error

	// Tìm tài khoản đã bị xóa theo phone, email hoặc username
	if req.Phone != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByPhone(c, req.Phone)
	} else if req.Email != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByEmail(c, req.Email)
	} else if req.Username != "" {
		deletedAuth, err = s.AuthMethodRepo.FindDeletedByUsername(c, req.Username)
	} else {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Vui lòng cung cấp phone, email hoặc username để khôi phục tài khoản",
		}
	}

	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy tài khoản đã bị xóa",
		}
	}

	if deletedAuth == nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy tài khoản đã bị xóa",
		}
	}

	// Kiểm tra xem tài khoản đã được khôi phục chưa
	existingAuth, err := s.AuthMethodRepo.FindByID(c, deletedAuth.ID)
	if err == nil && existingAuth != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tài khoản đã được khôi phục trước đó",
		}
	}

	// Khôi phục tài khoản bằng cách set deleted_at = null
	err = s.AuthMethodRepo.RestoreAccount(c, deletedAuth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi khi khôi phục tài khoản",
		}
	}

	return &dto.RestoreDeletedAccountResponse{
		AuthId:   deletedAuth.ID,
		Provider: deletedAuth.Provider,
		AuthName: deletedAuth.AuthName,
		Message:  "Tài khoản đã được khôi phục thành công",
	}, nil
}

// Xóa tài khoản
func (s *AuthUsecase) Delete(c context.Context) (*dto.DeleteAccountResponse, error) {
	principal := _jwt.GetPrincipalContext(c)
	profile, err := s.ProfileProvider.GetByProfileID(c, principal.ProfileId)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tài khoản không tồn tại",
		}
	}

	t := time.Now()
	profile.DeletedAt = &t

	if err := s.ProfileProvider.UpdateProfile(c, profile); err != nil {
		return nil, err
	}

	s.ProfileProvider.SoftDeleteProfile(c, principal.ProfileId)
	s.CookieProvider.ResetAll(c)
	s.CacheProvider.DeleteToken(c, principal.AuthID)

	return &dto.DeleteAccountResponse{
		// ProfileId: profile.ProfileID,
		Message: "Tài khoản đã được xóa",
	}, nil
}

func (s *AuthUsecase) RequestIsRegister(c context.Context, dto dto.OtpRequest) (*domain.AuthMethod, error) {
	// Tạo một UserAuthEntity mới
	e := &domain.AuthMethod{
		Provider: "PHONE",
		AuthName: dto.Phone,
		FullName: dto.Fullname,
	}

	// Lưu vào database
	createdAuth, err := s.AuthMethodRepo.Create(c, e)
	if err != nil {
		return nil, err
	}

	// Tạo AuthParam mới
	param, err := s.NewUser(c, createdAuth.ID)
	if err != nil {
		return nil, err
	}

	// Gửi OTP trong updateNewOTP luôn
	s.SendNewOTP(c, param.OTP, dto.Phone)

	return createdAuth, nil
}

func (s *AuthUsecase) NewUser(c context.Context, oauthId uint64) (*dto.AuthParam, error) {
	// Tạo UserStatusEntity mới
	statusEntity := &domain.UserStatusEntity{
		Active:   true,
		Verified: false,
		IdDomain: domain.IdDomain{
			AuthID: oauthId,
		},
	}

	// Lưu vào database
	if err := s.StatusRepo.CreateStatus(c, statusEntity); err != nil {
		return nil, err
	}

	// Tạo UserOTPEntity mới
	now := time.Now()
	otpEntity := &domain.UserOTPEntity{
		OTP: _utils.GenerateOTP(),
		IdDomain: domain.IdDomain{
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
	}, nil
}

func (s *AuthUsecase) SendNewOTP(c context.Context, otpEntity *domain.UserOTPEntity, phone string) string {
	ip := _utils.AuditIPFromContext(c)
	s.CheckOtpSpam(c, phone, ip)

	otp := _utils.GenerateOTP()
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

	// Gửi OTP qua ZNS
	if s.ZnsProvider != nil {
		result, err := s.ZnsProvider.SendOTPZNS(c, phone, otp)
		if err != nil {
			log.Printf("⚠️ [ZNS] Lỗi khi gửi OTP qua ZNS cho phone %s: %v", phone, err)
			// Log lỗi nhưng không block flow, vì OTP vẫn được lưu vào DB
		} else {
			if success, ok := result["success"].(bool); ok && success {
				log.Printf("✅ [ZNS] Đã gửi OTP thành công cho phone %s", phone)
			} else {
				log.Printf("⚠️ [ZNS] Gửi OTP thất bại cho phone %s: %v", phone, result["message"])
			}
		}
	} else {
		log.Printf("⚠️ [ZNS] ZnsProvider chưa được khởi tạo, không gửi OTP qua ZNS")
	}

	return otp
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
		return &_routes.Except{
			Code:    401,
			Message: "Bạn đã gửi quá nhiều OTP. Vui lòng thử lại sau.",
		}
	}

	// Kiểm tra giới hạn gửi OTP theo IP
	requestCountIp, err := s.CookieProvider.Get(c, redisKeyIp)
	countIp, _ := strconv.ParseInt(requestCountIp, 10, 32)
	if err != nil {
		return err
	}
	if countIp >= s.properties.LimitOtpIP {
		return &_routes.Except{
			Code:    401,
			Message: "IP của bạn đã bị chặn do gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
		}
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
func (s *AuthUsecase) GetParamValidate(c context.Context, oauth *domain.AuthMethod, otp string) (*dto.AuthParam, error) {
	if oauth == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng, vui lòng kiểm tra lại",
		}
	}

	statusEntity, err := s.StatusRepo.GetByID(c, oauth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy trạng thái người dùng",
		}
	}

	otpEntity, err := s.OTPRepo.GetByID(c, oauth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy OTP",
		}
	}

	return &dto.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: otp,
	}, nil
}

func (s *AuthUsecase) ResponseLogin(authEntity *domain.AuthMethod, profile *dto.ProfileDTO, sessionID uint64) *dto.AuthLoginResponse {
	if authEntity.UserID == 0 {
		return &dto.AuthLoginResponse{
			AuthID:   authEntity.ID,
			FullName: authEntity.FullName,
		}
	}
	accessToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
		SessionID: sessionID,
		Role:      enums.RoleMap[profile.RoleType],
		Type:      "ACCESS",
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	})

	return dto.MakeAuthLoginResponse(accessToken, authEntity, enums.RoleMap[profile.RoleType])
}

func (s *AuthUsecase) GetRefreshToken(authEntity *domain.AuthMethod, sessionID uint64) string {
	refreshToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
		// PlanID: ,
		SessionID: sessionID,
		Role:      "ROLE_USER",
		Type:      "REFRESH",
	})
	return refreshToken
}

func (s *AuthUsecase) CreateSession(c context.Context, e *domain.AuthMethod, otp *dto.OtpVerifyRequest) *domain.UserSessionEntity {
	var authId uint64
	if e != nil {
		authId = e.ID
	}
	if otp == nil {
		otp = &dto.OtpVerifyRequest{}
	}
	return &domain.UserSessionEntity{
		AuthID:     authId,
		Version:    otp.Version,
		Platform:   otp.Platform,
		OS:         otp.OS,
		DeviceName: otp.DeviceName,
		IPRequest:  _utils.AuditIPFromContext(c),
		UserAgent:  _utils.GetUserAgentFromContext(c),
		Activate:   false,
	}
}

// CreateAndSaveSession tạo session và lưu vào database
func (s *AuthUsecase) CreateAndSaveSession(c context.Context, e *domain.AuthMethod, otp *dto.OtpVerifyRequest) (*domain.UserSessionEntity, error) {
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
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
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
		return nil, &_routes.Except{
			Code:    400,
			Message: "Truy cập đạt giới hạn vui lòng thử lại sau",
		}
	}

	sessionEntity, err := s.SessionRepo.GetBySessionID(c, sessionId)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	if sessionEntity.FinishedDate != nil ||
		sessionEntity.Activate ||
		sessionEntity.SessionKey != e.SessionKey {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Phiên đăng nhập không chính xác",
		}
	}

	sessionEntity.Activate = true
	s.SessionRepo.UpdateSession(c, sessionEntity)

	authEntity, err := s.AuthMethodRepo.FindByID(c, sessionEntity.AuthID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	refreshToken := s.GenRefreshToken(c, authEntity, sessionId)
	profile, _ := s.ProfileProvider.GetByProfileID(c, authEntity.UserID)

	// Ghi log admin history
	go func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

		log.Printf("Creating login history for user %d (QR), IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent)

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
			log.Printf("Failed to log QR login history: %v", err)
		} else {
			log.Printf("Successfully logged QR login history for user %d", authEntity.UserID)
		}
	}()

	response := s.ResponseLogin(authEntity, profile, sessionEntity.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}

// ------------------------------------------------------------

func (s *AuthUsecase) GenRefreshToken(c context.Context, authEntity *domain.AuthMethod, sessionId uint64) string {
	// Lưu Unix timestamp (UTC) vào cache để tránh vấn đề timezone
	// Trừ 5 giây để đảm bảo token mới được issue sau thời điểm này sẽ hợp lệ
	validFromTimestamp := time.Now().UTC().Add(-5 * time.Second).Unix()
	s.CacheProvider.SaveToken(c, authEntity.ID, fmt.Sprintf("%d", validFromTimestamp), s.properties.JwtRefreshExpMinutes*60)

	refreshToken := s.GetRefreshToken(authEntity, sessionId)
	// Không set cookie nữa, trả refreshToken trong response body
	// s.CookieProvider.SetRefreshToken(c, refreshToken)
	log.Printf("✅ [GenRefreshToken] authID=%d, validFrom=%d (%s), refreshToken=%s",
		authEntity.ID, validFromTimestamp, time.Unix(validFromTimestamp, 0).UTC().Format(time.RFC3339), refreshToken)
	return refreshToken
}

// LogLoginHistory ghi lại lịch sử đăng nhập
func (s *AuthUsecase) LogLoginHistory(c context.Context, authEntity *domain.AuthMethod, profile *dto.ProfileDTO, title string) {
	go func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			log.Printf("Cannot log login history: profile is nil for user %d", authEntity.UserID)
			return
		}

		log.Printf("Creating login history for user %d, IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent)

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
			log.Printf("Failed to log login history: %v", err)
		} else {
			log.Printf("Successfully logged login history for user %d", authEntity.UserID)
		}
	}()
}

func (s *AuthUsecase) SwitchOrganization(c context.Context, organizationID uint64) (*dto.AuthLoginResponse, error) {
	// principal := _jwt.MetadataFromContext(c)
	authId := _utils.GetAuthIdFromContext(c)
	sessionId := _utils.GetSessionIdFromContext(c)
	profileId := _utils.GetProfileIdWithContext(c)
	if profileId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể xác thực phiên đăng nhập",
		}
	}
	clientResponse, err := s.OrganizationClient.SwitchOrganization(c, organizationID)
	fmt.Printf("clientResponse=%+v\n", clientResponse)
	if err != nil {
		return nil, err
	}
	if clientResponse.OrganizationId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể chuyển đổi tổ chức",
		}
	}
	// authEntity, err := s.AuthMethodRepo.GetByID(c, principal.AuthID)
	// if err != nil {
	// 	return nil, err
	// }
	profile, _ := s.ProfileProvider.GetByProfileID(c, profileId)
	payload := _jwt.JwtTokenProperties{
		AuthID:    authId,
		ProfileID: profileId,
		SessionID: sessionId,
		Role:      "ROLE_USER",
		Type:      "ACCESS",
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	}
	payload.AdditionalInfo = map[string]interface{}{
		"organization": clientResponse.OrganizationId,
	}
	newAccessToken, err := _jwt.GenerateToken(payload)
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
		OrganizationID: clientResponse.OrganizationId,
		ReferralCode:   profile.ReferralCode, // Mã giới thiệu
	}

	return response, nil
}

// AdminLogin xử lý đăng nhập admin bằng username và password
func (s *AuthUsecase) AdminLogin(c context.Context, adminLogin dto.AdminLoginRequest) (*dto.AuthLoginResponse, error) {
	// Validate input
	if adminLogin.Username == "" || adminLogin.Password == "" {
		return nil, _errors.ReturnError(
			int32(400),
			"Username và password không được để trống",
		)
	}

	// Tìm admin trong database
	adminAuthEntity, err := s.AuthMethodRepo.FindByUsername(c, adminLogin.Username)
	if err != nil {
		return nil, _errors.ReturnError(
			int32(401),
			"Username hoặc password không đúng",
		)
	}

	// Kiểm tra provider phải là ADMIN
	if adminAuthEntity.Provider != "ADMIN" {
		return nil, _errors.ReturnError(
			int32(401),
			"Tài khoản không có quyền admin",
		)
	}

	// Kiểm tra password
	if !_utils.CheckPasswordHash(adminLogin.Password, adminAuthEntity.Password) {
		return nil, _errors.ReturnError(
			int32(401),
			"Username hoặc password không đúng",
		)
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
			log.Printf("Failed to check account status: %v", err)
		}
		if isLocked || !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(c, adminAuthEntity.ID)
			return nil, _errors.ReturnError(
				int32(403),
				"Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên để được hỗ trợ",
			)
		}
	}

	// Tạo session cho admin
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
		Fullname:   adminAuthEntity.FullName,
	}

	newSession := s.createSession(c, adminAuthEntity, adminOtpRequest)
	s.SessionRepo.CreateSession(c, newSession)

	// Tạo refresh token
	refreshToken := s.GenRefreshToken(c, adminAuthEntity, newSession.SessionID)

	adminProfile.RoleType = enums.ERoleAdmin
	// Tạo response với đầy đủ thông tin
	response := s.ResponseLogin(adminAuthEntity, adminProfile, newSession.SessionID)
	response.RefreshToken = refreshToken              // Set refreshToken vào response body
	response.ReferralCode = adminProfile.ReferralCode // Set mã giới thiệu vào response
	response.Role = enums.RoleMap[adminProfile.RoleType]
	response.Username = adminAuthEntity.AuthName

	// Ghi log admin history
	go func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

		log.Printf("Creating admin history for user %d, IP: %s, UserAgent: %s",
			adminAuthEntity.UserID, ipAddress, userAgent)

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
			log.Printf("Failed to log admin login: %v", err)
		} else {
			log.Printf("Successfully logged admin login for user %d", adminAuthEntity.UserID)
		}
	}()

	return response, nil
}

// LockAccount khóa tài khoản (tạm thời hoặc vĩnh viễn)
func (s *AuthUsecase) LockAccount(ctx context.Context, req *dto.LockAccountRequest) (*dto.LockAccountResponse, error) {
	// Validate duration cho khóa tạm thời
	if req.LockType == domain.StatusTemporarilyLocked && req.Duration <= 0 {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Thời gian khóa tạm thời phải lớn hơn 0",
		}, fmt.Errorf("invalid duration for temporary lock")
	}

	// Kiểm tra user đã bị khóa chưa
	userInfo, err := s.UserInfoRepo.GetUserStatus(ctx, req.ProfileID)
	if err != nil {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Lỗi khi kiểm tra trạng thái user: " + err.Error(),
		}, err
	}

	if userInfo.IsLocked() {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "User đã bị khóa",
		}, fmt.Errorf("user already locked")
	}

	// Thực hiện khóa user
	err = s.UserInfoRepo.LockUser(ctx, req.ProfileID, req.LockType, req.Duration, req.Reason, req.LockedBy)
	if err != nil {
		return &dto.LockAccountResponse{
			Success: false,
			Message: "Lỗi khi khóa user: " + err.Error(),
		}, err
	}

	// Tạo response
	response := &dto.LockAccountResponse{
		Success:   true,
		Message:   "Khóa user thành công",
		ProfileID: req.ProfileID,
		LockType:  req.LockType,
		LockedAt:  time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Thêm thông tin thời gian khóa đến cho khóa tạm thời
	if req.LockType == domain.StatusTemporarilyLocked {
		lockedUntil := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		response.LockedUntil = lockedUntil.Format("2006-01-02T15:04:05Z07:00")
	}

	return response, nil
}

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
	case domain.StatusActive:
		statusText = "Hoạt động"
	case (domain.StatusInactive):
		statusText = "Vô hiệu hóa"
	case domain.StatusTemporarilyLocked:
		statusText = "Khóa tạm thời"
	case domain.StatusPermanentlyLocked:
		statusText = "Khóa vĩnh viễn"
	default:
		statusText = "Không xác định"
	}

	// Xác định loại khóa
	var lockType domain.AccountStatus
	if userInfo.IsTemporarilyLocked() {
		lockType = domain.StatusTemporarilyLocked
	} else if userInfo.IsPermanentlyLocked() {
		lockType = domain.StatusPermanentlyLocked
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
	if req.LockType == domain.StatusTemporarilyLocked && req.Duration <= 0 {
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
	if req.LockType == domain.StatusTemporarilyLocked {
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
		return nil, _errors.ReturnError(
			int32(401),
			"Không thể xác thực phiên đăng nhập",
		)
	}

	// Kiểm tra password và confirmPassword có khớp không
	if req.Password != req.ConfirmPassword {
		return nil, _errors.ReturnError(
			int32(400),
			"Mật khẩu xác nhận không khớp",
		)
	}

	// Kiểm tra độ dài password (tối thiểu 8 ký tự)
	if len(req.Password) < 8 {
		return nil, _errors.ReturnError(
			int32(400),
			"Mật khẩu phải có ít nhất 8 ký tự",
		)
	}

	// Kiểm tra user đã có auth_method với provider ADMIN chưa
	existingAuth, err := s.AuthMethodRepo.GetByUserIdAndProvider(ctx, profileId, "ADMIN")
	if err == nil && existingAuth != nil {
		return nil, _errors.ReturnError(
			int32(400),
			"Bạn đã có mật khẩu. Vui lòng sử dụng chức năng đổi mật khẩu nếu muốn thay đổi",
		)
	}

	// Hash password
	hashedPassword, err := _utils.HashPassword(req.Password)
	if err != nil {
		return nil, _errors.ReturnError(
			int32(500),
			"Lỗi khi xử lý mật khẩu",
		)
	}

	// Lấy thông tin profile để có username
	profile, err := s.ProfileProvider.GetByProfileID(ctx, profileId)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(
			int32(404),
			"Không tìm thấy thông tin người dùng",
		)
	}

	// Tạo auth_method mới với provider ADMIN
	newAuthMethod := &domain.AuthMethod{
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
		return nil, _errors.ReturnError(
			int32(500),
			"Lỗi khi tạo mật khẩu",
		)
	}

	// Tạo status và OTP cho auth_method mới (nếu cần)
	statusEntity := &domain.UserStatusEntity{
		Active:   true,
		Verified: true, // Đã verified vì user đã đăng nhập
		IdDomain: domain.IdDomain{
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
		return nil, _errors.ReturnError(
			int32(400),
			"Username và password không được để trống",
		)
	}

	// Tìm auth_method theo username (có thể là phone)
	// Thử tìm theo phone trước
	var authEntity *domain.AuthMethod
	var err error

	// Kiểm tra xem username có phải là số điện thoại không
	if _utils.ValidatePhoneNumber(req.Username) {
		// Tìm theo phone với provider ADMIN
		authEntity, err = s.AuthMethodRepo.FindByPhone(ctx, req.Username)
		if err != nil || authEntity == nil {
			return nil, _errors.ReturnError(
				int32(401),
				"Số điện thoại hoặc mật khẩu không đúng",
			)
		}
		// Kiểm tra có auth_method với provider ADMIN không
		if authEntity.Provider != "ADMIN" {
			// Tìm auth_method ADMIN của cùng user
			if authEntity.UserID != 0 {
				authEntity, err = s.AuthMethodRepo.GetByUserIdAndProvider(ctx, authEntity.UserID, "ADMIN")
				if err != nil || authEntity == nil {
					return nil, _errors.ReturnError(
						int32(401),
						"Bạn chưa thiết lập mật khẩu. Vui lòng đăng nhập bằng OTP",
					)
				}
			} else {
				return nil, _errors.ReturnError(
					int32(401),
					"Bạn chưa thiết lập mật khẩu. Vui lòng đăng nhập bằng OTP",
				)
			}
		}
	} else {
		// Tìm theo username với provider ADMIN
		authEntity, err = s.AuthMethodRepo.FindByAuthNameAndProvider(ctx, req.Username, "ADMIN")
		if err != nil || authEntity == nil {
			return nil, _errors.ReturnError(
				int32(401),
				"Username hoặc mật khẩu không đúng",
			)
		}
	}

	// Kiểm tra provider phải là ADMIN
	if authEntity.Provider != "ADMIN" {
		return nil, _errors.ReturnError(
			int32(401),
			"Phương thức đăng nhập không đúng",
		)
	}

	// Kiểm tra password
	if !_utils.CheckPasswordHash(req.Password, authEntity.Password) {
		return nil, _errors.ReturnError(
			int32(401),
			"Username hoặc mật khẩu không đúng",
		)
	}

	// Lấy profile của user
	profile, err := s.ProfileProvider.GetByProfileID(ctx, authEntity.UserID)
	if err != nil || profile == nil || profile.ProfileID == 0 {
		return nil, _errors.ReturnError(
			int32(404),
			"Không tìm thấy thông tin người dùng",
		)
	}

	// Check account status - nếu tài khoản bị khóa thì logout và báo lỗi
	if authEntity.UserID != 0 {
		isLocked, canLogin, err := s.ProfileProvider.CheckAccountStatus(ctx, authEntity.UserID)
		if err != nil {
			log.Printf("Failed to check account status: %v", err)
		}
		if isLocked || !canLogin {
			// Logout tài khoản nếu bị khóa
			s.CacheProvider.DeleteToken(ctx, authEntity.ID)
			return nil, _errors.ReturnError(
				int32(403),
				"Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên để được hỗ trợ",
			)
		}
	}

	// Tạo session
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
		Fullname:   authEntity.FullName,
	}

	newSession := s.createSession(ctx, authEntity, otpRequest)
	s.SessionRepo.CreateSession(ctx, newSession)

	// Tạo refresh token
	refreshToken := s.GenRefreshToken(ctx, authEntity, newSession.SessionID)

	// Lấy passwordAuthId (chính là auth_method hiện tại)
	response := s.ResponseLogin(authEntity, profile, newSession.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	response.PasswordAuthId = authEntity.ID

	// Ghi log admin history
	go func() {
		// Tạo context mới với timeout riêng, không dùng context gốc
		contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Lấy IP address và User Agent từ context gốc trước khi vào goroutine
		ipAddress := _utils.AuditIPFromContext(ctx)
		userAgent := _utils.GetUserAgentFromContext(ctx)

		// Nếu không có IP hoặc User Agent, sử dụng giá trị mặc định
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		if userAgent == "" {
			userAgent = "Admin-Client"
		}

		log.Printf("Creating admin history for user %d, IP: %s, UserAgent: %s",
			authEntity.UserID, ipAddress, userAgent)

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
			log.Printf("Failed to log admin login: %v", err)
		} else {
			log.Printf("Successfully logged admin login for user %d", authEntity.UserID)
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
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể xác thực phiên đăng nhập",
		}
	}

	// Kiểm tra quyền truy cập profile mới
	if req.ProfileID != currentProfileId {
		// TODO: Kiểm tra quyền truy cập profile khác (có thể là admin hoặc có quyền đặc biệt)
		// Hiện tại chỉ cho phép chuyển đổi trong cùng profile
		return nil, &_routes.Except{
			Code:    403,
			Message: "Không có quyền truy cập profile này",
		}
	}

	// Lấy thông tin profile
	profile, err := s.ProfileProvider.GetByProfileID(c, req.ProfileID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy profile",
		}
	}

	// Lấy thông tin auth hiện tại (sử dụng provider OTP làm mặc định)
	authEntity, err := s.AuthMethodRepo.GetByUserIdAndProvider(c, req.ProfileID, "PHONE")
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy thông tin xác thực",
		}
	}

	// Tạo JWT token mới với thông tin chuyển đổi
	payload := _jwt.JwtTokenProperties{
		AuthID:    authId,
		ProfileID: currentProfileId,
		SessionID: sessionId,
		Role:      "ROLE_USER",
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
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi tạo token",
		}
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

// GetLoginHistory lấy danh sách lịch sử đăng nhập của user
func (s *AuthUsecase) GetLoginHistory(ctx context.Context, page, size int32) ([]*domain.UserSessionEntity, int64, error) {
	// Lấy profileId từ context
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, 0, _errors.ReturnError(
			int32(401),
			"Không tìm thấy thông tin người dùng",
		)
	}

	// Lấy danh sách lịch sử đăng nhập từ repository (JOIN với auth_method)
	sessions, total, err := s.SessionRepo.GetLoginHistoryByUserID(ctx, profileID, page, size)
	if err != nil {
		return nil, 0, _errors.ReturnError(
			int32(500),
			"Không thể lấy lịch sử đăng nhập: "+err.Error(),
		)
	}

	return sessions, total, nil
}
