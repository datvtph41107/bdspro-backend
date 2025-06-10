package db

import (
	_db "common/db"
	"social/internal/domain"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		// &domain.NewsFeed{},
		// &domain.FriendTag{},
		// &domain.NewsFeedMedia{},
		// &domain.ReportReason{},
		// &domain.Report{},
		&domain.Comment{},
	// &domain.Like{},
	// &domain.NewsFeedShare{},
	)
}
