package postgres

import (
	"fmt"

	_db "common/db"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// TqdDB is the process-owned PostgreSQL pool for modules/tqd.
// DSN: database.tqd.dsn, falling back to database.dsn when empty (shared-DB develop).
type TqdDB struct {
	DB *gorm.DB
}

// NewTqdDB opens the TQD database pool and applies pending tqd-service/migrate files (userv4/db).
func NewTqdDB() (*TqdDB, func(), error) {
	dsn := firstNonEmptyDSN(
		viper.GetString("database.tqd.dsn"),
		viper.GetString("database.dsn"),
	)
	if dsn == "" {
		return nil, nil, fmt.Errorf("database.tqd.dsn or database.dsn is required for modules/tqd")
	}
	database, err := _db.Open(_db.DatabaseConfig{DSN: dsn})
	if err != nil {
		return nil, nil, fmt.Errorf("open tqd database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("tqd database pool: %w", err)
	}
	if err := ensureTqdLocationCatalogSchema(database); err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("tqd location catalog schema: %w", err)
	}
	return &TqdDB{DB: database}, func() { _ = sqlDB.Close() }, nil
}
