package dto

type PropertiesDTO struct {
	MaxTimesResend       int
	MaxTimesOtpEnter     int
	ExpiredAfterSec      int64
	LockVerifyAfter      int
	LockSendAfter        int
	LockRequestAfter     int
	LimitOtpDevice       int64
	LimitOtpIP           int64
	CooldownSeconds      int
	IpCooldownSeconds    int
	JwtRefreshExpMinutes uint64
	Admin                struct {
		Username string
		Password string
		Email    string
		FullName string
	}
	CSKH struct {
		TeamIDs []uint64
	}
	SMSTemplate string
}
