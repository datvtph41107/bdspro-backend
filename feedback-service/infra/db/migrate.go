package db

import (
	_db "common/db"
	"feedback/internal/domain"
)

func AutoMigrate() {
	db := _db.DB

	// Migrate rate related tables
	db.AutoMigrate(&domain.Rate{})
	db.AutoMigrate(&domain.AttachEntity{})

	// Migrate report related tables
	db.AutoMigrate(&domain.Report{})
	db.AutoMigrate(&domain.ReportReason{})
	db.AutoMigrate(&domain.ReportProof{})
	// db.AutoMigrate(&domain.ReportEntity{})
	db.AutoMigrate(&domain.FeedbackRegistry{})
}
