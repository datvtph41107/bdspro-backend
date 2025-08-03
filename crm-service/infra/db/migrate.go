package db

import (
	_db "common/db"
	"crm/internal/domain"
	"log"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// &FriendEntity{},
		// &domain.PipelineEntity{},
		&domain.ContactEntity{},
		// &FollowEntity{},
		// &BlockEntity{},
		// &GroupEntity{},
		// &domain.HistoryEntity{},
		// &domain.AppInvitedEntity{},
		// &domain.PipelineEntity{},
		// &domain.StageEntity{},
		// &domain.RuleEntity{},
		&domain.LeadEntity{},
		// &domain.ContactEntity{},
		// &domain.SharingEntity{},
		// &domain.LeadEntity{},
		// &domain.ProductCareEntity{},
		// &domain.DocumentEntity{},
		// &domain.ProductCareEntity{},
	)

	log.Println("Migrate success")
}
