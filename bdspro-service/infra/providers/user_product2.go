package providers

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	user_repo "bdspro/internal/repo/user"
	_utils "common/utils"
	"context"
)

type UProductProvider struct {
	ProductRepo user_repo.UProductRepo
}

func NewUProductProvider(productRepo user_repo.UProductRepo) *UProductProvider {
	return &UProductProvider{
		ProductRepo: productRepo,
	}
}

func (s *UProductProvider) ExportExcel(ctx context.Context, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return s.ProductRepo.SearchProductWithOwner(&profileId, dto)
}
