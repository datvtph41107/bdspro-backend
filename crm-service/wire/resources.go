package wire

import (
	_db "common/db"
	"fmt"

	"gorm.io/gorm"
)

// OpenDatabase gives Wire ownership of CRM's single PostgreSQL pool and its
// close order. Every repository in the graph receives this exact *gorm.DB.
func OpenDatabase() (*gorm.DB, func(), error) {
	database, err := _db.NewDB()
	if err != nil {
		return nil, nil, fmt.Errorf("open CRM database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("access CRM database pool: %w", err)
	}
	cleanup := func() {
		_ = sqlDB.Close()
	}
	return database, cleanup, nil
}
