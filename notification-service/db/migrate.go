package db

import (
	_db "common/db"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
	// &domain.HistoryEntity{},
	// &domain.NotificationEntity{},
	// &domain.AdminHistoryEntity{},
	// &domain.AccountWarningEntity{},
	// &domain.AccountWarningLogEntity{},
	// &domain.AccountWarningTemplateEntity{},
	// &domain.DealHistoryEntity{},
	// &domain.HistoryAuthEntity{},
	// &domain.HistoryFcmEntity{},
	// &domain.PersonConfigEntity{},
	// &domain.PropertyHistory{},
	)
}
