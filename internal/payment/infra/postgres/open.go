package postgres

import (
	"fmt"

	paymentconfig "payment/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg paymentconfig.DatabaseConfig) (*gorm.DB, func(), error) {
	gormConfig := &gorm.Config{Logger: newGormLogger()}
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: cfg.DSN, PreferSimpleProtocol: true}), gormConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("open Payment PostgreSQL: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get Payment PostgreSQL pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	cleanup := func() { _ = sqlDB.Close() }
	return db, cleanup, nil
}
