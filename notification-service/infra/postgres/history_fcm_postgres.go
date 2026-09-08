package postgres

import "gorm.io/gorm"

type HistoryRepository struct {
	DB *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{DB: db}
}
