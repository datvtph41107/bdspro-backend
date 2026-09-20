package _db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"common/logging"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultMaxOpenConns    = 300
	defaultMaxIdleConns    = 200
	defaultConnMaxLifetime = 5 * time.Minute
)

// DatabaseConfig is the process-owned input required to open PostgreSQL.
// Service-specific schema, queries, and transaction semantics do not belong here.
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

var DB *gorm.DB

// Open establishes one PostgreSQL pool from explicit process-owned config.
// It does not mutate DB; compatibility callers decide whether to expose the
// returned pool through the historical package-global handle.
func Open(config DatabaseConfig) (*gorm.DB, error) {
	dsn := strings.TrimSpace(config.DSN)
	if dsn == "" {
		return nil, errors.New("database DSN is required")
	}

	gormLogger := logger.New(
		logging.StdLogger("postgres"),
		logger.Config{
			SlowThreshold:        time.Second,
			LogLevel:             logger.Info,
			Colorful:             false,
			ParameterizedQueries: true,
		},
	)

	database, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{Logger: gormLogger},
	)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("get postgres pool: %w", err)
	}

	maxOpen := config.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = defaultMaxOpenConns
	}
	maxIdle := config.MaxIdleConns
	if maxIdle < 0 {
		maxIdle = 0
	}
	if maxIdle == 0 {
		maxIdle = defaultMaxIdleConns
	}
	if maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	lifetime := config.ConnMaxLifetime
	if lifetime <= 0 {
		lifetime = defaultConnMaxLifetime
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(lifetime)

	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return database, nil
}

// NewDB is the compatibility constructor for services that have not moved DB
// configuration into their process boundary.
func NewDB() (*gorm.DB, error) {
	database, err := Open(DatabaseConfig{DSN: viper.GetString("database.dsn")})
	if err != nil {
		return nil, err
	}
	DB = database
	return database, nil
}
