package auth_service

import (
	_db "common/db"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"user/config/property"
	"user/internal/dto"
	auth_model "user/internal/models/auth"

	"github.com/gin-gonic/gin"
)

type GoogleAuthService struct {
	// Properties   *UserProperties
	HttpClient  *http.Client
	AuthService *AuthService
	// UserAuthRepo *UserAuthRepository
}

func NewGoogleAuthService(httpClient *http.Client, authService *AuthService) *GoogleAuthService {
	return &GoogleAuthService{
		// Properties:  properties,
		HttpClient:  httpClient,
		AuthService: authService,
	}
}

func (g *GoogleAuthService) GetUrlLogin() string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + property.AppProperties.Google.Client.ID +
		"&redirect_uri=" + url.QueryEscape(property.AppProperties.Google.RedirectURI) +
		"&response_type=code&scope=openid%20email%20profile&state=xyz"
}

func (g *GoogleAuthService) HandleCallback(c *gin.Context, accessToken string, state string) (*dto.AuthLoginResponse, error) {
	// accessToken, err := g.getAccessToken(code)
	// if err != nil {
	// 	return nil, err
	// }

	userInfo, err := g.getUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return g.RegisterOrLogin(c, userInfo)
}

func (g *GoogleAuthService) getAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", property.AppProperties.Google.Client.ID)
	data.Set("client_secret", property.AppProperties.Google.Client.Secret)
	data.Set("code", code)
	data.Set("redirect_uri", property.AppProperties.Google.RedirectURI)
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

func (g *GoogleAuthService) getUserInfo(accessToken string) (map[string]interface{}, error) {
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

func (g *GoogleAuthService) RegisterOrLogin(c *gin.Context, userInfo map[string]interface{}) (*dto.AuthLoginResponse, error) {
	oauthId, ok := userInfo["sub"].(string)
	if !ok {
		return nil, errors.New("failed to get user ID")
	}
	existingUser, _ := auth_model.AuthFindByOAuthId(oauthId, "GOOGLE")
	// if err != nil {
	// 	return nil, err
	// }

	if existingUser == nil {
		oauthUser := &auth_model.UserAuthEntity{
			Provider: "GOOGLE",
			Username: oauthId,
			OAuthID:  oauthId,
			Name:     userInfo["name"].(string),
			Avatar:   userInfo["picture"].(string),
			Email:    userInfo["email"].(string),
		}
		if err := _db.SaveWithAudit(c, oauthUser); err != nil {
			return nil, err
		}
		g.AuthService.NewUser(c, oauthUser.ID)
		existingUser = oauthUser
	}

	session := g.AuthService.CreateSession(c, existingUser, &dto.OtpVerifyRequest{})
	// g.AuthService.CookieService.SetRefreshToken(c.Writer, refreshToken)
	// planID := g.AuthService.GetPlanID(existingUser.UserID)
	profile := g.AuthService.GetProfileByID(existingUser.UserID)
	g.AuthService.GenRefreshToken(c, existingUser, session.SessionID)
	return g.AuthService.ResponseLogin(existingUser, profile, session.SessionID), nil
}
