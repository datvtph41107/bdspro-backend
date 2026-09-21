package oauth

import (
	_errors "common/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"user/config"
	"user/internal"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/interface/repo"
	"user/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ZaloAuthService struct {
	// apiClient   *ApiClient
	// properties  *AuthProperties
	// userRepo    *UserAuthRepository
	AuthService    *usecase.AuthUsecase
	HttpClient     *http.Client
	AuthMethodRepo repo.AuthMethodRepository
}

func NewZaloAuthService(
	httpClient *http.Client,
	authService *usecase.AuthUsecase,
	authMethodRepo repo.AuthMethodRepository,
) *ZaloAuthService {
	return &ZaloAuthService{
		// apiClient:   apiClient,
		// properties:  properties,
		// userRepo:    userRepo,
		HttpClient:     httpClient,
		AuthService:    authService,
		AuthMethodRepo: authMethodRepo,
	}
}

func (s *ZaloAuthService) GetUrlLogin() string {
	zaloAuthUrl := config.Properties.OAuth.Zalo.PermissionURL +
		"app_id=" + config.Properties.OAuth.Zalo.AppID +
		"&redirect_uri=" + url.QueryEscape(config.Properties.OAuth.Zalo.RedirectURI) +
		"&state=xyz"
	return zaloAuthUrl
}

// func (s *ZaloAuthService) HandleCallback(code string, state string, request *http.Request) (map[string]interface{}, error) {
// 	params := url.Values{}
// 	params.Set("app_id", property.AppProperties.Zalo.App.ID)
// 	params.Set("code", code)
// 	params.Set("grant_type", "authorization_code")

// 	headers := map[string]string{
// 		"secret_key": property.AppProperties.Zalo.SecretKey,
// 	}

// 	resp, err := s.apiClient.RequestTo("https://oauth.zaloapp.com/v4/access_token", http.MethodPost, params, headers)
// 	if err != nil {
// 		return nil, errors.New("Không lấy được access_token từ Zalo")
// 	}

// 	accessToken, ok := resp["access_token"].(string)
// 	if !ok {
// 		return nil, errors.New("Invalid access_token response")
// 	}

// 	userInfo, err := s.fetchUserInfo(accessToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.registerOrLogin(userInfo, request)
// }

// func (s *ZaloAuthService) fetchUserInfo(accessToken string) (map[string]interface{}, error) {
// 	userInfoUrl := "https://graph.zalo.me/v2.0/me?fields=id,name,picture,is_sensitive&access_token=" + accessToken
// 	resp, err := s.apiClient.RequestTo(userInfoUrl, http.MethodGet, nil, nil)
// 	if err != nil {
// 		return nil, errors.New("Không lấy được thông tin user từ Zalo")
// 	}
// 	return resp, nil
// }

func (s *ZaloAuthService) HandleCallback(c *gin.Context, code string, state string) (*dto.AuthLoginResponse, error) {
	accessToken, err := s.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return s.RegisterOrLogin(c, userInfo, nil)
}

func (s *ZaloAuthService) HandleCallbackWithDeviceInfo(c context.Context, code string, state string, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	accessToken, err := s.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return s.RegisterOrLogin(c, userInfo, otpReq)
}

func (s *ZaloAuthService) getAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("app_id", config.Properties.OAuth.Zalo.AppID)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", "https://oauth.zaloapp.com/v4/access_token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("secret_key", config.Properties.OAuth.Zalo.SecretKey)

	// Add data to body
	req.URL.RawQuery = data.Encode()

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("không lấy được access_token từ Zalo: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		if errorMsg, ok := result["error_description"].(string); ok {
			return "", fmt.Errorf("Zalo OAuth error: %s", errorMsg)
		}
		return "", errors.New("không lấy được access_token từ Zalo")
	}
	return accessToken, nil
}

func (s *ZaloAuthService) fetchUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://graph.zalo.me/v2.0/me?fields=id,name,picture,is_sensitive&access_token="+accessToken, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *ZaloAuthService) RegisterOrLogin(c context.Context, userInfo map[string]interface{}, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	username, ok := userInfo["id"].(string)
	if !ok {
		return nil, _errors.ReturnError(service.ZaloUserIDUnavailable)
	}

	fullName := ""
	if name, ok := userInfo["name"].(string); ok {
		fullName = name
	}

	avatar := extractAvatar(userInfo)
	isSensitive := false
	if isS, ok := userInfo["is_sensitive"].(bool); ok {
		isSensitive = isS
	}

	existingUser, err := s.AuthMethodRepo.FindByOAuthID(c, username, "ZALO")
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("lỗi khi tìm kiếm user: %v", err)
	}

	if existingUser == nil {
		e := &auth.AuthMethod{
			Provider:    "ZALO",
			AuthName:    username,
			FullName:    fullName,
			Avatar:      avatar,
			IsSensitive: isSensitive,
		}
		createdUser, err := s.AuthMethodRepo.Create(c, e)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo user mới: %v", err)
		}
		_, err = s.AuthService.NewUser(c, createdUser.ID)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo profile mới: %v", err)
		}
		existingUser = createdUser
	} else {
		// Cập nhật thông tin nếu có thay đổi
		updated := false
		if existingUser.FullName != fullName && fullName != "" {
			existingUser.FullName = fullName
			updated = true
		}
		if existingUser.Avatar != avatar && avatar != "" {
			existingUser.Avatar = avatar
			updated = true
		}
		if existingUser.IsSensitive != isSensitive {
			existingUser.IsSensitive = isSensitive
			updated = true
		}

		if updated {
			_, _ = s.AuthMethodRepo.Update(c, existingUser)
		}
	}

	// Nếu chưa có UserID, tạo profile mới (không bắt buộc số điện thoại)
	if existingUser.UserID == 0 {
		profile, err := s.AuthService.ProfileProvider.CreateProfile(c, existingUser.ID, nil, fullName, "", "", avatar)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo profile: %v", err)
		}

		// Cập nhật UserID cho AuthMethod
		existingUser.UserID = profile.ProfileID
		_, err = s.AuthMethodRepo.Update(c, existingUser)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi cập nhật UserID: %v", err)
		}
	}

	// Nếu không có otpReq thì tạo mặc định
	if otpReq == nil {
		otpReq = &dto.OtpVerifyRequest{}
	}

	session, err := s.AuthService.CreateAndSaveSession(c, existingUser, otpReq)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo session: %v", err)
	}

	profile, err := s.AuthService.ProfileProvider.GetByProfileID(c, existingUser.UserID)
	if err != nil || profile == nil {
		// Nếu không tìm thấy profile, thử tạo lại
		profile, err = s.AuthService.ProfileProvider.CreateProfile(c, existingUser.ID, nil, fullName, "", "", avatar)
		if err != nil {
			return nil, _errors.ReturnError(service.UserNotFoundAndCreateFailed)
		}
	}

	refreshToken := s.AuthService.GenRefreshToken(c, existingUser, session.SessionID)

	// Ghi log admin history
	s.AuthService.LogLoginHistory(c, existingUser, profile, "Đăng nhập hệ thống bằng Zalo OAuth")

	response := s.AuthService.ResponseLogin(c, existingUser, profile, session.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}

func extractAvatar(userInfo map[string]interface{}) string {
	if picture, ok := userInfo["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			if url, ok := data["url"].(string); ok {
				return url
			}
		}
	}
	return ""
}
