package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyInfoRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyInfoRepo(db *gorm.DB) repo.PropertyInfoRepo {
	return &PostgrePropertyInfoRepo{DB: db}
}

func (r *PostgrePropertyInfoRepo) Save(ctx context.Context, entity *domain.PropertyInfo) error {
	return GetDB(ctx, r.DB).Save(entity).Error
}
