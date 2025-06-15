package providers

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
)

type UProductProvider struct {
	ProductRepo repo.ProductRepo
}

func NewUProductProvider(productRepo repo.ProductRepo) *UProductProvider {
	return &UProductProvider{
		ProductRepo: productRepo,
	}
}

func (s *UProductProvider) ExportExcel(ctx context.Context, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return s.ProductRepo.SearchProductWithOwner(&profileId, dto)
}
