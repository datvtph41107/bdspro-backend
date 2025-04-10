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

type ZaloAuthService struct {
	// apiClient   *ApiClient
	// properties  *AuthProperties
	// userRepo    *UserAuthRepository
	AuthService *AuthService
	HttpClient  *http.Client
}

func NewZaloAuthService(httpClient *http.Client, authService *AuthService) *ZaloAuthService {
	return &ZaloAuthService{
		// apiClient:   apiClient,
		// properties:  properties,
		// userRepo:    userRepo,
		HttpClient:  httpClient,
		AuthService: authService,
	}
}

func (s *ZaloAuthService) GetUrlLogin() string {
	zaloAuthUrl := config.AppProperties.Zalo.PermissionURL +
		"app_id=" + config.AppProperties.Zalo.App.ID +
		"&redirect_uri=" + url.QueryEscape(config.AppProperties.Zalo.RedirectURI) +
		"&state=xyz"
	return zaloAuthUrl
}

// func (s *ZaloAuthService) HandleCallback(code string, state string, request *http.Request) (map[string]interface{}, error) {
// 	params := url.Values{}
// 	params.Set("app_id", config.AppProperties.Zalo.App.ID)
// 	params.Set("code", code)
// 	params.Set("grant_type", "authorization_code")

// 	headers := map[string]string{
// 		"secret_key": config.AppProperties.Zalo.SecretKey,
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

func (s *ZaloAuthService) HandleCallback(c *gin.Context, code string, state string) (*data.AuthLoginResponse, error) {
	// params := url.Values{}
	// params.Set("app_id", config.AppProperties.Zalo.App.ID)
	// params.Set("code", code)
	// params.Set("grant_type", "authorization_code")

	// req, err := http.NewRequest("POST", "https://oauth.zaloapp.com/v4/access_token", bytes.NewBufferString(params.Encode()))
	// if err != nil {
	// 	return nil, err
	// }
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("secret_key", config.AppProperties.Zalo.SecretKey)

	// resp, err := s.HttpClient.Do(req)
	// if err != nil {
	// 	return nil, &_routes.Except{
	// 		Code:    401,
	// 		Message: "Không lấy được access_token từ Zalo",
	// 	}
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// var result map[string]interface{}
	// if err := json.Unmarshal(body, &result); err != nil {
	// 	return nil, err
	// }

	// accessToken, ok := result["access_token"].(string)
	// if !ok {
	// 	return nil, &_routes.Except{
	// 		Code:    401,
	// 		Message: "Token hết hạn",
	// 	}
	// }

	userInfo, err := s.fetchUserInfo(code)
	if err != nil {
		return nil, err
	}

	return s.RegisterOrLogin(c, userInfo)
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

func (s *ZaloAuthService) RegisterOrLogin(c *gin.Context, userInfo map[string]interface{}) (*data.AuthLoginResponse, error) {
	username, ok := userInfo["id"].(string)
	if !ok {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Không lấy được id từ Zalo",
		}
	}
	e, err := models.AuthFindByOAuthId(username, "ZALO")
	if err != nil {
		e = &models.UserAuthEntity{
			Provider:    "ZALO",
			Username:    username,
			OAuthID:     username,
			Name:        userInfo["name"].(string),
			Avatar:      extractAvatar(userInfo),
			IsSensitive: userInfo["is_sensitive"].(bool),
		}
		_db.SaveWithAudit(c, e)
		s.AuthService.NewUser(c, e.ID)
	}

	session := s.AuthService.CreateSession(c, e, &data.OtpVerifyRequest{})
	if err != nil {
		return nil, err
	}

	return s.AuthService.ResponseLogin(e, session.SessionID), nil
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
