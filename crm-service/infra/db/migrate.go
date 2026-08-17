package db

import (
	_db "common/db"
	seo_domain "crm/internal/domain/seo"
	"log"
)

func AutoMigrate() {
	_db.DB.AutoMigrate(
		&seo_domain.SeoRef{},
	)

	log.Println("Migrate success")
}
