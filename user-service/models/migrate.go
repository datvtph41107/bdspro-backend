package models

import _db "common/db"

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&UserProfileEntity{},
		&UserSessionEntity{},
	)
}
