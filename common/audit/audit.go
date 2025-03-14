package models

import (
	"gorm.io/gorm"
)

// AutoAudit gắn CreatedBy và UpdatedBy khi insert/update
func AutoAudit(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("set_created_by", func(tx *gorm.DB) {
		profileId := tx.Statement.Context.Value("profileId").(uint)
		if tx.Statement.Schema != nil {
			tx.Statement.SetColumn("created_by", profileId) // Thay bằng ID người dùng thực
		}
	})

	db.Callback().Update().Before("gorm:update").Register("set_updated_by", func(tx *gorm.DB) {
		profileId := tx.Statement.Context.Value("profileId").(uint)
		if tx.Statement.Schema != nil {
			tx.Statement.SetColumn("updated_by", profileId) // Thay bằng ID người dùng thực
		}
	})
}
