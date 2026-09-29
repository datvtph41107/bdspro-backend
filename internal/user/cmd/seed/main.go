package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_db "common/db"
	"user/config"
	"user/database/seeders"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	databaseConfig, err := config.LoadDatabaseRuntime()
	if err != nil {
		return fmt.Errorf("load User database config: %w", err)
	}
	database, err := _db.Open(_db.DatabaseConfig{
		DSN:             databaseConfig.DSN,
		MaxOpenConns:    databaseConfig.MaxOpenConns,
		MaxIdleConns:    databaseConfig.MaxIdleConns,
		ConnMaxLifetime: databaseConfig.ConnMaxLifetime,
	})
	if err != nil {
		return fmt.Errorf("open User database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("get User database pool: %w", err)
	}
	defer sqlDB.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	result, err := seeders.RunDatabaseSeeder(ctx, database)
	if err != nil {
		return fmt.Errorf("run User DatabaseSeeder: %w", err)
	}

	fmt.Printf("User DatabaseSeeder complete: seeders=%d\n", result.SeederCount)
	return nil
}
