package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"
	"errors"

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

func (r *PostgrePropertyIdentifyRepo) GetByID(ctx context.Context, id uint64) (*domain.PropertyIdentify, error) {
	var e domain.PropertyIdentify
	err := GetDB(ctx, r.DB).First(&e, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &e, err
}

func (r *PostgrePropertyIdentifyRepo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).Model(&domain.PropertyIdentify{}).Where("id = ?", id).Updates(fields).Error
}
