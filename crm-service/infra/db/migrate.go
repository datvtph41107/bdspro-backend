package db

import (
	_db "common/db"
	seo_domain "crm/internal/domain/seo"
	"log"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		// Core CRM entities
		// &domain.ContactEntity{},
		// &domain.LeadEntity{},
		// &domain.LeadContactEntity{},
		// &domain.OriginProfile{},
		// // Pipeline management
		// &domain.PipelineEntity{},
		// &domain.StageEntity{},
		// &domain.RuleEntity{},

		// // Social entities
		// &domain.FriendEntity{},
		// &domain.FollowEntity{},
		// &domain.BlockEntity{},
		// &domain.GroupEntity{},

		// // Sharing and documents
		// &domain.SharingEntity{},
		// &domain.DocumentEntity{},

		// // Product care
		// &domain.ProductCareEntity{},

		// // App invited
		// &domain.AppInvitedEntity{},

		// // Registry
		// &domain.FeedbackRegistry{},
		// &domain.Rate{},
		// &domain.Report{},
		// &domain.ReportReason{},
		// &domain.ReportProof{},
		// &domain.FeedbackAttachEntity{},

		// Marketing entities
		// &domain.AdvertisingSettings{},
		// &domain.BudgetConfig{},
		// &domain.Campaign{},
		// &domain.Package{},
		// &domain.BudgetUsage{},
		// &domain.BudgetAlert{},

		// &domain.Appointment{},
		// &domain.AppointmentParticipant{},
		// &domain.AppointmentReminder{},
		// &domain.ContactProductEntity{},
		&seo_domain.SeoDomain{},
		&seo_domain.SeoInternalLink{},
		&seo_domain.SeoRelative{},
		&seo_domain.SeoTrend{},
	)

	log.Println("Migrate success")
}