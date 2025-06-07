package db

import (
	_db "common/db"
	"crm/internal/domain"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// &FriendEntity{},
		// &ContactEntity{},
		// &FollowEntity{},
		// &BlockEntity{},
		// &GroupEntity{},
		// &domain.HistoryEntity{},
		// &domain.PipelineEntity{},
		// &domain.StageEntity{},
		// &domain.RuleEntity{},
		&domain.CustomerEntity{},
		// &domain.ContactEntity{},
	)
}
