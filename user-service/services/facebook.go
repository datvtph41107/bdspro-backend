package services

import (
	"auth/config"
	"auth/data"
	"auth/models"
	_db "common/db"
	_routes "common/routes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type FacebookAuthService struct {
	AuthService *AuthService
	HttpClient  *http.Client
}

func NewFacebookAuthService(authService *AuthService, httpClient *http.Client) *FacebookAuthService {
	return &FacebookAuthService{
		// apiClient:   apiClient,
		// properties:  properties,
		// userRepo:    userRepo,
		AuthService: authService,
		HttpClient:  httpClient,
	}
}

func (s *FacebookAuthService) GetUrlLogin() string {
	fbAuthUrl := "https://www.facebook.com/v22.0/dialog/oauth?" +
		"client_id=" + config.AppProperties.Facebook.Client.ID +
		"&redirect_uri=" + url.QueryEscape(config.AppProperties.Facebook.RedirectURI) +
		"&state=xyz" +
		"&scope=email"
	return fbAuthUrl
}

func (s *FacebookAuthService) HandleCallback(c *gin.Context, code string, state string) (*data.AuthLoginResponse, error) {
	params := url.Values{}
	params.Set("client_id", config.AppProperties.Facebook.Client.ID)
	params.Set("redirect_uri", config.AppProperties.Facebook.RedirectURI)
	params.Set("client_secret", config.AppProperties.Facebook.Client.Secret)
	params.Set("code", code)

	tokenURL := "https://graph.facebook.com/v22.0/oauth/access_token?" + params.Encode()
	resp, err := s.HttpClient.Get(tokenURL)
	if err != nil {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Không lấy được access_token từ Facebook",
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Lỗi đọc phản hồi từ Facebook",
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Lỗi parse JSON từ Facebook",
		}
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Không lấy được access_token từ Facebook",
		}
	}

	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return s.RegisterOrLogin(c, userInfo)
}

func (s *FacebookAuthService) fetchUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email&access_token="+accessToken, nil)
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

func (s *FacebookAuthService) RegisterOrLogin(c *gin.Context, userInfo map[string]interface{}) (*data.AuthLoginResponse, error) {
	username, ok := userInfo["id"].(string)
	if !ok {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Không lấy được id từ Facebook",
		}
	}
	e, err := models.AuthFindByOAuthId(username, "FACEBOOK")
	if err != nil {
		e = &models.UserAuthEntity{
			Provider: "FACEBOOK",
			Username: username,
			OAuthID:  username,
			Name:     userInfo["name"].(string),
			Email:    extractEmail(userInfo),
			Avatar:   ExtractAvatar(userInfo),
		}
		_db.SaveWithAudit(c, e)
		s.AuthService.NewUser(c, e.ID)
	}

	session := s.AuthService.CreateSession(c, e, &data.OtpVerifyRequest{})

	return s.AuthService.ResponseLogin(e, session.SessionID), nil
}

func extractEmail(userInfo map[string]interface{}) string {
	if email, ok := userInfo["email"].(string); ok {
		return email
	}
	return ""
}

func ExtractAvatar(userInfo map[string]interface{}) string {
	if picture, ok := userInfo["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			if url, ok := data["url"].(string); ok {
				return url
			}
		}
	}
	return ""
}
