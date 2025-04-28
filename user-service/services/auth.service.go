package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"auth/config"
	"auth/data"
	"auth/enums"
	"auth/models"
	"auth/validator"
	_db "common/db"
	_jwt "common/jwt"
	_redis "common/redis"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	// OAuthRepo        repository.UserAuthRepository
	// SessionRepo      repository.UserSessionRepository
	// OTPRepo          repository.UserOTPRepository
	// StatusRepo       repository.UserStatusRepository
	// InfoRepo         repository.UserInfoRepository
	// JWTProvider      utils.JwtProvider
	// AuthValidator    utils.AuthValidator
	// AuditService     utils.AuditService
	RedisService  *_redis.RedisService
	CookieService *CookieService

	// ConfigProperties utils.PropertiesConfiguration
	// Properties config.AppProperties
}

func NewAuthService(redisService *_redis.RedisService, cookie *CookieService) *AuthService {
	return &AuthService{
		RedisService:  redisService,
		CookieService: cookie,
	}
}

func (s *AuthService) RequestOtp(c *gin.Context, otp data.OtpRequest) (*data.RequestLoginResponse, error) {
	valid := _utils.ValidatePhoneNumber(otp.Phone)
	if !valid {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số điện thoại không đúng",
		}
	}
	oauth, err := models.AuthFindByPhone(otp.Phone)
	if err == nil && oauth != nil {
		param, _ := GetParamValidate(oauth, "")
		err = validator.ValidateAuth(c,
			[]enums.AuthCodeEnum{
				enums.LOCKED,
				enums.LIMIT_REQUEST_TIME,
			}, param)
		if err != nil {
			return nil, err
		}
		s.SendNewOTP(c, param.OTP, otp.Phone) // todo
	} else {
		oauth, err = s.RequestIsRegister(c, otp)
		if err != nil {
			return nil, err
		}
	}
	return &data.RequestLoginResponse{
		AuthID:   oauth.ID,
		Provider: "PHONE",
		OAuthID:  otp.Phone,
	}, nil
}

