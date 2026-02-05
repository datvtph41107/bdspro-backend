package db

import (
	"core/kernel/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

func New(cfg config.DatabaseConfig, log logger.Interface) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger:      log,
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}
