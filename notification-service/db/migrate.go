package db

import (
	_db "common/db"
	"noti/internal/domain"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&domain.HistoryEntity{},
		&domain.NotificationEntity{},
	)
}
