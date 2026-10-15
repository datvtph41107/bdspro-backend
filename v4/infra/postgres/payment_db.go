package postgres

import (
	"fmt"

	_db "common/db"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// PaymentDB is the PostgreSQL pool for modules/payment (banks catalog).
// DSN: database.payment.dsn, falling back to database.dsn when empty.
type PaymentDB struct {
	DB *gorm.DB
}

func NewPaymentDB() (*PaymentDB, func(), error) {
	dsn := firstNonEmptyDSN(
		viper.GetString("database.payment.dsn"),
		viper.GetString("database.dsn"),
	)
	if dsn == "" {
		return nil, nil, fmt.Errorf("database.payment.dsn or database.dsn is required for modules/payment")
	}
	database, err := _db.Open(_db.DatabaseConfig{DSN: dsn})
	if err != nil {
		return nil, nil, fmt.Errorf("open payment database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("payment database pool: %w", err)
	}
	return &PaymentDB{DB: database}, func() { _ = sqlDB.Close() }, nil
}
