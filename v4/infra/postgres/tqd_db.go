package postgres

import (
	"fmt"
	"strings"

	_db "common/db"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// TqdDB is the process-owned PostgreSQL pool for modules/tqd (database.tqd.dsn).
// It is distinct from the v4 user DSN so TQD province_v2 / ward_v2 do not collide
// with hub's province_v2 on db_v4.
type TqdDB struct {
	DB *gorm.DB
}

// NewTqdDB opens the TQD database pool from database.tqd.dsn.
func NewTqdDB() (*TqdDB, func(), error) {
	dsn := strings.TrimSpace(viper.GetString("database.tqd.dsn"))
	if dsn == "" {
		return nil, nil, fmt.Errorf("database.tqd.dsn is required for modules/tqd")
	}
	database, err := _db.Open(_db.DatabaseConfig{DSN: dsn})
	if err != nil {
		return nil, nil, fmt.Errorf("open tqd database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("tqd database pool: %w", err)
	}
	return &TqdDB{DB: database}, func() { _ = sqlDB.Close() }, nil
}
