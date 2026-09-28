package dto

import (
	"time"
)

// AdminAccessRequest đại diện cho request tạo/cập nhật quyền truy cập
type AdminAccessRequest struct {
	UserID        uint64     `json:"userId" binding:"required" validate:"required"`
	IPAddress     string     `json:"ipAddress" validate:"omitempty,ip"`
	IPRange       string     `json:"ipRange" validate:"omitempty,cidr"`
	DeviceID      string     `json:"deviceId" validate:"omitempty,min=1,max=255"`
	DeviceName    string     `json:"deviceName" validate:"omitempty,min=1,max=255"`
	DeviceType    string     `json:"deviceType" validate:"omitempty,oneof=mobile desktop tablet"`
	AuthType      string     `json:"authType" binding:"required" validate:"required,oneof=IP DEVICE OTP 2FA"`
	Status        string     `json:"status" validate:"omitempty,oneof=active blocked expired"`
	EffectiveFrom *time.Time `json:"effectiveFrom"`
	EffectiveTo   *time.Time `json:"effectiveTo"`
	MaxDevices    int        `json:"maxDevices" validate:"omitempty,min=1,max=10"`
	Require2FA    bool       `json:"require2fa"`
}

// AdminAccessResponse đại diện cho response quyền truy cập
type AdminAccessResponse struct {
	ID            uint64     `json:"id"`
	UserID        uint64     `json:"userId"`
	Username      string     `json:"username"`
	IPAddress     string     `json:"ipAddress"`
	IPRange       string     `json:"ipRange"`
	DeviceID      string     `json:"deviceId"`
	DeviceName    string     `json:"deviceName"`
	DeviceType    string     `json:"deviceType"`
	AuthType      string     `json:"authType"`
	Status        string     `json:"status"`
	EffectiveFrom *time.Time `json:"effectiveFrom"`
	EffectiveTo   *time.Time `json:"effectiveTo"`
	MaxDevices    int        `json:"maxDevices"`
	Require2FA    bool       `json:"require2fa"`
	CreatedBy     uint64     `json:"createdBy"`
	UpdatedBy     uint64     `json:"updatedBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// AdminAccessListRequest đại diện cho request lấy danh sách
type AdminAccessListRequest struct {
	Page      int                    `json:"page" validate:"omitempty,min=1"`
	Size      int                    `json:"size" validate:"omitempty,min=1,max=100"`
	UserID    *uint64                `json:"userId"`
	Status    string                 `json:"status" validate:"omitempty,oneof=active blocked expired"`
	AuthType  string                 `json:"authType" validate:"omitempty,oneof=IP DEVICE OTP 2FA"`
	IPAddress string                 `json:"ipAddress"`
	Filters   map[string]interface{} `json:"filters"`
}

// AdminAccessListResponse đại diện cho response danh sách
type AdminAccessListResponse struct {
	Data  []*AdminAccessResponse `json:"data"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}

// AdminAccessValidationRequest đại diện cho request validation
type AdminAccessValidationRequest struct {
	UserID    uint64 `json:"userId" binding:"required" validate:"required"`
	IPAddress string `json:"ipAddress" binding:"required" validate:"required,ip"`
	DeviceID  string `json:"deviceId" validate:"omitempty,min=1,max=255"`
}

// AdminAccessValidationResponse đại diện cho response validation
type AdminAccessValidationResponse struct {
	IsAllowed  bool   `json:"isAllowed"`
	Reason     string `json:"reason"`
	Require2FA bool   `json:"require2fa"`
	AccessID   uint64 `json:"accessId"`
	DeviceID   string `json:"deviceId"`
	IPAddress  string `json:"ipAddress"`
}

// AdminAccessLogRequest đại diện cho request tạo log
type AdminAccessLogRequest struct {
	UserID     uint64 `json:"userId" binding:"required" validate:"required"`
	IPAddress  string `json:"ipAddress" binding:"required" validate:"required,ip"`
	DeviceID   string `json:"deviceId" validate:"omitempty,min=1,max=255"`
	DeviceName string `json:"deviceName" validate:"omitempty,min=1,max=255"`
	Action     string `json:"action" binding:"required" validate:"required,oneof=login logout access_denied"`
	Status     string `json:"status" binding:"required" validate:"required,oneof=success failed blocked"`
	Reason     string `json:"reason" validate:"omitempty,min=1,max=500"`
	UserAgent  string `json:"userAgent" validate:"omitempty,min=1,max=500"`
}

// AdminAccessLogResponse đại diện cho response log
type AdminAccessLogResponse struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"`
	IPAddress  string    `json:"ipAddress"`
	DeviceID   string    `json:"deviceId"`
	DeviceName string    `json:"deviceName"`
	Action     string    `json:"action"`
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AdminAccessLogListRequest đại diện cho request lấy danh sách log
type AdminAccessLogListRequest struct {
	Page      int    `json:"page" validate:"omitempty,min=1"`
	Size      int    `json:"size" validate:"omitempty,min=1,max=100"`
	UserID    uint64 `json:"userId"`
	IPAddress string `json:"ipAddress"`
}

// AdminAccessLogListResponse đại diện cho response danh sách log
type AdminAccessLogListResponse struct {
	Data  []*AdminAccessLogResponse `json:"data"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"size"`
}
