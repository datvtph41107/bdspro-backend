package db

import (
	_db "common/db"
	"user/internal/domain"
	"user/internal/models"
)

func AutoMigrate() {
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
