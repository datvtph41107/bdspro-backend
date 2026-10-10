package main

import (
	"common/logging"
	logfabric "common/logging/fabric"
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	_db "common/db"
	_utils "common/utils"
	"userv4/config"
	bootstrapstore "userv4/infra/postgres/bootstrapadmin"
	"userv4/internal/usecase/bootstrapadmin"
)

const requiredConfirmation = "CREATE_ROOT_OPERATOR"

func main() {
	defer logging.MustStart("user-bootstrap-admin")()
	logfabric.Bridge()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Topology + admin credentials come from config/<profile>.yml only.
	databaseConfig, err := config.LoadDatabaseRuntime()
	if err != nil {
		return fmt.Errorf("load User database config: %w", err)
	}

	admin := config.Properties.Admin
	confirm := strings.TrimSpace(admin.BootstrapConfirm)
	if confirm != requiredConfirmation {
		return fmt.Errorf(
			"refusing bootstrap: set admin.bootstrap_confirm: %q in config/<profile>.yml",
			requiredConfirmation,
		)
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
		Username: admin.Username,
		Password: admin.Password,
		FullName: admin.FullName,
		Email:    admin.Email,
		Phone:    admin.Phone,
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
		"root operator created: profile_id=%d auth_id=%d role_id=%d permissions_granted=%d username=%s",
		result.ProfileID, result.AuthID, result.RoleID, result.PermissionCount, admin.Username,
	)
	return nil
}
