package db

import (
	"log"

	_db "common/db"
	"user/internal/domain"
	"user/internal/models"
)

func AutoMigrate() {
	// Bảng activity login/refresh cho dashboard active users.
	_ = _db.DB.AutoMigrate(&domain.UserActivityLog{})
	ensureUserProfileLastLoginColumns()
	return
	_db.DB.AutoMigrate(
		// Auth entities (migrated from auth-service)
		&domain.AuthMethod{},
		&domain.AuthUser{},
		&domain.UserOTPEntity{},
		&domain.UserSessionEntity{},
		&domain.DeviceEntity{},
		&domain.UserStatusEntity{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RoleGroup{},
		&domain.RoleProfile{},
		&domain.AdminAccessDomain{},
		&domain.AdminAccessLogDomain{},
		&domain.Color{},
		&domain.AuthConfig{},
		&domain.UserPINEntity{},

		// Core entities
		&models.UserProfileEntity{},
		&models.ProfileDeleted{},
		&models.AdminProfile{},

		// Social entities
		&models.BlockEntity{},
		&models.ContactEntity{},
		&models.FollowEntity{},
		&models.FriendEntity{},
		&models.GroupEntity{},
		&models.BookmarkUserEntity{},
		&models.MainAreaEntity{},
		&models.MainAreaProfileEntity{},
		&models.PurposeUseEntity{},
		&models.PurposeUseProfileEntity{},
		&models.TagEntity{},
		&models.TagUserEntity{},
		&models.KYCEntity{},
		&models.ProfileMediaEntity{},

		// Profile related entities
		&models.ProfessionEntity{},
		&models.CertificationEntity{},
		&models.AttachEntity{},

		// Membership & plan
		&models.Membership{},
		&models.Plan{},

		// Pack related entities
		&models.PriceTableDomain{},
		&models.PackItemDomain{},
	)
}

func ensureUserProfileLastLoginColumns() {
	stmts := []string{
		`ALTER TABLE user_profile ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ DEFAULT NULL`,
		`ALTER TABLE user_profile ADD COLUMN IF NOT EXISTS total_login INTEGER NOT NULL DEFAULT 0`,
		`CREATE INDEX IF NOT EXISTS idx_user_profile_last_login_at ON user_profile(last_login_at)`,
	}
	for _, stmt := range stmts {
		if err := _db.DB.Exec(stmt).Error; err != nil {
			log.Printf("⚠️  migrate user_profile last_login_at: %v", err)
		}
	}
}
