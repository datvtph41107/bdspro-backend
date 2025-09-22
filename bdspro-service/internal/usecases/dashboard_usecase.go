package usecases

import (
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
	"time"

	_enum "common/domain/enum"
	bdspropb "pb/types/bdspro"
)

type DashboardUsecase struct {
	ProductRepo repo.ProductRepo
	AssetRepo   repo.AssetRepo
	ProjectRepo repo.ProjectRepo
}

func NewDashboardUsecase(
	productRepo repo.ProductRepo,
	assetRepo repo.AssetRepo,
	projectRepo repo.ProjectRepo,
) *DashboardUsecase {
	return &DashboardUsecase{
		ProductRepo: productRepo,
		AssetRepo:   assetRepo,
		ProjectRepo: projectRepo,
	}
}

func (u *DashboardUsecase) GetCashProfit(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*bdspropb.CashProfit, error) {
	return nil, nil
}

func (u *DashboardUsecase) GetCountByOwner(ctx context.Context, ownerOf _enum.EOwnerOf, ownerId uint64) (*dto.CountByOwnerResponse, error) {
	// Đếm sản phẩm
	totalProduct, err := u.ProductRepo.CountByOwner(ctx, ownerOf, ownerId)
	if err != nil {
		return nil, err
	}

	// Đếm tài sản
	totalAsset, err := u.AssetRepo.CountByOwner(ctx, ownerOf, ownerId)
	if err != nil {
		return nil, err
	}

	// Đếm dự án
	totalProject, err := u.ProjectRepo.CountByOwner(ctx, ownerOf, ownerId)
	if err != nil {
		return nil, err
	}

	return &dto.CountByOwnerResponse{
		TotalProduct: totalProduct,
		TotalAsset:   totalAsset,
		TotalProject: totalProject,
	}, nil
}
