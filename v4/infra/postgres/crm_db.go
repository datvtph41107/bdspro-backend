package postgres

import (
	"fmt"
	"strings"

	_db "common/db"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// CrmDB is the process-owned PostgreSQL pool for modules/crm (database.crm.dsn).
// It is distinct from the v4 user DSN so CRM-owned tables (tb_contact, friend, …)
// stay on crm-service's db_crm without colliding with the thinner v4 stub contact table.
type CrmDB struct {
	DB *gorm.DB
}

// NewCrmDB opens the CRM database pool from database.crm.dsn.
func NewCrmDB() (*CrmDB, func(), error) {
	dsn := strings.TrimSpace(viper.GetString("database.crm.dsn"))
	if dsn == "" {
		return nil, nil, fmt.Errorf("database.crm.dsn is required for modules/crm")
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