func (s *AuthService) ResendOTP(c *gin.Context, otp data.OtpRequest) (any, error) {
	oauth, err := models.AuthFindByPhone(otp.Phone)
	if err != nil {
		return nil, err
	}
	param, err := GetParamValidate(oauth, "")
	if err != nil {
		return nil, err
	}
	err = validator.ValidateAuth(c, []enums.AuthCodeEnum{
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
	return nil, nil
}

func (s *AuthService) VerifyOtp(c *gin.Context, otp data.OtpVerifyRequest) (*data.AuthLoginResponse, error) {
	e, err := models.AuthFindByPhone(otp.Phone)
	if err != nil {
		return nil, err
	}
	param, _ := GetParamValidate(e, otp.Otp)
	err = validator.ValidateAuth(c, []enums.AuthCodeEnum{
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
	_db.SaveWithAudit(c, param.OTP)

	param.Status.LockedUntil = time.Time{}
	param.Status.Verified = true
	_db.SaveWithAudit(c, param.Status)

	newSession := s.createSession(c, e, otp)
	_db.SaveWithAudit(c, newSession)

	infoEntity, err := models.ProfileGetByProfileId(e.UserID)
	if err != nil {
		infoEntity = models.NewUserProfileEntity(e.ID, e.Name, e.Email, e.Phone, e.Avatar)
		_db.SaveWithAudit(c, infoEntity)

		e.UserID = infoEntity.ProfileID
		_db.SaveWithAudit(c, e)
	}

	// Liên kết tài khoản OAuth với hồ sơ người dùng
	if otp.AuthID != 0 {
		var oauthAcc *models.UserAuthEntity
		if err := _db.DB.First(&oauthAcc, otp.AuthID).Error; err == nil {
			if oauthAcc != nil && oauthAcc.UserID == 0 {
				var statusOauth models.UserStatusEntity
				if err := _db.DB.First(&statusOauth, oauthAcc.ID).Error; err == nil {
					statusOauth.Verified = true
					_db.SaveWithAudit(c, statusOauth)
				}

				oauthAcc.UserID = infoEntity.ProfileID
				_db.SaveWithAudit(c, oauthAcc)

				infoEntity.Avatar = oauthAcc.Avatar
				infoEntity.Email = oauthAcc.Email
				_db.SaveWithAudit(c, infoEntity)
			}
		}
	}

	// s.RedisService.SaveToken(e.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), config.AppProperties.Jwt.RefreshExpMinutes*60)
	// refreshToken := s.GetRefreshToken(e, newSession.SessionID)
	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	s.GenRefreshToken(c, e, newSession.SessionID)

	return s.ResponseLogin(e, newSession.SessionID), nil
}

func (s *AuthService) createSession(c *gin.Context, e *models.UserAuthEntity, otp data.OtpVerifyRequest) *models.UserSessionEntity {
	return &models.UserSessionEntity{
		AuthID:    e.ID,
		Version:   otp.Version,
		Platform:  otp.Platform,
		OS:        otp.OS,
		IPRequest: _utils.AuditIP(c),
	}
}

func (s *AuthService) RefreshToken(refreshToken string) (*data.RefreshResponse, error) {
	log.Println("refreshToken: %s", refreshToken, ":???")
	props, err := _jwt.GetProperties(refreshToken)
	if err != nil {
		return nil, &_routes.Except{
			Code:    401,
			Message: err.Error(),
		}
	}

	date, _ := s.RedisService.GetToken(props.AuthID)
	invalidToken := _utils.ValidTokenByIssueAt(props.IssuedAt, date)
	if invalidToken {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Phiên làm việc đã hết hạn vui lòng đăng nhập lại",
		}
	}

	if props.Type == "REFRESH" {
		newToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID:    props.AuthID,
			ProfileID: props.ProfileID,
			Role:      props.Role,
			SessionID: props.SessionID,
			Type:      "ACCESS",
		})
		return &data.RefreshResponse{AccessToken: newToken}, nil
	}

	return nil, &_routes.Except{
		Code:    401,
		Message: "Token không đúng",
	}
}

func (s *AuthService) Logout(c *gin.Context) (any, error) {
	sessionID := _jwt.GetSessionId(c)
	currentSession, err := models.SessionGetBySessionId(sessionID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	currentSession.FinishedDate = &now
	_db.SaveWithAudit(c, currentSession)
	s.RedisService.DeleteToken(currentSession.AuthID)
	return 0, nil
}

// Restore khôi phục tài khoản
func (s *AuthService) Restore(c *gin.Context, dto data.RestorePhoneRequest) (*data.RequestLoginResponse, error) {
	ok, err := models.AuthExistsByOAuth(dto.Phone)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số điện thoại đã liên kết với tài khoản khác. Vui lòng nhập số điện thoại khác",
		}
	}

	profileId := _jwt.GetProfileId(c)

	e := models.UserAuthEntity{
		Provider: "PHONE",
		OAuthID:  dto.Phone,
		UserID:   profileId,
	}

	if err := _db.SaveWithAudit(c, &e); err != nil {
		return nil, err
	}

	param, _ := s.NewUser(c, e.ID)
	if err := _db.SaveWithAudit(c, param.OTP); err != nil {
		return nil, err
	}

	return &data.RequestLoginResponse{
		AuthID:   e.ID,
		Provider: e.Provider,
		OAuthID:  e.OAuthID,
	}, nil
}

// Xóa tài khoản
func (s *AuthService) Delete(c *gin.Context) (uint64, error) {
	// profileId := _jwt.GetProfileId(c)
	principal := _jwt.GetPrincipal(c)
	profile, err := models.ProfileGetByProfileId(principal.ProfileId)
	if err != nil {
		return 0, err
	}
	if profile == nil {
		return 0, &_routes.Except{
			Code:    400,
			Message: "Tài khoản không tồn tại",
		}
	}

	t := time.Now()
	profile.DeletedAt = &t

	if err := _db.SaveWithAudit(c, profile); err != nil {
		return 0, err
	}

	models.DeleteProfile(principal.ProfileId)

	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		c.SetCookie(cookie.Name, "", -1, "/", "", false, true)
	}
	s.RedisService.DeleteToken(principal.UserID)

	return profile.ProfileID, nil
}

