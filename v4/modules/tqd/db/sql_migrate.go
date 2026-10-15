// Package db is reserved for TQD DB helpers; versioned SQL remains in tqd-service/migrate.
package db

import "gorm.io/gorm"

// RunPendingMigrations is intentionally not invoked from v4 boot (see infra/postgres/module_db.go).
func RunPendingMigrations(_ *gorm.DB) error {
	return nil
}
