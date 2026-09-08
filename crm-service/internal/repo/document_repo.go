package repo

import (
	"context"
	"crm/internal/domain"
)

type DocumentRepo interface {
	Bulk(ctx context.Context, deletedIds []uint64, documents []domain.DocumentEntity) error
	Create(ctx context.Context, documents []domain.DocumentEntity) error
}