package postgres

import (
	"context"
	user_model "user/internal/models/user"

	"gorm.io/gorm"
)

type AttachPostgres struct {
	db *gorm.DB
}

func NewAttachPostgres(db *gorm.DB) *AttachPostgres {
	return &AttachPostgres{db: db}
}

func (a *AttachPostgres) CreateBatchAttach(ctx context.Context, attach *[]user_model.AttachEntity) error {
	if err := GetDB(ctx, a.db).
		Save(attach).
		Error; err != nil {
		return err
	}
	return nil
}
