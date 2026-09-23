package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"user/config"
	"user/internal/domain/auth"
	"user/internal/domain/user"
	"user/internal/dto"
	"user/internal/interface/repo"
	"user/internal/usecase"
	usecases "user/internal/usecase/user"

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
	redirectURI := strings.TrimSpace(config.Properties.OAuth.Google.RedirectURI)
	log.Printf("[OAuthGoogle] GetUrlLogin redirect_uri=%s", redirectURI)
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + config.Properties.OAuth.Google.ClientID +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=openid%20email%20profile&state=xyz"
}

func (g *OAuthGoogle) HandleCallback(c *gin.Context, code string, state string) (*dto.AuthLoginResponse, error) {
	accessToken, err := g.getAccessToken(code, "")
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
	return g.HandleCallbackWithDeviceInfoAndRedirect(c, code, state, "", otpReq)
}

func (g *OAuthGoogle) HandleCallbackWithDeviceInfoAndRedirect(
	c context.Context,
	code string,
	state string,
	redirectURI string,
	otpReq *dto.OtpVerifyRequest,
) (*dto.AuthLoginResponse, error) {
	accessToken, err := g.getAccessToken(code, redirectURI)
	if err != nil {
		return nil, err
	}

	userInfo, err := g.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return g.RegisterOrLogin(c, userInfo, otpReq)
}

func (g *OAuthGoogle) getAccessToken(code string, redirectURI string) (string, error) {
	resolvedRedirect := strings.TrimSpace(redirectURI)
	if resolvedRedirect == "" {
		resolvedRedirect = strings.TrimSpace(config.Properties.OAuth.Google.RedirectURI)
	}
	data := url.Values{}
	data.Set("client_id", config.Properties.OAuth.Google.ClientID)
	data.Set("client_secret", config.Properties.OAuth.Google.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", resolvedRedirect)
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

	if token, ok := result["access_token"].(string); ok && token != "" {
		return token, nil
	}

	errCode, _ := result["error"].(string)
	errDesc, _ := result["error_description"].(string)
	if errCode == "" {
		errCode = "unknown_error"
	}
	if errDesc == "" {
		errDesc = "failed to get access token"
	}
	return "", fmt.Errorf("google oauth token exchange failed: %s (%s); redirect_uri=%s", errCode, errDesc, resolvedRedirect)
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
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("lỗi khi tìm kiếm user: %v", err)
		}
		existingUser = nil
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

	profile, err := g.ensureGoogleProfile(c, existingUser, fullName, email, avatar)
	if err != nil {
		return nil, err
	}

	if otpReq == nil {
		otpReq = &dto.OtpVerifyRequest{}
	}

	session, err := g.AuthService.CreateAndSaveSession(c, existingUser, otpReq)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo session: %v", err)
	}

	refreshToken := g.AuthService.GenRefreshToken(c, existingUser, session.SessionID)

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
	response.RefreshToken = refreshToken
	response.ReferralCode = profile.ReferralCode
	return response, nil
}

// ensureGoogleProfile returns an alive profile for the Google auth method.
// Missing profile (new user or orphaned user_id) creates a fresh row and relinks auth_method.
func (g *OAuthGoogle) ensureGoogleProfile(
	c context.Context,
	existingUser *auth.AuthMethod,
	fullName, email, avatar string,
) (*user.UserProfileEntity, error) {
	if existingUser == nil {
		return nil, errors.New("auth method is required")
	}
	if existingUser.UserID != 0 {
		profile, err := g.ProfileUsecase.ProfileRepo.GetByProfileID(existingUser.UserID)
		if err == nil && profile != nil {
			return profile, nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("lỗi khi đọc profile: %v", err)
		}
	}

	profile, err := g.ProfileUsecase.ProfileRepo.CreateUserProfile(c, &user.UserProfileEntity{
		FullName: fullName,
		Email:    email,
		Avatar:   avatar,
	})
	if err != nil {
		return nil, fmt.Errorf("lỗi khi đăng ký profile mới: %v", err)
	}
	existingUser.UserID = profile.ProfileID
	if _, err := g.AuthMethodRepo.Update(c, existingUser); err != nil {
		return nil, fmt.Errorf("lỗi khi gắn profile vào auth method: %v", err)
	}
	return profile, nil
}
