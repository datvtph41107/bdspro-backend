package postgresjob

import "gorm.io/gorm"

// AutoMigrate exposes this package's private persistence model to the single
// TQD schema bootstrap registry without leaking it into business packages.
func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(&model{})
}
