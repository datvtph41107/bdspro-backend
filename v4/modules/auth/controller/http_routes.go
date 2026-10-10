package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/auth/dto"
)

// APIVersionPrefixes mirrors gateway HTTP version mounts (/v2, /v2.1, /v3).
var APIVersionPrefixes = []string{"/v2", "/v2.1", "/v3"}

// RegisterHTTP mounts login routes under every API version prefix, same as gateway.
//
//	GET  {prefix}/auth/phone/check/{phone}
//	POST {prefix}/auth/otp/request
//	POST {prefix}/auth/otp/verify
//	POST {prefix}/auth/otp/resend
//	POST {prefix}/auth/password/login
//	POST {prefix}/auth/login/key
//	POST {prefix}/auth/admin/login
//	POST {prefix}/auth/token/refresh
func (c *LoginController) RegisterHTTP(mux *http.ServeMux) {
	for _, prefix := range APIVersionPrefixes {
		c.registerAuthHTTP(mux, prefix)
	}
}

func (c *LoginController) registerAuthHTTP(mux *http.ServeMux, prefix string) {
	auth := prefix
	httpapi.GET(mux, auth+"/auth/phone/check/{phone}", c.phoneCheck)
	httpapi.POST(mux, auth+"/auth/otp/request", c.requestOTP)
	httpapi.POST(mux, auth+"/auth/otp/verify", c.verifyOTP)
	httpapi.POST(mux, auth+"/auth/otp/resend", c.resendOTP)
	httpapi.POST(mux, auth+"/auth/password/login", c.loginWithPassword)
	httpapi.POST(mux, auth+"/auth/login/key", c.loginWithKey)
	httpapi.POST(mux, auth+"/auth/admin/login", c.adminLogin)
	httpapi.POST(mux, auth+"/auth/token/refresh", c.refreshToken)
}

func (c *LoginController) phoneCheck(ctx context.Context, req dto.PhoneCheckPath) (*httpapi.Result, error) {
	result, err := c.PhoneCheckEncoded(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) requestOTP(ctx context.Context, req dto.OtpRequest) (*httpapi.Result, error) {
	result, err := c.login.RequestOTP(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) verifyOTP(ctx context.Context, req dto.OtpVerifyRequest) (*httpapi.Result, error) {
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		return nil, _errors.ReturnError(400, "platform, version và os là bắt buộc")
	}
	result, err := c.login.VerifyOTP(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) resendOTP(ctx context.Context, req dto.OtpRequest) (*httpapi.Result, error) {
	result, err := c.login.ResendOTP(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) loginWithPassword(ctx context.Context, req dto.LoginWithPasswordRequest) (*httpapi.Result, error) {
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		return nil, _errors.ReturnError(400, "platform, version và os là bắt buộc")
	}
	result, err := c.login.LoginWithPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) loginWithKey(ctx context.Context, req dto.LoginWithKeyRequest) (*httpapi.Result, error) {
	result, err := c.login.LoginWithKey(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) adminLogin(ctx context.Context, req dto.AdminLoginRequest) (*httpapi.Result, error) {
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		return nil, _errors.ReturnError(400, "platform, version và os là bắt buộc")
	}
	result, err := c.login.AdminLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *LoginController) refreshToken(ctx context.Context, req dto.RefreshTokenHTTPRequest) (*httpapi.Result, error) {
	result, err := c.login.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}
