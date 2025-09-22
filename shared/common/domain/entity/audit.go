package _models

import (
	_enums "common/domain/enum"
	"log"

	"gorm.io/gorm"
)

// @swagger:model
type AuditBase struct {
	// gorm.Model
	CreatedBy *uint64 `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy *uint64 `gorm:"column:updated_by" json:"updatedBy,omitempty"`
}

// GORM Callback: Gán `CreatedBy` và `UpdatedBy` khi tạo mới
func (base *AuditBase) BeforeCreate(tx *gorm.DB) error {

	userID := GetCurrentUserID(tx) // Lấy user ID từ context hoặc request
	base.CreatedBy = &userID
	base.UpdatedBy = &userID

	return nil
}

// GORM Callback: Cập nhật `UpdatedBy` khi có thay đổi
func (base *AuditBase) BeforeUpdate(tx *gorm.DB) error {
	userID := GetCurrentUserID(tx)
	base.UpdatedBy = &userID
	return nil
}

func GetCurrentUserID(tx *gorm.DB) uint64 {
	// Lấy profileId từ context
	ginProfileId := tx.Statement.Context.Value("profileId")
	profileId := tx.Statement.Context.Value(_enums.ProfileIDKey)
	if ginProfileId != nil {
		return ginProfileId.(uint64)
	}
	if profileId == nil {
		return 0
	}
	log.Println("CALLBACK", profileId, ginProfileId)
	return profileId.(uint64) // Nếu không có user_id, trả về 0 (hoặc panic nếu cần)
}
