package db

import (
	_db "common/db"
	"social/internal/domain"
)

func MigrateDomain() {
	_db.DB.AutoMigrate(
		// Core news feed entities
		&domain.NewsFeed{},
		&domain.NewsFeedOfUser{},
		&domain.NewsFeedOfGroup{},

		// News feed interactions
		&domain.NewsFeedMedia{},
		&domain.NewsFeedShare{},
		&domain.Comment{},
		&domain.Like{},

		// Social features
		&domain.FriendTag{},

		// Report entities
		&domain.Report{},
		&domain.ReportReason{},
	)
}
