package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_db "common/db"
	"user/config"
	"user/database/seeders"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	seedFlag := flag.String("name", "", "named User seed: admin, acceptance-admin, acceptance-client, root-admin")
	flag.Parse()

	seedValue := strings.TrimSpace(*seedFlag)
	if seedValue == "" && flag.NArg() > 0 {
		seedValue = flag.Arg(0)
	}
	if seedValue == "" {
		seedValue = string(seeders.Admin)
	}
	name, err := seeders.ParseName(seedValue)
	if err != nil {
		return err
	}

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

	result, err := seeders.Run(ctx, database, name, seeders.Options{
		Environment: os.Getenv("QHPRO_ENVIRONMENT"),

		AcceptanceAdminPassword: os.Getenv("QHPRO_ACCEPTANCE_ADMIN_PASSWORD"),
		AcceptanceUsername: firstNonEmpty(
			os.Getenv("QHPRO_ACCEPTANCE_USERNAME"),
			"qhpro_acceptance_0000000000",
		),
		AcceptancePhone: firstNonEmpty(
			os.Getenv("QHPRO_ACCEPTANCE_PHONE"),
			"0390000000",
		),
		AcceptanceEmail: firstNonEmpty(
			os.Getenv("QHPRO_ACCEPTANCE_EMAIL"),
			"qhpro-0000000000@e.invalid",
		),
		AcceptanceUserPassword: os.Getenv("QHPRO_ACCEPTANCE_USER_PASSWORD"),

		RootConfirmation: os.Getenv("QHPRO_BOOTSTRAP_CONFIRM"),
		RootUsername:     os.Getenv("QHPRO_BOOTSTRAP_ADMIN_USERNAME"),
		RootPassword:     os.Getenv("QHPRO_BOOTSTRAP_ADMIN_PASSWORD"),
		RootFullName:     os.Getenv("QHPRO_BOOTSTRAP_ADMIN_FULL_NAME"),
		RootEmail:        os.Getenv("QHPRO_BOOTSTRAP_ADMIN_EMAIL"),
		RootPhone:        os.Getenv("QHPRO_BOOTSTRAP_ADMIN_PHONE"),
	})
	if err != nil {
		return fmt.Errorf("run User seed %q: %w", name, err)
	}

	log.Printf(
		"User seed complete: name=%s profile_id=%d auth_id=%d role_id=%d permissions=%d",
		result.Name, result.ProfileID, result.AuthID, result.RoleID, result.PermissionCount,
	)
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
