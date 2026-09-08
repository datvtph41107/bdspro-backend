package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DocumentRepository interface {
	CreateBatch(ctx context.Context, documents []domain.AttachDocument) ([]domain.AttachDocument, error)
}
