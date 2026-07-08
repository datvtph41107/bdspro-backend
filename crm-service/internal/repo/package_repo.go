package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type PackageRepo interface {
	Create(ctx context.Context, pkg *domain.Package) (*domain.Package, error)
	Update(ctx context.Context, id uint64, pkg *domain.Package) (*domain.Package, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.Package, error)
	Search(ctx context.Context, organizationID uint64, searchDTO dto.PackageSearchDTO) ([]*domain.Package, int64, error)
	GetByType(ctx context.Context, packageType domain.PackageType) ([]*domain.Package, error)
	GetPopularPackages(ctx context.Context, limit int) ([]*domain.Package, error)
	GetRecommendedPackages(ctx context.Context, limit int) ([]*domain.Package, error)
	GetBestValuePackages(ctx context.Context, limit int) ([]*domain.Package, error)
	GetPackagesByProductType(ctx context.Context, productType string) ([]*domain.Package, error)
	UpdateStatus(ctx context.Context, id uint64, status domain.PackageStatus) error
}