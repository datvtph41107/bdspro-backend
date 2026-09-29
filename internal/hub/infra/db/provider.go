package db

import (
	common_db "common/db"
	"fmt"
	"hub/config"

	"gorm.io/gorm"
)

// NewDB opens the single Hub PostgreSQL pool from the process-owned runtime
// snapshot. The compatibility package-global handle is maintained until the
// remaining shared DB lifecycle debt is retired in S6.
func NewDB(runtime config.Runtime) (*gorm.DB, func(), error) {
	database, err := common_db.Open(common_db.DatabaseConfig{
		DSN: runtime.Database.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open hub postgres: %w", err)
	}

	common_db.DB = database
	cleanup := func() {
		if sqlDB, err := database.DB(); err == nil {
			_ = sqlDB.Close()
		}
		if common_db.DB == database {
			common_db.DB = nil
		}
	}
	return database, cleanup, nil
}
