package postgre

import (
	"context"
	"crm/infra/impl"
	"crm/internal/domain"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.DocumentRepo
type DocumentPostgre struct {
	DB *gorm.DB
}

func NewDocumentPostgre(DB *gorm.DB) *DocumentPostgre {
	return &DocumentPostgre{DB: DB}
}

func (repo DocumentPostgre) Bulk(ctx context.Context, deletedIds []uint64, documents []domain.DocumentEntity) error {
	var err error
	if len(deletedIds) > 0 {
		err = impl.GetDB(ctx, repo.DB).
			Model(&domain.DocumentEntity{}).
			Where("id IN (?)", deletedIds).
			Update("deleted_at", time.Now()).Error
	}
	if len(documents) > 0 {
		err = impl.GetDB(ctx, repo.DB).
			Model(&domain.DocumentEntity{}).
			Save(&documents).Error
	}
	return err
}

func (repo DocumentPostgre) Create(ctx context.Context, documents []domain.DocumentEntity) error {
	return impl.GetDB(ctx, repo.DB).
		Model(&domain.DocumentEntity{}).
		Create(&documents).Error
}