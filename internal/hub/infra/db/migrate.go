package db

import (
	"fmt"
	"hub/internal/domain"

	"gorm.io/gorm"
)

// AutoMigrateModels là registry bootstrap GORM duy nhất của Hub. Registry này
// cũng được contract test đối chiếu với migrate/ để tránh thêm model nhưng quên
// cập nhật SQL production.
func AutoMigrateModels() []any {
	return []any{
		&domain.Applink{},
		&domain.EventQueueEntity{},
		&domain.Province{},
		&domain.District{},
		&domain.Ward{},
		&domain.ProvinceV2{},
		&domain.WardV2{},
		&domain.UserGuideEntity{},
		&domain.SystemConfigEntity{},
		&domain.FAQEntity{},
		&domain.VersionEntity{},
		&domain.ApiKeyEntity{},
		&domain.UserGuideStepEntity{},
		&domain.InteractiveEvent{},
		&domain.ErrorLog{},
		&domain.UpdateData{},
	}
}

// AutoMigrate chỉ chạy khi process owner chọn schema mode "automigrate".
// Chế độ mặc định "sql" để migrate/ là authority schema khi deploy.
func AutoMigrate(database *gorm.DB) error {
	if err := database.AutoMigrate(AutoMigrateModels()...); err != nil {
		return fmt.Errorf("auto migrate hub schema: %w", err)
	}
	return nil
}
