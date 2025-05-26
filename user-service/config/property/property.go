package property

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	// Zalo
	Zalo struct {
		App struct {
			ID     string `mapstructure:"id"`
			Secret string `mapstructure:"secret"`
		} `mapstructure:"app"`
		RedirectURI   string `mapstructure:"redirect_uri"`
		ProfileURL    string `mapstructure:"profile_url"`
		PermissionURL string `mapstructure:"permission_url"`
		SecretKey     string `mapstructure:"secret_key"`
	} `mapstructure:"zalo"`

	// Google
	Google struct {
		Client struct {
			ID     string `mapstructure:"id"`
			Secret string `mapstructure:"secret"`
		} `mapstructure:"client"`
		RedirectURI string `mapstructure:"redirect-uri"`
	} `mapstructure:"google"`

	// Facebook
	Facebook struct {
		Client struct {
			ID     string `mapstructure:"id"`
			Secret string `mapstructure:"secret"`
		} `mapstructure:"client"`
		RedirectURI string `mapstructure:"redirect-uri"`
	}

	// OTP Config
	Otp struct {
		MaxTimesResend    int   `mapstructure:"max-times-resend"`
		MaxTimesOtpEnter  int   `mapstructure:"max-times-otp-enter"`
		ExpiredAfterSec   int   `mapstructure:"expired-after-seconds"`
		LockVerifyAfter   int   `mapstructure:"lock-verify-after"`
		LockSendAfter     int   `mapstructure:"lock-send-after"`
		LockRequestAfter  int   `mapstructure:"lock-request-after"`
		LimitOtpDevice    int64 `mapstructure:"limit-otp-device"`
		LimitOtpIP        int64 `mapstructure:"limit-otp-ip"`
		CooldownSeconds   int   `mapstructure:"cooldown-seconds"`
		IpCooldownSeconds int   `mapstructure:"ip-cooldown-seconds"`
	} `mapstructure:"otp"`

	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
}

var AppProperties UserProperties

const C_SESSION_ID = "C_SESSION_ID"
