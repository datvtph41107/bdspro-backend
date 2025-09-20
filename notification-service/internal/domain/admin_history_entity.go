package domain

import (
	_models "common/models"
	shared_enum "pb/enums"

	"github.com/lib/pq"
)

type AdminHistoryEntity struct {
	_models.BaseEntity
	TargetId   uint64                     `json:"targetId"`
	TargetType shared_enum.ETargetHistory `json:"targetType"`
	ActionType shared_enum.EHistory       `json:"actionType"`
	Title      string                     `json:"title,omitempty"`
	Note       pq.StringArray             `gorm:"type:text[]" json:"note,omitempty"`
	PreStage   string                     `json:"preStage,omitempty"`
	AfterStage string                     `json:"afterStage,omitempty"`
	AdminID    uint64                     `json:"adminId"`             // ID của admin thực hiện action
	OwnerID    *uint64                    `json:"ownerId"`             // ID của owner bị tác động
	OwnerOf    shared_enum.EOwnerType     `json:"ownerType"`           // Loại owner bị tác động
	IsInternal bool                       `json:"isInternal"`          // Có phải action nội bộ không
	AdminRole  string                     `json:"adminRole,omitempty"` // Vai trò của admin (super_admin, admin, moderator, etc.)
	IPAddress  string                     `json:"ipAddress,omitempty"` // IP address của admin
	UserAgent  string                     `json:"userAgent,omitempty"` // User agent của admin
	// Thêm thông tin user (không lưu trong DB, chỉ để hiển thị)
	AdminName   *string `json:"adminName,omitempty" gorm:"-"`
	AdminAvatar *string `json:"adminAvatar,omitempty" gorm:"-"`
	// Thêm trạng thái thành công/thất bại (không lưu trong DB, chỉ để hiển thị)
	Status *string `json:"status,omitempty" gorm:"-"`
}

func (AdminHistoryEntity) TableName() string {
	return "admin_histories"
}
