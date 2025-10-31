package db

import (
	_db "common/db"
	models "user/internal/models"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// Core entities
		&models.UserProfileEntity{},
		&models.ProfileDeleted{},
		&models.AdminProfile{},
		&models.SystemConfigEntity{},

		// Social entities
		&models.BlockEntity{},
		&models.ContactEntity{},
		&models.FollowEntity{},
		&models.FriendEntity{},
		&models.GroupEntity{},
		&models.BookmarkUserEntity{},

		// Profile related entities
		&models.ProfessionEntity{},
		&models.CertificationEntity{},
		&models.AttachEntity{},
	)
}
