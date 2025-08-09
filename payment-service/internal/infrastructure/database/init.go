package database

import (
	"payment/internal/domain/entity"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Bank{})
}
