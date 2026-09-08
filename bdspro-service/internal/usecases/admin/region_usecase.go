package admin_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type RegionUsecase struct {
	RegionRepo repo.RegionRepo
}

func NewRegionUsecase(regionRepo repo.RegionRepo) *RegionUsecase {
	return &RegionUsecase{
		RegionRepo: regionRepo,
	}
}

// GetRegionList lấy danh sách địa chỉ theo level và parent
func (u *RegionUsecase) GetRegionList(c context.Context, request *dto.RegionRequest) ([]domain.Region, int64, error) {
	return u.RegionRepo.Search(c, request)
}

// GetRegionDetail lấy chi tiết địa chỉ
func (u *RegionUsecase) GetRegionDetail(c context.Context, id uint64) (*domain.Region, error) {
	return u.RegionRepo.GetByID(c, id)
}
