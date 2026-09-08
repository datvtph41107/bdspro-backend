package seeders

import (
	"context"
	"fmt"
	"strings"

	_utils "common/utils"

	"gorm.io/gorm"
)

func seedAcceptanceAdmin(ctx context.Context, db *gorm.DB, options Options) (Result, error) {
	if err := requireNonProduction(options.Environment); err != nil {
		return Result{}, err
	}
	password := strings.TrimSpace(options.AcceptanceAdminPassword)
	if len(password) < 12 {
		return Result{}, fmt.Errorf("QHPRO_ACCEPTANCE_ADMIN_PASSWORD must contain at least 12 characters")
	}
	passwordHash, err := _utils.HashPassword(password)
	if err != nil {
		return Result{}, fmt.Errorf("hash acceptance admin password: %w", err)
	}

	result, err := seedAdmin(ctx, db, adminSeedSpec{
		Username:             "qhpro.acceptance.admin",
		PasswordHash:         passwordHash,
		FullName:             "QHPRO Acceptance Commercial Admin",
		Email:                "qhpro.acceptance.admin@example.invalid",
		Phone:                "0399999999",
		RoleKey:              "QHPRO_ACCEPTANCE_COMMERCIAL_ADMIN",
		RoleName:             "QHPRO Acceptance Commercial Admin",
		RoleDescription:      "Acceptance fixture role for commercial operation verification",
		AdminRoleDescription: "Acceptance fixture for commercial UI verification",
		RoleType:             "ROLE_ADMIN",
		RoleCode:             20,
		PermissionKeys:       acceptanceAdminPermissionKeys,
	})
	if err != nil {
		return Result{}, err
	}
	result.Name = AcceptanceAdmin
	return result, nil
}
