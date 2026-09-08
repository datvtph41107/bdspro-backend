package seeders

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Name string

const (
	Admin            Name = "admin"
	AcceptanceAdmin  Name = "acceptance-admin"
	AcceptanceClient Name = "acceptance-client"
	RootAdmin        Name = "root-admin"
)

var ErrUnsafeEnvironment = errors.New("seed is not allowed in this environment")

type Options struct {
	Environment string

	AcceptanceAdminPassword string
	AcceptanceUsername      string
	AcceptancePhone         string
	AcceptanceEmail         string
	AcceptanceUserPassword  string

	RootConfirmation string
	RootUsername     string
	RootPassword     string
	RootFullName     string
	RootEmail        string
	RootPhone        string
}

type Result struct {
	Name            Name
	ProfileID       uint64
	AuthID          uint64
	RoleID          uint64
	PermissionCount int64
}

func Names() []Name {
	return []Name{Admin, AcceptanceAdmin, AcceptanceClient, RootAdmin}
}

func ParseName(value string) (Name, error) {
	name := Name(strings.ToLower(strings.TrimSpace(value)))
	for _, candidate := range Names() {
		if name == candidate {
			return name, nil
		}
	}
	return "", fmt.Errorf("unknown User seed %q", value)
}

func Run(ctx context.Context, db *gorm.DB, name Name, options Options) (Result, error) {
	if db == nil {
		return Result{}, errors.New("User seed database is not configured")
	}

	switch name {
	case Admin:
		return seedDevelopmentAdmin(ctx, db, options)
	case AcceptanceAdmin:
		return seedAcceptanceAdmin(ctx, db, options)
	case AcceptanceClient:
		return seedAcceptanceClient(ctx, db, options)
	case RootAdmin:
		return seedRootAdmin(ctx, db, options)
	default:
		return Result{}, fmt.Errorf("unknown User seed %q", name)
	}
}

func requireNonProduction(environment string) error {
	normalized := strings.ToLower(strings.TrimSpace(environment))
	switch normalized {
	case "development", "dev", "local", "test", "testing", "acceptance":
		return nil
	case "":
		return fmt.Errorf("%w: QHPRO_ENVIRONMENT is required", ErrUnsafeEnvironment)
	default:
		return fmt.Errorf("%w: %s", ErrUnsafeEnvironment, normalized)
	}
}
