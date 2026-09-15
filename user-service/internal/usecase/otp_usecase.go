package usecase

import (
	_errors "common/errors"
	_fault "common/fault"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"time"
	"user/internal/dto"
	"user/internal/enums"
	"user/internal/interface/factory"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
)

type OtpUsecase struct {
	otpRepo              repo.OTPRepository
	statusRepo           repo.StatusRepository
	properties           dto.PropertiesDTO
	notificationProvider providers.NotificationProvider
}

func NewOtpUsecase(
	otpRepo repo.OTPRepository,
	statusRepo repo.StatusRepository,
	properties factory.IFactory,
	notificationProvider providers.NotificationProvider,
) *OtpUsecase {
	return &OtpUsecase{
		otpRepo:              otpRepo,
		statusRepo:           statusRepo,
		properties:           properties.GetProperties(),
		notificationProvider: notificationProvider,
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
			s.logOtpLockHistory(c, param, "limit_otp_send", otpEntity.OTPSendTime, s.properties.MaxTimesResend)
			return _errors.ReturnError(
				int32(code),
				fmt.Sprintf("Đã quá %d lần gửi OTP. Vui lòng thử lại sau", s.properties.MaxTimesResend),
			)
		}
	case enums.LIMIT_OTP_ENTER:
		if otpEntity.OTPCheckTime > s.properties.MaxTimesOtpEnter {
			s.LockAccount(c, param)
			s.logOtpLockHistory(c, param, "limit_otp_enter", otpEntity.OTPCheckTime, s.properties.MaxTimesOtpEnter)
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
					" lần.",
			)
		}
	case enums.EXPIRED:
		if otpEntity.ExpiredTime.IsZero() && otpEntity.ExpiredTime.Before(time.Now()) {
			return _errors.ReturnError(
				int32(code),
				"Mã OTP đã hết hạn. Nhấn Gửi lại mã.",
			)
		}
	case enums.LIMIT_NEXT_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(s.properties.LockSendAfter) * time.Second)
			if time.Now().Before(nextTime) {
				seconds := int32(time.Until(nextTime).Seconds())
				return _fault.New(
					_fault.KindResourceExhausted,
					"user.otp.next_send_limited",
					"Hãy thử lại sau "+fmt.Sprintf("%d", seconds)+"s",
				).WithMetadata(map[string]string{
					"legacy_code": fmt.Sprintf("%d", code),
					"second":      fmt.Sprintf("%d", seconds),
				})
			}
		}
	case enums.LIMIT_REQUEST_TIME:
		if otpEntity.OTPDate != nil {
			nextTime := otpEntity.OTPDate.Add(time.Duration(s.properties.LockRequestAfter) * time.Second)
			if time.Now().Before(nextTime) {
				seconds := int32(time.Until(nextTime).Seconds())
				return _fault.New(
					_fault.KindResourceExhausted,
					"user.otp.request_limited",
					"Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau "+fmt.Sprintf("%d", seconds)+"s",
				).WithMetadata(map[string]string{
					"legacy_code": fmt.Sprintf("%d", code),
					"second":      fmt.Sprintf("%d", seconds),
				})
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

func (s *OtpUsecase) logOtpLockHistory(ctx context.Context, param *dto.AuthParam, reason string, current int, max int) {
	if s.notificationProvider == nil || param == nil {
		return
	}

	var userID uint64
	var phone string
	if param.Auth != nil {
		userID = param.Auth.UserID
		phone = param.Auth.AuthName
	}
	if userID == 0 && param.Status != nil {
		userID = param.Status.AuthID
	}

	authID := uint64(0)
	if param.Status != nil {
		authID = param.Status.AuthID
	}

	if userID == 0 && authID == 0 {
		return
	}

	ipAddress := _utils.AuditIPFromContext(ctx)
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}

	userAgent := _utils.GetUserAgentFromContext(ctx)
	if userAgent == "" {
		userAgent = "Mobile-Client"
	}

	description := "Tài khoản bị khóa do vượt giới hạn OTP"
	actionName := "OTP_LIMIT"
	switch reason {
	case "limit_otp_send":
		description = fmt.Sprintf("Tài khoản bị khóa do gửi OTP quá %d lần", max)
		actionName = "OTP_LIMIT_SEND"
	case "limit_otp_enter":
		description = fmt.Sprintf("Tài khoản bị khóa do nhập OTP sai quá %d lần", max)
		actionName = "OTP_LIMIT_ENTER"
	}

	metadata := map[string]interface{}{
		"authId":       authID,
		"reason":       reason,
		"currentCount": current,
		"maxCount":     max,
	}
	if phone != "" {
		metadata["phone"] = phone
	}
	if param.Status != nil && !param.Status.LockedUntil.IsZero() {
		metadata["lockedUntil"] = param.Status.LockedUntil.Format(time.RFC3339)
	}
	if param.OTP != nil {
		metadata["otpSendTime"] = param.OTP.OTPSendTime
		metadata["otpCheckTime"] = param.OTP.OTPCheckTime
	}

	success := false
	payload := dto.HistoryAuthCreateDTO{
		UserID:        userID,
		ActionType:    "otp_lock",
		ActionName:    actionName,
		Description:   description,
		Success:       &success,
		Reason:        reason,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SourceService: "auth-service",
		Metadata:      metadata,
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.notificationProvider.CreateHistoryAuth(timeoutCtx, &payload); err != nil {
		log.Printf("failed to log OTP lock history: %v", err)
	}
}
