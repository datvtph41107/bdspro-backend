package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"user/config"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/interface/repo"
	"user/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FacebookAuthService struct {
	AuthService    *usecase.AuthUsecase
	HttpClient     *http.Client
	AuthMethodRepo repo.AuthMethodRepository
}

func NewFacebookAuthService(authService *usecase.AuthUsecase,
	httpClient *http.Client,
	authMethodRepo repo.AuthMethodRepository) *FacebookAuthService {
	return &FacebookAuthService{
		AuthService:    authService,
		HttpClient:     httpClient,
		AuthMethodRepo: authMethodRepo,
	}
}

func (s *FacebookAuthService) GetUrlLogin() string {
	fbAuthUrl := "https://www.facebook.com/v22.0/dialog/oauth?" +
		"client_id=" + config.Properties.OAuth.Facebook.ClientID +
		"&redirect_uri=" + url.QueryEscape(config.Properties.OAuth.Facebook.RedirectURI) +
		"&state=xyz" +
		"&scope=email"
	return fbAuthUrl
}

func (s *FacebookAuthService) HandleCallback(c *gin.Context, code string, platform string, state string) (*dto.AuthLoginResponse, error) {
	var accessToken string
	var err error

	// For Android platform, code is already the access token
	if platform == "android" {
		accessToken = code
	} else {
		// Get access token from authorization code
		accessToken, err = s.getAccessToken(code)
		if err != nil {
			return nil, fmt.Errorf("không thể lấy access token: %v", err)
		}
	}

	// Get user info using access token
	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin user: %v", err)
	}

	return s.RegisterOrLogin(c, userInfo, nil)
}

func (s *FacebookAuthService) HandleCallbackWithDeviceInfo(c context.Context, code string, platform string, state string, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	var accessToken string
	var err error

	// For Android platform, code is already the access token
	if platform == "android" {
		accessToken = code
	} else {
		// Get access token from authorization code
		accessToken, err = s.getAccessToken(code)
		if err != nil {
			return nil, fmt.Errorf("không thể lấy access token: %v", err)
		}
	}

	// Get user info using access token
	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin user: %v", err)
	}

	return s.RegisterOrLogin(c, userInfo, otpReq)
}

func (s *FacebookAuthService) getAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", config.Properties.OAuth.Facebook.ClientID)
	params.Set("redirect_uri", config.Properties.OAuth.Facebook.RedirectURI)
	params.Set("client_secret", config.Properties.OAuth.Facebook.ClientSecret)
	params.Set("code", code)

	tokenURL := "https://graph.facebook.com/v22.0/oauth/access_token?" + params.Encode()
	resp, err := s.HttpClient.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("không lấy được access_token từ Facebook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Facebook OAuth API trả về lỗi: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("lỗi đọc phản hồi từ Facebook: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("lỗi parse JSON từ Facebook: %v", err)
	}

	if errMsg, exists := result["error"]; exists {
		return "", fmt.Errorf("Facebook OAuth error: %v", errMsg)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("không lấy được access_token từ Facebook")
	}
	return accessToken, nil
}

func (s *FacebookAuthService) fetchUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email,picture&access_token="+accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo request: %v", err)
	}

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gọi Facebook Graph API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Facebook Graph API trả về lỗi: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc phản hồi từ Facebook: %v", err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("lỗi parse user info: %v", err)
	}

	return userInfo, nil
}

func (s *FacebookAuthService) RegisterOrLogin(c context.Context, userInfo map[string]interface{}, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	username, ok := userInfo["id"].(string)
	if !ok {
		return nil, errors.New("không lấy được id từ Facebook")
	}

	fullName := s.extractString(userInfo, "name")
	email := s.extractString(userInfo, "email")
	avatar := s.extractAvatar(userInfo)

	// Check if user already exists
	existingUser, err := s.AuthMethodRepo.FindByOAuthID(c, username, "FACEBOOK")
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("lỗi khi tìm kiếm user: %v", err)
	}

	if existingUser == nil {
		// User doesn't exist, create new one
		oauthUser := &auth.AuthMethod{
			Provider: "FACEBOOK",
			AuthName: username,
			FullName: fullName,
			Email:    email,
			Avatar:   avatar,
		}

		createdUser, err := s.AuthMethodRepo.Create(c, oauthUser)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo user mới: %v", err)
		}

		// Create user profile
		_, err = s.AuthService.NewUser(c, createdUser.ID)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo profile: %v", err)
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
		if existingUser.Email != email && email != "" {
			existingUser.Email = email
			updated = true
		}

		if updated {
			_, _ = s.AuthMethodRepo.Update(c, existingUser)
		}
	}

	// Nếu chưa có UserID, tạo profile mới (không bắt buộc số điện thoại)
	if existingUser.UserID == 0 {
		profile, err := s.AuthService.ProfileProvider.CreateProfile(c, existingUser.ID, nil, fullName, email, "", avatar)
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

	// Create session and save to database
	session, err := s.AuthService.CreateAndSaveSession(c, existingUser, otpReq)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo session: %v", err)
	}

	// Generate refresh token
	refreshToken := s.AuthService.GenRefreshToken(c, existingUser, session.SessionID)

	// Get user profile
	profile, err := s.AuthService.ProfileProvider.GetByProfileID(c, existingUser.UserID)
	if err != nil || profile == nil {
		// Nếu không tìm thấy profile, thử tạo lại
		profile, err = s.AuthService.ProfileProvider.CreateProfile(c, existingUser.ID, nil, fullName, email, "", avatar)
		if err != nil {
			return nil, fmt.Errorf("không tìm thấy thông tin người dùng và không thể tạo mới")
		}
	}

	// Ghi log admin history
	s.AuthService.LogLoginHistory(c, existingUser, profile, "Đăng nhập hệ thống bằng Facebook OAuth")

	response := s.AuthService.ResponseLogin(c, existingUser, profile, session.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}

// extractString helper function để extract string từ map
func (s *FacebookAuthService) extractString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

// extractAvatar helper function để extract avatar URL từ Facebook response
func (s *FacebookAuthService) extractAvatar(userInfo map[string]interface{}) string {
	if picture, ok := userInfo["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			if url, ok := data["url"].(string); ok {
				return url
			}
		}
	}
	return ""
}
