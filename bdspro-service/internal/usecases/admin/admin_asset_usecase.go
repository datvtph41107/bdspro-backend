package admin_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
	"context"
)

type AdminAssetUsecase struct {
	crud.BaseUsecase[domain.Asset, admin_repo.AdminAssetRepo]
}

func NewAdminAssetUsecase(assetRepo admin_repo.AdminAssetRepo) *AdminAssetUsecase {
	return &AdminAssetUsecase{
		crud.BaseUsecase[domain.Asset, admin_repo.AdminAssetRepo]{
			Repo: assetRepo,
		},
	}
}

func (uc *AdminAssetUsecase) Approve(ctx context.Context, id uint64) error {
	return uc.Repo.Approve(ctx, id)
}

func (uc *AdminAssetUsecase) Reject(ctx context.Context, id uint64) error {
	return uc.Repo.Reject(ctx, id)
}

func (uc *AdminAssetUsecase) Archive(ctx context.Context, req *dto.ArchivedRequest) error {
	return uc.Repo.Archive(ctx, req)
}

func (uc *AdminAssetUsecase) Merge(ctx context.Context, req *dto.MergeAssetRequest) (*domain.Asset, error) {
	return uc.Repo.Merge(ctx, req)
}

func (uc *AdminAssetUsecase) Split(ctx context.Context, req *dto.SplitAssetDTO) (*domain.Asset, error) {
	return uc.Repo.Split(ctx, req)
}

func (uc *AdminAssetUsecase) GetAdminAssets(ctx context.Context, req *dto.AdminAssetSearchRequest) ([]domain.Asset, int64, error) {
	return uc.Repo.GetAdminAssets(ctx, req)
}

func (uc *AdminAssetUsecase) GetAssetDetail(ctx context.Context, id uint64) (*domain.Asset, error) {
	return uc.Repo.GetAssetDetail(ctx, id)
}
