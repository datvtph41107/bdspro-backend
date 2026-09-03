package db

import (
	_db "common/db"
	"hub/internal/domain"
	"log"
	"strings"
)

func AutoMigrate() {
	db := _db.DB

	models := []interface{}{
		&domain.Applink{},
		&domain.EventQueueEntity{},
		&domain.Province{},
		&domain.District{},
		&domain.UserGuideEntity{},
		&domain.SystemConfigEntity{},
		&domain.FAQEntity{},
		&domain.VersionEntity{},
		&domain.ApiKeyEntity{},
		&domain.UserGuideStepEntity{},
		&domain.InteractiveEvent{},
		&domain.ClickURL{},
		&domain.ErrorLog{},
		&domain.UpdateData{},
		&domain.RequestStatsLog{},
		&domain.UserInteractionLog{},
		&domain.RequestAuditLog{},
		&domain.KeywordSearchLog{},
		&domain.SearchLatencyLog{},
		// &domain.UserGuideGroupEntity{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			if isIgnorableMigrateError(err) {
				log.Printf("Auto migrate skip (already exists): %v", err)
				continue
			}
			log.Fatalf("Auto migrate error: %v", err)
		}
	}

	log.Println("Auto migrate success")
}

func isIgnorableMigrateError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "42p07")
}
