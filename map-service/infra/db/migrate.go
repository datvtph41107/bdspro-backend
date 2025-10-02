package db

import (
	"log"
	"map/internal/domain"

	"gorm.io/gorm"
)

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	db.AutoMigrate(
		&domain.MapPoint{},
	)

	log.Println("Database migrations completed successfully")
	return nil
}
