package usecase

import (
	"context"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

// IRegionUsecase defines business operations for region
// mainly simple CRUD and lookup.
type IRegionUsecase interface {
	GetByID(ctx context.Context, id uint64) (*domain.Region, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error)
	List(ctx context.Context) ([]domain.Region, error)
	Create(ctx context.Context, region *domain.Region) (*domain.Region, error)
	Update(ctx context.Context, region *domain.Region) (*domain.Region, error)
	Delete(ctx context.Context, id uint64) error
}

// RegionUsecase is the implementation of IRegionUsecase
type RegionUsecase struct {
	repo _repo.RegionRepository
}

func NewRegionUsecase(repo _repo.RegionRepository) IRegionUsecase {
	return &RegionUsecase{repo: repo}
}

func (u *RegionUsecase) GetByID(ctx context.Context, id uint64) (*domain.Region, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *RegionUsecase) GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error) {
	return u.repo.GetByIDs(ctx, ids)
}

func (u *RegionUsecase) List(ctx context.Context) ([]domain.Region, error) {
	return u.repo.List(ctx)
}

func (u *RegionUsecase) Create(ctx context.Context, region *domain.Region) (*domain.Region, error) {
	return u.repo.Create(ctx, region)
}

func (u *RegionUsecase) Update(ctx context.Context, region *domain.Region) (*domain.Region, error) {
	return u.repo.Update(ctx, region)
}

func (u *RegionUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}
