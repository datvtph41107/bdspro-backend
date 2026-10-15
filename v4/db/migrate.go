// Package db owns User Service's GORM bootstrap registry.
//
// Production schema evolution remains owned by user-service/migrate. The
// registry exists for the explicit QHPRO_USER_DB_SCHEMA_MODE=automigrate mode
// and contains only current, reachable GORM-owned entity.
package db

import (
	"gorm.io/gorm"
)

func AutoMigrate(database *gorm.DB) error {
	// if err := database.AutoMigrate(AutoMigrateModels()...); err != nil {
	// 	return err
	// }
	return ensureAuthStatusFromLegacy(database)
}

// ensureAuthStatusFromLegacy copies registration lock rows from user_status when
// that legacy table is still present (pre-000021 DBs under automigrate bootstrap).
func ensureAuthStatusFromLegacy(database *gorm.DB) error {
	if database == nil {
		return nil
	}
	var exists bool
	if err := database.Raw(`
SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = current_schema()
      AND table_name = 'user_status'
)`).Scan(&exists).Error; err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return database.Exec(`
INSERT INTO auth_status (auth_id, active, verified, locked_until, updated_at)
SELECT
    us.auth_id,
    us.active,
    us.verified,
    us.locked_until,
    COALESCE(us.updated_at, NOW())
FROM user_status us
ON CONFLICT (auth_id) DO NOTHING
`).Error
}
