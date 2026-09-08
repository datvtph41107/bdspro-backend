package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DocumentRepository interface {
	CreateBatch(ctx context.Context, documents []entity.AttachDocument) ([]entity.AttachDocument, error)
}
