package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_db "common/db"
	_utils "common/utils"
	"user/config"
	bootstrapstore "user/infra/postgres/bootstrapadmin"
	"user/internal/usecase/bootstrapadmin"
)

const requiredConfirmation = "CREATE_ROOT_OPERATOR"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if os.Getenv("QHPRO_BOOTSTRAP_CONFIRM") != requiredConfirmation {
		return fmt.Errorf("refusing bootstrap: set QHPRO_BOOTSTRAP_CONFIRM=%s", requiredConfirmation)
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

	service := bootstrapadmin.NewService(bootstrapstore.NewStore(database), _utils.HashPassword)
	result, err := service.Bootstrap(ctx, bootstrapadmin.Input{
		Username: os.Getenv("QHPRO_BOOTSTRAP_ADMIN_USERNAME"),
		Password: os.Getenv("QHPRO_BOOTSTRAP_ADMIN_PASSWORD"),
		FullName: os.Getenv("QHPRO_BOOTSTRAP_ADMIN_FULL_NAME"),
		Email:    os.Getenv("QHPRO_BOOTSTRAP_ADMIN_EMAIL"),
		Phone:    os.Getenv("QHPRO_BOOTSTRAP_ADMIN_PHONE"),
	})
	if err != nil {
		switch {
		case errors.Is(err, bootstrapadmin.ErrAlreadyBootstrapped):
			return errors.New("root operator already exists; use the authenticated Admin/IAM API for every later administrator")
		case errors.Is(err, bootstrapadmin.ErrUsernameInUse):
			return errors.New("requested administrator username is already in use")
		default:
			return err
		}
	}

	log.Printf(
		"root operator created: profile_id=%d auth_id=%d role_id=%d permissions_granted=%d",
		result.ProfileID, result.AuthID, result.RoleID, result.PermissionCount,
	)
	return nil
}
