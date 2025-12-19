package usecases

import (
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
	"math/rand"
	bdspropb "pb/types/bdspro"
	"time"

	_enum "common/domain/enum"
)

type DashboardUsecase struct {
	ProductRepo repo.ProductRepo
	AssetRepo   repo.AssetRepo
	ProjectRepo repo.ProjectRepo
	PostRepo    repo.PostRepo
}

func NewDashboardUsecase(
	productRepo repo.ProductRepo,
	assetRepo repo.AssetRepo,
	projectRepo repo.ProjectRepo,
	postRepo repo.PostRepo,
) *DashboardUsecase {
	return &DashboardUsecase{
		ProductRepo: productRepo,
		AssetRepo:   assetRepo,
		ProjectRepo: projectRepo,
		PostRepo:    postRepo,
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

	// Đếm post
	totalPost, err := u.PostRepo.CountByOwner(ctx, ownerOf, ownerId)
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
		TotalPost:    totalPost,
		TotalProject: totalProject,
	}, nil
}

func (u *DashboardUsecase) GetCountPostByTime(ctx context.Context, from time.Time, to time.Time) ([]*dto.CountPostByTimeResponse, error) {
	countPostByTime, err := u.PostRepo.CountPostByTime(ctx, from, to)
	if err != nil {
		return nil, err
	}

	countPostByTimeResponse := make([]*dto.CountPostByTimeResponse, len(countPostByTime))
	for i, count := range countPostByTime {
		countPostByTimeResponse[i] = &dto.CountPostByTimeResponse{
			Date:  count.Date,
			Count: count.Count,
		}
	}
	return countPostByTimeResponse, nil
}

func (u *DashboardUsecase) GetCountOfUser(ctx context.Context, profileId uint64) (*dto.CountOfUserResponse, error) {
	// Đếm số tin đã đăng
	totalPosted := u.PostRepo.CountPost(ctx, profileId)

	return &dto.CountOfUserResponse{
		TotalPosted:            uint32(totalPosted),
		TotalInterestedProduct: uint32(20 + rand.Intn(81)), // Random từ 20-100
		TotalSavedProduct:      uint32(20 + rand.Intn(81)), // Random từ 20-100
	}, nil
}
