package seeders

import (
	"context"
	"fmt"

	_utils "common/utils"

	"gorm.io/gorm"
)

func seedDevelopmentAdmin(ctx context.Context, db *gorm.DB, options Options) (Result, error) {
	if err := requireNonProduction(options.Environment); err != nil {
		return Result{}, err
	}

	passwordHash, err := _utils.HashPassword("admin123")
	if err != nil {
		return Result{}, fmt.Errorf("hash development admin password: %w", err)
	}

	result, err := seedAdmin(ctx, db, adminSeedSpec{
		Username:             "admin",
		PasswordHash:         passwordHash,
		FullName:             "Development Admin",
		RoleKey:              "QHPRO_DEVELOPMENT_ADMIN",
		RoleName:             "QHPRO Development Admin",
		RoleDescription:      "Development-only administrator provisioned by the User seed lifecycle",
		AdminRoleDescription: "Development-only administrator",
		RoleType:             "ROLE_ADMIN",
		RoleCode:             20,
		GrantAllPermissions:  true,
	})
	if err != nil {
		return Result{}, err
	}
	result.Name = Admin
	return result, nil
}
