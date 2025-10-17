package handler

import (
	"auth/infra/oauth"
	"auth/internal/dto"
	_routes "common/routes"
	"context"
	"net/http"
	authpb "pb/types/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// GetLoginURL - gRPC method để lấy URL đăng nhập OAuth
func (h *OAuthHandler) GetLoginURL(ctx context.Context, req *authpb.GetLoginURLRequest) (*authpb.GetLoginURLResponse, error) {
	var url string

	switch req.Provider {
	case "zalo":
		url = h.ZaloAuthService.GetUrlLogin()
	case "google":
		url = h.GoogleAuthService.GetUrlLogin()
	case "facebook":
		url = h.FacebookAuthService.GetUrlLogin()
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Provider không hợp lệ: %s", req.Provider)
	}

	return &authpb.GetLoginURLResponse{
		Url: url,
	}, nil
}

// HandleCallback - gRPC method để xử lý OAuth callback
func (h *OAuthHandler) HandleCallback(ctx context.Context, req *authpb.HandleCallbackRequest) (*authpb.HandleCallbackResponse, error) {
	var result *dto.AuthLoginResponse
	var err error

	// Tạo gin.Context giả để truyền vào các OAuth service (cần refactor sau)
	// TODO: Refactor OAuth services để không phụ thuộc vào gin.Context
	ginCtx := &gin.Context{}
	ginCtx.Request = &http.Request{}

	switch req.Provider {
	case "zalo":
		result, err = h.ZaloAuthService.HandleCallback(ginCtx, req.Code, req.State)
	case "google":
		result, err = h.GoogleAuthService.HandleCallback(ginCtx, req.Code, req.State)
	case "facebook":
		result, err = h.FacebookAuthService.HandleCallback(ginCtx, req.Code, req.Platform, req.State)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Provider không hợp lệ: %s", req.Provider)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi xử lý callback: %v", err)
	}

	return &authpb.HandleCallbackResponse{
		AuthId:       result.AuthID,
		ProfileId:    result.ProfileID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		FullName:     result.FullName, // OAuth fullname
		// Avatar:       result.Avatar,
		// Email:        result.Email,
		// Phone:        result.Phone,
		Provider: req.Provider,
	}, nil
}
