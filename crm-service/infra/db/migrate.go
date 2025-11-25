package db

import (
	_db "common/db"
	"crm/internal/domain"
	"log"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// Core CRM entities
		&domain.ContactEntity{},
		&domain.LeadEntity{},
		&domain.LeadContactEntity{},

		// Pipeline management
		&domain.PipelineEntity{},
		&domain.StageEntity{},
		&domain.RuleEntity{},

		// Social entities
		&domain.FriendEntity{},
		&domain.FollowEntity{},
		&domain.BlockEntity{},
		&domain.GroupEntity{},

		// Sharing and documents
		&domain.SharingEntity{},
		&domain.DocumentEntity{},

		// Product care
		&domain.ProductCareEntity{},

		// App invited
		&domain.AppInvitedEntity{},

		// Registry
		&domain.FeedbackRegistry{},
		&domain.Rate{},
		&domain.Report{},
		&domain.ReportReason{},
		&domain.ReportProof{},
		&domain.FeedbackAttachEntity{},
	)

	log.Println("Migrate success")
}
