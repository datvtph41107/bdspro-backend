package db

import (
	_db "common/db"
	models "user/internal/models"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// &models.BlockEntity{},
		// &models.ContactEntity{},
		// &models.FollowEntity{},
		// &models.FriendEntity{},
		// &models.GroupEntity{},
		// &models.UserProfileEntity{},
		// &models.ProfessionEntity{},
		// &models.CertificationEntity{},
		// &models.AttachEntity{},
		// &models.Report{},
		&models.UserProfileEntity{},
		&models.BookmarkUserEntity{},
		&models.AdminProfile{},
	// &models.ReportReason{},
	// &models.ReportProof{},
	// &models.Rate{},
	)
}
