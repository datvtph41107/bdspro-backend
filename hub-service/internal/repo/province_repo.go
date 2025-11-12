package repo

import (
	"context"
	"hub/internal/domain"
)

type IProvinceRepo interface {
	GetAll(ctx context.Context) ([]domain.Province, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.Province, error)
	Search(ctx context.Context, keyword string) ([]domain.Province, error)
	InferFromText(ctx context.Context, text string) (*domain.Province, error)
}
