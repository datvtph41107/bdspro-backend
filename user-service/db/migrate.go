package db

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
	// &auth_model.UserSessionEntity{},
	// &auth_model.UserAuthEntity{},
	// &user_model.BlockEntity{},
	// &user_model.ContactEntity{},
	// &user_model.FollowEntity{},
	// &user_model.FriendEntity{},
	// &user_model.GroupEntity{},
	// &user_model.UserProfileEntity{},
	// &user_model.ProfessionEntity{},
	// &user_model.CertificationEntity{},
	// &user_model.Rate{},
	)
}
