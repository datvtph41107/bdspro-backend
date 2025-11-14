package domain

import (
	_models "common/models"

	"gorm.io/datatypes"
)

// HistoryAuthEntity đại diện cho lịch sử thao tác của user trên auth-service
type HistoryAuthEntity struct {
	_models.BaseEntity
	ID             uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64            `gorm:"index;not null" json:"userId"`                                  // ID của user thực hiện hành động
	OrganizationID *uint64           `gorm:"index" json:"organizationId,omitempty"`                         // Tổ chức mà user đang thao tác
	ActionType     string            `gorm:"size:100;not null" json:"actionType"`                           // Nhóm hành động (login, logout, change_password,...)
	ActionName     string            `gorm:"size:255;not null" json:"actionName"`                           // Tên hành động hiển thị
	Description    string            `gorm:"type:text" json:"description,omitempty"`                        // Mô tả chi tiết hành động
	Success        bool              `gorm:"not null;default:true" json:"success"`                          // Trạng thái thực hiện
	Reason         string            `gorm:"size:255" json:"reason,omitempty"`                              // Lý do thất bại (nếu có)
	IPAddress      string            `gorm:"size:45" json:"ipAddress,omitempty"`                            // Địa chỉ IP thực hiện
	UserAgent      string            `gorm:"size:500" json:"userAgent,omitempty"`                           // User agent của phiên làm việc
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`             // Thông tin bổ sung dạng JSON
	PerformedBy    *uint64           `gorm:"index" json:"performedBy,omitempty"`                            // ID người thao tác (nếu khác userId, ví dụ admin)
	SessionID      string            `gorm:"size:255" json:"sessionId,omitempty"`                           // Mã session liên quan
	Channel        string            `gorm:"size:100" json:"channel,omitempty"`                             // Kênh thực hiện (web, mobile, api_key,...)
	DeviceID       string            `gorm:"size:255" json:"deviceId,omitempty"`                            // Device ID nếu có
	Location       string            `gorm:"size:255" json:"location,omitempty"`                            // Vị trí địa lý (geoip)
	AdditionalNote string            `gorm:"size:500" json:"additionalNote,omitempty"`                      // Ghi chú thêm
	SourceService  string            `gorm:"size:255;not null;default:'auth-service'" json:"sourceService"` // Service ghi log
}

// TableName xác định tên bảng cho GORM
func (HistoryAuthEntity) TableName() string {
	return "history_auth"
}
