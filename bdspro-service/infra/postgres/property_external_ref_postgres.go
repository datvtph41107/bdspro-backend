package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyExternalRefRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyExternalRefRepo(db *gorm.DB) repo.PropertyExternalRefRepo {
	return &PostgrePropertyExternalRefRepo{DB: db}
}

func (r *PostgrePropertyExternalRefRepo) Create(ctx context.Context, entity *domain.PropertyExternalRef) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}