func (s *AuthService) RequestIsRegister(c *gin.Context, dto data.OtpRequest) (*models.UserAuthEntity, error) {
	// Tạo một UserAuthEntity mới
	e := &models.UserAuthEntity{
		Provider: "PHONE",
		OAuthID:  dto.Phone,
		Name:     dto.Fullname,
	}

	// Lưu vào database
	if err := _db.SaveWithAudit(c, e); err != nil {
		return nil, err
	}

	// Tạo AuthParam mới
	param, err := s.NewUser(c, e.ID)
	if err != nil {
		return nil, err
	}

	// Gửi OTP trong updateNewOTP luôn
	s.SendNewOTP(c, param.OTP, dto.Phone)

	return e, nil
}
func (s *AuthService) NewUser(c *gin.Context, oauthId uint64) (*data.AuthParam, error) {
	// Tạo UserStatusEntity mới
	statusEntity := &models.UserStatusEntity{
		Active:   true,
		Verified: false,
		UserIdEntity: models.UserIdEntity{
			AuthID: oauthId,
		},
	}

	// Lưu vào database
	if err := _db.SaveWithAudit(c, statusEntity); err != nil {
		return nil, err
	}

	// Tạo UserOTPEntity mới
	otpEntity := &models.UserOTPEntity{
		OTP: _utils.GenerateOTP(),
		UserIdEntity: models.UserIdEntity{
			AuthID: oauthId,
		},
		// AuthID:       oauthId,
		OTPSendTime:  0,
		OTPCheckTime: 0,
	}

	return &data.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: "",
	}, nil
}

func (s *AuthService) SendNewOTP(c *gin.Context, otpEntity *models.UserOTPEntity, phone string) string {
	ip := _utils.AuditIP(c)
	s.CheckOtpSpam(phone, ip)

	otp := _utils.GenerateOTP()
	otpEntity.OTPSendTime += 1
	otpEntity.OTP = otp
	otpEntity.Activate = false
	t := time.Now().Add(time.Duration(config.AppProperties.Otp.ExpiredAfterSec) * time.Second)
	otpEntity.ExpiredTime = &t

	_db.SaveWithAudit(c, otpEntity)

	// TODO: Gọi sang notification để gửi OTP

	return otp
}

func (s *AuthService) CheckOtpSpam(phoneNumber, ip string) error {
	redisKeyPhone := fmt.Sprintf("otp:request:%s", phoneNumber)
	redisKeyIp := fmt.Sprintf("otp:ip:%s", ip)

	// Kiểm tra giới hạn gửi OTP theo số điện thoại
	requestCountPhone, err := s.RedisService.Get(redisKeyPhone)
	if err != nil {
		return err
	}
	countPhone, _ := strconv.ParseInt(requestCountPhone, 10, 32)
	if countPhone >= config.AppProperties.Otp.LimitOtpDevice {
		return &_routes.Except{
			Code:    401,
			Message: "Bạn đã gửi quá nhiều OTP. Vui lòng thử lại sau.",
		}
	}

	// Kiểm tra giới hạn gửi OTP theo IP
	requestCountIp, err := s.RedisService.Get(redisKeyIp)
	countIp, _ := strconv.ParseInt(requestCountIp, 10, 32)
	if err != nil {
		return err
	}
	if countIp >= config.AppProperties.Otp.LimitOtpIP {
		return &_routes.Except{
			Code:    401,
			Message: "IP của bạn đã bị chặn do gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
		}
	}

	// Tăng số lần gửi OTP trong Redis
	if err := s.RedisService.Increment(redisKeyPhone); err != nil {
		return err
	}
	if err := s.RedisService.Expire(redisKeyPhone, config.AppProperties.Otp.CooldownSeconds); err != nil {
		return err
	}

	if err := s.RedisService.Increment(redisKeyIp); err != nil {
		return err
	}
	if err := s.RedisService.Expire(redisKeyIp, config.AppProperties.Otp.IpCooldownSeconds); err != nil {
		return err
	}

	return nil
}

// GetParamValidate kiểm tra thông tin xác thực và trả về data.AuthParam
func GetParamValidate(oauth *models.UserAuthEntity, otp string) (*data.AuthParam, error) {
	if oauth == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng, vui lòng kiểm tra lại",
		}
	}

	statusEntity, err := models.StatusGetById(oauth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy trạng thái người dùng",
		}
	}

	otpEntity, err := models.OTPGetById(oauth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy OTP",
		}
	}

	return &data.AuthParam{
		Status:  statusEntity,
		OTP:     otpEntity,
		OTPCode: otp,
	}, nil
}

