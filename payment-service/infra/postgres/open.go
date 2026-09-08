package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	paymentconfig "payment/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg paymentconfig.DatabaseConfig) (*gorm.DB, func(), error) {
	gormConfig := &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)}
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
