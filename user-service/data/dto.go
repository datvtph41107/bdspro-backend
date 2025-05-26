package data

import auth_model "user/models/auth"

// AuthParam chứa các thông tin cần thiết để xác thực
type AuthParam struct {
	Status  *auth_model.UserStatusEntity
	OTP     *auth_model.UserOTPEntity
	OTPCode string
}
