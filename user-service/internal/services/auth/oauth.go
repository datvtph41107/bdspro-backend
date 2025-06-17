package auth_service

import (
	"net/http"
	"user/internal/dto"

	"github.com/gin-gonic/gin"
)

type OAuthController interface {
	HandleCallback(c *gin.Context, code string, state string) (*dto.AuthLoginResponse, error)
}

type OAuthService struct {
	GoogleAuthService   *GoogleAuthService
	ZaloAuthService     *ZaloAuthService
	FacebookAuthService *FacebookAuthService
}

func NewOAuthService(
	google *GoogleAuthService,
	zalo *ZaloAuthService,
	facebook *FacebookAuthService,
) *OAuthService {
	return &OAuthService{
		GoogleAuthService:   google,
		ZaloAuthService:     zalo,
		FacebookAuthService: facebook,
	}
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}

// type AuthLoginResponse struct {
// 	ID string `json:"id"`
// }
