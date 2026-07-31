package db

import (
	_db "common/db"
	"crm/internal/domain"
	seo_domain "crm/internal/domain/seo"
	"log"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&domain.SupportTicket{},
		&domain.SupportTicketNote{},
		&domain.SupportTicketEvent{},
		&domain.AdminOpportunityEvent{},
		&seo_domain.SeoRef{},
	)

	log.Println("Migrate success")
}
