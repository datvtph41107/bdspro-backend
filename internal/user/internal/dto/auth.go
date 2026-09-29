package dto

import "user/internal/domain/auth"

// Hàm khởi tạo AuthLoginResponse từ UserAuthEntity
func MakeAuthLoginResponse(accessToken string, authEntity *auth.AuthMethod, role string) *AuthLoginResponse {
	return &AuthLoginResponse{
		AccessToken: accessToken,
		// RefreshToken: refreshToken,
		Role:            role,
		Username:        authEntity.AuthName,
		AuthID:          authEntity.ID,
		ProfileID:       authEntity.UserID,
		FullName:        authEntity.FullName, // OAuth fullname
		ServerPublicKey: authEntity.PublicKey,
		// EncryptedAuthKey: authEntity.AuthKey,
	}
}

// OtpRequest đại diện cho yêu cầu gửi OTP
// @swagger:model
type OtpRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164"`
	Fullname string `json:"fullname,omitempty"`
	Mode     string `json:"mode" binding:"required" validate:"required"`
	AuthId   uint64 `json:"authId,omitempty"`
}

// OtpVerifyRequest mở rộng từ OtpRequest để xác minh OTP
// @swagger:model
type OtpVerifyRequest struct {
	OtpRequest
	Otp        string `json:"otp" binding:"required" validate:"required"`
	AuthID     uint64 `json:"authId,omitempty"` // đối với oauth thì truyền -> sẽ link oauth với tài khoản otp, tạo profile mới
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
	OS         string `json:"os,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
	Fullname   string `json:"fullname,omitempty"`

	ClientPublicKey string `json:"clientPublicKey,omitempty"` // Diffie-Hellman public key từ client
}

// RequestOTPResponse đại diện cho response gửi OTP
type RequestOTPResponse struct {
	AuthId  uint64 `json:"authId"`
	Message string `json:"message"`
}

// VerifyOTPResponse đại diện cho response xác minh OTP
type VerifyOTPResponse struct {
	AuthId       uint64 `json:"authId"`
	ProfileId    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

// ResendOTPRequest đại diện cho yêu cầu gửi lại OTP
type ResendOTPRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164"`
	Fullname string `json:"fullname,omitempty"`
}

// ResendOTPResponse đại diện cho response gửi lại OTP
type ResendOTPResponse struct {
	AuthId     uint64 `json:"authId"`
	Message    string `json:"message"`
	SmsChannel string `json:"smsChannel"` // "zalo" nếu gửi ZNS thành công, "sms" nếu ZNS lỗi
}

// RefreshTokenRequest đại diện cho yêu cầu làm mới token
// @swagger:model
type RefreshTokenRequest struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

// RefreshTokenResponse đại diện cho response làm mới token
// @swagger:model
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutRequest đại diện cho yêu cầu đăng xuất
// @swagger:model
type LogoutRequest struct {
	// Empty request
}

// LogoutResponse đại diện cho response đăng xuất
type LogoutResponse struct {
	Message string `json:"message"`
}

// DeleteAccountRequest đại diện cho yêu cầu xóa tài khoản
type DeleteAccountRequest struct {
	// Empty request
}

// DeleteAccountResponse đại diện cho response xóa tài khoản
type DeleteAccountResponse struct {
	Message string `json:"message"`
}

// RestorePhoneRequest đại diện cho yêu cầu khôi phục số điện thoại
// @swagger:model
type RestorePhoneRequest struct {
	Phone string `json:"phone" binding:"required" validate:"required,e164"`
}

// RestoreAccountResponse đại diện cho response khôi phục tài khoản
type RestoreAccountResponse struct {
	AuthId   uint64            `json:"authId"`
	Provider auth.ProviderType `json:"provider"`
	AuthName string            `json:"authName"`
}

// RestoreDeletedAccountRequest đại diện cho yêu cầu khôi phục tài khoản đã bị xóa
type RestoreDeletedAccountRequest struct {
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
}

