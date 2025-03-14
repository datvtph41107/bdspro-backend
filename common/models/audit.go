package cmodels

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditBase struct {
	CreatedBy *uint `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"column:updated_by" json:"updated_by,omitempty"`
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

func GetCurrentUserID(tx *gorm.DB) uint {
	ginCtx, ok := tx.Statement.Context.(*gin.Context)
	if !ok {
		return 0
	}

	// Lấy profileId từ context
	profileId, ok := ginCtx.Get("profileId")

	// Giả sử bạn lưu user_id vào context
	//  cjwt.GetProfileId(tx)
	// profileId, ok := tx.Statement.Context.Value("profileId").(uint)
	// val, ok := tx.Statement.Context.Value("claims").(*cjwt.Claims)
	println(profileId)
	if ok {
		return profileId.(uint)
	}
	return 0 // Nếu không có user_id, trả về 0 (hoặc panic nếu cần)
}
