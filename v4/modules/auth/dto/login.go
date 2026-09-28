package dto

// Login DTOs owned by modules/auth. Converted at the controller/service edge
// to/from legacy internal/dto while login migrates off AuthUsecase.

type OtpRequest struct {
	Phone    string `json:"phone"`
	Fullname string `json:"fullname,omitempty"`
	Mode     string `json:"mode"`
	AuthId   uint64 `json:"authId,omitempty"`
}

type OtpVerifyRequest struct {
	OtpRequest
	Otp             string `json:"otp"`
	AuthID          uint64 `json:"authId,omitempty"`
	Version         string `json:"version,omitempty"`
	Platform        string `json:"platform,omitempty"`
	OS              string `json:"os,omitempty"`
	DeviceName      string `json:"deviceName,omitempty"`
	DeviceID        string `json:"deviceId,omitempty"`
	Fullname        string `json:"fullname,omitempty"`
	ClientPublicKey string `json:"clientPublicKey,omitempty"`
}

type RequestLoginResponse struct {
	AuthID     uint64 `json:"authId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	OAuthID    string `json:"oauthId,omitempty"`
	Existed    bool   `json:"existed"`
	SmsChannel string `json:"smsChannel,omitempty"`
}

type ResendOTPResponse struct {
	AuthId     uint64 `json:"authId,omitempty"`
	Message    string `json:"message,omitempty"`
	SmsChannel string `json:"smsChannel,omitempty"`
}

type LoginWithPasswordRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
	OS         string `json:"os,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

type LoginWithKeyRequest struct {
	Phone            string `json:"phone"`
	EncryptedAuthKey string `json:"encryptedAuthKey"`
	ClientPublicKey  string `json:"clientPublicKey"`
	Version          string `json:"version,omitempty"`
	Platform         string `json:"platform,omitempty"`
	OS               string `json:"os,omitempty"`
	DeviceName       string `json:"deviceName,omitempty"`
	DeviceID         string `json:"deviceId,omitempty"`
}

type AdminLoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
	OS         string `json:"os,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

type AuthLoginResponse struct {
	AccessToken      string `json:"accessToken,omitempty"`
	AuthID           uint64 `json:"authId,omitempty"`
	Role             string `json:"role,omitempty"`
	Username         string `json:"username,omitempty"`
	ProfileID        uint64 `json:"profileId,omitempty"`
	RefreshToken     string `json:"refreshToken,omitempty"`
	OrganizationID   uint64 `json:"organizationId,omitempty"`
	FullName         string `json:"fullName,omitempty"`
	PasswordAuthId   uint64 `json:"passwordAuthId,omitempty"`
	ReferralCode     string `json:"referralCode,omitempty"`
	RequireFullname  bool   `json:"requireFullname"`
	TempToken        string `json:"tempToken,omitempty"`
	EncryptedAuthKey string `json:"encryptedAuthKey,omitempty"`
	ServerPublicKey  string `json:"serverPublicKey,omitempty"`
}

type RefreshTokenHTTPRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type PhoneCheckResponse struct {
	Existed bool    `json:"existed"`
	UserId  *uint64 `json:"userId,omitempty"`
}

// PhoneCheckPath is GET {prefix}/auth/phone/check/{phone}.
type PhoneCheckPath struct {
	Phone string `path:"phone"`
}
