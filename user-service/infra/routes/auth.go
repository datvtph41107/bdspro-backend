package routes

import (
	_routes "common/routes"
	"net/http"
	"strconv"

	"user/config/property"
	"user/internal/dto"
	auth_service "user/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	Service *auth_service.AuthService
}

func NewAuthRoutes(service *auth_service.AuthService) *AuthRouter {
	return &AuthRouter{
		Service: service,
	}
}

func AuthRoutes(r *gin.Engine, router *AuthRouter) {
	api := r.Group("/user/auth")

	api.POST("/otp/request", router.RequestOTP)

	api.POST("/otp/verify", router.VerifyOTP)

	api.POST("/otp/resend", router.ResendOTP)

	api.POST("/token/refresh", router.RefreshToken)

	api.DELETE("/logout", router.Logout)

	api.DELETE("/delete", router.Delete)

	api.PUT("/restore", router.Restore)

	api.POST("/qr/init", router.QRInit)

	api.POST("/qr/confirm", router.QRConfirm)

	api.POST("/qr/status", router.QRStatus)

	api.POST("/switch/organization", router.SwitchOrganization)
}

// @Tags Auth
// @Summary Gửi OTP
// @Description Gửi lại OTP hoặc submit lúc nhập SĐT
// @Accept  json
// @Produce  json
// @Param body body dto.OtpRequest true "OTP Request"
// @Router /v1/user/auth/otp/request [post]
func (router *AuthRouter) RequestOTP(c *gin.Context) {
	var dto dto.OtpRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := router.Service.RequestOtp(c, dto)
	// c.JSON(http.StatusOK, response)
	_routes.RouteResult(c, result, err)
}

// @Tags Auth
// @Summary Xác minh OTP
// @Description Xác minh OTP đã nhận
// @Accept  json
// @Produce  json
// @Param body body dto.OtpVerifyRequest true "OTP Verify Request"
// @Router /v1/user/auth/otp/verify [post]
func (router *AuthRouter) VerifyOTP(c *gin.Context) {
	var dto dto.OtpVerifyRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := router.Service.VerifyOtp(c, dto)
	// cookieService.SetRefreshToken(c.Writer, result.RefreshToken)
	// c.JSON(http.StatusOK, result)
	_routes.RouteResult(c, result, err)
}

// @Tags Auth
// @Summary Resend OTP
// @Description Gửi lại OTP mới
// @Accept  json
// @Produce  json
// @Param body body dto.OtpRequest true "OTP Request"
// @Router /v1/user/auth/otp/resend [post]
func (router *AuthRouter) ResendOTP(c *gin.Context) {
	var dto dto.OtpRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := router.Service.ResendOTP(c, dto)
	_routes.RouteResult(c, result, err)
}

// @Tags Auth
// @Summary Refresh token
// @Description Làm mới access token bằng refresh token trong cookie
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Router /v1/user/auth/token/refresh [post]
func (router *AuthRouter) RefreshToken(c *gin.Context) {
	refreshToken, _ := c.Cookie("rt")
	result, err := router.Service.RefreshToken(c, refreshToken)
	_routes.RouteResult(c, result, err)
}

// @Tags Auth
// @Summary Logout
// @Description Đăng xuất, xóa refresh token
// @Router /v1/user/auth/logout [delete]
func (router *AuthRouter) Logout(c *gin.Context) {
	// cookieService.DeleteRefreshToken(c.Writer)
	// c.JSON(http.StatusOK, response)
	result, err := router.Service.Logout(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Khôi phục tài khoản
// @Description Dùng để khôi phục số điện thoại khi mất hoặc đổi số mới
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body dto.RestorePhoneRequest true "Restore Phone Request"
// @Tags Auth
// @Router /v1/user/auth/restore [put]
func (router *AuthRouter) Restore(c *gin.Context) {
	var dto dto.RestorePhoneRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	result, err := router.Service.Restore(c, dto)
	_routes.RouteResult(c, result, err)
}

// QRInit tạo một QR session mới
// @Summary Tạo QR Login
// @Description API này tạo một phiên QR để người dùng có thể đăng nhập qua app
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SessionQRRequest true "Thông tin tạo QR"
// @Success 200 {object} dto.SessionQRResponse
// @Router /v1/user/auth/qr/init [post]
func (router *AuthRouter) QRInit(c *gin.Context) {
	var dto dto.SessionQRRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sessionStr, _ := c.Cookie(property.C_SESSION_ID)
	sessionId, _ := strconv.ParseUint(sessionStr, 10, 32)

	result, err := router.Service.CreateQRSession(c, &dto, sessionId)

	_routes.RouteResult(c, result, err)
}

// QRConfirm xác nhận QR login
// @Summary Xác nhận QR Login
// @Description API này xác nhận rằng app đã quét QR và người dùng đã đồng ý đăng nhập
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SessionQRConfirm true "Thông tin xác nhận QR"
// @Router /v1/user/auth/qr/confirm [post]
func (router *AuthRouter) QRConfirm(c *gin.Context) {
	var dto dto.SessionQRConfirm
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := router.Service.ConfirmQRSession(c, &dto)
	_routes.RouteResult(c, result, err)
}

// QRStatus kiểm tra trạng thái QR login
// @Summary Kiểm tra trạng thái QR
// @Description API này kiểm tra xem QR đã được xác nhận hay chưa
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SessionQRRequest true "Thông tin kiểm tra QR"
// @Router /v1/user/auth/qr/status [post]
func (router *AuthRouter) QRStatus(c *gin.Context) {
	var dto dto.SessionQRRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sessionStr, _ := c.Cookie(property.C_SESSION_ID)
	sessionId, _ := strconv.ParseUint(sessionStr, 10, 32)

	result, err := router.Service.VerifyQRSession(c, &dto, sessionId)
	_routes.RouteResult(c, result, err)
}

// Xóa tài khoản
// @Summary Xóa tài khoản
// @Description API này cho phép xóa người dùng
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PublicUserInfo "Cập nhật thành công"
// @Router /v1/user/auth/delete [delete]
func (p *AuthRouter) Delete(c *gin.Context) {
	result, err := p.Service.Delete(c)

	_routes.RouteResult(c, result, err)
}

// @Summary Chuyển đổi tổ chức
// @Description API này cho phép người dùng chuyển đổi tổ chức
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SwitchOrganizationRequest true "Thông tin chuyển đổi tổ chức"
// @Router /v1/user/auth/switch/organization [post]
func (router *AuthRouter) SwitchOrganization(c *gin.Context) {
	var dto dto.SwitchOrganizationRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := router.Service.SwitchOrganization(c, dto.OrganizationID)
	_routes.RouteResult(c, result, err)
}
