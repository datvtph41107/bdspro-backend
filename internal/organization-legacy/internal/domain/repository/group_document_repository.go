package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupDocumentRepository interface {
	Create(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	Update(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.GroupDocument, error)
	GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupDocument, uint32, error)
}
