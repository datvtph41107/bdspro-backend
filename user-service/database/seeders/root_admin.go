package seeders

import (
	"context"
	"fmt"

	_utils "common/utils"
	bootstrapstore "user/infra/postgres/bootstrapadmin"
	"user/internal/usecase/bootstrapadmin"

	"gorm.io/gorm"
)

const requiredRootConfirmation = "CREATE_ROOT_OPERATOR"

func seedRootAdmin(ctx context.Context, db *gorm.DB, options Options) (Result, error) {
	if options.RootConfirmation != requiredRootConfirmation {
		return Result{}, fmt.Errorf("refusing root-admin seed: set QHPRO_BOOTSTRAP_CONFIRM=%s", requiredRootConfirmation)
	}

	service := bootstrapadmin.NewService(bootstrapstore.NewStore(db), _utils.HashPassword)
	created, err := service.Bootstrap(ctx, bootstrapadmin.Input{
		Username: options.RootUsername,
		Password: options.RootPassword,
		FullName: options.RootFullName,
		Email:    options.RootEmail,
		Phone:    options.RootPhone,
	})
	if err != nil {
		return Result{}, err
	}
	return Result{
		Name:            RootAdmin,
		ProfileID:       created.ProfileID,
		AuthID:          created.AuthID,
		RoleID:          created.RoleID,
		PermissionCount: created.PermissionCount,
	}, nil
}
