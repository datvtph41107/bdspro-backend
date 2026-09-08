package config

// AppProperties cấu hình tổng thể của ứng dụng
type AppProperties struct {
	Server struct {
		TCPPort int `mapstructure:"tcp_port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
	OAuth struct {
		Google struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURI  string `mapstructure:"redirect_uri"`
		} `mapstructure:"google"`
		Facebook struct {
			ClientID     string `mapstructure:"client_id"`
			ClientSecret string `mapstructure:"client_secret"`
			RedirectURI  string `mapstructure:"redirect_uri"`
		} `mapstructure:"facebook"`
		Zalo struct {
			AppID         string `mapstructure:"app_id"`
			SecretKey     string `mapstructure:"secret_key"`
			RedirectURI   string `mapstructure:"redirect_uri"`
			PermissionURL string `mapstructure:"permission_url"`
		} `mapstructure:"zalo"`
	} `mapstructure:"oauth"`
	JWT struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
		Secret              string   `mapstructure:"secret"`
		AccessTokenExpiry   int      `mapstructure:"access_token_expiry"`
		RefreshTokenExpiry  int      `mapstructure:"refresh_token_expiry"`
	} `mapstructure:"jwt"`
	SMS struct {
		Endpoint    string `mapstructure:"endpoint"`
		Provider    string `mapstructure:"provider"`
		AccountSID  string `mapstructure:"account_sid"`
		AuthToken   string `mapstructure:"auth_token"`
		FromNumber  string `mapstructure:"from_number"`
		SMSTemplate string `mapstructure:"template"`
	} `mapstructure:"sms"`
	ZNS struct {
		AccessToken   string `mapstructure:"access_token"`
		RefreshToken  string `mapstructure:"refresh_token"`
		TemplateID    int    `mapstructure:"template_id"`
		OTPTemplateID int    `mapstructure:"otp_template_id"`
	} `mapstructure:"zns"`
	// OTP Config
	Otp struct {
		MaxTimesResend    int   `mapstructure:"max-times-resend"`
		MaxTimesOtpEnter  int   `mapstructure:"max-times-otp-enter"`
		ExpiredAfterSec   int64 `mapstructure:"expired-after-seconds"`
		LockVerifyAfter   int   `mapstructure:"lock-verify-after"`
		LockSendAfter     int   `mapstructure:"lock-send-after"`
		LockRequestAfter  int   `mapstructure:"lock-request-after"`
		LimitOtpDevice    int64 `mapstructure:"limit-otp-device"`
		LimitOtpIP        int64 `mapstructure:"limit-otp-ip"`
		CooldownSeconds   int   `mapstructure:"cooldown-seconds"`
		IpCooldownSeconds int   `mapstructure:"ip-cooldown-seconds"`
	} `mapstructure:"otp"`
	Admin struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Email    string `mapstructure:"email"`
		FullName string `mapstructure:"full_name"`
	} `mapstructure:"admin"`
	// CSKH Team IDs - danh sách profile ID của đội chăm sóc khách hàng
	CSKH struct {
		TeamIDs []uint64 `mapstructure:"team_ids"`
	} `mapstructure:"cskh"`
}

var Properties AppProperties

const C_SESSION_ID = "C_SESSION_ID"
