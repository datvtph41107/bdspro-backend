package oauth

import (
	"auth/config"
	"auth/infra/client"
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/enums"
	"auth/internal/interface/repo"
	"auth/internal/usecase"
	_errors "common/errors"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type OAuthGoogle struct {
	HttpClient      *http.Client
	AuthService     *usecase.AuthUsecase
	AuthMethodRepo  repo.AuthMethodRepository
	ProfileProvider *client.UserClient
}

func NewGoogleAuthService(httpClient *http.Client,
	authService *usecase.AuthUsecase,
	authMethodRepo repo.AuthMethodRepository,
	ProfileProvider *client.UserClient,
) *OAuthGoogle {
	return &OAuthGoogle{
		HttpClient:      httpClient,
		AuthService:     authService,
		AuthMethodRepo:  authMethodRepo,
		ProfileProvider: ProfileProvider,
	}
}

func (g *OAuthGoogle) GetUrlLogin() string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + config.Properties.OAuth.Google.ClientID +
		"&redirect_uri=" + url.QueryEscape(config.Properties.OAuth.Google.RedirectURI) +
		"&response_type=code&scope=openid%20email%20profile&state=xyz"
}

func (g *OAuthGoogle) HandleCallback(c *gin.Context, accessToken string, state string) (*dto.AuthLoginResponse, error) {
	// accessToken, err := g.getAccessToken(code)
	// if err != nil {
	// 	return nil, err
	// }

	userInfo, err := g.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return g.RegisterOrLogin(c, userInfo, &dto.OtpVerifyRequest{})
}

func (g *OAuthGoogle) HandleCallbackWithDeviceInfo(c context.Context, accessToken string, state string, otpReq *dto.OtpVerifyRequest) (*dto.AuthLoginResponse, error) {
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
		return nil, errors.New("failed to get user ID")
	}
	existingUser, _ := g.AuthMethodRepo.FindByOAuthID(c, oauthId, "GOOGLE")
	// if err != nil {
	// 	return nil, err
	// }

	if existingUser == nil {
		oauthUser := &domain.AuthMethod{
			Provider: "GOOGLE",
			AuthName: oauthId,
			FullName: userInfo["name"].(string),
			Avatar:   userInfo["picture"].(string),
			Email:    userInfo["email"].(string),
		}
		createdUser, err := g.AuthMethodRepo.Create(c, oauthUser)
		if err != nil {
			return nil, err
		}
		g.AuthService.NewUser(c, createdUser.ID)
		existingUser = createdUser
	}

	if existingUser.UserID == 0 {
		return g.AuthService.ResponseLogin(existingUser, nil, 0), nil
	}
	// Nếu không có otpReq thì tạo mặc định
	if otpReq == nil {
		otpReq = &dto.OtpVerifyRequest{}
	}

	session, _ := g.AuthService.CreateAndSaveSession(c, existingUser, otpReq)

	// g.AuthService.CookieService.SetRefreshToken(c.Writer, refreshToken)
	// planID := g.AuthService.GetPlanID(existingUser.UserID)
	profile, err := g.ProfileProvider.GetByProfileID(c, existingUser.UserID)
	if err != nil || profile == nil {
		return nil, _errors.ReturnError(int32(enums.NOT_FOUND_ACCOUNT), "không tìm thấy thông tin người dùng")
	}

	refreshToken := g.AuthService.GenRefreshToken(c, existingUser, session.SessionID)

	// Ghi log admin history
	g.AuthService.LogLoginHistory(c, existingUser, profile, "Đăng nhập hệ thống bằng Google OAuth")

	response := g.AuthService.ResponseLogin(existingUser, profile, session.SessionID)
	response.RefreshToken = refreshToken         // Set refreshToken vào response body
	response.ReferralCode = profile.ReferralCode // Set mã giới thiệu vào response
	return response, nil
}
