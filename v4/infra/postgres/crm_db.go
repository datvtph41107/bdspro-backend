package postgres

import (
	"fmt"

	_db "common/db"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// CrmDB is the process-owned PostgreSQL pool for modules/crm.
// DSN: database.crm.dsn, falling back to database.dsn when empty (shared-DB develop).
type CrmDB struct {
	DB *gorm.DB
}

// NewCrmDB opens the CRM database pool from database.crm.dsn (or database.dsn).
func NewCrmDB() (*CrmDB, func(), error) {
	dsn := firstNonEmptyDSN(
		viper.GetString("database.crm.dsn"),
		viper.GetString("database.dsn"),
	)
	if dsn == "" {
		return nil, nil, fmt.Errorf("database.crm.dsn or database.dsn is required for modules/crm")
	}
	database, err := _db.Open(_db.DatabaseConfig{DSN: dsn})
	if err != nil {
		return nil, nil, fmt.Errorf("open crm database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("crm database pool: %w", err)
	}
	return &CrmDB{DB: database}, func() { _ = sqlDB.Close() }, nil
}
