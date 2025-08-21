package db

import (
	"auth/internal/domain"
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&domain.AuthMethodDomain{},
		&domain.UserOTPEntity{},
		&domain.UserSessionEntity{},
		&domain.UserStatusEntity{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RoleGroup{},
	)
}
