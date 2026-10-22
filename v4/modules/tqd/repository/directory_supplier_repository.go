package repository

import (
	"context"
	"mdn4/modules/tqd/dto"
	domain "mdn4/modules/tqd/entity"
)

// DirectorySupplierRepository defines the interface for directory supplier repository
type DirectorySupplierRepository interface {
	// CRUD operations
	Create(ctx context.Context, supplier *domain.DirectorySupplier) error
	GetByID(ctx context.Context, id uint64) (*domain.DirectorySupplier, error)
	GetByCode(ctx context.Context, code string) (*domain.DirectorySupplier, error)
	Update(ctx context.Context, supplier *domain.DirectorySupplier) error
	Delete(ctx context.Context, id uint64) error

	// List operations
	List(ctx context.Context, req *dto.DirectorySupplierFilterDTO) ([]domain.DirectorySupplier, int64, error)
}
