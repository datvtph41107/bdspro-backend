package properties

import (
	"user/config"
	"user/internal/dto"
)

// @bind: user/internal/interface/factory.IFactory
type RuntimeProperties struct {
	Properties dto.PropertiesDTO
}

func NewRuntimeProperties() *RuntimeProperties {
	return &RuntimeProperties{
		Properties: dto.PropertiesDTO{
			MaxTimesResend:       config.Properties.Otp.MaxTimesResend,
			MaxTimesOtpEnter:     config.Properties.Otp.MaxTimesOtpEnter,
			ExpiredAfterSec:      config.Properties.Otp.ExpiredAfterSec,
			LockVerifyAfter:      config.Properties.Otp.LockVerifyAfter,
			LockSendAfter:        config.Properties.Otp.LockSendAfter,
			LockRequestAfter:     config.Properties.Otp.LockRequestAfter,
			LimitOtpDevice:       config.Properties.Otp.LimitOtpDevice,
			LimitOtpIP:           config.Properties.Otp.LimitOtpIP,
			CooldownSeconds:      config.Properties.Otp.CooldownSeconds,
			IpCooldownSeconds:    config.Properties.Otp.IpCooldownSeconds,
			JwtRefreshExpMinutes: config.Properties.JWT.RefreshExpMinutes,
			Admin: struct {
				Username string
				Password string
				Email    string
				FullName string
			}{
				Username: config.Properties.Admin.Username,
				Password: config.Properties.Admin.Password,
				Email:    config.Properties.Admin.Email,
				FullName: config.Properties.Admin.FullName,
			},
			CSKH: struct {
				TeamIDs []uint64
			}{
				TeamIDs: config.Properties.CSKH.TeamIDs,
			},
			SMSTemplate: config.Properties.SMS.SMSTemplate,
		},
	}
}

func (r *RuntimeProperties) GetProperties() dto.PropertiesDTO {
	return r.Properties
}
