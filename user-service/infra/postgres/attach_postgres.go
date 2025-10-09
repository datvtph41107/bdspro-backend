package postgres

import (
	"context"
	models "user/internal/models"

	"gorm.io/gorm"
)

type AttachPostgres struct {
	db *gorm.DB
}

func NewAttachPostgres(db *gorm.DB) *AttachPostgres {
	return &AttachPostgres{db: db}
}

func (a *AttachPostgres) CreateBatchAttach(ctx context.Context, attach *[]models.AttachEntity) error {
	if err := GetDB(ctx, a.db).
		Save(attach).
		Error; err != nil {
		return err
	}
	return nil
}
