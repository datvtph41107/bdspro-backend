package routes

// @Tags Auth
// @Summary Gửi OTP
// @Description Gửi lại OTP hoặc submit lúc nhập SĐT
// @Accept  json
// @Produce  json
// @Param body body data.OtpRequest true "OTP Request"
// @Router /user/auth/otp/request [post]
func SendOtp() {}

// @Tags Auth
// @Summary Xác minh OTP
// @Description Xác minh OTP đã nhận
// @Accept  json
// @Produce  json
// @Param body body data.OtpVerifyRequest true "OTP Verify Request"
// @Router /user/auth/otp/verify [post]
func VerifyOtp() {}

// @Tags Auth
// @Summary Resend OTP
// @Description Gửi lại OTP mới
// @Accept  json
// @Produce  json
// @Param body body data.OtpRequest true "OTP Request"
// @Router /user/auth/otp/resend [post]
func ResendOtp() {}

// @Tags Auth
// @Summary Refresh token
// @Description Làm mới access token bằng refresh token trong cookie
// @Accept  json
// @Produce  json
// @Router /user/auth/token/refresh [post]
func RefreshToken() {}

// @Tags Auth
// @Summary Logout
// @Description Đăng xuất, xóa refresh token
// @Router /user/auth/logout [delete]
func Logout() {}

// @Summary Khôi phục tài khoản
// @Description Dùng để khôi phục số điện thoại khi mất hoặc đổi số mới
// @Accept  json
// @Produce  json
// @Param body body data.RestorePhoneRequest true "Restore Phone Request"
// @Tags Auth
// @Router /user/auth/restore [put]
func RestoreAccount() {}

// UpdateUserInfo cập nhật thông tin người dùng
// @Summary Cập nhật thông tin người dùng
// @Description API này cho phép người dùng cập nhật thông tin cá nhân
// @Tags User
// @Accept json
// @Produce json
// @Param request body data.UserInfoRequest true "Thông tin cập nhật"
// @Success 200 {object} data.UserInfoResponse "Cập nhật thành công"
// @Router /user/auth/set-info [post]
func UpdateInfo() {}
