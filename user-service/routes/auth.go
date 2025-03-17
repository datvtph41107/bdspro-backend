package routes

import (
	"auth/data"
	"auth/services"
	_routes "common/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title User API
// @version 1.0
// @description API quản lý người dùng
// @host localhost:8080
// @BasePath /user
func AuthRoutes(r *gin.Engine, authService *services.AuthService) {
	api := r.Group("/user/auth")

	api.POST("/otp/request", func(c *gin.Context) {
		var dto data.OtpRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, _ := authService.RequestOtp(c, dto)
		c.JSON(http.StatusOK, response)
	})

	// @Summary Xác minh OTP
	// @Description Đến bước nhập OTP thì gọi API này để kích hoạt tài khoản
	// @Accept  json
	// @Produce  json
	// @Param body body data.OtpVerifyRequest true "OTP Verify Request"
	// @Success 200 {object} data.AuthLoginResponse
	// @Router /auth/otp/verify [post]
	api.POST("/otp/verify", func(c *gin.Context) {
		var dto data.OtpVerifyRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := authService.VerifyOtp(c, dto)
		// cookieService.SetRefreshToken(c.Writer, result.RefreshToken)
		// c.JSON(http.StatusOK, result)
		_routes.RouteResult(c, result, err)
	})

	// @Summary Resend OTP
	// @Description Gửi OTP mới
	// @Accept  json
	// @Produce  json
	// @Param body body data.OtpRequest true "OTP Request"
	// @Success 200 {object} data.ResponseDTO
	// @Router /auth/otp/resend [post]
	api.POST("/otp/resend", func(c *gin.Context) {
		var dto data.OtpRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := authService.ResendOTP(c, dto)
		_routes.RouteResult(c, result, err)
	})

	// @Summary Refresh token
	// @Description API làm mới access token
	// @Param rt cookie string false "Refresh Token"
	// @Success 200 {object} data.ResponseDTO
	// @Router /auth/token/refresh [post]
	api.POST("/token/refresh", func(c *gin.Context) {
		refreshToken, _ := c.Cookie("rt")
		result, err := authService.RefreshToken(refreshToken)
		_routes.RouteResult(c, result, err)
	})

	// @Summary Đăng xuất
	// @Description Xóa refresh token
	// @Success 200 {object} data.ResponseDTO
	// @Router /auth/logout [delete]
	api.DELETE("/logout", func(c *gin.Context) {
		// cookieService.DeleteRefreshToken(c.Writer)
		response := authService.Logout(c)
		c.JSON(http.StatusOK, response)
	})

	api.POST("/set-info", func(c *gin.Context) {
		var dto data.UserInfoRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := authService.UpdateUserInfo(c, dto)
		_routes.RouteResult(c, result, err)
	})

	api.PUT("/restore", func(c *gin.Context) {
		var dto data.RestorePhoneRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}
		result, err := authService.Restore(c, dto)
		_routes.RouteResult(c, result, err)
	})

}