// RestoreDeletedAccountResponse đại diện cho response khôi phục tài khoản đã bị xóa
type RestoreDeletedAccountResponse struct {
	AuthId   uint64            `json:"authId"`
	Provider auth.ProviderType `json:"provider"`
	AuthName string            `json:"authName"`
	Message  string            `json:"message"`
}

// SessionQRRequest đại diện cho yêu cầu tạo QR session
// @swagger:model
type SessionQRRequest struct {
	ClientID   string `json:"clientId"`
	UserAgent  string `json:"userAgent"`
	CsrfToken  string `json:"csrfToken"`
	Expiration int64  `json:"expiration"`
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// SessionQRResponse đại diện cho response QR session
// @swagger:model
type SessionQRResponse struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// SessionQRConfirm đại diện cho yêu cầu xác nhận QR
// @swagger:model
type SessionQRConfirm struct {
	SessionId  uint64 `json:"sessionId"`
	SessionKey string `json:"sessionKey"`
}

// QRConfirmResponse đại diện cho response xác nhận QR
type QRConfirmResponse struct {
	Message string `json:"message"`
}

// QRStatusRequest đại diện cho yêu cầu kiểm tra trạng thái QR
type QRStatusRequest struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// QRStatusResponse đại diện cho response trạng thái QR
type QRStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AuthLoginResponse đại diện cho response đăng nhập
type AuthLoginResponse struct {
	AccessToken    string `json:"accessToken"`
	AuthID         uint64 `json:"authId"`
	Role           string `json:"role,omitempty"`
	Username       string `json:"username,omitempty"`
	ProfileID      uint64 `json:"profileId"`
	RefreshToken   string `json:"refreshToken"`
	OrganizationID uint64 `json:"organizationId,omitempty"`
	FullName       string `json:"fullName"`                 // OAuth fullname
	PasswordAuthId uint64 `json:"passwordAuthId,omitempty"` // ID của auth_method có provider ADMIN
	ReferralCode   string `json:"referralCode,omitempty"`   // Mã giới thiệu của user

	RequireFullname bool   `json:"requireFullname"`     // true nếu user cần cập nhật họ tên
	TempToken       string `json:"tempToken,omitempty"` // Token tạm thời (10 phút) cho user chưa có họ tên

	// Diffie-Hellman key exchange fields
	EncryptedAuthKey string `json:"encryptedAuthKey,omitempty"` // Auth key đã được mã hóa bằng shared secret
	ServerPublicKey  string `json:"serverPublicKey,omitempty"`  // Diffie-Hellman public key của server
	// Avatar         string `json:"avatar"`
	// Phone          string `json:"phone"`
	// Email          string `json:"email"`
}

// OAuthLoginRequest đại diện cho yêu cầu đăng nhập OAuth
type OAuthLoginRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
	Code     string `json:"code" binding:"required" validate:"required"`
	State    string `json:"state"`
	Platform string `json:"platform"`
}

// OAuthLoginResponse đại diện cho response đăng nhập OAuth
type OAuthLoginResponse struct {
	AuthID       uint64 `json:"authId"`
	ProfileID    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Provider     string `json:"provider"`
}

// GetLoginURLRequest đại diện cho yêu cầu lấy URL đăng nhập OAuth
type GetLoginURLRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
}

// GetLoginURLResponse đại diện cho response URL đăng nhập OAuth
type GetLoginURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

// HandleCallbackRequest đại diện cho yêu cầu xử lý OAuth callback
type HandleCallbackRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
	Code     string `json:"code" binding:"required" validate:"required"`
	State    string `json:"state"`
	Platform string `json:"platform"`
}

// HandleCallbackResponse đại diện cho response xử lý OAuth callback
type HandleCallbackResponse struct {
	AuthId       uint64 `json:"authId"`
	ProfileId    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Provider     string `json:"provider"`
}

