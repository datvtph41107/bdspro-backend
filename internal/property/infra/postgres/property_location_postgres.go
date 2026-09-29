package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyLocationRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyLocationRepo(db *gorm.DB) repo.PropertyLocationRepo {
	return &PostgrePropertyLocationRepo{DB: db}
}

func (r *PostgrePropertyLocationRepo) Create(ctx context.Context, entity *domain.PropertyLocation) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}
