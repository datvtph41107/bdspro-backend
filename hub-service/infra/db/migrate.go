package db

import (
	_db "common/db"
	"hub/internal/domain"
	"log"
)

func AutoMigrate() {
	db := _db.DB

	err := db.AutoMigrate(
		&domain.EventQueueEntity{},
		&domain.Province{},
		&domain.District{},
		&domain.UserGuideEntity{},
		&domain.SystemConfigEntity{},
		&domain.FAQEntity{},
		// &domain.UserGuideGroupEntity{},
	)

	if err != nil {
		log.Fatalf("Auto migrate error: %v", err)
	}

	log.Println("Auto migrate success")
}