// AdminLoginRequest đại diện cho yêu cầu đăng nhập admin
type AdminLoginRequest struct {
	Username   string `json:"username" binding:"required" validate:"required"`
	Password   string `json:"password" binding:"required" validate:"required"`
	Version    string `json:"version" binding:"required" validate:"required"`
	Platform   string `json:"platform" binding:"required" validate:"required"`
	OS         string `json:"os" binding:"required" validate:"required"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

// LoginWithKeyRequest đại diện cho yêu cầu đăng nhập bằng auth key (Diffie-Hellman)
type LoginWithKeyRequest struct {
	Phone            string `json:"phone" binding:"required" validate:"required,e164"`
	EncryptedAuthKey string `json:"encryptedAuthKey" binding:"required" validate:"required"` // Auth key đã mã hóa
	ClientPublicKey  string `json:"clientPublicKey" binding:"required" validate:"required"`  // DH public key của client
	Version          string `json:"version" binding:"required" validate:"required"`
	Platform         string `json:"platform" binding:"required" validate:"required"`
	OS               string `json:"os" binding:"required" validate:"required"`
	DeviceName       string `json:"deviceName,omitempty"`
	DeviceID         string `json:"deviceId,omitempty"`
}

// Account Lock DTOs
// LockAccountRequest đại diện cho yêu cầu khóa tài khoản
type LockAccountRequest struct {
	ProfileID uint64             `json:"profileId" binding:"required" validate:"required"`
	LockType  auth.AccountStatus `json:"lockType" binding:"required" validate:"required"` // temporary hoặc permanent, 10: active, 40: inactive, 20: temporarily_locked, 30: permanently_locked
	Duration  int                `json:"duration,omitempty"`                              // Số giờ khóa (chỉ dùng cho temporary)
	Reason    string             `json:"reason" binding:"required" validate:"required,min=10,max=500"`
	LockedBy  uint64             `json:"lockedBy" binding:"required" validate:"required"` // ID admin thực hiện khóa
}

// UnlockAccountRequest đại diện cho yêu cầu mở khóa tài khoản
type UnlockAccountRequest struct {
	ProfileID  uint64 `json:"profileId" binding:"required" validate:"required"`
	Reason     string `json:"reason,omitempty" validate:"omitempty,min=5,max=500"`
	UnlockedBy uint64 `json:"unlockedBy" binding:"required" validate:"required"` // ID admin thực hiện mở khóa
}

// LockAccountResponse đại diện cho response khóa tài khoản
type LockAccountResponse struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	ProfileID   uint64             `json:"profileId"`
	LockType    auth.AccountStatus `json:"lockType"`
	LockedAt    string             `json:"lockedAt,omitempty"`
	LockedUntil string             `json:"lockedUntil,omitempty"` // Chỉ có khi lockType = temporary
}

// UnlockAccountResponse đại diện cho response mở khóa tài khoản
type UnlockAccountResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	ProfileID  uint64 `json:"profileId"`
	UnlockedAt string `json:"unlockedAt"`
}

// GetAccountStatusRequest đại diện cho yêu cầu lấy trạng thái tài khoản
type GetAccountStatusRequest struct {
	ProfileID uint64 `json:"profileId" binding:"required" validate:"required"`
}

// GetAccountStatusResponse đại diện cho response trạng thái tài khoản
type GetAccountStatusResponse struct {
	ProfileID   uint64             `json:"profileId"`
	Status      uint32             `json:"status"`
	StatusText  string             `json:"statusText"`
	IsLocked    bool               `json:"isLocked"`
	LockType    auth.AccountStatus `json:"lockType,omitempty"`
	LockedAt    string             `json:"lockedAt,omitempty"`
	LockedUntil string             `json:"lockedUntil,omitempty"`
	LockReason  string             `json:"lockReason,omitempty"`
	LockedBy    *uint64            `json:"lockedBy,omitempty"`
	CanLogin    bool               `json:"canLogin"`
	IsExpired   bool               `json:"isExpired,omitempty"` // Chỉ có khi lockType = temporary
}

// Batch Lock DTOs
// LockMultipleUsersRequest đại diện cho yêu cầu khóa nhiều user
type LockMultipleUsersRequest struct {
	ProfileIDs []uint64           `json:"profileIds" binding:"required" validate:"required,min=1"`
	LockType   auth.AccountStatus `json:"lockType" binding:"required" validate:"required"`
	Duration   int                `json:"duration,omitempty"` // Số giờ khóa (chỉ dùng cho temporary)
	Reason     string             `json:"reason" binding:"required" validate:"required,min=10,max=500"`
	LockedBy   uint64             `json:"lockedBy" binding:"required" validate:"required"`
}

// UnlockMultipleUsersRequest đại diện cho yêu cầu mở khóa nhiều user
type UnlockMultipleUsersRequest struct {
	ProfileIDs []uint64 `json:"profileIds" binding:"required" validate:"required,min=1"`
	Reason     string   `json:"reason,omitempty" validate:"omitempty,min=5,max=500"`
	UnlockedBy uint64   `json:"unlockedBy" binding:"required" validate:"required"`
}

// LockMultipleUsersResponse đại diện cho response khóa nhiều user
type LockMultipleUsersResponse struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	ProfileIDs  []uint64           `json:"profileIds"`
	LockType    auth.AccountStatus `json:"lockType"`
	LockedAt    string             `json:"lockedAt,omitempty"`
	LockedUntil string             `json:"lockedUntil,omitempty"`
}

// UnlockMultipleUsersResponse đại diện cho response mở khóa nhiều user
type UnlockMultipleUsersResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	ProfileIDs []uint64 `json:"profileIds"`
	UnlockedAt string   `json:"unlockedAt"`
}

// GetAuthUsersByProfileIDsRequest request lấy AuthUser theo danh sách profile IDs
type GetAuthUsersByProfileIDsRequest struct {
	ProfileIDs []uint64 `json:"profileIds" binding:"required" validate:"required,min=1"`
}

// GetAuthUsersByProfileIDsResponse response lấy AuthUser theo danh sách profile IDs
type GetAuthUsersByProfileIDsResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Users   []AuthUserStatusInfo `json:"users"`
}

// AuthUserStatusInfo thông tin status của AuthUser
type AuthUserStatusInfo struct {
	ProfileID   uint64             `json:"profileId"`
	Status      auth.AccountStatus `json:"status"`
	StatusText  string             `json:"statusText"`
	IsLocked    bool               `json:"isLocked"`
	LockType    auth.AccountStatus `json:"lockType,omitempty"`
	LockedAt    string             `json:"lockedAt,omitempty"`
	LockedUntil string             `json:"lockedUntil,omitempty"`
	LockReason  string             `json:"lockReason,omitempty"`
	LockedBy    uint64             `json:"lockedBy,omitempty"`
	CanLogin    bool               `json:"canLogin"`
	IsExpired   bool               `json:"isExpired,omitempty"`
}

// CreatePasswordRequest đại diện cho yêu cầu tạo mật khẩu
type CreatePasswordRequest struct {
	Password        string `json:"password" binding:"required" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirmPassword" binding:"required" validate:"required"`
}

