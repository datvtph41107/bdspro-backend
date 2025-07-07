package orm

import "gorm.io/gorm"

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *gormRepository {
	return &gormRepository{db: db}
}
