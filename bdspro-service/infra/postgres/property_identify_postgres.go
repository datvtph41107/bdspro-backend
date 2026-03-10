package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyIdentifyRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyIdentifyRepo(db *gorm.DB) repo.PropertyIdentifyRepo {
	return &PostgrePropertyIdentifyRepo{DB: db}
}

func (r *PostgrePropertyIdentifyRepo) Create(ctx context.Context, entity *domain.PropertyIdentify) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}
