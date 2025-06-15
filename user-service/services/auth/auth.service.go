package auth_service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	_db "common/db"
	_jwt "common/jwt"
	_redis "common/redis"
	_routes "common/routes"
	_utils "common/utils"
	"user/config/property"
	"user/data"
	"user/enums"
	auth_model "user/models/auth"
	user_model "user/models/user"
	"user/repo"
	"user/validator"

	organizationpb "pb/types/organization"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	RedisService       *_redis.RedisService
	CookieService      *CookieService
	PreePlanID         *uint64
	ProfileRepo        *repo.ProfileRepo
	OrganizationClient organizationpb.OrganizationMemberServiceClient
	// ConfigProperties utils.PropertiesConfiguration
	// Properties property.AppProperties
}

// NewRedisService khởi tạo RedisService
func NewRedisService() *_redis.RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"), // Địa chỉ Redis (ví dụ: "localhost:6379")
		Password: viper.GetString("redis.pass"),
	})
	return &_redis.RedisService{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func NewOrganizationClient() organizationpb.OrganizationMemberServiceClient {
	conn, err := grpc.NewClient(viper.GetString("grpc.organization.host"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to organization service: %v", err)
	}
	return organizationpb.NewOrganizationMemberServiceClient(conn)
}

func NewAuthService(redisService *_redis.RedisService,
	cookie *CookieService,
	profileRepo *repo.ProfileRepo,
) *AuthService {
	planId := uint64(1)
	return &AuthService{
		RedisService:       redisService,
		CookieService:      cookie,
		PreePlanID:         &planId,
		ProfileRepo:        profileRepo,
		OrganizationClient: NewOrganizationClient(),
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
	oauth, err := auth_model.AuthFindByPhone(otp.Phone)
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
	oauth, err := auth_model.AuthFindByPhone(otp.Phone)
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
	e, err := auth_model.AuthFindByPhone(otp.Phone)
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

	infoEntity, err := s.ProfileRepo.GetByProfileID(e.UserID)
	if err != nil || infoEntity == nil || infoEntity.ProfileID == 0 {
		infoEntity = user_model.MakeUserProfileEntity(e.ID, s.PreePlanID, e.Name, e.Email, e.Phone, e.Avatar)
		_db.SaveWithAudit(c, infoEntity)

		e.UserID = infoEntity.ProfileID
		_db.SaveWithAudit(c, e)
	}

	// Liên kết tài khoản OAuth với hồ sơ người dùng
	if otp.AuthID != 0 {
		var oauthAcc *auth_model.UserAuthEntity
		if err := _db.DB.First(&oauthAcc, otp.AuthID).Error; err == nil {
			if oauthAcc != nil && oauthAcc.UserID == 0 {
				var statusOauth auth_model.UserStatusEntity
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

	// s.RedisService.SaveToken(e.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), property.AppProperties.Jwt.RefreshExpMinutes*60)
	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	s.GenRefreshToken(c, e, newSession.SessionID)

	return s.ResponseLogin(e, infoEntity, newSession.SessionID), nil
}

func (s *AuthService) createSession(c *gin.Context, e *auth_model.UserAuthEntity, otp data.OtpVerifyRequest) *auth_model.UserSessionEntity {
	return &auth_model.UserSessionEntity{
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
		profile := s.GetProfileByID(props.ProfileID)
		newToken, _ := _jwt.GenerateToken(_jwt.JwtTokenProperties{
			AuthID:    props.AuthID,
			ProfileID: props.ProfileID,
			Role:      props.Role,
			SessionID: props.SessionID,
			PlanID:    profile.PlanID,
			PlanAt:    profile.PlanAt,
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
	currentSession, err := auth_model.SessionGetBySessionId(sessionID)
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
	ok, err := auth_model.AuthExistsByOAuth(dto.Phone)
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

	e := auth_model.UserAuthEntity{
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
	profile, err := s.ProfileRepo.GetByProfileID(principal.ProfileId)
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

	auth_model.DeleteProfile(principal.ProfileId)

	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		c.SetCookie(cookie.Name, "", -1, "/", "", false, true)
	}
	s.RedisService.DeleteToken(principal.UserID)

	return profile.ProfileID, nil
}

func (s *AuthService) RequestIsRegister(c *gin.Context, dto data.OtpRequest) (*auth_model.UserAuthEntity, error) {
	// Tạo một UserAuthEntity mới
	e := &auth_model.UserAuthEntity{
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
	statusEntity := &auth_model.UserStatusEntity{
		Active:   true,
		Verified: false,
		UserIdEntity: auth_model.UserIdEntity{
			AuthID: oauthId,
		},
	}

	// Lưu vào database
	if err := _db.SaveWithAudit(c, statusEntity); err != nil {
		return nil, err
	}

	// Tạo UserOTPEntity mới
	otpEntity := &auth_model.UserOTPEntity{
		OTP: _utils.GenerateOTP(),
		UserIdEntity: auth_model.UserIdEntity{
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

func (s *AuthService) SendNewOTP(c *gin.Context, otpEntity *auth_model.UserOTPEntity, phone string) string {
	ip := _utils.AuditIP(c)
	s.CheckOtpSpam(phone, ip)

	otp := _utils.GenerateOTP()
	otpEntity.OTPSendTime += 1
	otpEntity.OTP = otp
	otpEntity.Activate = false
	t := time.Now().Add(time.Duration(property.AppProperties.Otp.ExpiredAfterSec) * time.Second)
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
	if countPhone >= property.AppProperties.Otp.LimitOtpDevice {
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
	if countIp >= property.AppProperties.Otp.LimitOtpIP {
		return &_routes.Except{
			Code:    401,
			Message: "IP của bạn đã bị chặn do gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
		}
	}

	// Tăng số lần gửi OTP trong Redis
	if err := s.RedisService.Increment(redisKeyPhone); err != nil {
		return err
	}
	if err := s.RedisService.Expire(redisKeyPhone, property.AppProperties.Otp.CooldownSeconds); err != nil {
		return err
	}

	if err := s.RedisService.Increment(redisKeyIp); err != nil {
		return err
	}
	if err := s.RedisService.Expire(redisKeyIp, property.AppProperties.Otp.IpCooldownSeconds); err != nil {
		return err
	}

	return nil
}

// GetParamValidate kiểm tra thông tin xác thực và trả về data.AuthParam
func GetParamValidate(oauth *auth_model.UserAuthEntity, otp string) (*data.AuthParam, error) {
	if oauth == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Thông tin không đúng, vui lòng kiểm tra lại",
		}
	}

	statusEntity, err := auth_model.StatusGetById(oauth.ID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không tìm thấy trạng thái người dùng",
		}
	}

	otpEntity, err := auth_model.OTPGetById(oauth.ID)
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

func (s *AuthService) ResponseLogin(authEntity *auth_model.UserAuthEntity, profile *user_model.UserProfileEntity, sessionID uint64) *data.AuthLoginResponse {
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
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	})

	return data.MakeAuthLoginResponse(accessToken, authEntity, "ROLE_USER")
}

func (s *AuthService) GetRefreshToken(authEntity *auth_model.UserAuthEntity, sessionID uint64) string {
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

func (s *AuthService) CreateSession(c *gin.Context, e *auth_model.UserAuthEntity, otp *data.OtpVerifyRequest) *auth_model.UserSessionEntity {
	var authId uint64
	if e != nil {
		authId = e.ID
	}
	if otp == nil {
		otp = &data.OtpVerifyRequest{}
	}
	return &auth_model.UserSessionEntity{
		AuthID:    authId,
		Version:   otp.Version,
		Platform:  otp.Platform,
		OS:        otp.OS,
		IPRequest: _utils.AuditIP(c),
		Activate:  false,
	}
}

func (s *AuthService) CreateQRSession(c *gin.Context, e *data.SessionQRRequest, sessionId uint64) (*data.SessionQRResponse, error) {
	session, _ := auth_model.SessionGetBySessionId(sessionId)
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

	s.CookieService.Set(c, property.C_SESSION_ID, strconv.FormatUint(session.SessionID, 10))

	return &data.SessionQRResponse{
		SessionKey: key,
		SessionId:  session.SessionID,
	}, nil
}

const maxConcurrentRequests = 1000 // Giới hạn số request đồng thời

var semaphore = make(chan struct{}, maxConcurrentRequests)

func (s *AuthService) ConfirmQRSession(c *gin.Context, e *data.SessionQRConfirm) (*auth_model.UserSessionEntity, error) {
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

	sessionEntity, err := auth_model.SessionGetBySessionId(e.SessionId)
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

	sessionEntity, err := auth_model.SessionGetBySessionId(sessionId)
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

	authEntity, err := auth_model.AuthGetWithId(sessionEntity.AuthID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: err.Error(),
		}
	}

	// s.RedisService.SaveToken(authEntity.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), property.AppProperties.Jwt.RefreshExpMinutes*60)
	// s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	s.GenRefreshToken(c, authEntity, sessionId)
	profile := s.GetProfileByID(authEntity.UserID)
	return s.ResponseLogin(authEntity, profile, sessionEntity.SessionID), nil
}

func (s *AuthService) GetPlanID(profileID uint64) (*uint64, *time.Time) {
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileID)
	if err != nil {
		return nil, nil
	}
	return infoEntity.PlanID, infoEntity.PlanAt
}

func (s *AuthService) GetProfileByID(profileID uint64) *user_model.UserProfileEntity {
	infoEntity, err := s.ProfileRepo.GetByProfileID(profileID)
	if err != nil {
		return nil
	}
	return infoEntity
}

func (s *AuthService) GenRefreshToken(c *gin.Context, authEntity *auth_model.UserAuthEntity, sessionId uint64) string {
	s.RedisService.SaveToken(authEntity.ID, time.Now().Add(-5*time.Second).Format(time.RFC3339), property.AppProperties.Jwt.RefreshExpMinutes*60)
	refreshToken := s.GetRefreshToken(authEntity, sessionId)
	s.CookieService.SetRefreshToken(c.Writer, refreshToken)
	log.Println("refreshToken=", refreshToken)
	return refreshToken
}

func (s *AuthService) SwitchOrganization(c *gin.Context, organizationID uint32) (*data.AuthLoginResponse, error) {
	principal := _jwt.GetPrincipal(c)
	if principal == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể xác thực phiên đăng nhập",
		}
	}
	clientResponse, err := s.OrganizationClient.SwitchOrganization(c, &organizationpb.SwitchOrganizationRequest{
		OrganizationId: organizationID,
		UserId:         uint32(principal.ProfileId),
	})
	fmt.Printf("clientResponse=%+v\n", clientResponse)
	if err != nil {
		return nil, err
	}
	if !clientResponse.IsSuccess {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể chuyển đổi tổ chức",
		}
	}
	authEntity, err := auth_model.AuthGetWithId(principal.UserID)
	if err != nil {
		return nil, err
	}
	profile := s.GetProfileByID(principal.UserID)
	payload := _jwt.JwtTokenProperties{
		AuthID:    principal.UserID,
		ProfileID: principal.UserID,
		SessionID: uint64(principal.Session),
		Role:      "ROLE_USER",
		Type:      "ACCESS",
		PlanID:    profile.PlanID,
		PlanAt:    profile.PlanAt,
	}
	payload.AdditionalInfo = map[string]interface{}{
		"organization_id": clientResponse.Id,
	}
	newAccessToken, err := _jwt.GenerateToken(payload)
	if err != nil {
		return nil, err
	}
	response := &data.AuthLoginResponse{
		AuthID:         principal.UserID,
		AccessToken:    newAccessToken,
		Role:           "ROLE_USER",
		Username:       profile.FullName,
		ProfileID:      profile.ProfileID,
		RefreshToken:   s.GetRefreshToken(authEntity, uint64(principal.Session)),
		OrganizationID: clientResponse.Id,
	}

	return response, nil
}
