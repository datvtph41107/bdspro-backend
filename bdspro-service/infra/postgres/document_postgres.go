package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.DocumentRepository
type DocumentPostgresRepository struct {
	db *gorm.DB
}

func NewDocumentPostgresRepository(db *gorm.DB) *DocumentPostgresRepository {
	return &DocumentPostgresRepository{db: db}
}

func (r *DocumentPostgresRepository) CreateBatch(ctx context.Context, documents []domain.AttachDocument) ([]domain.AttachDocument, error) {
	if err := GetDB(ctx, r.db).Create(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}
