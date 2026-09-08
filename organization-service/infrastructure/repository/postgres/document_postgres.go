package postgres

import (
	"context"
	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.DocumentRepository
type DocumentPostgresRepository struct {
	db *gorm.DB
}

func NewDocumentPostgresRepository(db *gorm.DB) *DocumentPostgresRepository {
	return &DocumentPostgresRepository{db: db}
}

func (r *DocumentPostgresRepository) CreateBatch(ctx context.Context, documents []entity.AttachDocument) ([]entity.AttachDocument, error) {
	if err := GetDB(ctx, r.db).Create(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}
