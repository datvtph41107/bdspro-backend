package _models

import (
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
	// ginCtx, ok := tx.Statement.Context.(*gin.Context)
	// if !ok {
	// 	return 0
	// }

	// Lấy profileId từ context
	profileId := tx.Statement.Context.Value("profileId")
	if profileId == nil {
		return 0
	}

	// Giả sử bạn lưu user_id vào context
	//  cjwt.GetProfileId(tx)
	// profileId, ok := tx.Statement.Context.Value("profileId").(uint)
	// val, ok := tx.Statement.Context.Value("claims").(*cjwt.Claims)
	log.Println("CALLBACK", profileId)
	// if ok {
	// 	return profileId.(uint64)
	// }
	return profileId.(uint64) // Nếu không có user_id, trả về 0 (hoặc panic nếu cần)
}
