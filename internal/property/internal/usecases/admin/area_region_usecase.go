package admin_usecases

import (
	"context"
	"bdspro/internal/domain"
	_repo "bdspro/internal/repo"
)

type AreaRegionUsecase struct {
	repo _repo.AreaRegionRepo
}

func NewAreaRegionUsecase(repo _repo.AreaRegionRepo) *AreaRegionUsecase {
	return &AreaRegionUsecase{repo: repo}
}

func (u *AreaRegionUsecase) GetByID(ctx context.Context, id uint64) (*domain.AreaRegion, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *AreaRegionUsecase) GetByIDs(ctx context.Context, ids []uint64) ([]domain.AreaRegion, error) {
	return u.repo.GetByIDs(ctx, ids)
}

func (u *AreaRegionUsecase) List(ctx context.Context) ([]domain.AreaRegion, error) {
	return u.repo.List(ctx)
}

func (u *AreaRegionUsecase) Create(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error) {
	return u.repo.Create(ctx, region)
}

func (u *AreaRegionUsecase) Update(ctx context.Context, region *domain.AreaRegion) (*domain.AreaRegion, error) {
	return u.repo.Update(ctx, region)
}

func (u *AreaRegionUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}
