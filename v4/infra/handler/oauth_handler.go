package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	authpb "pb/types/auth"
	"userv4/infra/oauth"
	"userv4/internal/dto"
	"userv4/internal/httpapi"
	"userv4/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuthHandler struct {
	authpb.UnimplementedOAuthServiceServer
	GoogleAuthService   *oauth.OAuthGoogle
	ZaloAuthService     *oauth.ZaloAuthService
	FacebookAuthService *oauth.FacebookAuthService
	OAuthUsecase        *usecase.OAuthUsecase
}

func NewOAuthRouter(
	google *oauth.OAuthGoogle,
	zalo *oauth.ZaloAuthService,
	facebook *oauth.FacebookAuthService,
	oauthUsecase *usecase.OAuthUsecase,
) *OAuthHandler {
	return &OAuthHandler{
		GoogleAuthService:   google,
		ZaloAuthService:     zalo,
		FacebookAuthService: facebook,
		OAuthUsecase:        oauthUsecase,
	}
}

func OAuthRoutes(mux *http.ServeMux, router *OAuthHandler) {
	mux.HandleFunc("GET /user/oauth/{provider}/login", router.RedirectToLogin)
	mux.HandleFunc("GET /user/oauth/{provider}/callback", router.OauthCallback)
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
func (route *OAuthHandler) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	var url string

	switch provider {
	case "zalo":
		url = route.ZaloAuthService.GetUrlLogin()
	case "google":
		url = route.GoogleAuthService.GetUrlLogin()
	case "facebook":
		url = route.FacebookAuthService.GetUrlLogin()
	default:
		writeOAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "Vui lòng chọn provider hợp lệ"})
		return
	}

	writeOAuthJSON(w, http.StatusOK, map[string]string{"url": url})
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
func (route *OAuthHandler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	q := r.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	platform := q.Get("platform")
	ctx := r.Context()

	var result *dto.AuthLoginResponse
	var err error

	switch provider {
	case "zalo":
		result, err = route.ZaloAuthService.HandleCallback(ctx, code, state)
	case "google":
		result, err = route.GoogleAuthService.HandleCallback(ctx, code, state)
	case "facebook":
		result, err = route.FacebookAuthService.HandleCallback(ctx, code, platform, state)
	default:
		writeOAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "Provider không hợp lệ"})
		return
	}

	if err != nil {
		httpapi.WriteServiceError(w, err)
		return
	}
	httpapi.WriteOK(w, result)
}

func writeOAuthJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
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
	// Validate required fields for session creation
	if req.Platform == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Platform là bắt buộc (WEB, MOBILE, DESKTOP)")
	}
	if req.Version == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Version là bắt buộc")
	}
	if req.Os == "" {
		return nil, status.Errorf(codes.InvalidArgument, "OS là bắt buộc (iOS, Android, Windows, MacOS, Linux)")
	}

	var result *dto.AuthLoginResponse
	var err error

	otpReq := &dto.OtpVerifyRequest{
		Version:    req.Version,
		Platform:   req.Platform,
		OS:         req.Os,
		DeviceName: req.DeviceName,
	}

	switch req.Provider {
	case "zalo":
		result, err = h.ZaloAuthService.HandleCallbackWithDeviceInfo(ctx, req.Code, req.State, otpReq)
	case "google":
		result, err = h.GoogleAuthService.HandleCallbackWithDeviceInfoAndRedirect(ctx, req.Code, req.State, req.GetRedirectUri(), otpReq)
	case "facebook":
		result, err = h.FacebookAuthService.HandleCallbackWithDeviceInfo(ctx, req.Code, req.Platform, req.State, otpReq)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Provider không hợp lệ: %s", req.Provider)
	}

	if err != nil {
		log.Printf("[OAuth.HandleCallback] provider=%s failed: %v", req.Provider, err)
		msg := err.Error()
		switch {
		case strings.Contains(msg, "redirect_uri_mismatch"):
			return nil, status.Errorf(codes.InvalidArgument, "%s", msg)
		case strings.Contains(msg, "google oauth token exchange failed"):
			return nil, status.Errorf(codes.InvalidArgument, "%s", msg)
		case strings.Contains(msg, "Platform là bắt buộc"), strings.Contains(msg, "Version là bắt buộc"):
			return nil, status.Errorf(codes.InvalidArgument, "%s", msg)
		default:
			return nil, status.Errorf(codes.Internal, "Lỗi xử lý callback: %v", err)
		}
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
		Provider:     req.Provider,
		ReferralCode: result.ReferralCode, // Mã giới thiệu
	}, nil
}

// GetConnectedAccounts - gRPC method để lấy danh sách tài khoản MXH đã liên kết
// @Summary Lấy danh sách tài khoản MXH đã liên kết
// @Description Lấy danh sách các tài khoản mạng xã hội (Google, Facebook, Zalo) mà user đã liên kết
// @Tags OAuth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authpb.GetConnectedAccountsResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v2/auth/oauth/connected-accounts [get]
func (h *OAuthHandler) GetConnectedAccounts(ctx context.Context, req *authpb.GetConnectedAccountsRequest) (*authpb.GetConnectedAccountsResponse, error) {
	accounts, err := h.OAuthUsecase.GetConnectedAccounts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi khi lấy danh sách tài khoản: %v", err)
	}

	var pbAccounts []*authpb.ConnectedAccount
	for _, account := range accounts {
		pbAccounts = append(pbAccounts, &authpb.ConnectedAccount{
			AuthId:    account.AuthID,
			Provider:  account.Provider,
			AuthName:  account.AuthName,
			FullName:  account.FullName,
			Email:     account.Email,
			Avatar:    account.Avatar,
			CreatedAt: account.CreatedAt,
		})
	}

	return &authpb.GetConnectedAccountsResponse{
		Data: pbAccounts,
	}, nil
}
