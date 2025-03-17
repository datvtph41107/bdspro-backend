package data

import "auth/models"

// AuthParam chứa các thông tin cần thiết để xác thực
type AuthParam struct {
	Status  *models.UserStatusEntity
	OTP     *models.UserOTPEntity
	OTPCode string
}
