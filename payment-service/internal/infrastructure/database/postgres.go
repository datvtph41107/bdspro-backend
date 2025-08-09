package database

import (
	"log"
	"os"
	"time"

	"payment/appconfig"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB() *gorm.DB {
	cfg := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  appconfig.DB_URL,
		PreferSimpleProtocol: true,
	}), cfg)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	return db
}
