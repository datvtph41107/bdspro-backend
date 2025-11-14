package dto

import (
	_dto "common/domain/dto"
	"time"
)

// HistoryAuthCreateDTO đại diện cho request tạo lịch sử thao tác của user trên auth-service
type HistoryAuthCreateDTO struct {
	UserID         uint64                 `json:"userId" binding:"required"`     // ID của user thực hiện
	OrganizationID *uint64                `json:"organizationId,omitempty"`      // Tổ chức đang hoạt động
	ActionType     string                 `json:"actionType" binding:"required"` // Nhóm hành động (login, logout,...)
	ActionName     string                 `json:"actionName" binding:"required"` // Tên hành động hiển thị
	Description    string                 `json:"description,omitempty"`         // Mô tả chi tiết
	Success        *bool                  `json:"success,omitempty"`             // Trạng thái thực hiện
	Reason         string                 `json:"reason,omitempty"`              // Lý do thất bại
	IPAddress      string                 `json:"ipAddress,omitempty"`           // Địa chỉ IP thực hiện
	UserAgent      string                 `json:"userAgent,omitempty"`           // User agent của thiết bị
	Metadata       map[string]interface{} `json:"metadata,omitempty"`            // Thông tin bổ sung
	PerformedBy    *uint64                `json:"performedBy,omitempty"`         // Người thực hiện (nếu khác user)
	SessionID      string                 `json:"sessionId,omitempty"`           // Mã session liên quan
	Channel        string                 `json:"channel,omitempty"`             // Kênh thực hiện (web, mobile,...)
	DeviceID       string                 `json:"deviceId,omitempty"`            // Device ID
	Location       string                 `json:"location,omitempty"`            // Thông tin vị trí
	AdditionalNote string                 `json:"additionalNote,omitempty"`      // Ghi chú thêm
	SourceService  string                 `json:"sourceService,omitempty"`       // Service nguồn ghi log
}

// HistoryAuthSearchDTO filter & pagination khi lấy lịch sử bảo mật
type HistoryAuthSearchDTO struct {
	_dto.Pagable
	UserID         uint64
	OrganizationID *uint64
	ActionType     string
	Channel        string
	DeviceID       string
	Success        *bool
	FromDate       *time.Time
	ToDate         *time.Time
}
