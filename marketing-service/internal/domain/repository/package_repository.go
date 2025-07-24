package repository

import (
	"context"
	"marketing/internal/domain/model"
	"marketing/internal/dto"
)

type PackageRepository interface {
	Create(ctx context.Context, pkg *model.Package) (*model.Package, error)
	Update(ctx context.Context, id uint64, pkg *model.Package) (*model.Package, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Package, error)
	Search(ctx context.Context, organizationID uint64, searchDTO dto.PackageSearchDTO) ([]*model.Package, int64, error)
	GetByType(ctx context.Context, packageType model.PackageType) ([]*model.Package, error)
	GetPopularPackages(ctx context.Context, limit int) ([]*model.Package, error)
	GetRecommendedPackages(ctx context.Context, limit int) ([]*model.Package, error)
	GetBestValuePackages(ctx context.Context, limit int) ([]*model.Package, error)
	GetPackagesByProductType(ctx context.Context, productType string) ([]*model.Package, error)
	UpdateStatus(ctx context.Context, id uint64, status model.PackageStatus) error
} 