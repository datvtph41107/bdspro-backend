package models

import (
	"fmt"

	"gorm.io/gorm"
)

// AutoMigrate is File Service's complete GORM bootstrap registry. Production
// evolution remains owned by file-service/migrate.
func AutoMigrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(hashtext(?))",
			"qhpro:file:gorm-schema",
		).Error; err != nil {
			return fmt.Errorf("acquire File schema lock: %w", err)
		}
		return tx.AutoMigrate(&FileEntity{}, &AccessEntity{})
	})
}
