package controller

import (
	"net/http"

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
	auth := prefix + "/auth"
	mux.HandleFunc("GET "+auth+"/phone/check/{phone}", c.httpPhoneCheck)
	mux.HandleFunc("POST "+auth+"/otp/request", c.httpRequestOTP)
	mux.HandleFunc("POST "+auth+"/otp/verify", c.httpVerifyOTP)
	mux.HandleFunc("POST "+auth+"/otp/resend", c.httpResendOTP)
	mux.HandleFunc("POST "+auth+"/password/login", c.httpLoginWithPassword)
	mux.HandleFunc("POST "+auth+"/login/key", c.httpLoginWithKey)
	mux.HandleFunc("POST "+auth+"/admin/login", c.httpAdminLogin)
	mux.HandleFunc("POST "+auth+"/token/refresh", c.httpRefreshToken)
}

func (c *LoginController) httpPhoneCheck(w http.ResponseWriter, r *http.Request) {
	result, err := c.PhoneCheckEncoded(r.Context(), r.PathValue("phone"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpRequestOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.OtpRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.login.RequestOTP(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpVerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.OtpVerifyRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		writeHTTPError(w, http.StatusBadRequest, "platform, version và os là bắt buộc")
		return
	}
	result, err := c.login.VerifyOTP(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpResendOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.OtpRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.login.ResendOTP(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpLoginWithPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginWithPasswordRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		writeHTTPError(w, http.StatusBadRequest, "platform, version và os là bắt buộc")
		return
	}
	result, err := c.login.LoginWithPassword(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpLoginWithKey(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginWithKeyRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.login.LoginWithKey(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpAdminLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.AdminLoginRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Platform == "" || req.Version == "" || req.OS == "" {
		writeHTTPError(w, http.StatusBadRequest, "platform, version và os là bắt buộc")
		return
	}
	result, err := c.login.AdminLogin(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func (c *LoginController) httpRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenHTTPRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.login.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeHTTPOK(w, result)
}

func writeHTTPOK(w http.ResponseWriter, data any) {
	httpapi.WriteOK(w, data)
}

func writeHTTPError(w http.ResponseWriter, httpStatus int, message string) {
	httpapi.WriteError(w, httpStatus, message)
}

func writeServiceError(w http.ResponseWriter, err error) {
	httpapi.WriteServiceError(w, err)
}
