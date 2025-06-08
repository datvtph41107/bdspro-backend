package db

import (
	_db "common/db"
	"social/internal/domain"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		&domain.NewsFeed{},
		&domain.TagFriend{},
		&domain.Media{},
		&domain.ReportReason{},
		&domain.Report{},
		&domain.Comment{},
		&domain.Like{},
	)
}
