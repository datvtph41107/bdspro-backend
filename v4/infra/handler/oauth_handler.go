package handler

import (
	"common/httpapi"
	"context"
	"log"
	"mdn4/modules/auth/dto"
	"mdn4/modules/auth/oauth"
	usecase "mdn4/modules/auth/services"
	"net/http"
	authpb "pb/types/auth"
	"strings"

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

// RegisterHTTP mounts the OAuthService HTTP contract (shared/protobuf/schema/auth/oauth.proto).
//
//	GET /v2/auth/oauth/{provider}/login
//	GET /v2/auth/oauth/{provider}/callback
//	GET /v2/auth/oauth/connected-accounts
func (h *OAuthHandler) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/auth/oauth/connected-accounts", h.connectedAccountsHTTP)
	httpapi.GET(mux, "/v2/auth/oauth/{provider}/login", h.loginURLHTTP)
	httpapi.GET(mux, "/v2/auth/oauth/{provider}/callback", h.callbackHTTP)
}

type oauthLoginHTTPRequest struct {
	Provider string `path:"provider"`
}

type oauthCallbackHTTPRequest struct {
	Provider    string `path:"provider"`
	Code        string `form:"code"`
	State       string `form:"state"`
	Platform    string `form:"platform"`
	Version     string `form:"version"`
	OS          string `form:"os"`
	DeviceName  string `form:"deviceName"`
	RedirectURI string `form:"redirectUri"`
}

type oauthCallbackHTTPResponse struct {
	AuthID       uint64 `json:"authId"`
	ProfileID    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName,omitempty"`
	Provider     string `json:"provider"`
	ReferralCode string `json:"referralCode,omitempty"`
}

func (h *OAuthHandler) loginURLHTTP(ctx context.Context, req oauthLoginHTTPRequest) (*httpapi.Result, error) {
	resp, err := h.GetLoginURL(ctx, &authpb.GetLoginURLRequest{Provider: req.Provider})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(dto.GetLoginURLResponse{URL: resp.GetUrl(), State: resp.GetState()}), nil
}

func (h *OAuthHandler) callbackHTTP(ctx context.Context, req oauthCallbackHTTPRequest) (*httpapi.Result, error) {
	if strings.TrimSpace(req.Code) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "code là bắt buộc")
	}
	resp, err := h.HandleCallback(ctx, &authpb.HandleCallbackRequest{
		Provider:    req.Provider,
		Code:        req.Code,
		State:       req.State,
		Platform:    req.Platform,
		Version:     req.Version,
		Os:          req.OS,
		DeviceName:  req.DeviceName,
		RedirectUri: req.RedirectURI,
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(oauthCallbackHTTPResponse{
		AuthID:       resp.GetAuthId(),
		ProfileID:    resp.GetProfileId(),
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		FullName:     resp.GetFullName(),
		Provider:     resp.GetProvider(),
		ReferralCode: resp.GetReferralCode(),
	}), nil
}

func (h *OAuthHandler) connectedAccountsHTTP(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	accounts, err := h.OAuthUsecase.GetConnectedAccounts(ctx)
	if err != nil {
		return nil, err
	}
	if accounts == nil {
		accounts = []*dto.ConnectedAccountDTO{}
	}
	return httpapi.OK(accounts), nil
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
