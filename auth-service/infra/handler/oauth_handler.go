package handler

import (
	"auth/infra/oauth"
	"auth/internal/dto"
	_routes "common/routes"
	"net/http"
	authpb "pb/types/auth"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	authpb.UnimplementedOAuthServiceServer
	GoogleAuthService   *oauth.OAuthGoogle
	ZaloAuthService     *oauth.ZaloAuthService
	FacebookAuthService *oauth.FacebookAuthService
}

func NewOAuthRouter(
	google *oauth.OAuthGoogle,
	zalo *oauth.ZaloAuthService,
	facebook *oauth.FacebookAuthService,
) *OAuthHandler {
	return &OAuthHandler{
		GoogleAuthService:   google,
		ZaloAuthService:     zalo,
		FacebookAuthService: facebook,
	}
}

func OAuthRoutes(r *gin.Engine, router *OAuthHandler) {
	oauthGroup := r.Group("/user/oauth")
	{
		oauthGroup.GET("/:provider/login", router.RedirectToLogin)
		oauthGroup.GET("/:provider/callback", router.OauthCallback)
	}
}

// OAuthRoutes thiết lập các route liên quan đến OAuth.
// @Summary OAuth Login
// @Description Chuyển hướng người dùng đến trang đăng nhập của provider (Google, Facebook, Zalo, v.v.).
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "Tên của provider (google, facebook, zalo, ...)"
// @Success 302 "Chuyển hướng đến trang đăng nhập của provider"
// @Failure 400 {object} map[string]interface{} "Lỗi khi provider không hợp lệ"
// @Router /oauth/{provider}/login [get]
func (route *OAuthHandler) RedirectToLogin(c *gin.Context) {
	provider := c.Param("provider")
	var url string

	switch provider {
	case "zalo":
		url = route.ZaloAuthService.GetUrlLogin()
	case "google":
		url = route.GoogleAuthService.GetUrlLogin()
	case "facebook":
		url = route.FacebookAuthService.GetUrlLogin()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn provider hợp lệ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// @Summary OAuth Callback
// @Description Nhận callback từ provider sau khi người dùng đăng nhập và xử lý xác thực.
// @Tags OAuth
// @Accept json
// @Produce json
// @Param provider path string true "Tên của provider (google, facebook, zalo, ...)"
// @Param code query string true "Mã xác thực từ provider"
// @Param state query string false "Chuỗi state để bảo vệ chống tấn công CSRF"
// @Success 200 {object} map[string]interface{} "Thông tin người dùng từ provider"
// @Failure 400 {object} map[string]interface{} "Lỗi khi xử lý callback"
// @Router /oauth/{provider}/callback [get]
func (route *OAuthHandler) OauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	platform := c.Query("platform")
	var result *dto.AuthLoginResponse
	var err error

	switch provider {
	case "zalo":
		result, err = route.ZaloAuthService.HandleCallback(c, code, state)
	case "google":
		result, err = route.GoogleAuthService.HandleCallback(c, code, state)
	case "facebook":
		result, err = route.FacebookAuthService.HandleCallback(c, code, platform, state)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider không hợp lệ"})
		return
	}

	_routes.RouteResult(c, result, err)
	// c.JSON(http.StatusOK, response)
}
