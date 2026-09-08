package postgresusage

import "gorm.io/gorm"

// AutoMigrate exposes the private usage row only to TQD's schema bootstrap
// registry. Production indexes/constraints remain canonical SQL concerns.
func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(&model{})
}
