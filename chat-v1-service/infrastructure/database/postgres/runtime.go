package postgres

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	postgresInstance, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: newGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	// ✅ Get underlying sql.DB for connection pool configuration
	sqlDB, err := postgresInstance.DB()
	if err != nil {
		return nil, err
	}

	// ✅ Configure connection pool
	sqlDB.SetMaxOpenConns(25)                  // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Maximum lifetime of connections
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time of connections

	// ✅ Health check with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return postgresInstance, nil
}

// ✅ Add health check function
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// ✅ Add connection pool stats
func GetConnectionStats(db *gorm.DB) (open, idle int) {
	sqlDB, err := db.DB()
	if err != nil {
		return 0, 0
	}

	return sqlDB.Stats().OpenConnections, sqlDB.Stats().Idle
}
