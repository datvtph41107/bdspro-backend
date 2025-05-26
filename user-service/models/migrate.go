package models

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
	// &user_model.UserProfileEntity{},
	// &auth_model.UserSessionEntity{},
	// &auth_model.UserAuthEntity{},
	// &user_model.BlockEntity{},
	// &user_model.ContactEntity{},
	// &user_model.FollowEntity{},
	// &user_model.FriendEntity{},
	// &user_model.GroupEntity{},
	)
}
