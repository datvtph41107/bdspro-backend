package oauth

import (
	_errors "common/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"user/config"
	"user/internal"
	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/repo"
	"user/internal/models"
	"user/internal/usecase"
	"user/internal/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OAuthGoogle struct {
	HttpClient     *http.Client
	AuthService    *usecase.AuthUsecase
	AuthMethodRepo repo.AuthMethodRepository
	// ProfileProvider *client.UserClient
	ProfileUsecase *usecases.ProfileUsecase
}

func NewGoogleAuthService(httpClient *http.Client,
	authService *usecase.AuthUsecase,
	authMethodRepo repo.AuthMethodRepository,
	profileUsecase *usecases.ProfileUsecase,
	// ProfileProvider *client.UserClient,
) *OAuthGoogle {
	return &OAuthGoogle{
		HttpClient:     httpClient,
		AuthService:    authService,
		AuthMethodRepo: authMethodRepo,
		ProfileUsecase: profileUsecase,
	}
}

func (g *OAuthGoogle) GetUrlLogin() string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + config.Properties.OAuth.Google.ClientID +
		"&redirect_uri=" + url.QueryEscape(config.Properties.OAuth.Google.RedirectURI) +
		"&response_type=code&scope=openid%20email%20profile&state=xyz"
}

func (g *OAuthGoogle) HandleCallback(c *gin.Context, code string, state string) (*dto.AuthLoginResponse, error) {
	accessToken, err := g.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	userInfo, err := g.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return g.RegisterOrLogin(c, userInfo, &dto.OtpVerifyRequest{})
}

func (g *OAuthGoogle) HandleCallbackWithDeviceInfo(c context.Context, code string, state string, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	accessToken, err := g.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	userInfo, err := g.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return g.RegisterOrLogin(c, userInfo, otpReq)
}

func (g *OAuthGoogle) getAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", config.Properties.OAuth.Google.ClientID)
	data.Set("client_secret", config.Properties.OAuth.Google.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.Properties.OAuth.Google.RedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := g.HttpClient.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	fmt.Println("result", result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("failed to get access token")
	}
	return token, nil
}

func (g *OAuthGoogle) getUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := g.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (g *OAuthGoogle) RegisterOrLogin(c context.Context, userInfo map[string]interface{}, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
	oauthId, ok := userInfo["sub"].(string)
	if !ok {
		return nil, errors.New("failed to get user ID from Google")
	}

	fullName := ""
	if name, ok := userInfo["name"].(string); ok {
		fullName = name
	}

	avatar := ""
	if picture, ok := userInfo["picture"].(string); ok {
		avatar = picture
	}

	email := ""
	if e, ok := userInfo["email"].(string); ok {
		email = e
	}

	existingUser, err := g.AuthMethodRepo.FindByOAuthID(c, oauthId, "GOOGLE")
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("lỗi khi tìm kiếm user: %v", err)
	}

	if existingUser == nil {
		oauthUser := &auth.AuthMethod{
			Provider: "GOOGLE",
			AuthName: oauthId,
			FullName: fullName,
			Avatar:   avatar,
			Email:    email,
		}
		createdUser, err := g.AuthMethodRepo.Create(c, oauthUser)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo user mới: %v", err)
		}
		_, err = g.AuthService.NewUser(c, createdUser.ID)
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
		if existingUser.Email != email && email != "" {
			existingUser.Email = email
			updated = true
		}

		if updated {
			_, _ = g.AuthMethodRepo.Update(c, existingUser)
		}
	}

	// Nếu chưa có UserID, tạo profile mới (không bắt buộc số điện thoại)
	if existingUser.UserID == 0 {
		profile, err := g.ProfileUsecase.ProfileRepo.CreateUserProfile(c, &models.UserProfileEntity{
			FullName: fullName,
			Email:    email,
			Avatar:   avatar,
		})
		if err != nil {
			return nil, fmt.Errorf("lỗi khi tạo profile: %v", err)
		}

		// Cập nhật UserID cho AuthMethod
		existingUser.UserID = profile.ProfileID
		_, err = g.AuthMethodRepo.Update(c, existingUser)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi cập nhật UserID: %v", err)
		}
	}

	// Nếu không có otpReq thì tạo mặc định
	if otpReq == nil {
		otpReq = &dto.OtpVerifyRequest{}
	}

	session, err := g.AuthService.CreateAndSaveSession(c, existingUser, otpReq)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo session: %v", err)
	}

	profile, err := g.ProfileUsecase.ProfileRepo.GetByProfileID(existingUser.UserID)
	if err != nil || profile == nil {
		// Nếu không tìm thấy profile, thử tạo lại (trường hợp user_id đã có nhưng profile bị xóa)
		profile, err = g.ProfileUsecase.ProfileRepo.CreateUserProfile(c, &models.UserProfileEntity{
			ProfileID:    existingUser.UserID,
			FullName:     fullName,
			Email:        email,
			Avatar:       avatar,
			ReferralCode: "",
		})

		// profile = profileResponse
		if err != nil {
			return nil, _errors.ReturnError(service.UserNotFoundAndCreateFailedLower, _errors.WithLegacyCode(int32(enums.NOT_FOUND_ACCOUNT)))
		}
	}

	refreshToken := g.AuthService.GenRefreshToken(c, existingUser, session.SessionID)

	// Ghi log admin history
	g.AuthService.LogLoginHistory(c, existingUser, &dto.ProfileDTO{
		ProfileID:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
	}, "Đăng nhập hệ thống bằng Google OAuth")

	response := g.AuthService.ResponseLogin(c, existingUser, &dto.ProfileDTO{
		ProfileID:    profile.ProfileID,
		FullName:     profile.FullName,
		Email:        profile.Email,
		Avatar:       profile.Avatar,
		ReferralCode: profile.ReferralCode,
	}, session.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}
