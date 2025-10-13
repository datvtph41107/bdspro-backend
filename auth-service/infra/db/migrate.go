package db

import (
	"auth/internal/domain"
	_db "common/db"
	"log"
	"os"

	"gorm.io/gorm"
)

func RunSQL(db *gorm.DB, path string) {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read sql file %s: %v", path, err)
	}

	// Chạy toàn bộ SQL trong file
	if err := db.Exec(string(sqlBytes)).Error; err != nil {
		log.Fatalf("cannot execute sql file %s: %v", path, err)
	}
}

func AutoMigrate() {
	// RunSQL(_db.DB, "infra/db/migrate.sql")
	// RunSQL(_db.DB, "infra/db/migrate_user_info.sql")
	_db.DB.AutoMigrate(
		&domain.AuthMethod{},
		&domain.AuthUser{},
		&domain.UserOTPEntity{},
		&domain.UserSessionEntity{},
		&domain.UserStatusEntity{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RoleGroup{},
		&domain.RoleProfile{},
		&domain.AdminAccessDomain{},
		&domain.Color{},
	)
}
