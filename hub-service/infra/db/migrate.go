package db

import (
	_db "common/db"
	"fmt"
	"hub/internal/domain"
	"log"
)

func AutoMigrate() error {
	db := _db.DB

	err := db.AutoMigrate(
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
		&domain.ErrorLog{},
		&domain.UpdateData{},
		// &domain.UserGuideGroupEntity{},
	)

	if err != nil {
		return fmt.Errorf("auto migrate hub schema: %w", err)
	}

	log.Println("Auto migrate success")
	return nil
}
