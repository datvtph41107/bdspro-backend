package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

type PropertyUserPostgres struct {
	db *gorm.DB
}

func NewPropertyUserPostgres(db *gorm.DB) repo.PropertyUserRepo {
	return &PropertyUserPostgres{db: db}
}

func (r *PropertyUserPostgres) CreateOwner(ctx context.Context, propertyId uint64, profileID uint64) (*domain.PropertyUser, error) {
	originId := _utils.GetOriginIdFromContext(ctx)
	propertyOfUser := &domain.PropertyUser{
		PropertyLineageID: &propertyId,
		OwnerOriginID:     profileID,
		OwnerAt:           time.Now(),

		OriginProfileID: &originId,
	}
	err := GetDB(ctx, r.db).Create(propertyOfUser).Error
	if err != nil {
		return nil, err
	}
	return propertyOfUser, nil
}
