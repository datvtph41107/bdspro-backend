package repo

import (
	"context"
	"hub/internal/domain"
)

type IProvinceV2Repo interface {
	// GetAll returns all provinces
	GetAll(ctx context.Context) ([]domain.ProvinceV2, error)

	// GetList returns paginated list of provinces
	GetList(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error)

	// GetByID returns province by ID
	GetByID(ctx context.Context, id uint64) (*domain.ProvinceV2, error)

	// GetByCode returns province by code
	GetByCode(ctx context.Context, code int) (*domain.ProvinceV2, error)

	// GetByCodeWithWards returns province with its wards
	GetByCodeWithWards(ctx context.Context, code int) (*domain.ProvinceV2, error)

	// Search searches provinces by keyword
	Search(ctx context.Context, keyword string, page, size int) ([]domain.ProvinceV2, int64, error)

	InferFromText(ctx context.Context, text string) (*domain.ProvinceV2, error)
}
