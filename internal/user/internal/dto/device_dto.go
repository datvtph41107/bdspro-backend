package dto

import "time"

// CreateDeviceRequest chứa thông tin thiết bị gửi từ client
type CreateDeviceRequest struct {
	DeviceID     string  `json:"deviceId" validate:"required,max=255"`
	DeviceName   string  `json:"deviceName" validate:"omitempty,max=255"`
	DeviceType   string  `json:"deviceType" validate:"omitempty,oneof=mobile desktop tablet web"`
	Platform     string  `json:"platform" validate:"omitempty,max=100"`
	OSVersion    string  `json:"osVersion" validate:"omitempty,max=100"`
	AppVersion   string  `json:"appVersion" validate:"omitempty,max=50"`
	BuildNumber  string  `json:"buildNumber" validate:"omitempty,max=50"`
	Manufacturer string  `json:"manufacturer" validate:"omitempty,max=100"`
	Model        string  `json:"model" validate:"omitempty,max=100"`
	Locale       string  `json:"locale" validate:"omitempty,max=20"`
	Timezone     string  `json:"timezone" validate:"omitempty,max=50"`
	PushToken    string  `json:"pushToken" validate:"omitempty,max=255"`
	IPAddress    string  `json:"ipAddress" validate:"omitempty,max=64"`
	UserAgent    string  `json:"userAgent" validate:"omitempty,max=512"`
	LastSeenAt   *string `json:"lastSeenAt,omitempty"`
}

// DeviceResponse đại diện cho thông tin thiết bị trả về client
type DeviceResponse struct {
	ID             uint64     `json:"id"`
	DeviceID       string     `json:"deviceId"`
	DeviceName     string     `json:"deviceName,omitempty"`
	DeviceType     string     `json:"deviceType,omitempty"`
	Platform       string     `json:"platform,omitempty"`
	OSVersion      string     `json:"osVersion,omitempty"`
	AppVersion     string     `json:"appVersion,omitempty"`
	BuildNumber    string     `json:"buildNumber,omitempty"`
	Manufacturer   string     `json:"manufacturer,omitempty"`
	Model          string     `json:"model,omitempty"`
	Locale         string     `json:"locale,omitempty"`
	Timezone       string     `json:"timezone,omitempty"`
	PushToken      string     `json:"pushToken,omitempty"`
	IPAddress      string     `json:"ipAddress,omitempty"`
	UserAgent      string     `json:"userAgent,omitempty"`
	ProfileID      *uint64    `json:"profileId,omitempty"`
	OrganizationID *uint64    `json:"organizationId,omitempty"`
	AuthID         *uint64    `json:"authId,omitempty"`
	LastSeenAt     *time.Time `json:"lastSeenAt,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
}