func (s *AuthService) ResponseLogin(authEntity *models.UserAuthEntity, sessionID uint64) *data.AuthLoginResponse {
	if authEntity.UserID == 0 {
		return &data.AuthLoginResponse{
			AuthID: authEntity.ID,
		}
	}
	accessToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
		SessionID: sessionID,
		Role:      "ROLE_USER",
		Type:      "ACCESS",
	})

	return data.NewAuthLoginResponse(accessToken, authEntity, "ROLE_USER")
}

func (s *AuthService) GetRefreshToken(authEntity *models.UserAuthEntity, sessionID uint64) string {
	refreshToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
		SessionID: sessionID,
		Role:      "ROLE_USER",
		Type:      "REFRESH",
	})
	return refreshToken
}

func (s *AuthService) CreateSession(c *gin.Context, e *models.UserAuthEntity, otp *data.OtpVerifyRequest) *models.UserSessionEntity {
	var authId uint64
	if e != nil {
		authId = e.ID
	}
	if otp == nil {
		otp = &data.OtpVerifyRequest{}
	}
	return &models.UserSessionEntity{
		AuthID:    authId,
		Version:   otp.Version,
		Platform:  otp.Platform,
		OS:        otp.OS,
		IPRequest: _utils.AuditIP(c),
		Activate:  false,
	}
}

func (s *AuthService) CreateQRSession(c *gin.Context, e *data.SessionQRRequest, sessionId uint64) (*data.SessionQRResponse, error) {
	session, _ := models.SessionGetBySessionId(sessionId)
	if session == nil {
		session = s.CreateSession(c, nil, nil)
	} else {

		// copier.Copy(&session, &e)
		// session.Activate
		// e.SessionId = sessionId
		return &data.SessionQRResponse{
			SessionKey: session.SessionKey,
			SessionId:  session.SessionID,
		}, nil
	}

	key := _utils.RandomString(24)
	session.SessionKey = key

	err := _db.SaveWithAudit(c, &session)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	s.CookieService.Set(c, config.C_SESSION_ID, strconv.FormatUint(session.SessionID, 10))

	return &data.SessionQRResponse{
		SessionKey: key,
		SessionId:  session.SessionID,
	}, nil
}

const maxConcurrentRequests = 1000 // Giới hạn số request đồng thời

var semaphore = make(chan struct{}, maxConcurrentRequests)

func (s *AuthService) ConfirmQRSession(c *gin.Context, e *data.SessionQRConfirm) (*models.UserSessionEntity, error) {
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

	props := _jwt.GetPrincipal(c)
	if props == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể xác thực phiên đăng nhập",
		}
	}

	sessionEntity, err := models.SessionGetBySessionId(e.SessionId)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	if sessionEntity.FinishedDate != nil || sessionEntity.AuthID != 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Phiên đăng nhập không chính xác",
		}
	}

	sessionEntity.AuthID = props.UserID
	_db.SaveWithAudit(c, sessionEntity)

	return sessionEntity, nil
}

func (s *AuthService) VerifyQRSession(c *gin.Context, e *data.SessionQRRequest, sessionId uint64) (*data.AuthLoginResponse, error) {
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

	sessionEntity, err := models.SessionGetBySessionId(sessionId)
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
	_db.SaveWithAudit(c, sessionEntity)

	authEntity, err := models.AuthGetWithId(sessionEntity.AuthID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	// s.RedisService.SaveToken(authEntity.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), config.AppProperties.Jwt.RefreshExpMinutes*60)
	// refreshToken := s.GetRefreshToken(authEntity, sessionId)
	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	s.GenRefreshToken(c, authEntity, sessionId)
	return s.ResponseLogin(authEntity, sessionEntity.SessionID), nil
}

func (s *AuthService) GenRefreshToken(c *gin.Context, authEntity *models.UserAuthEntity, sessionId uint64) string {
	s.RedisService.SaveToken(authEntity.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), config.AppProperties.Jwt.RefreshExpMinutes*60)
	refreshToken := s.GetRefreshToken(authEntity, sessionId)
	s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	return refreshToken
}
