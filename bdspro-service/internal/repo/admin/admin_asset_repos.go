package admin_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	"context"
)

type AdminAssetRepo interface {
	crud.ICrudRepo[domain.Asset]
	Approve(ctx context.Context, id uint64) error
	Reject(ctx context.Context, id uint64) error
	Archive(ctx context.Context, req *dto.ArchivedRequest) error
	Merge(ctx context.Context, req *dto.MergeAssetRequest) (*domain.Asset, error)
	Split(ctx context.Context, req *dto.SplitAssetDTO) (*domain.Asset, error)
	GetAdminAssets(ctx context.Context, req *dto.AdminAssetSearchRequest) ([]domain.Asset, int64, error)
	GetAssetDetail(ctx context.Context, id uint64) (*domain.Asset, error)
}
