package db

import (
	_db "common/db"
	"fmt"
	"log/slog"
	"strings"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
	// &domain.SupportTicket{},
	// &domain.SupportTicketNote{},
	// &domain.SupportTicketEvent{},
	// &domain.AdminOpportunityEvent{},
	// &seo_domain.SeoRef{},
	)
	slog.Info(strings.TrimSuffix(fmt.Sprintln("Migrate success"), "\n"))
}
