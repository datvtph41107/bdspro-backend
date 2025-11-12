package repo

import (
	"context"
	"hub/internal/domain"
)

type IWardRepo interface {
	GetAll(ctx context.Context) ([]domain.Ward, error)
	GetByDistrictID(ctx context.Context, districtID string) ([]domain.Ward, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.Ward, error)
	Search(ctx context.Context, keyword string) ([]domain.Ward, error)
	SearchByDistrictID(ctx context.Context, districtID string, keyword string) ([]domain.Ward, error)
	InferFromText(ctx context.Context, text string) (*domain.Ward, error)
}
