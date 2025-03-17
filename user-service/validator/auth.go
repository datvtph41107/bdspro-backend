package validator

import (
	"auth/config"
	"auth/data"
	"auth/enums"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthValidator struct {
	Properties *config.UserProperties
	// OtpRepository   repository.UserOTPRepository
	// StatusRepository repository.UserStatusRepository
}

func NewAuthValidator(props *config.UserProperties) *AuthValidator {
	return &AuthValidator{
		Properties: props,
		// OtpRepository:   otpRepo,
		// StatusRepository: statusRepo,
	}
}

func ValidateAuth(c *gin.Context, authCodes []enums.AuthCodeEnum, param *data.AuthParam) error {
	for _, code := range authCodes {
		if err := Validate(c, code, param); err != nil {
			return err
		}
	}
	return nil
}

func Validate(c *gin.Context, validCode enums.AuthCodeEnum, param *data.AuthParam) error {
	if param == nil {
		return &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}
	statusEntity := param.Status
	otpEntity := param.OTP
	otp := param.OTPCode

	code := int(validCode)
	switch validCode {
	case enums.LOCKED:
		if !statusEntity.LockedUntil.IsZero() && statusEntity.LockedUntil.After(time.Now()) {
			minutes := time.Until(statusEntity.LockedUntil).Minutes()
			return &_routes.Except{
				Code:    code,
				Message: "Tài khoản của bạn bị tạm khóa. Vui lòng thử lại sau " + fmt.Sprintf("%.0f phút", minutes),
			}
		}
	case enums.LIMIT_OTP_SEND:
		if otpEntity.OTPSendTime >= config.AppProperties.Otp.MaxTimesResend {
			LockAccount(c, param)
			return &_routes.Except{
				Code:    code,
				Message: "Đã quá 5 lần gửi OTP. Vui lòng thử lại sau " + fmt.Sprintf("%.0f phút", time.Until(statusEntity.LockedUntil).Minutes()),
			}
		}
	case enums.LIMIT_OTP_ENTER:
		if otpEntity.OTPCheckTime >= config.AppProperties.Otp.MaxTimesOtpEnter {
			LockAccount(c, param)
			return &_routes.Except{
				Code:    code,
				Message: "Đã quá 3 lần nhập. Vui lòng thử lại sau " + fmt.Sprintf("%.0f phút", time.Until(statusEntity.LockedUntil).Minutes()),
			}
		}
	case enums.CHECK_OTP:
		if otpEntity.OTP != otp {
			otpEntity.OTPCheckTime++
			_db.SaveWithAudit(c, otpEntity)
			return &_routes.Except{
				Code: code,
				Message: "Mã không chính xác. Bạn còn " + fmt.Sprintf("%d",
					config.AppProperties.Otp.MaxTimesOtpEnter+1-otpEntity.OTPCheckTime) +
					" lần. Mã OTP của bạn là " + otpEntity.OTP,
			}
		}
	case enums.EXPIRED:
		if otpEntity.ExpiredTime.IsZero() && otpEntity.ExpiredTime.Before(time.Now()) {
			return &_routes.Except{
				Code:    code,
				Message: "Mã đã hết hạn",
			}
		}
	case enums.LIMIT_NEXT_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(config.AppProperties.Otp.LockSendAfter) * time.Second)
			if time.Now().Before(nextTime) {
				return &_routes.Except{
					Code:    code,
					Message: "Hãy thử lại sau " + fmt.Sprintf("%.0f", time.Until(nextTime).Seconds()) + "s",
				}
			}
		}
	case enums.LIMIT_REQUEST_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(config.AppProperties.Otp.LockRequestAfter) * time.Second)
			if time.Now().Before(nextTime) {
				return &_routes.Except{
					Code:    code,
					Message: "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau " + fmt.Sprintf("%.0f", time.Until(nextTime).Seconds()) + "s",
				}
			}
		}
	case enums.ACTIVATED:
		if otpEntity.Activate {
			return &_routes.Except{
				Code:    code,
				Message: "OTP đã kích hoạt",
			}
		}
	case enums.INACTIVE_ACCOUNT:
		if otpEntity.Activate {
			return &_routes.Except{
				Code:    code,
				Message: "Tài khoản đang bị vô hiệu hóa",
			}
		}
	}
	return nil
}

func LockAccount(c *gin.Context, param *data.AuthParam) {
	statusEntity := param.Status
	otpEntity := param.OTP

	if statusEntity.LockedUntil.IsZero() || statusEntity.LockedUntil.Before(time.Now()) {
		otp := _utils.GenerateOTP()
		otpEntity.OTP = otp
		otpEntity.OTPCheckTime = 0
		otpEntity.OTPSendTime = 0
		statusEntity.LockedUntil = time.Now().Add(time.Duration(config.AppProperties.Otp.LockVerifyAfter) * time.Second)
		_db.SaveWithAudit(c, otpEntity)
		_db.SaveWithAudit(c, statusEntity)
	}
}