// CreatePasswordResponse đại diện cho response tạo mật khẩu
type CreatePasswordResponse struct {
	AuthId  uint64 `json:"authId"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// LoginWithPasswordRequest đại diện cho yêu cầu đăng nhập bằng mật khẩu
type LoginWithPasswordRequest struct {
	Username   string `json:"username" binding:"required" validate:"required"` // Phone hoặc username
	Password   string `json:"password" binding:"required" validate:"required"`
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
	OS         string `json:"os,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

// SwitchAccountRequest đại diện cho yêu cầu chuyển đổi tài khoản
type SwitchAccountRequest struct {
	ProfileID      uint64  `json:"profileId" binding:"required" validate:"required"`
	OrganizationID *uint64 `json:"organizationId,omitempty"`
	GroupID        *uint32 `json:"groupId,omitempty"`
	OwnerType      *uint32 `json:"ownerType,omitempty"`
}

// SwitchAccountResponse đại diện cho response chuyển đổi tài khoản
type SwitchAccountResponse struct {
	AuthID         uint64  `json:"authId"`
	ProfileID      uint64  `json:"profileId"`
	OrganizationID *uint64 `json:"organizationId,omitempty"`
	GroupID        *uint32 `json:"groupId,omitempty"`
	Role           string  `json:"role"`
	FullName       string  `json:"fullName"`
	Avatar         string  `json:"avatar"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	AccessToken    string  `json:"accessToken"`
}

// GetLoginHistoryRequest đại diện cho request lấy lịch sử đăng nhập
type GetLoginHistoryRequest struct {
	Page int32 `json:"page" form:"page"`
	Size int32 `json:"size" form:"size"`
}

// GetLoginHistoryResponse đại diện cho response lịch sử đăng nhập
type GetLoginHistoryResponse struct {
	Data  []LoginHistoryItem `json:"data"`
	Total int64              `json:"total"`
}

// LoginHistoryItem đại diện cho một mục lịch sử đăng nhập
type LoginHistoryItem struct {
	SessionID    uint64 `json:"sessionId"`
	AuthID       uint64 `json:"authId"`
	DeviceID     string `json:"deviceId"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
	OS           string `json:"os"`
	DeviceName   string `json:"deviceName"`
	CreatedDate  string `json:"createdDate"`
	FinishedDate string `json:"finishedDate"`
	IPRequest    string `json:"ipRequest"`
	TotalRequest uint64 `json:"totalRequest"`
	UserAgent    string `json:"userAgent"`
	Activate     bool   `json:"activate"`
}

// ConnectedAccountDTO đại diện cho một tài khoản MXH đã liên kết
type ConnectedAccountDTO struct {
	AuthID    uint64 `json:"authId"`
	Provider  string `json:"provider"` // GOOGLE, FACEBOOK, ZALO
	AuthName  string `json:"authName"` // OAuth ID từ provider
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"createdAt"`
}

// PIN DTOs

// CreatePINRequest đại diện cho yêu cầu tạo mã PIN
type CreatePINRequest struct {
	PIN   string `json:"pin" binding:"required" validate:"required,len=6,numeric"`
	Phone string `json:"phone" binding:"required" validate:"required,e164"`
	OTP   string `json:"otp" binding:"required" validate:"required,len=6,numeric"`
}

// CreatePINResponse đại diện cho response tạo mã PIN
type CreatePINResponse struct {
	AuthId  uint64 `json:"authId"`
	Message string `json:"message"`
}

// VerifyPINRequest đại diện cho yêu cầu xác minh mã PIN
type VerifyPINRequest struct {
	Phone      string `json:"phone" binding:"required" validate:"required,e164"`
	PIN        string `json:"pin" binding:"required" validate:"required,len=6,numeric"`
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
	OS         string `json:"os,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
}

// UpdatePINRequest đại diện cho yêu cầu cập nhật mã PIN
type UpdatePINRequest struct {
	OldPIN string `json:"oldPin" binding:"required" validate:"required,len=6,numeric"`
	NewPIN string `json:"newPin" binding:"required" validate:"required,len=6,numeric"`
}

// UpdatePINResponse đại diện cho response cập nhật mã PIN
type UpdatePINResponse struct {
	Message string `json:"message"`
}

// ActivePINRequest đại diện cho yêu cầu bật/tắt PIN
type ActivePINRequest struct {
	Active bool   `json:"active"`
	PIN    string `json:"pin" binding:"required" validate:"required,len=6,numeric"`
}

// CheckPINExistsRequest đại diện cho yêu cầu kiểm tra mã PIN
type CheckPINExistsRequest struct {
	// Empty - lấy từ token
}

// CheckPINExistsResponse đại diện cho response kiểm tra mã PIN
type CheckPINExistsResponse struct {
	Exists bool   `json:"exists"`
	AuthId uint64 `json:"authId"`
}

// AuthParam dùng để validate PIN
type PINAuthParam struct {
	Status  *auth.UserStatusEntity
	PIN     *auth.UserPINEntity
	PINCode string
}

// PinVerifyMeta chứa thông tin bổ sung để log/validate PIN
type PinVerifyMeta struct {
	Phone      string
	Platform   string
	OS         string
	DeviceID   string
	DeviceName string
}
