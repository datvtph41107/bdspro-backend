package usecase

import (
	"auth/internal/dto"
	"auth/internal/enums"
	"auth/internal/interface/factory"
	"auth/internal/interface/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"time"
)

type OtpUsecase struct {
	otpRepo    repo.OTPRepository
	statusRepo repo.StatusRepository
	properties dto.PropertiesDTO
}

func NewOtpUsecase(
	otpRepo repo.OTPRepository,
	statusRepo repo.StatusRepository,
	properties factory.IFactory,
) *OtpUsecase {
	return &OtpUsecase{
		otpRepo:    otpRepo,
		statusRepo: statusRepo,
		properties: properties.GetProperties(),
	}
}

func (s *OtpUsecase) ValidateAuth(c context.Context, authCodes []enums.AuthCodeEnum, param *dto.AuthParam) error {
	for _, code := range authCodes {
		if err := s.Validate(c, code, param); err != nil {
			return err
		}
	}
	return nil
}

func (s *OtpUsecase) Validate(c context.Context, validCode enums.AuthCodeEnum, param *dto.AuthParam) error {
	if param == nil {
		return _errors.ReturnError(
			int32(400),
			"Người dùng không tồn tại",
		)
	}
	statusEntity := param.Status
	otpEntity := param.OTP
	otp := param.OTPCode

	code := int(validCode)
	switch validCode {
	case enums.LOCKED:
		if !statusEntity.LockedUntil.IsZero() && statusEntity.LockedUntil.After(time.Now()) {
			minutes := time.Until(statusEntity.LockedUntil).Minutes()
			return _errors.ReturnError(int32(code), "Tài khoản của bạn bị tạm khóa. Vui lòng thử lại sau "+fmt.Sprintf("%.0f phút", minutes))
		}
	case enums.LIMIT_OTP_SEND:
		if otpEntity.OTPSendTime >= s.properties.MaxTimesResend {
			s.LockAccount(c, param)
			return _errors.ReturnError(
				int32(code),
				fmt.Sprintf("Đã quá %d lần gửi OTP. Vui lòng thử lại sau", s.properties.MaxTimesResend),
			)
		}
	case enums.LIMIT_OTP_ENTER:
		if otpEntity.OTPCheckTime >= s.properties.MaxTimesOtpEnter {
			s.LockAccount(c, param)
			return _errors.ReturnError(
				int32(code),
				"Đã quá 3 lần nhập. Vui lòng thử lại sau "+fmt.Sprintf("%.0f phút", time.Until(statusEntity.LockedUntil).Minutes()),
			)
		}
	case enums.CHECK_OTP:
		if otpEntity.OTP != otp {
			otpEntity.OTPCheckTime++
			s.otpRepo.UpdateOTP(c, otpEntity)
			return _errors.ReturnError(
				int32(code),
				"Mã không chính xác. Bạn còn "+fmt.Sprintf("%d",
					s.properties.MaxTimesOtpEnter+1-otpEntity.OTPCheckTime)+
					" lần. Mã OTP của bạn là "+otpEntity.OTP,
			)
		}
	case enums.EXPIRED:
		if otpEntity.ExpiredTime.IsZero() && otpEntity.ExpiredTime.Before(time.Now()) {
			return _errors.ReturnError(
				int32(code),
				"Mã đã hết hạn",
			)
		}
	case enums.LIMIT_NEXT_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(s.properties.LockSendAfter) * time.Second)
			if time.Now().Before(nextTime) {
				return _errors.ReturnError(
					int32(code),
					"Hãy thử lại sau "+fmt.Sprintf("%.0f", time.Until(nextTime).Seconds())+"s",
				)
			}
		}
	case enums.LIMIT_REQUEST_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(s.properties.LockRequestAfter) * time.Second)
			if time.Now().Before(nextTime) {
				return _errors.ReturnError(
					int32(code),
					"Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau "+fmt.Sprintf("%.0f", time.Until(nextTime).Seconds())+"s",
				)
			}
		}
	case enums.ACTIVATED:
		if otpEntity.Activate {
			return _errors.ReturnError(
				int32(code),
				"OTP đã kích hoạt",
			)
		}
	case enums.INACTIVE_ACCOUNT:
		if otpEntity.Activate {
			return _errors.ReturnError(
				int32(code),
				"Tài khoản đang bị vô hiệu hóa",
			)
		}
	}
	return nil
}

func (s *OtpUsecase) LockAccount(c context.Context, param *dto.AuthParam) {
	statusEntity := param.Status
	otpEntity := param.OTP

	if statusEntity.LockedUntil.IsZero() || statusEntity.LockedUntil.Before(time.Now()) {
		otp := _utils.GenerateOTP()
		otpEntity.OTP = otp
		otpEntity.OTPCheckTime = 0
		otpEntity.OTPSendTime = 0
		statusEntity.LockedUntil = time.Now().Add(time.Duration(
			s.properties.LockVerifyAfter) * time.Second,
		)
		s.otpRepo.UpdateOTP(c, otpEntity)
		s.statusRepo.UpdateStatus(c, statusEntity)
	}
}
