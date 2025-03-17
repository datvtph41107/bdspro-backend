package routes

import (
	"auth/config"
	"auth/data"
	"auth/services"
	_jwt "common/jwt"
	_routes "common/routes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	Service *services.AuthService
}

func NewAuthRoutes(service *services.AuthService) *AuthRouter {
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

	api.PUT("/restore", router.Restore)

	api.POST("/qr/init", router.QRInit)

	api.POST("/qr/confirm", _jwt.JWTAuthMiddleware(publicRoutes), router.QRConfirm)

	api.POST("/qr/status", router.QRStatus)
}

// @Tags Auth
// @Summary Gửi OTP
// @Description Gửi lại OTP hoặc submit lúc nhập SĐT
// @Accept  json
// @Produce  json
// @Param body body data.OtpRequest true "OTP Request"
// @Router /auth/otp/request [post]
func (router *AuthRouter) RequestOTP(c *gin.Context) {
	var dto data.OtpRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, _ := router.Service.RequestOtp(c, dto)
	c.JSON(http.StatusOK, response)
}

// @Tags Auth
// @Summary Xác minh OTP
// @Description Xác minh OTP đã nhận
// @Accept  json
// @Produce  json
// @Param body body data.OtpVerifyRequest true "OTP Verify Request"
// @Router /auth/otp/verify [post]
func (router *AuthRouter) VerifyOTP(c *gin.Context) {
	var dto data.OtpVerifyRequest
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
// @Param body body data.OtpRequest true "OTP Request"
// @Router /auth/otp/resend [post]
func (router *AuthRouter) ResendOTP(c *gin.Context) {
	var dto data.OtpRequest
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
// @Router /auth/token/refresh [post]
func (router *AuthRouter) RefreshToken(c *gin.Context) {
	refreshToken, _ := c.Cookie("rt")
	result, err := router.Service.RefreshToken(refreshToken)
	_routes.RouteResult(c, result, err)
}

// @Tags Auth
// @Summary Logout
// @Description Đăng xuất, xóa refresh token
// @Router /auth/logout [delete]
func (router *AuthRouter) Logout(c *gin.Context) {
	// cookieService.DeleteRefreshToken(c.Writer)
	response := router.Service.Logout(c)
	c.JSON(http.StatusOK, response)
}

// @Summary Khôi phục tài khoản
// @Description Dùng để khôi phục số điện thoại khi mất hoặc đổi số mới
// @Accept  json
// @Produce  json
// @Param body body data.RestorePhoneRequest true "Restore Phone Request"
// @Tags Auth
// @Router /auth/restore [put]
func (router *AuthRouter) Restore(c *gin.Context) {
	var dto data.RestorePhoneRequest
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
// @Param request body data.SessionQRRequest true "Thông tin tạo QR"
// @Success 200 {object} data.SessionQRResponse
// @Router /auth/qr/init [post]
func (router *AuthRouter) QRInit(c *gin.Context) {
	var dto data.SessionQRRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sessionStr, _ := c.Cookie(config.C_SESSION_ID)
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
// @Param request body data.SessionQRConfirm true "Thông tin xác nhận QR"
// @Router /auth/qr/confirm [post]
func (router *AuthRouter) QRConfirm(c *gin.Context) {
	var dto data.SessionQRConfirm
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
// @Param request body data.SessionQRRequest true "Thông tin kiểm tra QR"
// @Router /auth/qr/status [post]
func (router *AuthRouter) QRStatus(c *gin.Context) {
	var dto data.SessionQRRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	sessionStr, _ := c.Cookie(config.C_SESSION_ID)
	sessionId, _ := strconv.ParseUint(sessionStr, 10, 32)

	result, err := router.Service.VerifyQRSession(c, &dto, sessionId)
	_routes.RouteResult(c, result, err)
}
